package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)
const (
	CallStatusQeued = "queued"
	CallStatusRinging = "ringing"
	CallStatusInProgress = "in-progress"
	CallStatusCompleted = "completed"
	CallStatusFailed = "failed"
	CallStatusBusy = "busy"
	CallStatusNoAnswer = "no-answer"
)
type CallInfo struct {
	Sid string `json:"sid"`
	Status string `json:"status"`
	Duration float64 `json:"duration"`
	Uri string `json:"uri"`
}

type IncomingPhoneNumber struct {
	Sid         string `json:"sid"`
	PhoneNumber string `json:"phone_number"`
}

type RestcommApi struct {
	Server string
	User   string
	Pass   string
}

func NewRestcommApi(server string, user string, pass string) RestcommApi {
	return RestcommApi{Server: server, User: user, Pass: pass}
}

func NewIncomingPhoneNumber(sid string, phoneNumber string) IncomingPhoneNumber {
	return IncomingPhoneNumber{Sid: sid, PhoneNumber: phoneNumber}
}

func (n IncomingPhoneNumber) Find(api RestcommApi) *IncomingPhoneNumber {
	acc := api.User + ":" + api.Pass
	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/IncomingPhoneNumbers.json", acc, api.Server, api.User)

	resp, err := http.Get(path)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(errors.New(fmt.Sprintf("Can't execute request %d", resp.StatusCode)))
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	numbers := make([]IncomingPhoneNumber, 0)
	json.Unmarshal(body, &numbers)

	fmt.Println("!numbers.len", len(numbers))
	for _, e := range numbers {
		fmt.Println("\t", e.Sid, e.PhoneNumber)
		if e.PhoneNumber == n.PhoneNumber {
			return &e
		}
	}
	return nil
}
func (n IncomingPhoneNumber) Update(api RestcommApi, callBack string) error {
	acc := api.User + ":" + api.Pass

	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/IncomingPhoneNumbers/%s.json", acc, api.Server, api.User, n.Sid)
	data := url.Values{
		"isSIP":    {"true"},
		"VoiceUrl": {callBack},
		"SmsUrl":   {callBack}}

	_, err := api.Post(path, data)
	return err
}

func (n IncomingPhoneNumber) Create(api RestcommApi, callBack string) error {
	acc := api.User + ":" + api.Pass

	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/IncomingPhoneNumbers.json", acc, api.Server, api.User)
	data := url.Values{
		"isSIP":       {"true"},
		"VoiceUrl":    {callBack},
		"SmsUrl":      {callBack},
		"PhoneNumber": {n.PhoneNumber}}

	_, err := api.Post(path, data)
	return err
}

func (n IncomingPhoneNumber) CreateOrUpdate(api RestcommApi, callBack string) error {
	e := n.Find(api)
	if e != nil {
		fmt.Println("Number was found", n)
		return e.Update(api, callBack)
	} else {
		return n.Create(api, callBack)
	}
}

func (*RestcommApi) Post(path string, params url.Values) (int, error) {
	Trace.Println("Make POST req: url =", path)

	data := params.Encode()

	client := &http.Client{}
	r, _ := http.NewRequest(
		"POST",
		path,
		bytes.NewBufferString(data))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(r)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return resp.StatusCode, errors.New("Resp code is not 200 for " + path + "; status = " + resp.Status)
	}
	return 200, nil
}

func (api RestcommApi) CompleteCallByUri(callUri string) bool {
	return api.UpdateCallByUri(callUri, url.Values{"Status": {"completed"}})
}

func (api RestcommApi) UpdateCallByUri(callUri string, params url.Values) bool {
	if api.Server == "" {
		Error.Println("RestcommApi.UpdateCallByUri: api.Server is empty")
		return false
	}
	acc := api.User + ":" + api.Pass
	path := fmt.Sprintf("http://%s@%s/restcomm%s", acc, api.Server, callUri)
	status, err := api.Post(path, params)

	if err != nil && status == 0 {
		Error.Println("RestcommApi.UpdateCallByUri: post error - ", err)
		return false
	}
	return true
}

func (api RestcommApi) MakeCall(from string, to string, callback string, statusCallback string) (*CallInfo, error) {
	Trace.Println("\tapi.MakeCall: from =", from, " to =", to, " callback =", callback)
	acc := api.User + ":" + api.Pass
	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/Calls.json", acc, api.Server, api.User)
	resp, err := http.PostForm(path,
		url.Values{
			"From":           {from},
			"To":             {to},
			"Url":            {callback},
			"StatusCallback": {statusCallback}})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Resp code is not 200 for " + path + "; StatusCode = " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var call CallInfo
	json.Unmarshal(body, &call)

	return &call, nil
}

func (api RestcommApi) GetCallInfo(sid string) (*CallInfo, error) {

	acc := api.User + ":" + api.Pass
	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/Calls/%s.json", acc, api.Server, api.User, sid)
	resp, err := http.Get(path)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Resp code is not 200 for " + path + "; StatusCode = " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	
	var call CallInfo
	json.Unmarshal(body, &call)

	return &call, nil
}

func (api RestcommApi) GetTo(sid string) (string, error) {
	acc := api.User + ":" + api.Pass
	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/Calls/%s.json", acc, api.Server, api.User, sid)
	resp, err := http.Get(path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return "", err2
	}

	call := make(map[string]interface{})
	json.Unmarshal(body, &call)

	return call["to"].(string), nil
}

// will overrid StatusCallback and Url fields of call, please use it for inprogress calls only
func (api RestcommApi) SubscribeStatus(sid string, statusCallback string) error {

	acc := api.User + ":" + api.Pass
	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/Calls/%s", acc, api.Server, api.User, sid)
	params := url.Values{"StatusCallback": {statusCallback}, "Url": {statusCallback}}

	_, errPost := api.Post(path, params)
	if errPost != nil {
		return errPost
	}
	return nil
}

func (api RestcommApi) SendSms(to string, from string, msg string) error {
	Trace.Println("TODO: send sms")

	acc := api.User + ":" + api.Pass
	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/SMS/Messages", acc, api.Server, api.User)
	params := url.Values{"To": {to}, "From": {from}, "Body": {msg}}

	_, err := api.Post(path, params)
	if err != nil {
		Error.Println("Send message error", err)
	}
	return err
}

func GetClientName(client string) string {
	if strings.HasPrefix(client, "sip:") {
		return client[4:len(client)]
	}
	if strings.HasPrefix(client, "client:") {
		return client[7:len(client)]
	}
	return client
}

var digistRegExp = regexp.MustCompile("[0-9]+")

func IsPhoneNumber(from string) bool {
	return digistRegExp.MatchString(from)
}

func ConvertToSipCall(from string, didProvider string) string {
	if digistRegExp.MatchString(from) {
		return "sip:" + from + "@" + didProvider
	}
	return "client:" + from
}

func ConvertToSipSms(from string, didProvider string) string {
	//for sms we will have settings in config
	return from

	/*    if digistRegExp.MatchString(from) {
		return "sip:" + from + "@" + didProvider
	}
	// no need client: prefix for sms
	return from*/
}
