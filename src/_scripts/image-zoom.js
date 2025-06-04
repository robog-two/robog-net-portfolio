document.addEventListener("DOMContentLoaded", function () {
  // Only add zoom functionality to blog post images
  const blogImages = document.querySelectorAll(".post-content img");

  blogImages.forEach((img) => {
    // This is kind of messy but I'm not writing a whole eleventy plugin that does this. Not right now, anyways
    const caption =
      img.parentElement.parentElement.nextSibling.nextElementSibling;
    if (caption.tagName == "FIGCAPTION") {
      const newCaption = document.createElement("figcaption");
      newCaption.innerHTML = caption.innerHTML; // Innerhtml in case there is text formatting

      const newFigure = document.createElement("figure");
      const oldInner = img.parentElement.parentElement.innerHTML;
      img.parentElement.parentElement.replaceWith(newFigure);
      newFigure.innerHTML = oldInner;

      caption.remove();
      newFigure.appendChild(newCaption);

      newFigure.childNodes[0].childNodes[1].addEventListener(
        "click",
        function () {
          openImageModal(this);
        },
      );
    } else {
      img.addEventListener("click", function () {
        openImageModal(this);
      });
    }
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

    // Close on Escape key
    document.addEventListener("keydown", function (e) {
      if (e.key === "Escape") {
        closeModal();
      }
    });

    function closeModal() {
      document.body.removeChild(modal);
      document.body.style.overflow = "";
    }

    // Animate in
    setTimeout(() => {
      modal.classList.add("show");
    }, 10);
  }
});
