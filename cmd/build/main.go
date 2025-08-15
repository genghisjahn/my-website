package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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

var (
	articleTpl *template.Template
	listTpl    *template.Template

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
	}
	tpl, err := template.New(filepath.Base(path)).Funcs(funcs).Parse(string(b))
	if err != nil {
		log.Fatalf("parse template %s: %v", path, err)
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

func humanDate(t time.Time) string {
	return t.Format("January 2, 2006")
}

type articleView struct {
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
}
type listView struct {
	Title    string
	Subtitle string
	Items    []listItem
}

func main() {
	root := "."
	srcDir := filepath.Join(root, "articles")
	outDir := filepath.Join(root, "public")
	articleTpl = mustTemplate(filepath.Join(root, "templates", "article.html.tmpl"))
	listTpl = mustTemplate(filepath.Join(root, "templates", "list.html.tmpl"))

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
			arts[i].Prev = &arts[i-1]
		}
		if i < len(arts)-1 {
			arts[i].Next = &arts[i+1]
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
		av := articleView{
			Title:        a.Title,
			Date:         a.Date,
			DateHuman:    humanDate(a.t),
			Author:       a.Author,
			Tags:         a.Tags,
			ContentHTML:  template.HTML(a.ContentHTML),
			CanonicalURL: a.CanonicalURL,
			Hero:         a.Hero,
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
