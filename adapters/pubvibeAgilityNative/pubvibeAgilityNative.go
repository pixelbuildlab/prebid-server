package pubvibeAgilityNative

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prebid/openrtb/v20/openrtb2"
	"github.com/prebid/prebid-server/v4/adapters"
	"github.com/prebid/prebid-server/v4/config"
	"github.com/prebid/prebid-server/v4/errortypes"
	"github.com/prebid/prebid-server/v4/openrtb_ext"
	"github.com/prebid/prebid-server/v4/util/jsonutil"
)

// extraInfo holds additional endpoint URLs, parsed from config.Adapter.ExtraAdapterInfo.
// config.Adapter.Endpoint           → Banner Web MS endpoint (default)
// extra_info.banner_ia_endpoint     → Banner IA endpoint
// extra_info.native_endpoint        → Native endpoint
type extraInfo struct {
	BannerIAEndpoint string `json:"banner_ia_endpoint"`
	NativeEndpoint   string `json:"native_endpoint"`
}

type adapter struct {
	bannerWebMSEndpoint string
	bannerIAEndpoint    string
	nativeEndpoint      string
}

// Builder creates a new pubvibeAgilityNative adapter.
func Builder(_ openrtb_ext.BidderName, adapterConfig config.Adapter, _ config.Server) (adapters.Bidder, error) {
	info := extraInfo{}
	if adapterConfig.ExtraAdapterInfo != "" {
		if err := json.Unmarshal([]byte(adapterConfig.ExtraAdapterInfo), &info); err != nil {
			return nil, fmt.Errorf("pubvibeAgilityNative: failed to parse extra_info: %w", err)
		}
	}

	return &adapter{
		bannerWebMSEndpoint: adapterConfig.Endpoint,
		bannerIAEndpoint:    info.BannerIAEndpoint,
		nativeEndpoint:      info.NativeEndpoint,
	}, nil
}

// MakeRequests splits the OpenRTB request by media type and routes each impression
// to the correct AgilityNative endpoint based on the publisher-supplied mediatype param.
// If mediatype is absent, the adapter defaults to:
//   - banner imps  → Banner Web MS endpoint
//   - native imps  → Native endpoint
func (a *adapter) MakeRequests(request *openrtb2.BidRequest, _ *adapters.ExtraRequestInfo) ([]*adapters.RequestData, []error) {
	var bannerWebMSImps, bannerIAImps, nativeImps []openrtb2.Imp
	var errs []error

	for _, imp := range request.Imp {
		mediaType, err := parseMediaType(imp)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		switch mediaType {
		case "banner_ia":
			if imp.Banner != nil {
				bannerIAImps = append(bannerIAImps, imp)
			}
		case "native":
			if imp.Native != nil {
				nativeImps = append(nativeImps, imp)
			}
		default:
			// "banner_web_ms" or no mediatype — route by imp type, default to banner_web_ms
			if imp.Banner != nil {
				bannerWebMSImps = append(bannerWebMSImps, imp)
			} else if imp.Native != nil {
				nativeImps = append(nativeImps, imp)
			}
		}
	}

	var requests []*adapters.RequestData

	if len(bannerWebMSImps) > 0 {
		reqData, err := buildRequest(request, bannerWebMSImps, a.bannerWebMSEndpoint)
		if err != nil {
			errs = append(errs, err)
		} else {
			requests = append(requests, reqData)
		}
	}

	if len(bannerIAImps) > 0 {
		reqData, err := buildRequest(request, bannerIAImps, a.bannerIAEndpoint)
		if err != nil {
			errs = append(errs, err)
		} else {
			requests = append(requests, reqData)
		}
	}

	if len(nativeImps) > 0 {
		reqData, err := buildRequest(request, nativeImps, a.nativeEndpoint)
		if err != nil {
			errs = append(errs, err)
		} else {
			requests = append(requests, reqData)
		}
	}

	return requests, errs
}

// parseMediaType extracts the mediatype value from imp.ext.bidder.
// Returns empty string (not an error) when the field is absent — callers treat that as default routing.
func parseMediaType(imp openrtb2.Imp) (string, error) {
	var bidderExt adapters.ExtImpBidder
	if err := json.Unmarshal(imp.Ext, &bidderExt); err != nil {
		return "", fmt.Errorf("imp %s: failed to unmarshal ext: %w", imp.ID, err)
	}

	if bidderExt.Bidder == nil {
		return "", nil
	}

	var pubvibeExt openrtb_ext.ExtImpPubvibeAgilityNative
	if err := json.Unmarshal(bidderExt.Bidder, &pubvibeExt); err != nil {
		return "", fmt.Errorf("imp %s: failed to unmarshal bidder ext: %w", imp.ID, err)
	}

	return pubvibeExt.MediaType, nil
}

func buildRequest(request *openrtb2.BidRequest, imps []openrtb2.Imp, endpoint string) (*adapters.RequestData, error) {
	reqCopy := *request
	reqCopy.Imp = imps

	body, err := json.Marshal(reqCopy)
	if err != nil {
		return nil, fmt.Errorf("pubvibeAgilityNative: marshal bidRequest: %w", err)
	}

	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	headers.Add("Accept", "application/json")

	return &adapters.RequestData{
		Method:  http.MethodPost,
		Uri:     endpoint,
		Body:    body,
		Headers: headers,
		ImpIDs:  openrtb_ext.GetImpIDs(reqCopy.Imp),
	}, nil
}

// MakeBids unpacks the AgilityNative server response into Prebid bid objects.
func (a *adapter) MakeBids(request *openrtb2.BidRequest, _ *adapters.RequestData, responseData *adapters.ResponseData) (*adapters.BidderResponse, []error) {
	if adapters.IsResponseStatusCodeNoContent(responseData) {
		return nil, nil
	}

	if err := adapters.CheckResponseStatusCodeForErrors(responseData); err != nil {
		return nil, []error{err}
	}

	var response openrtb2.BidResponse
	if err := jsonutil.Unmarshal(responseData.Body, &response); err != nil {
		return nil, []error{&errortypes.BadServerResponse{
			Message: fmt.Sprintf("pubvibeAgilityNative: failed to unmarshal response: %s", err),
		}}
	}

	bidResponse := adapters.NewBidderResponseWithBidsCapacity(len(request.Imp))
	if response.Cur != "" {
		bidResponse.Currency = response.Cur
	}

	var errs []error

	for _, seatBid := range response.SeatBid {
		for i, bid := range seatBid.Bid {
			bidType, err := getMediaTypeForBid(bid)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			bidResponse.Bids = append(bidResponse.Bids, &adapters.TypedBid{
				Bid:     &seatBid.Bid[i],
				BidType: bidType,
			})
		}
	}

	return bidResponse, errs
}

func getMediaTypeForBid(bid openrtb2.Bid) (openrtb_ext.BidType, error) {
	switch bid.MType {
	case openrtb2.MarkupBanner:
		return openrtb_ext.BidTypeBanner, nil
	case openrtb2.MarkupNative:
		return openrtb_ext.BidTypeNative, nil
	default:
		return "", fmt.Errorf("pubvibeAgilityNative: unsupported MType %d for bid %s", bid.MType, bid.ID)
	}
}
