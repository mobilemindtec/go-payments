
package pagarme

import (
	"github.com/mobilemindtec/go-utils/beego/validator"  
	"github.com/mobilemindtec/go-payments/support"  
  "github.com/beego/beego/v2/core/validation"
  "github.com/mobilemindtec/go-payments/api"
  "github.com/beego/i18n"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
  "errors"
	"bytes"
	"fmt"	
)

type ResultProcessor func(data []byte, response *Response) error

const (

	PAGARME_URL = "https://api.pagar.me/1"
	PagarmeContry = "Brasil"
	PagarmeTypeIndividual = "individual"
	PagarmeTypeCorporation = "corporation"	

)

type Pagarme struct {
  Lang string  
  ApiKey string
  CryptoKey string
  EntityValidator *validator.EntityValidator
  EntityValidatorResult *validator.EntityValidatorResult
  ValidationErrors map[string]string
  HasValidationError bool    
  Debug bool
}

func NewPagarme(lang string, apiKey string, cryptoKey string) *Pagarme{
	entityValidator := validator.NewEntityValidator(lang, "Pagarme")
	entityValidatorResult := new(validator.EntityValidatorResult)
	entityValidatorResult.Errors = map[string]string{}
  return &Pagarme{ Lang: lang, ApiKey: apiKey, CryptoKey: cryptoKey, EntityValidator: entityValidator, EntityValidatorResult: entityValidatorResult }
}

func (this *Pagarme) SetDebug() {
	this.Debug = true
}

func (this *Pagarme) GetCardHashKey() (*CardHashKey, error) {

  resultProcessor := func(data []byte, response *Response) error {    
  	response.CardHashKey = new(CardHashKey)
    return json.Unmarshal(data, response.CardHashKey)
  }

	result, err := this.get(fmt.Sprintf("transactions/card_hash_key?encryption_key=%v", this.CryptoKey), resultProcessor)

	if err != nil {
		return nil, err
	}

	return result.CardHashKey, nil
}

func (this *Pagarme) EncryptCard(card *Card)(*CardHashKey, error){

	CardHashKey, err := this.GetCardHashKey()

	if err != nil {
		return nil, err
	}


	params := url.Values{}
	params.Add("card_number", card.Number)
	params.Add("card_holder_name", card.HolderName)
	params.Add("card_expiration_date", card.ExpirationDate)
	params.Add("card_cvv", card.Cvv)	


	encodedCard := params.Encode()

	if this.Debug {
		fmt.Println("encodedCard = %v", encodedCard)	
		fmt.Println("hasKeyInfo.PublicKey = %v", CardHashKey.PublicKey)
	}

	encryptedData, err := support.RsaEncrypt([]byte(encodedCard), []byte(CardHashKey.PublicKey))

	if err != nil {
		return nil, err
	}

	//fmt.Println("encryptedData = %v", encryptedData)

	encryptedText := base64.StdEncoding.EncodeToString(encryptedData)

	if this.Debug {
		fmt.Println("encryptedText = %v", encryptedText)
	}

	CardHashKey.Hash =  fmt.Sprintf("%v_%v", CardHashKey.Id, encryptedText)
	
	return CardHashKey, nil

}


func (this *Pagarme) TokenCreate(card *Card) (*Response, error) {
	
	cardHash, err := this.EncryptCard(card)

	if err != nil {
		fmt.Println("error EncryptCard: %v", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	card.ApiKey = this.ApiKey
	card.Hash = cardHash.Hash

	if !this.onValidEntity(card) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}

  resultProcessor := func(data []byte, response *Response) error {    
  	response.CardResult = new(CardResult)
    return json.Unmarshal(data, response.CardResult)
  }

	response, err := this.post(card, "cards", resultProcessor)

  if err != nil || response.Error {
    return response, err
  }

  response.Status = api.PagarmeSuccess

  return response, err

}

func (this *Pagarme) PaymentCreate(payment *Payment) (*Response, error) {
	

	if payment.PaymentMethod == api.PaymentTypeCreditCard {

		if len(payment.CardId) == 0 {

			card := new(Card)
			card.Number = payment.CardNumber
			card.HolderName = payment.CardHolderName
			card.ExpirationDate = payment.CardExpirationDate
			card.Cvv = payment.CardCvv

			if !this.onValidEntity(card) {
				return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
			}

			cardInfo, err := this.EncryptCard(card)

			if err != nil {
				fmt.Println("error EncryptCard ", err.Error())
				return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))			
			}

			payment.CardHash = cardInfo.Hash

		} 

	}

	payment.ApiKey = this.ApiKey

	if !this.onValidPayment(payment) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}

	return this.post(payment, "transactions", nil)
}

