package main

import "tad-demo/common"

var cfg = common.NewConfig()

// This module collects incomming phone numbers from the database and adds them to the conference call
// We need to have webserver to handle callback urls from RestComm and provide commands for it.

func main() {

	conference := Conference{}

	webServer := WebServer{callback: conference}
	webServer.Start()

	conference.RegisterNumber()

	conference.Subscribe()

    common.WaitCtrlC()
	common.Info.Println("Finished")
}
