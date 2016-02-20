package main

import (
	"fmt"
	"net/http"
	"tad-demo/common"
	"time"
)

type BillingWebServer struct {
}

func (w BillingWebServer) Start() {
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

func (s BillingWebServer) startWebServer() {
	fmt.Println("\tStart conference web server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("\t<- http request:", r.URL)

		w.Header().Set("Content-Type", "text/xml")
		if r.URL.Path == "/conference" {
			fmt.Sprintf(w,
				`<!DOCTYPE html>
				<html>
				  <head>
				    <title>Opensell statistics</title>
				    <style type="text/css"> html,body,#tbl_wrap{height:100%;width:100%;padding:0;margin:0}#td_wrap{vertical-align:middle;text-align:center}</style>
				  </head>
				<body>
				  <table id="tbl_wrap"><tbody><tr><td id="td_wrap">
				    <div id="centered_div1"><font size="40" face="Arial">
				      Conference call ended
				    </font></div>
				    <div id="centered_div2"><font size="40" face="Arial">
				      Total attenders: %d
				    </font></div>
				    <div id="centered_div3"><font size="40" face="Arial">
				      Call price: %d$
				    </font></div>
				  </td></tr></tbody></table>
				</body>
				</html>`,
			0, 0)
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Conference), nil)
	if err != nil {
		panic(err)
	}
}

