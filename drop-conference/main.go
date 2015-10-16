package main

import (
	"tad-demo/common"
)

// Drop current conference call

func main() {
	common.Info.Println("Drop conference")

	conference := Conference{}

	participants := conference.GetParticipants()

	common.Info.Println("Participants count", len(participants))

	conference.Drop()

	common.Info.Println("Press Ctrl+C")

	common.WaitCtrlC()
}
