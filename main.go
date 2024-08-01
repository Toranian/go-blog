package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	fileName := name + ".md"
	if name == "" {
		fileName = "index.md"
		name = "Blog"
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

	md, err := os.ReadFile("blog/" + fileName)
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	htmlContent := template.HTML(mdToHTML(md))
	tmpl, err := template.ParseFiles("components/template.html")

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

func main() {
	const port uint16 = 3000
	portStr := strconv.Itoa(int(port))
	url := "http://localhost:" + portStr

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)

	fmt.Printf("Server running at %s\n", url)
	log.Fatal(http.ListenAndServe(":"+portStr, nil))
}
