function updateWinstons() {
  shuffleWinstons();
  window.setInterval(shuffleWinstons, 1000);
}

async function shuffleWinstons() {
  let winstons = document.getElementById("winstons").children;
  for (winston of winstons) {
    winston.style.position = "relative";
    winston.style.left = `${Math.floor(Math.random() * 200)}px`;
    winston.style.top = `${Math.floor(Math.random() * 50)}px`;
  }
}
