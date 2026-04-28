---
slug: ai-journey-3
title: AI Journey Part 3
date: 2026-04-28
author:
  name: Jon Wear
summary:
tags:
  - name: ai
    slug: ai
  - name: claude
    slug: claude
  - name: anthropic
    slug: anthropic
hero:
  src: /images/hero_article_6b.png
  alt: AI Part 3
css: /css/retro-sci-fi.css
draft: false
---

And then I met Claude Code.  It had been out a while but my new job leveraged it a lot and I figured out I might as well get to know the tool.  I'd started using it around November of 2025, but not really doing a lot with it.  Still writing a lot of code by hand.  I would use it and check it's outputs and I then I would _type_ the code in myself.  I just didn't trust it at first.  But after a while I realized that if I gave better prompts (and by better I mean smaller focus with more detail) it could do the "typing" faster than I could.  I just needed to make sure I covered tests, read over the code and basically understood what I was doing.  It felt dangerious...and scary.  But I noticed I was spending more and more time typing in the terminal window at the bottom of VS Code and less time looking at the big windows of the code edititor that actually contained the source code.  Most of my code review was happening in Github Pull Requests.  It was some time in mid January 2026 that I tiled six terminal windows across my wide screen and never opened VS Code at all.  It was one of my most productive days ever.  Not in the LOC written sense, but in the features shipped sense.  I was still very nervous about this but I realized that something had changed fundamentally and the way I've coded in the past is not going to be the way I coded in the future.  Or...maybe it is.  But it's not going to look the same to an outsider.  

I had spent the day coding without opening a code editor.  I just sat there for a while staring at my screen, wondering what had just happepend.  I was reminded of this quote from Zen and the Art of Motorcycle Maintenaince:

>An untrained observer will see only physical labor and often get the idea that physical labor is mainly what the mechanic does. Actually the physical labor is the smallest and easiest part of what the mechanic does. By far the greatest part of his work is careful observation and precise thinking.

I was still thinking in terms of systems.  My brain still felt like I was coding, but my hands were like "we're not really doing anything..."  I was thinking a lot more about how the system worked, why things were ordered as they were (and if they should be that way at all) rather than working on getting the exact syntax of my code right.  It was a very interesting shift.  And that's just the coding part of it.

I started to learn about skills and realized I could have claude make github tickets based on error logs.  Or better yet I could take the notes that Gemini makes from Google Meetings, paste that into claude, no wait, I could write skill that just auto downloads those meeting notes, gives me a list of candidate github issues to create, let me review and then make them.  The feedback loop is so much smaller now with meetings.  The action items don't get forgotten.  When some one says, "we should make a ticket for..." the ticket actually gets made.  There's a checklist to review the issue before it's created (does it already exist, are the requirements for new issues met, etc.).  If Claude code did _nothing_ at all regarding writing code and just ingested meeting transcripts and turned that into github issue it would be a huge benefit all by itself.