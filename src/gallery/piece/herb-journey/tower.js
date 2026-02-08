// North = +z
// East = +x
// South = -z
// West = -x

towerSizes = [
  // Width, height, w/h * 100
  [46, 107, 43], // aspect .43
  [49, 175, 28], // aspect .28
  [63, 39, 162], // aspect 1.62
  [68, 109, 62], // aspect 0.62
  [73, 68, 107], // aspect 1.07
]

towerCache = {};

function loadTowerTexture(tWidth, tHeight) {
  // TEXTURE MANAGEMENT
  let aspect = round((tWidth/tHeight) * 100);
  if (towerCache[towerCache[aspect]]) {
    if (towerCache[towerCache[aspect]].isLoaded) {
      texture(towerCache[towerCache[aspect]]);
    } else {
      fill(0);
    }
  } else {
    let sizeArr = towerSizes.reduce((acc, next) => {
      if (abs(next[2] - aspect) < abs(acc[2] - aspect)) {
        return next;
      } else {
        return acc;
      }
    });
    let resStr = "towers/tower" + sizeArr[0] + "x" + sizeArr[1] + ".jpg";
    
    towerCache[aspect] = resStr;
    if (towerCache[resStr] == undefined) {
      towerCache[resStr] = loadImage(
        resStr,
        () => {
          towerCache[resStr].isLoaded = true;
          console.log("INFORMATION [i] Tower " + sizeArr[0] + "x" + sizeArr[1] + " now loaded!");
        }
      );
      fill(0);
    } else {
      texture(towerCache[resStr]);
    }
  }
}

function renderTower(tWidth, tHeight) {
  push();
  if (round(tHeight) % 21 == 0 && millis() % 500 > 400) {
    pointLight(color("#882C09"), 0, tHeight, 0); // twinkling light
  }
  
  let tRadius = tWidth / 2;
  
  noStroke();
  
  textureMode(NORMAL);
  loadTowerTexture(tWidth, tHeight);
  
  beginGeometry(); // begin tower model
  beginShape(TRIANGLE_STRIP);
  
  // North side (+z)
  //vertex(-tRadius, 0, tRadius, 0, 1);
  //vertex(-tRadius, -tHeight, tRadius, 0, 0);
  vertex(tRadius, 0, tRadius, 1, 1);
  vertex(tRadius, -tHeight, tRadius, 1, 0);
  
  // East side (+x)
  vertex(tRadius, 0, -tRadius, 0, 1);
  vertex(tRadius, -tHeight, -tRadius, 0, 0);
  
  // South side (-z)
  vertex(-tRadius, 0, -tRadius, 1, 1);
  vertex(-tRadius, -tHeight, -tRadius, 1, 0);
  
  // West side (-x)
  vertex(-tRadius, 0, tRadius, 0, 1);
  vertex(-tRadius, -tHeight, tRadius, 0, 0);
  
  
  // Calculate normals & display tower model
  endShape();
  let towerModel = endGeometry();
  towerModel.computeNormals();
  model(towerModel);
  
  if (round(tHeight) % 2 == 0) {
    texture(top1Tex);
  } else {
    texture(top2Tex);
  }
  
  beginGeometry(); // begin roof model
  
  // Tower roof
  beginShape(TRIANGLE_STRIP);
  vertex(-tRadius, -tHeight, -tRadius, 0, 0);
  vertex(tRadius, -tHeight, -tRadius, 1, 0);
  vertex(-tRadius, -tHeight, tRadius, 0, 1);
  vertex(tRadius, -tHeight, tRadius, 1, 1);
  endShape();
  
  // Calculate normals and display roof model
  let towerTopModel = endGeometry();
  towerTopModel.computeNormals();
  model(towerTopModel);
  
  
  pop();
}

function mousePressed() {
  fullscreen(!fullscreen());
  setTimeout(() => {
    resizeCanvas(windowWidth, windowHeight);
  }, 1000);
}