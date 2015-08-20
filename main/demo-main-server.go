package main

// This module registers phone number in RestComm 
// and collects all incomming phone numbers to the databse.

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
