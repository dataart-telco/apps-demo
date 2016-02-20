package opencell_api

import (
	"bytes"
	"errors"
	"net/http"
	"fmt"
	"io/ioutil"
)

type HttpUtils struct {
	BasicAuthString string
}

func NewHttpUtils(basicAuthStr string) HttpUtils {
	return HttpUtils{BasicAuthString: basicAuthStr}
}

func (this HttpUtils) DoPostSoap(url string, xml string) (int, error, http.Response) {
	
	request, err := http.NewRequest("POST", url, bytes.NewBufferString(xml))
	
	if err != nil {
		panic(err)
	}
	
	request.Header.Add("Authorization", this.BasicAuthString)

	client := &http.Client{}
	resp, err := client.Do(request)
	
	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	
	if err != nil {
		fmt.Println(err)
		return resp.StatusCode, err, *resp
	}

	if resp.StatusCode != 200 {
		return resp.StatusCode, errors.New("Resp code is not 200 for " + url + "; status = " + resp.Status), *resp
	}
	
	return 0, err, *resp
}

func (this HttpUtils) DoPostJson(url string, jsonData string) (int, error, http.Response) {
	
	request, err := http.NewRequest("POST", url, bytes.NewBufferString(jsonData))
	
	if err != nil {
		panic(err)
	}
	
	request.Header.Add("Authorization", this.BasicAuthString)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	
	if err != nil {
		return resp.StatusCode, err, *resp
	}

	if resp.StatusCode != 200 {
		return resp.StatusCode, errors.New("Resp code is not 200 for " + url + "; status = " + resp.Status), *resp
	}
	
	return 0, err, *resp
}