package pubvibeISCream

import (
	"testing"

	"github.com/prebid/prebid-server/v4/adapters/adapterstest"
	"github.com/prebid/prebid-server/v4/config"
	"github.com/prebid/prebid-server/v4/openrtb_ext"
)

const testsDir = "pubvibeISCreamtest"
const testsBannerIAEndpoint = "http://srv.agilityadvsrv.com/rtb/1931448329/2138715324"
const testsBannerWebMSEndpoint = "http://srv.agilityadvsrv.com/rtb/1931448329/1272911745"
const testsVideoIAEndpoint = "http://srv.agilityadvsrv.com/rtb/1931448329/3367928980"

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(
		openrtb_ext.BidderPubvibeISCream,
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
