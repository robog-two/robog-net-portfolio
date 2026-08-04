const scl = 10;
let wingRot = 0;

const bpm = 118 * 4;

let music = undefined;
let timeAdjust = bpm / 60 / 1000;
let timeAdjustSine = timeAdjust * 3.14159265359;

let parsleyTex;
let sunTex;
let top1Tex;
let top2Tex;

const cityLength = 20;
let city = [];

const boidCount = 100;
let boids = [];

let cityTravelDist = 0;

function preload() {
  parsleyTex = loadImage("parsley.png");
  sunTex = loadImage("sunsphere2.jpg");
  top1Tex = loadImage("towers/tower_top.jpg");
  top2Tex = loadImage("towers/tower_top_2.jpg");
}

function randomizeHeight() {
  if (random() > 0.6) {
    return random() ** 2 * 26 * scl + 4 * scl;
  } else {
    return -1;
  }
}

function setup() {
  createCanvas(windowWidth, windowHeight, WEBGL);

  for (let row = 0; row < cityLength; row++) {
    city.push([]);
    for (let col = 0; col < 10; col++) {
      city[row][col] = randomizeHeight();
    }
  }

  for (let i = 0; i < boidCount; i++) {
    boids.push({});
    boids[i].pos = createVector(
      (random(20) - 10) * scl,
      (random(20) - 30) * scl,
      random(20 * scl),
    );
    boids[i].vel = createVector(random(-1, 1), random(-1, 1), random(-1, 1))
      .normalize();
  }
}

function draw() {
  let parsleyHeight = (sin(millis() / 10000) / 2 + 0.5) * scl * -10;
  camera(
    0,
    -10 * scl + parsleyHeight,
    -10 * scl,
    0,
    -8 * scl + parsleyHeight,
    0,
    0,
    1,
    0,
  );
  perspective(2 * atan(height / 2 / 800), width / height, 0.01, 220 * scl);
  //orbitControl();

  background(255);

  // Skybox
  push();
  rotateY(3.14);
  noStroke();
  texture(sunTex);
  sphere(200 * scl);
  pop();

  imageLight(sunTex);
  directionalLight(color("#EBCA80"), 0, 0.8, 0.5);

  // Have to translate there & back to keep lights global
  let buildingScl = 7;
  let blockScl = 8;

  rotateX(HALF_PI);
  noStroke();
  fill(color("#090620"));
  plane(300 * scl);
  ambientMaterial(255);
  rotateX(-HALF_PI);

  cityTravelDist += 1;
  if (cityTravelDist > blockScl * scl) {
    cityTravelDist = 0;
    city.shift();
    city.push([]);
    for (let col = 0; col < city[0].length; col++) {
      city[city.length - 1][col] = randomizeHeight();
    }
  }

  translate(0, 0, -cityTravelDist);

  translate((city[0].length / -2) * scl * blockScl, 0, -2 * scl * blockScl);
  for (let row = 0; row < city.length; row++) {
    for (let col = 0; col < city[row].length; col++) {
      if (city[row][col] != -1) {
        translate(scl * blockScl * col, 0, scl * blockScl * row);
        translate((blockScl / 2) * scl, 0, (blockScl / 2) * scl);
        renderTower(scl * buildingScl, city[row][col]);
        translate(-scl * blockScl * col, 0, -scl * blockScl * row);
        translate((blockScl / -2) * scl, 0, (blockScl / -2) * scl);
      }
    }
  }
  translate((city[0].length / 2) * scl * blockScl, 0, 2 * scl * blockScl);

  translate(0, 0, cityTravelDist);

  push();
  translate(0, -12 * scl + parsleyHeight, 0);
  rotateX(cos(millis() / 300) * 0.05);
  rotateY(sin(millis() / 400) * 0.05);
  rotateZ(cos(millis() / 500) * 0.05);
  translate(0, 4 * scl, scl / -2);
  renderParsley();
  pop();

  updateBoids(boids, 200 * scl);
  for (let boid of boids) {
    push();
    translate(boid.pos.x, boid.pos.y, boid.pos.z);
    stroke(255, 0, 255);
    rotateY(createVector(boid.vel.x, boid.vel.z).heading() - HALF_PI);
    renderParsley();
    pop();
  }

  if (music && music.isLoaded() && !music.isPlaying()) {
    try {
      music.play();
    } catch (_) {
      //ignore
    }
  }
}

function mousePressed() {
  music = music ?? loadSound("chive.mp3");
}
