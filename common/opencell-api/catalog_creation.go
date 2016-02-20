package opencell_api

import (
	"fmt"
)
const (
	ACCOUNT_URL = "/meveo/AccountWs"
	SETTINGS_URL = "/meveo/SettingsWs"
	CATALOG_URL = "/meveo/CatalogWs"
)

type Catalog struct {
	BasicAuthString string
	ServerUrl string
	
	ioUtils *IOUtils
	httpUtils HttpUtils
}

func NewCatalog(basicAuthStr string, serverUrl string) Catalog {
	return Catalog { 
		BasicAuthString: basicAuthStr, 
		ServerUrl: serverUrl, 
		ioUtils: new (IOUtils), 
		httpUtils: NewHttpUtils(basicAuthStr) }
}

func (this Catalog) doPostToSettings(rawXml []byte) bool {

	url := this.ServerUrl + SETTINGS_URL
	status, err, _ := this.httpUtils.DoPostSoap(url, string(rawXml))
	if err != nil && status == 0 {
		fmt.Println(err)
		return false
	}
	return true
}

func (this Catalog) doPostToCatalog(rawXml []byte) bool {
	
	url := this.ServerUrl + CATALOG_URL 
	status, err, _ := this.httpUtils.DoPostSoap(url, string(rawXml))
	if err != nil && status == 0 {
		fmt.Println(err)
		return false
	}
	return true
}

func (this Catalog) doPostToAccount(rawXml []byte) bool {
	
	url := this.ServerUrl + ACCOUNT_URL
	status, err, _ := this.httpUtils.DoPostSoap(url, string(rawXml))
	if err != nil && status == 0 {
		fmt.Println(err)
		return false
	}
	return true
}

func (this Catalog) CreateInvoiceCategory() bool {
	xmlPath := "xml/init/01_create_invoice_category.xml"
	rawXml := this.ioUtils.GetFileData(this.ioUtils.GetAbsolutePath(xmlPath))     
    if !this.doPostToSettings(rawXml) {
    	fmt.Println("Cannot create invoice category")
    	return false
    }
    return true
}

func (this Catalog) CreateInvoiceSubCategory() bool {
	xmlPath := "xml/init/02_create_invoice_sub_category.xml"  
	rawXml := this.ioUtils.GetFileData(this.ioUtils.GetAbsolutePath(xmlPath))   
    if !this.doPostToSettings(rawXml) {
    	fmt.Println("Cannot create invoice sub category")
    	return false
    }
    return true
}

func (this Catalog) CreateInvoiceSubCategoryCountry() bool {
    xmlPath := "xml/init/03_create_invoice_subcategory_country.xml" 
    rawXml := this.ioUtils.GetFileData(this.ioUtils.GetAbsolutePath(xmlPath))   
    if !this.doPostToSettings(rawXml) {
    	fmt.Println("Cannot create invoice sub category country")
    	return false
    }
    return true
}	

func (this Catalog) CreateCharge() bool {
	xmlPath := "xml/init/04_create_charge.xml"
	rawXml := this.ioUtils.GetFileData(this.ioUtils.GetAbsolutePath(xmlPath))  
    if !this.doPostToCatalog(rawXml) {
    	fmt.Println("Cannot create charge")
    	return false
    }
    return true
}

func (this Catalog) CreateService() bool {
	xmlPath := "xml/init/05_create_service.xml" 
	rawXml := this.ioUtils.GetFileData(this.ioUtils.GetAbsolutePath(xmlPath)) 
    if !this.doPostToCatalog(rawXml) {
    	fmt.Println("Cannot create service")
    	return false
    }
    return true
}

func (this Catalog) CreateOffer() bool {
	xmlPath := "xml/init/06_create_offer.xml" 
	rawXml := this.ioUtils.GetFileData(this.ioUtils.GetAbsolutePath(xmlPath)) 
    if !this.doPostToCatalog(rawXml) {
    	fmt.Println("Cannot create offer")
    	return false
    }
    return true
}

func (this Catalog) CreatePricePlan() bool {
	
    xmlPath := "xml/init/07_create_price_plan.xml" 
    rawXml := this.ioUtils.GetFileData(this.ioUtils.GetAbsolutePath(xmlPath)) 
    if !this.doPostToCatalog(rawXml) {
    	fmt.Println("Cannot create price plan")
    	return false
    }
    return true
}

func (this Catalog) CreateCustomerBrand() bool {   
    
    xmlPath := "xml/init/08_create_customer_brand.xml"
    rawXml := this.ioUtils.GetFileData(this.ioUtils.GetAbsolutePath(xmlPath))  
    if !this.doPostToAccount(rawXml) {
    	fmt.Println("Cannot create customer brand")
    	return false
    }
    return true
}

func (this Catalog) CreateSeller() bool {       
    
    xmlPath := "xml/init/09_create_seller.xml" 
    rawXml := this.ioUtils.GetFileData(this.ioUtils.GetAbsolutePath(xmlPath)) 
    if !this.doPostToSettings(rawXml) {
    	fmt.Println("Cannot create seller")
    	return false
    }
    return true
}

func (this Catalog) CreateBillingCycle() bool {   
    
    xmlPath := "xml/init/10_create_billing_cycle.xml"
    rawXml := this.ioUtils.GetFileData(this.ioUtils.GetAbsolutePath(xmlPath))  
    if !this.doPostToSettings(rawXml) {
    	fmt.Println("Cannot create billing cycle")
    	return false
    }
    return true
}