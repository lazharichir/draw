package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/lazharichir/draw/core"
)

func (h *handlers) PollAreaPixels(w http.ResponseWriter, r *http.Request) {
	// fail a quarter of the time
	if rand.Intn(4) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	fmt.Println("GET /poll", len(pixels), "#pixels", from, topLeft, bottomRight)

	// send empty json object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	by, _ := json.Marshal(pixels)
	w.Write(by)
}
