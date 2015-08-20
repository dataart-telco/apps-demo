package main
import (
	"tad-demo/common"
)

// This module gets participants of current conference, drops it and makes advertising call to them
// We need to have webserver to handle callback urls from RestComm and provide commands for it.

func main() {

	webServer := NewWebServer()
	webServer.Start()

	conference := Conference{}

	participants := conference.GetParticipants()

	conference.Drop()

	for _, participant := range participants{

		call := NewAdvertisingCall(participant)
		call.Prompt(cfg.Messages.Question)

		call.Variant("1").Text(cfg.Messages.Answer1).Confirmation(cfg.Messages.ThanksForAnswer)
		call.Other().Text(cfg.Messages.ThanksForAttention)
		call.Exec()
	}

	common.WaitCtrlC()
}
