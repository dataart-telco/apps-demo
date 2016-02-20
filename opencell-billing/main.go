package main

import (
	"tad-demo/common"
)

// Drop current conference call

func main() {
	common.Info.Println("Opencell billing")
	callListener := BillingListener{}
	callListener.Subscribe()
}
