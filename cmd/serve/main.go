package main

import (
	"compress/gzip"
	"context"
	"flag"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func init() {
	_ = mime.AddExtensionType(".webmanifest", "application/manifest+json")
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	publicDir := flag.String("public", "./public", "public dir")
	cssDir := flag.String("css", "./css", "css dir")
	imagesDir := flag.String("images", "./images", "images dir")
	flag.Parse()

	mux := http.NewServeMux()

	// / -> public
	mux.Handle("/", gzipWrap(logWrap(cacheWrap(http.FileServer(http.Dir(*publicDir))))))

	// /css -> css
	mux.Handle("/css/",
		gzipWrap(logWrap(cacheWrap(http.StripPrefix("/css/",
			http.FileServer(http.Dir(*cssDir)))))))

	// /images -> images (if present)
	if dirExists(*imagesDir) {
		mux.Handle("/images/",
			logWrap(cacheWrap(http.StripPrefix("/images/",
				http.FileServer(http.Dir(*imagesDir))))))
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Serving %s at http://localhost%s", abs(*publicDir), *addr)
	log.Printf("Mount /css -> %s", abs(*cssDir))
	if dirExists(*imagesDir) {
		log.Printf("Mount /images -> %s", abs(*imagesDir))
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("Server stopped")
}

func cacheWrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hasExt(r.URL.Path, ".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".ico", ".woff", ".woff2", ".xml", ".webmanifest") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			w.Header().Set("Expires", time.Now().AddDate(1, 0, 0).UTC().Format(http.TimeFormat))
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		next.ServeHTTP(w, r)
	})
}

func logWrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func hasExt(p string, exts ...string) bool {
	for _, e := range exts {
		if strings.HasSuffix(strings.ToLower(p), e) {
			return true
		}
	}
	return false
}

func dirExists(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

func abs(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return a
}

// gzip compression middleware
func gzipWrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if client doesn't accept gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		// Only compress compressible types
		if !isCompressible(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length") // Length changes after compression

		gz := gzip.NewWriter(w)
		defer gz.Close()

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

func isCompressible(path string) bool {
	return hasExt(path, ".html", ".css", ".js", ".json", ".xml", ".svg", ".txt", ".webmanifest")
}
