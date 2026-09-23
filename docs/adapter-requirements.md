# OpenRTB Endpoint → Prebid Server Adapter: Requirements & Reference

This document captures everything required to build a Prebid Server (Go) bid adapter
from an OpenRTB-compatible endpoint. Use this as a template when integrating any new
RTB endpoint in the future.

---

## Overview

A Prebid Server adapter is the bridge between an incoming OpenRTB 2.x auction request
and an external RTB endpoint. It has two responsibilities:

1. **MakeRequests** — transform the Prebid OpenRTB request into one or more HTTP
   requests for the endpoint
2. **MakeBids** — parse the endpoint's HTTP response back into Prebid bid objects

---

## What You Need From the RTB Endpoint Provider

Before writing an adapter, collect the following from the endpoint provider:

| # | Information Needed | Example (pubvibeXenon) |
|---|-------------------|-------------------|
| 1 | Endpoint URL(s) | `http://rtb.xenrtb.com/?pid=...` |
| 2 | HTTP method | `POST` |
| 3 | Request format | OpenRTB 2.x JSON |
| 4 | Response format | OpenRTB 2.x JSON |
| 5 | Authentication method | `pid` query parameter in URL |
| 6 | Supported media types | banner, video |
| 7 | Separate endpoints per media type? | Yes (different `pid` per type) |
| 8 | Required request headers | `Content-Type: application/json` |
| 9 | Required impression fields | none (standard OpenRTB) |
| 10 | Required bid response fields | `seatbid[].bid[].mtype` (to detect banner vs video) |
| 11 | No-bid response | HTTP 204 or empty body |
| 12 | Error response codes | Standard HTTP 4xx/5xx |
| 13 | Publisher/account identifier | `pid` in URL (per media type) |
| 14 | GVL vendor ID (for GDPR) | none required for this adapter |

---

## Files Required for Every Adapter

### 1. `adapters/{name}/{name}.go` — Core adapter logic

Must implement the `adapters.Bidder` interface:

```go
type Bidder interface {
    MakeRequests(*openrtb2.BidRequest, *ExtraRequestInfo) ([]*RequestData, []error)
    MakeBids(*openrtb2.BidRequest, *RequestData, *ResponseData) (*BidderResponse, []error)
}
```

Must export a `Builder` function with this exact signature:

```go
func Builder(_ openrtb_ext.BidderName, config config.Adapter, _ config.Server) (adapters.Bidder, error)
```

**MakeRequests checklist:**
- [ ] Read endpoint URL from `config.Adapter.Endpoint`
- [ ] Read extra config (e.g. second endpoint) from `config.Adapter.ExtraAdapterInfo`
- [ ] Copy the request before modifying (`reqCopy := *request`)
- [ ] Set `ImpIDs` on `RequestData` using `openrtb_ext.GetImpIDs()`
- [ ] Set `Content-Type: application/json;charset=utf-8` header
- [ ] Set `Accept: application/json` header
- [ ] Marshal request body with `json.Marshal()`

**MakeBids checklist:**
- [ ] Handle HTTP 204 with `adapters.IsResponseStatusCodeNoContent()`
- [ ] Handle HTTP errors with `adapters.CheckResponseStatusCodeForErrors()`
- [ ] Unmarshal response with `jsonutil.Unmarshal()`
- [ ] Set `bidResponse.Currency` from `response.Cur`
- [ ] Detect bid type from `bid.MType` (1=banner, 2=video, 4=native)
- [ ] Use `&seatBid.Bid[i]` (pointer to slice element, not loop variable)

---

### 2. `static/bidder-info/{name}.yaml` — Bidder metadata

```yaml
endpoint: "https://endpoint.example.com/bid"
# optional: second endpoint or extra config
extra_info: '{"video_endpoint":"https://endpoint.example.com/video"}'
maintainer:
  email: your@email.com
gvlVendorID: 123          # only if GDPR consent required
openrtb:
  version: 2.6
capabilities:
  site:
    mediaTypes:
      - banner
      - video
  app:                    # include only if endpoint supports app traffic
    mediaTypes:
      - banner
```

---

### 3. `static/bidder-params/{name}.json` — JSON schema for imp ext params

If the endpoint needs publisher-supplied params (e.g. placement ID):

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "title": "MyAdapter Params",
  "type": "object",
  "properties": {
    "placement_id": {
      "type": "string",
      "description": "Publisher placement ID",
      "minLength": 1
    }
  },
  "required": ["placement_id"]
}
```

If no params needed (like pubvibeXenon — routing handled server-side):

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "title": "MyAdapter Params",
  "type": "object",
  "properties": {}
}
```

---

### 4. `openrtb_ext/imp_{name}.go` — Imp ext type

Maps to the JSON schema above:

```go
package openrtb_ext

type ExtImpMyAdapter struct {
    PlacementID string `json:"placement_id"`
}
```

If no params:
```go
package openrtb_ext

type ExtImpMyAdapter struct{}
```

---

### 5. `openrtb_ext/bidders.go` — Bidder registration (2 places)

**In `coreBidderNames` slice** (alphabetical order):
```go
BidderMyAdapter,
```

**In const block** (alphabetical order):
```go
BidderMyAdapter BidderName = "myadapter"
```

---

### 6. `exchange/adapter_builders.go` — Builder registration (2 places)

**In imports** (alphabetical order):
```go
"github.com/prebid/prebid-server/v4/adapters/myadapter"
```

