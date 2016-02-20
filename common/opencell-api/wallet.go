package opencell_api

import (
    "io/ioutil"
    "fmt"
    "encoding/json"
    "strconv"
)

const (
	WALLET_URL = "/meveo/api/rest/billing/wallet/balance/open"
)

type Wallet struct {
	basicAuthString string
	serverUrl string
}

func NewWallet(basicAuthStr string, srvUrl string) Wallet {
	return Wallet { 
		basicAuthString: basicAuthStr, 
		serverUrl: srvUrl }
}

type ResponseBalance struct {
	Status string     `json:"status"`
	ErrorCode int     `json:"errorCode"`
	Message string    `json:"message"`
} 

func (this Wallet) GetOpenBalance(clientID string) (bool, float64) {
	
	httpUtils := NewHttpUtils(this.basicAuthString)
	ioUtils := new (IOUtils)
	
	xmlPath := "xml/open_balance.xml"
	rawJSON := ioUtils.GetFileData(ioUtils.GetAbsolutePath(xmlPath))

	 
	openBalanceJSON := fmt.Sprintf(string(rawJSON), clientID)
	//fmt.Println(openBalanceJSON);
	
	url := this.serverUrl + WALLET_URL
	status, err, resp := httpUtils.DoPostJson(url, openBalanceJSON)
	if err != nil && status == 0 {
		fmt.Println("Payemnt.ChargeClient: post error - ", err)
		return false, -1
	}
	
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return false, -1
    }
	
	balance := this.parseResponse(respBody)
	
	return true, balance
}

func (this Wallet) parseResponse(respBody []byte) float64 {
	
	//fmt.Println("Response: \n")
	//fmt.Println(string(respBody))

	var respBalance ResponseBalance
	err := json.Unmarshal(respBody, &respBalance)
	if err != nil {
		fmt.Printf(err.Error())
		return -1.0;
	}
	
	//fmt.Println(respBalance)

	if  respBalance.ErrorCode != 0 {
		fmt.Printf("Cannot get custoner balance. Error code = %i", respBalance.ErrorCode)
		return -1.0
	}
	
	balanceFloat,_ := strconv.ParseFloat(respBalance.Message, 64)
	
	return balanceFloat
}