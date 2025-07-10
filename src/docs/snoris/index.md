---
title: Snoris
layout: docs.njk
tags:
  - post
description: A simple and elegant fan controller for Linux systems.
---

For most use cases, you will simply need to run:

```sh
sudo python3 snoris_setup.py
```

## General Troubleshooting

It's important to read all errors and logs carefully as they will generally
indicate where the problem lies. You can find the logs of snoris with:

```sh
sudo systemctl status snoris
```

Or with:

```sh
sudo journalctl -u snoris
```

Snoris needs three things to function:

- Sensor data from your kernel
- Ability to control all fans with pwm
- RPM data from the fans

All of these are achieved through device files in `/sys/class/hwmon/`. You can
verify fans are working by looking for files like `pwm1` and `fan1_input`. For
temperature sensors, it's easier to use `lm_sensors` from your distribution's
package manager. First, run

```sh
sudo sensors-detect
```

After detecting all available sensors and installing appropriate kernel modules,
try:

```sh
sensors
```

This command should give you some temperatures in celsius if everything works.

## Dell Systems

You may need to enable dell-smm-utils kernel module.

## Lenovo Systems

Snoris is not necessary for these devices. If you must use Snoris, you will
probably need to disable the built-in fan controller.

```sh
sudo rmmod thinkpad_acpi
sudo modprobe thinkpad_acpi fan_control=1
sudo bash -c "echo 1 > /sys/class/hwmon/hwmon4/pwm1_enable"
```
