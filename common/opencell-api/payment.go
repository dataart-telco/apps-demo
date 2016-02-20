package opencell_api

import (
    "io/ioutil"
    "fmt"
    "strconv"
)

type Payment struct {
	BasicAuthString string
	MediationUrl string
}

func NewPayment(basicAuthStr string, mediationUrl string) Payment {
	return Payment{ BasicAuthString: basicAuthStr, MediationUrl: mediationUrl }
}

func (this Payment) ChargeCustomer(xmlPath string, clientID string, time string, minutes int) bool {
	
	httpUtils := NewHttpUtils(this.BasicAuthString)

	rawXml, err := ioutil.ReadFile(xmlPath)
    check(err)
	 
	chargeXml := fmt.Sprintf(string(rawXml), time, strconv.Itoa(minutes), clientID)
	fmt.Println(chargeXml);
	
	status, err, _ := httpUtils.DoPostSoap(this.MediationUrl, chargeXml)
	if err != nil && status == 0 {
		fmt.Println(err)
		return false
	}
	return true
}