package main
import "tad-demo/common"

var cfg = common.NewConfig()

func main() {

	webServer := WebServer{}
	webServer.Start()

	conference := Conference{}
    conference.RegisterNumber(cfg.Callback.Conference)

	storage := Storage{}
	subscription := storage.Subscribe()

	for{
		pstn := subscription.Receive()
		conference.Add(pstn)
	}
}
