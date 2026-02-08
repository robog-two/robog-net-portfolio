const neighborDistance = 100;
const maxForce = 0.8;
const maxSpeed = 5;

function updateBoids(boids, cubeSize) {
  for (let thisBoid of boids) {
    thisBoid.acc = createVector(0, 0, 0);
    let separation = separateB(thisBoid, boids);
    let alignment = alignB(thisBoid, boids);
    let cohesion = cohesionB(thisBoid, boids);

    // Arbitrarily weight these forces
    separation.mult(1.5);
    alignment.mult(1.0);
    cohesion.mult(1.0);
    
    // Add the force vectors to acceleration
    thisBoid.acc.add(separation);
    thisBoid.acc.add(alignment);
    thisBoid.acc.add(cohesion);

    thisBoid.vel.add(thisBoid.acc);

    // Limit speed
    thisBoid.vel.limit(maxSpeed);
    thisBoid.pos.add(thisBoid.vel);
    
    
    // Limit positions
    if (thisBoid.pos.x > cubeSize/2) {
      thisBoid.pos.x = -cubeSize/2;
    }
    if (thisBoid.pos.x < -cubeSize/2) {
      thisBoid.pos.x = cubeSize/2;
    }
    if (thisBoid.pos.y < -cubeSize) {
      thisBoid.pos.y = 0;
    }
    if (thisBoid.pos.y > 0) {
      thisBoid.pos.y = -cubeSize;
    }
    if (thisBoid.pos.z > cubeSize) {
      thisBoid.pos.z = 0;
    }
    if (thisBoid.pos.z < 0) {
      thisBoid.pos.z = cubeSize;
    }
  }
}

// thanks shiffman
function separateB(thisBoid, boids) {
  let desiredSeparation = 25.0;
  let steer = createVector(0, 0, 0);
  let count = 0;

  // For every boid in the system, check if it's too close
  for (let boid of boids) {
    let distanceToNeighbor = p5.Vector.dist(thisBoid.pos, boid.pos);

    // If the distance is greater than 0 and less than an arbitrary amount (0 when you are yourself)
    if (distanceToNeighbor > 0 && distanceToNeighbor < desiredSeparation) {
      // Calculate vector pointing away from neighbor
      let diff = p5.Vector.sub(thisBoid.pos, boid.pos);
      diff.normalize();

      // Scale by distance
      diff.div(distanceToNeighbor);
      steer.add(diff);

      // Keep track of how many
      count++;
    }
  }

  // Average -- divide by how many
  if (count > 0) {
    steer.div(count);
  }

  // As long as the vector is greater than 0
  if (steer.mag() > 0) {
    // Implement Reynolds: Steering = Desired - Velocity
    steer.normalize();
    steer.mult(maxSpeed);
    steer.sub(thisBoid.vel);
    steer.limit(maxForce);
  }
  return steer;
}

// Alignment
// For every nearby boid in the system, calculate the average velocity
function alignB(thisBoid, boids) {
  let neighborDistance = 50;
  let sum = createVector(0, 0, 0);
  let count = 0;
  for (let i = 0; i < boids.length; i++) {
    let d = p5.Vector.dist(thisBoid.pos, boids[i].pos);
    if (d > 0 && d < neighborDistance) {
      sum.add(boids[i].vel);
      count++;
    }
  }
  if (count > 0) {
    sum.div(count);
    sum.normalize();
    sum.mult(maxSpeed);
    let steer = p5.Vector.sub(sum, thisBoid.vel);
    steer.limit(maxForce);
    return steer;
  } else {
    return createVector(0, 0, 0);
  }
}

// Cohesion
// For the average location (i.e., center) of all nearby boids, calculate steering vector towards that location
function cohesionB(thisBoid, boids) {
  let neighborDistance = 50;
  let sum = createVector(0, 0, 0); // Start with empty vector to accumulate all locations
  let count = 0;
  for (let i = 0; i < boids.length; i++) {
    let d = p5.Vector.dist(thisBoid.pos, boids[i].pos);
    if (d > 0 && d < neighborDistance) {
      sum.add(boids[i].pos); // Add location
      count++;
    }
  }
  if (count > 0) {
    sum.div(count);
    return seekB(thisBoid, sum); // Steer towards the location
  } else {
    return createVector(0, 0, 0);
  }
}

function seekB(thisBoid, target) {
    // A vector pointing from the location to the target
    let desired = p5.Vector.sub(target, thisBoid.pos);

    // Normalize desired and scale to maximum speed
    desired.normalize();
    desired.mult(maxSpeed);

    // Steering = Desired minus Velocity
    let steer = p5.Vector.sub(desired, thisBoid.velocity);

    // Limit to maximum steering force
    steer.limit(maxForce);
    return steer;
  }
