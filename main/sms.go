package main
import (
	"fmt"
	"net/http"
	"tad-demo/common"
)

var cfg = common.NewConfig()
var db = common.NewDbClient(cfg.Service.Redis)
var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)

type Subscription struct{
	acceptedQueue chan string
}

func (sms Subscription)Receive()(string){
	return <-sms.acceptedQueue
}

type Sms struct{
	subscription Subscription
}

func (sms Sms)RegisterNumber(number string){
	fmt.Println("\tRegister number:", number)

	callBack := fmt.Sprintf("http://%s/register", cfg.GetExternalAddress(cfg.ServerPort.Main))
	common.NewIncomingPhoneNumber("", cfg.Callback.Phone).CreateOrUpdate(restcommApi, callBack)
}

func (sms Sms)Subscribe()(Subscription){
	fmt.Println("Start main web app v.1.0")

	sms.subscription = Subscription{acceptedQueue: make(chan string, 100)}

	go func() {
		http.HandleFunc("/", sms.handler)
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort.Main), nil)
		if (err != nil) {
			panic(err)
		}
	}()
	return sms.subscription;
}

func (sms Sms)handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<Response><Hangup/></Response>")

	from := r.PostFormValue("From")

	fmt.Println("\tReceive ", r.Method, " call from ", from)

	from = common.ConvertToSip(from)
	sms.subscription.acceptedQueue <- from
}

func (sms Sms)Await(){
    fmt.Println("wait for ctrl+c")
	common.WaitCtrlC()
}

type Storage struct{

}

func (storage Storage)Save(from string){
	key := "conf:" + from;
	db.Set(key, from, 0)
	db.Publish(cfg.Redis.Channel, from)
}
