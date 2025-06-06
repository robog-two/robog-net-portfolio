let peeps = [];
let rateTotals = 0;

const friction = 0.5; // Multiply velocity by this each frame (1 = frictionless, 0 = no velocity)
const dominantProp = 1; // Proportion of circular flow (and 1-prop is how much perlin noise)
const differentProp = 0.8; // The proportion for the special peep
let inward = 0.8; // How much to rotate inward direction vector for circular flow
let magn = 0.5; // Magnitude of movement vectors (affects the acceleration)
let holeRadius = 0.06; // How big is the sun in the center as a proportion of the area of the screen

let bgR = 135;
let bgG = 206;
let bgB = 235;

// Shepard tone configuration
const bottomF = 40;
const topF = bottomF * 4;
let changeSpeed = 0.999; // Percent per frame
let oscillators = [];
let hasStartedAudio = false;

// Configure frequencies in noise
const hi = 0.1;
const med = 0.2;
const low = 0.3;
const ultralow = 0.4;

// Clip noise between values to increase variation
const noiseMin = 0.4;
const noiseMax = 0.6;

function myNoise(x, y) {
  let noiseAt = map(noise(x * 2, y * 2), 0, 1, 0, hi) +
    map(noise(x, y), 0, 1, 0, med) +
    map(noise(x / 2, y / 2), 0, 1, 0, low) +
    map(noise(x / 4, y / 4), 0, 1, 0, ultralow);
  return Math.min(Math.max(noiseAt, noiseMin), noiseMax);
}

function preload() {
  sound = loadSound("birds.mp3");
}

function setup() {
  createCanvas(windowWidth - 100, windowHeight - 300);

  holeRadius = sqrt((holeRadius * width * height) / PI);

  for (let i = 0; i < 10000; i++) {
    let dir = random(-PI, PI);
    let vel = random(holeRadius, Math.max(height / 2, width / 2) * 1.414);
    let x = (cos(dir) * vel) + (width / 2);
    let y = (sin(dir) * vel) + (height / 2);
    peeps.push({ x, y, vx: 0, vy: 0, c: 0 });
  }

  //peeps[peeps.length - 1].c = 1;

  frameRate(30);
  background(bgR, bgG, bgB);
}

function draw() {
  for (let i = 0; i < oscillators.length; i++) {
    let currentFreq = oscillators[i].getFreq() * changeSpeed;
    if (currentFreq < bottomF || currentFreq > topF) {
      oscillators[i].disconnect();
      oscillators[i].stop();
      if (random() > 0.5) {
        oscillators[i] = new p5.SawOsc();
      } else {
        oscillators[i] = new p5.TriOsc();
      }
      oscillators[i].pan(random(-0.5, 0.5));
      oscillators[i].amp(0);
      oscillators[i].start();
      if (currentFreq < bottomF) {
        currentFreq = topF;
      } else {
        currentFreq = bottomF;
      }
    }
    let currentAmp = 1 -
      Math.abs(
        (2 / (topF - bottomF)) *
          ((currentFreq - bottomF) - (0.5 * (topF - bottomF))),
      ); // 1 - |(1/5)(x-5)|
    oscillators[i].freq(currentFreq);
    oscillators[i].amp(currentAmp);
  }

  rateTotals += Math.round(frameRate());
  if (frameCount % 8 == 0) {
    if (rateTotals / 8 > 29) {
      for (let i = 0; i < 100; i++) {
        let dir = random(-PI, PI);
        let vel = Math.max(height / 2, width / 2) * 1.414;
        vel += random(40);
        let x = (cos(dir) * vel) + (width / 2);
        let y = (sin(dir) * vel) + (height / 2);
        peeps.push({ x, y, vx: 0, vy: 0, c: 0 });
      }
      console.log("Increasing! FPS is OK. New count:", peeps.length);
    }
    if (rateTotals / 8 < 25) {
      peeps = peeps.slice(200);
      console.log(
        "Reducing!",
        Math.round(rateTotals / 8),
        "fps too low. New count:",
        peeps.length,
      );
    }
    rateTotals = 0;
  }
  background(bgR, bgG, bgB, 20);

  for (let peep of peeps) {
    if (peep.c == 1) {
      stroke(0);
      strokeWeight(3);
    } else {
      stroke(230);
      strokeWeight(1);
    }
    let x = peep.x;
    let y = peep.y;
    let flowAngle = flow(
      peep.x,
      peep.y,
      peep.c == 1 ? differentProp : dominantProp,
    ).mult(magn);
    peep.vx = (peep.vx * friction) + flowAngle.x;
    peep.vy = (peep.vy * friction) + flowAngle.y;
    x += peep.vx;
    y += peep.vy;
    let cdist = dist(peep.x, peep.y, width / 2, height / 2);
    if (cdist < holeRadius) {
      let dir = random(-PI, PI);
      let vel = Math.max(height / 2, width / 2) * 1.414;
      x = (cos(dir) * vel) + (width / 2);
      y = (sin(dir) * vel) + (height / 2);
      vx = 0;
      vy = 0;
    } else {
      strokeWeight(
        map(cdist, holeRadius, Math.max(width / 2, height / 2), 0, 4),
      );
      line(x, y, peep.x, peep.y);
    }
    peep.x = x;
    peep.y = y;
  }

  noStroke();
  fill(252, 255, 0);
  ellipse(
    width / 2,
    height / 2,
    (holeRadius * 2) * (1.1 + (sin(millis() / 5000) * 0.1)),
  );
}

function flow(x, y, prop) {
  return (
    createVector(-(x - (width / 2)), -(y - height / 2)).normalize().rotate(
      inward,
    )
  ).mult(prop).add((
    createVector(1, 0).rotate(map(
      myNoise(x * 0.08, y * 0.08),
      noiseMin,
      noiseMax,
      -PI,
      PI,
    ))
  ).mult(1 - prop));
}
let filtereq;
function mousePressed() {
  if (!hasStartedAudio) {
    sound.loop();
    hasStartedAudio = true;
    filtereq = new p5.Distortion();
    filtereq.set(0.2);

    let currentFreq = bottomF;
    let i = 0;
    while (currentFreq <= topF) {
      let currentAmp = 1 -
        Math.abs(
          (2 / (topF - bottomF)) *
            ((currentFreq - bottomF) - (0.5 * (topF - bottomF))),
        ); // 1 - |(1/5)(x-5)|
      if (random() > 0.5) {
        oscillators.push(new p5.SawOsc());
      } else {
        oscillators.push(new p5.TriOsc());
      }
      oscillators[i].pan(random(-0.5, 0.5));
      oscillators[i].freq(currentFreq);
      oscillators[i].amp(currentAmp);
      oscillators[i].start();
      i++;
      currentFreq = bottomF * (i + 1);
    }
  }
}
