package v4

import (	
	"github.com/mobilemindtec/go-utils/beego/validator"	
	"github.com/mobilemindtec/go-payments/payzen"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/i18n"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"	
	"strings"
	"errors"
	"bytes"	
	"fmt"	
)

const (
	ApiUrl = "https://api.payzen.com.br/api-payment/V4"
)

type Authentication struct {
	UserName string
	PasswordProd string
	PasswordTest string
} 

func NewAuthentication(username string, passwordProd string, passwordTest string) *Authentication{
	return &Authentication{ UserName: username, PasswordProd: passwordProd, PasswordTest: passwordTest }
}

func (this *Authentication) CreateBasicAuth(mode ApiMode) string {

	auth := ""

	switch mode {
		case Prod:
			auth = fmt.Sprintf("%v:%v", this.UserName, this.PasswordProd)
		case Test:
			auth = fmt.Sprintf("%v:%v", this.UserName, this.PasswordTest)
	}

	return fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(auth)))

}

type PayZen struct {
	Mode ApiMode
	Authentication *Authentication

  EntityValidator *validator.EntityValidator
  EntityValidatorResult *validator.EntityValidatorResult

  ValidationErrors map[string]string
  HasValidationError bool
  Lang string 

	Debug bool
}

func NewPayZen(lang string, mode ApiMode, authentication *Authentication) *PayZen{
	entityValidator := validator.NewEntityValidator(lang, "PayZen")
	return &PayZen{ Mode: mode, Authentication: authentication, Lang: lang, EntityValidator: entityValidator }
}

func (this *PayZen) SetDebug() { this.Debug = true }

func (this *PayZen) PaymentCreate(payment *Payment) (*PayZenResult, error) { 

	if !this.onValidPayment(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))				
	}

	result := new(PaymentResponse)
	info, err := this.request(payment, "PCI/Charge/CreatePayment", result)
	result.Response = info["Response"]
	result.Request = info["Request"]	
	return NewPayZenResultWithResponse(result), err
}

func (this *PayZen) PaymentCancelOrRefund(transactionUuid string, amount float64) (*PayZenResult, error) { 
	
  if len(transactionUuid) == 0 {
    this.SetValidationError("transactionUuid", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

  if amount <= 0 {
    this.SetValidationError("amount", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

	result := new(TransationResponse)
	data := make(map[string]interface{})

	data["uuid"] = transactionUuid
	data["amount"] = ConvertAmount(amount)
	data["currency"] = payzen.Currency

	info, err := this.request(data, "Transaction/CancelOrRefund", result)

	result.Response = info["Response"]
	result.Request = info["Request"]	

	return NewPayZenResultWithTransaction(result), err
}

func (this *PayZen) PaymentCapture(transactionUuid string) (*PayZenResult, error) { 
	
  if len(transactionUuid) == 0 {
    this.SetValidationError("transactionUuid", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }


	result := new(TransationResponse)
	data := make(map[string][]string)

	data["uuids"] = []string{transactionUuid,}


	info, err := this.request(data, "Transaction/Capture", result)

	result.Response = info["Response"]
	result.Request = info["Request"]	

	return NewPayZenResultWithTransaction(result), err
}

func (this *PayZen) TokenCreate(payment *Payment) (*PayZenResult, error) { 

	if !this.onValidPaymentWithCreateToken(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))				
	}

	result := new(PaymentResponse)
	info, err := this.request(payment, "PCI/Charge/CreateToken", result)
	result.Response = info["Response"]
	result.Request = info["Request"]	
	return NewPayZenResultWithResponse(result), err
}

func (this *PayZen) TokenUpdate(payment *Payment) (*PayZenResult, error) { 

	if !this.onValidPaymentWithUpdateToken(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))				
	}

	result := new(PaymentResponse)
	info, err := this.request(payment, "Token/Update", result)
	result.Response = info["Response"]
	result.Request = info["Request"]	
	return NewPayZenResultWithResponse(result), err
}

func (this *PayZen) TokenGet(paymentMethodToken string) (*PayZenResult, error) { 

  if len(paymentMethodToken) == 0 {
    this.SetValidationError("paymentMethodToken", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

	result := new(PaymentResponse)
	data := make(map[string]string)
	data["paymentMethodToken"] = paymentMethodToken

	info, err := this.request(data, "Token/Get", result)

	result.Response = info["Response"]
	result.Request = info["Request"]	

	return NewPayZenResultWithResponse(result), err
}

func (this *PayZen) TokenCancel(paymentMethodToken string) (*PayZenResult, error) { 

  if len(paymentMethodToken) == 0 {
    this.SetValidationError("paymentMethodToken", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

	result := new(PaymentResponse)
	data := make(map[string]string)
	data["paymentMethodToken"] = paymentMethodToken

	info, err := this.request(data, "Token/Cancel", result)

	result.Response = info["Response"]
	result.Request = info["Request"]	

	return NewPayZenResultWithResponse(result), err
}

func (this *PayZen) SubscriptionCreate(subscription *Subscription) (*PayZenResult, error) { 

	error := subscription.BuildRule()

	if len(error) > 0 {
		this.SetValidationError("Rrule", error)
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))				
	}

	if !this.onValidSubscription(subscription) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))				
	}	

	result := new(PaymentResponse)
	info, err := this.request(subscription, "Charge/CreateSubscription", result)
	result.Response = info["Response"]
	result.Request = info["Request"]	
	return NewPayZenResultWithResponse(result), err
}

