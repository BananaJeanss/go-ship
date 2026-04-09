package handlers

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
)

func GuidesHandler(w http.ResponseWriter, r *http.Request) {
	source := strings.TrimSpace(r.URL.Query().Get("guide")) // e.g. ?guide=some-guide.md

	// fallback for URLs like /guides?hi
	if source == "" {
		raw := strings.TrimSpace(r.URL.RawQuery)
		if raw != "" && !strings.Contains(raw, "=") && !strings.Contains(raw, "&") {
			source = raw
		}
	}

	if source == "" {
		http.Error(w, "Guide not specified", http.StatusBadRequest)
		return
	}

	source = filepath.Base(source)
	if filepath.Ext(source) == "" {
		source += ".md"
	}

	// verify guide actually exists as md file in public/guides/
	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	sourcePath := filepath.Join(cwd, "public", "guides", source)
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		http.Error(w, "Guide not found", http.StatusNotFound)
		return
	}

	// read the md file
	mdContent, err := os.ReadFile(sourcePath)
	if err != nil {
		http.Error(w, "Failed to read guide", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(mdContent, &buf); err != nil {
		http.Error(w, "Failed to render guide", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/guides/index.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	type GuideData struct {
		Title   string
		Content string
	}

	if err := tmpl.Execute(w, GuideData{
		Title:   source,
		Content: buf.String(),
	}); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
