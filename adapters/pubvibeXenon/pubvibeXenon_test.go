package pubvibeXenon

import (
	"testing"

	"github.com/prebid/prebid-server/v4/adapters/adapterstest"
	"github.com/prebid/prebid-server/v4/config"
	"github.com/prebid/prebid-server/v4/openrtb_ext"
)

const testsDir = "pubvibeXenontest"
const testsBannerEndpoint = "http://rtb.xenrtb.com/?pid=9d203cb398882b94a15db0aaf0c08470"
const testsVideoEndpoint = "http://rtb.xenrtb.com/?pid=7d7a1a086cfa5105a4bdee1bb73f90ac"

func TestJsonSamples(t *testing.T) {
	bidder, buildErr := Builder(
		openrtb_ext.BidderPubvibeXenon,
		config.Adapter{
			Endpoint:         testsBannerEndpoint,
			ExtraAdapterInfo: `{"video_endpoint":"` + testsVideoEndpoint + `"}`,
		},
		config.Server{ExternalUrl: "http://hosturl.com", GvlID: 1, DataCenter: "2"})
	if buildErr != nil {
		t.Fatalf("Builder returned unexpected error %v", buildErr)
	}

	adapterstest.RunJSONBidderTest(t, testsDir, bidder)
}
