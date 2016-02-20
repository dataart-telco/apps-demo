package main

import (
	"gopkg.in/redis.v3"
	"tad-demo/common"
	"tad-demo/common/opencell-api"
	"time"
	"fmt"
	"encoding/json"
)

const (
	opencellUser = "DataArt Conference"
	JSON_STATS_KEY="json_stats"
)
var cfg = common.NewConfig()
var db = common.NewDbClient(cfg.Service.Redis)
var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)
var opencellApi = opencell_api.NewOpencellAPI(cfg.Opencell.Password, cfg.Opencell.Host)

type Caller struct {
	Sid string
	Time int
}
func NewCaller(sid string, time int) Caller {
	return Caller {
		Sid : sid,
		Time : time}
}
type ConfStats struct {
	Sum float64
	Callers []Caller
}

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

func (l BillingListener) SubscribeCallerList() {
	//subscribe for events from call_status queue
	opencellApi.InitOpenCell()
	opencellApi.CreateNewCustomer(opencellUser)

	sub, _ := db.Subscribe(common.CHANNEL_CONF_DROPPED)
	for {
		msg, e2 := sub.Receive()
		if e2 != nil {
			common.Error.Println("opencell-billing: Error receiving message from Redis 'conf_dropped' queue")
			panic(e2)
		}
		switch v := msg.(type) {
		case *redis.Message:
			go func(){
				l.ParseCallerList(v)
			}()
		}
	}
}

func (l BillingListener) ParseCallerList(v *redis.Message) {
	cList := make([]String, 0)
	json.Unmarshal([]byte(v.Payload), &cList)
	callerList := ConfStats{}
	total := 0
	for _, element := range cList {
		callInfo, err := restcommApi.GetCallInfo(element)
		if err != nil {
			common.Error.Println("opencell-billing: Failed to query restcomm for call sid:%s", callStatus.CallSid)
		} else {
			duration := callInfo.Duration
			total += duration
			caller := NewCaller(element, duration)
			callerList.Callers = append(callerList.Callers, caller)
		}
	}
	callerList.Sum = float64(total)
	jsonStr, _ := json.Marshal(callerList)
	db.Set(JSON_STATS_KEY, jsonStr)
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