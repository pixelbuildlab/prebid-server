# Pubvive Bid Adapter Integration

## Overview

This document describes the integration of the **pubvive** bid adapter into the Prebid Server (Go). The adapter connects to the [xenrtb.com](http://rtb.xenrtb.com) RTB endpoints for web banner and video ad serving.

---

## Endpoints

| Format | Endpoint |
|--------|----------|
| Banner | `http://rtb.xenrtb.com/?pid=9d203cb398882b94a15db0aaf0c08470` |
| Video  | `http://rtb.xenrtb.com/?pid=7d7a1a086cfa5105a4bdee1bb73f90ac` |

There is also a financial/reporting API (`http://api.xenrtb.com/ssp/financial/`) which is **not** integrated — it is out of scope for the current phase (ad serving only).

---

## Files Created / Modified

### New Files

| File | Description |
|------|-------------|
| `adapters/pubvive/pubvive.go` | Core adapter logic |
| `adapters/pubvive/pubvive_test.go` | JSON test runner |
| `adapters/pubvive/params_test.go` | Bidder params schema validation tests |
| `adapters/pubvive/pubvivetest/exemplary/simple-banner.json` | Test case: banner only |
| `adapters/pubvive/pubvivetest/exemplary/simple-video.json` | Test case: video only |
| `adapters/pubvive/pubvivetest/exemplary/banner-and-video.json` | Test case: both types in one request |
| `static/bidder-info/pubvive.yaml` | Bidder metadata and endpoint config |
| `static/bidder-params/pubvive.json` | JSON schema for bidder params (no required params) |
| `openrtb_ext/imp_pubvive.go` | Imp-level ext type definition |

### Modified Files

| File | Change |
|------|--------|
| `openrtb_ext/bidders.go` | Added `BidderPubvive` to `coreBidderNames` slice and const block |
| `exchange/adapter_builders.go` | Added `pubvive` import and `BidderPubvive: pubvive.Builder` mapping |

---

## Adapter Design

### Routing Logic

The adapter splits incoming impressions by media type and sends each group to the appropriate endpoint:

- Impressions with `imp.banner` → banner endpoint (`pid=9d203...`)
- Impressions with `imp.video` → video endpoint (`pid=7d7a1...`)

If a single bid request contains both banner and video impressions, **two separate HTTP calls** are made — one per endpoint. Each returns its own `BidderResponse` which Prebid Server merges into the final auction response.

### Configuration

The endpoints are stored in `static/bidder-info/pubvive.yaml`:

```yaml
endpoint: "http://rtb.xenrtb.com/?pid=9d203cb398882b94a15db0aaf0c08470"
extra_info: '{"video_endpoint":"http://rtb.xenrtb.com/?pid=7d7a1a086cfa5105a4bdee1bb73f90ac"}'
```

- `endpoint` → banner endpoint (standard `config.Adapter.Endpoint` field)
- `extra_info` → JSON string parsed by the adapter to extract `video_endpoint` (stored in `config.Adapter.ExtraAdapterInfo`)

### Publisher-Side Parameters

The adapter requires **no publisher-supplied parameters** in the impression ext. Publishers simply include:

```json
"ext": { "prebid": { "bidder": { "pubvive": {} } } }
```

---

## Running the Server

```bash
cd prebid-server
PBS_GDPR_DEFAULT_VALUE=0 go run .
```

Server starts on port **8000** by default. The two GDPR vendor list warnings at startup are non-fatal and only affect cookie syncing, not ad serving.

---

## Example Bid Request

```bash
curl -X POST http://localhost:8000/openrtb2/auction \
  -H "Content-Type: application/json" \
  -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" \
  -H "X-Forwarded-For: 1.2.3.4" \
  -d '{
    "id": "test-request",
    "imp": [{
      "id": "imp-1",
      "banner": { "format": [{"w": 300, "h": 250}] },
      "bidfloor": 0.01,
      "bidfloorcur": "USD",
      "ext": { "prebid": { "bidder": { "pubvive": {} } } }
    }],
    "site": {
      "page": "https://yoursite.com/article",
      "domain": "yoursite.com",
      "publisher": { "id": "your-publisher-id" }
    },
    "device": {
      "ua": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
      "ip": "1.2.3.4",
      "language": "en"
    }
  }'
```

For video:

```bash
curl -X POST http://localhost:8000/openrtb2/auction \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test-request",
    "imp": [{
      "id": "imp-1",
      "video": {
        "mimes": ["video/mp4"],
        "minduration": 5,
        "maxduration": 30,
        "protocols": [1, 2, 5],
        "w": 640,
        "h": 480
      },
      "ext": { "prebid": { "bidder": { "pubvive": {} } } }
    }],
    "site": {
      "page": "https://yoursite.com/article",
      "domain": "yoursite.com",
      "publisher": { "id": "your-publisher-id" }
    }
  }'
```

---

## Test Results

All adapter tests pass:

```
=== RUN   TestValidParams   --- PASS
=== RUN   TestInvalidParams --- PASS
=== RUN   TestJsonSamples   --- PASS
ok   github.com/prebid/prebid-server/v4/adapters/pubvive
```

Full project build: clean (`go build ./...` — no errors).

Config tests: `TestFullConfig` — PASS.

---

## Current Status

| Task | Status |
|------|--------|
| Adapter code | ✅ Done |
| Banner routing | ✅ Done |
| Video routing | ✅ Done |
| Bidder registration | ✅ Done |
| Tests (unit + JSON) | ✅ Done |
| Server running locally | ✅ Done |
| Live bid request verified | ✅ Done (204ms response from xenrtb.com) |
| Analytics / reporting API | ⏸ Out of scope (current phase) |
| Production deployment | ⏸ Not started |

---

## Next Steps (Optional)

1. **Production config** — create a `prebid_server_config.yaml` with production settings (host, port, cache, rate limits, etc.)
2. **Docker deployment** — build and deploy using the existing `Dockerfile`
3. **Analytics module** — if financial reporting from `api.xenrtb.com` is needed in a future phase
4. **Usersync** — add a usersync endpoint for xenrtb if cookie-based audience targeting is needed

---

## Site Integration (Prebid.js + Deployment)

### Architecture Overview

```
Your Website (browser)
    │
    │  Prebid.js (client-side JS)
    │
    ▼
Prebid Server (Go) ← this repo, running on your server
    │
    ▼
xenrtb.com RTB endpoints (pubvive adapter)
```

ads.txt is already done. The remaining steps are below.

---

### Step 1 — Build Prebid.js

Prebid.js is a separate client-side library. Build a custom bundle at:

**https://docs.prebid.org/download.html**

Make sure to include:
- **Prebid Server Bid Adapter** (`prebidServer`) — connects your page to this Go server
- Any other client-side adapters you want alongside it

Download the generated `prebid.js` file and host it on your site (or use a CDN).

---

### Step 2 — Add Prebid.js to Your Page

```html
<!-- Load Prebid.js -->
<script src="/path/to/prebid.js"></script>

<script>
  var pbjs = pbjs || {};
  pbjs.que = pbjs.que || [];

  pbjs.que.push(function() {

    // Point Prebid.js at your Prebid Server
    pbjs.setConfig({
      s2sConfig: {
        accountId: '1',
        bidders: ['pubvive'],
        defaultVendor: 'appnexus',
        timeout: 1000,
        endpoint: {
          p1Consent: 'https://YOUR_SERVER_DOMAIN/openrtb2/auction',
          noConsent:  'https://YOUR_SERVER_DOMAIN/openrtb2/auction'
        }
      }
    });

    // Define your ad units
    pbjs.addAdUnits([
      {
        code: 'banner-300x250',           // must match div id below
        mediaTypes: {
          banner: { sizes: [[300, 250], [728, 90]] }
        },
        bids: [{
          bidder: 'pubvive',
          params: {}
        }]
      },
      {
        code: 'video-unit',
        mediaTypes: {
          video: {
            playerSize: [640, 480],
            context: 'instream',
            mimes: ['video/mp4'],
            protocols: [1, 2, 5],
            minduration: 5,
            maxduration: 30
          }
        },
        bids: [{
          bidder: 'pubvive',
          params: {}
        }]
      }
    ]);

    // Request bids
    pbjs.requestBids({
      bidsBackHandler: function(bids) {
        // If using Google Ad Manager, set targeting and refresh slots:
        // googletag.cmd.push(function() {
        //   pbjs.setTargetingForGPTAsync();
        //   googletag.pubads().refresh();
        // });
        console.log('Bids received:', bids);
      }
    });

  });
</script>

<!-- Banner ad slot -->
<div id="banner-300x250">
  <script>
    googletag.cmd.push(function() {
      googletag.display('banner-300x250');
    });
  </script>
</div>
```

Replace `YOUR_SERVER_DOMAIN` with the public domain/IP of your Prebid Server.

---

### Step 3 — Deploy Prebid Server Publicly

The server currently runs only on `localhost`. To serve real traffic:

**Option A — Direct (simplest)**
```bash
# On your server
PBS_GDPR_DEFAULT_VALUE=0 ./prebid-server
# Runs on port 8000
```

**Option B — Docker**
```bash
docker build --platform linux/amd64 -t prebid-server .
docker run -e PBS_GDPR_DEFAULT_VALUE=0 -p 8000:8000 prebid-server
```

**Option C — Behind nginx with HTTPS (recommended for production)**

Nginx config:
```nginx
server {
    listen 443 ssl;
    server_name prebid.yoursite.com;

    ssl_certificate     /etc/letsencrypt/live/prebid.yoursite.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/prebid.yoursite.com/privkey.pem;

    location / {
        proxy_pass         http://127.0.0.1:8000;
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
        proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

Get a free SSL cert with:
```bash
certbot --nginx -d prebid.yoursite.com
```

---

### Step 4 — Connect to Google Ad Manager (optional)

If you use GAM (DFP) to serve ads:

1. Set up **Prebid price bucket line items** in GAM — use the [Prebid Line Item Manager](https://docs.prebid.org/adops/step-by-step.html) or do it manually
2. Add the GPT integration to your page alongside Prebid.js:

```html
<!-- Google Publisher Tag -->
<script async src="https://securepubads.g.doubleclick.net/tag/js/gpt.js"></script>
<script>
  window.googletag = window.googletag || { cmd: [] };
  googletag.cmd.push(function() {
    googletag.defineSlot('/YOUR_NETWORK_ID/banner', [300, 250], 'banner-300x250')
      .addService(googletag.pubads());
    googletag.pubads().enableSingleRequest();
    googletag.enableServices();
  });
</script>
```

Then in `bidsBackHandler`:
```javascript
bidsBackHandler: function() {
  googletag.cmd.push(function() {
    pbjs.setTargetingForGPTAsync();
    googletag.pubads().refresh();
  });
}
```

---

### Checklist

| Step | Status |
|------|--------|
| ads.txt | ✅ Done |
| Prebid Server running locally | ✅ Done |
| pubvive adapter integrated | ✅ Done |
| Build Prebid.js bundle | ⏸ Pending |
| Add Prebid.js to site pages | ⏸ Pending |
| Deploy Prebid Server to public host | ⏸ Pending |
| Add HTTPS (SSL) | ⏸ Pending |
| Connect to ad server (GAM) | ⏸ Optional |
