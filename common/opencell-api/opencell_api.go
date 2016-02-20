package opencell_api

import (
	"fmt"
	"time"
)

type OpencellAPI struct {
	catalog Catalog
	customer Customer
	wallet Wallet
	payment Payment
}

const TIME_FORMAT = "2006-01-02T15:04:05.000Z"

func NewOpencellAPI(basicAuthStr string, srvUrl string) OpencellAPI {
	srvUrl = "http://" + srvUrl
	return OpencellAPI {
		catalog: NewCatalog(basicAuthStr, srvUrl),
		customer: NewCustomer(basicAuthStr, srvUrl),
		wallet: NewWallet(basicAuthStr, srvUrl),
		payment: NewPayment(basicAuthStr, srvUrl)}
}

func (this OpencellAPI) InitOpenCell() bool {
	
	if (!this.catalog.CreateInvoiceCategory()) { 
		return false
	}
	
	if (!this.catalog.CreateInvoiceSubCategory()) { 
		return false
	}
	
	if (!this.catalog.CreateInvoiceSubCategoryCountry()) { 
		return false
	}
	
	if (!this.catalog.CreateCharge()) { 
		return false
	}
	
	if (!this.catalog.CreateService()) { 
		return false
	}	
	
	if (!this.catalog.CreateOffer()) { 
		return false
	}
	
	if (!this.catalog.CreatePricePlan()) { 
		return false
	}
	
	if (!this.catalog.CreateCustomerBrand()) { 
		return false
	}
	
	if (!this.catalog.CreateSeller()) { 
		return false
	}		
	
	if (!this.catalog.CreateBillingCycle()) { 
		return false
	}		
	return true
}

func (this OpencellAPI) CreateNewCustomer(clientID string) bool {
	
	if this.customer.CreateCustomerHierarchy(clientID) {
		fmt.Printf("OK : customer with id %s created \n", clientID)
		return true
	}
	fmt.Printf("Error: cannot create customer \n")
	return false
}

func (this OpencellAPI) ChargeCustomer(clientID string, time string, minutes float64) bool {
	if this.payment.ChargeCustomer(clientID, time, minutes) {
		fmt.Printf("OK : customer with id %s was charged with minutes = %d \n", clientID, minutes)
		return true
	}
	fmt.Printf("Error: cannot charge a customer \n")
	return false
}

func (this OpencellAPI) GetBalance(clientID string) float64 {
	result, balance := this.wallet.GetOpenBalance(clientID)
	if result {
		fmt.Printf("OK : customer balance is %f \n", balance)
		return balance
	}
	return -1.0
}

func (this OpencellAPI) GetBalanceWithRange(clientID string, from time.Time, to time.Time) float64 {
	result, balance := this.wallet.GetOpenBalanceTimeRange(clientID, from, to)
	if result {
		fmt.Printf("OK : customer balance is %f \n", balance)
		return balance
	}
	return -1.0
}