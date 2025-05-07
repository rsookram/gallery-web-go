package main

import (
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

//go:embed "image.html"
var imageHtml string

type DirectoryMetadata struct {
	Count int
	Hash  string
}

func main() {
	port := flag.Uint("port", 4255, "the port to serve on")
	flag.Parse()

	dir := flag.Arg(0)
	if dir == "" {
		dir = "."
	}

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	slog.Info("starting", "addr", addr, "dir", dir)

	images, err := loadImages(dir)
	if err != nil {
		slog.Error("failed to read", "dir", dir, "error", err)
		os.Exit(1)
	}

	if len(images) == 0 {
		slog.Error("no files", "dir", dir)
		os.Exit(1)
	}

	count := len(images)
	slog.Info("loaded images", "count", count)

	// Provide a hash to prevent browser-level caching from using stale entries
	// after restarting the server with a different dir.
	hashBytes := md5.Sum([]byte(dir))
	hash := hex.EncodeToString(hashBytes[:])

	meta := DirectoryMetadata{
		Count: count,
		Hash:  hash,
	}

	t := template.Must(template.New("image.html").Parse(imageHtml))

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, meta); err != nil {
			slog.Warn("failed to execute template", "meta", meta, "error", err)
		}
	})

	http.HandleFunc("GET /{hash}/image/{index}/", func(w http.ResponseWriter, r *http.Request) {
		index := r.PathValue("index")

		idx, err := strconv.ParseInt(index, 10, 32)
		if err != nil {
			slog.Warn("invalid image index", "path", r.URL.Path, "error", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, images[idx])
	})

	slog.Error("exiting", "error", http.ListenAndServe(addr, nil))
}

func loadImages(dir string) ([]string, error) {
	images := make([]string, 0, 128)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		images = append(images, path)
		return nil
	})

	return images, err
}
