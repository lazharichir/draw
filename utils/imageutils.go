package utils

import (
	"image"

	"golang.org/x/image/draw"
)

func ResizeImage(src image.Image, newWidth, newHeight int64) (image.Image, error) {
	// Create a new image with the specified dimensions.
	dst := image.NewNRGBA(image.Rect(0, 0, int(newWidth), int(newHeight)))

	// Draw the source image onto the destination image, scaled to fit.
	scaler := draw.BiLinear.NewScaler(src.Bounds().Dx(), src.Bounds().Dy(), dst.Bounds().Dx(), dst.Bounds().Dy())

	// Scale(dst Image, dr image.Rectangle, src image.Image, sr image.Rectangle, op Op, opts *Options)
	scaler.Scale(
		dst,
		dst.Bounds(),
		src,
		src.Bounds(),
		draw.Over,
		nil,
	)

	// done
	return dst, nil
}
