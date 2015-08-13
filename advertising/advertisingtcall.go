package main
import (
	"fmt"
	common "tad-demo/common"
)

var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)
var cfg = common.NewConfig()
var db = common.NewDbClient(cfg.Service.Redis)

var pendingCalls = make(map[string]*AdvertisingCall)

type Variant struct{
	text string
	confirmation string
}

func (v *Variant)Text(text string)(*Variant){
	v.text = text
	return v
}

func (v *Variant)Confirmation(text string)(*Variant){
	v.confirmation = text
	return v
}

type AdvertisingCall struct{
	Sid string
	Variants map[int]*Variant

	Question string
}

func NewAdvertisingCall(participant string)(*AdvertisingCall){
	return &AdvertisingCall{Sid: participant, Variants: make(map[int]*Variant)}
}

func (call *AdvertisingCall)Variant(val int)(*Variant){
	call.Variants[val] = &Variant{}
	return call.Variants[val]
}

func (call *AdvertisingCall)Other()(*Variant){
	call.Variants[-1] = &Variant{}
	return call.Variants[-1];
}

func (call *AdvertisingCall)Prompt(question string){
	call.Question = question;
}

func (call AdvertisingCall)Exec(){
	url := fmt.Sprintf("http://%s/play-sound.xml", cfg.GetExternalAddress(cfg.ServerPort.Advertising))

	to := db.Get(call.Sid).Val()
	if to == ""{
		fmt.Println("\t to is EMPTY for", call.Sid)
		return
	}

	pendingCalls[common.GetClientName(to)] = &call
	restcommApi.MakeCall(cfg.Messages.DialFrom, to, url)
}