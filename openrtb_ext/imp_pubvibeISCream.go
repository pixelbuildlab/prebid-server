package openrtb_ext

// ExtImpPubvibeISCream defines the publisher-supplied parameters for the pubvibeISCream bidder.
// All fields are optional. When mediatype is omitted the adapter routes by impression type:
// banner → Banner IA endpoint, video → Video IA endpoint.
type ExtImpPubvibeISCream struct {
	// MediaType selects the ISCream endpoint to use.
	// Accepted values: "banner_ia", "banner_web_ms", "video_ia"
	MediaType string `json:"mediatype,omitempty"`
}
