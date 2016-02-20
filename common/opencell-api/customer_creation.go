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

type Customer struct {
	BasicAuthString string
	AccountUrl string
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func NewCustomer(basicAuthStr string, accountUrl string) Customer {
	return Customer{ BasicAuthString: basicAuthStr, AccountUrl: accountUrl }
}

func (this Customer) CreateCustomerHierarchy(xmlPath string, clientID string, firstName string, lastName string) bool {
	
	rawXml, err := ioutil.ReadFile(xmlPath)
    check(err)
	 
	customerCreationXml := fmt.Sprintf(string(rawXml), clientID, firstName, lastName, clientID, clientID, clientID, clientID, clientID, clientID, clientID, clientID, clientID, clientID)
	fmt.Println(customerCreationXml);
	
	return this.doPostToAccount(customerCreationXml)
	
}

func (this Customer) doPostToAccount(xmlData string) bool {
	
	httpUtils := NewHttpUtils(this.BasicAuthString)
	status, err, _ := httpUtils.DoPostSoap(this.AccountUrl, xmlData)
	if err != nil && status == 0 {
		fmt.Println(err)
		return false
	}
	return true
}