func (this *PayZen) SubscriptionUpdate(subscription *Subscription) (*PayZenResult, error) { 


  if len(subscription.SubscriptionId) == 0 {
    this.SetValidationError("SubscriptionId", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

	error := subscription.BuildRule()

	if len(error) > 0 {
		this.SetValidationError("Rrule", error)
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))				
	}

	if !this.onValidSubscription(subscription) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))				
	}	

	result := new(PaymentResponse)
	info, err := this.request(subscription, "Subscription/Update", result)
	result.Response = info["Response"]
	result.Request = info["Request"]	
	return NewPayZenResultWithResponse(result), err
}

func (this *PayZen) SubscriptionGet(subscriptionId string, paymentMethodToken string) (*PayZenResult, error) { 

  if len(subscriptionId) == 0 {
    this.SetValidationError("subscriptionId", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

  if len(paymentMethodToken) == 0 {
    this.SetValidationError("paymentMethodToken", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

	result := new(PaymentResponse)
	data := make(map[string]string)
	data["subscriptionId"] = subscriptionId
	data["paymentMethodToken"] = paymentMethodToken

	info, err := this.request(data, "Subscription/Get", result)

	result.Response = info["Response"]
	result.Request = info["Request"]

	return NewPayZenResultWithResponse(result), err
}

func (this *PayZen) SubscriptionCancel(subscriptionId string, paymentMethodToken string) (*PayZenResult, error) { 

  if len(subscriptionId) == 0 {
    this.SetValidationError("subscriptionId", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

  if len(paymentMethodToken) == 0 {
    this.SetValidationError("paymentMethodToken", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

	result := new(PaymentResponse)
	data := make(map[string]string)
	data["subscriptionId"] = subscriptionId
	data["paymentMethodToken"] = paymentMethodToken

	info, err := this.request(data, "Subscription/Cancel", result)

	result.Response = info["Response"]
	result.Request = info["Request"]

	return NewPayZenResultWithResponse(result), err
}


func (this *PayZen) TransactionGet(transactionUuid string) (*PayZenResult, error) { 

  if len(transactionUuid) == 0 {
    this.SetValidationError("transactionUuid", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

	result := new(TransationResponse)
	data := make(map[string]string)
	data["uuid"] = transactionUuid

	info, err := this.request(data, "Transaction/Get", result)

	result.Response = info["Response"]
	result.Request = info["Request"]

	return NewPayZenResultWithTransaction(result), err
}



func (this *PayZen) request(requestData interface{}, action string, result interface{}) (map[string]string, error) {

	//result := new(PaymentResponse)
	//result.RequestObject = soap

	info := make(map[string]string)

	jsonContent, err := json.MarshalIndent(requestData, " ", "  ")

	if err != nil {
		fmt.Println("** request: error MarshalIndent: %v", err)
		return info, errors.New(fmt.Sprintf("request - MarshalIndent: %v", err.Error()))
	}


	data := bytes.NewBuffer(jsonContent)
	info["Request"] = string(data.Bytes())

	url := fmt.Sprintf("%v/%v", ApiUrl, action)

	if this.Debug {
		fmt.Println("**********************************************")
		fmt.Println(" ----------------- JSON REQUEST --------------")
		fmt.Println("URL: %v", url)
		fmt.Println("%v", string(data.Bytes()))
		fmt.Println("**********************************************")
	}

	client := new(http.Client)
	req, err := http.NewRequest("POST", url, data)

	if err != nil {
		return info,  errors.New(fmt.Sprintf("http.NewRequest: %v", err.Error()))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", this.Authentication.CreateBasicAuth(this.Mode))

	resp, err := client.Do(req)

	if err != nil {
		return info, errors.New(fmt.Sprintf("client.Do: %v", err.Error()))
	}

	response, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return info, errors.New(fmt.Sprintf("ioutil.ReadAll: %v", err.Error()))		
	}

	fmt.Println("***********************************")
	fmt.Println("***** PAYZEN START RESPONSE ******")	
	fmt.Println("***** STATUS CODE: %v", resp.StatusCode)
	if this.Debug && resp.StatusCode != 200 {
		fmt.Println("***** RESPONSE: %v", string(response))		
	}
	fmt.Println("***** PAYZEN END RESPONSE *********")
	fmt.Println("***********************************")

	info["Response"] = string(response)

	switch resp.StatusCode {
		case 200:

			if this.Debug {

				jsonResult := make(map[string]json.RawMessage)
				json.Unmarshal(response, &jsonResult)
				jsonContent, _ := json.MarshalIndent(jsonResult, " ", "  ")
				fmt.Println("***********************************************")
				fmt.Println(" ----------------- JSON RESPONSE --------------")
				fmt.Println(string(jsonContent))
				fmt.Println("***********************************************")
			}

			err := json.Unmarshal(response, result)

			if err != nil {				
				return info, errors.New(fmt.Sprintf("request - Unmarshal Response: %v", err.Error()))
			}

			return info, nil

		default:
			return info, errors.New(fmt.Sprintf("PayZen API error - Code %v, Status: %v", resp.StatusCode, resp.Status))
	}
}

func (this *PayZen) onValidSubscription(subscription *Subscription) bool {
	
	this.EntityValidatorResult, _ = this.EntityValidator.Valid(subscription, func (validator *validation.Validation) {
		if len(strings.TrimSpace(subscription.PaymentMethodToken)) == 0 {
			validator.SetError(this.getMessage("PaymentMethodToken"), this.getMessage("PayZen.rquired"))
		}  		
	})
  
  if this.EntityValidatorResult.HasError {
  	this.onValidationErrors()
  	return false
  }

  return true	
}

func (this *PayZen) onValidPayment(payment *Payment) bool {
	return this.validPayment(payment, false, false, false)
}

func (this *PayZen) onValidPaymentWithCreateToken(payment *Payment) bool {
	return this.validPayment(payment, false, true, false)
}

func (this *PayZen) onValidPaymentWithUpdateToken(payment *Payment) bool {
	return this.validPayment(payment, false, true, true)
}

func (this *PayZen) validPayment(payment *Payment, checkPaymentId bool, tokenOperation bool, tokenUpdate bool) bool {

	items := []interface{}{
		payment,
		payment.Customer,
		payment.Customer.BillingDetails,
		payment.Device,
	}

  if payment.Card != nil && len(strings.TrimSpace(payment.Card.PaymentMethodToken)) == 0  {
	  items = append(items, payment.Card)
	}

  this.EntityValidatorResult, _ = this.EntityValidator.ValidMult(items, func (validator *validation.Validation) {

  	if checkPaymentId {
			if len(strings.TrimSpace(payment.PaymentOrderId)) == 0 {
				validator.SetError(this.getMessage("PayZen.PaymentOrderId"), this.getMessage("PayZen.rquired"))
			}  		
  	}

  	if !tokenOperation {

			if len(strings.TrimSpace(payment.OrderId)) == 0 {
				validator.SetError(this.getMessage("PayZen.OrderId"), this.getMessage("PayZen.rquired"))
			}

			if payment.Card == nil && !checkPaymentId {
				validator.SetError(this.getMessage("PayZen.Card"), this.getMessage("PayZen.rquired"))
			}
			
			if payment.Card != nil {
				if payment.Card.InstallmentNumber <= 0 {
					validator.SetError(this.getMessage("PayZen.Installments"), this.getMessage("PayZen.rquired"))
				}
			}

			if payment.Amount <= 0 {
				validator.SetError(this.getMessage("PayZen.Amount"), this.getMessage("PayZen.rquired"))
			}

			if payment.Customer == nil && !checkPaymentId {
				validator.SetError(this.getMessage("PayZen.Customer"), this.getMessage("PayZen.rquired"))
			} 

			if payment.Customer != nil {			
				if len(strings.TrimSpace(payment.Customer.BillingDetails.IdentityCode)) == 0 {
					validator.SetError(this.getMessage("PayZen.IdentityCode"), this.getMessage("PayZen.rquired"))
				}
			}

  	}

  	if tokenUpdate {
				if len(strings.TrimSpace(payment.PaymentMethodToken)) == 0 {
					validator.SetError(this.getMessage("PaymentMethodToken"), this.getMessage("PayZen.rquired"))
				}  		
  	}  	

  })

  if this.EntityValidatorResult.HasError {
  	this.onValidationErrors()
  	return false
  }

  return true
}

func (this *PayZen) onValidationErrors(){
	this.HasValidationError = true
	data := make(map[interface{}]interface{})
  this.EntityValidator.CopyErrorsToView(this.EntityValidatorResult, data)
  fmt.Println("DATA ERRORS = %v", data)
  this.ValidationErrors = data["errors"].(map[string]string)
}

func (this *PayZen) getMessage(key string, args ...interface{}) string{
  return i18n.Tr(this.Lang, key, args)
}

func (this *PayZen) SetValidationError(key string, value string){
  this.HasValidationError = true
  if this.ValidationErrors == nil {
    this.ValidationErrors = make(map[string]string)
  }
  this.ValidationErrors[key]= value
}
