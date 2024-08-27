# Go Blog

This blogging application was created to leverage the ease of use of Markdown documents, with the speed of Go, and the simplicity of basic styling.

![Beach Image](https://lh3.googleusercontent.com/pw/AP1GczM47iCWwqJF9EMK4uXT6bvUUGilIPKd8CJWR6gdd7i_DshBczQ1UaDWYbmXbwL9R5lKOhpqKYOUaVtiTCI-7IbEZD-IMDHgJUJ5zju3zpV7CP_nZnQ_wIxARp3e9-wwrQpDKDB3Bav9K1r2AsuQ7gmLPQ=w1739-h978-s-no-gm?authuser=0)

### The Tech Stack

This tech stack is:

- Go: Basic handling, no fancy frameworks
- Markdown: All markdown is handled by the gomarkdown package, with some additional handlers to add more functionality.
- SCSS: Basic styling, but I hate having to manage large CSS files.
- TOML: Allows for user configuration

### Features

- Speed! Really super fast.
- Simple styling! No need for crazy amounts of JavaScript and CSS libraries, just simple SCSS.
- Markdown editing! Write new pages quickly and effectively. Can link to other Markdown pages with ease using file-based routing.
- Code highlighting! Spiffy code highlighting using Chroma.
- Code expansion! Less horizontal scrolling for visitors if you write long code lines.
- Responsive! It's a single column! Wow!

### Code highlighting

Checkout this piece of Go code that renders this piece of Go code! So humorous and meta! Wow!
Too wide? Don't worry, try clicking that "Expand Code" button!

```go
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

```

### Custom Configuration

Hate the way this site looks? No problem! You can change that with the `blog-config.toml` file!

```toml
[metadata]
title = "Blog"
description = "Blog built in Go!"
author = "Isaac Morrow"
keywords = ["blog", "go", "golang", "static", "site", "generator"]


[cssvariables] # Optional
lightthemebackground = "#faf3e0"
lightthemetext = "#444"
lightthemeaccent = "#333"
lightthemelink = "#333"

darkthemebackground = "#1b2432"
darkthemetext = "#f5f5f5"
darkthemeaccent = "#9ecfff"
darkthemelink = "#9ec"

maincontentwidth = "800px"
codeblockwidth = "1000px"

[BlogData]
contentdirectory = "./web/docs"
sourcefile = "index.md"
staticdirectory = "./web/static/"
```

### File (Based) Routing

To navigate through the project, file-based routing is used. For example, the entrance file is `index.md`, but that file can link to `examples/supercool.md`. In the browser, you can then access it through [examples/supercool](/examples/supercool)

here's what that looks like:

```
docs
 │ examples
 │ └ supercool.md
 └ index.md

```
