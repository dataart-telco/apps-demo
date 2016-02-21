package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"tad-demo/common"
	"time"
)

var cfg = common.NewConfig()
var db = common.NewDbClient(cfg.Service.Redis)
var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)
var truPhone = common.Truphone{}
type SmsService struct {
}

func (s SmsService) SendSms(to string) {
	common.Info.Println("Send sms to", to)

	to = common.ConvertToSipSms(to, cfg.Sip.DidProvider)

	common.Trace.Println("To is converted:", to)
	//TODO cfg.Callback.Sms
	record := s.GetRecordUrl()
	sms := cfg.Messages.SmsMessageSimple
	if record != "" {
		record = common.GetShortLink(record)
		if record != "" {
			sms = fmt.Sprintf(cfg.Messages.SmsMessage, record)
		}
	}
	common.Trace.Println("Sms message:", sms)

	if common.IsPhoneNumber(to) {
		common.Info.Println("use truPhone", to)
		err := truPhone.SendSms(to, "DataArt", sms)
		if err != nil {
			common.Error.Println("Send sms to truphone error", err)
		}
	} else {
		err := restcommApi.SendSms(to, "DataArt", sms)
		if err != nil {
			common.Error.Println("Send sms to restcomm error", err)
		}
	}
}

func (s SmsService) GetRecordUrl() string {
	if cfg.Service.Recorder == "" {
		return ""
	}
	host := "http://" + cfg.Service.Recorder
	resp, err := http.Get(host + "/last")
	if err != nil {
		common.Error.Println("ger record file name from", host, "error", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	body, _ := ioutil.ReadAll(resp.Body)
	wav := string(body)
	return host + "/" + wav
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
	common.Info.Println("\tStart web server: recorder =", cfg.Service.Recorder)

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
