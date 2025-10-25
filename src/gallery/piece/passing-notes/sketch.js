let isDrawing = false;

function setup() {
  createCanvas(250, 250);
  background(255);
  let btn = createButton("Print to my Desk");

  btn.mousePressed(sendCanvas);
}

function draw() {
  if (mouseIsPressed && mouseY < height) {
    stroke(0);
    strokeWeight(4);
    line(pmouseX, pmouseY, mouseX, mouseY);
  }
}

function sendCanvas() {
  // Get the canvas as a Base64-encoded PNG image
  const imgData = canvas.toDataURL("image/png");

  // Prepare JSON payload.
  const data = {
    bitmap: imgData,
  };

  // Send POST request
  fetch("https://printer-reverse-proxy.deno.dev", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((res) => {
      if (res.ok) {
        alert(
          "Your note has been received and will be printed on my desk shortly.",
        );
        background(255);
      } else {
        throw new Error("Server error!");
      }
    })
    .catch((err) => alert("Error: " + err));
}
