package services

import (
	"context"
	"errors"
	"image"

	"github.com/lazharichir/draw/core"
	"github.com/lazharichir/draw/storage"
)

// Create a function that generates updated cached tiles

type TileCache struct {
	store storage.PixelStore
}

func NewTileCache(store storage.PixelStore) *TileCache {
	return &TileCache{store}
}

func (cache *TileCache) PutTile(ctx context.Context, canvasID int64, tile core.Tile, img image.Image) error {
	x := tile.GetMinX()
	y := tile.GetMinY()
	side := tile.Width
	_, _, _ = x, y, side

	if tile.Width != tile.Height {
		return errors.New("tile is not a square")
	}

	return nil
}

func (cache *TileCache) GetTile(ctx context.Context, canvasID int64, area core.Area) (image.Image, error) {
	return nil, nil
}
