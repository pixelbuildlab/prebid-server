# Pubvibe Adapters — Prebid.js Integration Guide

**For:** Prebid.js developer  
**Site:** Your site (reference: jobsedutimes.com model)  
**Server:** Prebid Server (Go) running your custom adapters  
**Mode:** Server-to-Server (s2sConfig) only — all bidding goes through your Prebid Server  

---

## Adapters Available

| Adapter Name | Provider | Banner IA | Banner Web MS | Video IA |
|---|---|---|---|---|
| `pubvibeXenon` | xenrtb.com | ✅ | ❌ | ✅ |
| `pubvibeAniview` | aniview.com | ✅ | ✅ | ✅ |
| `pubvibeISCream` | agilityadvsrv.com | ✅ | ✅ | ✅ |

---

## 1. Prebid Server s2sConfig (Required — set once)

This connects Prebid.js to your Prebid Server. All three adapters route through it.

```javascript
pbjs.setConfig({
  s2sConfig: {
    accountId: '1',
    bidders: ['pubvibeXenon', 'pubvibeAniview', 'pubvibeISCream'],
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

## 2. Bidder Parameters Reference

### pubvibeXenon

| Parameter | Type | Required | Values | Description |
|---|---|---|---|---|
| *(none)* | — | No | — | No params needed. Routing is handled server-side by impression type |

Routing behavior (automatic, no params needed):
- `imp.banner` → Banner endpoint
- `imp.video` → Video endpoint

---

### pubvibeAniview

| Parameter | Type | Required | Values | Description |
|---|---|---|---|---|
| `mediatype` | string | No | `banner_ia`, `banner_web_ms`, `video_ia` | Selects which Aniview endpoint to use |

Default routing when `mediatype` is omitted:
- `imp.banner` → `banner_ia`
- `imp.video` → `video_ia`

---

### pubvibeISCream

| Parameter | Type | Required | Values | Description |
|---|---|---|---|---|
| `mediatype` | string | No | `banner_ia`, `banner_web_ms`, `video_ia` | Selects which ISCream endpoint to use |

Default routing when `mediatype` is omitted:
- `imp.banner` → `banner_ia`
- `imp.video` → `video_ia`

---

## 3. Single Adapter Usage

### pubvibeXenon — Banner only

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-300x250',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [{
    bidder: 'pubvibeXenon',
    params: {}
  }]
}]);
```

### pubvibeXenon — Video only

```javascript
pbjs.addAdUnits([{
  code: 'div-video-unit',
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
    bidder: 'pubvibeXenon',
    params: {}
  }]
}]);
```

---

