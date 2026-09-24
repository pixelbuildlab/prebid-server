package pubvibeAniview

import (
	"testing"

	"github.com/prebid/prebid-server/v4/adapters/adapterstest"
	"github.com/prebid/prebid-server/v4/config"
	"github.com/prebid/prebid-server/v4/openrtb_ext"
)

const testsDir = "pubvibeAniviewtest"
const testsBannerIAEndpoint = "https://rtb.aniview.com/sspRTB2?tagid=6aae3b475877fabb4000b596"
const testsBannerWebMSEndpoint = "https://rtb.aniview.com/sspRTB2?tagid=6aae3ba9df8f4d78430013d6"
const testsVideoIAEndpoint = "https://rtb.aniview.com/sspRTB2?tagid=6aae3c03b9da07c777013436"

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(
		openrtb_ext.BidderPubvibeAniview,
		config.Adapter{
			Endpoint:         testsBannerIAEndpoint,
			ExtraAdapterInfo: `{"banner_web_ms_endpoint":"` + testsBannerWebMSEndpoint + `","video_ia_endpoint":"` + testsVideoIAEndpoint + `"}`,
		},
		config.Server{ExternalUrl: "http://hosturl.com", GvlID: 1, DataCenter: "2"})
	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, testsDir, bidder)
}
