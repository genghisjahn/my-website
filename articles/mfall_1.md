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

I finally launched the second version of Mercury Falling.  Which, as things would have it, is actualy the third version but almost no one outside my family ever say the first version (it was terrible).  Mercury Falling is service that watches weather for specific conditions and a specific time you care about, and when those conditions occur, you get notified.  

At first the notification was an SMS message.  In fact the entire "app" was SMS messages.  You signed up vis SMS, configured your alerts via SMS, choose your time zone, time of alert, alert criteria all through SMS.  If you kept head, arms and legs inside the vehicale at all times and stayed strictly to the happy path, it worked.  But maintain state via SMS isn't easy (as I quickly discovered) plus it the UI isn't very forgiving AND let's suppose it "took off" and _tens_ of users used it.  Well, SMS gets expensive at a certain scale.  But I never got to that because it was just too hard to setup.  My daughter who was 14 at the time actually set it up quickly, but no one over the age of 25 could make heads or tails of it.  So it was scrapped.

The second version (that actually had some users) used SMS as a sign up and alert mechanism, but all of the configuration of an alert was done via the web.  That actually worked well enough.  And I built this thing to SCALE.  I had a binary for the front end, the back end.  I used mysql for the datastore, redis to cache, rabbitMQ to handle writes and a series of workers to keep it all humming.  It worked fine, but it was a pain to update.  Also it ran on a server in my basement so whenever the weather got bad and the power went out, not only did my weather alert service go down (really bad timing) but when the power came back I had to reteach myself how to bring all the various services, docker containers and queue workers back on line.  Like I said, this version worked fine but it also ran into the same problem as the first version: SMS gets expensive.  I'd give my Twilio account some cash and a few weeks later it would be gone.  When that account ran out of money, no one got alerts and that's how it would be until I put money into it again.  That got annoying.  Plus sometimes one worker would go down and the others would be fine.  For example, login would stop working but the alerts would go out, or vice versa.  I decided to scrap the whole thing and write it anew.  Only this time I wasn't going to overcomplicate things for a service that really wasn't that complicated.  For second version I leveraged _a lot_ of chatGPT.  This was back in the copy/paste days.  I did so much copy pasting between chatGPT's web UI and VS Code.  I didn't really know how to use LLMs at that time.  I just knew they could help me get unstuck and mostly write code that worked if I guided it well enough.

Then I discovered Claude Code.

I wrote out what the new Mercury Falling was going to do.  I designed the database, I spec'd out the web pages needed, the routes the API would need, all that kind of thing.  How I'd do auth (magic links or passkeys).  What I'd use for payments (Stripe) and then had Claude turn that into a design document.