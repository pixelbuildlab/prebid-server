package openrtb_ext

// ExtImpPubvibeAniview defines the publisher-supplied parameters for the pubvibeAniview bidder.
// All fields are optional. When mediatype is omitted the adapter routes by impression type:
// banner → Banner IA endpoint, video → Video IA endpoint.
type ExtImpPubvibeAniview struct {
	// MediaType selects the Aniview endpoint to use.
	// Accepted values: "banner_ia", "banner_web_ms", "video_ia"
	MediaType string `json:"mediatype,omitempty"`
}
