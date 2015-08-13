package main
import (
	"net/http"
	"fmt"
	"strings"
	"strconv"
)

const PATH_GATHER  = "/gather.html"

const KEY_GATHER = "gather"

type HttpHandler func(http.ResponseWriter, *http.Request)

type WebServer struct{
}

func NewWebServer()(WebServer){
    return WebServer{}
}

func (web WebServer)httpHandlerPlay(w http.ResponseWriter, r *http.Request){
	fmt.Println("\t<- http request -", r.URL)

	if strings.HasPrefix(r.URL.Path, PATH_GATHER){
		web.handleAnswer(w, r)
		return
	}
	web.handleAskQuestion(w, r)
}

func (web WebServer)handleAnswer(w http.ResponseWriter, r *http.Request){
	phone := r.FormValue("To")
	digits, err := strconv.Atoi(r.FormValue("Digits"))

	if err != nil {
		fmt.Fprintf(w, "<Response><Say>%s</Say><Hangup/></Response>", cfg.Messages.ThanksForAttention)
		return
	}

	scenario, ok := pendingCalls[phone]

	if !ok {
		fmt.Fprintf(w, "<Response><Say>%s</Say><Hangup/></Response>", cfg.Messages.ThanksForAttention)
		return
	}

	variant, ok := scenario.Variants[digits]
	if !ok {
		fmt.Fprintf(w, "<Response><Say>%s</Say><Hangup/></Response>",  cfg.Messages.ThanksForAttention)
		return
	}
	db.RPush(KEY_GATHER, phone)

	fmt.Fprintf(w, "<Response><Say>%s</Say><Hangup/></Response>",  variant.confirmation)
}

func (web WebServer)handleAskQuestion(w http.ResponseWriter, r *http.Request){

	phone := r.FormValue("To")
	scenario, ok := pendingCalls[phone]

	if !ok{
		fmt.Fprintf(w, "<Response><Say>%s</Say></Response>", cfg.Messages.ThanksForAttention)
		return
	}

	resp := fmt.Sprintf("<Say>%s</Say>", scenario.Question)
	for k, v := range scenario.Variants{
		if k == -1 {
			continue
		}
		resp += fmt.Sprintf("<Say>%s</Say>", v.text)
	}

	another := ""

	variant, ok := scenario.Variants[-1]
	if ok {
		another = fmt.Sprintf("<Say>%s</Say>", variant.text)
	}

	fmt.Fprintf(w,
		"<Response><Gather action=\"%s\" method=\"GET\" timeout=\"5\" numDigits=\"1\">%s</Gather>%s<Hangup/></Response>",
		fmt.Sprintf("http://%s%s", cfg.GetExternalAddress(cfg.ServerPort.Advertising), PATH_GATHER),
		resp,
		another)
}

func (web WebServer) Start() {
	go func() {
		fmt.Println("Start advertising web server")

		http.HandleFunc("/", web.httpHandlerPlay)
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Advertising), nil)
		if (err != nil) {
			panic(err)
		}
	}()
}
