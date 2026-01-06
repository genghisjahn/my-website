package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

func dirExists(p string) bool { fi, err := os.Stat(p); return err == nil && fi.IsDir() }

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		out := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(out, 0o755)
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(out, b, 0o644)
	})
}

type Author struct {
	Name string  `json:"name"`
	URL  *string `json:"url"`
}
type Tag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}
type Hero struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}
type Article struct {
	Slug           string  `json:"slug"`
	Title          string  `json:"title"`
	Subtitle       *string `json:"subtitle"`
	Date           string  `json:"date"` // YYYY-MM-DD
	Updated        *string `json:"updated"`
	Author         Author  `json:"author"`
	Summary        *string `json:"summary"`
	Tags           []Tag   `json:"tags"`
	Hero           *Hero   `json:"hero"`
	CanonicalURL   *string `json:"canonical_url"`
	CSS            *string `json:"css"`
	Draft          bool    `json:"draft"`
	ReadingTimeMin *int    `json:"reading_time_min"`
	ContentHTML    string  `json:"content_html"`
	// derived
	t    time.Time
	Prev *Article `json:"-"`
	Next *Article `json:"-"`
}

type markdownArticle struct {
	Slug           string  `yaml:"slug"`
	Title          string  `yaml:"title"`
	Subtitle       *string `yaml:"subtitle"`
	Date           string  `yaml:"date"`
	Updated        *string `yaml:"updated"`
	Author         Author  `yaml:"author"`
	Summary        *string `yaml:"summary"`
	Tags           []Tag   `yaml:"tags"`
	Hero           *Hero   `yaml:"hero"`
	CanonicalURL   *string `yaml:"canonical_url"`
	CSS            *string `yaml:"css"`
	Draft          bool    `yaml:"draft"`
	ReadingTimeMin *int    `yaml:"reading_time_min"`
}

// Note represents a short public note (like a gist)
type Note struct {
	Slug        string  `yaml:"slug"`
	Title       string  `yaml:"title"`
	Date        string  `yaml:"date"` // YYYY-MM-DD or YYYY-MM-DDTHH:MM
	Author      Author  `yaml:"author"`
	Tags        []Tag   `yaml:"tags"`
	Source      *string `yaml:"source"` // optional: URL, book name, or person
	Draft       bool    `yaml:"draft"`
	ContentHTML string
	t           time.Time
}

