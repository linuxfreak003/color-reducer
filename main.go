package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/linuxfreak003/color-reducer/adapters"
)

var (
	nColors         int
	inputFilename   string
	reducedFilename string
	outlineFilename string
	outline         bool
)

func init() {
	flag.IntVar(&nColors, "colors", 24, "Number of colors to reduce to")
	flag.StringVar(&inputFilename, "input", "", "Input file (.jpg)")
	flag.StringVar(&reducedFilename, "reduced-output", "reduced.png", "Filname of reduced image (png)")
	flag.StringVar(&outlineFilename, "outline-output", "outline.png", "Filename of outline image (png)")
	flag.BoolVar(&outline, "outline", false, "Specify whether to include outline image")
}

func main() {
	flag.Parse()

	// Load an image (replace with your image loading logic)
	file, err := os.Open(inputFilename)
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
	fmt.Println("Sampling Colors...")
	colors := reducer.SampleColors(img, nColors)

	// Save the modified image
	fmt.Println("Generating reduced image...")
	newImg := reducer.ReduceImage(img, colors)
	reducedFile, err := os.Create(reducedFilename)
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
	if outline {
		fmt.Println("Generating Outline...")
		outlineImg := reducer.OutlineImage(newImg)
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
	}

	fmt.Println("Done!")
}
