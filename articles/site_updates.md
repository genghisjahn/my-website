---
slug: site-updates-1
title: Site Updates
date: 2026-01-05
author:
  name: Jon Wear
summary: Added some updates I've been wanting to the deploy pipline as features I thought were missing.
tags:
  - name: web
    slug: web
hero:
  src: /images/llm_danger.png
  alt: Site Updates
css: /css/retro-sci-fi.css
draft: false
reading_time_min: 2
---

I've read a lot about how blogging can be helpful to ones career as well as enjoyable in and of itself.  I have no idea who's reading this because I don't have Google Analytics on here or any other analytics([Website Privacy Report](https://pro.nixondigital.io/preview/jonwear.com)).  I did just add [_web mentions_](https://en.wikipedia.org/wiki/Webmention) using [webmention.io](https://webmention.io) so I guess there is some tracking, but that's only for users who explicity reference my content, or I theirs.

Also, I've been wanting my own quick notes reference for a while.  This would mostly be quick and dirty TIL type content, but also just quick notes.  I think of it as articles are things I'm writing with the intention of sharing (whether anyone reads it or not) and Notes are thing I'm writing for my own reference (whether anyone else reads them or not).  There might be a kind of short cut, postbin kind of thing I might add to that, but that would involve interacting with github so I can make sure I don't lose the content.  I'll think about it.

Since this little blog engine of mine is a Hugo lite engine of sorts, I realized that I should stop tracking the `/public` directory since it can be rebuilt at any time as long as I have the source `.md` files.  People say that RSS is dead and at the same time is making a comeback, so I added RSS feeds for articles and notes (which you'll see at the top).  That site-nav at the top is new as well.

For a while I was just making images from chatGPT and putting the PNG as a hero image at the top of the articles.  That worked fine but those images are huge, so as part of the deploy process now I convert all .png, .jpg and other images files to [.webp](https://en.wikipedia.org/wiki/WebP) using [cwebp](https://developers.google.com/speed/webp/docs/cwebp).  That really helps a lot with image size.  For example: 

```bash
hero_article_1.webp: 658KB → 22KB (96%)
```

I run through a website grader and saw I set expiration headers set and use gzip.  Done & done.

Webmentions are an interesting way to see who is referencing your content (as mentioned before).  It wasn't the easist thing to setup and test, but I'm pretty sure I have it configured properly now.  I was able to send a test mention and a mention from my site to my site on this [article](articles/human-text-and-robot-html/).  I hope it's a way to help foster the spirit of the intenet communities of yore.

Lastly, claude code was very helpful in getting this setup.  Being able to talk to code and github is a real game changer for development.  This little bullet list at the end is the output the prompt:

> okay, look at all the commits on main for today and let me know what features we implemented.

> give me a markdown list of these items, oldest to newest.

  - Add Notes section for quick reference snippets with datetime sorting and source field
  - Add top nav, untrack public/ from git
  - Add RSS feeds for posts and notes with autodiscovery links
  - Add WebP image conversion during deploy with per-file compression stats
  - Add Expires header and more cacheable file types
  - Add gzip compression for HTML, CSS, JS, and text files
  - Add webmentions support with facepile display
  - Add webmention.app sending for posts and notes feeds
  - Add h-entry microformat markup and h-card with author photo
  - Move nav bar above title on article and note page