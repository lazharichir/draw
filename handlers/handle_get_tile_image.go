package handlers

import (
	"fmt"
	"image"
	"image/png"
	"net/http"

	"github.com/lazharichir/draw/core"
)

func (h *handlers) respondWithImage(w http.ResponseWriter, r *http.Request, img image.Image) {
	// render the tile as a png image
	w.Header().Set("Content-Type", "image/png")
	encoder := png.Encoder{}
	encoder.CompressionLevel = png.NoCompression
	if err := encoder.Encode(w, img); err != nil {
		fmt.Println("respondWithImage", err)
	}
}

func (h *handlers) GetTileImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	canvasID := int64(0)
	x := chiURLParamInt64(r, "x")
	y := chiURLParamInt64(r, "y")
	d := chiURLParamInt64(r, "d")
	area := core.NewAreaSquare(core.Pt(x, y), d)

	// check if the tile is in the cache
	cached, err := h.tileCache.GetTile(ctx, canvasID, area)
	if err != nil {
		fmt.Println("GetTileImage.tileCache.GetTile", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if a tile was cached, return it
	if cached != nil {
		h.respondWithImage(w, r, cached)
		return
	}

	// if not, get the pixels from the storage
	pixels, err := h.storage.GetPixelsFromTopLeft(canvasID, x, y, d)
	if err != nil {
		fmt.Println("GetTileImage.storage.GetPixelsFromTopLeft", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create a new tile and add the pixels to it
	newTile := core.NewTile(core.Point{X: x, Y: y}, d, d)
	newTile.AddPixels(pixels...)

	// render the tile as a png image
	cached = newTile.AsImage()

	// store the tile in the cache
	if err := h.tileCache.PutTile(ctx, canvasID, newTile, cached); err != nil {
		fmt.Println("GetTileImage.tileCache.PutTile", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// render the tile as a png image
	h.respondWithImage(w, r, cached)
}
