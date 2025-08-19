---
slug: baseball_1
title: Calculating the Optimal Lineup
subtitle: It's not going well..
date: 2025-08-19
author:
  name: Jon Wear
summary: Trying to use some software and some Monte Carlo to figure out the best lineup for the Phillies
tags:
  - name: baseball
    slug: baseball
  - name: Monte Carlo
    slug: monte_carlo
hero:
  src: /images/hero.png
  alt: alt name
css: /css/retro-sci-fi.css
draft: true
---

For the past few years my family has been following Major League Baseball.  It started when we all had covid and we started watching the Phillies playoff run in 2023.  We quickly adopted favorite players even though we didn't really know much about any of them.  It was fun family times. 

My two sons play little league now and the five of us still watch a lot of baseball.  My wife and I look at each other sometimes in bewilderment.  "Who are we?  We just talked for 20 minutes about possible minor leauge prospects."  I say all that to let you know that while I've been aware of baseball for my whole life (it was my father's favorite sport), I haven't really followed it until recently.  

Baseball is a game of stats and the historical record for those stats that goes way, way back (100+ years).  One of the most referenced stats from the game is the Batting Average, or BA.  The batting average is the number of total hits a player has divided by number of batting attempts.  It doesn't tell you what kind of hit the player got, just how that they hit the ball and get to remain on the diamond somewhere instead of getting an out and heading back to the dugout.  Before the game starts each team submits a lineup card that listing what order their nine starting players will bat.  This order proceeds, 1-9 for entire game. There's a lot of thought that goes into this.  Here is my understanding of how a line up is constructed, just from watching games and listening to commentary.

* 1st - Lead off Hitter
    * This is someone who good at getting on base.  They don't necessarily hit for power (BA does not measure power, just hits) but they can get on base with a single, or they can be patient and get a walk (and possible hit by a pitch).  You just want some one that can get on base in this spot.
* 2nd - (two hole)
    * Pretty much the same thing.  Maybe this player has a little more power but you want some one that is likely to get a hit.  Would probably be a player near the top of your batting average.
* 3rd - (three hole)