package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

const (
	MAXHEIGHT = 800
)

func main() {
	// Define the relative path to list files from
	relativePath := "web/static/images/"

	// Get the current working directory
	cwd, err := os.Getwd()

	// Combine the current working directory with the relative path
	dir := filepath.Join(cwd, relativePath)

	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %v", err)
	}

	// Walk through the directory and list files
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fmt.Printf("%s\t%d bytes\n", path, info.Size())
			img, err := imaging.Open(path)

			if err != nil {
				log.Printf("Failed to open image %s", err)
				return nil
			}

			// Get the image resolutions
			bounds := img.Bounds()
			width, height := bounds.Dx(), bounds.Dy()
			fmt.Printf("Image: %s\tResolution: %dx%d\n", path, width, height)

			if height > MAXHEIGHT {
				resizedImg := imaging.Resize(img, MAXHEIGHT, 0, imaging.Lanczos)
				// outputPath := filepath.Join(filepath.Dir(path), )
				err = imaging.Save(resizedImg, path)
				if err != nil {
					fmt.Printf("Resized image.")
				}
			}

		}
		return nil
	})

	if err != nil {
		fmt.Printf("Failed to walk directory: %v", err)
	}
}
