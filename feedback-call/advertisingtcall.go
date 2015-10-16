package main

import (
	"encoding/json"
	"fmt"
	common "tad-demo/common"
	"time"
)

var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)
var cfg = common.NewConfig()
var db = common.NewDbClient(cfg.Service.Redis)

const PENDING_PREFIX = "pending:"
const UNDEFINED_VARIANT = "-1"

type Variant struct {
	Message     string
	ConfMessage string
}

func (v *Variant) Text(text string) *Variant {
	v.Message = text
	return v
}

func (v *Variant) Confirmation(text string) *Variant {
	v.ConfMessage = text
	return v
}

type AdvertisingCall struct {
	Sid      string
	Variants map[string]*Variant

	Question string
}

func NewAdvertisingCall(participant string) *AdvertisingCall {
	return &AdvertisingCall{Sid: participant, Variants: make(map[string]*Variant)}
}

func (call *AdvertisingCall) Variant(val string) *Variant {
	call.Variants[val] = &Variant{}
	return call.Variants[val]
}

func (call *AdvertisingCall) Other() *Variant {
	call.Variants[UNDEFINED_VARIANT] = &Variant{}
	return call.Variants[UNDEFINED_VARIANT]
}

func (call *AdvertisingCall) Prompt(question string) {
	call.Question = question
}

func (call AdvertisingCall) Exec() {
	url := fmt.Sprintf("http://%s/play-sound.xml", cfg.GetExternalAddress(cfg.ServerPort.Advertising))

	to := db.Get(call.Sid).Val()
	if to == "" {
		fmt.Println("\t to is EMPTY for", call.Sid)
		return
	}

	bytes, err := json.Marshal(&call)
	if err != nil {
		fmt.Println("\tcan't convert to json", call.Sid)
		fmt.Println("\tERROR:", err)
		return
	}

	db.Set(PENDING_PREFIX+common.GetClientName(to), string(bytes), 1*time.Hour)

	restcommApi.MakeCall(cfg.Messages.DialFrom, to, url, "")
}
