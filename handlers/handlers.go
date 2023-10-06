package handlers

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lazharichir/draw/core"
	"github.com/lazharichir/draw/services"
	"github.com/lazharichir/draw/storage"
)

func New(
	storage storage.PixelStore,
	landRegistry *services.LandRegistry,
	tileCache *services.TileCache,
) *handlers {
	return &handlers{
		storage:      storage,
		landRegistry: landRegistry,
		tileCache:    tileCache,
	}
}

type handlers struct {
	storage      storage.PixelStore
	landRegistry *services.LandRegistry
	tileCache    *services.TileCache
}

func strToInt64(str string) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

func chiURLQueryInt64(r *http.Request, key string) int64 {
	str := r.URL.Query().Get(key)
	if len(str) == 0 {
		return -1
	}
	return strToInt64(str)
}

func chiURLParamInt64(r *http.Request, key string) int64 {
	str := chi.URLParam(r, key)
	if len(str) == 0 {
		return -1
	}
	return strToInt64(str)
}

func buildTileFromImage(x, y int64, img image.Image) core.Tile {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	tile := core.NewTilePWH(core.Point{X: x, Y: y}, int64(width), int64(height))

	for i := int64(0); i < int64(width); i++ {
		for j := int64(0); j < int64(height); j++ {
			imagePx := img.At(int(i), int(j))
			tile.NewPixel(x+i, y+j, imagePx)
		}
	}

	return tile
}

func loadImageFromURL(URL string) (image.Image, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("received non 200 response code")
	}

	switch response.Header.Get("Content-Type") {
	case "image/png":
		return png.Decode(response.Body)
	case "image/jpeg":
		return jpeg.Decode(response.Body)
	case "image/jpg":
		return jpeg.Decode(response.Body)
	default:
		return nil, errors.New("unsupported content type (only jpg and png are supported)")
	}

}
