package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"tad-demo/common"
)

var cfg = common.NewConfig()
var db = common.NewDbClient(cfg.Service.Redis)
var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)

type Subscription struct {
	acceptedQueue chan string
}

func (sms Subscription) Receive() string {
	return <-sms.acceptedQueue
}

type Sms struct {
	subscription Subscription
}

func (sms Sms) RegisterNumber(number string) {
	common.Info.Println("\tRegister number:", number)

	callBack := fmt.Sprintf("http://%s/register", cfg.GetExternalAddress(cfg.ServerPort.Main))
	common.NewIncomingPhoneNumber("", cfg.Callback.Phone).CreateOrUpdate(restcommApi, callBack)
}

func (sms Sms) Subscribe() Subscription {
	common.Info.Println("Start main web app v.1.0")

	sms.subscription = Subscription{acceptedQueue: make(chan string, 100)}

	go func() {
		http.HandleFunc("/", sms.handler)
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Main), nil)
		if err != nil {
			panic(err)
		}
	}()
	return sms.subscription
}

func (sms Sms) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")

	if r.URL.Path == "/test.xml" {
		var buffer bytes.Buffer
		for i := 0; i < 10000; i++ {
			buffer.WriteString("long text here")
		}
		fmt.Fprintf(w, "<Test>%x</Test>", md5.Sum(buffer.Bytes()))
		buffer.Reset()
		return
	}
	fmt.Fprintf(w, "<Response><Hangup/></Response>")

	from := r.PostFormValue("From")

	common.Trace.Println("\tReceive ", r.Method, " call from ", from)
	if from != "" {
		from = common.ConvertToSip(from)
		sms.subscription.acceptedQueue <- from
	}
}

func (sms Sms) Await() {
	common.Info.Println("wait for ctrl+c")
	common.WaitCtrlC()
}

type Storage struct {
}

func (storage Storage) Save(from string) {
	key := "conf:" + from
	db.Set(key, from, 0)
	db.Publish(cfg.Redis.MainChannel, from)
}
