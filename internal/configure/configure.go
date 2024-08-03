package configure

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	CSSVariables CSSVariables
	MetaData     MetaData
}

type MetaData struct {
	Title       string
	Description string
	Author      string
	Keywords    []string
}

type CSSVariables struct {
	// Light theme colors
	LightThemeBackground string
	LightThemeText       string
	LightThemeAccent     string
	LightThemeLink       string

	// Dark theme colors
	DarkThemeBackground string
	DarkThemeText       string
	DarkThemeAccent     string
	DarkThemeLink       string

	// Widths for the main content and expanded codeblocks
	MainContentWidth string
	CodeBlockWidth   string
}

func GetConfigFromTOML() {
	var config Config
	if _, err := toml.DecodeFile("blog-config.toml", &config); err != nil {
		fmt.Println("Error decoding TOML file:", err)
		return
	}

	fmt.Println("Config MetaData:", config.MetaData)
	fmt.Println("Config CSSVariables:", config.CSSVariables)
}
