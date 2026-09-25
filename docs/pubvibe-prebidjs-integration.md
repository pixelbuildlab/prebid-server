# Pubvibe Adapters — Prebid.js Integration Guide

**For:** Prebid.js developer  
**Server:** Prebid Server (Go) running your custom adapters  
**Mode:** Server-to-Server (s2sConfig) only — all bidding goes through your Prebid Server  

---

## Adapters Available

| Adapter Name | Provider | Banner Web MS | Banner IA | Native | Video IA |
|---|---|---|---|---|---|
| `pubvibeXenon` | xenrtb.com | ❌ | ✅ (auto, no param) | ❌ | ✅ (auto, no param) |
| `pubvibeAniview` | aniview.com | ✅ | ✅ | ❌ | ✅ |
| `pubvibeISCream` | agilityadvsrv.com | ✅ | ✅ | ❌ | ✅ |
| `pubvibeAgilityNative` | agilitydigitalmedia.com | ✅ (default) | ✅ | ✅ | ❌ |

> **Important for developers:** Always specify `mediatype` explicitly in your bid params for `pubvibeAniview`, `pubvibeISCream`, and `pubvibeAgilityNative`. For `pubvibeXenon`, never pass `mediatype` — use `params: {}` only.

---

## Mediatype Values Per Adapter

> **Param key spelling:** The param is `mediatype` (all lowercase). Do NOT write `mediaType` (camelCase) — it will be ignored by the server and the default routing will apply.

| Adapter | Banner Web MS | Banner IA | Native | Video IA | Default when `params: {}` |
|---|---|---|---|---|---|
| `pubvibeXenon` | ❌ not supported | ❌ no param needed — banner auto-routes | ❌ | ❌ no param needed — video auto-routes | banner imp → banner endpoint, video imp → video endpoint |
| `pubvibeAniview` | `mediatype: 'banner_web_ms'` | `mediatype: 'banner_ia'` | ❌ | `mediatype: 'video_ia'` | banner imp → `banner_ia`, video imp → `video_ia` |
| `pubvibeISCream` | `mediatype: 'banner_web_ms'` | `mediatype: 'banner_ia'` | ❌ | `mediatype: 'video_ia'` | banner imp → `banner_ia`, video imp → `video_ia` |
| `pubvibeAgilityNative` | `mediatype: 'banner_web_ms'` | `mediatype: 'banner_ia'` | `mediatype: 'native'` | ❌ not supported | banner imp → `banner_web_ms`, native imp → `native` |

**Rule of thumb:**
- `pubvibeXenon` → always use `params: {}` (no mediatype needed or accepted)
- All others → always specify `mediatype` explicitly to avoid ambiguity

---

## 1. Prebid Server s2sConfig (Required — set once)

```javascript
pbjs.setConfig({
  s2sConfig: {
    accountId: '1',
    bidders: ['pubvibeXenon', 'pubvibeAniview', 'pubvibeISCream', 'pubvibeAgilityNative'],
    defaultVendor: 'appnexus',
    timeout: 1500,
    endpoint: {
      p1Consent: 'https://YOUR_PREBID_SERVER_DOMAIN/openrtb2/auction',
      noConsent:  'https://YOUR_PREBID_SERVER_DOMAIN/openrtb2/auction'
    }
  }
});
```

> Replace `YOUR_PREBID_SERVER_DOMAIN` with your actual server domain/IP.

---

## 2. Banner Web MS Ad Units

Banner Web MS is available on `pubvibeAniview`, `pubvibeISCream`, and `pubvibeAgilityNative`.
**Always pass `mediatype: 'banner_web_ms'` explicitly.**

### pubvibeAniview — Banner Web MS

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-web-ms',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [{
    bidder: 'pubvibeAniview',
    params: { mediatype: 'banner_web_ms' }
  }]
}]);
```

### pubvibeISCream — Banner Web MS

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-web-ms',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [{
    bidder: 'pubvibeISCream',
    params: { mediatype: 'banner_web_ms' }
  }]
}]);
```

### pubvibeAgilityNative — Banner Web MS

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-web-ms',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [{
    bidder: 'pubvibeAgilityNative',
    params: { mediatype: 'banner_web_ms' }
  }]
}]);
```

### All Three Competing on Banner Web MS

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-web-ms',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [
    { bidder: 'pubvibeAniview',       params: { mediatype: 'banner_web_ms' } },
    { bidder: 'pubvibeISCream',       params: { mediatype: 'banner_web_ms' } },
    { bidder: 'pubvibeAgilityNative', params: { mediatype: 'banner_web_ms' } }
  ]
}]);
```

---

## 3. Banner IA Ad Units

Banner IA is available on `pubvibeXenon`, `pubvibeAniview`, `pubvibeISCream`, and `pubvibeAgilityNative`.

