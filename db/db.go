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
var subscriptionsMux sync.Mutex

// Connect to database
func Connect(url string) error {
	dbConnection, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	db = dbConnection

	subscriptions = make(map[string][]subscription)
	listener = pq.NewListener(url,
		10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
			if err != nil {
				panic(err)
			}
		})

	go func() {
		for n := range listener.NotificationChannel() {
			if channels, ok := subscriptions[n.Channel]; ok {
				for _, channel := range channels {
					channel.c <- n.Extra
				}
			}
		}
	}()

	return nil
}

func subscribe(name string) subscription {
	subscriptionsMux.Lock()
	defer subscriptionsMux.Unlock()

	if subscriptions[name] == nil {
		subscriptions[name] = []subscription{}
		if err := listener.Listen(name); err != nil {
			panic(err)
		}
	}

	c := subscription{
		name: name,
		c:    make(chan string, 256),
	}

	subscriptions[name] = append(subscriptions[name], c)
	return c
}

func (c *subscription) close() {
	subscriptionsMux.Lock()
	defer subscriptionsMux.Unlock()

	j := 0
	for _, subscriptionChannel := range subscriptions[c.name] {
		if subscriptionChannel.c != c.c {
			subscriptions[c.name][j] = subscriptionChannel
			j++
		}
	}
	subscriptions[c.name] = subscriptions[c.name][:j]
	close(c.c)

	if len(subscriptions[c.name]) == 0 {
		if err := listener.Unlisten(c.name); err != nil {
			panic(err)
		}

		subscriptions[c.name] = nil
	}
}
