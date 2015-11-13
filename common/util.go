package common

import (
	"bytes"
	"encoding/json"
	"errors"
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
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Make json request respose is %d", resp.StatusCode))
	}
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

func GetGoogleShortLink(url string) string {
	resp, err := MakeJsonPost("https://www.googleapis.com/urlshortener/v1/url?key=AIzaSyCsMGYHHdwvQYgOAaskd-GZfGSPe1Tk66w", fmt.Sprintf("{\"longUrl\": \"%s\"}", url))
	if err != nil {
		Error.Println("Get short url error", err)
		return ""
	}
	return resp
}

func GetShortLink(url string) string {
	Trace.Println("make short link:", url)
	resp, err := http.Get("https://api-ssl.bitly.com/v3/shorten?access_token=1beca81b7f09ddd49d541fa802042cd19c468529&longUrl=" + url)
	if err != nil {
		Error.Println("Make shortlink error", err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		Error.Println("Make shortlink respose is", resp.StatusCode)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error.Println("Mae shortlink error - read body", err)
		return ""
	}

	bundle := make(map[string]interface{})
	defer func() {
		r := recover()
		if r != nil {
			Error.Println("Catch panic in GetShortLink:", string(body), " | panic:", r)
		}
	}()
	json.Unmarshal(body, &bundle)
	if bundle["status_code"].(float64) != 200 {
		return ""
	}
	shortLink := bundle["data"].(map[string]interface{})["url"].(string)
	return shortLink
}
