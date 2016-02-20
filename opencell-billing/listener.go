package main

import (
	"gopkg.in/redis.v3"
	"tad-demo/common"
	"tad-demo/common/opencell-api"
	"time"
	"fmt"
)

const (
	opencellUser = "DataArt Conference"
)
var cfg = common.NewConfig()
var db = common.NewDbClient(cfg.Service.Redis)
var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)
var opencellApi = opencell_api.NewOpencellAPI(cfg.Opencell.Password, cfg.Opencell.Host)


type BillingListener struct {
}

func (l BillingListener) Subscribe() {
	//subscribe for events from call_status queue
	opencellApi.InitOpenCell()
	opencellApi.CreateNewCustomer(opencellUser)
	
	sub, _ := db.Subscribe(common.CHANNEL_CALL_STATUS)
	for {
		msg, e2 := sub.Receive()
		if e2 != nil {
			common.Error.Println("opencell-billing: Error receiving message from Redis 'call_status' queue")
			panic(e2)
		}
		switch v := msg.(type) {
		case *redis.Message:
			go func(){
				l.DoMessage(v)
			}()
		}
	}
}

func (l BillingListener) DoMessage(v *redis.Message) {
	callStatus := common.NewCallStatus(v.Payload)
	if callStatus.CallStatus == common.CallStatusCompleted {
		callInfo, err := restcommApi.GetCallInfo(callStatus.CallSid)
		if err != nil {
			common.Error.Println("opencell-billing: Failed to query restcomm for call sid:%s", callStatus.CallSid)
		} else {
			duration := callInfo.Duration
			timestamp := time.Now()
			timestring := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.000Z",
				timestamp.Year(),
				timestamp.Month(),
				timestamp.Day(),
				timestamp.Hour(),
				timestamp.Minute(),
				timestamp.Second())
			opencellApi.ChargeCustomer(opencellUser, timestring, float64(duration)/60.0 )
		}
	}
}