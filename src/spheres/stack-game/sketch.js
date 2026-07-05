// The size / pos of the current block
let left, right, back, front;

let slomo = 7;

// Is it coming from the front or right side?
let comingFromFront = true;

let lost = false;

let leniency = 10;
let clack = 0; // Countdown for perfect block effect

// (displayHeight*dpi) of the blocks & where we're at
let y = 0;
let blockHeight = 75;
let dpi = 1;

// All the blocks in the stack
let blocks = [];
let targetColor;
let currentColor;

let fallingBlocks = [undefined, undefined, undefined];

let lastBG;

let isreset = false;
let roboto;
let full = false;
let hasResetOnce = false;

let fallingArrWindow = 0;

let hitSound;
let perfectSound;
let slidingSound;

let musicSound;

function addFallingBlock(block) {
  fallingBlocks[fallingArrWindow] = block;
  fallingArrWindow = int((fallingArrWindow + 1) % 3);
}

class Block {
  constructor(left, right, front, back, c, y) {
    this.left = left;
    this.right = right;
    this.front = front;
    this.back = back;
    this.c = c;
    this.y = y;
  }
}

class FallingBlock extends Block {
  constructor(left, right, front, back, c, y, fromFront, over) {
    super(left, right, front, back, c, y);
    this.vel = createVector(0, 0, 0);
    this.rotVel = createVector(0, 0, 0);
    this.rot = createVector(0, 0, 0);
    this.fromFront = fromFront;
    this.over = over;
  }

  update() {
    if (lost && !isreset) {
      this.left += this.vel.x / slomo;
      this.right += this.vel.x / slomo;
      this.front += this.vel.z / slomo;
      this.back += this.vel.z / slomo;
      this.y += this.vel.y / slomo;
      this.vel.y -= 2 / slomo;
      this.rot.add(this.rotVel.copy().mult(1 / slomo));
    } else {
      this.left += this.vel.x;
      this.right += this.vel.x;
      this.front += this.vel.z;
      this.back += this.vel.z;
      this.y += this.vel.y;
      this.vel.y -= 2;
      this.rot.add(this.rotVel);
    }
  }
}

function preload() {
  roboto = loadFont(
    "https://unpkg.com/roboto-regular-woff@0.7.1/Roboto-Regular.woff"
  );
  
  hitSound = loadSound("chop.mp3");
  perfectSound = loadSound("spoton.mp3");
  slidingSound = loadSound("slide.mp3");
  slidingSound.setVolume(0.3);

  musicSound = loadSound("music.mp3");
}

function setup() {
  createCanvas(windowWidth, windowHeight, WEBGL);
  frameRate(60);
  textFont(roboto);
  lost = true;
  lastBG = [0, 0, 0];
  currentColor = [0, 0, 0];
  targetColor = [0, 0, 0];
}

