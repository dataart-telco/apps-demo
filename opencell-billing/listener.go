package main

import (
	"gopkg.in/redis.v3"
	"tad-demo/common"
)
var cfg = common.NewConfig()
var db = common.NewDbClient(cfg.Service.Redis)
var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)

type BillingListener struct {
}

func (l BillingListener) Subscribe() {
	//subscribe for events from call_status queue
	go func() {
		sub, _ := db.Subscribe(common.CHANNEL_CALL_STATUS)
		for {
			msg, e2 := sub.Receive()
			if e2 != nil {
				common.Error.Println("opencell-billing: Error receiving message from Redis 'call_status' queue")
				panic(e2)
			}
			switch v := msg.(type) {
			case *redis.Message:
				callStatus := common.NewCallStatus(msg)
				callInfo, err := restcommApi.GetCallInfo(msg.CallSid)
				if err != nil {
					common.Error.Println("opencell-billing: Failed to query restcomm for call sid:%s", msg.CallSid)
				} else {
					//opencell goes here

				}
			}
		}
	}()
	go func() {
		sub, _ := db.Subscribe(common.CHANNEL_CONF_DROPPED)
		for {
			msg, e2 := sub.Receive()
			if e2 != nil {
				common.Error.Println("opencell-billing: Error receiving message from Redis 'call dropped' queue")
				panic(e2)
			}
			switch v := msg.(type) {
			case *redis.Message:
				//do somethin on call drop
			}
		}
	}()
}

