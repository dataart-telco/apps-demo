package main

import (
	"gopkg.in/redis.v3"
	"tad-demo/common"
)

var db = common.NewDbClient(cfg.Service.Redis)

type Listener struct {
}

func (l Listener) Subscribe() {
	l.subscription = Subscription{acceptedQueue: make(chan string, 100)}

	//subscribe for events from call_status queue
	go func() {
		sub, _ := db.Subscribe("call_status")
		for {
			msg, e2 := sub.Receive()
			if e2 != nil {
				common.Error.Println("opencell-billing: Error receiving message from Redis 'call_status' queue")
				panic(e2)
			}
			switch v := msg.(type) {
			case *redis.Message:
				//poke restcomm for call duration
				//poke opencell for billing using duration
			}
		}
	}()
}

