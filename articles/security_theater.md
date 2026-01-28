---
slug: security-theater
title: Security Theater
date: 2026-01-27
author:
  name: Jon Wear
summary:
tags:
  - name: security
    slug: security
hero:
  src: /images/security_theater.png
  alt: Security Theater
css: /css/retro-sci-fi.css
draft: false
---

When I was 11 years old, writing TI BASIC on my TI-99/4A, I had this idea that I wanted to be able to hide my code.  I had invented this scenario where someone might break into my room and load my boring (but coded by me) games from the cassette tape they were saved on.  I needed to secure my code, but how?  I knew PRINT statements could write to the command line for a user to read something and I knew INPUT statements could accept what the user typed on the command line and let the program do something with it.  And then I thought, "Wait, I could PRINT out "Enter Password" and then INPUT to get what they typed, then I could compare that with a stored password and see if they match.  If they match, then they can use the program and if not, the program will end."  Brilliant!  The only problem was that all a user had to do was type LIST into the command line before running the program and they'd see all of the code as well as what the password was supposed to be.  Then I decided to just put a bunch of REM statements at the top of the code (comments basically) so that if a person did type LIST the code would scroll up ...

You know, as I write this post, I'm getting bored.  There is too much explaining of 1980s consumer tech before I get to a point that just about everyone in tech will agree with already.  Let's just stop here.