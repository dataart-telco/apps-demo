package main

import (
	"gopkg.in/redis.v3"
	"tad-demo/common"
	"tad-demo/common/opencell-api"
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
				callStatus := common.NewCallStatus(v.Payload)
				if callStatus.CallStatus == common.CallStatusCompleted {
					callInfo, err := restcommApi.GetCallInfo(callStatus.CallSid)
					if err != nil {
						common.Error.Println("opencell-billing: Failed to query restcomm for call sid:%s", callStatus.CallSid)
					} else {
						time := callInfo.Duration
						timestamp = 

					}
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

