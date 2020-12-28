package db

import (
	"github.com/juseongkr/websocket-chat-go/chat"
)

func CreateRoom(name string) (id int, err error) {
	err = db.QueryRow(
		`INSERT INTO chatrooms (name) VALUES ($1) RETURNING id`,
		name,
	).Scan(&id)

	return
}

func IsRoomExists(id int) (exists bool, err error) {
	err = db.QueryRow(
		`SELECT EXISTS ( SELECT 1 FROM chatrooms WHERE id = $1)`,
		id,
	).Scan(&exists)

	return
}

func GetRooms() ([]chat.Room, error) {
	rooms := []chat.Room{}
	rows, err := db.Query(`SELECT id, name FROM chatrooms`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var room chat.Room
		if err := rows.Scan(&room.Id, &room.Name); err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}
