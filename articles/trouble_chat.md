---
slug: trouble-with-the-robot
title: Trouble with the Robot
date: 2025-08-15
author:
  name: Jon Wear
summary: It saves so much time.  It wastes so much time.
tags:
  - name: chatGPT
    slug: chatGPT
hero:
  src: /images/hero_article_2.png
  alt: Corridor lights
css: /css/retro-sci-fi.css
draft: false
reading_time_min: 2
---

Where was I? Oh yes. My super simple/fast web server was working. My content was building, it was rendering just like I wanted. Then I asked for a script that would copy the server binary from my local machine to... the server. And it took 4-Evah!

Back and forth, forgetting command line args, forgetting the directory, making it way more complicated than it needed to be and it never quite working. It would even do some kind of unicode ... ellipses thing that prevented the script from running at all. I finally threw it all away and had it do each part of the script one at a time. First build, then stop remotely, then copy out... etc.

It works now, but jeez.