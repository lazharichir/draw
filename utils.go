package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"

	"github.com/lazharichir/draw/core"
)

// get one tile
var tileCache = sync.Map{}

func getTile(x, y int64, d int64) core.Tile {
	cacheKey := fmt.Sprintf("%dx%d_%d", x, y, d)

	// check cache
	if tile, ok := tileCache.Load(cacheKey); ok {
		return tile.(core.Tile)
	}

	// build the tile
	tile := core.NewTile(core.Point{X: x, Y: y}, d, d)
	for i := int64(0); i < int64(d); i++ {
		for j := int64(0); j < int64(d); j++ {
			tile.AddPixels(generatePixel(x+i, y+j))
		}
	}

	// cache the tile
	tileCache.Store(cacheKey, tile)

	return tile
}

// generate a random color
var randomColors = []color.RGBA{
	{R: 100, G: 100, B: 100, A: 255},
	{R: 100, G: 100, B: 100, A: 255},
	{R: 0, G: 0, B: 0, A: 255},
	{R: 0, G: 0, B: 0, A: 255},
	{R: 0, G: 0, B: 0, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
}

func generateRandomColor() color.RGBA {
	return randomColors[rand.Intn(len(randomColors))]
}

func generatePixel(x, y int64) core.Pixel {
	return core.Pixel{X: x, Y: y, RGBA: generateRandomColor()}
}
