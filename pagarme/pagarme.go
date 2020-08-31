
package pagarme

import (
	"github.com/mobilemindtec/go-utils/beego/validator"  
	"github.com/mobilemindtec/go-payments/support"  
  "github.com/astaxie/beego/validation"
  "github.com/leekchan/accounting"
  "github.com/beego/i18n"
  "encoding/base64"
	"encoding/json"
  "encoding/hex"
  "crypto/sha1"
  "crypto/hmac"	
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
	"net/url"
  "errors"
	"bytes"
	"fmt"	
)

const (

	PAGARME_URL = "https://api.pagar.me/1"
	PAGARME_TRANSACTION_URL = "https://api.pagar.me/1/transactions"
	PAGARME_TRANSACTION_CAPTURE_URL = "https://api.pagar.me/1/transactions/%v/capture"
	PAYMENT_METHOD_BOLETO = "boleto"
	PAYMENT_METHOD_CREDIT_CARD = "credit_card"

	PagarmeContry = "Brasil"
	PagarmeTypeIndividual = "individual"
	PagarmeTypeCorporation = "corporation"	

)

type PagarmeStatus int

/*

processing	Transação está processo de autorização.
authorized	Transação foi autorizada. Cliente possui saldo na conta e este valor foi reservado para futura captura, que deve acontecer em até 5 dias para transações criadas com api_key. Caso não seja capturada, a autorização é cancelada automaticamente pelo banco emissor, e o status da transação permanece authorized.
paid	Transação paga. Foi autorizada e capturada com sucesso, e para boleto, significa que nossa API já identificou o pagamento de seu cliente.
refunded	Transação estornada completamente.
waiting_payment	Transação aguardando pagamento (status válido para boleto bancário).
pending_refund	Transação do tipo boleto e que está aguardando para confirmação do estorno solicitado.
refused	Transação recusada, não autorizada.
chargedback	Transação sofreu chargeback. Mais em nossa central de ajuda
analyzing	Transação encaminhada para a análise manual feita por um especialista em prevenção a fraude.
pending_review	Transação pendente de revisão manual por parte do lojista. Uma transação ficar

*/

const (
  PagarmeProcessing PagarmeStatus = 1 + iota
  PagarmeAuthorized         
  PagarmePaid 
  PagarmeRefunded
  PagarmeWaitingPayment 
  PagarmePendingRefund 
  PagarmeRefused
	PagarmeChargedback
	PagarmeAnalyzing
	PagarmePendingReview
)

type PagarmeError struct {
	Message string `json:"message"`
	ParameterName string `json:"parameter_name"`
	Type string `json:"type"`
}

type PagarmeResultError struct {
	Method string `json:"method"`
	Url string `json:"url"`
	Errors []*PagarmeError `json:"errors"`	
}