function draw() {
  if (!full) {
    background(0);
    noStroke();
    fill(255);
    textAlign(LEFT, TOP);
    textSize(60);
    text("Tap to enter fullscreen.", 0, 0, width / 2, height / 2);
    return;
  }
  if (lost && hasResetOnce) {
    let curr =
      blocks[int(constrain(floor(y / blockHeight), 0, blocks.length - 1))].c;
    lastBG = color(255 - red(curr), 255 - green(curr), 255 - blue(curr));
  } else {
    lastBG = color(
      255 - red(currentColor),
      255 - green(currentColor),
      255 - blue(currentColor)
    );
  }
  background(lastBG);

  camera(
    -(displayWidth * dpi),
    (displayHeight * dpi) / 2 + y,
    -(displayWidth * dpi),
    0,
    y,
    0,
    0,
    -1,
    0
  );
  pointLight(
    255,
    255,
    255,
    -(displayWidth * dpi),
    blockHeight * 10 + y,
    -(displayWidth * dpi) / 2
  );
  pointLight(
    255,
    255,
    255,
    -(displayWidth * dpi) / 4,
    blockHeight * -5 + y,
    -(displayWidth * dpi)
  );
  pointLight(255, 255, 255, 0, blockHeight * 30 + y, 0);

  if (lost || !isreset) {
    if (y <= blockHeight / 7) {
      if (!isreset) {
        y = 0;
        blocks = [];
        left = -(displayWidth * dpi) / 6;
        right = (displayWidth * dpi) / 6;
        front = (-(displayWidth * dpi) * 2) / 3;
        back = -(displayWidth * dpi) / 3;

        targetColor = color(
          round(random(50, 200)),
          round(random(50, 200)),
          round(random(50, 200))
        );
        currentColor = targetColor;
        lastBG = color(255 - currentColor);
        blocks.push(
          new Block(
            -(displayWidth * dpi) / 6,
            (displayWidth * dpi) / 6,
            -(displayWidth * dpi) / 6,
            (displayWidth * dpi) / 6,
            currentColor,
            -blockHeight
          )
        );

        comingFromFront = true;
        isreset = true;
        musicSound.stop();
        hasResetOnce = true;
        return;
      }
    } else {
      y = fallingBlocks[(fallingArrWindow + 2) % 3].y;
    }
  } else if (round(y) % blockHeight != 0) {
    y += blockHeight / 7;
  } else {
    let ratebutnotp5 = 10 + floor(y / blockHeight / 25);
    leniency = ratebutnotp5 * 2;
    if (comingFromFront) {
      front += ratebutnotp5;
      back += ratebutnotp5;
    } else {
      left += ratebutnotp5;
      right += ratebutnotp5;
    }
  }

  if (round(hue(currentColor) / 5) == round(hue(targetColor) / 5)) {
    targetColor = color(
      round(random(50, 200)),
      round(random(50, 200)),
      round(random(50, 200))
    );
    adjustColor();
  }

  if (!lost) {
    push();
    fill(currentColor);
    noStroke();
    translate((back + front) / 2, y + blockHeight / 2, (left + right) / 2);
    box(front - back, blockHeight, right - left);
    pop();
  }

  if (clack > 0) {
    let c = blocks[blocks.length - 1];
    push();
    fill(255);
    noStroke();
    translate((c.back + c.front) / 2, c.y, (c.left + c.right) / 2);
    box(c.front - c.back - clack, 0.1, c.right - c.left + clack);
    pop();
    clack--;
  }

  // Render the block tower
  for (
    let i = min(int(ceil(y - (height*2) / blockHeight)),0);
    i <= ((y + (height*3)) / blockHeight);
    i++
  ) {
    if (i >= 0 && i < blocks.length) {
      let b = blocks[i];
      push();
      fill(b.c);
      noStroke();
      translate(
        (b.back + b.front) / 2,
        b.y + blockHeight / 2,
        (b.left + b.right) / 2
      );
      box(b.front - b.back, blockHeight, b.right - b.left);
      pop();
    }
  }

  // Simulate the blocks that are falling with gravity
  for (let b of fallingBlocks) {
    if (b !== undefined) {
      push();
      fill(b.c);
      noStroke();
      translate(
        (b.back + b.front) / 2,
        b.y + blockHeight / 2,
        (b.left + b.right) / 2
      );
      push();
      rotateX(b.rot.x);
      rotateY(b.rot.y);
      rotateZ(b.rot.z);
      box(b.front - b.back, blockHeight, b.right - b.left);
      pop();
      pop();
      b.update();
      for (let b2 of blocks) {
        if (b2.y + blockHeight >= b.y && b.y >= b2.y) {
          if (
            (b.fromFront && b.front < b2.back && b.back > b2.front) ||
            (!b.fromFront && b.left < b2.right && b.right > b2.left)
          ) {
            if (b.fromFront) {
              if (b.over) {
                if (!hitSound.isPlaying()) hitSound.play();
                b.rotVel.z -= PI / (lost ? 90 : 60);
                b.vel.z += lost ? 1 : 6;
              } else {
                if (!hitSound.isPlaying()) hitSound.play();
                b.rotVel.z += PI / (lost ? 90 : 60);
                b.vel.z -= lost ? 1 : 6;
              }
            } else {
              if (b.over) {
                if (!hitSound.isPlaying()) hitSound.play();
                b.rotVel.x += PI / (lost ? 90 : 60);
                b.vel.x += lost ? 1 : 6;
              } else {
                if (!hitSound.isPlaying()) hitSound.play();
                b.rotVel.x -= PI / (lost ? 90 : 60);
                b.vel.x -= lost ? 1 : 6;
              }
            }
            b.vel.y *= -0.5;
            b.vel.y += 10;
            b.y += 2;
          }
        }
      }
    }
  }

  noStroke();
  textSize(100);
  textAlign(LEFT, TOP);
  push();
  let b =
    blocks[int(constrain(floor(y / blockHeight) + 2, 0, blocks.length - 1))];
  if (brightness(b.c) > 220) {
    fill(0);
  } else {
    fill(255);
  }
  translate(b.front, y + blockHeight * 2 + 20, b.left);
  push();
  rotateZ(PI);
  rotateY(PI);
  text(
    blocks.length > 1 ? blocks.length + " blocks" : "1 block",
    0,
    0,
    displayWidth * dpi,
    displayHeight * dpi
  );
  pop();
  pop();
  
  b = blocks[blocks.length - 1];
  
  if ((comingFromFront && front < b.back && back > b.front) ||  (!comingFromFront && left < b.right && right > b.left)) {
    if (!slidingSound.isPlaying()) slidingSound.loop();
  } else {
    if (slidingSound.isPlaying()) slidingSound.stop();
  }

  if (!isreset || lost) {
    musicSound.rate(-0.9)
    noStroke();
    textSize(75);
    b =
      blocks[int(constrain(floor(y / blockHeight) - 1, 0, blocks.length - 1))];
    if (brightness(b.c) > 220) {
      fill(0);
    } else {
      fill(255);
    }
    textAlign(LEFT, TOP);
    push();
    translate(b.front, y - blockHeight - 20, b.left);
    push();
    rotateZ(PI);
    rotateY(PI);
    text("tap to start", 0, 0, displayWidth * dpi, displayHeight * dpi);
    pop();
    pop();
  }
}

