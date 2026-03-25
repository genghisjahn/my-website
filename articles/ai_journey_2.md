---
slug: ai-journey-2
title: AI Journey Part 2
date: 2026-03-25
author:
  name: Jon Wear
summary:
tags:
  - name: ai
    slug: ai
  - name: chatGPT
    slug: chatGPT
hero:
  src: /images/hero_article_6a.png
  alt: AI Part 2
css: /css/retro-sci-fi.css
draft: false
---

I could have continued copy/pasting my way to more AI software, but surely there was a better way.  And there was!  Turns out that OpenAI had a ChatGPT plugin for visual studio that gave you a little chatGPT window that could actually _see_ your code and update it while you watched.  This was great!  It actually took me closer to what I wanted because I could see it update my code right away in an environment I was very used to working in. I felt like I was flying.  I'd just code along and if I came to something I wasn't sure about or a bug that seemed particularly nasty, I highlighted the text and it would come up with a change (not always a fix) and ask if I wanted to apply it.  It looked something like this:

<img src="/images/vs_plugin1.png" style="width:50%;" alt="desc2">

Things have changed so much since then that I can't find the exact version of that plugin/extension anyway, but you get the idea.  ChatGPT was now inside my code editor and I could type code, have it help out, apply or not, and keep going.  Really sped things along.  In fact I used that to make the site [Mercury Falling](https://mercuryfalling.net).  It's a site that sends SMS weather alerts to your phone at a time you specify.  So if you wanted to know if it was going to be below freezing in the morning, you could have it check at 8pm the night before and it would send you a text message if it was going to be frigid.  I used this to help write the Go backend (which I mostly wrote by hand) but it helped a lot with the front end because I'm not as well versed in html,css and javascript.  This simple site used `mysql`, `rabbitMQ` & `redis`.  Writes to the site went to a queue([`rabbitMQ`](https://www.rabbitmq.com/)) and workers updated [mysql](https://www.mysql.com/).  Reads would try to pull from the cache ([redis](https://redis.io/)) and if it wasn't there, pull from the database (mysql) and then write to the cache so it would be there next time.  This is where having chatGPT helped me _too much_. I ended up writing a way overly complicated backend when I could have gotten by with just [postgres](https://www.postgresql.org/) and called it a day.  But still, freaking great.  It helped a lot with the integration with Twilio and OpenWeather.  Now, people have been integrating with Twilio and other 3rd party APIs for decades.  The chatGPT plugin didn't make that possible.  But it sure made all the boilerplate hookup code way faster to write and troubleshoot.

I was able to launch a thing (a _product_) if you will.  I used the plugin for infrastructure, backend, frontend, design, datastore.  The whole thing.  It has actual users(some of them I don't even know!).  I've learned _a lot_ since then and I would do this totally differently now (in fact I am, Mercury Falling v2 is in the works).  But the key bit was the VS Code extension from OpenAI:

<img src="/images/vs_plugin2.png" style="width:50%;" alt="desc2">

So that's the second big step I took into AI or LLM assisted programming.  Nothing like what I do today, but really cool at the time (way back in 2025).  I progressed from copy/pasting back and forth from my IDE and a web page to having a floating tool window that updated my code for me as I watched.

Blog post written by me with minor spelling/grammar changes suggested by Claude Code.

Hero image created by ChatGPT.