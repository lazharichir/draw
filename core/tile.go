package core

import (
	"image"
	"image/color"
)

func NewTile(area Area) Tile {
	return Tile{Area: area, Pixels: []Pixel{}}
}

func NewTilePWH(pt Point, w, h int64) Tile {
	area := NewAreaWH(pt, w, h)
	return Tile{Area: area, Pixels: []Pixel{}}
}

type Tile struct {
	Area
	Pixels []Pixel
}

func (t Tile) GetMinX() int64 {
	return t.Min.X
}

func (t Tile) GetMaxX() int64 {
	return t.Max.X
}

func (t Tile) GetMinY() int64 {
	return t.Min.Y
}

func (t Tile) GetMaxY() int64 {
	return t.Max.Y
}

func (t *Tile) NewPixel(x, y int64, c color.Color) {
	t.Pixels = append(t.Pixels, NewPixel(x, y, c))
}

func (t *Tile) AddPixels(pixels ...Pixel) {
	t.Pixels = append(t.Pixels, pixels...)
}

func (t Tile) AsImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(t.Width()), int(t.Height())))

	for _, pixel := range t.Pixels {
		localX := Abs(t.GetMinX() - pixel.X)
		localY := Abs(t.GetMinY() - pixel.Y)
		img.Set(int(localX), int(localY), pixel.RGBA)
	}

	return img
}
