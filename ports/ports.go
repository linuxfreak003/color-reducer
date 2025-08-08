package ports

import (
	"image"
	"image/color"
)

type ColorReducer interface {
	SampleColors(image.Image, int) []color.Color
	ReduceImage(image.Image, []color.Color) image.Image
	OutlineImage(image.Image) image.Image
}
