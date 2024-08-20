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
)

type BlogHandler struct {
	// config is the configuration for the blog
	Config configure.Config
}

func (h *BlogHandler) Handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	fileName := name + ".md"
	config := h.Config
	if name == "" {
		fileName = config.BlogData.SourceFile
		name = config.MetaData.Title
	} else {
		// Set the name to the title of the blog post
		// Extract the path from the URL
		parsedURL, err := url.Parse(name)
		if err != nil {
			panic(err)
		}
		urlPath := parsedURL.Path

		// Get the last part of the path
		lastPart := path.Base(urlPath)
		name = lastPart
	}

	missingPage, err := template.ParseFiles("web/components/404.html")

	fmt.Println("Looking for: ", config.BlogData.ContentDirectory+fileName)

	md, err := os.ReadFile(config.BlogData.ContentDirectory + "/" + fileName)

	if err != nil {
		missingPage.Execute(w, nil)
		return
	}

	html := renderMD.MDToHTML(md)
	htmlContent := template.HTML(html)

	tmpl, err := template.ParseFiles("web/components/template.html")

	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Body  template.HTML
	}{
		Body:  htmlContent,
		Title: name,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}
