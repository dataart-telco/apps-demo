package opencell_api

import (
    "io/ioutil"
    "fmt"
    "strconv"
)

const (
	MEDIATION_URL = "/meveo/MediationWs"
)

type Payment struct {
	basicAuthString string
	serverUrl string
}

func NewPayment(basicAuthStr string, srvUrl string) Payment {
	return Payment { 
		basicAuthString: basicAuthStr, 
		serverUrl: srvUrl }
}

func (this Payment) ChargeCustomer(clientID string, time string, minutes float64) bool {
	
	ioUtils := new (IOUtils)
	xmlPath := "xml/charge.xml"
	rawXml, err := ioutil.ReadFile(ioUtils.GetAbsolutePath(xmlPath))
    check(err)
	 
	chargeXml := fmt.Sprintf(string(rawXml), time, strconv.FormatFloat(minutes, 'f', 6, 64), clientID)
	//fmt.Println(chargeXml)
	
	return this.doPostToMediation(chargeXml)
}

func (this Payment) doPostToMediation(xmlData string) bool {
	
	httpUtils := NewHttpUtils(this.basicAuthString)
	url := this.serverUrl + MEDIATION_URL
	status, err, _ := httpUtils.DoPostSoap(url, xmlData)
	if err != nil && status == 0 {
		fmt.Println(err)
		return false
	}
	return true
}