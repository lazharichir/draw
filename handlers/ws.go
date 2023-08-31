package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

func (h *handlers) TheWS(w http.ResponseWriter, r *http.Request) {
	// canvasID := chiURLParamInt64(r, "canvasID")
	fmt.Println("TheWS")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("TheWS.err1", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	fmt.Println("TheWS.beforeForLoop")
	for {
		_, body, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("TheWS.err2", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("Recieved: %v", string(body))
		if err := conn.WriteMessage(1, []byte("Server Received the message ")); err != nil {
			fmt.Println("TheWS.err3", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
