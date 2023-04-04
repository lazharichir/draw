package core

import (
	"image"
	"image/color"
)

type Point struct {
	X int64
	Y int64
}

type Pixel struct {
	X    int64
	Y    int64
	RGBA color.RGBA
}

func NewTile(topLeft Point, width, height int64) Tile {
	return Tile{TopLeft: topLeft, Width: width, Height: height, Pixels: []Pixel{}}
}

type Tile struct {
	TopLeft Point
	Width   int64
	Height  int64
	Pixels  []Pixel
}

func (t Tile) GetMinX() int64 {
	return t.TopLeft.X
}

func (t Tile) GetMaxX() int64 {
	return t.TopLeft.X + t.Width
}

func (t Tile) GetMinY() int64 {
	return t.TopLeft.Y
}

func (t Tile) GetMaxY() int64 {
	return t.TopLeft.Y + t.Height
}

func (t *Tile) AddPixels(pixels ...Pixel) {
	t.Pixels = append(t.Pixels, pixels...)
}

func (t Tile) AsImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(t.Width), int(t.Height)))

	for _, pixel := range t.Pixels {
		localX := Abs(t.GetMinX() - pixel.X)
		localY := Abs(t.GetMinY() - pixel.Y)
		img.Set(int(localX), int(localY), pixel.RGBA)
	}

	return img
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
