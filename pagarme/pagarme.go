
package pagarme

import (
	"github.com/mobilemindtec/go-utils/beego/validator"
  "github.com/mobilemindtec/go-utils/beego/db"
  "github.com/astaxie/beego/validation"
  "github.com/leekchan/accounting"
  "github.com/beego/i18n"
	"encoding/json"
  "encoding/hex"
  "crypto/sha1"
  "crypto/hmac"	
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
  "errors"
	"bytes"
	"fmt"
)

const (

	PAGARME_TRANSACTION_URL = "https://api.pagar.me/1/transactions"
	PAGARME_TRANSACTION_CAPTURE_URL = "https://api.pagar.me/1/transactions/%v/capture"
	PAYMENT_METHOD_BOLETO = "boleto"
	PAYMENT_METHOD_CREDIT_CARD = "credit_card"

)

type PagarmeResult struct {

	Errors []map[string]string `json:"errors"`

	ResponseValues map[string]interface{}
	Response string

}

type PagarmeAddress struct {
	
	Neighborhood string ` json:"neighborhood" valid:"Required" `
	Street string `json:"street" valid:"Required" `
	StreetNumber string `json:"street_number" valid:"Required" `
	ZipCode string `json:"zipcode" valid:"Required" `

	City string `json:"city" valid:"Required" `
	State string `json:"state" valid:"Required" `

}

type PagarmePhone struct {

	Ddd string `json:"ddd" valid:"Required;MaxSize(2)" `
	Number string `json:"number" valid:"Required;MaxSize(9);MinSize(9)" `

}

type PagarmeCustomer struct {

	DocumentNumber string `json:"document_number" valid:"Required;MaxSize(11);MinSize(11)"`
	Email string `json:"email" valid:"Required;Email"`
	Name string `json:"name" valid:"Required"`

	Address *PagarmeAddress `json:"address" valid:"Required"`
	Phone *PagarmePhone `json:"phone" valid:"Required"`

}

type PagarmePayment struct {

	Amount int `json:"amount" valid:"Required"`
	ApiKey string `json:"api_key" valid:"Required"`	
 	Installments int `json:"Installments" valid:"Required"`	// parcelas
	Customer *PagarmeCustomer `json:"customer" valid:"Required"`
	PaymentMethod string `json:"payment_method" valid:"Required"`

	CardId string `json:"card_id" valid:""`
	
	CardHash string `json:"card_hash" valid:""`

	CardHolderName string `json:"card_holder_name" valid:""`
	CardExpirationDate string `json:"card_expiration_date" valid:""`
	CardNumber	string `json:"card_number" valid:""`
	CardCvv string `json:"card_cvv" valid:""`

	PostbackUrl string `json:"postback_url" valid:"Required"`
	SoftDescriptor string `json:"soft_descriptor" valid:"Required"`
	Metadata map[string]string `json:"metadata"`

	BoletoExpirationDate string `json:"boleto_expiration_date" valid:""`
	BoletoInstructions string `json:"boleto_instructions" valid:""`

}


func NewPagarmePaymentCard(apiKey string, cardHash string, amount float64) *PagarmePayment {
	return &PagarmePayment{ ApiKey: apiKey, CardHash: cardHash, Amount: formatAmount(amount), Installments: 1, PaymentMethod: PAYMENT_METHOD_CREDIT_CARD }
}


func NewPagarmePaymentBoleto(apiKey string, amount float64) *PagarmePayment {
	return &PagarmePayment{ PaymentMethod: PAYMENT_METHOD_BOLETO, ApiKey: apiKey, Amount: formatAmount(amount)  }
}

type PagarmeService struct {
  Lang string
  Session *db.Session
  ApiKey string
  EntityValidator *validator.EntityValidator
  EntityValidatorResult *validator.EntityValidatorResult
}

type CaptureData struct {
	ApiKey string `json:"api_key" valid:"Required"`	
	IdOuToken string `json:"-" valid:"Required"` 
	Amount int `json:"amount" valid:"Required"` 
	SplitRules []string `json:"split_rules"` 
	Metadata map[string]string `json:"metadata"`
}

func NewCaptureData(apiKey string, idOuToken string, amount float64) *CaptureData {	
	return &CaptureData{ ApiKey: apiKey, IdOuToken: idOuToken, Amount:  formatAmount(amount) }
}


func NewPagarmeService(session *db.Session, lang string) *PagarmeService{
	entityValidator := validator.NewEntityValidator(lang, "Pagarme")
  return &PagarmeService{ Lang: lang, Session: session, EntityValidator: entityValidator }
}


