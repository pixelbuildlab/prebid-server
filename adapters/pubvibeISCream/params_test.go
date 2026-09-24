package pubvibeISCream

import (
	"encoding/json"
	"testing"

	"github.com/prebid/prebid-server/v4/openrtb_ext"
)

const schemaDirectory = "../../static/bidder-params"

// TestValidParams ensures the pubvibeISCream schema accepts valid params.
func TestValidParams(t *testing.T) {
	validator, err := openrtb_ext.NewBidderParamsValidator(schemaDirectory)
	if err != nil {
		t.Fatalf("Failed to fetch validator: %v", err)
	}

	for _, p := range validParams {
		if err := validator.Validate(openrtb_ext.BidderPubvibeISCream, json.RawMessage(p)); err != nil {
			t.Errorf("Schema should allow valid params: %s\n Error: %v", p, err)
		}
	}
}

// TestInvalidParams ensures the pubvibeISCream schema rejects invalid params.
func TestInvalidParams(t *testing.T) {
	validator, err := openrtb_ext.NewBidderParamsValidator(schemaDirectory)
	if err != nil {
		t.Fatalf("Failed to fetch validator: %v", err)
	}

	for _, p := range invalidParams {
		if err := validator.Validate(openrtb_ext.BidderPubvibeISCream, json.RawMessage(p)); err == nil {
			t.Errorf("Schema should reject invalid params: %s", p)
		}
	}
}

// All mediatype values are valid; omitting mediatype is also valid (defaults to banner_ia / video_ia).
var validParams = []string{
	`{}`,
	`{"mediatype": "banner_ia"}`,
	`{"mediatype": "banner_web_ms"}`,
	`{"mediatype": "video_ia"}`,
}

var invalidParams = []string{
	`null`,
	`true`,
	`5`,
	`"invalid"`,
	`[]`,
	`{"mediatype": "unknown_type"}`,
}
