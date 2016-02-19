package main

import (
	"fmt"
	"encoding/json"
	common "tad-demo/common"
)

var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)
var db = common.NewDbClient(cfg.Service.Redis)

type CallStatus struct {
	To string
	CallStatus string
	CallSid string
}

type Conference struct {
}

func (conf* Conference) RegisterNumber() {
	fmt.Println("\tRegister number:", cfg.Callback.Phone, cfg.Callback.Conference)

	callBack := fmt.Sprintf("http://%s/register", cfg.GetExternalAddress(cfg.ServerPort.Conference))
	common.NewIncomingPhoneNumber("", cfg.Callback.Phone).CreateOrUpdate(restcommApi, callBack)

	common.NewIncomingPhoneNumber("", cfg.Callback.Conference).CreateOrUpdate(restcommApi, "")
}

func (conf* Conference) Subscribe() {
	common.Info.Println("Start conference hub app v.2.0")
	/*conf.incomingChannel = make(chan string, 100)
	conf.statusChannel = make(chan CallStatus, 100)

	go func() {
		for {
			from := <- conf.incomingChannel
			conf.Add(from)
		}
	}()

	go func() {
		for {
			callsStatus := <- conf.statusChannel
			conf.FireCallStatus(callsStatus)
		}
	}()*/
}

func (conf* Conference) FireCallStatus(callStatus CallStatus) {
	bytes, _ := json.Marshal(&callStatus)
	db.Publish("callStatus", string(bytes))
	/*if callStatus == "in-progress" {
		//common.Trace.Println("add in-progress call to db:", callSid, "to =", to)
		//db.Set(cfg.Redis.InProgressKey+":"+callSid, callSid, 1*time.Hour)
	} else if callStatus == "completed" {
		common.Trace.Println("add completed call to stream:", callSid, "to =", to)
		//db.Publish(cfg.Redis.ConfChannel, to)
	}*/
}

func (conf* Conference) Add(from string) {
	to := common.ConvertToSipCall(from, cfg.Sip.DidProvider)

	/*call, err := restcommApi.MakeCall(cfg.Messages.DialFrom, to,
		fmt.Sprintf("http://%s/make-conference.xml", cfg.GetExternalAddress(cfg.ServerPort.Conference)),
		fmt.Sprintf("http://%s/call-status.xml", cfg.GetExternalAddress(cfg.ServerPort.Conference)))

	if err != nil {
		fmt.Println("ERROR: Call to", to, " with erorr", err)
		return
	}*/

	//db.RPush(common.DB_KEY_URI, call.Uri)
	//db.Set(call.Sid, from, 0)

	key := "conf:" + from
	db.Set(key, from, 0)
}

func (conf Conference) HandleCallStatusChanged(to string, callStatus string, callSid string){
	common.Trace.Println("call status chnaged: ", callSid, "to =", to, "status=", callStatus)
	go func() {
		conf.FireCallStatus(CallStatus{To: to, CallStatus: callStatus, CallSid: callSid})
	}()
}

func (conf Conference) HandleIncomingCall(from string) string {
	common.Trace.Println("HandleIncomingCall:", from)
	if from != "" {
		go func() {
			conf.Add(from)
		}()
	}
	resp := fmt.Sprintf("<Response><Say>%s</Say><Dial><Conference startConferenceOnEnter=\"true\">%s</Conference></Dial></Response>",
				cfg.Messages.ConferenceWelcome, cfg.Callback.Conference)
	common.Trace.Println("HandleIncomingCall resp:", resp)
	return resp
}




