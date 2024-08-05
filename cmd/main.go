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
	"strings"

	chromaHTML "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	configure "goblog/internal/configure"
)

// renderParagraph handles the rendering of paragraph nodes
func renderParagraph(w io.Writer, _ *ast.Paragraph, entering bool) {
	if entering {
		io.WriteString(w, "<p>")
	} else {
		io.WriteString(w, "</p>")
	}
}

// func renderHeading(w io.Writer, h *ast.Heading, entering bool, path string) {
//
// 	fmt.Printf("%s", string(h.Level))
// 	if entering {
// 		header := fmt.Sprintf("<a href=\"%s#%s\">", path, h.Content)
// 		io.WriteString(w, header)
// 	} else {
// 		io.WriteString(w, "</a>")
// 	}
// }

// renderHeading renders a heading element, wrapping it with an anchor link.
func renderHeading(w io.Writer, h *ast.Heading, entering bool, path string) {
	if entering {
		// Extract the text content of the heading
		var contentBuilder strings.Builder
		for _, node := range h.Children {
			if textNode, ok := node.(*ast.Text); ok {
				contentBuilder.Write(textNode.Literal)
			}
		}

		fmt.Printf("%d", h.Level)

		content := contentBuilder.String()
		// Generate a valid HTML id for the heading link
		id := strings.ReplaceAll(strings.ToLower(content), " ", "-")
		// header := fmt.Sprintf(`<h%d id="%s"><a href="%s#%s" class="heading-link">`, h.Level, id, path, id)
		header := fmt.Sprintf(`<a href="%s#%s" class="heading-link"><h%d id="%s">`, path, id, h.Level, id)
		io.WriteString(w, header)
	} else {
		io.WriteString(w, fmt.Sprintf("</h%d></a>", h.Level))
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
	case *ast.Heading:
		renderHeading(w, n, entering, "/")
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

	fmt.Println("Request for", fileName)
	missingPage, err := template.ParseFiles("web/components/404.html")

	md, err := os.ReadFile("web/blog/" + fileName)
	if err != nil {
		missingPage.Execute(w, nil)
		return
	}

	html := mdToHTML(md)
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

func main() {
	const port uint16 = 3000
	portStr := strconv.Itoa(int(port))
	url := "http://localhost:" + portStr

	// Load the configuration from the TOML file
	config, err := configure.GetConfigFromTOML()
	if err != nil {
		fmt.Printf("Error loading configuration file. Error: %s", err)
		return
	}
	fmt.Println("Configuration loaded successfully.")

	content := configure.GenerateSCSSVariables(config.CSSVariables)

	// Write the SCSS variables to a file
	err = os.WriteFile("web/static/scss/_variables.scss", []byte(content), 0644)

	if err != nil {
		log.Fatal("Error writing SCSS variables to file.")
		return
	}

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)

	fmt.Printf("Server running at %s\n", url)
	log.Fatal(http.ListenAndServe(":"+portStr, nil))
}
