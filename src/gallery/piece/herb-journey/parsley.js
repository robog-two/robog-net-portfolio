function renderParsley() {
  pointLight(color("#DDFFD6"), 0, 0, 0);
  
  push(); // begin butterfly

  let wingRot = abs(sin(millis() * timeAdjust));
  
  texture(parsleyTex);
  textureMode(NORMAL);

  noStroke();
  beginGeometry();
  beginShape(TRIANGLE_STRIP);
  // Left wing
  vertex(cos(PI + wingRot) * scl, sin(PI + wingRot) * scl, scl / 2, 1, 0);
  vertex(0, 0, 0, 0, 0);
  vertex(0, 0, scl, 1, 1);
  // Right wing
  vertex(cos(-wingRot) * scl, sin(-wingRot) * scl, scl / 2, 0, 1);

  endShape();
  let parsley = endGeometry();
  parsley.computeNormals();
  drawTransparent(() => {
    model(parsley);
  });
  pop(); // end butterfly
}
