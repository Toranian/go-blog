package markdown

import (
	"fmt"
	"io"
	"strings"

	chromaHTML "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// renderHeading renders a heading element, wrapping it with an anchor link.
func renderHeading(w io.Writer, h *ast.Heading, entering bool) {
	if entering {
		// Extract the text content of the heading
		var contentBuilder strings.Builder
		for _, node := range h.Children {
			if textNode, ok := node.(*ast.Text); ok {
				contentBuilder.Write(textNode.Literal)
			}
		}

		content := contentBuilder.String()

		// Generate a valid HTML id for the heading link
		id := strings.ReplaceAll(strings.ToLower(content), " ", "-")

		// header := fmt.Sprintf(`<h%d id="%s"><a href="%s#%s" class="heading-link">`, h.Level, id, path, id)
		header := fmt.Sprintf(`<a href="#%s" class="heading-link"><h%d id="%s">`, id, h.Level, id)
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
	case *ast.CodeBlock:
		renderCodeBlock(w, n, entering)
		return ast.GoToNext, true
	case *ast.Heading:
		renderHeading(w, n, entering)
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

func MDToHTML(md []byte) []byte {
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
