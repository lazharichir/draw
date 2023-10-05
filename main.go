package main

import (
	"compress/gzip"
	"fmt"
	"image/color"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lazharichir/draw/core"
	"github.com/lazharichir/draw/handlers"
	"github.com/lazharichir/draw/services"
	"github.com/lazharichir/draw/storage"
)

var (
	lastTopLeft     = core.Point{X: -1000, Y: 1000}
	lastBottomRight = core.Point{X: -1000, Y: 1000}
)

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

	db := storage.NewPG()
	storage := storage.NewPGPixelStore(db, nil)
	landRegistry := services.NewLandRegistry(db)
	tileCache := services.NewTileCache(storage)

	handlers := handlers.New(storage, landRegistry, tileCache)

	r := chi.NewRouter()

	r.Use(
		middleware.Compress(5, "gzip"),
		// middleware.Logger,
		middleware.StripSlashes,
	)
	r.Use(cors.Handler(
		cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		},
	))

	r.Get("/tile/{x}x{y}_{d}.png", Gzip(handlers.GetTileImage))
	r.Put("/pixel/{canvasID}/{x}/{y}/{r}/{g}/{b}/{a}", handlers.DrawPixel)
	r.Delete("/pixel/{canvasID}/{x}/{y}", handlers.ErasePixel)
	r.Get("/image", handlers.DrawImage)
	r.Get("/ws", handlers.TheWS)
	r.Get("/poll", func(w http.ResponseWriter, r *http.Request) {
		lastTopLeft = core.Point{X: chiURLQueryInt64(r, "tlx"), Y: chiURLQueryInt64(r, "tly")}
		lastBottomRight = core.Point{X: chiURLQueryInt64(r, "brx"), Y: chiURLQueryInt64(r, "bry")}
		handlers.PollAreaPixels(w, r)
	})

	ticker := time.NewTicker(900 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:

				randomPixels := []core.Pixel{}
				for i := 0; i < rand.Intn(2); i++ {
					px := generatePixel(
						randomBetweenInts(lastTopLeft.X, lastBottomRight.X),
						randomBetweenInts(lastTopLeft.Y, lastBottomRight.Y),
					)
					randomPixels = append(randomPixels, px)
				}

				if err := storage.DrawPixels(1, randomPixels); err != nil {
					fmt.Println(err)
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// start the server
	http.ListenAndServe(":1001", r)
}

func randomBetweenInts(min, max int64) int64 {
	if min == max {
		return min
	}
	return rand.Int63n(max-min) + min
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

// get one tile
var tileCache = sync.Map{}

func getTile(x, y int64, d int64) core.Tile {
	cacheKey := fmt.Sprintf("%dx%d_%d", x, y, d)

	// check cache
	if tile, ok := tileCache.Load(cacheKey); ok {
		return tile.(core.Tile)
	}

	// build the tile
	tile := core.NewTile(core.Point{X: x, Y: y}, d, d)
	for i := int64(0); i < int64(d); i++ {
		for j := int64(0); j < int64(d); j++ {
			tile.AddPixels(generatePixel(x+i, y+j))
		}
	}

	// cache the tile
	tileCache.Store(cacheKey, tile)

	return tile
}

// generate a random color
func randomRGBA() color.RGBA {
	rgba := color.RGBA{}
	rgba.R = uint8(rand.Intn(255))
	rgba.G = uint8(rand.Intn(255))
	rgba.B = uint8(rand.Intn(255))
	rgba.A = 255
	return rgba
}

var randomColors = []color.RGBA{
	{R: 100, G: 100, B: 100, A: 255},
	{R: 100, G: 100, B: 100, A: 255},
	{R: 0, G: 0, B: 0, A: 255},
	{R: 0, G: 0, B: 0, A: 255},
	{R: 0, G: 0, B: 0, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
	{R: 255, G: 242, B: 232, A: 255},
}

func generateRandomColor() color.RGBA {
	return randomColors[rand.Intn(len(randomColors))]
}

func generatePixel(x, y int64) core.Pixel {
	return core.NewPixel(x, y, randomRGBA())
}
