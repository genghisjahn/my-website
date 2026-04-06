---
slug: anglican-history-podcast
title: A Survey of Anglican History
date: 2026-04-06T14:00
author:
  name: Jon Wear
tags:
  - name: podcast
    slug: podcast
  - name: anglican
    slug: anglican
  - name: episcopal
    slug: episcopal
source: https://anglicanhistory.jonwear.com
draft: false
---

I was confirmed in the Episcopal Church on January 31st, 2026.  I'd been attending services for several years prior to that at [St. Paul's Chestnut Hill](https://stpaulschestnuthill.org).  I read books (lots of books) and attended confirmation classes.  Anyway, during Lent this year (2026) I decided to give up espresso and music.  Not as easy as I thought it would be. I ended up listening to a lot of podcasts and drinking _a lot_ of coffee.  I looked for Episcopal podcasts, and while there are several, many of them no longer publish new episodes or do not cover the history of the church that I'm interested in.

I started reading Wikipedia articles and then realized I could make a podcast for this.  I could pull [Wikipedia articles](https://en.wikipedia.org/wiki/Christianity_in_Roman_Britain) (with proper attribution of course), clean them up so they would read better, clone my voice using [pocket-tts-server](https://github.com/ai-joe-git/pocket-tts-server), make some background music via [suno.com](https://suno.com), write some scripts that merge the text with the music at the right volumes using [ffmpeg](https://ffmpeg.org) scripts and...done (there are a few other steps, but nothing [Claude Code](https://claude.com/product/claude-code) can't handle).  It's hosted on a small [Hetzner](https://www.hetzner.com) server and I get caching from [Cloudflare](https://www.cloudflare.com) since the RSS feed is served via a Cloudflare tunnel.  I have a podcast that I want to listen to covering the topics I choose.  Takes about 3 minutes to make a new episode.

I'm working to make the voice cloning more natural.  There are odd cadences (too fast, too slow) now and then and sometimes the pronunciation is off, but it works for me.  If it works for anyone else, great.  But it's a podcast tailored to what I wanted to learn about _and_ I learned a lot of technical things getting it all set up.

Podcast is here, [A Survey of Anglican History](https://podcasts.apple.com/us/podcast/a-survey-of-anglican-history/id1890756025), but also on Spotify, Youtube, Amazon and other podcasts platforms.