package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/linuxfreak003/util/maps"
	"github.com/linuxfreak003/util/slice"
)

var colors = []color.Color{}

func sampleColors(img image.Image, n int) {
	bounds := img.Bounds()
	colorMap := make(map[color.Color]int)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := img.At(x, y)
			colorMap[color] = colorMap[color] + 1
		}
	}

	type KV struct {
		K color.Color
		V int
	}
	cls := maps.ToSlice(colorMap, func(k color.Color, v int) KV {
		return KV{K: k, V: v}
	})
	slice.Sort(cls, func(a, b KV) bool {
		return a.V < b.V
	})
	for i, kv := range cls {
		if i > n {
			return
		}
		colors = append(colors, kv.K)
	}
}

func main() {
	// Load an image (replace with your image loading logic)
	file, err := os.Open("example.jpg")
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

	sampleColors(img, 16)
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(x, y, getClosestColor(img.At(x, y)))
		}
	}

	// Save the modified image
	outFile, err := os.Create("outout.png")
	if err != nil {
		fmt.Println("Error: Could not create output file", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, newImg)
	if err != nil {
		fmt.Println("Error: Could not encode image", err)
		return
	}

	fmt.Println("Done!")
}

func getClosestColor(in color.Color) color.Color {
	minDist := -1
	var closest color.Color
	r, g, b, _ := in.RGBA()

	// Normalize to [0, 255]
	r >>= 8
	g >>= 8
	b >>= 8

	for _, c := range colors {
		cr, cg, cb, _ := c.RGBA()
		// Normalize to [0, 255]
		cr >>= 8
		cg >>= 8
		cb >>= 8
		dr := int(r) - int(cr)
		dg := int(g) - int(cg)
		db := int(b) - int(cb)
		dist := dr*dr + dg*dg + db*db
		if minDist == -1 || dist < minDist {
			minDist = dist
			closest = c
		}
	}
	return closest
}
