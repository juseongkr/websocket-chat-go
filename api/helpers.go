package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/juseongkr/websocket-chat-go/db"
	"github.com/juseongkr/websocket-chat-go/token"
	"io"
	"log"
	"net/http"
	"strconv"
)

func must(err error) {
	if err == db.ErrUnauthorized {
		panic(unauthorizedError)
	} else if err != nil {
		log.Println("internal error:", err)
		panic(internalError)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	must(json.NewEncoder(w).Encode(v))
}

func parseJSON(r io.Reader, v interface{}) {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		log.Println("parse json body:", err)
		panic(malformedInputError)
	}
}

func parseIntParam(r *http.Request, key string) int {
	vars := mux.Vars(r)
	if v, ok := vars[key]; ok {
		index, err := strconv.Atoi(v)
		if err == nil {
			return index
		}
	}
	panic(malformedInputError)
}

func userId(r *http.Request) int {
	t := r.URL.Query().Get("token")
	id, err := token.Parse(t)
	if err != nil {
		log.Println(err)
		panic(unauthorizedError)
	}

	return id
}
