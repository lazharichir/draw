package handlers

import (
	"fmt"
	"net/http"
)

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
