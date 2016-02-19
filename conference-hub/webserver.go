package main

import (
	"fmt"
	"net/http"
	"time"
	"bytes"
	"crypto/md5"
	common "tad-demo/common"
)

type WebServerCallback interface{
	HandleCallStatusChanged(to string, callStatus string, callSid string)
	HandleIncomingCall(from string) string 
}

type WebServer struct {
	callback WebServerCallback
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

func (self WebServer) startWebServer() {
	common.Info.Println("\tStart conference web server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		if r.URL.Path == "/test.xml" {
			var buffer bytes.Buffer
			for i := 0; i < 100; i++ {
				buffer.WriteString("long text here")
			}
			fmt.Fprintf(w, "<Test>%x</Test>", md5.Sum(buffer.Bytes()))
			buffer.Reset()
			return
		}else if r.URL.Path == "/call-status.xml" {
			common.Info.Println("\t<- http request:", r.URL)
			to := r.FormValue("To")
			callStatus := r.FormValue("CallStatus")
			callSid := r.FormValue("CallSid")
			self.callback.HandleCallStatusChanged(to, callStatus, callSid)
		} else if r.URL.Path == "/register"{
			common.Info.Println("\t<- register request:", r.URL)
			from := r.PostFormValue("From")
			common.Info.Println("\t<- http request:", r.URL, from)
			resp := self.callback.HandleIncomingCall(from)
			fmt.Fprint(w, resp)
		} else {
			common.Info.Println("\t<- http request:", r.URL)
			fmt.Fprintf(w, "<Response><Hangup/></Response>")
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Conference), nil)
	if err != nil {
		panic(err)
	}
}
