package main

import (
	"cmp"
	"fmt"
	"log"
	"net/http"
	"os"

	configure "goblog/internal/configure"
	"goblog/internal/handlers"
)

func main() {
	port := cmp.Or(os.Getenv("PORT"), "3000")
	url := "http://localhost:" + port

	// Load the configuration from the TOML file
	config, err := configure.GetConfigFromTOML()
	if err != nil {
		log.Fatalf("Error loading configuration file. Error: %s", err)
		return
	}
	fmt.Println("Configuration loaded successfully.")

	content := configure.GenerateSCSSVariables(config.CSSVariables)
	blogData := config.BlogData

	// Write the SCSS variables to a file
	err = os.WriteFile("web/static/scss/_variables.scss", []byte(content), 0644)

	if err != nil {
		log.Fatal("Error writing SCSS variables to file.")
		return
	}

	fmt.Printf("Initialization successful.\n")

	fs := http.FileServer(http.Dir(blogData.StaticDirectory))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Create an instance of BlogHandler
	blogHandler := &handlers.BlogHandler{Config: config}

	// Define a wrapper function for http.HandleFunc
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		blogHandler.Handler(w, r)
	})

	fmt.Printf("Server running at %s\n", url)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
