package pubvive

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

// extraInfo holds the video endpoint URL, parsed from config.Adapter.ExtraAdapterInfo.
type extraInfo struct {
	VideoEndpoint string `json:"video_endpoint"`
}

type adapter struct {
	bannerEndpoint string
	videoEndpoint  string
}

// Builder creates a new pubvive adapter.
// config.Adapter.Endpoint        → banner endpoint
// config.Adapter.ExtraAdapterInfo → JSON: {"video_endpoint": "..."}
func Builder(_ openrtb_ext.BidderName, adapterConfig config.Adapter, _ config.Server) (adapters.Bidder, error) {
	info := extraInfo{}
	if adapterConfig.ExtraAdapterInfo != "" {
		if err := json.Unmarshal([]byte(adapterConfig.ExtraAdapterInfo), &info); err != nil {
			return nil, fmt.Errorf("pubvive: failed to parse extra_info: %w", err)
		}
	}

	return &adapter{
		bannerEndpoint: adapterConfig.Endpoint,
		videoEndpoint:  info.VideoEndpoint,
	}, nil
}

// MakeRequests splits the OpenRTB request by media type and routes banner imps
// to the banner endpoint and video imps to the video endpoint.
func (a *adapter) MakeRequests(request *openrtb2.BidRequest, _ *adapters.ExtraRequestInfo) ([]*adapters.RequestData, []error) {
	var bannerImps, videoImps []openrtb2.Imp

	for _, imp := range request.Imp {
		if imp.Banner != nil {
			bannerImps = append(bannerImps, imp)
		} else if imp.Video != nil {
			videoImps = append(videoImps, imp)
		}
	}

	var requests []*adapters.RequestData
	var errs []error

	if len(bannerImps) > 0 {
		reqData, err := a.buildRequest(request, bannerImps, a.bannerEndpoint)
		if err != nil {
			errs = append(errs, err)
		} else {
			requests = append(requests, reqData)
		}
	}

	if len(videoImps) > 0 {
		reqData, err := a.buildRequest(request, videoImps, a.videoEndpoint)
		if err != nil {
			errs = append(errs, err)
		} else {
			requests = append(requests, reqData)
		}
	}

	return requests, errs
}

func (a *adapter) buildRequest(request *openrtb2.BidRequest, imps []openrtb2.Imp, endpoint string) (*adapters.RequestData, error) {
	reqCopy := *request
	reqCopy.Imp = imps

	body, err := json.Marshal(reqCopy)
	if err != nil {
		return nil, fmt.Errorf("marshal bidRequest: %w", err)
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

// MakeBids unpacks the server's response into Bids.
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
			Message: fmt.Sprintf("failed to unmarshal response: %s", err),
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
	case openrtb2.MarkupVideo:
		return openrtb_ext.BidTypeVideo, nil
	default:
		return "", fmt.Errorf("unsupported MType %d", bid.MType)
	}
}
