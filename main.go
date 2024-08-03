package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	chromaHTML "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// renderParagraph handles the rendering of paragraph nodes
func renderParagraph(w io.Writer, _ *ast.Paragraph, entering bool) {
	if entering {
		io.WriteString(w, "<p>")
	} else {
		io.WriteString(w, "</p>")
	}
}

// Render a code block using Chroma to add styling.
func renderCodeBlock(w io.Writer, c *ast.CodeBlock, entering bool) {
	if entering {
		// Extract the programming language from the Info field
		language := string(c.Info)

		// Get the appropriate lexer for the language
		lexer := lexers.Get(language)
		if lexer == nil {
			lexer = lexers.Fallback
		}

		// Use Chroma to highlight the code block
		iterator, err := lexer.Tokenise(nil, string(c.Literal))
		if err != nil {
			io.WriteString(w, `<div class="code-block"><pre><code>`)
			w.Write(c.Literal)
			io.WriteString(w, `</code></pre></div>`)
			return
		}

		formatter := chromaHTML.New(chromaHTML.WithClasses(true))
		style := styles.Get("github")
		if style == nil {
			style = styles.Fallback
		}

		// Write the opening tags with the language class
		if language != "" {
			io.WriteString(w, fmt.Sprintf(`<div class="code-block"><pre><code class="language-%s">`, language))
		} else {
			io.WriteString(w, `<div class="code-block"><pre><code>`)
		}

		// Write the highlighted code
		err = formatter.Format(w, style, iterator)
		if err != nil {
			w.Write(c.Literal)
		}

		io.WriteString(w, `</code></pre></div>`)
	}
}

// myRenderHook is the custom render hook for handling different node types
func customRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch n := node.(type) {
	case *ast.Paragraph:
		renderParagraph(w, n, entering)
		return ast.GoToNext, true
	case *ast.CodeBlock:
		renderCodeBlock(w, n, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func newCustomizedRender() *html.Renderer {
	opts := html.RendererOptions{
		RenderNodeHook: customRenderHook,
	}
	return html.NewRenderer(opts)
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Tables | parser.FencedCode | parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.HeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.LazyLoadImages
	opts := html.RendererOptions{Flags: htmlFlags, RenderNodeHook: customRenderHook}
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

	missingPage, err := template.ParseFiles("components/404.html")

	md, err := os.ReadFile("blog/" + fileName)
	if err != nil {
		missingPage.Execute(w, nil)
		return
	}

	html := mdToHTML(md)
	htmlContent := template.HTML(html)

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
