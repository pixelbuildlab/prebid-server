package pubvive

import (
	"encoding/json"
	"testing"

	"github.com/prebid/prebid-server/v4/openrtb_ext"
)

const schemaDirectory = "../../static/bidder-params"

// TestValidParams ensures the pubvive schema accepts valid params.
// The pubvive adapter has no required params so an empty object is valid.
func TestValidParams(t *testing.T) {
	validator, err := openrtb_ext.NewBidderParamsValidator(schemaDirectory)
	if err != nil {
		t.Fatalf("Failed to fetch validator: %v", err)
	}

	for _, p := range validParams {
		if err := validator.Validate(openrtb_ext.BidderPubvive, json.RawMessage(p)); err != nil {
			t.Errorf("Schema should allow valid params: %s\n Error: %v", p, err)
		}
	}
}

// TestInvalidParams ensures the pubvive schema rejects invalid params.
func TestInvalidParams(t *testing.T) {
	validator, err := openrtb_ext.NewBidderParamsValidator(schemaDirectory)
	if err != nil {
		t.Fatalf("Failed to fetch validator: %v", err)
	}

	for _, p := range invalidParams {
		if err := validator.Validate(openrtb_ext.BidderPubvive, json.RawMessage(p)); err == nil {
			t.Errorf("Schema should reject invalid params: %s", p)
		}
	}
}

var validParams = []string{
	`{}`,
}

var invalidParams = []string{
	`null`,
	`true`,
	`5`,
	`"invalid"`,
	`[]`,
}
