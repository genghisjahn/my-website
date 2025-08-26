---
slug: no-hitters-are-difficult-to-simulate
title: No-hitters are difficult (to simulate)
date: 2025-08-26
author:
  name: Jon Wear
summary: No-hitters in MLB are rare, and they are hard to simulate
tags:
  - name: baseball
    slug: baseball
  - name: no-hitter
    slug: no-hitter
hero:
  src: /images/hero_article_4.png
  alt: no-hitters
css: /css/retro-sci-fi.css
draft: false
reading_time_min: 2
---

Baseball is a fun sport to play with from a data modeling point of view because there are so many numbers associated with each game and there are over 150 seasons of baseball stats to look over.  A single season has 2,430 individual games.  That's at least 131,220 plate appearances per season and don't get me started on the number of pitches, swings, foulballs, etc there are in addition.  But let's get to the point, the no-hitter.  

A no-hitter is when a pitcher goes the entire game without allowing the other team to record a single hit. Most MLB games are 9 innings, each inning has 3 outs so that means the pitcher records 27 consecutive outs without allowing a hit.  That is very hard to do.  There have only been 326 no hitters in the approximately 150 years of MLB [^1].  The last one was September 4th, 2024 (although the Cubs used 3 pitchers to do it [^2] ).

Anyhoo, I thought it would be fun to write a simple model of a baseball season to see how many no hitters occurred just from a normal random chance based on the hitting abilities of the players.  Even though a no-hitter is primarily a feat showing the prowess of the _pitcher_, I would not take pitching into account _at all_ (I told you this would be a simple model).  My model works like this.  I take a line up of 9 players.  Each player would have a batting average (BA) and an On Base Percentage (OBP).  BA is the number of actual hits a player has divided by the number of overall at bats.  OBP is the number of times the player got on base (due to a hit, a walk, or getting hit by a pitch).  

1. I loop through all 9 players a total of 27 times (9 innings, 3 outs each inning).  
1. For each loop, I generate a random number from 0 to 1.  
    1. If that number is less than the current player's BA, then we have a hit.  We stop the model, the game can't be a no hitter.
    1. If that number is greater than the current player's BA but less than the current player's OBP, that means the player got walked, or got hit by a pitch.  So they are not out, but they didn't hit the ball either.
        1. Yes, it is possible to pitch 4 walks in a row, lose a game 0-1 and still record a no hitter.
    1. If the number is greater than the current players OBP, the player is out and we move on.

There is some extra code in there too because sometimes a player can get to first base on a walk for example, and then the next player could hit into a double play (where two outs are record in one play).  So the model tracks if a runner is on first base due to a non-hit and if any following player in the inning is out, there is a 10% chance that two outs will be recorded.  Note, this would only make sense if there are 0 or 1 outs at the time.  That's it.  Check BA and OBP.  If say, [**Hu's on first**](http://cdn2.sbnation.com/imported_assets/983249/405694_178315385608995_100002916031264_290280_582254303_n_medium.jpg) and the next player gets out there's a 10% both players are out.

How many no-hitters does my simulation come up with?  I've run it a bunch of times and the range is usually between 98 and 115 in a 150 years of seasons.  What this tells us is the pitcher has a big impact on how often a no hitter happens.  That was _obvious_ before I ran the model.  But it does quantify for me _how much_ of an impact the pitcher has in this situation.  I now know that a pitcher's ability more than doubles the chances of a no hitter over and above the ability of the hitters on the other team.  Well, I think I know that.  Baseball in general is hard.  No hitters are harder and statistics are more difficult yet.  Here's my code if you want to yell at me about it: [baseball_nohit](https://github.com/genghisjahn/baseball_nohit/blob/main/main.go).

Some interesting notes on no-hitters in baseball:

1. There was a game were _both_ pitchers threw a no-hitter through the regulation 9 innings.  [May 2nd, 1917](https://sabr.org/gamesproj/game/may-2-1917-fred-toney-and-reds-prevail-1-0-in-double-no-hitter-against-cubs-hippo-vaughn/).  It ended in the 10th inning with a final score of 1-0.
1. [Johnny Vander Meer](https://en.wikipedia.org/wiki/Johnny_Vander_Meer) threw no hitters in back to back games in June of 1938.
1. [Nolan Ryan](https://en.wikipedia.org/wiki/Nolan_Ryan) has the most no hitters (7) all time , is tied for second for the most one-hitters (12) and also holds the record for most strike outs with 5,714.  Here's [the video](https://www.youtube.com/watch?v=L9m_Kk4kzAY) of him getting the last out of his 7th no hitter.  He was 44 years old at the time.



[^1]: MLB is very old or very young depending on how you look at it.  The National League started in 1876, the American League in 1901, however the to leagues did not officially combine until 2000. [nytimes](https://www.nytimes.com/1999/09/16/sports/baseball-league-presidents-out-as-baseball-centralizes.html)

[^2]: Shota's 7 frames start Cubs' first no-no at Wrigley since 1972 [mlb.com](https://www.mlb.com/news/shota-imanaga-starts-cubs-combined-no-hitter)

AI Notice: I used chatGPT for grammer and spelling checks only.  No LLM was used to create the text of this post.  The hero image at the top was created by chatGPT.