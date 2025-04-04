package handlers

import (
	"fmt"
	"goblog/internal/configure"
	renderMD "goblog/internal/markdown"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type BlogHandler struct {
	Config configure.Config
}

func (h *BlogHandler) Handler(w http.ResponseWriter, r *http.Request) {
	config := h.Config
	requestPath := r.URL.Path

	// Strip the leading slash
	trimmedPath := requestPath[1:]
	filePath := filepath.Join(config.BlogData.ContentDirectory, trimmedPath)

	// Default to index.md if root
	if trimmedPath == "" {
		filePath = filepath.Join(config.BlogData.ContentDirectory, config.BlogData.SourceFile)
	}

	// Try to stat the file/folder
	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		// If it's a directory, look for index.md
		filePath = filepath.Join(filePath, "index.md")
	} else if err == nil && !info.IsDir() {
		// It's a file, check if it ends with .md
		if filepath.Ext(filePath) == "" {
			filePath += ".md"
		}
	} else if os.IsNotExist(err) {
		// If file or directory doesn't exist, try adding .md
		filePath += ".md"
	}

	// Final check if the file actually exists
	mdBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Markdown file not found:", filePath)
		tmpl404, _ := template.ParseFiles("web/components/404.html")
		tmpl404.Execute(w, nil)
		return
	}

	html := renderMD.MDToHTML(mdBytes)
	htmlContent := template.HTML(html)

	tmpl, err := template.ParseFiles("web/components/template.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	// Use the last part of the path as the title
	parsedURL, _ := url.Parse(requestPath)
	pageTitle := path.Base(parsedURL.Path)
	if pageTitle == "" {
		pageTitle = config.MetaData.Title
	}

	data := struct {
		Title string
		Body  template.HTML
	}{
		Title: pageTitle,
		Body:  htmlContent,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}
