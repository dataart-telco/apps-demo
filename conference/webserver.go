package main
import (
	"fmt"
	"net/http"
)

type WebServer struct{

}

func (w WebServer)Start(){
	go w.startWebServer()
}

func (w WebServer)startWebServer() {
	fmt.Println("\tStart conference web server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Println("\t<- http request:", r.URL)
		fmt.Fprintf(w, "<Response><Say>%s</Say><Dial><Conference startConferenceOnEnter=\"true\">%s</Conference></Dial></Response>", cfg.Messages.ConferenceWelcome, cfg.Callback.Conference)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Conference), nil)
	if(err != nil){
		panic(err)
	}
}
