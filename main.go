package main

import (
	"github.com/juseongkr/websocket-chat-go/ws"
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8080", ws.Handler()); err != nil {
		log.Fatal(err)
	}
}
