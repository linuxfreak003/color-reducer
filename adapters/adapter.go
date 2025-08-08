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

func (*Reducer) SampleColors(img image.Image, n int) []color.Color {
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
	for i, kv := range cls {
		if i > n {
			return colors
		}
		colors = append(colors, kv.K)
	}
	return colors
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