function adjustColor() {
  let targetHue = red(targetColor);
  let targetSat = green(targetColor);
  let targetBgt = blue(targetColor);

  let hue = red(currentColor);
  let sat = green(currentColor);
  let bgt = blue(currentColor);

  for (let i = 0; i < 5; i++) {
    if (targetHue > hue) hue++;
    if (targetHue < hue) hue--;
    if (targetSat > sat) sat++;
    if (targetSat < sat) sat--;
    if (targetBgt > bgt) bgt++;
    if (targetBgt < bgt) bgt--;
  }

  currentColor = color(round(hue), round(sat), round(bgt));
}

function keyPressed() {
  press();
}

function mousePressed() {
  
  if (mousePressed) {
    press();
  }
}

function press() {
  if (!full) {
    userStartAudio();
    musicSound.setVolume(0.7);
    musicSound.loop();
    fullscreen(true);
    full = true;
    dpi = ceil(1000 / displayWidth);
    resizeCanvas(displayWidth, displayHeight);
  }
  
  if (lost && isreset) {
    musicSound.setVolume(0.7);
    musicSound.rate(1);
    musicSound.loop();
    lost = false;
    return;
  }
  if (lost && !isreset) {
    musicSound.setVolume(0.7);
    musicSound.rate(1);
    musicSound.loop();
    fallingBlocks = [];
    y = 0;
    lost = false;
    isreset = false;
    return;
  }
  
  if (round(y) % blockHeight == 0 && !lost) {
    let b = blocks[blocks.length - 1];
    if (
      abs(left - b.left) < leniency &&
      abs(right - b.right) < leniency &&
      abs(front - b.front) < leniency &&
      abs(back - b.back) < leniency
    ) {
      perfectSound.play();
      if (slidingSound.isPlaying()) slidingSound.stop();
      clack = 40;
      blocks.push(new Block(b.left, b.right, b.front, b.back, currentColor, y));
    } else {
      hitSound.play();
      if (slidingSound.isPlaying()) slidingSound.stop();
      if (comingFromFront) {
        if (front > b.back || back < b.front) {
          lost = true;
          isreset = false;
          if (max(front, b.front) == front) {
            addFallingBlock(
              new FallingBlock(
                left,
                right,
                front,
                back,
                currentColor,
                y,
                comingFromFront,
                true
              )
            );
          } else {
            addFallingBlock(
              new FallingBlock(
                left,
                right,
                front,
                back,
                currentColor,
                y,
                comingFromFront,
                false
              )
            );
          }
          return;
        }
      } else {
        if (left > b.right || right < b.left) {
          lost = true;
          isreset = false;
          if (max(left, b.left) == left) {
            addFallingBlock(
              new FallingBlock(
                left,
                right,
                front,
                back,
                currentColor,
                y,
                comingFromFront,
                true
              )
            );
          } else {
            addFallingBlock(
              new FallingBlock(
                left,
                right,
                front,
                back,
                currentColor,
                y,
                comingFromFront,
                false
              )
            );
          }
          return;
        }
      }

      if (!comingFromFront) {
        if (max(left, b.left) == left) {
          addFallingBlock(
            new FallingBlock(
              b.right,
              right,
              front,
              back,
              currentColor,
              y,
              comingFromFront,
              true
            )
          );
        } else {
          addFallingBlock(
            new FallingBlock(
              left,
              b.left,
              front,
              back,
              currentColor,
              y,
              comingFromFront,
              false
            )
          );
        }
      } else {
        if (max(front, b.front) == front) {
          addFallingBlock(
            new FallingBlock(
              left,
              right,
              b.back,
              back,
              currentColor,
              y,
              comingFromFront,
              true
            )
          );
        } else {
          addFallingBlock(
            new FallingBlock(
              left,
              right,
              front,
              b.front,
              currentColor,
              y,
              comingFromFront,
              false
            )
          );
        }
      }

      blocks.push(
        new Block(
          max(left, b.left),
          min(right, b.right),
          max(front, b.front),
          min(back, b.back),
          currentColor,
          y
        )
      );
    }

    comingFromFront = !comingFromFront;
    let b2 = blocks[blocks.length - 1];

    if (comingFromFront) {
      front = b2.front - (displayWidth * dpi) / 3;
      back = b2.back - (displayWidth * dpi) / 3;
      left = b2.left;
      right = b2.right;
    } else {
      front = b2.front;
      back = b2.back;
      left = b2.left - (displayWidth * dpi) / 3;
      right = b2.right - (displayWidth * dpi) / 3;
    }

    adjustColor();
    y += blockHeight / 7;
  }
}
