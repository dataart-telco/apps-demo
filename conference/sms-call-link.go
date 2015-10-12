package main

import (
	"fmt"
)

type SmsCallLink struct {
}

func (sms SmsCallLink) HandleStatusChanged(sid string, to string, status string) error {
	recId, err := restcommApi.GetRecording(sid)
	if err != nil {
		return err
	}
	recUrl := fmt.Sprintf("http://%s:8080/restcomm/recordings/%s.wav", cfg.GetExternalAddress(cfg.ServerPort.Conference), recId)
	restcommApi.SendSms(to, fmt.Sprintf("Call record here %s", recUrl))
	return nil
}
