// Analytics for every page on the site, including the static gallery pieces and
// spheres. Sends the same data to Plausible and Umami.
//   Plausible docs: https://plausible.io/docs/script-extensions
//   Umami docs:     https://docs.umami.is/docs/tracker-functions
//
// Built in to both: pageviews, sources/UTMs, entry pages, devices.
// Built in to Plausible only: time on page, scroll depth, outbound links and
// file downloads (enabled below). Umami gets those from this file instead.
//
// Events:
//   "First Click"           first link clicked in a session {to, text, kind, area}
//   "Outbound Link: Click"  (Umami; Plausible sends its own) {url}
//   "File Download"         (Umami; Plausible sends its own) {url}
//   "Page Engagement"       (Umami) visible seconds + max scroll % for a pageview
//   "Gallery Piece"         (Umami) {piece, pieces_seen}; Plausible gets
//                           pieces_seen as a pageview property
//
// Fault tolerance: each service is loaded and called independently, and every
// call into it is wrapped in try/catch. If one is blocked (ad blocker, network
// error, outage) or throws, the other keeps working and the page is unaffected.
// A link click that is held for tracking always navigates.
(function () {
  "use strict";

  // ---------------------------------------------------------------------------
  // Helpers
  // ---------------------------------------------------------------------------

  // sessionStorage can throw (private mode, blocked storage); fail quietly
  function readStore(key) {
    try { return sessionStorage.getItem(key); } catch (e) { return null; }
  }
  function writeStore(key, value) {
    try { sessionStorage.setItem(key, value); } catch (e) { /* ignore */ }
  }

  // Run fn, swallowing any error so one failure can't stop the rest
  function safe(fn) {
    try { fn(); } catch (e) { /* analytics must never break the page */ }
  }

  function loadScript(src, attrs, onload, onerror) {
    safe(function () {
      var script = document.createElement("script");
      script.async = true;
      script.src = src;
      for (var name in attrs) script.setAttribute(name, attrs[name]);
      script.onload = function () { safe(onload); };
      script.onerror = function () { safe(onerror); };
      (document.head || document.documentElement).appendChild(script);
    });
  }

  // Same extension list Plausible uses for file downloads
  var DOWNLOAD_EXT = /\.(pdf|xlsx?|docx?|txt|rtf|csv|exe|key|pp[st]x?|7z|pkg|rar|gz|zip|avi|mov|mp4|mpe?g|wmv|midi?|mp3|wav|wma|dmg|iso|msi)$/i;

  // ---------------------------------------------------------------------------
  // Umami. It has no command queue, so calls made before script.js loads are
  // held here and sent once it has.
  // ---------------------------------------------------------------------------

  var umamiQueue = [];
  var umamiFailed = false;
  function umamiSend(name, data) {
    safe(function () { window.umami.track(name, data); });
  }
  function umamiTrack(name, data) {
    if (umamiFailed) return;
    if (window.umami && typeof window.umami.track === "function") umamiSend(name, data);
    else umamiQueue.push([name, data]);
  }

  loadScript("https://cloud.umami.is/script.js", {
    "data-website-id": "5489e229-f2eb-4bf0-afd1-115a5206bbda",
    // Only count the live site, not localhost or preview deploys
    "data-domains": "robog.net",
  }, function () {
    var queued = umamiQueue;
    umamiQueue = [];
    // If the script loaded but didn't set up window.umami (some blockers
    // serve an empty file), treat it as blocked
    if (!window.umami || typeof window.umami.track !== "function") {
      umamiFailed = true;
      return;
    }
    queued.forEach(function (call) { umamiSend(call[0], call[1]); });
  }, function () {
    umamiFailed = true;
    umamiQueue = [];
  });

  // ---------------------------------------------------------------------------
  // Gallery pieces: count distinct pieces opened in this browser session
  // ---------------------------------------------------------------------------

  var pieceMatch = location.pathname.match(/^\/gallery\/piece\/([^/]+)/);
  var piecesSeen = null;
  if (pieceMatch) safe(function () {
    var seen = (readStore("pa-pieces") || "").split(",").filter(Boolean);
    if (seen.indexOf(pieceMatch[1]) === -1) seen.push(pieceMatch[1]);
    writeStore("pa-pieces", seen.join(","));
    piecesSeen = seen.length;
    umamiTrack("Gallery Piece", { piece: pieceMatch[1], pieces_seen: piecesSeen });
  });

  // ---------------------------------------------------------------------------
  // Plausible
  // ---------------------------------------------------------------------------

  var plausibleFailed = false;
  function plausibleTrack(name, options) {
    if (plausibleFailed) return false;
    try {
      window.plausible(name, options);
      return true;
    } catch (e) {
      return false;
    }
  }

  loadScript("https://plausible.io/js/pa-zOwKZXLFpAVB8BbtYJb3V.js", {}, null, function () {
    plausibleFailed = true;
  });

  safe(function () {
    window.plausible = window.plausible || function () {
      (plausible.q = plausible.q || []).push(arguments);
    };
    plausible.init = plausible.init || function (i) { plausible.o = i || {}; };

    plausible.init({
      outboundLinks: true,
      fileDownloads: true,
      customProperties: function (eventName) {
        if (eventName !== "pageview" || !piecesSeen) return {};
        return { pieces_seen: String(piecesSeen) };
      },
    });
  });

  // ---------------------------------------------------------------------------
  // Link clicks: "First Click" for both, outbound links and downloads for Umami
  // ---------------------------------------------------------------------------

  function onClick(event) {
    var link = event.target.closest && event.target.closest("a[href]");
    if (!link) return;

    var url;
    try { url = new URL(link.href, location.href); } catch (e) { return; }
    var internal = url.origin === location.origin;
    // Ignore in-page jumps like the "Skip to content" link
    if (internal && url.pathname === location.pathname && url.hash) return;

    if (!internal && /^https?:$/.test(url.protocol)) {
      umamiTrack("Outbound Link: Click", { url: link.href });
    }
    if (DOWNLOAD_EXT.test(url.pathname) || link.hasAttribute("download")) {
      umamiTrack("File Download", { url: link.href });
    }

    if (readStore("pa-first-click")) return;
    writeStore("pa-first-click", "1");

    var props = {
      to: internal ? url.pathname : url.hostname + url.pathname,
      text: (link.textContent || link.getAttribute("aria-label") || "").trim().replace(/\s+/g, " ").slice(0, 100),
      kind: internal ? "internal" : "external",
      area: link.closest("nav") ? "nav" : "content",
    };

    // Umami sends with fetch keepalive, which survives navigation
    umamiTrack("First Click", props);

    // A same-tab navigation can cancel Plausible's request, so hold the
    // navigation until it confirms (or briefly times out), like its own
    // outbound link tracking does
    var sameTab = event.type === "click" && event.button === 0 &&
      !event.metaKey && !event.ctrlKey && !event.shiftKey && !event.altKey &&
      (!link.target || link.target === "_self") && !link.hasAttribute("download");

    // Only hold the click when Plausible might still deliver; if it's
    // blocked, let the browser navigate right away
    if (!sameTab || event.defaultPrevented || plausibleFailed) {
      plausibleTrack("First Click", { props: props });
      return;
    }

    var done = false;
    function go() {
      if (done) return;
      done = true;
      location.href = link.href;
    }
    // Schedule the fallback before anything else can fail, so a held click
    // always navigates
    setTimeout(go, 500);
    event.preventDefault();
    if (!plausibleTrack("First Click", { props: props, callback: go })) go();
  }

  function handleClick(event) {
    try { onClick(event); } catch (e) { /* never block a click */ }
  }
  document.addEventListener("click", handleClick);
  // Middle-click opens a new tab
  document.addEventListener("auxclick", function (event) {
    if (event.button === 1) handleClick(event);
  });

  // ---------------------------------------------------------------------------
  // Time on page and scroll depth for Umami (Plausible measures these itself).
  // Sent once per pageview, the first time the page is hidden: navigating away,
  // closing the tab, or switching apps on mobile.
  // ---------------------------------------------------------------------------

  var visibleMs = 0;
  var visibleSince = document.visibilityState === "visible" ? Date.now() : null;
  var maxScroll = 0;
  var engagementSent = false;

  function measureScroll() {
    var doc = document.documentElement;
    var height = Math.max(doc.scrollHeight, document.body ? document.body.scrollHeight : 0);
    var seen = height <= window.innerHeight ? 1 : (window.scrollY + window.innerHeight) / height;
    maxScroll = Math.max(maxScroll, Math.min(1, seen));
  }

  function timeBucket(seconds) {
    if (seconds < 10) return "0-10s";
    if (seconds < 30) return "10-30s";
    if (seconds < 60) return "30-60s";
    if (seconds < 180) return "1-3m";
    if (seconds < 600) return "3-10m";
    return "10m+";
  }

  function sendEngagement() {
    if (visibleSince !== null) {
      visibleMs += Date.now() - visibleSince;
      visibleSince = null;
    }
    if (engagementSent) return;
    engagementSent = true;
    measureScroll();
    var seconds = Math.round(visibleMs / 1000);
    umamiTrack("Page Engagement", {
      seconds: seconds,
      time: timeBucket(seconds),
      // Rounded down to the nearest 10% so the breakdown stays readable
      scroll: Math.floor(maxScroll * 10) * 10,
    });
  }

  window.addEventListener("scroll", function () { safe(measureScroll); }, { passive: true });
  window.addEventListener("load", function () { safe(measureScroll); });
  document.addEventListener("visibilitychange", function () {
    safe(function () {
      if (document.visibilityState === "hidden") sendEngagement();
      else if (visibleSince === null) visibleSince = Date.now();
    });
  });
  window.addEventListener("pagehide", function () { safe(sendEngagement); });
})();
