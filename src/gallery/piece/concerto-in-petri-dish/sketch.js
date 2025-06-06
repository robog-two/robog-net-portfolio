let board = [];
let palettes;
let colors;

let quarterNote;
let wholeNote;
updateBPM(80);

let piano;
let lastMelodicNote = 0;
let orchestra;
let lastOrchestralNote = 0;
let leeway = 40; // Frame times aren't exactly aligned with bpm

let repetition;

let orchestraOffset = 20;
let pianoOffset = 0;


let binary = ["000", "001", "010", "011", "100", "101", "110", "111"];
function createRules(ruleNum) {
  let ruleNumBinary = ruleNum.toString(2);
  let ruleset = {};
  for (let i = 0; i < binary.length; i++) {
    ruleset[binary[i]] =
      ruleNumBinary.charAt(ruleNumBinary.length - (i + 1)) == "1" ? 1 : 0;
  }
  return ruleset;
}

// Try 110, 90, 222
let rules;

let bassSynth;
let melodySynth;

function updateBPM(bpm) {
  quarterNote = 60000 / bpm;
  wholeNote = quarterNote * 4;
}

function playNamedNote(midiNumber, soundFile) {
  const ratio = Math.pow(2, (midiNumber - 60) / 12); // 60 = C4 MIDI
  soundFile.rate(ratio);
  soundFile.play();
}

function playPiano(note) {
  if (millis() - lastMelodicNote > quarterNote - leeway) {
    lastMelodicNote = millis();
  } else {
    return
  }
  playNamedNote(note, piano)
}

function playOrchestra(note) {
  if (millis() - lastOrchestralNote > wholeNote - leeway) {
    lastOrchestralNote = millis();
  } else {
    return
  }
  playNamedNote(note, orchestra);
}

function resetBoard() {
  rules = createRules(round(random(0, 254)));

  for (let col = 0; col < width / 10; col++) {
    board[0][col] = random() > 0.2 ? random(Object.keys(colors)) : "W";
  }
}

function preload() {
  piano = loadSound("middle_c_oboe.mp3");
  orchestra = loadSound("middle_c_orchestra.mp3");
  orchestra.playMode("sustain");
}

function setup() {
  createCanvas(800, 800);
  
  colors = {
    Y: color(255, 188, 10),
    B: color(64, 144, 242),
    R: color(217, 48, 96),
    O: color(240, 93, 35),
    G: color(64, 221, 159),
    P: color(76, 44, 144),
    W: color(255, 248, 240),
  };

  for (let row = 0; row < height / 10; row++) {
    board.push([]);
    for (let col = 0; col < width / 10; col++) {
      board[row].push("W");
    }
  }

  resetBoard();
}

let currentRow = 1;

let lastTime = -1;
let interval = 100;

let lastNumBlanks = -1;
let maxRepeats = 10;

function draw() {
  if (frameCount % 7200 == 10) {
    // Randomize parameters every two minutes-ish
    maxRepeats = random(1, 10);
    orchestraOffset = round(random(-5, 20));
    pianoOffset = round(random(-10, 10));
    updateBPM(random(40, 210));
    interval = quarterNote;
  }

  let hadColor = false;
  background(255, 248, 240);

  // draw the board
  for (let row = 0; row < height / 10; row++) {
    for (let col = 0; col < width / 10; col++) {
      noStroke();
      fill(colors[board[row][col]]);
      let bacteriaSize =
        noise(row, col) * sin(millis() / 300 + row / 4) * 4 + 10;
      if (board[row][col] != "W") {
        ellipse(col * 10 + 5, row * 10 + 5, bacteriaSize, bacteriaSize);
      }
    }
  }

  if (floor(millis() / interval) != lastTime) {
    let numBlanks = 0;

    lastTime = floor(millis() / interval);
    for (let row = 0; row < height / 10; row++) {
      for (let col = 0; col < width / 10; col++) {
        if (row == currentRow) {
          let firstC = board[row - 1][col - 1] ?? random(Object.keys(colors));
          let secondC = board[row - 1][col] ?? random(Object.keys(colors));
          let thirdC = board[row - 1][col + 1] ?? random(Object.keys(colors));

          let ruleToApply = `${firstC == "W" ? 0 : 1}${secondC == "W" ? 0 : 1}${
            thirdC == "W" ? 0 : 1
          }`; //`
          if (rules[ruleToApply] == 1) {
            hadColor = true;
            let newColor = -1;
            let allColors = [];
            if (firstC != "W") allColors.push(firstC);
            if (secondC != "W" && !allColors.includes(secondC))
              allColors.push(secondC);
            if (thirdC != "W" && !allColors.includes(thirdC))
              allColors.push(thirdC);
            if (allColors.length == 3) allColors[1] = allColors[2];
            if (allColors.length == 1) newColor = allColors[0];

            if (allColors.includes("B") && allColors.includes("R"))
              newColor = "P";
            if (allColors.includes("R") && allColors.includes("Y"))
              newColor = "O";
            if (allColors.includes("Y") && allColors.includes("B"))
              newColor = "G";
            if (allColors.includes("G") && allColors.includes("P"))
              newColor = "B";
            if (allColors.includes("P") && allColors.includes("O"))
              newColor = "R";
            if (allColors.includes("O") && allColors.includes("G"))
              newColor = "Y";
            if (newColor == -1)
              newColor = random(allColors) ?? random(Object.keys(colors));

            board[row][col] = newColor;
          } else {
            board[row][col] = "W";
          }
          if (board[row][col] == "W") {
            numBlanks++;
          }
        }
      }
    }

    if (lastNumBlanks == numBlanks) {
      repetition++;
    }
  
    if (repetition > maxRepeats || numBlanks > 70) {
      repetition = 0;
      rules = createRules(round(random(0, 254)));
    }

    lastNumBlanks = numBlanks;

    try {
      //if (random() > 0.7) {
        if (numBlanks < 50) {
          playOrchestra(numBlanks + orchestraOffset);
        } else {
          playPiano(numBlanks + pianoOffset);
        }
      //}
    } catch (e) {}

    if (/*!hadColor || */ currentRow == board.length) {
      currentRow = 0;
      repetition = 0;
      resetBoard();
    }
    currentRow += 1;
  }
}

function mousePressed() {
  userStartAudio();
}
