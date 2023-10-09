package services_test

import (
	"context"
	"testing"

	"github.com/lazharichir/draw/config"
	"github.com/lazharichir/draw/core"
	"github.com/lazharichir/draw/services"
	"github.com/lazharichir/draw/utils"
	"github.com/stretchr/testify/assert"
)

func TestTileCache_Flow(t *testing.T) {
	cfg := config.GetConfig()

	// Create a new mock storage.
	s3 := utils.MustNewS3Client(cfg.R2_ACCOUNT_ID, cfg.R2_ACCESS_KEY_ID, cfg.R2_ACCESS_KEY_SECRET)
	cache := services.NewTileCache(s3, cfg.R2_TILECACHE_BUCKET_NAME_TEST)

	// Create a new mock image.
	mockTile := core.NewTile(core.NewArea(core.Pt(-5, -5), core.Pt(5, 5)))
	mockImg := mockTile.AsImage()

	// Call the PutTile method and check the error.
	err := cache.PutTile(context.Background(), 1, mockTile, mockImg)
	assert.NoError(t, err)

	// Call the GetTile method and check the image and error.
	img, err := cache.GetTile(context.Background(), 1, mockTile.Area)
	assert.NoError(t, err)
	assert.Equal(t, mockImg.Bounds(), img.Bounds())

	// test that all pixels are equal
	for x := mockTile.GetMinX(); x < mockTile.GetMaxX(); x++ {
		for y := mockTile.GetMinY(); y < mockTile.GetMaxY(); y++ {
			left := mockImg.At(int(x), int(y))
			right := img.At(int(x), int(y))
			assert.True(t, utils.CompareColors(left, right))
		}
	}

	// delete the tile
	err = cache.DeleteTile(context.Background(), 1, mockTile)
	assert.NoError(t, err)
}
