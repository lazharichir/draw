package core

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTile(t *testing.T) {
	topLeft := Point{X: 10, Y: 20}
	width := int64(30)
	height := int64(40)
	expected := Tile{
		Area:   Area{Min: topLeft, Max: Pt(40, 60)},
		Pixels: []Pixel{},
	}
	actual := NewTilePWH(topLeft, width, height)
	assert.Equal(t, expected, actual)
}

func TestTile_GetMinX(t *testing.T) {
	tile := NewTilePWH(Pt(10, 20), 30, 40)
	expected := int64(10)
	actual := tile.GetMinX()
	assert.Equal(t, expected, actual)
}

func TestTile_GetMaxX(t *testing.T) {
	tile := NewTilePWH(Pt(10, 20), 30, 40)
	expected := int64(40)
	actual := tile.GetMaxX()
	assert.Equal(t, expected, actual)
}

func TestTile_GetMinY(t *testing.T) {
	tile := NewTilePWH(Pt(10, 20), 30, 40)
	expected := int64(20)
	actual := tile.GetMinY()
	assert.Equal(t, expected, actual)
}

func TestTile_GetMaxY(t *testing.T) {
	tile := NewTilePWH(Pt(10, 20), 30, 40)
	expected := int64(60)
	actual := tile.GetMaxY()
	assert.Equal(t, expected, actual)
}

func TestTile_NewPixel(t *testing.T) {
	tile := Tile{}
	tile.NewPixel(10, 20, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	expected := []Pixel{
		NewPixel(10, 20, color.RGBA{R: 255, G: 0, B: 0, A: 255}),
	}
	assert.Equal(t, expected, tile.Pixels)
}

func TestTile_AddPixels(t *testing.T) {
	tile := Tile{}
	tile.AddPixels(
		NewPixel(10, 20, color.RGBA{R: 255, G: 0, B: 0, A: 255}),
		NewPixel(20, 30, color.RGBA{R: 0, G: 255, B: 0, A: 255}),
	)
	expected := []Pixel{
		NewPixel(10, 20, color.RGBA{R: 255, G: 0, B: 0, A: 255}),
		NewPixel(20, 30, color.RGBA{R: 0, G: 255, B: 0, A: 255}),
	}
	assert.Equal(t, expected, tile.Pixels)
}

func TestTile_AsImage(t *testing.T) {
	tile := NewTilePWH(Pt(10, 20), 30, 40)
	tile.AddPixels(
		NewPixel(10, 20, color.RGBA{R: 255, G: 0, B: 0, A: 255}),
		NewPixel(20, 30, color.RGBA{R: 0, G: 255, B: 0, A: 255}),
	)
	expected := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	actual := tile.AsImage().At(0, 0)
	assert.Equal(t, expected, actual)
}
