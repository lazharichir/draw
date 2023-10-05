package handlers

import (
	"fmt"
	"image/color"
	"net/http"

	"github.com/lazharichir/draw/core"
	"github.com/lazharichir/draw/services"
)

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

	pixel := core.NewPixel(x, y, color)

	// check if the pixel can be drawn
	if ok, err := h.landRegistry.CanDrawPixel(r.Context(), 0, 0, pixel); err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !ok {
		err := services.ErrCannotDrawInArea(0, pixel.Point, pixel.Point)
		fmt.Println(err)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.storage.DrawPixels(canvasID, []core.Pixel{pixel}); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
