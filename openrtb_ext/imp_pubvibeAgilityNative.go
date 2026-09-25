package openrtb_ext

// ExtImpPubvibeAgilityNative defines the publisher-supplied parameters for the pubvibeAgilityNative bidder.
// All fields are optional. When mediatype is omitted the adapter routes by impression type:
// banner → Banner Web MS endpoint (default), native → Native endpoint.
type ExtImpPubvibeAgilityNative struct {
	// MediaType selects the AgilityNative endpoint to use.
	// Accepted values: "banner_web_ms", "banner_ia", "native"
	MediaType string `json:"mediatype,omitempty"`
}
