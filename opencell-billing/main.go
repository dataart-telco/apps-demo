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
	go func() {
		callListener.Subscribe()
	}()
	go func() {
		callListener.SubscribeCallerList()
	}()
	common.WaitCtrlC()
}
