package main

import (
	"fmt"
	common "tad-demo/common"
)

var restcommApi = common.NewRestcommApi(cfg.Service.Restcomm, cfg.Auth.User, cfg.Auth.Pass)

type Conference struct {
}

func (conf Conference) RegisterNumber(phone string) {
	common.Info.Println("\tRegister number:", phone)

	common.NewIncomingPhoneNumber("", cfg.Callback.Conference).CreateOrUpdate(restcommApi, "")
}

func (conf Conference) Add(client string) {
	to := common.ConvertToRestcommNumber(client)

	call, err := restcommApi.MakeCallWithStatus(
		cfg.Messages.DialFrom,
		to,
		fmt.Sprintf("http://%s/make-conference.xml", cfg.GetExternalAddress(cfg.ServerPort.Conference)),
		fmt.Sprintf("http://%s/call-status.xml", cfg.GetExternalAddress(cfg.ServerPort.Conference)))

	if err != nil {
		fmt.Println("ERROR: Call to", to, " with erorr", err)
		return
	}

	db.RPush(common.DB_KEY_URI, call.Uri)
	db.Set(call.Sid, to, 0)
}
