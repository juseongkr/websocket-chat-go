package db

import (
	"encoding/json"
	"fmt"
	"github.com/juseongkr/websocket-chat-go/chat"
	"log"
)

type ChatroomSubscription struct {
	sub subscription
	C   <-chan chat.Message
}

func sendExistsMessages(chatroomId int, c chan<- chat.Message, limit int) error {
	rows, err := db.Query(
		`WITH msgs AS (
			SELECT msg.sender_id, usr.username, msg.text, msg.sent_on
			FROM messages msg
			JOIN users usr ON msg.sender_id = usr.id
			WHERE msg.chatroom_id = $1
			ORDER BY msg.sent_on DESC
			LIMIT $2
		) SELECT * FROM msgs ORDER BY sent_on ASC`,
		chatroomId, limit,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var msg chat.Message
		if err := rows.Scan(&msg.SenderId, &msg.Sender, &msg.Text, &msg.SentOn); err != nil {
			return err
		}
		c <- msg
	}

	return nil
}

func SendMessage(senderId, chatroomId int, text string) error {
	_, err := db.Exec(
		`INSERT INTO messages (sender_id, chatroom_id, text) VALUES ($1, $2, $3)`,
		senderId, chatroomId, text,
	)

	return err
}

func NewChatroomSubscription(chatroomId int) (ChatroomSubscription, error) {
	c := make(chan chat.Message, 128)
	if err := sendExistsMessages(chatroomId, c, 100); err != nil {
		return ChatroomSubscription{}, err
	}

	chatroomSubscription := ChatroomSubscription{
		sub: subscribe(fmt.Sprintf("new_message_%d", chatroomId)),
		C:   c,
	}

	go func() {
		defer close(c)
		for msg := range chatroomSubscription.sub.c {
			var parsedMessage chat.Message
			if err := json.Unmarshal([]byte(msg), &parsedMessage); err != nil {
				log.Println("fail to parse:", err)
				continue
			}
			c <- parsedMessage
		}
	}()

	return chatroomSubscription, nil
}

func (c ChatroomSubscription) Close() {
	c.sub.close()
}
