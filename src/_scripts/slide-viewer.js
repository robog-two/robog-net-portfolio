let currentSrc = null;

function updateSlideViewer() {
  const article = document.querySelector("article.post-content");
  if (!article) return;

  const centerFold = window.innerHeight / 2;
  const images = Array.from(article.querySelectorAll("img"));

  let currentImage = images[0];
  for (const img of images) {
    const rect = img.getBoundingClientRect();
    if (rect.bottom <= centerFold) {
      currentImage = img;
    }
  }
  const container = document.querySelector("#slideViewer");
  const viewer = document.querySelector("#slideViewer img:not(.transitioning)");
  if (currentImage && viewer) {
    if (currentImage.src !== currentSrc) {
      currentSrc = currentImage.src;

      const next = viewer.cloneNode(true);
      next.style.display = "block";
      next.src = currentImage.src;
      next.alt = currentImage.alt;
      container.appendChild(next);

      viewer.classList.add("transitioning");

      setTimeout(() => {
        for (
          const child of document.querySelectorAll(
            "#slideViewer .transitioning",
          )
        ) {
          container.removeChild(child);
        }
      }, 100);
    }
  }
}

window.addEventListener("scroll", updateSlideViewer, true);
window.addEventListener("load", updateSlideViewer);
updateSlideViewer();
