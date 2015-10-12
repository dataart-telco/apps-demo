package main

import "tad-demo/common"

var cfg = common.NewConfig()

// This module collects incomming phone numbers from the database and adds them to the conference call
// We need to have webserver to handle callback urls from RestComm and provide commands for it.

func main() {

	webServer := WebServer{StatusHandler: SmsCallLink{}}
	webServer.Start()

	conference := Conference{}
	conference.RegisterNumber(cfg.Callback.Conference)

	storage := Storage{}
	subscription := storage.Subscribe()

	for {
		pstn := subscription.Receive()
		conference.Add(pstn)

	}
}
