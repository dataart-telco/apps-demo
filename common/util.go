package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func WaitCtrlC() {
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 2)
	signal.Notify(signal_channel, os.Interrupt, syscall.SIGTERM)
	<-signal_channel
}

func MakeJsonPost(url string, jsonBody string) (string, error) {
	Trace.Println("send body: ", jsonBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Trace.Println("Send request error", err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	Trace.Println("Response body", string(body))

	bundle := make(map[string]interface{})
	defer func() {
		r := recover()
		if r != nil {
			Error.Println("Catch panic in MakeJsonPost:", string(body), " | panic:", r)
		}
	}()
	json.Unmarshal(body, &bundle)

	return bundle["id"].(string), nil
}

func GetShortLink(url string) string {
	resp, err := MakeJsonPost("https://www.googleapis.com/urlshortener/v1/url?key=AIzaSyCsMGYHHdwvQYgOAaskd-GZfGSPe1Tk66w", fmt.Sprintf("{\"longUrl\": \"%s\"}", url))
	if err != nil {
		Error.Println("Get short url error", err)
		return ""
	}
	return resp
}
