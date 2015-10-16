package main

import (
	"gopkg.in/redis.v3"
	"tad-demo/common"
)

type Conference struct {
	subscription Subscription
}

func (conf Conference) Subscribe() Subscription {
	common.Trace.Println("Subscribe for new calls")

	conf.subscription = Subscription{acceptedQueue: make(chan string, 100)}

	//subscribe for events from main service
	go func() {
		sub, _ := db.Subscribe(cfg.Redis.ConfChannel)
		for {
			msg, e2 := sub.Receive()
			if e2 != nil {
				common.Error.Println("receive message error", e2)
				continue
			}
			common.Trace.Println("Message ", msg)
			switch v := msg.(type) {
			case *redis.Message:
				conf.subscription.acceptedQueue <- v.Payload
			}
		}
	}()

	return conf.subscription
}

type Subscription struct {
	acceptedQueue chan string
}

func (s Subscription) Receive() string {
	return <-s.acceptedQueue
}
