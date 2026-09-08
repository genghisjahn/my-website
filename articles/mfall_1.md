---
slug: mfall-1
title: How Mercury Falling Came To Be
date: 2026-08-22
author:
  name: Jon Wear
summary:
tags:
  - name: claudecode
    slug: claudecode
  - name: weather
    slug: weather
  - name: mercuryfalling
    slug: mercuryfalling
hero:
  src: /images/mfall.png
  alt: Mercury Falling
css: /css/retro-sci-fi.css
draft: true
---

I finally launched the second version of Mercury Falling.  Which, as it turns out, is actualy the _third_ version but almost no one outside my family ever say the first version (it was terrible).  Mercury Falling is, by the way, a service that watches weather for specific conditions at a specific time you care about, and when those conditions occur, you get notified.  

At first, as in version 1, the notification was an SMS message.  In fact the entire "app" was SMS messages.  You signed up vis SMS, configured your alerts via SMS, chose your time zone, time of alert, alert criteria...all through SMS.  If you kept your head, arms and legs inside the vehicle at all times and stayed strictly on the happy path, it worked.  But maintaining state via SMS isn't easy (as I quickly discovered) plus the UI isn't very forgiving. Let's suppose it "took off" and _tens_ of users used it.  Well, SMS gets expensive at a certain scale.  But I never got to that because it was just too hard to setup.  My daughter who was 14 at the time actually set it up quickly, but no one over the age of 25 could make heads or tails of it.  So it was scrapped.

The second version (that actually had some users) used SMS as a sign up and alert mechanism, but all of the configuration of an alert was done via the web.  That actually worked well enough.  And I built this thing to SCALE.  I had a binary for the front end and the back end.  I used mysql for the data store, redis to cache, rabbitMQ to handle writes and a series of workers to keep it all humming.  It worked fine, but it was a pain to update.  Also it ran on a server in my basement so whenever the weather got bad and the power went out, not only did my weather alert service go down (really bad timing) but when the power came back I had to reteach myself how to bring all the various services, docker containers and queue workers back on line.  Like I said, this version worked fine but it also ran into the same problem as the first version: SMS gets expensive.  I'd give my Twilio account some cash and a few weeks later it would be gone.  When that account ran out of money, no one got alerts and that's how it would be until I put money into it again.  That got annoying.  Plus sometimes one worker would go down and the others would be fine.  For example, login would stop working but the alerts would go out, or vice versa.  Also I realized early on that getting an alert on a single weather metric was useful but not _that_ useful.  I had users ask for the ability to specify multiple criteria so that they could know if it's good weather for camping, running, skiing, etc.  There were lots of cominbations of weather that a user might want to experience or avoid.  But my database design wasn't great and if I changed the mysql tables, which was easy enough, I'd also have to change redis code, API code and code in the various service workers.  It just wasn't going to work.  I leveraged _a lot_ of chatGPT to get version 2 done.  I'd say probably a 3rd of it was me coding by hand and the rest pasted in from ChatGPT.  This was back in my copy/paste days.  I did _so_ much copy pasting between chatGPT's web UI and VS Code.  I didn't really know how to use LLMs at that time.  I just knew they could help me get unstuck and mostly write code that worked if I guided it well enough.  Anyhoo...

I decided to scrap the whole thing and write it anew.  Only this time I wasn't going to overcomplicate things for a service that really wasn't that complicated.  

(I discovered Claude Code.  Everything went faster.)

For the third version I decided to have one binary listening on two ports (one for the client request and one for the backend incase I want to have multiple clients later).  I dropped mysql, redis, rabbitMQ and docker and did everything in postgres.  There just wasn't any need for all the infracstructure I had in the version two setup.  I also decided to scrap SMS as an alert mechanism and go with just email (using Resend) and push notifications for Progressive Web Apps(PWAs).  I was able to implement magic link logins via email & passkeys as well as the multi-criteria weather alerts, which worked great!  Also, instead of having a user specify each alert criteria, they could instead choose to _describe_ the kind of weather they wanted to be alerted about and a cheap LLM would built the alert entry for them, done!

There parts of the version 2 wheel that I didn't want to reinvent, so I created a subdirectory called "version_1" (really version 2) in the new repo and pointed claude code to that from time to time whenever I wanted to pull something from the old version.  Development took about one month to go from start to finish.  