func NewPagarmeResultError() *PagarmeResultError{
	return &PagarmeResultError{ Errors: []*PagarmeError{} }
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

type PagarmeCustomerDocument struct {
	Type string `json:"type"`
	Number string `json:"number"`
}

type PagarmeCustomer struct {

	DocumentNumber string `json:"document_number" valid:"Required;MaxSize(11);MinSize(11)"`
	Email string `json:"email" valid:"Required;Email"`
	Name string `json:"name" valid:"Required"`

	Address *PagarmeAddress `json:"address" valid:"Required"`
	Phone *PagarmePhone `json:"phone" valid:"Required"`
	ApiKey string `json:"api_key" valid:"Required"`	

	Id int64 `json:"id"` //

}

func NewPagarmeCustomer() *PagarmeCustomer {
	entity := new(PagarmeCustomer)
	return entity
}

/*

	Exemplo de uma transação de R$ 100, onde 99 vai para o cliente e 1 real vai para a mobile mind

	cliente := new(PagarmeSplitRule)
	cliente.Liable = true
	cliente.ChargeProcessingFee = true
	//cliente.Percentage = 100 // apenas no caso de percentual
	cliente.ChargeRemainderFee = true
	cliente.RecipientId = id do recebedor no pagarme
	cliente.Amount = 99 // 99 reais

	mobilemind := new(PagarmeSplitRule)
	mobilemind.Liable = false
	mobilemind.ChargeProcessingFee = false
	//mobilemind.Percentage = 0 // apenas no caso de percentual
	mobilemind.ChargeRemainderFee = false
	mobilemind.RecipientId = id do recebedor mobile mind no pagarme
	mobilemind.Amount = 1 // 1 real


*/

type PagarmeSplitRule struct {
	Liable bool `json:"liable"` // Se o recebedor é responsável ou não pelo chargeback. Default true para todos os recebedores da transação.
	ChargeProcessingFee bool `json:"charge_processing_fee"` // Se o recebedor será cobrado das taxas da criação da transação. Default true para todos os recebedores da transação.
	Percentage int64 `json:"percentage"` // Qual a porcentagem que o recebedor receberá. Deve estar entre 0 e 100. Se amount já está preenchido, não é obrigatório
	Amount int64 `json:"amount"` // Qual o valor da transação o recebedor receberá. Se percentage já está preenchido, não é obrigatório

	ChargeRemainderFee bool `json:"charge_remainder_fee"` //Se o recebedor deverá pagar os eventuais restos das taxas, calculadas em porcentagem. Sendo que o default vai para o primeiro recebedor definido na regra.
	RecipientId string `json:"recipient_id"` // Id do recebedor
}

type PagarmePayment struct {

	Amount int64 `json:"amount" valid:"Required"`
	ApiKey string `json:"api_key" valid:"Required"`	
 	Installments int `json:"Installments,omitempty" valid:"Required"`	// parcelas
	Customer *PagarmeCustomer `json:"customer" valid:"Required"`
	PaymentMethod string `json:"payment_method" valid:"Required"`

	
	PostbackUrl string `json:"postback_url,omitempty" valid:""`
	SoftDescriptor string `json:"soft_descriptor,omitempty" valid:"Required"` // nome que aparece na fatura do cliente
	Metadata map[string]string `json:"metadata,omitempty"`

	Capture bool `json:"capture"`

	BoletoExpirationDate string `json:"boleto_expiration_date,omitempty" valid:""`
	BoletoInstructions string `json:"boleto_instructions,omitempty" valid:""`

	CardId string `json:"card_id,omitempty" valid:""`
	CardHolderName string `json:"card_holder_name,omitempty" valid:""`
	CardExpirationDate string `json:"card_expiration_date,omitempty" valid:""`
	CardNumber	string `json:"card_number,omitempty" valid:""`
	CardCvv string `json:"card_cvv,omitempty" valid:""`	
	CardHash string `json:"card_hash,omitempty" valid:""`

	SplitRules []*PagarmeSplitRule `json:"split_rules,omitempty"`

}

type PagarmeCard struct {
	Id string `json:"card_id" valid:""`
	HolderName string `json:"card_holder_name" valid:""`
	ExpirationDate string `json:"card_expiration_date" valid:""`
	Number	string `json:"card_number" valid:""`
	Cvv string `json:"card_cvv" valid:""`
	CustomerId string `json:"customer_id" valid:""`
	Hash string `json:"card_hash" valid:""`
	ApiKey string `json:"api_key" valid:"Required"`	
}

type CardHashKey struct {
	Id int64 `json:"id"`
	PublicKey string `json:"public_key"`
	Hash string
}

/*
processing	Transação está processo de autorização.
authorized	Transação foi autorizada. Cliente possui saldo na conta e este valor foi reservado para futura captura, que deve acontecer em até 5 dias para transações criadas com api_key. Caso não seja capturada, a autorização é cancelada automaticamente pelo banco emissor, e o status da transação permanece authorized.
paid	Transação paga. Foi autorizada e capturada com sucesso, e para boleto, significa que nossa API já identificou o pagamento de seu cliente.
refunded	Transação estornada completamente.
waiting_payment	Transação aguardando pagamento (status válido para boleto bancário).
pending_refund	Transação do tipo boleto e que está aguardando para confirmação do estorno solicitado.
refused	Transação recusada, não autorizada.
chargedback	Transação sofreu chargeback. Mais em nossa central de ajuda
analyzing	Transação encaminhada para a análise manual feita por um especialista em prevenção a fraude.
pending_review	Transação pendente de revisão manual por parte do lojista. Uma transação ficará com esse status por até 48 horas corridas.
*/

type PagarmeResponse struct {

	Object string `json:"object"`
	StatusStr string `json:"status"` // processing, authorized, paid, refunded, waiting_payment, pending_refund, refused
	RefuseReason string `json:"refuse_reason"` // acquirer, antifraud, internal_error, no_acquirer, acquirer_timeout
	StatusReason string `json:"status_reason"` // acquirer, antifraud, internal_error, no_acquirer, acquirer_timeout
	AcquirerName string `json:"acquirer_name"` // tone, cielo, rede	
	AcquirerId string `json:"acquirer_id"`
	AcquirerResponseCode string `json:"acquirer_response_code"`
	AuthorizationCode string `json:"authorization_code"`
	SoftDescriptor string `json:"soft_descriptor"`
	Tid string `json:"tid"`
	Nsu string `json:"nsu"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	Amount int64 `json:"amount"`
	AuthorizedAmount int64 `json:"authorized_amount"`
	PaidAmount int64 `json:"paid_amount"`
	RefundedAmount int64 `json:"refunded_amount"`
	Installments int64 `json:"installments"`
	TransactionId int64 `json:"id"`
	Cost float64 `json:"cost"`
	CardHolderName string `json:"card_holder_name"`
	CardLastDigits string `json:"card_last_digits"`
	CardFirstDigits string `json:"card_first_digits"`
	CardBrand string `json:"card_brand"`
	CardPinMode string `json:"card_pin_mode"`
	PostbackUrl string `json:"postback_url"`
	PaymentMethod string `json:"payment_method"`
	CaptureMethod string `json:"capture_method"`
	AntifraudScore string `json:"antifraud_score"`
	BoletoUrl string `json:"boleto_url"`
	BoletoBarcode string `json:"boleto_barcode"`
	BoletoExpirationDate string `json:"boleto_expiration_date"`
	Referer string `json:"referer"`
	Ip string `json:"ip"`
	ReferenceKey string `json:"reference_key"`

	Errors []*PagarmeError `json:"errors"`
	ResponseValues map[string]interface{}
	Response string	
	Request string	
	Message string	
	Error bool

	Status PagarmeStatus


}

func NewPagarmePaymentCard(amount float64) *PagarmePayment {
	return &PagarmePayment{ Amount: FormatAmount(amount), Installments: 1, PaymentMethod: PAYMENT_METHOD_CREDIT_CARD }
}


func NewPagarmePaymentBoleto(amount float64) *PagarmePayment {
	return &PagarmePayment{ PaymentMethod: PAYMENT_METHOD_BOLETO, Amount: FormatAmount(amount)  }
}

type PagarmeService struct {
  Lang string  
  ApiKey string
  CryptoKey string
  EntityValidator *validator.EntityValidator
  EntityValidatorResult *validator.EntityValidatorResult
  Debug bool
}

type CaptureData struct {
	ApiKey string `json:"api_key" valid:"Required"`	
	TransactionId int64 `json:"-" valid:"Required"` 
	Amount int64 `json:"amount" valid:"Required"` 
	SplitRules []*PagarmeSplitRule `json:"split_rules,omitempty"`
	Metadata map[string]string `json:"metadata"`
}

func NewCaptureData(transactionId int64, amount float64) *CaptureData {	
	return &CaptureData{ TransactionId: transactionId, Amount:  FormatAmount(amount) }
}


func NewPagarmeService(lang string) *PagarmeService{
	entityValidator := validator.NewEntityValidator(lang, "Pagarme")
  return &PagarmeService{ Lang: lang, EntityValidator: entityValidator }
}

func NewPagarmeServiceWithCert(lang string, apiKey string, cryptoKey string) *PagarmeService{
	entityValidator := validator.NewEntityValidator(lang, "Pagarme")
  return &PagarmeService{ Lang: lang, EntityValidator: entityValidator, ApiKey: apiKey, CryptoKey: cryptoKey }
}

func (this *PagarmeService) GetCardHashKey() (*CardHashKey, error) {


	r, err := http.Get(fmt.Sprintf("%v/transactions/card_hash_key?encryption_key=%v", PAGARME_URL, this.CryptoKey))

	if err != nil {
		fmt.Println("error http.Get ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	response, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("error ioutil.ReadAll ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	fmt.Println(string(response))	


	result := new(CardHashKey)

	err = json.Unmarshal(response, result)

	if err != nil {
		fmt.Println("error json.Unmarshal ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	return result, nil
}

func (this *PagarmeService) EncryptCard(card *PagarmeCard)(*CardHashKey, error){

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

	fmt.Println("encodedCard = %v", encodedCard)	
	//fmt.Println("hasKeyInfo.PublicKey = %v", CardHashKey.PublicKey)

	encryptedData, err := support.RsaEncrypt([]byte(encodedCard), []byte(CardHashKey.PublicKey))

	if err != nil {
		return nil, err
	}

	//fmt.Println("encryptedData = %v", encryptedData)

	encryptedText := base64.StdEncoding.EncodeToString(encryptedData)

	fmt.Println("encryptedText = %v", encryptedText)

	CardHashKey.Hash =  fmt.Sprintf("%v_%v", CardHashKey.Id, encryptedText)
	
	return CardHashKey, nil

}


func (this *PagarmeService) CreateCard(card *PagarmeCard) (map[string]string, error) {

	//if !this.onValid(card) {
	//	return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	//}

	
	cardHash, err := this.EncryptCard(card)

	if err != nil {
		fmt.Println("error EncryptCard: %v", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	card.ApiKey = this.ApiKey
	card.Hash = cardHash.Hash

	jsonData, err := json.Marshal(card)

	if err != nil {
		fmt.Println("error json.Marshal ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	data := bytes.NewBuffer(jsonData)

	r, err := http.Post(fmt.Sprintf("%v/cards", PAGARME_URL), "application/json", data)

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


func (this *PagarmeService) CreateCustomer(customer *PagarmeCustomer) (error) {

	//if !this.onValid(card) {
	//	return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	//}

	customer.ApiKey = this.ApiKey
	

	jsonData, err := json.Marshal(customer)

	if err != nil {
		fmt.Println("error json.Marshal ", err.Error())
		return errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	data := bytes.NewBuffer(jsonData)

	r, err := http.Post(fmt.Sprintf("%v/customers", PAGARME_URL), "application/json", data)

	if err != nil {
		fmt.Println("error http.Post ", err.Error())
		return errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	response, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("error ioutil.ReadAll ", err.Error())
		return errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	fmt.Println(string(response))	

	
	err = json.Unmarshal(response, customer)

	if err != nil {
		fmt.Println("error json.Unmarshal ", err.Error())
		return errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	return nil

}

func (this *PagarmeService) CreatePayment(payment *PagarmePayment) (*PagarmeResponse, error) {
	

	this.EntityValidatorResult = new(validator.EntityValidatorResult)
	this.EntityValidatorResult.Errors = map[string]string{}

	if payment.PaymentMethod == PAYMENT_METHOD_CREDIT_CARD {

		if len(payment.CardHash) == 0 {

			card := new(PagarmeCard)
			card.Number = payment.CardNumber
			card.HolderName = payment.CardHolderName
			card.ExpirationDate = payment.CardExpirationDate
			card.Cvv = payment.CardCvv
			cardInfo, err := this.EncryptCard(card)

			if err != nil {
				fmt.Println("error EncryptCard ", err.Error())
				return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))			
			}

			payment.CardHash = cardInfo.Hash

		}

	}

	payment.ApiKey = this.ApiKey

	if !this.onValid(payment) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}




	jsonData, err := json.Marshal(payment)

	if err != nil {
		fmt.Println("error json.Marshal ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	reqtestString := string(jsonData)

	if this.Debug {
		x, _ := json.MarshalIndent(payment, "", "    ")	
		fmt.Println("##################################################################")
		fmt.Println(string(x))
		fmt.Println("##################################################################")
	}

	data := bytes.NewBuffer(jsonData)

	r, err := http.Post(PAGARME_TRANSACTION_URL, "application/json", data)

	if err != nil {
		fmt.Println("error http.Post ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	response, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("error ioutil.ReadAll ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Println("***** PAGARME PAYMENT START RESPONSE ****** ")	
	fmt.Println("***** STATUS CODE: %v", r.StatusCode)	
	fmt.Println("***** RESPONSE: %v", string(response))	
	fmt.Println("***** PAGARME PAYMENT END RESPONSE ****** ")
	fmt.Println("------------------------------------------------------------------------------------")

	switch r.StatusCode {
		case 200:

			result := &PagarmeResponse{ Response: string(response), Request: reqtestString  }
			//values := make(map[string]interface{})

			if err := json.Unmarshal(response, &result); err != nil {
				fmt.Println("error json.Unmarshal: %v", err)
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}

			//result.ResponseValues = values


			switch result.StatusStr {
			  case "processing":
			  	result.Status = PagarmeProcessing
			  case "authorized":
			  	result.Status = PagarmeAuthorized
			  case "paid":
			  	result.Status = PagarmePaid
			  case "refunded":
			  	result.Status = PagarmeRefunded
			  case "waiting_payment":
			  	result.Status = PagarmeWaitingPayment
			  case "pending_refund":
			  	result.Status = PagarmeRefunded
			  case "chargedback":
			  	result.Status = PagarmeChargedback
			  case "analyzing":
			  	result.Status = PagarmeAnalyzing
			  case "pending_review":
			  	result.Status = PagarmePendingReview
			  case "refused":
			  	result.Message = "A transação foi recusada, verifique os dados cartão"
			  	result.Error = true
			  	result.Status = PagarmeRefused
			  default:
			  	return result, errors.New(fmt.Sprintf("Problemas na trasanção. Status %v não reconhecido.", result.StatusStr))				
			}		

			
			return result, nil

		case 400:
			

			result := &PagarmeResponse{ Response: string(response), Request: reqtestString  }
			
			if err := json.Unmarshal(response, result); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}
		
			
			for _, it := range result.Errors {
				 this.EntityValidatorResult.Errors[it.ParameterName] = fmt.Sprintf("%v -  %v", it.Type, it.Message)
			}
			this.EntityValidatorResult.HasError = true

			if(len(result.Errors) > 0){				
				return nil, errors.New(fmt.Sprintf("Pagarme %v: %v, %v", result.Errors[0].ParameterName, result.Errors[0].Type, result.Errors[0].Message))				
			}

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

func (this *PagarmeService) Capture(captureData *CaptureData) (*PagarmeResponse, error) {

		this.EntityValidatorResult = new(validator.EntityValidatorResult)
		this.EntityValidatorResult.Errors = map[string]string{}

	captureData.ApiKey = this.ApiKey

	if !this.onValidCapture(captureData) {
		return nil, errors.New(this.getMessage("Pagarme.ValidationError"))
	}

	jsonData, err := json.Marshal(captureData)

	if err != nil {
		fmt.Println("error json.Marshal ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	reqtestString := string(jsonData)

	if this.Debug {
		fmt.Println("##################################################################")
		fmt.Println(reqtestString)	
		fmt.Println("##################################################################")
	}

	data := bytes.NewBuffer(jsonData)

	url := fmt.Sprintf(PAGARME_TRANSACTION_CAPTURE_URL, captureData.TransactionId) 

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

			result := &PagarmeResponse{ Response: string(response), Request: reqtestString  }
			/*
			values := make(map[string]interface{})

			if err := json.Unmarshal(response, &values); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}*/

			if err := json.Unmarshal(response, &result); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}

			//result.ResponseValues = values

			switch result.StatusStr {
			  case "processing":
			  	result.Status = PagarmeProcessing
			  case "authorized":
			  	result.Status = PagarmeAuthorized
			  case "paid":
			  	result.Status = PagarmePaid
			  case "refunded":
			  	result.Status = PagarmeRefunded
			  case "waiting_payment":
			  	result.Status = PagarmeWaitingPayment
			  case "pending_refund":
			  	result.Status = PagarmeRefunded
			  case "chargedback":
			  	result.Status = PagarmeChargedback
			  case "analyzing":
			  	result.Status = PagarmeAnalyzing
			  case "pending_review":
			  	result.Status = PagarmePendingReview
			  case "refused":
			  	result.Message = "A transação foi recusada, verifique os dados cartão"
			  	result.Error = true
			  	result.Status = PagarmeRefused
			  default:
			  	return result, errors.New(fmt.Sprintf("Problemas na trasanção. Status %v não reconhecido.", result.StatusStr))				
			}			

			
			return result, nil
		case 400:
			
			result := &PagarmeResponse{ Response: string(response), Request: reqtestString  }
			
			if err := json.Unmarshal(response, result); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}
		
			for _, it := range result.Errors {
				 this.EntityValidatorResult.Errors[it.ParameterName] = fmt.Sprintf("%v -  %v", it.Type, it.Message)
			}
			this.EntityValidatorResult.HasError = true

			if(len(result.Errors) > 0){				
				return nil, errors.New(fmt.Sprintf("Pagarme %v: %v, %v", result.Errors[0].ParameterName, result.Errors[0].Type, result.Errors[0].Message))				
			}

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

func (this *PagarmeService) FindPayment(id string) (*PagarmeResponse, error) {
	

	this.EntityValidatorResult = new(validator.EntityValidatorResult)
	this.EntityValidatorResult.Errors = map[string]string{}

	if this.Debug {
		fmt.Println("##################################################################")
		fmt.Println("Transaction Id = %v", id)	
		fmt.Println("##################################################################")
	}

	url := fmt.Sprintf("%v/transactions/%v?api_key=%v", PAGARME_URL, id, this.ApiKey) 

	r, err := http.Get(url)

	if err != nil {
		fmt.Println("error http.Get ", err.Error())
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

			result := &PagarmeResponse{ Response: string(response)  }
			/*
			values := make(map[string]interface{})

			if err := json.Unmarshal(response, &values); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}
			*/

			if err := json.Unmarshal(response, &result); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}

			//result.ResponseValues = values

			switch result.StatusStr {
			  case "processing":
			  	result.Status = PagarmeProcessing
			  case "authorized":
			  	result.Status = PagarmeAuthorized
			  case "paid":
			  	result.Status = PagarmePaid
			  case "refunded":
			  	result.Status = PagarmeRefunded
			  case "waiting_payment":
			  	result.Status = PagarmeWaitingPayment
			  case "pending_refund":
			  	result.Status = PagarmeRefunded
			  case "chargedback":
			  	result.Status = PagarmeChargedback
			  case "analyzing":
			  	result.Status = PagarmeAnalyzing
			  case "pending_review":
			  	result.Status = PagarmePendingReview
			  case "refused":
			  	result.Message = "A transação foi recusada, verifique os dados cartão"
			  	result.Error = true
			  	result.Status = PagarmeRefused
			  default:
			  	return result, errors.New(fmt.Sprintf("Problemas na trasanção. Status %v não reconhecido.", result.StatusStr))				
			}			

			
			return result, nil
		case 400:
			
			result := &PagarmeResponse{ Response: string(response)  }
			
			if err := json.Unmarshal(response, result); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}
		
			for _, it := range result.Errors {
				 this.EntityValidatorResult.Errors[it.ParameterName] = fmt.Sprintf("%v -  %v", it.Type, it.Message)
			}
			this.EntityValidatorResult.HasError = true

			if(len(result.Errors) > 0){				
				return nil, errors.New(fmt.Sprintf("Pagarme %v: %v, %v", result.Errors[0].ParameterName, result.Errors[0].Type, result.Errors[0].Message))				
			}

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

func (this *PagarmeService) RefundPayment(id string) (*PagarmeResponse, error) {
	
	

	this.EntityValidatorResult = new(validator.EntityValidatorResult)
	this.EntityValidatorResult.Errors = map[string]string{}

	payload := map[string]string{}
	payload["api_key"] = this.ApiKey

	jsonData, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("error json.Marshal ", err.Error())
		return nil, errors.New(this.getMessage("Pagarme.Error", err.Error()))
	}

	data := bytes.NewBuffer(jsonData)

	url := fmt.Sprintf("%v/transactions/%v/refund", PAGARME_URL, id)

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

			result := &PagarmeResponse{ Response: string(response)  }
			/*
			values := make(map[string]interface{})

			if err := json.Unmarshal(response, &values); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}*/

			if err := json.Unmarshal(response, &result); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}

			//result.ResponseValues = values

			switch result.StatusStr {
			  case "processing":
			  	result.Status = PagarmeProcessing
			  case "authorized":
			  	result.Status = PagarmeAuthorized
			  case "paid":
			  	result.Status = PagarmePaid
			  case "refunded":
			  	result.Status = PagarmeRefunded
			  case "waiting_payment":
			  	result.Status = PagarmeWaitingPayment
			  case "pending_refund":
			  	result.Status = PagarmeRefunded
			  case "chargedback":
			  	result.Status = PagarmeChargedback
			  case "analyzing":
			  	result.Status = PagarmeAnalyzing
			  case "pending_review":
			  	result.Status = PagarmePendingReview
			  case "refused":
			  	result.Message = "A transação foi recusada, verifique os dados cartão"
			  	result.Error = true
			  	result.Status = PagarmeRefused
			  default:
			  	return result, errors.New(fmt.Sprintf("Problemas na trasanção. Status %v não reconhecido.", result.StatusStr))				
			}			

			
			return result, nil
		case 400:
			
			result := &PagarmeResponse{ Response: string(response)  }
			
			if err := json.Unmarshal(response, result); err != nil {
				return nil, errors.New(fmt.Sprintf("Pagarme: Error on converte response to json - %v", err.Error()))
			}
		
			for _, it := range result.Errors {
				 this.EntityValidatorResult.Errors[it.ParameterName] = fmt.Sprintf("%v -  %v", it.Type, it.Message)
			}
			this.EntityValidatorResult.HasError = true

			if(len(result.Errors) > 0){				
				return nil, errors.New(fmt.Sprintf("Pagarme %v: %v, %v", result.Errors[0].ParameterName, result.Errors[0].Type, result.Errors[0].Message))				
			}

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



func (this *PagarmeService) onValid(payment *PagarmePayment) bool {
  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(payment, func (validator *validation.Validation) {
  	
  	emptyCardHash := len(strings.TrimSpace(payment.CardHash)) == 0
  	emptyCardId := len(strings.TrimSpace(payment.CardId)) == 0


  	if payment.PaymentMethod == PAYMENT_METHOD_CREDIT_CARD {
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
		}else{
			if len(strings.TrimSpace(payment.BoletoExpirationDate)) == 0 {
				validator.SetError("BoletoExpirationDate", this.getMessage("Pagarme.rquired"))
			} 					
		}		

  })
  return this.EntityValidatorResult.HasError == false
}

func (this *PagarmeService) onValidCapture(captureData *CaptureData) bool {
  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(captureData, nil)
  return this.EntityValidatorResult.HasError == false
}

func FormatAmount(amount float64) int64 {
	ac := accounting.Accounting{Symbol: "", Precision: 2, Thousand: "", Decimal: ""}
	text := strings.Replace(ac.FormatMoney(amount), ",", "", -1)
	text = strings.Replace(text, ".", "", -1)
	val, _ := strconv.Atoi(text)
	return int64(val)
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