var (
	articleTpl  *template.Template
	listTpl     *template.Template
	noteTpl     *template.Template
	noteListTpl *template.Template

	reScriptStyle = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>|<style[^>]*>.*?</style>`)
	reTags        = regexp.MustCompile(`(?s)<[^>]+>`)
	reSpace       = regexp.MustCompile(`\s+`)
)

func mustTemplate(path string) *template.Template {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read template %s: %v", path, err)
	}
	// allow raw HTML insertion via template.HTML
	funcs := template.FuncMap{
		"split": strings.Split,
		"isURL": func(s *string) bool {
			if s == nil {
				return false
			}
			return strings.HasPrefix(*s, "http://") || strings.HasPrefix(*s, "https://")
		},
		"deref": func(s *string) string {
			if s == nil {
				return ""
			}
			return *s
		},
	}
	tpl, err := template.New(filepath.Base(path)).Funcs(funcs).Parse(string(b))
	if err != nil {
		log.Fatalf("parse template %s: %v", path, err)
	}
	// Parse partials
	partials := []string{"styles.html.tmpl", "favicons.html.tmpl", "feeds.html.tmpl", "webmention.html.tmpl", "theme-toggle.html.tmpl", "nav.html.tmpl"}
	for _, partial := range partials {
		partialPath := filepath.Join(filepath.Dir(path), partial)
		if _, err := os.Stat(partialPath); err == nil {
			partialB, err := os.ReadFile(partialPath)
			if err != nil {
				log.Fatalf("read partial %s: %v", partial, err)
			}
			_, err = tpl.Parse(string(partialB))
			if err != nil {
				log.Fatalf("parse partial %s: %v", partial, err)
			}
		}
	}
	return tpl
}

func readingTimeMinutes(contentHTML string) int {
	t := reScriptStyle.ReplaceAllString(contentHTML, "")
	t = reTags.ReplaceAllString(t, "")
	t = strings.TrimSpace(reSpace.ReplaceAllString(t, " "))
	if t == "" {
		return 1
	}
	words := len(strings.Fields(t))
	mins := (words + 219) / 220 // ceil
	if mins < 1 {
		mins = 1
	}
	return mins
}

func mustParseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Fatalf("bad date %q: %v", s, err)
	}
	return t
}

func mustParseDateTime(s string) time.Time {
	// Try datetime format first (YYYY-MM-DDTHH:MM)
	t, err := time.Parse("2006-01-02T15:04", s)
	if err == nil {
		return t
	}
	// Fall back to date-only format (YYYY-MM-DD)
	t, err = time.Parse("2006-01-02", s)
	if err != nil {
		log.Fatalf("bad date/datetime %q: %v", s, err)
	}
	return t
}

func humanDate(t time.Time) string {
	return t.Format("January 2, 2006")
}

type articleView struct {
	Slug         string
	Title        string
	Date         string
	DateHuman    string
	Author       Author
	Tags         []Tag
	ContentHTML  template.HTML
	CanonicalURL *string
	Hero         *Hero
	Prev         *Article
	Next         *Article
}

type listItem struct {
	Title     string
	URL       string
	ISODate   string
	HumanDate string
	Type      string // "article" or "note"
}
type listView struct {
	Title    string
	Subtitle string
	Items    []listItem
}

type noteView struct {
	Slug        string
	Title       string
	Date        string
	DateHuman   string
	Author      Author
	Tags        []Tag
	Source      *string
	ContentHTML template.HTML
}

type paginatedListView struct {
	Title       string
	Subtitle    string
	Items       []listItem
	CurrentPage int
	TotalPages  int
	PrevURL     string
	NextURL     string
}

// RSS feed types
type rssChannel struct {
	XMLName       xml.Name  `xml:"channel"`
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

func writeRSSFeed(outPath, title, link, description string, items []rssItem) error {
	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:         title,
			Link:          link,
			Description:   description,
			Language:      "en-us",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         items,
		},
	}
	buf := new(bytes.Buffer)
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(buf)
	enc.Indent("", "  ")
	if err := enc.Encode(feed); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, buf.Bytes(), 0o644)
}

func main() {
	root := "."
	srcDir := filepath.Join(root, "articles")
	notesSrcDir := filepath.Join(root, "notes")
	outDir := filepath.Join(root, "public")
	articleTpl = mustTemplate(filepath.Join(root, "templates", "article.html.tmpl"))
	listTpl = mustTemplate(filepath.Join(root, "templates", "list.html.tmpl"))
	noteTpl = mustTemplate(filepath.Join(root, "templates", "note.html.tmpl"))
	noteListTpl = mustTemplate(filepath.Join(root, "templates", "note_list.html.tmpl"))

	var arts []Article
	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if strings.HasSuffix(d.Name(), ".json") {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var a Article
			if err := json.Unmarshal(b, &a); err != nil {
				return err
			}
			if a.Draft {
				return nil
			}
			a.t = mustParseDate(a.Date)
			if a.ReadingTimeMin == nil || *a.ReadingTimeMin < 1 {
				rt := readingTimeMinutes(a.ContentHTML)
				a.ReadingTimeMin = &rt
			}
			arts = append(arts, a)
			return nil
		}
		if strings.HasSuffix(d.Name(), ".md") {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			content := string(b)
			parts := strings.SplitN(content, "---", 3)
			if len(parts) < 3 {
				return nil // skip invalid
			}
			var meta markdownArticle
			if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
				return err
			}
			if meta.Draft {
				return nil
			}
			htmlBuf := new(bytes.Buffer)

			var md = goldmark.New(
				goldmark.WithExtensions(
					extension.Strikethrough,
					extension.Table,
					extension.TaskList,
					extension.Footnote,
				),
				goldmark.WithRendererOptions(
					html.WithUnsafe(),
					html.WithXHTML(),
				),
			)

			if err := md.Convert([]byte(parts[2]), htmlBuf); err != nil {
				return err
			}

			htmlStr := htmlBuf.String()

			a := Article{
				Slug:           meta.Slug,
				Title:          meta.Title,
				Subtitle:       meta.Subtitle,
				Date:           meta.Date,
				Updated:        meta.Updated,
				Author:         meta.Author,
				Summary:        meta.Summary,
				Tags:           meta.Tags,
				Hero:           meta.Hero,
				CanonicalURL:   meta.CanonicalURL,
				CSS:            meta.CSS,
				Draft:          false,
				ReadingTimeMin: meta.ReadingTimeMin,
				ContentHTML:    htmlStr,
				t:              mustParseDate(meta.Date),
			}
			if a.ReadingTimeMin == nil || *a.ReadingTimeMin < 1 {
				rt := readingTimeMinutes(a.ContentHTML)
				a.ReadingTimeMin = &rt
			}
			arts = append(arts, a)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(arts, func(i, j int) bool { return arts[i].t.After(arts[j].t) })

	for i := range arts {
		if i > 0 {
			arts[i].Next = &arts[i-1]
		}
		if i < len(arts)-1 {
			arts[i].Prev = &arts[i+1]
		}
	}

	// Prepare maps
	tagMap := map[string]struct {
		Name  string
		Items []listItem
	}{}
	ymMap := map[string][]listItem{} // key "YYYY/MM"

	// Render articles
	for _, a := range arts {
		// Convert hero image to WebP
		var heroWebP *Hero
		if a.Hero != nil {
			heroWebP = &Hero{
				Src: toWebP(a.Hero.Src),
				Alt: a.Hero.Alt,
			}
		}
		av := articleView{
			Slug:         a.Slug,
			Title:        a.Title,
			Date:         a.Date,
			DateHuman:    humanDate(a.t),
			Author:       a.Author,
			Tags:         a.Tags,
			ContentHTML:  template.HTML(convertContentImagesToWebP(a.ContentHTML)),
			CanonicalURL: a.CanonicalURL,
			Hero:         heroWebP,
			Prev:         a.Prev,
			Next:         a.Next,
		}
		out := new(bytes.Buffer)
		if err := articleTpl.Execute(out, av); err != nil {
			log.Fatalf("render article %s: %v", a.Slug, err)
		}
		adir := filepath.Join(outDir, "articles", a.Slug)
		if err := os.MkdirAll(adir, 0o755); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(adir, "index.html"), out.Bytes(), 0o644); err != nil {
			log.Fatal(err)
		}

		item := listItem{
			Title:     a.Title,
			URL:       "/articles/" + a.Slug + "/",
			ISODate:   a.Date,
			HumanDate: humanDate(a.t),
			Type:      "article",
		}

		// tags
		for _, tg := range a.Tags {
			entry := tagMap[tg.Slug]
			entry.Name = tg.Name
			entry.Items = append(entry.Items, item)
			tagMap[tg.Slug] = entry
		}

		// archive buckets
		ym := a.t.Format("2006/01")
		ymMap[ym] = append(ymMap[ym], item)
	}

	// Process notes
	var notes []Note
	if dirExists(notesSrcDir) {
		err = filepath.WalkDir(notesSrcDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			if !strings.HasSuffix(d.Name(), ".md") {
				return nil
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			content := string(b)
			parts := strings.SplitN(content, "---", 3)
			if len(parts) < 3 {
				return nil // skip invalid
			}
			var note Note
			if err := yaml.Unmarshal([]byte(parts[1]), &note); err != nil {
				return err
			}
			if note.Draft {
				return nil
			}
			htmlBuf := new(bytes.Buffer)
			md := goldmark.New(
				goldmark.WithExtensions(
					extension.Strikethrough,
					extension.Table,
					extension.TaskList,
					extension.Footnote,
				),
				goldmark.WithRendererOptions(
					html.WithUnsafe(),
					html.WithXHTML(),
				),
			)
			if err := md.Convert([]byte(parts[2]), htmlBuf); err != nil {
				return err
			}
			note.ContentHTML = htmlBuf.String()
			note.t = mustParseDateTime(note.Date)
			notes = append(notes, note)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	sort.Slice(notes, func(i, j int) bool { return notes[i].t.After(notes[j].t) })

	// Render notes
	var noteItems []listItem
	for _, n := range notes {
		nv := noteView{
			Slug:        n.Slug,
			Title:       n.Title,
			Date:        n.Date,
			DateHuman:   humanDate(n.t),
			Author:      n.Author,
			Tags:        n.Tags,
			Source:      n.Source,
			ContentHTML: template.HTML(convertContentImagesToWebP(n.ContentHTML)),
		}
		out := new(bytes.Buffer)
		if err := noteTpl.Execute(out, nv); err != nil {
			log.Fatalf("render note %s: %v", n.Slug, err)
		}
		ndir := filepath.Join(outDir, "notes", n.Slug)
		if err := os.MkdirAll(ndir, 0o755); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(ndir, "index.html"), out.Bytes(), 0o644); err != nil {
			log.Fatal(err)
		}

		item := listItem{
			Title:     n.Title,
			URL:       "/notes/" + n.Slug + "/",
			ISODate:   n.Date,
			HumanDate: humanDate(n.t),
			Type:      "note",
		}
		noteItems = append(noteItems, item)

		// add to tags
		for _, tg := range n.Tags {
			entry := tagMap[tg.Slug]
			entry.Name = tg.Name
			entry.Items = append(entry.Items, item)
			tagMap[tg.Slug] = entry
		}
	}

	// Render tag pages
	for slug, v := range tagMap {
		sort.Slice(v.Items, func(i, j int) bool { return v.Items[i].ISODate > v.Items[j].ISODate })
		lv := listView{Title: "Tag: " + v.Name, Items: v.Items}
		writeList(listTpl, filepath.Join(outDir, "tag", slug, "index.html"), lv)
	}

	// Render archive month pages and archive index
	type ymEntry struct {
		Key   string
		Items []listItem
	}
	var months []ymEntry
	for k, items := range ymMap {
		sort.Slice(items, func(i, j int) bool { return items[i].ISODate > items[j].ISODate })
		months = append(months, ymEntry{Key: k, Items: items})
	}
	sort.Slice(months, func(i, j int) bool { return months[i].Key > months[j].Key })

	// month pages
	for _, m := range months {
		title := "Archive " + humanMonth(m.Key)
		lv := listView{Title: title, Items: m.Items}
		writeList(listTpl, filepath.Join(outDir, "archive", m.Key, "index.html"), lv)
	}
	// archive index
	var idxItems []listItem
	for _, m := range months {
		idxItems = append(idxItems, listItem{
			Title:     humanMonth(m.Key),
			URL:       "/archive/" + m.Key + "/",
			ISODate:   m.Key,
			HumanDate: humanMonth(m.Key),
		})
	}
	writeList(listTpl, filepath.Join(outDir, "archive", "index.html"), listView{
		Title:    "Archive",
		Subtitle: "By month",
		Items:    idxItems,
	})

	// Render notes list with pagination
	const notesPerPage = 20
	totalNotePages := (len(noteItems) + notesPerPage - 1) / notesPerPage
	if totalNotePages < 1 {
		totalNotePages = 1
	}
	for page := 1; page <= totalNotePages; page++ {
		start := (page - 1) * notesPerPage
		end := start + notesPerPage
		if end > len(noteItems) {
			end = len(noteItems)
		}
		pageItems := noteItems[start:end]

		var prevURL, nextURL string
		if page > 1 {
			if page == 2 {
				prevURL = "/notes/"
			} else {
				prevURL = "/notes/page/" + strconv.Itoa(page-1) + "/"
			}
		}
		if page < totalNotePages {
			nextURL = "/notes/page/" + strconv.Itoa(page+1) + "/"
		}

		plv := paginatedListView{
			Title:       "Notes",
			Subtitle:    "Quick reference notes",
			Items:       pageItems,
			CurrentPage: page,
			TotalPages:  totalNotePages,
			PrevURL:     prevURL,
			NextURL:     nextURL,
		}

		var outPath string
		if page == 1 {
			outPath = filepath.Join(outDir, "notes", "index.html")
		} else {
			outPath = filepath.Join(outDir, "notes", "page", strconv.Itoa(page), "index.html")
		}

		buf := new(bytes.Buffer)
		if err := noteListTpl.Execute(buf, plv); err != nil {
			log.Fatalf("render notes list page %d: %v", page, err)
		}
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
			log.Fatal(err)
		}
	}

	// Generate RSS feeds
	const siteURL = "https://jonwear.com"

	// Posts RSS feed
	var postRSSItems []rssItem
	for _, a := range arts {
		postRSSItems = append(postRSSItems, rssItem{
			Title:       a.Title,
			Link:        siteURL + "/articles/" + a.Slug + "/",
			Description: a.ContentHTML,
			PubDate:     a.t.Format(time.RFC1123Z),
			GUID:        siteURL + "/articles/" + a.Slug + "/",
		})
	}
	if err := writeRSSFeed(
		filepath.Join(outDir, "feed.xml"),
		"jonwear.com - Posts",
		siteURL,
		"Articles and posts from jonwear.com",
		postRSSItems,
	); err != nil {
		log.Fatalf("write posts RSS: %v", err)
	}

	// Notes RSS feed
	var noteRSSItems []rssItem
	for _, n := range notes {
		noteRSSItems = append(noteRSSItems, rssItem{
			Title:       n.Title,
			Link:        siteURL + "/notes/" + n.Slug + "/",
			Description: n.ContentHTML,
			PubDate:     n.t.Format(time.RFC1123Z),
			GUID:        siteURL + "/notes/" + n.Slug + "/",
		})
	}
	if err := writeRSSFeed(
		filepath.Join(outDir, "notes", "feed.xml"),
		"jonwear.com - Notes",
		siteURL+"/notes/",
		"Quick reference notes from jonwear.com",
		noteRSSItems,
	); err != nil {
		log.Fatalf("write notes RSS: %v", err)
	}

	// simple home index (latest N)
	var homeItems []listItem
	for i, it := range allItems(arts) {
		if i >= 12 {
			break
		}
		homeItems = append(homeItems, it)
	}
	writeList(listTpl, filepath.Join(outDir, "index.html"), listView{
		Title:    "jonwear.com",
		Subtitle: "jon@jonwear.com",
		Items:    homeItems,
	})

	// Generate 404 page
	tpl404 := mustTemplate(filepath.Join(root, "templates", "404.html.tmpl"))
	buf404 := new(bytes.Buffer)
	if err := tpl404.Execute(buf404, nil); err != nil {
		log.Fatalf("render 404: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "404.html"), buf404.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}

	if dirExists("css") {
		if err := copyDir("css", filepath.Join(outDir, "css")); err != nil {
			log.Fatal(err)
		}
	}
	if dirExists("images") {
		if err := copyDir("images", filepath.Join(outDir, "images")); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Build complete -> public/")
}

func allItems(arts []Article) []listItem {
	var items []listItem
	for _, a := range arts {
		items = append(items, listItem{
			Title:     a.Title,
			URL:       "/articles/" + a.Slug + "/",
			ISODate:   a.Date,
			HumanDate: humanDate(a.t),
		})
	}
	return items
}

func writeList(tpl *template.Template, outPath string, lv listView) {
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, lv); err != nil {
		log.Fatalf("render list %s: %v", outPath, err)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}

func humanMonth(key string) string {
	// key "YYYY/MM"
	t, err := time.Parse("2006/01", key)
	if err != nil {
		return key
	}
	return t.Format("January 2006")
}

// toWebP converts image path extensions to .webp
func toWebP(path string) string {
	for _, ext := range []string{".png", ".PNG", ".jpg", ".JPG", ".jpeg", ".JPEG"} {
		if strings.HasSuffix(path, ext) {
			return strings.TrimSuffix(path, ext) + ".webp"
		}
	}
	return path
}

// convertContentImagesToWebP replaces image extensions in HTML content
var reImageSrc = regexp.MustCompile(`(src=["']/images/[^"']+)\.(png|PNG|jpg|JPG|jpeg|JPEG)(["'])`)

func convertContentImagesToWebP(html string) string {
	return reImageSrc.ReplaceAllString(html, "${1}.webp${3}")
}