### pubvibeXenon — Banner IA (automatic, no params needed)

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-ia',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [{
    bidder: 'pubvibeXenon',
    params: {}
  }]
}]);
```

### pubvibeAniview — Banner IA

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-ia',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [{
    bidder: 'pubvibeAniview',
    params: { mediatype: 'banner_ia' }
  }]
}]);
```

### pubvibeISCream — Banner IA

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-ia',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [{
    bidder: 'pubvibeISCream',
    params: { mediatype: 'banner_ia' }
  }]
}]);
```

### pubvibeAgilityNative — Banner IA

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-ia',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [{
    bidder: 'pubvibeAgilityNative',
    params: { mediatype: 'banner_ia' }
  }]
}]);
```

### All Four Competing on Banner IA

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-ia',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [
    { bidder: 'pubvibeXenon',         params: {} },
    { bidder: 'pubvibeAniview',       params: { mediatype: 'banner_ia' } },
    { bidder: 'pubvibeISCream',       params: { mediatype: 'banner_ia' } },
    { bidder: 'pubvibeAgilityNative', params: { mediatype: 'banner_ia' } }
  ]
}]);
```

---

## 4. Native Ad Units

Native is only available on `pubvibeAgilityNative`.
**Always pass `mediatype: 'native'` explicitly.**

```javascript
pbjs.addAdUnits([{
  code: 'div-native-unit',
  mediaTypes: {
    native: {
      title:       { required: true, len: 80 },
      image:       { required: true, sizes: [[300, 250]] },
      sponsoredBy: { required: true },
      body:        { required: false }
    }
  },
  bids: [{
    bidder: 'pubvibeAgilityNative',
    params: { mediatype: 'native' }
  }]
}]);
```

---

## 5. Video IA Ad Units

Video IA is available on `pubvibeXenon`, `pubvibeAniview`, and `pubvibeISCream`.

### All Three Competing on Video IA

```javascript
pbjs.addAdUnits([{
  code: 'div-video-instream',
  mediaTypes: {
    video: {
      playerSize: [640, 480],
      context: 'instream',
      mimes: ['video/mp4', 'video/webm'],
      protocols: [1, 2, 3, 5, 6],
      minduration: 5,
      maxduration: 60,
      startdelay: 0,
      linearity: 1,
      skip: 1,
      skipmin: 5
    }
  },
  bids: [
    { bidder: 'pubvibeXenon',   params: {} },
    { bidder: 'pubvibeAniview', params: { mediatype: 'video_ia' } },
    { bidder: 'pubvibeISCream', params: { mediatype: 'video_ia' } }
  ]
}]);
```

---

## 6. Full Page Setup — All Adapters

```javascript
var pbjs = pbjs || {};
pbjs.que = pbjs.que || [];

