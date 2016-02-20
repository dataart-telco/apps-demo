package opencell_api

import (
	//"net/http"
	//"encoding/xml"
	"fmt"
	//"io"
    "io/ioutil"
    //"bytes"
    //"errors"
)

const (
	//ACCOUNT_URL = "/meveo/AccountWs"
)

type Customer struct {
	BasicAuthString string
	ServerUrl string
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func NewCustomer(basicAuthStr string, serverUrl string) Customer {
	return Customer{ BasicAuthString: basicAuthStr, ServerUrl: serverUrl }
}

func (this Customer) CreateCustomerHierarchy(clientID string) bool {
	
	ioUtils := new (IOUtils)
	
	xmlPath := "xml/customer_hierarchy.xml"
	rawXml, err := ioutil.ReadFile(ioUtils.GetAbsolutePath(xmlPath))
    check(err)
	 
	customerCreationXml := fmt.Sprintf(string(rawXml), clientID, clientID, clientID, clientID, clientID, clientID, clientID, clientID, clientID, clientID, clientID)
	fmt.Println(customerCreationXml);
	
	return this.doPostToAccount(customerCreationXml)
}

func (this Customer) doPostToAccount(xmlData string) bool {
	
	url := this.ServerUrl + ACCOUNT_URL
	httpUtils := NewHttpUtils(this.BasicAuthString)
	status, err, _ := httpUtils.DoPostSoap(url, xmlData)
	if err != nil && status == 0 {
		fmt.Println(err)
		return false
	}
	return true
}
