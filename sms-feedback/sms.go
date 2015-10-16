package main

import (
	"fmt"
	"net/http"
	"strings"
	"tad-demo/common"
	"time"
)

var cfg = common.NewConfig()
var db = common.NewDbClient(cfg.Service.Redis)
var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)

type SmsService struct {
}

func (s SmsService) SendSms(to string) {
	common.Trace.Println("Send sms to", to)
	//TODO cfg.Callback.Sms
	err := restcommApi.SendSms(to, "DataArt", cfg.Messages.SmsMessage)
	if err != nil {
		common.Error.Println("Send sms error", err)
	}
}

func (s SmsService) Start() {
	go s.startWebServer()

	url := fmt.Sprintf("http://%s/ping.xml", cfg.GetExternalAddress(cfg.ServerPort.Sms))
	for i := 0; i < 15; i++ {
		common.Info.Println("Wait until server is ready...")
		time.Sleep(1 * time.Second)

		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			return
		}
	}
}

func (s SmsService) startWebServer() {
	common.Info.Println("\tStart web server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		common.Trace.Println("\t<- http request:", r.URL)

		if strings.HasPrefix(r.URL.Path, "/call-status.xml") {
			to := r.FormValue("To")
			callStatus := r.FormValue("CallStatus")
			callSid := r.FormValue("CallSid")

			common.Trace.Println("call-status: Form: ", r.Form)
			common.Trace.Println("to=", to, "| callStatus=", callStatus, "| callSid=", callSid)
			if callStatus == "completed" {
				s.SendSms(to)
			}
			return
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Sms), nil)
	if err != nil {
		panic(err)
	}
}