pbjs.que.push(function() {

  // ── Step 1: Connect all adapters to your Prebid Server ─────────────────
  pbjs.setConfig({
    s2sConfig: {
      accountId: '1',
      bidders: ['pubvibeXenon', 'pubvibeAniview', 'pubvibeISCream', 'pubvibeAgilityNative'],
      defaultVendor: 'appnexus',
      timeout: 1500,
      endpoint: {
        p1Consent: 'https://YOUR_PREBID_SERVER_DOMAIN/openrtb2/auction',
        noConsent:  'https://YOUR_PREBID_SERVER_DOMAIN/openrtb2/auction'
      }
    },
    priceGranularity: 'medium',
    currency: { adServerCurrency: 'USD' }
  });

  // ── Step 2: Define all ad units ─────────────────────────────────────────
  pbjs.addAdUnits([

    // --- Leaderboard Banner IA — all 4 competing ---
    {
      code: 'div-leaderboard-728x90',
      mediaTypes: {
        banner: { sizes: [[728, 90], [970, 90]] }
      },
      bids: [
        { bidder: 'pubvibeXenon',         params: {} },
        { bidder: 'pubvibeAniview',       params: { mediatype: 'banner_ia' } },
        { bidder: 'pubvibeISCream',       params: { mediatype: 'banner_ia' } },
        { bidder: 'pubvibeAgilityNative', params: { mediatype: 'banner_ia' } }
      ]
    },

    // --- Medium Rectangle Banner IA ---
    {
      code: 'div-mrec-300x250',
      mediaTypes: {
        banner: { sizes: [[300, 250], [300, 600], [336, 280]] }
      },
      bids: [
        { bidder: 'pubvibeXenon',         params: {} },
        { bidder: 'pubvibeAniview',       params: { mediatype: 'banner_ia' } },
        { bidder: 'pubvibeISCream',       params: { mediatype: 'banner_ia' } },
        { bidder: 'pubvibeAgilityNative', params: { mediatype: 'banner_ia' } }
      ]
    },

    // --- Mid-page Banner Web MS — incremental fill ---
    {
      code: 'div-mid-banner-web-ms',
      mediaTypes: {
        banner: { sizes: [[300, 250], [320, 100]] }
      },
      bids: [
        { bidder: 'pubvibeAniview',       params: { mediatype: 'banner_web_ms' } },
        { bidder: 'pubvibeISCream',       params: { mediatype: 'banner_web_ms' } },
        { bidder: 'pubvibeAgilityNative', params: { mediatype: 'banner_web_ms' } }
      ]
    },

    // --- Mobile Banner IA ---
    {
      code: 'div-mobile-320x50',
      mediaTypes: {
        banner: { sizes: [[320, 50], [300, 50]] }
      },
      bids: [
        { bidder: 'pubvibeXenon',         params: {} },
        { bidder: 'pubvibeAniview',       params: { mediatype: 'banner_ia' } },
        { bidder: 'pubvibeISCream',       params: { mediatype: 'banner_ia' } },
        { bidder: 'pubvibeAgilityNative', params: { mediatype: 'banner_ia' } }
      ]
    },

    // --- Native ---
    {
      code: 'div-native-unit',
      mediaTypes: {
        native: {
          title:       { required: true, len: 80 },
          image:       { required: true, sizes: [[300, 250]] },
          sponsoredBy: { required: true },
          body:        { required: false }
        }
      },
      bids: [
        { bidder: 'pubvibeAgilityNative', params: { mediatype: 'native' } }
      ]
    },

    // --- Instream Video ---
    {
      code: 'div-video-instream',
      mediaTypes: {
        video: {
          playerSize: [640, 480],
          context: 'instream',
          mimes: ['video/mp4', 'video/webm'],
          protocols: [1, 2, 3, 5, 6],
          minduration: 5,
          maxduration: 60,
          startdelay: 0,
          linearity: 1,
          skip: 1,
          skipmin: 5
        }
      },
      bids: [
        { bidder: 'pubvibeXenon',   params: {} },
        { bidder: 'pubvibeAniview', params: { mediatype: 'video_ia' } },
        { bidder: 'pubvibeISCream', params: { mediatype: 'video_ia' } }
      ]
    }

  ]); // end addAdUnits

  // ── Step 3: Request bids and pass to your ad server ────────────────────
  pbjs.requestBids({
    bidsBackHandler: function() {
      // If using Google Ad Manager:
      googletag.cmd.push(function() {
        pbjs.setTargetingForGPTAsync();
        googletag.pubads().refresh();
      });
    }
  });

}); // end pbjs.que
```

---

## 7. HTML Slot Divs

```html
<!-- Leaderboard -->
<div id="div-leaderboard-728x90">
  <script>googletag.cmd.push(function(){ googletag.display('div-leaderboard-728x90'); });</script>
</div>

<!-- Medium Rectangle -->
<div id="div-mrec-300x250">
  <script>googletag.cmd.push(function(){ googletag.display('div-mrec-300x250'); });</script>
</div>

<!-- Mid-page Banner Web MS -->
<div id="div-mid-banner-web-ms">
  <script>googletag.cmd.push(function(){ googletag.display('div-mid-banner-web-ms'); });</script>
</div>

<!-- Mobile Banner -->
<div id="div-mobile-320x50">
  <script>googletag.cmd.push(function(){ googletag.display('div-mobile-320x50'); });</script>
</div>

<!-- Native -->
<div id="div-native-unit"></div>

<!-- Video Player -->
<div id="div-video-instream"></div>
```

---

## 8. Notes for Developers

- **Always specify `mediatype`** — do not rely on defaults unless you have confirmed the default matches the endpoint you want. The defaults differ per adapter (see table in section 2).
- **Preferred default is `banner_web_ms`** — if you are unsure which to use for a banner slot, use `banner_web_ms` for `pubvibeAniview`, `pubvibeISCream`, and `pubvibeAgilityNative`.
- **`pubvibeXenon` has no `mediatype` param** — it routes automatically: banner imps go to the banner endpoint, video imps go to the video endpoint. Pass `params: {}`.
- **No client-side adapter files needed** — all four adapters live in Prebid Server. Prebid.js only needs the `prebidServerBidAdapter` module (`modules/prebidServerBidAdapter`) in the build.
- **Native is exclusive to `pubvibeAgilityNative`** — the other three adapters do not support native.
- **Video is not supported by `pubvibeAgilityNative`** — use `pubvibeXenon`, `pubvibeAniview`, or `pubvibeISCream` for video.
- **Timeout** — `1500ms` is recommended for all adapters.
- **Currency** — all bids return USD. No currency conversion needed.
- **Best CPM strategy** — put all eligible adapters on every slot with the same `mediatype`. They compete and Prebid picks the highest bid automatically.
