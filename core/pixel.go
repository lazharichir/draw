package core

import (
	"image/color"
)

func NewPixel(x, y int64, c color.Color) Pixel {
	return Pixel{Point: Point{X: x, Y: y}, RGBA: color.RGBAModel.Convert(c).(color.RGBA)}
}

type Pixel struct {
	Point
	RGBA color.RGBA
}

func (pixel *Pixel) SetPoint(p Point) {
	pixel.Point = p
}

func (pixel *Pixel) SetColor(c color.Color) {
	pixel.RGBA = color.RGBAModel.Convert(c).(color.RGBA)
}
