package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/linuxfreak003/color-reducer/adapters"
)

func main() {
	// Load an image (replace with your image loading logic)
	file, err := os.Open("sunset.jpeg")
	if err != nil {
		fmt.Println("Error opening image:", err)
		return
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	reducer := adapters.NewReducer()
	colors := reducer.SampleColors(img, 24)
	newImg := reducer.ReduceImage(img, colors)
	outlineImg := reducer.OutlineImage(newImg)

	// Save the modified image
	reducedFile, err := os.Create("reduced.png")
	if err != nil {
		fmt.Println("Error: Could not create output file", err)
		return
	}
	defer reducedFile.Close()

	err = png.Encode(reducedFile, newImg)
	if err != nil {
		fmt.Println("Error: Could not encode image", err)
		return
	}

	// Save the outline image
	outlineFile, err := os.Create("outline.png")
	if err != nil {
		fmt.Println("Error: Could not create outline file", err)
		return
	}
	defer outlineFile.Close()

	err = png.Encode(outlineFile, outlineImg)
	if err != nil {
		fmt.Println("Error: Could not encode outline image", err)
		return
	}

	fmt.Println("Done!")
}
