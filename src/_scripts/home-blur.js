document.addEventListener("DOMContentLoaded", function () {
  const navButtons = document.querySelectorAll(".corner-nav");
  const bizCard = document.querySelector(".bizCard");

  navButtons.forEach((button) => {
    function activate() {
      if (button.timeout) {
        clearTimeout(button.timeout);
      }
      button.classList.add("highest-nav");
      bizCard.classList.add("nav-hovered");
    }

    function deactivate() {
      button.timeout = setTimeout(() => {
        button.classList.remove("highest-nav");
      }, 500);
      bizCard.classList.remove("nav-hovered");
    }

    // Keyboard focus gets the same treatment as mouse hover
    button.addEventListener("mouseenter", activate);
    button.addEventListener("focus", activate);
    button.addEventListener("mouseleave", deactivate);
    button.addEventListener("blur", deactivate);
  });
});
