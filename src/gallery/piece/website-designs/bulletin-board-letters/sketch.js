let letters = [];
let letterChars = "samknight";
let offsets = [
  // offsetX, offsetY, centerX, centerY, inertia data
  [0.2, 15, 16, 39],
  [34.4, 17, 51, 38],
  [70, 17, 97, 39],
  [136.5, 0, 150, 34],
  [168, 17, 186, 38],
  [202, 1, 208, 34],
  [215.4, 16, 233, 37],
  [252, 0, 268, 36],
  [284, 3, 296, 25],
];

let trace;

const cursorRepel = 80; // radius around cursor
const cursorForce = 1.5; // force to repel
const returnFrac = 0.003; // % attraction of original position vs force of cursor
const drag = 0.92;

// Adjust to place correctly on the page
let scaler = 1;

let globalOffset = [200, 50];

function setup() {
  noCanvas();
  for (const letter of letterChars) {
    const letterEl = createImg(
      "letter-" + letter + ".svg",
      letter.toUpperCase(),
    );
    letterEl.elt.style.transformOrigin = "top left";
    letters.push(letterEl);
  }
}

function draw() {
  scaler = (windowWidth * 0.7) / 312;
  globalOffset[0] = windowWidth * 0.15;

  for (let i = 0; i < letters.length; i++) {
    letters[i].elt.style.transform = "scale(" + scaler + ")";

    let x = offsets[i][0] * scaler;
    let y = offsets[i][1] * scaler;
    let cx = offsets[i][2] * scaler;
    let cy = offsets[i][3] * scaler;

    x += globalOffset[0];
    y += globalOffset[1];
    cx += globalOffset[0];
    cy += globalOffset[1];

    let inertiaData = offsets[i][4];
    if (inertiaData == undefined) {
      inertiaData = {
        x: x + random(300, -300),
        y: y + random(-50, -300),
        velX: 0,
        velY: 0,
      };
      offsets[i].push(inertiaData);
    }

    if (
      dist(mouseX, mouseY, inertiaData.x + (cx - x), inertiaData.y + (cy - y)) <
      cursorRepel * scaler
    ) {
      let vec = createVector(
        mouseX - (inertiaData.x + (cx - x)),
        mouseY - (inertiaData.y + (cy - y))
      );
      vec.normalize();
      vec.mult(-1 * cursorForce);
      inertiaData.velX += vec.x;
      inertiaData.velY += vec.y;
    }
    
    inertiaData.velX = (inertiaData.velX * (1-returnFrac)) + ((x - inertiaData.x) * returnFrac);
    inertiaData.velY = (inertiaData.velY * (1-returnFrac)) + ((y - inertiaData.y) * returnFrac);

    inertiaData.x += inertiaData.velX;
    inertiaData.y += inertiaData.velY;

    inertiaData.velX *= drag;
    inertiaData.velY *= drag;

    offsets[i][4] = inertiaData;
    letters[i].position(inertiaData.x, inertiaData.y);
  }
}
