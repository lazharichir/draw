package core

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPixel(t *testing.T) {
	x := int64(10)
	y := int64(20)
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	expected := Pixel{Point: Point{X: x, Y: y}, RGBA: c}
	actual := NewPixel(x, y, c)
	assert.Equal(t, expected, actual)
}

func TestPixel_SetPoint(t *testing.T) {
	pixel := Pixel{Point: Point{X: 10, Y: 20}, RGBA: color.RGBA{R: 255, G: 0, B: 0, A: 255}}
	expected := Point{X: 30, Y: 40}
	pixel.SetPoint(expected)
	assert.Equal(t, expected, pixel.Point)
}

func TestPixel_SetColor(t *testing.T) {
	pixel := Pixel{Point: Point{X: 10, Y: 20}, RGBA: color.RGBA{R: 255, G: 0, B: 0, A: 255}}
	expected := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	pixel.SetColor(expected)
	assert.Equal(t, expected, pixel.RGBA)
}
