document.addEventListener("DOMContentLoaded", function () {
  const navButtons = document.querySelectorAll(".corner-nav");
  const bizCard = document.querySelector(".bizCard");

  navButtons.forEach((button) => {
    button.addEventListener("mouseenter", function () {
      if (button.timeout) {
        clearTimeout(button.timeout);
        console.log("yep");
      }
      button.classList.add("highest-nav");
      bizCard.classList.add("nav-hovered");
    });

    button.addEventListener("mouseleave", function () {
      button.timeout = setTimeout(() => {
        button.classList.remove("highest-nav");
      }, 500);
      bizCard.classList.remove("nav-hovered");
    });
  });
});
