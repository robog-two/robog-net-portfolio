---

title: I wrote a fan controller! layout: blog.njk tags:

- post description: Snoris is a simple and elegant fan controller for Linux
  systems.

---

# What is a fan controller?

A fan controller should:

- Make fans faster when temperature gets hotter.
- Make fans slower when temperature gets colder.

Fan curves can be really complicated, but many computer fans only reach a few
speeds. So, in reality, fan "curves" behave more like discrete steps:

<!-- Yes there is HTML in my markdown. Yes it is gross. I needed a diagram, though! -->

<p style="text-align: center; display: flex; flex-direction: column; color: white; background-color: #050507; border-radius: 10px; padding: 5px">
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
with a complex curve, since most fans will only achieve a handful of specific
speeds, and this reduces the complexity of the program controlling them to a
simple state machine. This means that the fan controller can be more reliable
(very important if it's running on a pricey server) and more understandable
(makes it easier to debug and configure).

Additionally, because many modern fan controllers report RPMs, we can use the
hardware to calibrate the fan automatically. This is good because no two fans
are alike, and it's extremely difficult to create a custom fan curve for every
individual fan hardware and configuration. A simple 2-step process can be used
to determine how PWMs map to different fan speeds, and how we can use this to
control them accurately and configure them without user input.

1. Find how stable the fan's RPMs are:
   - Set the fan to the highest speed.
   - Wait 1 minute for the fans to achieve this speed.
   - Measure the fan RPMs for 10 seconds.
   - Calculate the range of RPM values based on the measured data.
2. Find every discrete fan speed:
   - Set the fan PWM controller to 0 and wait for the RPM values to stay within
     the range we found.
   - Increase the value sent to the PWM controller by some amount.
   - Wait for the RPM to stabilize within the normal range.
   - Check if the difference between this speed and the previous speed is more
     than the normal range:
     - If yes, this is a new speed! Remember this PWM value.
     - If no, the fans don't support a speed at this resolution.
   - Check every x values (configurable, but I found that 50 is generally fine)
     until 255.

# Why a new fan controller?

The current state-of-the-art fan controller on Linux systems is the kernel. 99%
of users do not need to mess with their fan settings, because the kernel modules
for their specific motherboard or device generally come with built-in settings
and configurations that are faster, better, and easier to set up. This is great
for boards that have them, and if you are one of these people, a fan controller
like Snoris will not help you, and will probably make your fans louder and your
computer hotter.

However, for devices that don't have kernel modules which control their fans,
like the two Dell OptiPlex 5050s that power parts of this website, your only
other option is `fancontrol`. I think `fancontrol` is a package that has grown
over the years to address too many different use cases since the features are
plentiful but the configuration is confusing and very manual.

I wanted something that is simple, so that my servers would turn on their fans
to improve performance and longevity. I wanted a tool that could configure
itself, using built-in sensors on the device, and an associated fan daemon that
is simple and responds quickly to sensor data.

I created Snoris with these constraints in mind. It's a simple collection of two
Python scripts and one systemd unit file. It configures itself automatically.
Most importantly: _When the computer gets hot, the fans turn on._

Install it yourself from the
[GitHub repository](https://github.com/robog-two/snoris).

[More detailed documentation.](https://robog.net/docs/snoris)
