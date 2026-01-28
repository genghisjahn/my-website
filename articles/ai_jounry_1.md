---
slug: ai-journey-part-1
title: My AI Journey Part 1
date: 2026-01-28
author:
  name: Jon Wear
summary:
tags:
  - name: ai
    slug: ai
hero:
  src: /images/ai_part_1.png
  alt: AI Part 1
css: /css/retro-sci-fi.css
draft: false
---

I think I started using ChatGPT when it was at 3.5.  It was fun to mess around with.  I'd ask it questions about history, some it'd get right, some it'd get wrong.  I rememeber one conversation where it was sure that Slash had played on a Michael Jackson album.  This was after I'd asked it about Eddie Van Halen.  It was fun seeing what it could get right and how confidently it would tell me something that was wrong.  Then I got the idea to write some code for me.  The prompt was something like this: _Write two rest APIs Golang.  One is a Post that accepts a json body for movieID and title.  It writes then takes that data and writes it to a mysql database.  The other is a Get that accepts a movieID, looks in a Redis cache for that Id and if it's not there, then it pulls from a mysql database, writes the data to redis and then gives a json resopnse to the caller._

Ha! Let the hilarity ensuse!  

Except it didn't ensue.  The code it wrote was just fine.  It left place holders for the connection strings to the database and the redis server.  It made a simple struct for the movie object. It made a create table .sql file for the movie table.  The code would have run fine. 

I asked it to write some tests.  It wrote some, not great, but okay.  I asked it to make interfaces for all of the data sources so I could have more unit tests to cover business logic...and it redid the code so that it could mock the database and the redis server.  Not bad. I copied the code out of the browser, pasted it into VS Code and a few minutes later...it ran.  Just fine.  Writes went to the DB, read came from the DB the first time and then from Redis (it set the TLL to 60 seconds without me telling it to).

No, this isn't the most complicated code in the world.  Most backend developers could write this without any real effort at all.  But the tests?  The dependecy injection?  And there's no way they could type it as fast as the LLM printed it out a line at a time like Joshua/WOPR in Wargames(I think this is my third post where I mention that movie).  I rememeber sitting and thinking, well this could change things.  If I have a code base that already has a comprehensive test suite, I could use this to make changes.  Golang is compiled, so if the code is just _bad_ then I'll know the change didn't work right awawy.  Since I have a lot of tests already I'll know if it made a breaking change that compiled.

And that's how I worked with LLMs for a while.  I'd type descriptions of functions into ChatGPT and then copy the code into my IDE.  Sometimes I'd ask it for syntax for obscure commands.  Sometimes I'd paste code into ChatGPT and ask it to tell me what it did, and if there were any errors or ways it could be improved.  At the time the sugguestions for making the code better often wasn't better, but it got really good at explaining code to me.  Copy to ChatGPT or copy back to VSCode.  The two tools were separate from each other and didn't interact directly at all.