package api

import (
	"github.com/gorilla/mux"
	"github.com/juseongkr/websocket-chat-go/db"
	"github.com/juseongkr/websocket-chat-go/token"
	"github.com/juseongkr/websocket-chat-go/ws"
	"net/http"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Handler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/signup", signup).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/signin", signin).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/room", createRoom).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/rooms", getRooms).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/room/{id}", connectToRoom).Methods(http.MethodGet, http.MethodOptions)

	router.Use(handlePanic)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})
	router.Use(mux.CORSMethodMiddleware(router))

	return router
}

func signup(w http.ResponseWriter, r *http.Request) {
	var req user

	parseJSON(r.Body, &req)
	id, err := db.CreateUser(req.Username, req.Password)
	must(err)

	t, err := token.CreateNewToken(id)
	must(err)

	writeJSON(w, struct {
		Token string `json:"token"`
	}{t})
}

func signin(w http.ResponseWriter, r *http.Request) {
	var req user

	parseJSON(r.Body, &req)
	id, err := db.FindUser(req.Username, req.Password)
	must(err)

	t, err := token.CreateNewToken(id)
	must(err)

	writeJSON(w, struct {
		Token string `json:"token"`
	}{t})
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	parseJSON(r.Body, &req)

	id, err := db.CreateRoom(req.Name)
	must(err)

	writeJSON(w, struct {
		Id int `json:"id"`
	}{id})
}

func getRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := db.GetRooms()
	must(err)

	writeJSON(w, rooms)
}

func connectToRoom(w http.ResponseWriter, r *http.Request) {
	uid := userId(r)
	roomId := parseIntParam(r, "id")
	exists, err := db.IsRoomExists(roomId)
	must(err)

	if !exists {
		panic(notFoundError)
	}

	ws.ChatHandler(roomId, uid).ServeHTTP(w, r)
}