**In `newAdapterBuilders()` map** (alphabetical order):
```go
openrtb_ext.BidderMyAdapter: myadapter.Builder,
```

---

### 7. Test files

| File | Purpose |
|------|---------|
| `adapters/{name}/{name}_test.go` | Runs JSON test cases |
| `adapters/{name}/params_test.go` | Validates bidder params schema |
| `adapters/{name}/{name}test/exemplary/simple-banner.json` | Banner bid test |
| `adapters/{name}/{name}test/exemplary/simple-video.json` | Video bid test |
| `adapters/{name}/{name}test/exemplary/banner-and-video.json` | Mixed type test |

**Test JSON structure:**
```json
{
  "mockBidRequest": { },        // OpenRTB BidRequest sent by publisher
  "httpCalls": [
    {
      "expectedRequest": {
        "uri": "https://endpoint.example.com/bid",
        "body": { },            // what adapter should send to endpoint
        "impIDs": ["imp-1"]
      },
      "mockResponse": {
        "status": 200,
        "body": { }             // fake endpoint response
      }
    }
  ],
  "expectedBidResponses": [
    {
      "currency": "USD",
      "bids": [
        {
          "bid": { },           // expected bid object
          "type": "banner"      // or "video"
        }
      ]
    }
  ]
}
```

---

## Endpoint Routing Patterns

### Single endpoint (most common)
```
All imps → one URL
config.Adapter.Endpoint = "https://endpoint.example.com/bid"
```

### Separate endpoints per media type (pubvibeXenon pattern)
```
Banner imps → bannerEndpoint
Video imps  → videoEndpoint

config.Adapter.Endpoint        = banner URL
config.Adapter.ExtraAdapterInfo = '{"video_endpoint":"..."}'
```

### Dynamic endpoint (param in URL)
```
URL built from publisher-supplied param in imp ext
e.g. "https://endpoint.example.com/?pub_id=" + impExt.PublisherID
```

---

## Bid Type Detection

Prefer `bid.MType` (OpenRTB 2.6+):

```go
switch bid.MType {
case openrtb2.MarkupBanner: return openrtb_ext.BidTypeBanner
case openrtb2.MarkupVideo:  return openrtb_ext.BidTypeVideo
case openrtb2.MarkupNative: return openrtb_ext.BidTypeNative
}
```

Fallback for older endpoints (pre-2.6) — infer from imp:

```go
func getMediaTypeForImp(impID string, imps []openrtb2.Imp) (openrtb_ext.BidType, error) {
    for _, imp := range imps {
        if imp.ID == impID {
            if imp.Banner != nil { return openrtb_ext.BidTypeBanner, nil }
            if imp.Video != nil  { return openrtb_ext.BidTypeVideo, nil }
            if imp.Native != nil { return openrtb_ext.BidTypeNative, nil }
        }
    }
    return "", fmt.Errorf("unknown imp id: %s", impID)
}
```

---

## Common Mistakes to Avoid

| Mistake | Correct approach |
|---------|-----------------|
| `bid := seatBid.Bid[i]; &bid` | Use `&seatBid.Bid[i]` directly |
| Modifying `request` directly | Always copy: `reqCopy := *request` |
| Using `../../../static/bidder-params` in params_test.go | Use `../../static/bidder-params` |
| `config.ExtraInfo` | Correct field is `config.ExtraAdapterInfo` |
| `openrtb_ext.NewBidderParamsValidator(dir, "")` | Only one argument: `NewBidderParamsValidator(dir)` |
| Forgetting `ImpIDs` in `RequestData` | Always set via `openrtb_ext.GetImpIDs()` |
| Binary compiled on macOS (darwin/arm64) | Cross-compile: `GOOS=linux GOARCH=amd64 go build` |

---

## Quick Checklist for a New Adapter

```
[ ] Collected endpoint URL(s), auth method, media types from provider
[ ] Created adapters/{name}/{name}.go with Builder, MakeRequests, MakeBids
[ ] Created static/bidder-info/{name}.yaml
[ ] Created static/bidder-params/{name}.json
[ ] Created openrtb_ext/imp_{name}.go
[ ] Added BidderName to coreBidderNames slice in openrtb_ext/bidders.go
[ ] Added BidderName const to openrtb_ext/bidders.go
[ ] Added import to exchange/adapter_builders.go
[ ] Added Builder to newAdapterBuilders() map in exchange/adapter_builders.go
[ ] Created test data JSON files in {name}test/exemplary/
[ ] Created {name}_test.go and params_test.go
[ ] go build ./... passes with no errors
[ ] go test ./adapters/{name}/... passes
```

---

## Reference: pubvibeXenon Adapter (xenrtb.com)

This adapter was the first one built using this pattern. Key decisions made:

- **Two endpoints** — xenrtb uses different `pid` values for banner vs video, so
  the adapter splits imps by media type and makes separate requests
- **No publisher params** — routing is handled entirely by the server config,
  publishers just send `"ext": {"bidder": {"pubvibeXenon": {}}}`
- **mtype-based bid detection** — endpoint returns OpenRTB 2.6 `mtype` field,
  so bid type is read directly from the response
- **extra_info pattern** — video endpoint stored as JSON in `extra_info` field,
  following the same pattern as `beachfront` adapter

Source: `adapters/pubvibeXenon/pubvibeXenon.go`
