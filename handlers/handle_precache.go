package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lazharichir/draw/core"
)

func (h *handlers) Precache(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	from, to := core.GetMinuteRangeForTime(time.Now().UTC())

	// load recently-changed areas
	lookup, err := h.storage.FindRecentlyChangedAreasBetweenDates(ctx, from, to)
	if err != nil {
		fmt.Println(ctx, "error loading recently-changed areas", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(lookup) == 0 {
		fmt.Println(ctx, "no recently-changed areas")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handlers) PrecacheArea(ctx context.Context, canvasID int64, area core.Area) error {
	x := area.Min.X
	y := area.Min.Y
	d := area.Width()

	// load pixels
	pixels, err := h.storage.GetPixelsFromTopLeft(canvasID, x, y, d)
	if err != nil {
		return err
	}

	// create a new tile and add the pixels to it
	newTile := core.NewTile(core.Pt(x, y), d, d)
	newTile.AddPixels(pixels...)

	// build the image
	img := newTile.AsImage()

	// store the tile in the cache
	if err := h.tileCache.PutTile(ctx, canvasID, newTile, img); err != nil {
		return err
	}

	return nil
}
