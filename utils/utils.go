package utils

import "image/color"

func Ptr[T any](input T) *T {
	return &input
}

// ensure two rgba colors are equal
func CompareColors(c1, c2 color.Color) bool {
	// c1 rgba
	c1R, c1G, c1B, c1A := c1.RGBA()
	// c2 rgba
	c2R, c2G, c2B, c2A := c2.RGBA()
	// compare
	return c1R == c2R && c1G == c2G && c1B == c2B && c1A == c2A
}
