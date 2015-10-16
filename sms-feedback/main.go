package main

func main() {
	smsService := SmsService{}
	smsService.Start()

	conference := Conference{}

	subscription := conference.Subscribe()
	for {
		sid := subscription.Receive()
		smsService.SendSms(sid)
	}
}
