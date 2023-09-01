package handlers

import (
	"encoding/json"
	"fmt"
	"image/color"
	"image/png"
	"net/http"
	"time"

	"github.com/lazharichir/draw/core"
	"github.com/lazharichir/draw/storage"
)

func New(storage storage.PixelStore) *handlers {
	return &handlers{
		storage: storage,
	}
}

type handlers struct {
	storage storage.PixelStore
}

func (h *handlers) PollAreaPixels(w http.ResponseWriter, r *http.Request) {
	canvasID := chiURLQueryInt64(r, "cid")
	from, err := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("invalid from date"))
		return
	}

	tlX := chiURLQueryInt64(r, "tlx")
	tlY := chiURLQueryInt64(r, "tly")
	brX := chiURLQueryInt64(r, "brx")
	brY := chiURLQueryInt64(r, "bry")
	topLeft := core.Point{X: tlX, Y: tlY}
	bottomRight := core.Point{X: brX, Y: brY}

	pixels, err := h.storage.GetLatestPixelsForArea(canvasID, topLeft, bottomRight, from)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// fmt.Println("GET /poll", canvasID, from, topLeft, bottomRight)

	// send empty json object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	by, _ := json.Marshal(pixels)
	w.Write(by)
}

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

func (h *handlers) ErasePixel(w http.ResponseWriter, r *http.Request) {
	canvasID := chiURLParamInt64(r, "canvasID")
	x := chiURLParamInt64(r, "x")
	y := chiURLParamInt64(r, "y")

	if err := h.storage.ErasePixel(canvasID, x, y); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handlers) DrawPixel(w http.ResponseWriter, r *http.Request) {
	canvasID := chiURLParamInt64(r, "canvasID")
	x := chiURLParamInt64(r, "x")
	y := chiURLParamInt64(r, "y")
	red := chiURLParamInt64(r, "r")
	green := chiURLParamInt64(r, "g")
	blue := chiURLParamInt64(r, "b")
	alpha := chiURLParamInt64(r, "a")

	color := color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: uint8(alpha),
	}

	if err := h.storage.DrawPixelRGBA(canvasID, x, y, color); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handlers) GetTileImage(w http.ResponseWriter, r *http.Request) {
	canvasID := int64(0)
	x := chiURLParamInt64(r, "x")
	y := chiURLParamInt64(r, "y")
	d := chiURLParamInt64(r, "d")

	pixels, err := h.storage.GetPixelsFromTopLeft(canvasID, x, y, d)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newTile := core.NewTile(core.Point{X: x, Y: y}, d, d)
	newTile.AddPixels(pixels...)

	img := newTile.AsImage()

	// render the tile as a png image
	w.Header().Set("Content-Type", "image/png")
	encoder := png.Encoder{}
	encoder.CompressionLevel = png.NoCompression
	if err := encoder.Encode(w, img); err != nil {
		fmt.Println(err)
	}
}
