package main

import (
	"gopkg.in/redis.v3"
	"tad-demo/common"
)

var db = common.NewDbClient(cfg.Service.Redis)

type Storage struct {
	subscription Subscription
}

func (s Storage) Subscribe() Subscription {
	//clear list of calls
	db.Del(common.DB_KEY_URI)

	s.subscription = Subscription{acceptedQueue: make(chan string, 100)}

	//subscribe for events from main service
	go func() {
		sub, _ := db.Subscribe(cfg.Redis.MainChannel)
		for {
			msg, e2 := sub.Receive()
			if e2 != nil {
				panic(e2)
			}
			switch v := msg.(type) {
			case *redis.Message:
				s.subscription.acceptedQueue <- v.Payload
			}
		}
	}()

	return s.subscription
}

type Subscription struct {
	acceptedQueue chan string
}

func (sms Subscription) Receive() string {
	return <-sms.acceptedQueue
}
