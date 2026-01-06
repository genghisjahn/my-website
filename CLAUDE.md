# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Build the static site (outputs to public/)
go run ./cmd/build

# Run local development server (serves public/ on port 8080)
go run ./cmd/serve

# Full deploy (build + WebP conversion + rsync to server)
./deploy.sh

# Deploy just the server binary (rebuilds and restarts remote Go server)
./server_deploy.sh
```

## Architecture

This is a custom static site generator written in Go that produces a personal website with articles and notes.

**Two main binaries:**
- `cmd/build/main.go` - Static site generator that reads content and produces HTML in `public/`
- `cmd/serve/main.go` - Development HTTP server with gzip, caching, and custom 404 handling

**Content format:**
- Articles: `articles/*.md` - Markdown files with YAML frontmatter (slug, title, date, author, tags, etc.)
- Notes: `notes/*.md` - Shorter Markdown posts with similar frontmatter
- Both support a `draft: true` flag to exclude from builds

**Templates:** Go HTML templates in `templates/` with partials (styles, favicons, feeds, webmention, theme-toggle)

**Output structure:**
- `public/articles/{slug}/index.html` - Individual article pages
- `public/notes/{slug}/index.html` - Individual note pages
- `public/tag/{slug}/` - Tag archive pages
- `public/archive/{YYYY/MM}/` - Monthly archive pages
- `public/feed.xml` and `public/notes/feed.xml` - RSS feeds

**Image handling:** Build process converts PNG/JPG images to WebP and rewrites `<img src>` paths automatically. Deploy script runs `cwebp` on images in `public/images/`.

**CSS cache busting:** When updating `css/retro-sci-fi.css`, increment the `?v=` query parameter in `templates/styles.html.tmpl` to bust browser caches.

**Dependencies:** goldmark (Markdown), gopkg.in/yaml.v3 (YAML frontmatter)
