package main

import (
	"fmt"
	"net/http"
	"time"
)

type OpencellWebServer struct {
}

func (w OpencellWebServer) Start() {
	go w.startWebServer()

	url := fmt.Sprintf("http://%s/ping.xml", cfg.GetExternalAddress(8080))
	for i := 0; i < 15; i++ {
		fmt.Println("Wait until server is ready...")
		time.Sleep(1 * time.Second)

		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			return
		}
	}
}

func (self OpencellWebServer) startWebServer() {
	fmt.Println("\tStart conference web server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("\t<- http request:", r.URL)

		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/statistics.json" {
			jsonString := db.Get(JSON_STATS_KEY).String()
			fmt.Fprint(w, jsonString)
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
	if err != nil {
		panic(err)
	}
}

