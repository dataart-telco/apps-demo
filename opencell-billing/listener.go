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
	Sid string `json:"sid"`
	Time int `json:"duration"`
}
func NewCaller(sid string, time int) Caller {
	return Caller {
		Sid : sid,
		Time : time}
}
type ConfStats struct {
	Sum float64 `json:"sum"`
	Callers []Caller `json:"callers"`
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
	_, err := db.Get("last_conf_time").Result()
	if err != nil {
		layout := "2006-01-02T15:04:05.000Z"
		db.Set("last_conf_time", time.Date(1970, 1, 1, 1, 1, 1, 1, time.Now().Location()).Format(layout),0)
	}
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
	cList := make([]string, 0)
	json.Unmarshal([]byte(v.Payload), &cList)
	callerList := ConfStats{}
	total := 0
	for _, element := range cList {
		callInfo, err := restcommApi.GetCallInfo(element)
		if err != nil {
			common.Error.Println("opencell-billing: Failed to query restcomm for call sid:%s", element)
		} else {
			duration := callInfo.Duration
			total += duration
			caller := NewCaller(element, duration)
			callerList.Callers = append(callerList.Callers, caller)
		}
	}
	layout := "2006-01-02T15:04:05.000Z"
	last_time_str, err := db.Get("last_conf_time").Result()
	var lastConfTime time.Time
	if err != nil {
		lastConfTime = time.Date(1970, 1, 1, 1, 1, 1, 1, time.Now().Location())
	} else {
		lastConfTime, _ = time.Parse(layout, last_time_str)
	}

	currentTime := time.Now()
	callerList.Sum = opencellApi.GetBalanceWithRange(opencellUser, lastConfTime, currentTime)
	jsonStr, _ := json.Marshal(callerList)
	db.Set(JSON_STATS_KEY, jsonStr,0)
	db.Set("last_conf_time", time.Now().Format(layout),0)
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