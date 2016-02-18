package main

import (
	"fmt"
	"net/http"
	"time"
	"bytes"
	"crypto/md5"
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
			fmt.Println("\t<- http request:", r.URL)
			to := r.FormValue("To")
			callStatus := r.FormValue("CallStatus")
			callSid := r.FormValue("CallSid")
			self.callback.HandleCallStatusChanged(to, callStatus, callSid)
		} else {
			fmt.Println("\t<- http request:", r.URL)
			from := r.FormValue("From")
			resp := self.callback.HandleIncomingCall(from)
			fmt.Fprint(w, resp)
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Conference), nil)
	if err != nil {
		panic(err)
	}
}
