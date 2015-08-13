package main
import (
	"time"
	"tad-demo/common"
)

func main() {

	webServer := NewWebServer()
	webServer.Start()

	conference := Conference{}

	participants := conference.GetParticipants()

	conference.Drop()

	time.Sleep(5 * time.Second)
	for _, participant := range participants{

		call := NewAdvertisingCall(participant)
		call.Prompt(cfg.Messages.Question)

		call.Variant(1).Text(cfg.Messages.Answer1).Confirmation(cfg.Messages.ThanksForAnswer)
		call.Other().Text(cfg.Messages.ThanksForAttention)
		call.Exec()
	}

	common.WaitCtrlC()
}
