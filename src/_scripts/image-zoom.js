document.addEventListener("DOMContentLoaded", function () {
  // Only add zoom functionality to blog post images
  const blogImages = document.querySelectorAll(".post-content img");

  blogImages.forEach((img) => {
    img.addEventListener("click", function () {
      openImageModal(img);
    });
  });

  function openImageModal(img) {
    // Create modal container
    const modal = document.createElement("div");
    modal.className = "image-modal";
    modal.innerHTML = `
            <div class="modal-backdrop">
                <div class="modal-content">
                    <img src="${img.src}" alt="${
      img.alt || "Image"
    }" class="modal-image">
                    <button class="modal-close" aria-label="Close image">×</button>
                </div>
            </div>
        `;

    // Add to page
    document.body.appendChild(modal);
    document.body.style.overflow = "hidden";

    // Add event listeners
    const closeBtn = modal.querySelector(".modal-close");
    const backdrop = modal.querySelector(".modal-backdrop");

    closeBtn.addEventListener("click", closeModal);
    backdrop.addEventListener("click", function (e) {
      if (e.target === backdrop) {
        closeModal();
      }
    });

    // if (modal.requestFullscreen) {
    //   modal.requestFullscreen();
    // } else if (document.webkitRequestFullscreen) { // Safari
    //   modal.webkitRequestFullscreen();
    // } else if (document.msRequestFullscreen) { // IE11
    //   modal.msRequestFullscreen();
    // }

    // // Listen for exiting fullscreen and close the modal when that happens
    // document.addEventListener("fullscreenchange", function () {
    //   closeImageModal();
    // });

    // Close on Escape key
    document.addEventListener("keydown", function (e) {
      if (e.key === "Escape") {
        closeModal();
      }
    });

    function closeModal() {
      if (document.fullscreenElement) {
        if (document.exitFullscreen) {
          document.exitFullscreen();
        } else if (document.webkitExitFullscreen) { // Safari
          document.webkitExitFullscreen();
        } else if (document.msExitFullscreen) { // IE11
          document.msExitFullscreen();
        }
      }
      document.body.removeChild(modal);
      document.body.style.overflow = "";
    }

    // Animate in
    setTimeout(() => {
      modal.classList.add("show");
    }, 10);
  }
});
