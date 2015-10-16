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

func (w WebServer) startWebServer() {
	fmt.Println("\tStart conference web server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("\t<- http request:", r.URL)

		w.Header().Set("Content-Type", "text/xml")
		if r.URL.Path == "/call-status.xml" {
			to := r.FormValue("To")
			callStatus := r.FormValue("CallStatus")
			callSid := r.FormValue("CallSid")
			if callStatus == "in-progress" {
				//common.Trace.Println("add in-progress call to db:", callSid, "to =", to)
				//db.Set(cfg.Redis.InProgressKey+":"+callSid, callSid, 1*time.Hour)
				//db.Publish(cfg.Redis.ConfChannel, callSid)
			} else if callStatus == "completed" {
				common.Trace.Println("add completed call to stream:", callSid, "to =", to)
				db.Publish(cfg.Redis.ConfChannel, to)
			}
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
