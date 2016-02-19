package main

import (
	"fmt"
	"net/http"
	"tad-demo/common"
	"time"
)

type WebServer struct {
}

func (w WebServer) Start() {
	go w.startWebServer()

	url := fmt.Sprintf("http://%s/ping.xml", cfg.GetExternalAddress(cfg.ServerPort.Conference))
	for i := 0; i < 15; i++ {
		fmt.Println("Wait until server is ready...")
		time.Sleep(1 * time.Second)

		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			return
		}
	}
}

func (self WebServer) startWebServer() {
	fmt.Println("\tStart conference web server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("\t<- http request:", r.URL)

		w.Header().Set("Content-Type", "text/xml")
		if r.URL.Path == "/dial-status.xml" {
			//TODO handle dial status ???
			common.Info.Println("\t<- http request:", r.URL)
		} else if r.URL.Path == "/call-status.xml" {
			common.Info.Println("\t<- http request:", r.URL)
			to := r.FormValue("To")
			callStatus := r.FormValue("CallStatus")
			callSid := r.FormValue("CallSid")
			go func(){
				self.handleCallStatusChanged(to, callStatus, callSid)
			}()
		} else {
			fmt.Fprintf(w,
				"<Response><Say>%s</Say><Dial><Conference startConferenceOnEnter=\"true\">%s</Conference></Dial></Response>",
				cfg.Messages.ConferenceWelcome, cfg.Callback.Conference)
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Conference), nil)
	if err != nil {
		panic(err)
	}
}

func (self WebServer) handleCallStatusChanged(to string, callStatus string, callSid string) {
	status := common.CallStatus{To: to, CallStatus: callStatus, CallSid: callSid}
	db.Publish(common.CHANNEL_CALL_STATUS, status.ToJson())
}
