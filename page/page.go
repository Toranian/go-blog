package page

import (
	"fmt"
	"os"
)

// --- Page components ---
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	filename := p.Title + ".txt"

	// 0600 = Read/write permissions for current user
	return os.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)

	if err != nil {
		fmt.Printf("Failed to load page: %v", filename)
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}
