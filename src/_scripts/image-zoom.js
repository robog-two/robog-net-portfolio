document.addEventListener("DOMContentLoaded", function () {
  // Only add zoom functionality to blog post images (not the slide
  // viewer's own display image, which is decorative)
  const blogImages = document.querySelectorAll(
    ".post-content img:not(#slideViewer img)",
  );

  blogImages.forEach((img) => {
    // Let keyboard and screen reader users open the zoomed view too
    img.setAttribute("role", "button");
    img.setAttribute(
      "aria-label",
      img.alt ? `Enlarge image: ${img.alt}` : "Enlarge image",
    );

    img.addEventListener("click", function () {
      openImageModal(img);
    });
    img.addEventListener("keydown", function (e) {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        openImageModal(img);
      }
    });
  });

  // Images that are hidden (e.g. collapsed by the slide viewer) must not be
  // tab stops, or keyboard focus would land somewhere invisible.
  function updateTabStops() {
    blogImages.forEach((img) => {
      img.tabIndex = img.getBoundingClientRect().height > 0 ? 0 : -1;
    });
  }
  updateTabStops();
  window.addEventListener("load", updateTabStops);
  window.addEventListener("resize", updateTabStops);

  function openImageModal(img) {
    // A native <dialog> moves focus inside, traps it, closes on Escape,
    // and is announced as a dialog by screen readers.
    const modal = document.createElement("dialog");
    modal.className = "image-modal";
    modal.setAttribute("aria-label", img.alt || "Enlarged image");
    modal.innerHTML = `
            <div class="modal-backdrop">
                <div class="modal-content">
                    <img class="modal-image">
                    <button class="modal-close" aria-label="Close image">×</button>
                </div>
            </div>
        `;
    const modalImage = modal.querySelector(".modal-image");
    modalImage.src = img.currentSrc || img.src;
    modalImage.alt = img.alt || "";

    // Add to page
    document.body.appendChild(modal);
    document.body.style.overflow = "hidden";
    modal.showModal();

    const closeBtn = modal.querySelector(".modal-close");
    const backdrop = modal.querySelector(".modal-backdrop");
    closeBtn.focus();

    closeBtn.addEventListener("click", () => modal.close());
    backdrop.addEventListener("click", function (e) {
      if (e.target === backdrop) {
        modal.close();
      }
    });

    // Fires for the close button, backdrop click and the Escape key
    modal.addEventListener("close", function () {
      modal.remove();
      document.body.style.overflow = "";
      img.focus();
    });

    // Animate in
    setTimeout(() => {
      modal.classList.add("show");
    }, 10);
  }
});
