---
title: SNORIS
layout: docs.njk
tags:
  - docroot
description: A well-designed fan controller for Linux systems.
---

(pronounced like "snore" plus "[loris](https://en.wikipedia.org/wiki/Loris)")

# What is a fan controller?

A fan controller should:

- Make fans faster when temperature gets hotter
- Make fans slower when temperature gets colder
- Determine what PWM values affect the fan's actual speed in RPMs, i.e., what is
  low-medium-high speed for your fan
- Shutdown the computer if temperature reaches unsafe values

Fan curves can be really complicated, but most computer fans actually only have
a few true speeds they will reach. So, in reality, fan "curves" behave more like
discrete steps:

<!-- Yes there is HTML in my markdown. Yes it is gross. I needed a diagram, though! -->
<p style="text-align: center; display: flex; flex-direction: column; color: white; background-color: #0d0748; border-radius: 10px; padding: 5px">
<span><em>Below a normal operating temp:</em></span>
<span>Turn fans off or set them to the lowest speed.</span>
<span>⬇⬆︎</span>
<span><em>Somewhere in the middle:</em></span>
<span>Set fans to a medium speed.</span>
<span>⬇⬆︎</span>
<span><em>Could be more than one speed setting in the middle.</em></span>
<span>⬇⬆︎</span>
<span><em>Above a certain high temp:</em></span>
<span>Set fans to maximum speed.</span>
<span>⬇⬆︎</span>
<span><em>Above maximum operating temp:</em></span>
<span>Protect hardware by turning the computer off.</span>
</p>

This means that fans are actually better controlled with discrete states than
with a complex curve, as most fans will only achieve a handful of specific
speeds, and this reduces the complexity of the program controlling them to a
simple state machine. This means that the fan controller can be more reliable
(very important if it's running on a pricey server) and more understandable
(makes it easier to debug and configure).

Additionally, because many modern fan controllers report RPMs, we can use the
hardware to calibrate the fan automatically. This is good because no two fans
are alike, and it's extremely difficult to create a custom fan curve for every
individual fan hardware and configuration. A simple 3-step process can be used
to determine how PWMs map to different fan speeds, and how we can use this to
control them accurately and configure them automatically.

1. Find how stable the fan's RPMs are
   - Set the fan to a "medium" speed like 128 (in a typical 0-255 PWM range)
   - Wait for the user to confirm auditorily that the fan has stopped changing
     speed (or, 5 minutes is probably more than fine for most fans in a headless
     environment)
   - Measure the fan RPMs for 10s
   - Calculate [2SEM](https://en.wikipedia.org/wiki/Standard_error) to determine
     a good range for the fan's "stable" RPM measurements
   - Congratulations, we now don't need any more user input!
2. Find out how long it takes the fan to "spin up" exactly
   - Now that we know when the fan is stable, set the fan's PWM controller to 0
   - Wait for the fan to remain within a 2SEM range for 1s
   - Start a timer, and set the fan's PWM controller to 255
   - Wait for the fan to remain within a 2SEM range for 1s
   - Stop the timer and subtract 1s: This is how long it takes the fan to spin
     up
   - This process should be repeated in reverse, because it might take longer to
     slow down than to spin up, and the maximum value should be used
3. Find every discrete fan speed
   - Set the fan PWM controller to 0 and wait for the spin up time
   - Increase the value sent to the PWM controller by some amount
   - Wait for the spin up time again
   - Check if the current RPMs have moved outside of the 2SEM window
     - If yes, this is a new speed! Remember this PWM value.
     - If no, the fans don't support a speed at this resolution.
   - Check every x values (configurable but I found that 50 is generally fine)
     until 255

# Introducing SNORIS

## A well-designed fan controller for Linux systems.

SNORIS is a simple collection of Python scripts and systemd unit files that runs
the above steps. It configures itself automatically. Most importantly: _When the
computer gets hot, the fans turn on._

The main SNORIS controller script runs in a separate service as the emergency
stop script, meaning that a problem with your fans or sensor configuration that
causes an error or crash will not prevent your computer from shutting down to
save its components. Most computers probably have safeguards in place at the
hardware/firmware level, but this is an extra layer of protection just in case

**P.S.** if SNORIS does crash this way, please submit any information printed in
`systemctl status snoris` or in the complete journalctl logs to
[the GitHub repository](https://github.com/robog-two/snoris). If you're not sure
what the cause is or you don't have complete information to report, that's OK.
Please submit an issue anyways. I would rather close unrelated issues than miss
issues that are mysterious yet fatal.
