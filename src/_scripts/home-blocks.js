new p5(function (p) {
  const BLOCK = 40;
  const PADDING = 80;
  const COLORS = ["#ff6fab", "#85f7ff", "#aaf751"];
  let gap;
  let colorBlocks = [];

  const Q = {
    identity: () => ({ w: 1, x: 0, y: 0, z: 0 }),
    fromAxisAngle: (ax, ay, az, angle) => {
      const s = Math.sin(angle / 2);
      return { w: Math.cos(angle / 2), x: ax * s, y: ay * s, z: az * s };
    },
    multiply: (a, b) => ({
      w: a.w * b.w - a.x * b.x - a.y * b.y - a.z * b.z,
      x: a.w * b.x + a.x * b.w + a.y * b.z - a.z * b.y,
      y: a.w * b.y - a.x * b.z + a.y * b.w + a.z * b.x,
      z: a.w * b.z + a.x * b.y - a.y * b.x + a.z * b.w,
    }),
    slerp: (a, b, t) => {
      let dot = a.w * b.w + a.x * b.x + a.y * b.y + a.z * b.z;
      if (dot < 0) { b = { w: -b.w, x: -b.x, y: -b.y, z: -b.z }; dot = -dot; }
      if (dot > 0.9995) {
        return Q.normalize({
          w: a.w + t * (b.w - a.w), x: a.x + t * (b.x - a.x),
          y: a.y + t * (b.y - a.y), z: a.z + t * (b.z - a.z)
        });
      }
      const theta0 = Math.acos(dot);
      const theta = theta0 * t;
      const s0 = Math.cos(theta) - dot * Math.sin(theta) / Math.sin(theta0);
      const s1 = Math.sin(theta) / Math.sin(theta0);
      return {
        w: s0 * a.w + s1 * b.w, x: s0 * a.x + s1 * b.x,
        y: s0 * a.y + s1 * b.y, z: s0 * a.z + s1 * b.z,
      };
    },
    normalize: (q) => {
      const len = Math.sqrt(q.w * q.w + q.x * q.x + q.y * q.y + q.z * q.z);
      return { w: q.w / len, x: q.x / len, y: q.y / len, z: q.z / len };
    },
    toMatrix: (q) => {
      const { w, x, y, z } = q;
      return [
        1 - 2*(y*y + z*z),     2*(x*y + z*w),     2*(x*z - y*w), 0,
            2*(x*y - z*w), 1 - 2*(x*x + z*z),     2*(y*z + x*w), 0,
            2*(x*z + y*w),     2*(y*z - x*w), 1 - 2*(x*x + y*y), 0,
                        0,                 0,                  0, 1
      ];
    },
    apply: (q) => {
      const m = Q.toMatrix(q);
      p.applyMatrix(
        m[0],  m[1],  m[2],  m[3],
        m[4],  m[5],  m[6],  m[7],
        m[8],  m[9],  m[10], m[11],
        m[12], m[13], m[14], m[15]
      );
    }
  };

  const ease = {
    inOutCubic: t => t < 0.5 ? 4 * t * t * t : (t - 1) * (2 * t - 2) * (2 * t - 2) + 1,
  };

  const Q_neutral = Q.identity();
  const Q_tipped = Q.multiply(
    Q.fromAxisAngle(0, 0, 1, Math.PI / 4),
    Q.fromAxisAngle(1, 0, 0, Math.atan(1 / Math.sqrt(2)))
  );

  class ColorBlock {
    constructor(color, index) {
      this.color = color;
      this.index = index;
      this.animations = [];
      this.activeAnimations = [];

      this.twirl = this.twirl.bind(this);
      this.jumpAndFlip = this.jumpAndFlip.bind(this);
      this.aboutFace = this.aboutFace.bind(this);
    }

    animate(fn, start, duration) {
      this._aboutFaceAxis = undefined;
      this.animations.push({ fn, start, end: start + duration });
    }

    update() {
      const now = p.millis();
      this.animations = this.animations.filter(a => now < a.end);
      this.activeAnimations = this.animations.filter(a => now >= a.start && now < a.end);
    }

    // TWIRL
    // A real spinning coin has two coupled motions:
    //   - precession: slow spin of the tilt axis around Y (what you see as the coin "wobbles")
    //   - tilt: how upright the coin stands
    //
    // As a coin falls, tilt → 0 (flat) and precession slows.
    // As it rises, tilt → 90° (upright) and precession speeds up.
    //
    // We model this with a single bell-curve tilt envelope and a spin rate
    // that is proportional to tilt — so they are always physically coupled.
    twirl(progress) {
      // Bell curve for tilt: 0 at start/end, peaks at 1 in the middle.
      // sin(π·t) gives a smooth symmetric arc — no phases, no seams.
      const tilt = Math.sin(Math.PI * progress);

      // Corner-up tilt angle: slerp neutral → Q_tipped by `tilt`.
      // At tilt=0 the cube is flat (neutral). At tilt=1 it's fully corner-up.
      const Q_currentTilt = Q.slerp(Q_neutral, Q_tipped, tilt);

      // Spin: integrate spin rate over time. Spin rate is proportional to tilt
      // so the cube spins fast when upright and barely moves when nearly flat —
      // exactly like a coin. We want ~4 full rotations at peak, so we use the
      // integral of sin(π·t) from 0 to t, scaled to give 4 rotations total.
      //
      // ∫₀ᵗ sin(π·u) du = (1 - cos(π·t)) / π
      // At t=1 this equals 2/π. Scale by 4·π/2 = 2π to get 4 full rotations.
      const spinAngle = (1 - Math.cos(Math.PI * progress)) * 4 * Math.PI;

      // Compose: first apply tilt, then spin around world Y.
      // Left-multiplying by Q_yRot keeps the spin in world space.
      const Q_yRot = Q.fromAxisAngle(0, 1, 0, spinAngle);
      const q = Q.multiply(Q_yRot, Q_currentTilt);

      Q.apply(q);
    }

    // JUMP AND FLIP
    jumpAndFlip(progress) {
      const arc = -4 * progress * (progress - 1);
      const jumpHeight = BLOCK * 1.5;
      p.translate(0, -arc * jumpHeight, 0);
      p.rotateX(progress * p.TWO_PI);
      p.rotateZ(arc * 0.52);
      p.rotateY(progress * p.PI * 0.5);
    }

    // ABOUT-FACE
    aboutFace(progress) {
      if (!this._aboutFaceAxis) {
        const axes = [
          { v: [1, 0, 0], diagonal: false },
          { v: [0, 1, 0], diagonal: false },
          { v: [0, 0, 1], diagonal: false },
          { v: [1, 1, 0], diagonal: true },
          { v: [1, -1, 0], diagonal: true },
          { v: [1, 0, 1], diagonal: true },
          { v: [1, 0, -1], diagonal: true },
        ];
        this._aboutFaceAxis = axes[Math.floor(Math.random() * axes.length)];
      }

      const { v: [ax, ay, az], diagonal } = this._aboutFaceAxis;
      const len = Math.sqrt(ax * ax + ay * ay + az * az);
      const targetAngle = diagonal ? p.PI : p.HALF_PI;
      const angle = ease.inOutCubic(progress) * targetAngle;
      const q = Q.fromAxisAngle(ax / len, ay / len, az / len, angle);
      Q.apply(q);
    }

    display() {
      this.update();

      p.push();

      const x = PADDING + this.index * (BLOCK + gap) + BLOCK / 2;
      const y = PADDING + BLOCK / 2;
      p.translate(-p.width / 2 + x, -p.height / 2 + y, 0);

      let shouldDraw = true;
      for (const anim of this.activeAnimations) {
        const progress = (p.millis() - anim.start) / (anim.end - anim.start);
        const result = anim.fn(progress);
        if (result === false) shouldDraw = false;
      }

      if (shouldDraw) {
        p.fill(this.color);
        p.box(BLOCK);
      }

      p.pop();
    }
  }

  p.setup = function () {
    const remPx = parseFloat(getComputedStyle(document.documentElement).fontSize);
    gap = Math.round(0.5 * remPx);

    const sketchW = BLOCK * 3 + gap * 2;
    const sketchH = BLOCK;

    const container = document.getElementById("color-blocks-sketch");
    container.style.width = sketchW + "px";
    container.style.height = sketchH + "px";

    const canvas = p.createCanvas(sketchW + PADDING * 2, sketchH + PADDING * 2, p.WEBGL);
    canvas.parent("color-blocks-sketch");
    canvas.style("position", "absolute");
    canvas.style("top", -PADDING + "px");
    canvas.style("left", -PADDING + "px");

    p.ortho();

    COLORS.forEach((color, i) => {
      colorBlocks.push(new ColorBlock(color, i));
    });

    p.noStroke();
  };

  const BLOCK_ANIMATIONS = ["twirl", "jumpAndFlip", "aboutFace"];

  function pickAnimation(block) {
    const r = Math.random();
    if (r < 0.9)  return { fn: block.aboutFace,    duration: 1500 };
    if (r < 0.98)  return { fn: block.jumpAndFlip,  duration: 1500 };
    return           { fn: block.twirl,         duration: 3000 };
  }

  p.draw = function () {
    p.clear();

    p.ambientLight(70);
    p.directionalLight(255, 255, 255, 0, 0, -1);

    const now = p.millis();
    if (now > 15000 && Math.floor(now / 5000) > Math.floor((now - p.deltaTime) / 5000)) {
      colorBlocks.forEach((block, i) => {
        const { fn, duration } = pickAnimation(block);
        block.animate(fn, now + i * 200, duration);
      });
    }

    colorBlocks.forEach(block => block.display());
  };
});
