package main

import (
	"tad-demo/common"
)

// Drop current conference call

func main() {
	common.Info.Println("Opencell billing")
	webserver := OpencellWebServer{}
	webserver.Start()
	callListener := BillingListener{}
	callListener.Subscribe()
	callListener.SubscribeCallerList()
}
