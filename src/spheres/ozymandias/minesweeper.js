// Daniel Shiffman
// The Coding Train
// Minesweeper
// https://thecodingtrain.com/challenges/71-minesweeper
// Video: https://youtu.be/LFU5ZlrR21E

// ES6 version: https://editor.p5js.org/codingtrain/sketches/Xap-KQuO_

// Daniel Shiffman
// The Coding Train
// Minesweeper
// https://thecodingtrain.com/challenges/71-minesweeper
// Video: https://youtu.be/LFU5ZlrR21E

function Cell(i, j, w) {
  this.i = i;
  this.j = j;
  this.x = i * w + window.minesweeper.gXOff;
  this.y = j * w + window.minesweeper.gYOff;
  this.w = w;
  this.neighborCount = 0;

  this.bee = false;
  this.revealed = false;
}

Cell.prototype.show = function() {
  noStroke();
  if (this.revealed) {
    if (this.bee) {
      noFill();
      stroke(255);
      rect(this.x, this.y, this.w, this.w);
      ellipse(this.x + this.w * 0.5, this.y + this.w * 0.5, this.w * 0.5);
    } else {
      noFill();
      stroke(255);
      rect(this.x, this.y, this.w, this.w);
      rect(this.x, this.y, this.w, this.w);
      if (this.neighborCount > 0) {
        textAlign(CENTER);
        fill(255);
        text(this.neighborCount, this.x + this.w * 0.5, this.y + this.w - 6);
      }
    }
  } else {
    stroke("#8795db");
    fill("#8795db");
    rect(this.x, this.y, this.w, this.w);
  }
}

Cell.prototype.countBees = function() {
  if (this.bee) {
    this.neighborCount = -1;
    return;
  }
  let total = 0;
  for (let xoff = -1; xoff <= 1; xoff++) {
    let i = this.i + xoff;
    if (i < 0 || i >= window.minesweeper.cols) continue;

    for (let yoff = -1; yoff <= 1; yoff++) {
      let j = this.j + yoff;
      if (j < 0 || j >= window.minesweeper.rows) continue;

      let neighbor = window.minesweeper.grid[i][j];
      if (neighbor.bee) {
        total++;
      }
    }
  }
  this.neighborCount = total;
}

Cell.prototype.contains = function(x, y) {
  return (x > this.x && x < this.x + this.w && y > this.y && y < this.y + this.w);
}

Cell.prototype.reveal = function() {
  this.revealed = true;
  if (this.neighborCount == 0) {
    // flood fill time
    this.floodFill();
  }
}

Cell.prototype.floodFill = function() {
  for (let xoff = -1; xoff <= 1; xoff++) {
    let i = this.i + xoff;
    if (i < 0 || i >= window.minesweeper.cols) continue;

    for (let yoff = -1; yoff <= 1; yoff++) {
      let j = this.j + yoff;
      if (j < 0 || j >= window.minesweeper.rows) continue;

      let neighbor = window.minesweeper.grid[i][j];
      // Note the neighbor.bee check was not required.
      // See issue #184
      if (!neighbor.revealed) {
        neighbor.reveal();
      }
    }
  }
}

// MAIN MINESWEEPER CODE ==========================================

window.minesweeper = {}

window.minesweeper.make2DArray = (cols, rows) => {
  let arr = new Array(cols);
  for (let i = 0; i < arr.length; i++) {
    arr[i] = new Array(rows);
  }
  return arr;
}

window.minesweeper.totalBees = 30;

window.minesweeper.setupSweep = (sweepSideLen, x, y) => {
  do {
    window.minesweeper.gXOff = x;
    window.minesweeper.gYOff = y;
    window.minesweeper.w = sweepSideLen/20;
    window.minesweeper.cols = 20;
    window.minesweeper.rows = 20;
    window.minesweeper.grid = window.minesweeper.make2DArray(window.minesweeper.cols, window.minesweeper.rows);
    for (let i = 0; i < window.minesweeper.cols; i++) {
      for (let j = 0; j < window.minesweeper.rows; j++) {
        window.minesweeper.grid[i][j] = new Cell(i, j, window.minesweeper.w);
      }
    }

    // Pick totalBees spots
    let options = [];
    for (let i = 0; i < window.minesweeper.cols; i++) {
      for (let j = 0; j < window.minesweeper.rows; j++) {
        if (!(i === 10 && j === 10)) {
          options.push([i, j]);
        }
      }
    }


    for (let n = 0; n < window.minesweeper.totalBees; n++) {
      let index = floor(random(options.length));
      let choice = options[index];
      let i = choice[0];
      let j = choice[1];
      // Deletes that spot so it's no longer an option
      options.splice(index, 1);
      window.minesweeper.grid[i][j].bee = true;
    }

    for (let i = 0; i < window.minesweeper.cols; i++) {
      for (let j = 0; j < window.minesweeper.rows; j++) {
        window.minesweeper.grid[i][j].countBees();
      }
    }

    // Reveal the center
    window.minesweeper.grid[10][10].reveal();


    for (let i = 0; i < window.minesweeper.cols; i++) {
      for (let j = 0; j < window.minesweeper.rows; j++) {
        window.minesweeper.grid[i][j].countBees();
      }
    }
  } while (window.minesweeper.grid[10][10].neighborCount !== 0)
}

window.minesweeper.gameOver = () => {
  for (let i = 0; i < window.minesweeper.cols; i++) {
    for (let j = 0; j < window.minesweeper.rows; j++) {
      window.minesweeper.grid[i][j].revealed = true;
    }
  }
  
  window.onGameOver()
}

window.minesweeper.win = () => {
  for (let i = 0; i < window.minesweeper.cols; i++) {
    for (let j = 0; j < window.minesweeper.rows; j++) {
      window.minesweeper.grid[i][j].revealed = true;
    }
  }
  
  window.onWin()
}

window.minesweeper.mousePressedSweep = () => {
  for (let i = 0; i < window.minesweeper.cols; i++) {
    for (let j = 0; j < window.minesweeper.rows; j++) {
      if (window.minesweeper.grid[i][j].contains(mouseX, mouseY)) {
        window.minesweeper.grid[i][j].reveal();

        if (window.minesweeper.grid[i][j].bee) {
          window.minesweeper.gameOver();
          return;
        }
      }
    }
  }
  
  let hasWon = true;
  for (let i = 0; i < window.minesweeper.cols; i++) {
    for (let j = 0; j < window.minesweeper.rows; j++) {
      if (!(window.minesweeper.grid[i][j].revealed || window.minesweeper.grid[i][j].bee)) {
        hasWon = false;
      }
    }
  }
  if (hasWon) {
    window.minesweeper.win();
  }
}

window.minesweeper.drawSweep = () => {
  for (let i = 0; i < window.minesweeper.cols; i++) {
    for (let j = 0; j < window.minesweeper.rows; j++) {
      window.minesweeper.grid[i][j].show();
    }
  }
}