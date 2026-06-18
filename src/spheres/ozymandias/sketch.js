// disable context menu for right click
document.oncontextmenu = function() {
        return false;
    }

let font;

let boardSideLength;

let time = 4000;
let startMillis = 0;

let sub = 0;
let fontscale = 1;

let bg;

// 0 - minesweeper screen, 1 - tutorial screen, 2 - game over screen, 3 - win screen
let screen = 1

let scores = -1;

function initvars() {
  time = 4000;
  startMillis = -1;
  sub = 0;
  fontscale = 1;
  screen = 1
  scores = -1;
  
  if (windowWidth >= 1116) {
    sub = 100;
    boardSideLength = Math.floor((Math.min(windowWidth, windowHeight) * 0.8)/20) * 20;
  } else {
    fontscale = 0.75;
    boardSideLength = Math.floor((Math.min(windowWidth, windowHeight) * 0.95)/20) * 20;
  }
  
  window.minesweeper.setupSweep(
    boardSideLength,
    (windowWidth - boardSideLength)/2.0,
    (windowHeight - boardSideLength)/2.0
  );
  
  console.log(window.minesweeper)
}

function setup() {
  font = loadFont("RomanceA.ttf");
  createCanvas(windowWidth, windowHeight, P2D);
  bg = createGraphics(width * 0.5, height * 0.5, WEBGL)
  initvars();
}

function updateScores() {
  scores = 0
fetch("https://robog-two-ozysb-91.deno.dev/highscores").then((r) => {
  r.json().then((r) => {
    scores = r;
  })
})
}

function drawScores() {
      fill(255);
      textSize(32 * fontscale);
      if (scores === -1) {
        updateScores();
      }
      if (scores === 0) {
        noStroke();
        let midpoint = width/2 - 25
        let componentA = map(sin(millis()/500.0), -1, 1, 50, midpoint)
        rect(
          componentA,
          200-sub,
          map(sin(millis()/500.0), -1, 1, 0, width - componentA),
          25
         );
        return;
      }
      if (scores.one.score !== 0) text("1. " + scores.one.username + " finished in " + scores.one.score + " B.C.", 50, 200-sub, width - 100, height);
      if (scores.two.score !== 0) text("2. " + scores.two.username + " finished in " + scores.two.score + " B.C.", 50, 250-sub, width - 100, height);
      if (scores.three.score !== 0) text("3. " + scores.three.username + " finished in " + scores.three.score + " B.C.", 50, 300-sub, width - 100, height);
}

function draw() {
  background(25);
  bg.background(25);
  textFont(font);
  
  bg.push();
  bg.rotateX(70);
  bg.stroke(117, 106, 69);
  bg.fill(25);
  desertOcean();
  bg.pop();
  
  image(bg, 0, 0, width, height);
  
  switch (screen) {
    case 0:
      fill(255);
      textSize(32 * fontscale);
      time = Math.floor(4000.0 - ((millis() - startMillis) / 30));
      if (time < 0) {
        screen = 2;
      }
      text(time + " B.C.", 50, 75, width - 100, height - 125);
      window.minesweeper.drawSweep();
      break;
    case 1:
      fill(234, 212, 138);
      textSize(50 * fontscale);
      text("Excavating Ozymandias' Empire", 50, 75, width - 100, height - 125);
      fill(255);
      textSize(32 * fontscale);
      text("Click or tap to excavate squares. The number indicates how many relics are in the surrounding squares. Do not excavate relics, they are delicate! You have 4000 years to find the remnants before they erode. Click or tap to begin.", 50, 200, width - 100, height - 250);
      break;
    case 4:
      fill(234, 212, 138);
      textSize(64 * fontscale);
      text("Scores to beat:", 50, 75, width - 100, height - 125);
      drawScores();
      break;
    case 2:
      window.minesweeper.drawSweep();
      textAlign(LEFT);
      fill(255);
      textSize(32  * fontscale);
      if (sub === 0) {
        text("Click or tap to try again.", 50, 75, width - 100, height - 125);
      } else {
        text("Ozymandias' empire was no match for the sands of time - or your terrible archaeology skills, apparently. Click or tap to try again.", 50, 75, width - 100, height - 125);
      }
      break;
    case 3:
      textAlign(LEFT);
      fill(234, 212, 138);
      textSize(64  * fontscale);
      text("Dig Completed by " + time + " B.C.", 50, 75, width - 100, height - 125);
      drawScores();
      break;
  }
}

function mousePressed() {
  if (screen === 2 || screen === 3) {
    initvars();
    screen = 1;
    return;
  }
  
  if (screen === 0) window.minesweeper.mousePressedSweep();
  
  if (screen === 4) {
    startMillis = millis();
    scores = -1;
    screen = 0;
    return;
  }
  
  if (screen === 1) screen = 4;
}

window.onGameOver = () => {
  console.log("Game lost.");
  startMillis = millis();
  // TODO: run stuff on game over
  screen = 2;
}

window.onWin = () => {
  console.log("Game won!");
  startMillis = millis();
  screen = 3;
  
  let pname = window.prompt("Please type your first name and last initial for the leaderboard. (ex: John D.)");
  if (pname == "" || pname == undefined || pname.length > 16) return; // very sorry to people with long names but unfortunately I'm terrible at programming so this will have to do. You can always just go around this check if you really need but it will ruin the display of the scoreboard.
  fetch(
     "https://robog-two-ozysb-91.deno.dev/setnewhighscore", {
     method: "POST",
     cache: "no-cache",
     headers: {
       "Content-Type": "application/json",
     },
     redirect: "follow",
     referrerPolicy: "no-referrer",
     body: JSON.stringify({
       username: pname,
       score: time,
     }),
  }).then((r) => {
    scores = -1;
  });
}

const oceanScale = 60;
const oceanWH = 16;
let shadowAngle = 4;
function desertOcean() {
  for (let y = (-oceanWH/2); y < oceanWH/2; y++) {
    bg.beginShape(TRIANGLE_STRIP);
    for (let x = (-oceanWH/2); x <= oceanWH/2; x++) {
      let z1 = getZ(x * oceanScale, y * oceanScale);
      let z2 = getZ(x * oceanScale, (y + 1) * oceanScale);
      bg.vertex(x * oceanScale, y * oceanScale, z1);
      bg.vertex(x * oceanScale, (y + 1) * oceanScale, z2);
    }
    bg.endShape();
  }
}

function getZ(x, y) {
  let time = 0;
  if (screen === 0) {
      time = millis() - startMillis;
    } else {
      time = startMillis;
    }
  let zVal = (noise(x/50.0, y/50.0, time/8000.0) * 200.0) - 100.0;
  return zVal;
}
