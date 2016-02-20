package opencell_api

import (
	"fmt"
)

type Catalog struct {
	BasicAuthString string
	SettingsUrl string
	CatalogUrl string
	AccountUrl string
}

func NewCatalog(basicAuthStr string, settingsUrl string, catalogUrl string, accountUrl string) Catalog {
	return Catalog { BasicAuthString: basicAuthStr, SettingsUrl: settingsUrl, CatalogUrl: catalogUrl, AccountUrl: accountUrl }
}
/*
func (this Catalog) CreateInvoice(xmlPath) bool {
	
	rawXml, err := ioutil.ReadFile(xmlPath)
    check(err)
	
	if !catalog.DoPostToSettings(rawXml) {
		fmt.Println("cannot create invoice category")
		return false
	}
}*/

func (this Catalog) DoPostToSettings(xmlPath string) bool {

	ioUtils := new(IOUtils)
	rawXml := ioUtils.GetFileData(xmlPath)  
	httpUtils := NewHttpUtils(this.BasicAuthString)
	status, err, _ := httpUtils.DoPostSoap(this.SettingsUrl, string(rawXml))
	if err != nil && status == 0 {
		fmt.Println(err)
		return false
	}
	return true
}

func (this Catalog) DoPostToCatalog(xmlPath string) bool {
	ioUtils := new(IOUtils)
	rawXml := ioUtils.GetFileData(xmlPath)  
	httpUtils := NewHttpUtils(this.BasicAuthString)
	status, err, _ := httpUtils.DoPostSoap(this.CatalogUrl, string(rawXml))
	if err != nil && status == 0 {
		fmt.Println(err)
		return false
	}
	return true
}

func (this Catalog) DoPostToAccount(xmlPath string) bool {
	ioUtils := new(IOUtils)
	rawXml := ioUtils.GetFileData(xmlPath)  
	httpUtils := NewHttpUtils(this.BasicAuthString)
	status, err, _ := httpUtils.DoPostSoap(this.AccountUrl, string(rawXml))
	if err != nil && status == 0 {
		fmt.Println(err)
		return false
	}
	return true
}