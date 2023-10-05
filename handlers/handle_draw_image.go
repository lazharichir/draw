package handlers

import (
	"fmt"
	"image/png"
	"net/http"
)

func (h *handlers) DrawImage(w http.ResponseWriter, r *http.Request) {
	// get a tile (e.g., http://localhost:1001/image?cid=0&x=-1000&y=-1000&src=https://freshman.tech/images/dp-illustration.png)

	canvasID := chiURLQueryInt64(r, "cid")
	x := chiURLQueryInt64(r, "x")
	y := chiURLQueryInt64(r, "y")
	src := r.URL.Query().Get("src")
	fmt.Println("GET /image", canvasID, x, y, src)

	img, err := loadImageFromURL(src)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("image size", img.Bounds().Max.X, img.Bounds().Max.Y)

	// get the pixels from the image
	tile := buildTileFromImage(int64(x), int64(y), img)
	if err := h.storage.DrawPixels(canvasID, tile.Pixels); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write the image to the response
	w.Header().Set("Content-Type", "image/png")
	encoder := png.Encoder{}
	encoder.CompressionLevel = png.NoCompression
	if err := encoder.Encode(w, img); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
