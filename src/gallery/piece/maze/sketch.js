const mazeRows = 20;
const mazeCols = 20;

const margin = 3;
const twoMargin = margin * 2;
let squareWidth;
let squareHeight;

let seenSquares = new PriorityQueue();

let topWalls = twoDArray(mazeRows + 1, mazeCols, true);
let leftWalls = twoDArray(mazeRows, mazeCols + 1, true);
let rooms = twoDArray(mazeRows, mazeCols, 0); // 0 = enclosed, 1 = enclosed but adjacent to an open room, 2 = a wall is open (this room is now fixed)

let generating = true;

function setup() {
  createCanvas(400, 400);

  squareWidth = (width - twoMargin) / mazeCols;
  squareHeight = (height - twoMargin) / mazeRows;

  // Initial conditions ================
  // The top left is locked in place
  rooms[0][0] = 2;

  // Around the top left we can explore
  rooms[0][1] = 1;
  seenSquares.push(floor(random(1000)), [0, 1]);
  rooms[1][0] = 1;
  seenSquares.push(floor(random(1000)), [1, 0]);

  // The start/end of the maze should be open
  destroyWall(0, 0, "left");
  destroyWall(mazeRows - 1, mazeCols - 1, "right");
  //frameRate(999999);
}

function draw() {
  if (generating) {
    background(color("#C1E4EE"));
  } else {
    if (mouseIsPressed) {
      stroke(color("#F26419"));
      strokeWeight(4);
      line(mouseX, mouseY, pmouseX, pmouseY);
    }
  }

  // Display the walls of the maze
  for (let row = 0; row <= mazeRows; row++) {
    for (let col = 0; col <= mazeCols; col++) {
      if (col != mazeCols && row != mazeRows) {
        if (rooms[row][col] != 2) {
          noStroke();
          if (rooms[row][col] == 0) fill(color("#81C3D7"));
          if (rooms[row][col] == 1) fill(color("#AE2012"));
          rect(
            map(col, 0, mazeCols, margin, width - margin),
            map(row, 0, mazeRows, margin, height - margin),
            squareWidth,
            squareHeight
          );
        }
      }

      stroke(color("#16425B"));
      strokeWeight(generating ? 3 : 1);
      // Top wall
      if (col != mazeCols && topWalls[row][col]) {
        let leftX = map(col, 0, mazeCols, margin, width - margin);
        let rightX = leftX + squareWidth;
        let y = map(row, 0, mazeRows, margin, height - margin);
        line(leftX, y, rightX, y);
      }
      if (row != mazeRows && leftWalls[row][col]) {
        let x = map(col, 0, mazeCols, margin, width - margin);
        let topY = map(row, 0, mazeRows, margin, height - margin);
        let bottomY = topY + squareHeight;
        line(x, topY, x, bottomY);
      }
    }
  }

  if (generating) {
    // maze generating algorithm
    if (seenSquares.length != 0) {
      let currentSq = seenSquares.pop();
      let currentRow = currentSq[0];
      let currentCol = currentSq[1];

      let options = [];
      udlr(currentRow, currentCol, (row, col, direction) => {
        if ((rooms[row] ?? [])[col] == 2) {
          options.push(direction);
        }
      });

      let choice = options[floor(random(0, options.length))];
      destroyWall(currentRow, currentCol, choice);

      rooms[currentRow][currentCol] = 2; // lock this room from editing

      udlr(currentRow, currentCol, (row, col, direction) => {
        if ((rooms[row] ?? [])[col] == 0) {
          seenSquares.push(floor(random(1000)), [row, col]);
          rooms[row][col] = 1;
        }
      });
    } else {
      generating = false;
      console.log(
        `Took ${frameCount} frames to generate maze, there are ${
          mazeRows * mazeCols
        } squares in this maze.`
      );
    }
  }
}

function destroyWall(row, col, direction) {
  switch (direction) {
    case "up": {
      topWalls[row][col] = false;
      break;
    }
    case "down": {
      topWalls[row + 1][col] = false;
      break;
    }
    case "left": {
      leftWalls[row][col] = false;
      break;
    }
    case "right": {
      leftWalls[row][col + 1] = false;
      break;
    }
    default: {
      throw new Error("Direction must be up/down/left/right");
    }
  }
}

function udlr(row, col, toExecute) {
  toExecute(row - 1, col, "up");
  toExecute(row + 1, col, "down");
  toExecute(row, col - 1, "left");
  toExecute(row, col + 1, "right");
}

function twoDArray(rows, cols, fillWith = undefined) {
  return Array.from({ length: rows }, () =>
    Array.from({ length: cols }, () => fillWith)
  );
}

function keyPressed() {
  if (key == " ") {
    redrawMaze();
  }
}

function redrawMaze() {
  // re-renders the maze if it has finished generating, effectively clearing the screen
  generating = true;
}