### pubvibeAniview — Banner IA

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-300x250',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [{
    bidder: 'pubvibeAniview',
    params: { mediatype: 'banner_ia' }
  }]
}]);
```

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

### pubvibeAniview — Video IA

```javascript
pbjs.addAdUnits([{
  code: 'div-video-unit',
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
    bidder: 'pubvibeAniview',
    params: { mediatype: 'video_ia' }
  }]
}]);
```

---

### pubvibeISCream — Banner IA

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-300x250',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90]] }
  },
  bids: [{
    bidder: 'pubvibeISCream',
    params: { mediatype: 'banner_ia' }
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

### pubvibeISCream — Video IA

```javascript
pbjs.addAdUnits([{
  code: 'div-video-unit',
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
    bidder: 'pubvibeISCream',
    params: { mediatype: 'video_ia' }
  }]
}]);
```

---

## 4. Combined Usage — All 3 Adapters on One Site

This is the recommended setup for maximum fill and best CPM. All 3 adapters compete
on the same ad unit. Prebid picks the highest bid.

### Banner Ad Unit — All 3 Adapters Competing

```javascript
pbjs.addAdUnits([{
  code: 'div-banner-300x250',
  mediaTypes: {
    banner: { sizes: [[300, 250], [728, 90], [320, 50]] }
  },
  bids: [
    // Xenon — no params needed
    {
      bidder: 'pubvibeXenon',
      params: {}
    },
    // Aniview — Banner IA endpoint
    {
      bidder: 'pubvibeAniview',
      params: { mediatype: 'banner_ia' }
    },
    // ISCream — Banner IA endpoint
    {
      bidder: 'pubvibeISCream',
      params: { mediatype: 'banner_ia' }
    }
  ]
}]);
```

### Banner Ad Unit — With Web MS variants competing too

Use a second ad unit or duplicate bids with `banner_web_ms` to get more demand:

```javascript
pbjs.addAdUnits([
  // Unit 1: Banner IA demand
  {
    code: 'div-banner-top',
    mediaTypes: {
      banner: { sizes: [[728, 90], [970, 90]] }
    },
    bids: [
      { bidder: 'pubvibeXenon',   params: {} },
      { bidder: 'pubvibeAniview', params: { mediatype: 'banner_ia' } },
      { bidder: 'pubvibeISCream', params: { mediatype: 'banner_ia' } }
    ]
  },
  // Unit 2: Banner Web MS demand (different inventory)
  {
    code: 'div-banner-mid',
    mediaTypes: {
      banner: { sizes: [[300, 250], [336, 280]] }
    },
    bids: [
      { bidder: 'pubvibeAniview', params: { mediatype: 'banner_web_ms' } },
      { bidder: 'pubvibeISCream', params: { mediatype: 'banner_web_ms' } }
    ]
  }
]);
```

### Video Ad Unit — All 3 Adapters Competing

```javascript
pbjs.addAdUnits([{
  code: 'div-video-player',
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

## 5. Full Page Setup — jobsedutimes.com Style

A complete multi-slot setup similar to a news/jobs site with 6-8 demand sources
across banner and video. All demand flows through your single Prebid Server.

```javascript
var pbjs = pbjs || {};
pbjs.que = pbjs.que || [];

pbjs.que.push(function() {

  // ── Step 1: Connect all adapters to your Prebid Server ─────────────────
  pbjs.setConfig({
    s2sConfig: {
      accountId: '1',
      bidders: ['pubvibeXenon', 'pubvibeAniview', 'pubvibeISCream'],
      defaultVendor: 'appnexus',
      timeout: 1500,
      endpoint: {
        p1Consent: 'https://YOUR_PREBID_SERVER_DOMAIN/openrtb2/auction',
        noConsent:  'https://YOUR_PREBID_SERVER_DOMAIN/openrtb2/auction'
      }
    },
    // Price granularity — important for GAM line items
    priceGranularity: 'medium',
    // Currency
    currency: { adServerCurrency: 'USD' }
  });

  // ── Step 2: Define all ad units ─────────────────────────────────────────
  pbjs.addAdUnits([

    // --- Leaderboard (728x90) — top of page ---
    {
      code: 'div-leaderboard-728x90',
      mediaTypes: {
        banner: { sizes: [[728, 90], [970, 90]] }
      },
      bids: [
        { bidder: 'pubvibeXenon',   params: {} },
        { bidder: 'pubvibeAniview', params: { mediatype: 'banner_ia' } },
        { bidder: 'pubvibeISCream', params: { mediatype: 'banner_ia' } }
      ]
    },

    // --- Medium Rectangle (300x250) — sidebar/inline ---
    {
      code: 'div-mrec-300x250',
      mediaTypes: {
        banner: { sizes: [[300, 250], [300, 600], [336, 280]] }
      },
      bids: [
        { bidder: 'pubvibeXenon',   params: {} },
        { bidder: 'pubvibeAniview', params: { mediatype: 'banner_ia' } },
        { bidder: 'pubvibeISCream', params: { mediatype: 'banner_ia' } }
      ]
    },

    // --- Mid-page Banner — Web MS demand (higher fill on content pages) ---
    {
      code: 'div-mid-banner-300x250',
      mediaTypes: {
        banner: { sizes: [[300, 250], [320, 100]] }
      },
      bids: [
        { bidder: 'pubvibeAniview', params: { mediatype: 'banner_web_ms' } },
        { bidder: 'pubvibeISCream', params: { mediatype: 'banner_web_ms' } }
      ]
    },

    // --- Mobile Banner (320x50) ---
    {
      code: 'div-mobile-320x50',
      mediaTypes: {
        banner: { sizes: [[320, 50], [300, 50]] }
      },
      bids: [
        { bidder: 'pubvibeXenon',   params: {} },
        { bidder: 'pubvibeAniview', params: { mediatype: 'banner_ia' } },
        { bidder: 'pubvibeISCream', params: { mediatype: 'banner_ia' } }
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
    bidsBackHandler: function(bids) {

      // If using Google Ad Manager:
      googletag.cmd.push(function() {
        pbjs.setTargetingForGPTAsync();
        googletag.pubads().refresh();
      });

      // If NOT using GAM — render directly:
      // pbjs.renderAd(document, bids['div-mrec-300x250'].bids[0].adId);
    }
  });

}); // end pbjs.que
```

---

## 6. HTML Slot Divs

Place these in your page HTML where ads should appear:

```html
<!-- Leaderboard -->
<div id="div-leaderboard-728x90">
  <script>googletag.cmd.push(function(){ googletag.display('div-leaderboard-728x90'); });</script>
</div>

<!-- Medium Rectangle -->
<div id="div-mrec-300x250">
  <script>googletag.cmd.push(function(){ googletag.display('div-mrec-300x250'); });</script>
</div>

<!-- Mid-page Banner (Web MS) -->
<div id="div-mid-banner-300x250">
  <script>googletag.cmd.push(function(){ googletag.display('div-mid-banner-300x250'); });</script>
</div>

<!-- Mobile Banner -->
<div id="div-mobile-320x50">
  <script>googletag.cmd.push(function(){ googletag.display('div-mobile-320x50'); });</script>
</div>

<!-- Video Player -->
<div id="div-video-instream"></div>
```

---

## 7. Quick Reference — mediatype Values per Adapter

| Adapter | Banner IA | Banner Web MS | Video IA |
|---|---|---|---|
| `pubvibeXenon` | automatic | ❌ not supported | automatic |
| `pubvibeAniview` | `"banner_ia"` | `"banner_web_ms"` | `"video_ia"` |
| `pubvibeISCream` | `"banner_ia"` | `"banner_web_ms"` | `"video_ia"` |

---

## 8. Notes for Developer

- **No client-side adapter files needed** — all three adapters live in Prebid Server. Prebid.js only needs the `prebidServer` module in the build.
- **Prebid.js build** must include the module `modules/prebidServerBidAdapter`. No other custom modules needed for these three adapters.
- **Timeout** — `1500ms` is recommended. ISCream and Aniview are RTB endpoints so they respond fast, but allow headroom.
- **mediatype is optional** — if you omit it, the server auto-routes: banner imps go to Banner IA, video imps go to Video IA. Only specify it explicitly when you want the Web MS endpoint.
- **Best CPM strategy** — put all 3 adapters on every banner unit with the same `mediatype`. They compete and Prebid picks the winner automatically.
- **Web MS slots** — add `banner_web_ms` as a second bid on the same slot (or a separate mid-content slot) for incremental fill.
- **Currency** — all bids return USD. No currency conversion needed.
