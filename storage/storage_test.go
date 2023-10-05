package storage_test

import (
	"context"
	"testing"

	"github.com/lazharichir/draw/core"
	storage "github.com/lazharichir/draw/storage"
	"github.com/stretchr/testify/assert"
)

var store storage.PixelStore

func init() {
	store = storage.NewPGPixelStore(storage.NewPG(), nil)
}

func TestDeleteLastChangedForAreas(t *testing.T) {
	var err error
	ctx := context.Background()

	// Insert some test data.
	canvasID := int64(0)
	side := int64(1024)
	pts := []core.Point{
		core.Pt(100, 100),
		core.Pt(12222, 7),
		core.Pt(122221, 2047),
		core.Pt(200, 200),
		core.Pt(-300, -300),
		core.Pt(600, 600),
	}
	err = store.SetLastChangedForPoints(ctx, canvasID, side, pts...)
	assert.NoError(t, err)

	// Delete the tile changes in the specified areas.
	areas := core.GetTileAreasFromPoints(side, pts...)
	err = store.DeleteLastChangedForAreas(ctx, canvasID, side, areas...)
	assert.NoError(t, err)
}
