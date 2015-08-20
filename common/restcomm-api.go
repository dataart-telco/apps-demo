package common
import (
	"net/http"
	"fmt"
	"io/ioutil"
	"net/url"
	"bytes"
	"strconv"
	"errors"
	"encoding/json"
	"strings"
)

type CallInfo struct{
	Sid string  `json:"sid"`    
    Uri string  `json:"uri"`    
}

type IncomingPhoneNumber struct{
	Sid string			`json:"sid"`
	PhoneNumber string	`json:"phone_number"`
}

type RestcommApi struct{

	Server string
	User string
	Pass string
}

func NewRestcommApi(server string, user string, pass string)(RestcommApi){
	return RestcommApi{Server:server, User:user, Pass:pass}
}

func NewIncomingPhoneNumber(sid string, phoneNumber string)(IncomingPhoneNumber) {
	return IncomingPhoneNumber{Sid:sid, PhoneNumber:phoneNumber}
}

func (n IncomingPhoneNumber)Find(api RestcommApi)(*IncomingPhoneNumber){
	acc := api.User + ":" + api.Pass
	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/IncomingPhoneNumbers.json", acc, api.Server, api.User)

	resp, err := http.Get(path)

	if(err != nil){
		panic(err)
	}

	if(resp.StatusCode != 200){
		panic(errors.New(fmt.Sprintf("Can't execute request %d", resp.StatusCode)))
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	numbers := make([]IncomingPhoneNumber,0)
	json.Unmarshal(body, &numbers)

	fmt.Println("!numbers.len", len(numbers))
	for _, e := range numbers{
		fmt.Println("\t", e.Sid, e.PhoneNumber)
		if e.PhoneNumber == n.PhoneNumber {
			return &e
		}
	}
	return nil
}

func (n IncomingPhoneNumber)Update(api RestcommApi, callBack string)(error){
	acc := api.User + ":" + api.Pass

	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/IncomingPhoneNumbers/%s.json", acc, api.Server, api.User, n.Sid)
	data := url.Values{
		"isSIP"		: {"true"},
		"VoiceUrl"	: {callBack},
		"SmsUrl"	: {callBack}}

	return api.Post(path, data)
}

func (n IncomingPhoneNumber)Create(api RestcommApi, callBack string)(error){
	acc := api.User + ":" + api.Pass

	path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/IncomingPhoneNumbers.json", acc, api.Server, api.User)
	data := url.Values{
		"isSIP"			: {"true"},
		"VoiceUrl"		: {callBack},
		"SmsUrl"		: {callBack},
		"PhoneNumber"	: {n.PhoneNumber}}

	return api.Post(path, data)
}

func (n IncomingPhoneNumber)CreateOrUpdate(api RestcommApi, callBack string)(error){
	e := n.Find(api)
	if(e != nil){
		fmt.Println("Number was found", n)
		return e.Update(api, callBack)
	}else{
		return n.Create(api, callBack)
	}
}

func (*RestcommApi) Post(path string, params url.Values)(error){
	data := params.Encode();

	client := &http.Client{}
	r, _ := http.NewRequest(
		"POST",
		path,
		bytes.NewBufferString(data))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(r)
	if(err != nil){
		return err
	}

	if(resp.StatusCode != 200){
		return errors.New("Resp code is not 200 for " + path)
	}
	return nil
}

func (api RestcommApi) CompleteCallByUri(callUri string)(bool){
	return api.UpdateCallByUri(callUri, url.Values{"Status" : {"completed"}})
}

func (api RestcommApi) UpdateCallByUri(callUri string, params url.Values)(bool){
	acc := api.User + ":" + api.Pass
	path := fmt.Sprintf("http://%s@%s/restcomm%s", acc, api.Server, callUri)
	err := api.Post(path, params)
	if(err != nil){
		return false
	}
	return true
}

func (api RestcommApi)MakeCall(from string, to string, callback string)(*CallInfo, error){
    fmt.Println("\tapi.MakeCall: from =", from, " to =", to, " callback =", callback)
	acc := api.User + ":" + api.Pass
    path := fmt.Sprintf("http://%s@%s/restcomm/2012-04-24/Accounts/%s/Calls.json", acc, api.Server, api.User)
	resp,err := http.PostForm(path,
		url.Values{
			"From" : {from},
			"To"   : {to},
			"Url"  : {callback}})

	if(err != nil){
		return nil, err
	}

    if(resp.StatusCode != 200){
		return nil, errors.New("Resp code is not 200 for " + path)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

    var call CallInfo
	json.Unmarshal(body, &call)

	return &call, nil;
}

func GetClientName(client string)(string){
	if(strings.HasPrefix(client, "sip:")){
		return client[4 : len(client)];
	}
	if(strings.HasPrefix(client, "client:")){
		return client[7 : len(client)];
	}
	return client
}

func ConvertToRestcommNumber(client string)(string){
	if(strings.HasPrefix(client, "sip:")){
		return "client:" + client[4 : len(client)];
	}
	return client
}

func ConvertToSip(from string)(string){
	if(!strings.HasPrefix(from, "+")){
		return "sip:" + from
	}
	return from;
}
