---
slug: ai-journey-3
title: AI Journey Part 3
date: 2026-07-28
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
  src: /images/hero_article_7.png
  alt: AI Part 3
css: /css/retro-sci-fi.css
draft: false
---

And then I met Claude Code.  It had been out a while but my new job leveraged it a lot and I figured out I might as well get to know the tool.  I'd started using it around November of 2025, but not really doing a lot with it.  Still writing a lot of code by hand.  I would use it and check its outputs and then I would _type_ the code in myself.  I just didn't trust it at first.  But after a while I realized that if I gave better prompts (and by better I mean smaller focus with more detail) it could do the "typing" faster than I could.  I just needed to make sure I covered tests, read over the code and basically understood what I was doing.  It felt dangerous...and scary.  But I noticed I was spending more and more time typing in the terminal window at the bottom of VS Code and less time looking at the big windows of the code editor that actually contained the source code.  Most of my code review was happening in GitHub Pull Requests.  

It was some time in mid January 2026 that I tiled six terminal windows across my widescreen and never opened VS Code at all.  It was one of my most productive days ever.  Not in the LOC written sense, but in features shipped.  I was still very nervous about this but I realized that something had changed fundamentally and the way I've coded in the past is not going to be the way I code in the future.  Or...maybe it is.  But it's not going to look the same to an outsider.  

I had spent the day coding without opening a code editor.  I just sat there for a while staring at my screen, wondering what had just happened.  I was reminded of this quote from Zen and the Art of Motorcycle Maintenance:

>An untrained observer will see only physical labor and often get the idea that physical labor is mainly what the mechanic does. Actually the physical labor is the smallest and easiest part of what the mechanic does. By far the greatest part of his work is careful observation and precise thinking.

I was still thinking in terms of systems.  My brain still felt like I was coding, but my hands were like "we're not really doing anything..."  I was thinking a lot more about how the system worked, why things were ordered as they were (and if they should be that way at all) rather than working on getting the exact syntax of my code right.  It was a very interesting shift.  And that's just the coding part of it.

I started to learn about skills and realized I could have Claude make GitHub tickets based on error logs.  Or better yet I could take the notes that Gemini makes from Google Meetings, paste that into Claude, no wait, I could write a skill that just auto downloads those meeting notes, gives me a list of candidate GitHub issues to create, let me review and then make them.  The feedback loop is so much smaller now with meetings.  The action items don't get forgotten.  When some one says, "we should make a ticket for..." the ticket actually gets made.  There's a checklist to review the issue before it's created (does it already exist, are the requirements for new issues met, etc.).  If Claude Code did _nothing_ at all regarding writing code and just ingested meeting transcripts and turned those into GitHub issues it would be a huge benefit all by itself.

And then I could ask Claude about the history of a function, the history of a stored procedure.  Anything that was committed to GitHub I could now query from the command line and get back a reasonably accurate answer.  I say reasonably because asking the architect that built the systems doesn't guarantee an accurate answer either.  Sometimes they misremember, sometimes another dev made a change they weren't aware of.  Sometimes the code they wrote doesn't do what they thought it did.  What I'm saying is, LLMs do get things wrong, but it's not like we biological language models got it right all the time either.  Plus the LLM doesn't get annoyed when I ask it the same question over and over or get defensive when I point out it got something wrong.

But wait!  There's _more_.  Since Claude has access to GitHub, it can do things like track when an issue is created, when it's assigned, and if we link it to a PR, we can then see when the first commit happened and track all the way through to when it gets merged and deployed.  Give it access to a few weeks of that and then we can say "Hey Claude, how long will it take to do this new GitHub issue" and it has a pretty good idea.  We can do more accurate forecasting for the weeks ahead and be able to say when we think issues will actually land in production.  Sure it's not exactly right, but it's way better than sitting around figuring out story/function points and bikeshedding a planning poker session.  

Anyway, I thought it would be helpful, for myself at least, to document how I got to the inflection point.  When I got to the "there's no going back" realization.  For me it was January of 2026.  It started out with pasting incorrect dynamoDB queries into chatGPT and it ended with Claude Code being able to take meeting notes, turn those into GitHub issues and then having a swarm of agents take on the job of building the spec, build the backend, build the front end, and fill in missing tests and check that the final solution met the spec, repeat until we got to a _ship it_ verdict that a human could look at and say yes or no.

And there's still plenty for a human developer to do in all this.  I don't know what the future holds, but it's going to be interesting.  Anyway, future posts will be about specifics of what I'm doing in this space and what I'm actually shipping.  Til then.

All text written by me with _minor_ spelling/grammar changes from Claude.  All images created by chatGPT.