func (this *PagarmeService) Pay(customer *PagarmePayment) (map[string]string, error) {
	
	if !this.onValid(customer) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}

	jsonData, err := json.Marshal(customer)

	if err != nil {
		fmt.Println("error json.Marshal ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	data := bytes.NewBuffer(jsonData)

	r, err := http.Post(PAGARME_TRANSACTION_URL, "text/json", data)

	if err != nil {
		fmt.Println("error http.Post ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	response, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("error ioutil.ReadAll ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	fmt.Println(string(response))	

	return nil, nil
}

func (this *PagarmeService) Capture(captureData *CaptureData) (*PagarmeResult, error) {
	
	if !this.onValidCapture(captureData) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}

	jsonData, err := json.Marshal(captureData)

	if err != nil {
		fmt.Println("error json.Marshal ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	fmt.Println("** jsonData = %v", string(jsonData))

	data := bytes.NewBuffer(jsonData)

	url := fmt.Sprintf(PAGARME_TRANSACTION_CAPTURE_URL, captureData.IdOuToken) 

	r, err := http.Post(url, "application/json", data)

	if err != nil {
		fmt.Println("error http.Post ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	response, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("error ioutil.ReadAll ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	fmt.Println("***** PAGARME CAPTURE START RESPONSE ****** ")	
	fmt.Println("***** STATUS CODE: %v", r.StatusCode)	
	fmt.Println("***** RESPONSE: %v", string(response))	
	fmt.Println("***** PAGARME CAPTURE END RESPONSE ****** ")

	switch r.StatusCode {
		case 200:

			values := make(map[string]interface{})

			if err := json.Unmarshal(response, &values); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}

			return &PagarmeResult{ Response: string(response), ResponseValues: values  }, nil

		case 400:
			
			result := new(PagarmeResult)
			
			if err := json.Unmarshal(response, result); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}
		
			this.EntityValidatorResult = new(validator.EntityValidatorResult)
			this.EntityValidatorResult.Errors = result.Errors[0]
			this.EntityValidatorResult.HasError = true

			return nil, errors.New("Pagarme: Erro de validação")

		case 401:
			return nil, errors.New("Pagarme: Access Denied")
		case 404:
			return nil, errors.New("Pagarme: Not Found")
		case 500:
			return nil, errors.New("Pagarme: Unknow error")
		default:
			return nil, errors.New(fmt.Sprintf("Pagarme: API error - %v", r.StatusCode))
	}

}

func (this *PagarmeService) onValid(customer *PagarmePayment) bool {
  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(customer, func (validator *validation.Validation) {
  	
  	emptyCardHash := len(strings.TrimSpace(customer.CardHash)) == 0
  	emptyCardId := len(strings.TrimSpace(customer.CardId)) == 0

		if emptyCardHash && emptyCardId {
			if len(strings.TrimSpace(customer.CardHolderName)) == 0 {
				validator.SetError("CardHolderName", this.getMessage("Pagarme.rquired"))
			} 

			if len(strings.TrimSpace(customer.CardExpirationDate)) == 0 {
				validator.SetError("CardExpirationDate", this.getMessage("Pagarme.rquired"))
			} 

			if len(strings.TrimSpace(customer.CardNumber)) == 0 {
				validator.SetError("CardNumber", this.getMessage("Pagarme.rquired"))
			} 

			if len(strings.TrimSpace(customer.CardCvv)) == 0 {
				validator.SetError("CardCvv", this.getMessage("Pagarme.rquired"))
			} 
		}		

  })
  return this.EntityValidatorResult.HasError == false
}

func (this *PagarmeService) onValidCapture(captureData *CaptureData) bool {
  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(captureData, nil)
  return this.EntityValidatorResult.HasError == false
}

func formatAmount(amount float64) int {
	ac := accounting.Accounting{Symbol: "", Precision: 2, Thousand: "", Decimal: ""}
	text := strings.Replace(ac.FormatMoney(amount), ",", "", -1)
	text = strings.Replace(text, ".", "", -1)
	val, _ := strconv.Atoi(text)
	return val
}

func (this *PagarmeService) getMessage(key string, args ...interface{}) string{
  return i18n.Tr(this.Lang, key, args)
}

func CheckPostbackSignature(apiKey string, hubSignature string, requestBody []byte) bool {

  pagarmeSignature := hubSignature

  if !strings.Contains(pagarmeSignature, "="){
    fmt.Println("************************************************")
    fmt.Println("** Pagarme Signature not has =")
    fmt.Println("************************************************")        
    return false
  }

  fmt.Println("************************************************")
  fmt.Println("**  X-Hub-Signature = ", pagarmeSignature)
  fmt.Println("************************************************")    

  /*
  this.Log("************************************************")
  this.Log("**  RequestBody = ", string(this.Ctx.Input.RequestBody))
  this.Log("************************************************")    
  */

  cleanedSignature := strings.Split(pagarmeSignature, "=")[1]

  fmt.Println("************************************************")
  fmt.Println("**  cleanedSignature = ", cleanedSignature)
  fmt.Println("************************************************")    

  mac := hmac.New(sha1.New, []byte(apiKey))
  mac.Write(requestBody)
  rawBodyMAC := mac.Sum(nil)
  computedHash := hex.EncodeToString(rawBodyMAC)

  fmt.Println("************************************************")
  fmt.Println("**  computedHash = ", computedHash)
  fmt.Println("************************************************")      

  if !hmac.Equal([]byte(cleanedSignature), []byte(computedHash)){
    fmt.Println("************************************************")
    fmt.Println("** Inválid Pagarme Signature: Expected: %v, Received: %v", string(cleanedSignature), string(computedHash))
    fmt.Println("************************************************")    
    return false
  }
  
  fmt.Println("************************************************")
  fmt.Println("** Válido Pagarme Signature: Expected: %v, Received: %v", string(cleanedSignature), string(computedHash))
  fmt.Println("************************************************")        
  
  return true
}