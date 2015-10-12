package main

import (
	"fmt"
	"net/http"
	"strings"
	"tad-demo/common"
	"time"
)

type WebServer struct {
	StatusHandler CallStatusHandler
}

type CallStatusHandler interface {
	HandleStatusChanged(sid string, to string, status string) error
}

func (w WebServer) Start() {
	go w.startWebServer()

	url := fmt.Sprintf("http://%s/ping.xml", cfg.GetExternalAddress(cfg.ServerPort.Conference))
	for i := 0; i < 15; i++ {
		common.Info.Println("Wait until server is ready...")
		time.Sleep(1 * time.Second)

		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			return
		}
	}
}

func (web WebServer) startWebServer() {
	common.Info.Println("\tStart conference web server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		common.Trace.Println("\t<- http request:", r.URL)

		w.Header().Set("Content-Type", "text/xml")
		if strings.HasPrefix(r.URL.Path, "/dial-status.xml") {
			r.ParseForm()
			common.Trace.Println("dial-status: Form: ", r.Form)
		} else if strings.HasPrefix(r.URL.Path, "/call-status.xml") {
			to := r.FormValue("To")
			callStatus := r.FormValue("CallStatus")
			callSid := r.FormValue("CallSid")

			common.Trace.Println("call-status: Form: ", r.Form)
			common.Trace.Println("to=", to, "| callStatus=", callStatus, "| callSid=", callSid)
			if callStatus == "completed" && web.StatusHandler != nil {
				web.StatusHandler.HandleStatusChanged(callSid, to, callStatus)
			}
		} else {
			fmt.Fprintf(w,
				//				"<Response><Say>%s</Say><Dial record=\"true\"><Conference startConferenceOnEnter=\"true\">%s</Conference></Dial></Response>",
				"<Response><Say>Test recording</Say><Record maxLength=\"30\"/></Response>")
			//cfg.Messages.ConferenceWelcome,
			//cfg.Callback.Conference)
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Conference), nil)
	if err != nil {
		panic(err)
	}
}
