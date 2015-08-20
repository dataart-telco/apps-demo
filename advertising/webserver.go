package main
import (
	"time"
    "encoding/json"
	"net/http"
	"fmt"
	"strings"
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
    
    w.Header().Set("Content-Type", "text/xml")

	if strings.HasPrefix(r.URL.Path, PATH_GATHER){
		web.handleAnswer(w, r)
		return
	}
	web.handleAskQuestion(w, r)
}

func (web WebServer)handleAnswer(w http.ResponseWriter, r *http.Request){
	phone := r.FormValue("To")
	digits := r.FormValue("Digits")

	if digits == "" {
		fmt.Fprintf(w, "<Response><Say>%s</Say><Hangup/></Response>", cfg.Messages.ThanksForAttention)
		return
	}

	scenario := getPendingCall(phone)

	if scenario == nil {
		fmt.Fprintf(w, "<Response><Say>%s</Say><Hangup/></Response>", cfg.Messages.ThanksForAttention)
		return
	}

	variant, ok := scenario.Variants[digits]
	if !ok {
		fmt.Fprintf(w, "<Response><Say>%s</Say><Hangup/></Response>",  cfg.Messages.ThanksForAttention)
		return
	}
	db.RPush(KEY_GATHER, phone)

	fmt.Fprintf(w, "<Response><Say>%s</Say><Hangup/></Response>",  variant.ConfMessage)
}

func (web WebServer)handleAskQuestion(w http.ResponseWriter, r *http.Request){

	phone := r.FormValue("To")
	scenario := getPendingCall(phone)

	if scenario == nil {
		fmt.Fprintf(w, "<Response><Say>%s</Say></Response>", cfg.Messages.ThanksForAttention)
		return
	}

	resp := fmt.Sprintf("<Say>%s</Say>", scenario.Question)
	for k, v := range scenario.Variants{
		if k == UNDEFINED_VARIANT {
			continue
		}
		resp += fmt.Sprintf("<Say>%s</Say>", v.Message)
	}

	another := ""

	variant, ok := scenario.Variants[UNDEFINED_VARIANT]
	if ok {
		another = fmt.Sprintf("<Say>%s</Say>", variant.Message)
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

    url := fmt.Sprintf("http://%s/play-sound.xml", cfg.GetExternalAddress(cfg.ServerPort.Advertising))
    for i := 0; i < 15; i++ {
        fmt.Println("Wait until server is ready...")        
        time.Sleep(1 * time.Second)

        resp, err := http.Get(url)
	    if(err == nil && resp.StatusCode == 200){
		    return
	    }
    }
}

func getPendingCall(client string)(*AdvertisingCall){
    jStr := db.Get(PENDING_PREFIX + client).Val()

    if jStr == "" {
        return nil
    }

    result := AdvertisingCall{}
    err := json.Unmarshal([]byte(jStr), &result)
    if err != nil {
        return nil
    }
    return &result
}
