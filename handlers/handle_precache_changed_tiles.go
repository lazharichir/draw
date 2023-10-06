package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lazharichir/draw/core"
)

func (h *handlers) PrecacheChangedTiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	from, to := core.GetMinuteRangeForTime(time.Now().UTC())

	// load recently-changed areas
	lookup, err := h.storage.FindRecentlyChangedAreasBetweenDates(ctx, from, to)
	if err != nil {
		fmt.Println(ctx, "error loading recently-changed areas", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// bail early if there are no recently-changed areas
	if len(lookup) == 0 {
		fmt.Println(ctx, "no recently-changed areas")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// iterate over the areas
	for canvasID, areas := range lookup {
		fmt.Println(`- Prechaching canvasID`, canvasID, `with`, len(areas), `areas`)
		// iterate over the areas
		for _, area := range areas {
			// precache the area
			if err := h.PrecacheArea(ctx, canvasID, area); err != nil {
				fmt.Println(ctx, "error precaching area", err)
			}
		}
	}

	// respond with a 200
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) PrecacheArea(ctx context.Context, canvasID int64, area core.Area) error {
	x := area.Min.X
	y := area.Min.Y
	d := area.Width()

	fmt.Println(`--- Prechaching area`, x, y, d, area.String())

	// load pixels
	pixels, err := h.storage.GetPixelsFromTopLeft(canvasID, x, y, d)
	if err != nil {
		return err
	}

	fmt.Println(`------- with`, len(pixels), `pixels`)

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
