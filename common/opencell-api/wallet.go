package opencell_api

import (
    "io/ioutil"
    "fmt"
    "encoding/json"
    "strconv"
)

type Wallet struct {
	BasicAuthString string
	WalletUrl string
}

func NewWallet(basicAuthStr string, walletUrl string) Wallet {
	return Wallet{ BasicAuthString: basicAuthStr, WalletUrl: walletUrl }
}

type ResponseBalance struct {
	Status string     `json:"status"`
	ErrorCode int     `json:"errorCode"`
	Message string    `json:"message"`
} 

func (this Wallet) GetOpenBalance(xmlPath string, clientID string) (bool, float64) {
	
	httpUtils := NewHttpUtils(this.BasicAuthString)
	
	rawJSON, err := ioutil.ReadFile(xmlPath)
    check(err)
	 
	openBalanceJSON := fmt.Sprintf(string(rawJSON), clientID)
	fmt.Println(openBalanceJSON);
	
	status, err, resp := httpUtils.DoPostJson(this.WalletUrl, openBalanceJSON)
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
	
	fmt.Println("Response: \n")
	fmt.Println(string(respBody))

	var respBalance ResponseBalance
	err := json.Unmarshal(respBody, &respBalance)
	if err != nil {
		fmt.Printf(err.Error())
		return -1.0;
	}
	
	fmt.Println(respBalance)

	if  respBalance.ErrorCode != 0 {
		fmt.Printf("Cannot get custoner balance. Error code = %i", respBalance.ErrorCode)
		return -1.0
	}
	
	balanceFloat,_ := strconv.ParseFloat(respBalance.Message, 64)
	
	return balanceFloat
}