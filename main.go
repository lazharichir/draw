package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lazharichir/draw/core"
	"github.com/lazharichir/draw/storage"
)

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

// Gzip Compression
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		handler.ServeHTTP(gzw, r)
	})
}

func main() {

	fmt.Println("Server started:", "http://localhost:1001")

	// storage
	db := storage.NewPG()
	pixelStore := storage.NewPGPixelStore(db, nil)

	// create a new router
	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "gzip"), middleware.Logger, middleware.StripSlashes)
	// cors allow all
	r.Use(cors.Handler(
		cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		},
	))

	// get a tile (e.g., http://localhost:1001/tile/0x0_10.jpg)
	r.Get("/tile/{x}x{y}_{d}.png", Gzip(func(w http.ResponseWriter, r *http.Request) {
		canvasID := int64(0)
		x := chiURLParamInt64(r, "x")
		y := chiURLParamInt64(r, "y")
		d := chiURLParamInt64(r, "d")

		pixels, err := pixelStore.GetPixelsFromTopLeft(canvasID, x, y, d)
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
	}))

	// draw a pixel
	r.Put("/pixel/{canvasID}/{x}/{y}/{r}/{g}/{b}/{a}", func(w http.ResponseWriter, r *http.Request) {

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

		fmt.Println("draw pixel", canvasID, x, y, color)

		if err := pixelStore.DrawPixelRGBA(canvasID, x, y, color); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	// erase a pixel
	r.Delete("/pixel/{canvasID}/{x}/{y}", func(w http.ResponseWriter, r *http.Request) {

		canvasID := chiURLParamInt64(r, "canvasID")
		x := chiURLParamInt64(r, "x")
		y := chiURLParamInt64(r, "y")

		fmt.Println("erase pixel", canvasID, x, y)

		if err := pixelStore.ErasePixel(canvasID, x, y); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	// get a tile (e.g., http://localhost:1001/image?cid=0&x=-1000&y=-1000&src=https://freshman.tech/images/dp-illustration.png)
	r.Get("/image", func(w http.ResponseWriter, r *http.Request) {
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
		if err := pixelStore.DrawPixels(canvasID, tile.Pixels); err != nil {
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
	})

	// start the server
	http.ListenAndServe(":1001", r)
}

func buildTileFromImage(x, y int64, img image.Image) core.Tile {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	tile := core.NewTile(core.Point{X: x, Y: y}, int64(width), int64(height))

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
