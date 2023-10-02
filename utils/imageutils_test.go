package utils

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResizeImage(t *testing.T) {
	// Create a test image.
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			src.SetRGBA(x, y, color.RGBA{uint8(x), uint8(y), 0, 255})
		}
	}

	// Test resizing the image.
	dst, err := ResizeImage(src, 50, 50)
	assert.NoError(t, err)
	assert.Equal(t, 50, dst.Bounds().Dx())
	assert.Equal(t, 50, dst.Bounds().Dy())

	// Test resizing the image to the same size.
	dst, err = ResizeImage(src, 100, 100)
	assert.NoError(t, err)
	assert.Equal(t, 100, dst.Bounds().Dx())
	assert.Equal(t, 100, dst.Bounds().Dy())

	// Test resizing the image to a larger size.
	dst, err = ResizeImage(src, 200, 200)
	assert.NoError(t, err)
	assert.Equal(t, 200, dst.Bounds().Dx())
	assert.Equal(t, 200, dst.Bounds().Dy())
}
