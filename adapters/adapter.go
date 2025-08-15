package adapters

import (
	"image"
	"image/color"

	"github.com/linuxfreak003/color-reducer/ports"
	"github.com/linuxfreak003/util/maps"
	"github.com/linuxfreak003/util/slice"
)

type Reducer struct{}

func NewReducer() ports.ColorReducer {
	return &Reducer{}
}

func (r *Reducer) SampleColors(img image.Image, n int) []color.Color {
	colors := []color.Color{}
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
	_ = r.SampleColors2(img, n)
	for i, kv := range cls {
		if i > n {
			return colors
		}
		colors = append(colors, kv.K)
	}
	return colors
}

func (*Reducer) SampleColors2(img image.Image, n int) []color.Color {
	bounds := img.Bounds()
	rgbCube := make([][][]int, 256)
	for i := range rgbCube {
		rgbCube[i] = make([][]int, 256)
		for j := range rgbCube[i] {
			rgbCube[i][j] = make([]int, 256)
		}
	}
	// Get all colors from image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := img.At(x, y)
			rC, gC, bC, _ := color.RGBA()
			r := rC >> 8
			g := gC >> 8
			b := bC >> 8
			rgbCube[r][g][b]++
		}
	}
	// Get box bounds
	var minR, maxR, minG, maxG, minB, maxB int
	var totalColors int
	for r, plane := range rgbCube {
		for g, line := range plane {
			for b, n := range line {
				if n > 0 {
					totalColors++
					minR = min(r, minR)
					maxR = max(r, maxR)
					minG = min(g, minG)
					maxG = max(g, maxG)
					minB = min(b, minB)
					maxB = max(b, maxB)
				}
			}
		}
	}
	// Loop through the cube by creating cube boundaries
	// that get progressively smaller, adding in the colors
	// that we come across. Giving priority to more popular colors.
	for r := minR; r <= maxR; r++ {
		for g := minG; r <= maxG; g++ {
			for b := minB; b <= maxB; b++ {
				i := rgbCube[r][g][b]
				if i > 0 {
					// Do something
				}
			}
		}
	}
	// Create box that closes in around points.
	return nil
}

func (*Reducer) ReduceImage(img image.Image, colors []color.Color) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(x, y, getClosestColor(img.At(x, y), colors))
		}
	}
	return newImg
}

func (*Reducer) OutlineImage(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if !matchingAdjacentColors(img, x, y) {
				newImg.Set(x, y, color.Black)
			} else {
				newImg.Set(x, y, color.White)
			}
		}
	}
	return newImg
}

func matchingAdjacentColors(img image.Image, x, y int) bool {
	bounds := img.Bounds()
	color := img.At(x, y)
	for j := max(y-1, bounds.Min.Y); j < min(y+1, bounds.Max.Y); j++ {
		for i := max(x-1, bounds.Min.X); i < min(x+1, bounds.Max.X); i++ {
			if img.At(i, j) != color {
				return false
			}
		}
	}
	return true
}

func getClosestColor(in color.Color, colors []color.Color) color.Color {
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
