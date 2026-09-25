package pubvibeAgilityNative

import (
	"testing"

	"github.com/prebid/prebid-server/v4/adapters/adapterstest"
	"github.com/prebid/prebid-server/v4/config"
	"github.com/prebid/prebid-server/v4/openrtb_ext"
)

const testsDir = "pubvibeAgilityNativetest"
const testsBannerWebMSEndpoint = "http://rtb-useast.agilitydigitalmedia.com/rtb?zone=413281"
const testsBannerIAEndpoint = "http://rtb-useast.agilitydigitalmedia.com/rtb?zone=413283"
const testsNativeEndpoint = "http://rtb-useast.agilitydigitalmedia.com/rtb?zone=413284"

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(
		openrtb_ext.BidderPubvibeAgilityNative,
		config.Adapter{
			Endpoint:         testsBannerWebMSEndpoint,
			ExtraAdapterInfo: `{"banner_ia_endpoint":"` + testsBannerIAEndpoint + `","native_endpoint":"` + testsNativeEndpoint + `"}`,
		},
		config.Server{ExternalUrl: "http://hosturl.com", GvlID: 1, DataCenter: "2"})
	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, testsDir, bidder)
}
