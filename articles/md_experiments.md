---
slug: markdown-experiments
title: Markdown Experiments
date: 2025-08-15
author:
  name: Jon Wear
summary: Let's see what works and doesn't work using markdown.
tags:
  - name: markdown
    slug: markdown
hero:
  src: /images/hero_article_3.png
  alt: Corridor lights
css: /css/retro-sci-fi.css
draft: false
reading_time_min: 2
---

First off, let's so some *basic* things with text to see if **markdown** renders it correctly.  So far so good.  I'll put that in <u>underline</u>? I want to see what we can ~~strike out~~.  Lastly, I want to see if we can do a link and or an image.  The link will be to [jonwear.com](https://jonwear.com) and the image will be <img src="/images/hero2.png" alt="hero image" width="400">.

> Let's look at quotes and things like that.  This is where I'm doing lots of quoting about something someone said somewhere.  It should be obvious it's a quote.  Let's see how it looks.

```go
    var md = goldmark.New(
        goldmark.WithExtensions(extension.Strikethrough),
        goldmark.WithRendererOptions(
            html.WithUnsafe(),
        ),
    )
```

`This could be kind of set apart.`

```js

alert('Hello!');

```

That is all, let's see how it renders.

1. We also need a list of things
1. Here's the second thing
1. Here's the 3rd thing.
    1. Here's a sub item
    1. Here's another one.
1. Back to the main list.

* This is not numbered.
* This is just bullets.
    * This is a sub bullet
* This is a regular bullet

| Player Name| Avg. |
|----------|----------|
| Jon Wear    | .384     |
| Bob "Whif" Johnson    | .221     |


- [x] Task 1
- [ ] Task 2