package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lazharichir/draw/handlers"
	"github.com/lazharichir/draw/services"
	"github.com/lazharichir/draw/storage"
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
		middleware.Logger,
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
	r.Get("/poll", handlers.PollAreaPixels)
	r.Get("/precache", handlers.PrecacheChangedTiles)

	// start the server
	http.ListenAndServe(":1001", r)
}
