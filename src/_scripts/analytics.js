// Plausible analytics for every page on the site, including the static gallery
// pieces and spheres. Docs: https://plausible.io/docs/script-extensions
//
// Built in (no code needed): pageviews, sources/UTMs, entry pages, devices,
// time on page and scroll depth.
// Extra here:
//   - Outbound link clicks and file downloads (goal "Outbound Link: Click",
//     "File Download")
//   - "First Click": the first link a visitor clicks in a session, with the
//     link destination and text as properties
//   - "pieces_seen" property on gallery piece pageviews: how many different
//     pieces this visitor has opened so far in the session
(function () {
  "use strict";

  var script = document.createElement("script");
  script.async = true;
  script.src = "https://plausible.io/js/pa-zOwKZXLFpAVB8BbtYJb3V.js";
  document.head.appendChild(script);

  // sessionStorage can throw (private mode, blocked storage); fail quietly
  function readStore(key) {
    try { return sessionStorage.getItem(key); } catch (e) { return null; }
  }
  function writeStore(key, value) {
    try { sessionStorage.setItem(key, value); } catch (e) { /* ignore */ }
  }

  // Count distinct gallery pieces opened in this browser session
  function piecesSeen() {
    var match = location.pathname.match(/^\/gallery\/piece\/([^/]+)/);
    if (!match) return null;
    var seen = (readStore("pa-pieces") || "").split(",").filter(Boolean);
    if (seen.indexOf(match[1]) === -1) seen.push(match[1]);
    writeStore("pa-pieces", seen.join(","));
    return seen.length;
  }

  window.plausible = window.plausible || function () {
    (plausible.q = plausible.q || []).push(arguments);
  };
  plausible.init = plausible.init || function (i) { plausible.o = i || {}; };

  plausible.init({
    outboundLinks: true,
    fileDownloads: true,
    customProperties: function (eventName) {
      if (eventName !== "pageview") return {};
      var count = piecesSeen();
      return count ? { pieces_seen: String(count) } : {};
    },
  });

  // "First Click": only the first link click per session is sent
  function onClick(event) {
    var link = event.target.closest && event.target.closest("a[href]");
    if (!link || readStore("pa-first-click")) return;

    var url = new URL(link.href, location.href);
    // Ignore in-page jumps like the "Skip to content" link
    if (url.origin === location.origin && url.pathname === location.pathname && url.hash) return;

    writeStore("pa-first-click", "1");

    var internal = url.origin === location.origin;
    var props = {
      to: internal ? url.pathname : url.hostname + url.pathname,
      text: (link.textContent || link.getAttribute("aria-label") || "").trim().replace(/\s+/g, " ").slice(0, 100),
      kind: internal ? "internal" : "external",
      area: link.closest("nav") ? "nav" : "content",
    };

    // A same-tab navigation can cancel the request, so hold the navigation
    // until Plausible confirms (or briefly times out), like its own outbound
    // link tracking does
    var sameTab = event.type === "click" && event.button === 0 &&
      !event.metaKey && !event.ctrlKey && !event.shiftKey && !event.altKey &&
      (!link.target || link.target === "_self") && !link.hasAttribute("download");

    if (!sameTab || event.defaultPrevented) {
      plausible("First Click", { props: props });
      return;
    }

    event.preventDefault();
    var done = false;
    function go() {
      if (done) return;
      done = true;
      location.href = link.href;
    }
    plausible("First Click", { props: props, callback: go });
    setTimeout(go, 500);
  }

  document.addEventListener("click", onClick);
  // Middle-click opens a new tab
  document.addEventListener("auxclick", function (event) {
    if (event.button === 1) onClick(event);
  });
})();
