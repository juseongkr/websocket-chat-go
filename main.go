package main

import (
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/juseongkr/websocket-chat-go/api"
	"github.com/juseongkr/websocket-chat-go/db"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Loading .env file error:", err)
	}

	if err := db.Connect(os.Getenv("DB_URL")); err != nil {
		log.Fatalln("DB Connection error:", err)
	}

	log.Println("hello db:", os.Getenv("DB_URL"))

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), api.Handler()); err != nil {
		log.Fatalln(err)
	}
}