func (this *Pagarme) PaymentCapture(captureData *CaptureData) (*Response, error) {

	captureData.ApiKey = this.ApiKey

	if !this.onValidEntity(captureData) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}

	return this.post(captureData, fmt.Sprintf("transactions/%v/capture", captureData.TransactionId), nil)
}

func (this *Pagarme) PaymentGet(id string) (*Response, error) {
	
  if len(id) == 0 {
    this.SetValidationError("id", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

  return this.get(fmt.Sprintf("transactions/%v?api_key=%v", id, this.ApiKey), nil)
}

func (this *Pagarme) PaymentRefund(id string, amount float64) (*Response, error) {

  if len(id) == 0 {
    this.SetValidationError("id", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

  val := FormatAmount(amount)

  if val <= 0 {
    this.SetValidationError("Amount", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

	data := make(map[string]interface{})
  data["amount"] = val

  response, err := this.post(data, fmt.Sprintf("transactions/%v/refund?api_key=%v", id, this.ApiKey), nil)	

  if err != nil || response.Error {
    return response, err
  }

  response.Status = api.PagarmeRefunded

  return response, err
}

func (this *Pagarme) PlanoCreate(plano *Plano) (*Response, error) {

	plano.ApiKey = this.ApiKey

	if !this.onValidEntity(plano) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}

  if plano.Days <= 0 {
    this.SetValidationError("Days", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }	

  if plano.InvoiceReminder <= 0 {
    this.SetValidationError("InvoiceReminder", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }	

  resultProcessor := func(data []byte, response *Response) error {      	
    return json.Unmarshal(data, response.Plano)
  }

	return this.post(plano, "plans", resultProcessor)
}

func (this *Pagarme) PlanoUpdate(plano *Plano) (*Response, error) {

	plano.ApiKey = this.ApiKey

	if !this.onValidEntity(plano) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}

  if plano.Id <= 0 {
    this.SetValidationError("id", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }	

  resultProcessor := func(data []byte, response *Response) error {      	
    return json.Unmarshal(data, response.Plano)
  }

	return this.put(plano, fmt.Sprintf("plans/%v", plano.Id), resultProcessor)
}

func (this *Pagarme) PlanoGet(id int64) (*Response, error) {
	
  if id <= 0 {
    this.SetValidationError("id", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

  resultProcessor := func(data []byte, response *Response) error {      	
    return json.Unmarshal(data, response.Plano)
  }  

  return this.get(fmt.Sprintf("plans/%v?api_key=%v", id, this.ApiKey), resultProcessor)
}

func (this *Pagarme) SubscriptionCreate(subscription *Subscription) (*Response, error) {

	subscription.ApiKey = this.ApiKey

	if !this.onValidSubscription(subscription) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}

  if subscription.PaymentMethod == api.PaymentTypeCreditCard {

    if len(subscription.CardId) == 0 {

      card := new(Card)
      card.Number = subscription.CardNumber
      card.HolderName = subscription.CardHolderName
      card.ExpirationDate = subscription.CardExpirationDate
      card.Cvv = subscription.CardCvv

      if !this.onValidEntity(card) {
        return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
      }

      cardInfo, err := this.EncryptCard(card)

      if err != nil {
        fmt.Println("error EncryptCard ", err.Error())
        return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))     
      }

      subscription.CardHash = cardInfo.Hash

    } 

  }  


	return this.post(subscription, "subscriptions", nil)
}

func (this *Pagarme) SubscriptionUpdate(subscription *Subscription) (*Response, error) {

	subscription.ApiKey = this.ApiKey

	if !this.onValidSubscription(subscription) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}

  if subscription.Id <= 0 {
    this.SetValidationError("id", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }	

  if subscription.PaymentMethod == api.PaymentTypeCreditCard {

    if len(subscription.CardId) == 0 {

      card := new(Card)
      card.Number = subscription.CardNumber
      card.HolderName = subscription.CardHolderName
      card.ExpirationDate = subscription.CardExpirationDate
      card.Cvv = subscription.CardCvv

      if !this.onValidEntity(card) {
        return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
      }

      cardInfo, err := this.EncryptCard(card)

      if err != nil {
        fmt.Println("error EncryptCard ", err.Error())
        return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))     
      }

      subscription.CardHash = cardInfo.Hash

    } 

  }  

	return this.put(subscription, fmt.Sprintf("subscriptions/%v", subscription.Id), nil)
}

func (this *Pagarme) SubscriptionCancel(id string) (*Response, error) {
	
  if len(id) == 0 {
    this.SetValidationError("id", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

	data := map[string]string{}
	data["api_key"] = this.ApiKey

  response, err := this.post(data, fmt.Sprintf("subscriptions/%v/cancel", id), nil)

  if err != nil || response.Error {
    return response, err
  }

  response.Status = api.PagarmeCancelled

  return response, err  
}

func (this *Pagarme) SubscriptionGet(id string) (*Response, error) {
	
  if len(id) == 0 {
    this.SetValidationError("id", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

  return this.get(fmt.Sprintf("subscriptions/%v?api_key=%v", id, this.ApiKey), nil)
}

func (this *Pagarme) SubscriptionTransactionsGet(id string) (*Response, error) {
	
  if len(id) <= 0 {
    this.SetValidationError("id", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

  resultProcessor := func(data []byte, response *Response) error {      	
    return json.Unmarshal(data, &response.Transactions)
  }   

  return this.get(fmt.Sprintf("subscriptions/%v/transactions?api_key=%v", id, this.ApiKey), resultProcessor)
}


func (this *Pagarme) SubscriptionSkip(id string, charges int64) (*Response, error) {
	
  if len(id) == 0 {
    this.SetValidationError("id", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

  if charges <= 0 {
    this.SetValidationError("charges", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

	data := make(map[string]interface{})
	data["api_key"] = this.ApiKey
	data["charges"] = charges

  response, err := this.post(data, fmt.Sprintf("subscriptions/%v/settle_charge", id), nil)

  if err != nil || response.Error {
    return response, err
  }

  response.Status = api.PagarmeSuccess

  return response, err    
}

func (this *Pagarme) CurrentBalance(recebedorId string) (*Response, error) {
	
  if len(recebedorId) == 0 {
    this.SetValidationError("recebedorId", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

  return this.get(fmt.Sprintf("balance?recipient_id=%v&api_key=%v", recebedorId, this.ApiKey), nil)
}

func (this *Pagarme) Movements(filter *Filter) (*Response, error) {
	

  resultProcessor := func(data []byte, response *Response) error {      	
    return json.Unmarshal(data, &response.Movements)
  }   

  filter.ApiKey = this.ApiKey

  url := fmt.Sprintf("balance/operations%v", this.urlQuery(filter.ToMap()))

  return this.get(url, resultProcessor)
}

func (this *Pagarme) TransferCreate(transfer *Transfer) (*Response, error) {
	

  resultProcessor := func(data []byte, response *Response) error {      	
    return json.Unmarshal(data, &response.TransferResult)
  }   

  transfer.ApiKey = this.ApiKey

  return this.post(transfer, "transfers", resultProcessor)
}

func (this *Pagarme) TransferList(filter *Filter) (*Response, error) {
	

  resultProcessor := func(data []byte, response *Response) error {      	
    return json.Unmarshal(data, &response.TransferResults)
  }   

  filter.ApiKey = this.ApiKey

  url := fmt.Sprintf("transfers%v", this.urlQuery(filter.ToMap()))

  return this.get(url, resultProcessor)
}

func (this *Pagarme) TransferGet(id string) (*Response, error) {
	
  if len(id) == 0 {
    this.SetValidationError("id", "is required")
    return nil, errors.New(this.getMessage("Pagarme.ValidationError"))       
  }

  resultProcessor := func(data []byte, response *Response) error {      	
    return json.Unmarshal(data, &response.TransferResult)
  }   


  return this.get(fmt.Sprintf("transfers/%v?api_key=%v", id, this.ApiKey), resultProcessor)
}

func (this *Pagarme) get(action string, resultProcessor ResultProcessor) (*Response, error) {
  return this.request(nil, action, "GET", resultProcessor)
}

func (this *Pagarme) delete(action string) (*Response, error) {
  return this.request(nil, action, "DELETE", nil)
}

func (this *Pagarme) post(data interface{}, action string, resultProcessor ResultProcessor) (*Response, error) {
  return this.request(data, action, "POST", resultProcessor)
}

func (this *Pagarme) put(data interface{}, action string, resultProcessor ResultProcessor) (*Response, error) {
  return this.request(data, action, "PUT", resultProcessor)
}

func (this *Pagarme) request(data interface{}, action string, method string, resultProcessor ResultProcessor) (*Response, error) {

	result := NewResponse()
  

  var req *http.Request
  var err error

  client := new(http.Client)
  apiUrl := fmt.Sprintf("%v/%v", PAGARME_URL, action)

  this.Log("URL %v, METHOD = %v", apiUrl, method)

  if (method == "POST" || method == "PUT") && data != nil {

    payload, err := json.Marshal(data)

    if err != nil {
      fmt.Println("error json.Marshal ", err.Error())    
      return result, err
    }

    postData := bytes.NewBuffer(payload)


    result.Request = string(payload)

    if this.Debug {
      fmt.Println("****************** Pagarme Request ******************")
      fmt.Println(result.Request)
      fmt.Println("****************** Pagarme Request ******************")
    }


    req, err = http.NewRequest(method, apiUrl, postData)

  } else {
    req, err = http.NewRequest(method, apiUrl, nil)    
  }

  if err != nil {
    fmt.Println("err = ", err)
    return nil, errors.New(fmt.Sprintf("error on http.NewRequest: %v", err))
  }

  req.Header.Add("Content-Type", "application/json")

  res, err := client.Do(req)

  if err != nil {
    fmt.Println("err = %v", err)
    return nil, errors.New(fmt.Sprintf("error on client.Do: %v", err))
  }

  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)

  if err != nil {
    fmt.Println("err = %v", err)
    return nil, errors.New(fmt.Sprintf("error on ioutil.ReadAll: %v", err))
  }

  result.Response = string(body)

  if this.Debug {
    fmt.Println("****************** Pagarme Response ******************")
    fmt.Println("STATUS CODE ", res.StatusCode)
    fmt.Println(result.Response)
    fmt.Println("****************** Pagarme Response ******************")
  }

  if res.StatusCode == 200 || res.StatusCode == 400 {
    if resultProcessor != nil && res.StatusCode == 200 {
      if err := resultProcessor(body, result); err != nil {
        fmt.Println("err = %v", err)
        return nil, errors.New(fmt.Sprintf("error on resultProcessor: %v", err))      
      }
    } else {

      err = json.Unmarshal(body, result)

      if err != nil {
        fmt.Println("err = %v", err)
        return nil, errors.New(fmt.Sprintf("error on json.Unmarshal: %v", err))
      }

    }
  }

  if res.StatusCode == 400 {
    result.Error = true

    if result.ErrorsCount() == 1 {
      result.Message = result.FirstError()
    } else {
      result.Message = fmt.Sprintf("Pagarme validation errror")
    }
    
    return result, nil
  }

  if res.StatusCode != 200 {
    result.Error = true
    result.Message = fmt.Sprintf("Pagarme error. Status: %v", res.StatusCode)
    return result, errors.New(result.Message) 
  }

  result.Error = result.HasError()
  result.BuildStatus()

  return result, nil
}

func (this *Pagarme) onValidPayment(payment *Payment) bool {
  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(payment, func (validator *validation.Validation) {
  	
  	emptyCardHash := len(strings.TrimSpace(payment.CardHash)) == 0
  	emptyCardId := len(strings.TrimSpace(payment.CardId)) == 0



  	switch payment.PaymentMethod {
      case api.PaymentTypeCreditCard:
  			if emptyCardHash && emptyCardId {
  				if len(strings.TrimSpace(payment.CardHolderName)) == 0 {
  					validator.SetError("CardHolderName", this.getMessage("Pagarme.rquired"))
  				} 

  				if len(strings.TrimSpace(payment.CardExpirationDate)) == 0 {
  					validator.SetError("CardExpirationDate", this.getMessage("Pagarme.rquired"))
  				} 

  				if len(strings.TrimSpace(payment.CardNumber)) == 0 {
  					validator.SetError("CardNumber", this.getMessage("Pagarme.rquired"))
  				} 

  				if len(strings.TrimSpace(payment.CardCvv)) == 0 {
  					validator.SetError("CardCvv", this.getMessage("Pagarme.rquired"))
  				} 
  			}
        break
      case api.PaymentTypeBoleto:
        if len(strings.TrimSpace(payment.BoletoExpirationDate)) == 0 {
          validator.SetError("BoletoExpirationDate", this.getMessage("Pagarme.rquired"))
        }
        break        
      case api.PaymentTypePix:
        if len(strings.TrimSpace(payment.PixExpirationDate)) == 0 {
          validator.SetError("PixExpirationDate", this.getMessage("Pagarme.rquired"))
        }
        break        
		}		

  })

  if this.EntityValidatorResult.HasError {
    this.onValidationErrors()
    return false
  }

  return true
}

func (this *Pagarme) onValidSubscription(subscription *Subscription) bool {

  items := []interface{}{
    subscription,
    subscription.Customer,
  }

  this.EntityValidatorResult, _ = this.EntityValidator.ValidMult(items, func (validator *validation.Validation) {
  	
  	emptyCardHash := len(strings.TrimSpace(subscription.CardHash)) == 0
  	emptyCardId := len(strings.TrimSpace(subscription.CardId)) == 0


  	if subscription.PaymentMethod == api.PaymentTypeCreditCard {
			if emptyCardHash && emptyCardId {
				if len(strings.TrimSpace(subscription.CardHolderName)) == 0 {
					validator.SetError("CardHolderName", this.getMessage("Pagarme.rquired"))
				} 

				if len(strings.TrimSpace(subscription.CardExpirationDate)) == 0 {
					validator.SetError("CardExpirationDate", this.getMessage("Pagarme.rquired"))
				} 

				if len(strings.TrimSpace(subscription.CardNumber)) == 0 {
					validator.SetError("CardNumber", this.getMessage("Pagarme.rquired"))
				} 

				if len(strings.TrimSpace(subscription.CardCvv)) == 0 {
					validator.SetError("CardCvv", this.getMessage("Pagarme.rquired"))
				} 
			}
		}		

  })

  if this.EntityValidatorResult.HasError {
    this.onValidationErrors()
    return false
  }

  return true
}

func (this *Pagarme) onValidEntity(entity interface{}) bool {
  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(entity, nil)

  if this.EntityValidatorResult.HasError {
    this.onValidationErrors()
    return false
  }

  return true
}

func (this *Pagarme) getMessage(key string, args ...interface{}) string{
  return i18n.Tr(this.Lang, key, args)
}

func (this *Pagarme) onValidationErrors(){
  this.HasValidationError = true
  this.ValidationErrors = this.EntityValidator.GetValidationErrors(this.EntityValidatorResult)
}

func (this *Pagarme) SetValidationError(key string, value string){
  this.HasValidationError = true
  if this.ValidationErrors == nil {
    this.ValidationErrors = make(map[string]string)
  }
  this.ValidationErrors[key]= value
}

func (this *Pagarme) Log(message string, args ...interface{}) {
	if this.Debug {
    fmt.Println("Pagarme: ", fmt.Sprintf(message, args...))
  }
}

func (this *Pagarme) urlQuery(filter map[string]string) string {
  url := ""
  if filter != nil && len(filter) > 0 {
    url = fmt.Sprintf("%v?", url)    

    for k, v := range filter {
      url = fmt.Sprintf("%v%v=%v", url, k, v)    
      url = fmt.Sprintf("%v&", url)    
    }
  }  

  return url
}