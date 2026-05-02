---
title: Podman + Systemd = Awesome!
layout: blog.njk
date: 2026-01-09
tags:
  - post
description: yes
---

Recently in my work for [Localify.org](https://localify.org/) I needed to connect an AWS-managed Prometheus database to a Spring Boot application. I learned about OpenTelemetry, specifically, AWS's own fork of it called ADOT. It turns out that we already had some pre-paid EC2 instances on our account, so I decided I would run the container using a setup I've used in the past on my own personal homelab infrastructure. I figured that it would be pretty useful to publish somewhere on the web, because I have a pretty great workflow that requires some specific little tweaks.

I love this setup because it allows you to keep configuration for an entire project in one user, which means that you can place multiple compute resources on the same machine, using different methods of containerization or orchestration (or none at all). Plus, it adds an extra layer of isolation to reduce the impact of sandbox escape issues with Podman, and simply helps with keeping project files organized with all necessary files easy to back-up in one folder.

# Setting Up

Dependencies, as they are on Debian 13.
```bash
sudo apt install podman podman-compose systemd-container
```

Create a user with the default settings from your distro

```bash
sudo adduser <user>
```

Fill out the password with `123` (or anything) and then press enter for all other information.
Next, lock the account's password.

```bash
sudo passwd -l <user>
```

Then we can make that account stick around, so that it's user systemd services will keep running even when we leave the machine unattended.

```bash
sudo loginctl enable-linger <user>
```

Now we will switch to the user we just created and configure the rootless container. **Do not add this user to sudoers, and do not use sudo on this user!** The user we created essentially acts as a second layer of defense. If there is a vulnerability in podman, an attacker must then also escape the restrictions imposed by the linux kernel in addition to the security of podman's containerization. Also, if the machine needs to be repurposed, we can simply delete the user and it will clean up all files, services, and containers at once. No mess!

```bash
sudo machinectl shell <user>@
```

OK, then we can set up podman. Below is an example configuration that I used, you can modify it according to your needs.
```bash
mkdir podman-compose; cd podman-compose;
touch podman-compose.yml; touch podman-compose-daemon.service;
```

You would place your regular docker-compose file in podman-compose.yml. Then, place the following in podman-compose-daemon.service:

```systemd
[Unit]
Description=Run podman-compose under this user

[Service]
Type=simple
WorkingDirectory=/home/the-usr/podman-compose/
ExecStartPre=/usr/bin/podman-compose down
ExecStart=/usr/bin/podman-compose up
ExecStop=/usr/bin/podman-compose down

[Install]
WantedBy=default.target
```

# Bonus Tips

- Check all the files in ~/podman-compose into git to back them up and track changes over time
- Use journalctl to get historic logs from podman for debugging purposes
- Remember that podman containers are rootless! If you are having issues accessing graphics cards, lower ports, or other features in the container that work with docker, check podman documentation to see how to work around it.
- Verify that the user service continues running even when logged out if it's critical. This is related to that `loginctl enable-linger` command
- It's always useful to use `podman run -it debian bash` if you need a simple shell to check what the container can access by default
- You can use `systemd` timer units to regularly restart, pull, or build containers, which is useful in scenarios where you might want to regularly deploy code (on my personal server, code is deployed nightly)
