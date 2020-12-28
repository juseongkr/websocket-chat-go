package db

import (
	"database/sql"
	"github.com/lib/pq"
	"sync"
	"time"
)

var db *sql.DB
var listener *pq.Listener

type subscription struct {
	name string
	c    chan string
}

var subscriptions map[string][]subscription
var subscriptionMux sync.Mutex

func Connect(url string) error {
	c, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	db = c

	subscriptions = make(map[string][]subscription)
	listener = pq.NewListener(url, time.Second*10, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			panic(err)
		}
	})

	go func() {
		for n := range listener.NotificationChannel() {
			if channels, ok := subscriptions[n.Channel]; ok {
				for _, c := range channels {
					c.c <- n.Extra
				}
			}
		}
	}()

	return nil
}

func subscribe(name string) subscription {
	subscriptionMux.Lock()
	defer subscriptionMux.Unlock()

	if subscriptions[name] == nil {
		subscriptions[name] = []subscription{}
		if err := listener.Listen(name); err != nil {
			panic(err)
		}
	}

	s := subscription{
		name: name,
		c:    make(chan string, 256),
	}

	subscriptions[name] = append(subscriptions[name], s)

	return s
}

func (s *subscription) close() {
	subscriptionMux.Lock()
	defer subscriptionMux.Unlock()

	idx := 0
	for _, subscriptionChannel := range subscriptions[s.name] {
		if subscriptionChannel.c != s.c {
			subscriptions[s.name][idx] = subscriptionChannel
			idx++
		}
	}

	subscriptions[s.name] = subscriptions[s.name][:idx]
	close(s.c)

	if len(subscriptions[s.name]) == 0 {
		if err := listener.Unlisten(s.name); err != nil {
			panic(err)
		}

		subscriptions[s.name] = nil
	}
}
