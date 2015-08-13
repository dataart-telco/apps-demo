package main

func main() {

	storage := Storage{}
	sms := Sms{}
    sms.RegisterNumber(cfg.Callback.Phone)

	subscription := sms.Subscribe()

	for{
		pstn := subscription.Receive()
		storage.Save(pstn)
	}
}
