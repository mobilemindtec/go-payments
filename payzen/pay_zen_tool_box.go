package payzen


import (	
	"github.com/mobilemindtec/go-payments/support"
	"github.com/mobilemindtec/go-utils/app/util"
	"github.com/mobilemindtec/go-payments/api"
	"github.com/leekchan/accounting"
	"github.com/satori/go.uuid"
	"encoding/base64"
	"crypto/sha256"
	"encoding/xml"
	_ "encoding/hex"
	"crypto/hmac"
	_ "crypto/sha1"
	"io/ioutil"
	"net/http"	
	"net/url"	
	"strconv"
	"strings"
  "errors"
	"bytes"	
	"sort"
	"fmt"	
	_ "time"
)

const (
  PAYZEN_API_URL_TEST_URL = "https://secure.payzen.com.br/vads-ws/v5?wsdl" //URL of the PayZen SOAP V5 WSDL
  PAYZEN_API_URL_PROD_URL = "https://secure.payzen.com.br/vads-ws/v5?wsdl" //URL of the PayZen SOAP V5 WSDL

	PAYZEN_FORM_URL = "https://secure.payzen.com.br/vads-payment/"

  ns = "http://v5.ws.vads.lyra.com/"  //Namespace of the service
  hns = "http://v5.ws.vads.lyra.com/Header" //Namespace ot the service header	
  CurrencyCode = "986"
  CountryBR = "BR"
  Currency = "BRL"
)

const (
	PayZenModeTest = "TEST"
	PayZenModeProduction = "PRODUCTION"
	UUIDBase = "1546058f-5a25-4334-85ae-e68f2a44bbaf"
	SchemeBoleto = "BOLETO"	
	DateTimeLayout = "2006-01-02T15:04:05Z"

	TimeZoneDateTimeLayout = "2006-01-02T15:04:05-07:00"	

	DateWithZeroTimeLayout = "2006-01-02T00:00:00Z"
)

type PayZenToolBox struct {
	Account *api.PayZenAccount
	Debug bool
}

func NewPayZenToolBox(account *api.PayZenAccount) *PayZenToolBox{
	return &PayZenToolBox{ Account: account }
}

func (this *PayZenToolBox) GenerateAuthToken(requestId string, timestamp string, format string) string{
	
	data := ""

	if format == "request" {
		data = requestId + timestamp
	} else {
		data = timestamp + requestId
	}

	mac := hmac.New(sha256.New, []byte(this.Account.Cert))
	mac.Write([]byte(data))
	raw := mac.Sum(nil)
	//computedHash := hex.EncodeToString(raw)

	base64Content := base64.StdEncoding.EncodeToString(raw)

	if this.Debug {
		fmt.Println("*********************************************")
		fmt.Println("** %v", base64Content)
		fmt.Println("*********************************************")
	}

	return base64Content
}

func (this *PayZenToolBox) GenerateFormAuthToken(data string) string{
	
	/*	
	hasher := sha1.New()
	hasher.Write([]byte(data))
	content := hex.EncodeToString(hasher.Sum(nil))
	fmt.Println("SHA1 content = %s", content)
	return content
	*/

	mac := hmac.New(sha256.New, []byte(this.Account.Cert))
	mac.Write([]byte(data))
	raw := mac.Sum(nil)
	//computedHash := hex.EncodeToString(raw)

	base64Content := base64.StdEncoding.EncodeToString(raw)

	if this.Debug {
		fmt.Println("*********************************************")
		fmt.Println("** %v", base64Content)
		fmt.Println("*********************************************")
	}

	return base64Content
	
}

func (this *PayZenToolBox) ValidateResponse(header *SOAPResponseHeader) bool{
	authToken := this.GenerateAuthToken(header.RequestId, header.Timestamp, "response")
	return authToken == header.AuthToken
}

func (this *PayZenToolBox) FillRequestHeader(header *SOAPHeader){
	header.ShopId = this.Account.ShopId
	header.Mode = this.Account.Mode
	header.Timestamp = util.DateNow().UTC().Format(DateTimeLayout) //time.Now().UTC().Format("2006-01-02T15:04:05Z")

	uuidBase, err := uuid.FromString(UUIDBase)

	if err != nil {
		fmt.Println("**** error on get uuid from string: %v", err)
	}

	header.RequestId = uuid.NewV5(uuidBase, header.Timestamp).String()
	header.AuthToken = this.GenerateAuthToken(header.RequestId, header.Timestamp, "request")
}


func (this *PayZenToolBox) FillCardInfo(cardRequest *SOAPCardRequest, card *api.Card) {
	cardRequest.Number = card.Number //"4970100000000007"
	cardRequest.Scheme = card.Scheme //"VISA"
	cardRequest.ExpiryMonth = card.ExpiryMonth //12
	cardRequest.ExpiryYear = card.ExpiryYear //2018
	cardRequest.CardSecurityCode = card.CardSecurityCode //"235"	
	cardRequest.CardHolderName = card.CardHolderName	
	cardRequest.PaymentToken = card.Token
}

func (this *PayZenToolBox) FillCustomerInfo(billingDetails *SOAPBillingDetails, customer *api.Customer) {	
	billingDetails.FirstName = customer.FirstName
	billingDetails.LastName = customer.LastName
	billingDetails.PhoneNumber = customer.PhoneNumber
	billingDetails.Email = customer.Email
	billingDetails.StreetNumber = customer.StreetNumber
	billingDetails.Address = customer.Address
	billingDetails.District = customer.District
	billingDetails.ZipCode = customer.ZipCode
	billingDetails.City = customer.City
	billingDetails.State = customer.State
	billingDetails.Country = customer.Country
	billingDetails.IdentityCode = customer.IdentityCode
}


/*
	Cria um token para compra com um click
*/
func (this *PayZenToolBox) CreatePaymentToken(paymentData *api.Payment) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run CreatePaymentToken: %v", paymentData)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	opCreateToken := NewSOAPCreateToken()
	
	opCreateToken.CommonRequest.SubmissionDate = soap.Header.Timestamp

	this.FillCardInfo(opCreateToken.CardRequest, paymentData.Card)
	this.FillCustomerInfo(opCreateToken.CustomerRequest.BillingDetails, paymentData.Customer)
	 
	soap.Body.OperationRequest = opCreateToken	

	soapResponse := NewSOAPResponseEnvelop()
	createTokenResponse := NewSOAPCreateTokenResponse()
	soapResponse.Body.CreateTokenResponse = createTokenResponse


	result, err := this.MakeRequest(soap, soapResponse)

	if err != nil || result.Error {
		return result, err
	}

	switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Success
			break
	}

	result.TokenInfo.Token = createTokenResponse.CreateTokenResult.CommonResponse.PaymentToken

	return result, err
}

/*
	Atualiza um token para compra com um click
*/

func (this *PayZenToolBox) UpdatePaymentToken(paymentData *api.Payment) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run UpdatePaymentToken: %v", paymentData)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	opUpdateToken := NewSOAPUpdateToken()
	
	opUpdateToken.CommonRequest.SubmissionDate = soap.Header.Timestamp

	this.FillCardInfo(opUpdateToken.CardRequest, paymentData.Card)
	this.FillCustomerInfo(opUpdateToken.CustomerRequest.BillingDetails, paymentData.Customer)
	
	opUpdateToken.QueryRequest.PaymentToken = paymentData.Card.Token

	soap.Body.OperationRequest = opUpdateToken	

	soapResponse := NewSOAPResponseEnvelop()
	updateTokenResponse := NewSOAPUpdateTokenResponse()
	soapResponse.Body.UpdateTokenResponse = updateTokenResponse


	result, err := this.MakeRequest(soap, soapResponse)

	result.TokenInfo.NotFound = strings.Contains(result.ResponseCodeDetail, "PaymentToken not found")

	if err != nil || result.Error {
		return result, err
	}


	switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Success
			break
	}

	result.TokenInfo.Token = updateTokenResponse.UpdateTokenResult.CommonResponse.PaymentToken

	return result, err
}

/*
	Cancela um token para compra com um click
*/
func (this *PayZenToolBox) CancelPaymentToken(paymentToken string) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run CancelPaymentToken: %v", paymentToken)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	opCancelToken := NewSOAPCancelToken()
	
	opCancelToken.CommonRequest.SubmissionDate = soap.Header.Timestamp
	
	opCancelToken.QueryRequest.PaymentToken = paymentToken

	soap.Body.OperationRequest = opCancelToken

	soapResponse := NewSOAPResponseEnvelop()
	cancelTokenResponse := NewSOAPCancelTokenResponse()
	soapResponse.Body.CancelTokenResponse = cancelTokenResponse


	result, err := this.MakeRequest(soap, soapResponse)

	if err != nil || result.Error {
		return result, err
	}

	switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Success
			break
	}

	result.TokenInfo.NotFound = strings.Contains(result.ResponseCodeDetail, "PaymentToken not found")

	return result, err
}

/*
	Reativa um token para compra com um click
*/
func (this *PayZenToolBox) ReactivePaymentToken(paymentToken string) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run ReactivePaymentToken: %v", paymentToken)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	opReactiveToken := NewSOAPReactiveToken()
	
	//opReactiveToken.CommonRequest.SubmissionDate = soap.Header.Timestamp
	
	opReactiveToken.QueryRequest.PaymentToken = paymentToken

	soap.Body.OperationRequest = opReactiveToken

	soapResponse := NewSOAPResponseEnvelop()
	reactiveTokenResponse := NewSOAPReactiveTokenResponse()
	soapResponse.Body.ReactiveTokenResponse = reactiveTokenResponse


	result, err := this.MakeRequest(soap, soapResponse)

	if err != nil || result.Error {
		return result, err
	}

	switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Success
			break
	}

	result.TokenInfo.NotFound = strings.Contains(result.ResponseCodeDetail, "PaymentToken not found")

	return result, err
}

/*
	Consulta um token para compra com um click
*/
func (this *PayZenToolBox) GetDetailsPaymentToken(paymentToken string) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run GetDetailsPaymentToken: %v", paymentToken)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	opGetDetailsToken := NewSOAPGetTokenDetails()
	
	//opReactiveToken.CommonRequest.SubmissionDate = soap.Header.Timestamp
	
	opGetDetailsToken.QueryRequest.PaymentToken = paymentToken

	soap.Body.OperationRequest = opGetDetailsToken

	soapResponse := NewSOAPResponseEnvelop()
	getDetailsTokenTokenResponse := NewSOAPGetTokenDetailsResponse()
	soapResponse.Body.GetTokenDetailsResponse = getDetailsTokenTokenResponse


	result, err := this.MakeRequest(soap, soapResponse)

	result.TokenInfo.NotFound = strings.Contains(result.ResponseCodeDetail, "PaymentToken not found")

	switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Success
			break
	}

	if !result.TokenInfo.NotFound {

		getTokenDetailsResult := getDetailsTokenTokenResponse.GetTokenDetailsResult

		result.TokenInfo.Number = getTokenDetailsResult.CardResponse.Number	
		result.TokenInfo.Brand = getTokenDetailsResult.CardResponse.Brand
		result.TokenInfo.CreationDate, _ =  util.DateParse(TimeZoneDateTimeLayout, getTokenDetailsResult.TokenResponse.CreationDate) 
		result.TokenInfo.CancellationDate, _ = util.DateParse(TimeZoneDateTimeLayout, getTokenDetailsResult.TokenResponse.CancellationDate) 

		result.TokenInfo.Cancelled = !result.TokenInfo.CancellationDate.IsZero()
		result.TokenInfo.Active = result.TokenInfo.CancellationDate.IsZero()		
	}


	return result, err
}

/*
	Retorna o status de um pagamento
*/
func (this *PayZenToolBox) FindPayment(orderId string) (*api.PaymentResult, error){

	if this.Debug {
		fmt.Println("** run FindPayment: %v", orderId)
	}

	// cria request
	soap := NewSOAPEnvelope()
	this.FillRequestHeader(soap.Header)	 
	
	findPayments := NewSOAPFindPayments()
	findPayments.QueryRequest.OrderId = orderId
	soap.Body.OperationRequest = findPayments  

	// cria response
	soapResponse := NewSOAPResponseEnvelop()
	soapResponse.Body.FindPaymentsResponse = NewSOAPFindPaymentsResponse()

	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentNotFound = strings.Contains(result.ResponseCodeDetail, "Transaction was not found")

	if !result.PaymentNotFound {				
		items := soapResponse.Body.FindPaymentsResponse.FindPaymentsResult.TransactionItem			
		for _, item := range items {

			res, _ := switchStatus(item.TransactionStatusLabel)

			result.TransactionStatus = res.TransactionStatus
			
			trans := new(api.TransactionItemResult)
			trans.TransactionUuid = item.TransactionUuid
			trans.TransactionStatus = res.TransactionStatus
			trans.TransactionStatusLabel = item.TransactionStatusLabel
			trans.Amount = parseAmount(item.Amount)			
			trans.ExpectedCaptureDate, _ = util.DateParse(TimeZoneDateTimeLayout, item.ExpectedCaptureDate)

			result.Transactions = append(result.Transactions, trans)

			if this.Debug {
				fmt.Println("********************************************************************************************************************")
				fmt.Println("********************************************************************************************************************")
				fmt.Println("***** TransactionItem: %v", item)
				fmt.Println("********************************************************************************************************************")
				fmt.Println("********************************************************************************************************************")
			}
		}		
	}

	return result, err
}

/*
	Retorna os detalhes de um pagamento
*/

func (this *PayZenToolBox) GetPaymentDetails(transactionUuid string) (*api.PaymentResult, error){
	return this.getPaymentDetails(transactionUuid, false)
}

func (this *PayZenToolBox) GetPaymentDetailsWithNsu(transactionUuid string) (*api.PaymentResult, error){
	return this.getPaymentDetails(transactionUuid, true)
}

func (this *PayZenToolBox) getPaymentDetails(transactionUuid string, withNsu bool) (*api.PaymentResult, error){

	if this.Debug {
		fmt.Println("** run GetPaymentDetails: %v", transactionUuid)
	}

	// cria request
	soap := NewSOAPEnvelope()
	this.FillRequestHeader(soap.Header)	 
	
	var getPaymentDetails *SOAPGetPaymentDetails

	if withNsu {
		getPaymentDetails = NewSOAPGetPaymentDetailsWithNsu()
	} else {
		getPaymentDetails = NewSOAPGetPaymentDetails()
	}
	
	getPaymentDetails.QueryRequest.Uuid = transactionUuid
	soap.Body.OperationRequest = getPaymentDetails  

	// cria response
	soapResponse := NewSOAPResponseEnvelop()
	getPaymentDetailsResponse := NewSOAPGetPaymentDetailsResponse()
	soapResponse.Body.GetPaymentDetailsResponse = getPaymentDetailsResponse

	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentNotFound = strings.Contains(result.ResponseCodeDetail, "Transaction was not found")

	if !result.PaymentNotFound {		
		res, erro := switchStatus(result.TransactionStatusLabel)

		if erro != nil {
			return result, errors.New(fmt.Sprintf("%v. Código de erro %v: %v", erro.Error(), result.ResponseCode, result.ResponseCodeDetail))
		}

		result.TransactionStatus = res.TransactionStatus

		paymentResponse := getPaymentDetailsResponse.GetPaymentDetailsResult.PaymentResponse

		trans := new(api.TransactionItemResult)
		trans.TransactionUuid = paymentResponse.TransactionUuid
		trans.TransactionId = paymentResponse.TransactionId
		trans.TransactionStatus = res.TransactionStatus
		trans.TransactionStatusLabel = result.TransactionStatusLabel
		trans.Amount = parseAmount(paymentResponse.Amount)			
		trans.ExpectedCaptureDate, _ = util.DateParse(TimeZoneDateTimeLayout, paymentResponse.ExpectedCaptureDate)
		trans.CreationDate, _ = util.DateParse(TimeZoneDateTimeLayout, paymentResponse.CreationDate)
		trans.ExternalTransactionId = paymentResponse.ExternalTransactionId
		result.Transactions = append(result.Transactions, trans)
		
	}	

	return result, err
}

/*
	Retorna os detalhes de um pagamento
*/
func (this *PayZenToolBox) ValidatePayment(transactionUuid string) (*api.PaymentResult, error){

	if this.Debug {
		fmt.Println("** run validatePaymentResponse: %v", transactionUuid)
	}

	// cria request
	soap := NewSOAPEnvelope()
	this.FillRequestHeader(soap.Header)	 
	
	operation := NewSOAPValidatePayment()
	operation.QueryRequest.Uuid = transactionUuid
	soap.Body.OperationRequest = operation  

	// cria response
	soapResponse := NewSOAPResponseEnvelop()
	respOp := NewSOAPValidatePaymentResponse()
	soapResponse.Body.ValidatePaymentResponse = respOp

	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentNotFound = strings.Contains(result.ResponseCodeDetail, "Transaction was not found")

	res, erro := switchStatus(result.TransactionStatusLabel)
	
	result.Error = res.Error
	result.Message = res.Message
	result.TransactionStatus = res.TransactionStatus

	if erro != nil { // TransactionStatusLabel not found
		return result, errors.New(fmt.Sprintf("%v. Código de erro %v: %v", erro.Error(), result.ResponseCode, result.ResponseCodeDetail))
	}

	return result, err
}

func (this *PayZenToolBox) DuplicatePayment(paymentData *api.Payment) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run DuplicatePayment: %v", paymentData)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	operation := NewSOAPDuplicatePayment()
	
	operation.CommonRequest.SubmissionDate = soap.Header.Timestamp

	operation.QueryRequest.Uuid = paymentData.TransactionUuid

	operation.OrderRequest.OrderId = paymentData.OrderId

	operation.PaymentRequest.Amount = formatAmount(paymentData.Amount)
	operation.PaymentRequest.Currency = CurrencyCode
	operation.PaymentRequest.ManualValidation = fmt.Sprintf("%v", paymentData.ValidationType)
		 
	soap.Body.OperationRequest = operation
	
	soapResponse := NewSOAPResponseEnvelop()
	respOp := NewSOAPDuplicatePaymentResponse()
	soapResponse.Body.DuplicatePaymentResponse = respOp	

	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentNotFound = strings.Contains(result.ResponseCodeDetail, "Transaction was not found")

	res, erro := switchStatus(result.TransactionStatusLabel)
	
	result.Error = res.Error
	result.Message = res.Message
	result.TransactionStatus = res.TransactionStatus

	if erro != nil { // TransactionStatusLabel not found
		return result, errors.New(fmt.Sprintf("A transação foi recusada. Status da transação: %v. Código de erro %v: %v", result.TransactionStatusLabel, result.ResponseCode, result.ResponseCodeDetail))
	}

	if err != nil || result.Error {
		return result, err
	}

	opResult := respOp.DuplicatePaymentResult
	
	if opResult.CardResponse.Scheme == SchemeBoleto {
		extraInfos := opResult.OrderResponse.ExtInfo
		extraInfo := extraInfos[0]
		result.BoletoUrl = extraInfo.Value
	}

	result.TransactionId = opResult.PaymentResponse.TransactionId
	result.TransactionUuid = opResult.PaymentResponse.TransactionUuid
	
	return result, err	
}

func (this *PayZenToolBox) RefundPayment(paymentData *api.Payment) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run RefundPayment: %v", paymentData)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	operation := NewSOAPRefundPayment()
	
	operation.CommonRequest.SubmissionDate = soap.Header.Timestamp

	operation.QueryRequest.Uuid = paymentData.TransactionUuid

	operation.PaymentRequest.Amount = formatAmount(paymentData.Amount)
	operation.PaymentRequest.Currency = CurrencyCode
	operation.PaymentRequest.ManualValidation = fmt.Sprintf("%v", paymentData.ValidationType)
		 
	soap.Body.OperationRequest = operation
	
	soapResponse := NewSOAPResponseEnvelop()
	respOp := NewSOAPRefundPaymentResponse()
	soapResponse.Body.RefundPaymentResponse = respOp	

	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentNotFound = strings.Contains(result.ResponseCodeDetail, "Transaction was not found")

	res, erro := switchStatus(result.TransactionStatusLabel)
	
	result.Error = res.Error
	result.Message = res.Message
	result.TransactionStatus = res.TransactionStatus

	if erro != nil { // TransactionStatusLabel not found
		return result, errors.New(fmt.Sprintf("%v. Código de erro %v: %v", erro.Error(), result.ResponseCode, result.ResponseCodeDetail))
	}

	if err != nil || result.Error {
		return result, err
	}

	opResult := respOp.RefundPaymentResult
	result.TransactionId = opResult.PaymentResponse.TransactionId
	result.TransactionUuid = opResult.PaymentResponse.TransactionUuid
	
	return result, err	
}

/*
	Captura uma transação
*/
func (this *PayZenToolBox) CapturePayment(captureObject *api.PaymentCapture) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run CapturePayment: %v", captureObject.TransactionUuids)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	capturePayment := NewSOAPCapturePayment()

	capturePayment.SettlementRequest.TransactionUuids = captureObject.TransactionUuids
	capturePayment.SettlementRequest.Date = util.DateNow().Format(DateWithZeroTimeLayout) 
	capturePayment.SettlementRequest.Commission = formatAmount(captureObject.Commission)
	soap.Body.OperationRequest = capturePayment

	soapResponse := NewSOAPResponseEnvelop()
	captureResponse := NewSOAPCreateCaptureResponse()
	soapResponse.Body.CapturePaymentResponse = captureResponse

	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentNotFound = strings.Contains(result.ResponseCodeDetail, "Transaction was not found")

	switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Captured
			break
	}


	return result, err
}

/*
	Autoriza um pagamento
*/
func (this *PayZenToolBox) CreatePayment(paymentData *api.Payment) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run CreatePayment: %v", paymentData)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	opCreatePayment := NewSOAPCreatePayment()
	
	opCreatePayment.CommonRequest.SubmissionDate = soap.Header.Timestamp

	opCreatePayment.PaymentRequest.PaymentOptionCode = int(paymentData.Installments) //parcelas
		
	opCreatePayment.PaymentRequest.Amount = formatAmount(paymentData.Amount)
	opCreatePayment.PaymentRequest.Currency = CurrencyCode

	opCreatePayment.PaymentRequest.ManualValidation = fmt.Sprintf("%v", paymentData.ValidationType)
	
	opCreatePayment.OrderRequest.OrderId = paymentData.OrderId
	
	this.FillCardInfo(opCreatePayment.CardRequest, paymentData.Card)
	this.FillCustomerInfo(opCreatePayment.CustomerRequest.BillingDetails, paymentData.Customer)
	 
	soap.Body.OperationRequest = opCreatePayment
	
	soapResponse := NewSOAPResponseEnvelop()
	createPaymentResponse := NewSOAPCreatePaymentResponse()
	soapResponse.Body.CreatePaymentResponse = createPaymentResponse	

	result, err := this.MakeRequest(soap, soapResponse)

	res, erro := switchStatus(result.TransactionStatusLabel)
	
	result.Error = res.Error
	result.Message = res.Message
	result.TransactionStatus = res.TransactionStatus

	if erro != nil { // TransactionStatusLabel not found
		return result, errors.New(fmt.Sprintf("%v. Código de erro %v: %v", erro.Error(), result.ResponseCode, result.ResponseCodeDetail))
	}

	if err != nil || result.Error {
		return result, err
	}

	if paymentData.Card.Scheme == SchemeBoleto {
		extraInfos := createPaymentResponse.CreatePaymentResult.OrderResponse.ExtInfo
		extraInfo := extraInfos[0]
		result.BoletoUrl = extraInfo.Value
	}

	result.TransactionId = createPaymentResponse.CreatePaymentResult.PaymentResponse.TransactionId
	result.TransactionUuid = createPaymentResponse.CreatePaymentResult.PaymentResponse.TransactionUuid
	
	return result, err	

}

func (this *PayZenToolBox) CreatePaymentBoletoOnline(paymentData *api.Payment) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run CreatePaymentBoletoItau: %v", paymentData)
	}

	payloadTmp := make(map[string]string)

	payloadTmp["vads_action_mode"] = "INTERACTIVE"
	payloadTmp["vads_page_action"] = "PAYMENT"
	payloadTmp["vads_payment_config"] = "SINGLE"
	payloadTmp["vads_language"] = "pt"
	payloadTmp["vads_validation_mode"] = "0"
	payloadTmp["vads_version"] = "V2"

	payloadTmp["vads_cust_national_id"] = paymentData.Customer.IdentityCode
	payloadTmp["vads_cust_name"] = fmt.Sprintf("%v %v", paymentData.Customer.FirstName, paymentData.Customer.LastName)
	payloadTmp["vads_cust_address"] = paymentData.Customer.Address
	payloadTmp["vads_cust_address_number"] = paymentData.Customer.StreetNumber
	payloadTmp["vads_cust_state"] = paymentData.Customer.State
	payloadTmp["vads_cust_zip"] = paymentData.Customer.ZipCode
	payloadTmp["vads_cust_city"] = paymentData.Customer.City
	payloadTmp["vads_cust_district"] = paymentData.Customer.District
	payloadTmp["vads_capture_delay"] = fmt.Sprintf("%v", paymentData.BoletoOnlineDaysDalay)
	payloadTmp["vads_amount"] = formatAmount(paymentData.Amount)
	
	payloadTmp["vads_cust_email"] = paymentData.Customer.Email
	payloadTmp["vads_cust_country"] = CountryBR
	payloadTmp["vads_currency"] = fmt.Sprintf("%v", CurrencyCode)
	payloadTmp["vads_cust_phone"] = paymentData.Customer.PhoneNumber
  
  payloadTmp["vads_payment_cards"] = string(paymentData.BoletoOnline) // tipo de boleto

  if this.Account.Mode == PayZenModeProduction {
  	payloadTmp["vads_ctx_mode"] = PayZenModeProduction
	} else {
		payloadTmp["vads_ctx_mode"] = PayZenModeTest
	}

	payloadTmp["vads_ext_info_soft_descriptor"] = "Mobile Mind Soluções Tecnológicas"
	payloadTmp["vads_order_info"] = paymentData.BoletoOnlineTexto
	payloadTmp["vads_order_info2"] = paymentData.BoletoOnlineTexto2
	payloadTmp["vads_order_info3"] = paymentData.BoletoOnlineTexto3

	payloadTmp["vads_trans_date"] = util.DateNow().UTC().Format("20060102150405")

	payloadTmp["vads_site_id"] = this.Account.ShopId

	payloadTmp["vads_order_id"] = paymentData.OrderId
	payloadTmp["vads_trans_id"] = paymentData.VadsTransId


	keys := []string{}
	for k := range payloadTmp {

		if len(strings.TrimSpace(payloadTmp[k])) > 0 { 
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	payload := make(map[string]string)
	values := ""
	for _, k := range keys {
		payload[k] = payloadTmp[k]
		values = fmt.Sprintf("%v%v+", values, payloadTmp[k])
	}

	values = fmt.Sprintf("%v%v", values, this.Account.Cert)

	if this.Debug {
		fmt.Println("signature values = %v", values)
	}

	signature := this.GenerateFormAuthToken(values)

	payload["signature"] = signature

	if this.Debug {
		fmt.Println("************ HTML FORM DATA ********************** ")
		for _, k := range keys {
			val := fmt.Sprintf("<input type='hidden' name='%v' value='%v'/>", k, payload[k])
			fmt.Println(val)		
		}
		val := fmt.Sprintf("<input type='hidden' name='signature' value='%v'/>", signature)
		fmt.Println(val)		
		fmt.Println("************ HTML FORM DATA ********************** ")		
	}
	
	result, err := this.MakeFormRequest(payload)
	
	if err != nil || result.Error {
		return result, err
	}


	if len(paymentData.SaveBoletoAtPath) > 0 {
		
		result.BoletoFileName = fmt.Sprintf("Boleto_%v_%v.pdf", payload["vads_order_id"], util.DateNow().Format("20060102150405"))
		result.BoletoUrl = fmt.Sprintf("%v/%v", paymentData.SaveBoletoAtPath, result.BoletoFileName)
		

		fmt.Println("save boleto at %v, size = %v", result.BoletoUrl, len(result.BoletoOutputContent))

		if err := ioutil.WriteFile(result.BoletoUrl, result.BoletoOutputContent, 0644); err != nil {
			return result, errors.New(fmt.Sprintf("error on create local file: %v", err))
		}

		//if this.Debug {
			fmt.Println("saved boleto = %v", result.BoletoUrl)
		//}
	}


	result.TransactionId = paymentData.VadsTransId //createPaymentResponse.CreatePaymentResult.PaymentResponse.TransactionId
	result.TransactionUuid = paymentData.VadsTransId //createPaymentResponse.CreatePaymentResult.PaymentResponse.TransactionUuid
	
	result.TransactionStatus = api.WaitingAuthorisation
	result.TransactionStatusLabel = "Boleto criado"
	result.ResponseCode = "200"
	result.ResponseCodeDetail  = "Boleto criado"

	
	return result, nil	

}

/*
	update payment
*/
func (this *PayZenToolBox) UpdatePayment(paymentData *api.Payment) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run UpdatePayment: %v", paymentData)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	operation := NewSOAPUpdatePayment()
	
	operation.CommonRequest.SubmissionDate = soap.Header.Timestamp
	
	operation.PaymentRequest.Amount = formatAmount(paymentData.Amount)
	operation.PaymentRequest.Currency = CurrencyCode
	operation.PaymentRequest.ManualValidation = fmt.Sprintf("%v", paymentData.ValidationType)
		 
	operation.QueryRequest.Uuid = paymentData.TransactionUuid

	soap.Body.OperationRequest = operation
	
	soapResponse := NewSOAPResponseEnvelop()
	respOp := NewSOAPUpdatePaymentResponse()
	soapResponse.Body.UpdatePaymentResponse = respOp	

	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentNotFound = strings.Contains(result.ResponseCodeDetail, "Transaction was not found")

	switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Authorised
			break
	}

	res, erro := switchStatus(result.TransactionStatusLabel)
	
	result.Error = res.Error
	result.Message = res.Message
	result.TransactionStatus = res.TransactionStatus

	if erro != nil { // TransactionStatusLabel not found
		return result, errors.New(fmt.Sprintf("%v. Código de erro %v: %v", erro.Error(), result.ResponseCode, result.ResponseCodeDetail))
	}

	if err != nil || result.Error {
		return result, err
	}

	opResult := respOp.UpdatePaymentResult
	
	if paymentData.Card.Scheme == SchemeBoleto {
		extraInfos := opResult.OrderResponse.ExtInfo
		extraInfo := extraInfos[0]
		result.BoletoUrl = extraInfo.Value
	}

	result.TransactionId = opResult.PaymentResponse.TransactionId
	result.TransactionUuid = opResult.PaymentResponse.TransactionUuid
	
	return result, err	

}

/* cancela um pagamento */

func (this *PayZenToolBox) CancelPayment(transactionUuid string) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run CancelPayment: %v", transactionUuid)
	}

	// cria request
	soap := NewSOAPEnvelope()
	this.FillRequestHeader(soap.Header)	 
	
	operation := NewSOAPCancelPayment()
	operation.QueryRequest.Uuid = transactionUuid
	soap.Body.OperationRequest = operation  

	// cria response
	soapResponse := NewSOAPResponseEnvelop()
	respOp := NewSOAPCancelPaymentResponse()
	soapResponse.Body.CancelPaymentResponse = respOp

	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentNotFound = strings.Contains(result.ResponseCodeDetail, "Transaction was not found")

	res, erro := switchStatus(result.TransactionStatusLabel)

	if erro != nil {
		return result, errors.New(fmt.Sprintf("%v. Código de erro %v: %v", erro.Error(), result.ResponseCode, result.ResponseCodeDetail))
	}

	result.Error = res.Error
	result.Message = res.Message
	result.TransactionStatus = res.TransactionStatus

	return result, err
}

/*
	Cria uma recorrência
*/

func (this *PayZenToolBox) CreateSubscription(subscription *api.Subscription) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run CreateSubscription: %v", subscription)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	operation := NewSOAPCreateSubscription()
	
	operation.CommonRequest.SubmissionDate = soap.Header.Timestamp

	operation.OrderRequest.OrderId = subscription.OrderId


	operation.SubscriptionRequest.EffectDate = subscription.EffectDate.Format(DateTimeLayout)
	operation.SubscriptionRequest.Amount = formatAmount(subscription.Amount)
	operation.SubscriptionRequest.Currency = CurrencyCode
	operation.SubscriptionRequest.InitialAmount = formatAmount(subscription.InitialAmount)
	operation.SubscriptionRequest.InitialAmountNumber = fmt.Sprintf("%v", subscription.InitialAmountNumber)
	operation.SubscriptionRequest.SubscriptionId = subscription.SubscriptionId
	operation.SubscriptionRequest.Description = subscription.Description

	rule, err := this.buildSubscriptionRule(subscription)
	
	if err != nil {
		return nil, err
	}
	
	operation.SubscriptionRequest.Rrule = rule

	operation.CardRequest.PaymentToken = subscription.Token

			 
	soap.Body.OperationRequest = operation
	
	soapResponse := NewSOAPResponseEnvelop()
	operationResponse := NewSOAPCreateSubscriptionResponse()
	soapResponse.Body.CreateSubscriptionResponse = operationResponse

	result, err := this.MakeRequest(soap, soapResponse)

	if err != nil || result.Error {
		return result, err
	}

	switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Success
			break
	}

	result.SubscriptionInfo.SubscriptionId = operationResponse.CreateSubscriptionResult.SubscriptionResponse.SubscriptionId
	
	return result, err	

}	

func (this *PayZenToolBox) GetSubscriptionDetails(subscriptionId string, paymentToken string) (*api.PaymentResult, error){

	if this.Debug {
		fmt.Println("** run GetSubscriptionDetails: %v", subscriptionId)
	}

	// cria request
	soap := NewSOAPEnvelope()
	this.FillRequestHeader(soap.Header)	 
	
	operation := NewSOAPGetSubscriptionDetails()
	operation.QueryRequest.SubscriptionId = subscriptionId
	operation.QueryRequest.PaymentToken = paymentToken

	soap.Body.OperationRequest = operation  

	// cria response
	soapResponse := NewSOAPResponseEnvelop()
	getSubscriptionDetailsResponse := NewSOAPGetSubscriptionDetailsResponse()
	soapResponse.Body.GetSubscriptionDetailsResponse = getSubscriptionDetailsResponse

	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentTokenNotFound = strings.Contains(result.ResponseCodeDetail, "PaymentToken not found")
	result.SubscriptionIdNotFound = strings.Contains(result.ResponseCodeDetail, "SubscriptionID was not found")
	result.SubscriptionInvalid = strings.Contains(result.ResponseCodeDetail, "Invalid Subscription")


	if err != nil {
		return result, err
	}

	if result.Error {
		return result, err
	}

	switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Success
			break
	}

	getSubscriptionDetailsResult := getSubscriptionDetailsResponse.GetSubscriptionDetailsResult
	subscriptionResponse := getSubscriptionDetailsResult.SubscriptionResponse


	result.SubscriptionInfo.SubscriptionId = subscriptionResponse.SubscriptionId
	result.SubscriptionInfo.PastPaymentsNumber = subscriptionResponse.PastPaymentsNumber
	result.SubscriptionInfo.TotalPaymentsNumber = subscriptionResponse.TotalPaymentsNumber
	result.SubscriptionInfo.EffectDate, _ = util.DateParse(TimeZoneDateTimeLayout, subscriptionResponse.EffectDate)
	result.SubscriptionInfo.CancelDate, _ = util.DateParse(TimeZoneDateTimeLayout, subscriptionResponse.CancelDate) 
	result.SubscriptionInfo.InitialAmountNumber = subscriptionResponse.InitialAmountNumber
	result.SubscriptionInfo.Rule = subscriptionResponse.Rrule
	result.SubscriptionInfo.Description = subscriptionResponse.Description
	result.SubscriptionInfo.Active = result.SubscriptionInfo.CancelDate.IsZero()
	result.SubscriptionInfo.Cancelled = !result.SubscriptionInfo.CancelDate.IsZero()
	result.SubscriptionInfo.Started = result.SubscriptionInfo.EffectDate.After(util.DateNow())


	return result, err
}

func (this *PayZenToolBox) CancelSubscription(subscriptionId string, paymentToken string) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run CancelSubscription: %v", subscriptionId)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	operation := NewSOAPCancelSubscription()
	
	operation.CommonRequest.SubmissionDate = soap.Header.Timestamp
	
	operation.QueryRequest.SubscriptionId = subscriptionId
	operation.QueryRequest.PaymentToken = paymentToken

	soap.Body.OperationRequest = operation

	soapResponse := NewSOAPResponseEnvelop()
	cancelSubscriptionResponse := NewSOAPCancelSubscriptionResponse()
	soapResponse.Body.CancelSubscriptionResponse = cancelSubscriptionResponse


	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentTokenNotFound = strings.Contains(result.ResponseCodeDetail, "PaymentToken not found")
	result.SubscriptionIdNotFound = strings.Contains(result.ResponseCodeDetail, "SubscriptionID was not found")
	result.SubscriptionInvalid = strings.Contains(result.ResponseCodeDetail, "Invalid Subscription")

	switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Success
			break
	}

	return result, err
}

func (this *PayZenToolBox) UpdateSubscription(subscription *api.Subscription) (*api.PaymentResult, error) {

	if this.Debug {
		fmt.Println("** run UpdateSubscription: %v", subscription)
	}

	soap := NewSOAPEnvelope()

	this.FillRequestHeader(soap.Header)

	operation := NewSOAPUpdateSubscription()
	
	operation.CommonRequest.SubmissionDate = soap.Header.Timestamp

	operation.QueryRequest.SubscriptionId = subscription.SubscriptionId
	operation.QueryRequest.PaymentToken = subscription.Token


	operation.SubscriptionRequest.EffectDate = subscription.EffectDate.Format(DateTimeLayout)
	operation.SubscriptionRequest.Amount = formatAmount(subscription.Amount)
	operation.SubscriptionRequest.Currency = CurrencyCode
	operation.SubscriptionRequest.InitialAmount = formatAmount(subscription.InitialAmount)
	operation.SubscriptionRequest.InitialAmountNumber = fmt.Sprintf("%v", subscription.InitialAmountNumber)
	operation.SubscriptionRequest.SubscriptionId = subscription.SubscriptionId
	operation.SubscriptionRequest.Description = subscription.Description

	//operation.PaymentRequest.ManualValidation = fmt.Sprintf("%v", subscription.ValidationType)

	rule, err := this.buildSubscriptionRule(subscription)
	
	if err != nil {
		return nil, err
	}
	
	operation.SubscriptionRequest.Rrule = rule
			 
	soap.Body.OperationRequest = operation
	
	soapResponse := NewSOAPResponseEnvelop()
	operationResponse := NewSOAPUpdateSubscriptionResponse()
	soapResponse.Body.UpdateSubscriptionResponse = operationResponse

	result, err := this.MakeRequest(soap, soapResponse)

	result.PaymentTokenNotFound = strings.Contains(result.ResponseCodeDetail, "PaymentToken not found")
	result.SubscriptionIdNotFound = strings.Contains(result.ResponseCodeDetail, "SubscriptionID was not found")
	result.SubscriptionInvalid = strings.Contains(result.ResponseCodeDetail, "Invalid Subscription")

	if err != nil || result.Error {
		return result, err
	}

	result.SubscriptionInfo.Token = operationResponse.UpdateSubscriptionResult.CommonResponse.PaymentToken
	
		switch result.ResponseCodeDetail {
		case "Action successfully completed":
			result.TransactionStatus = api.Success
			break
	}

	return result, err	

}	

func (this *PayZenToolBox) buildSubscriptionRule(subscription *api.Subscription) (string, error) {
	var rule string

	if len(subscription.Rule) > 0 {
		rule = subscription.Rule
	} else {
		
		rules := []string{}
		
		if subscription.Count > 0 {
			rules = append(rules, fmt.Sprintf("COUNT=%v", subscription.Count))
		}	

		switch subscription.Cycle {
			case api.Weekly: // semanal
				rules = append(rules, "FREQ=WEEKLY")
				break
			case api.Biweekly: // quinzenal
				rules = append(rules, "FREQ=WEEKLY")
				rules = append(rules, "INTERVAL=2")
				break
			case api.Monthly: // mensal
				rules = append(rules, "FREQ=MONTHLY")
				break
			case api.Quarterly: // trimestral
				rules = append(rules, "FREQ=MONTHLY")
				rules = append(rules, "INTERVAL=4")
				break
			case api.Semiannually: // semestral
				rules = append(rules, "FREQ=MONTHLY")
				rules = append(rules, "INTERVAL=6")
				break				
			case api.Yearly:
				rules = append(rules, "FREQ=YEARLY")
				break
			default:
		    return "", errors.New("cycle is required")
		}

		if subscription.PaymentAtLastDayOfMonth && subscription.PaymentAtDayOfMonth > 0 {
	    return "", errors.New("use PaymentAtLastDayOfMonth or PaymentAtDayOfMonth")
		}

		if subscription.PaymentAtDayOfMonth > 0 {
			rules = append(rules, fmt.Sprintf("BYMONTHDAY=%v", subscription.PaymentAtDayOfMonth))
		}

		if subscription.PaymentAtLastDayOfMonth {
			rules = append(rules, "BYMONTHDAY=28,29,30,31")	
			rules = append(rules, "BYSETPOS=-1")	
		}

		if len(rules) == 0 {
	    return "", errors.New("is required")
		}

		rule = "RRULE"

		for i, it := range rules {
			if i  == 0 {
				rule = fmt.Sprintf("%v:%v", rule, it)
			} else {
				rule = fmt.Sprintf("%v;%v", rule, it)
			}
		}

	}	

	return rule, nil
}

func (this *PayZenToolBox) MakeRequest(soap *SOAPEnvelope, soapResponse *SOAPResponseEnvelop) (*api.PaymentResult, error) {

	result := api.NewPaymentResult()
	//result.RequestObject = soap

	xmlContent, err := xml.MarshalIndent(soap, "", "")

	if err != nil {
		fmt.Println("** MakeRequest: error MarshalIndent: %v", err)
		return result, errors.New(fmt.Sprintf("MakeRequest - MarshalIndent: %v", err.Error()))
	}

	result.Request = string(xmlContent)	


	data := bytes.NewBuffer(xmlContent)

	url := ""

	fmt.Println("***********************************")
	if this.Account.Mode == PayZenModeProduction {
		url = PAYZEN_API_URL_PROD_URL
		fmt.Println("**** PayZen MODE PRODUCTION")
	} else {
		url = PAYZEN_API_URL_TEST_URL
		fmt.Println("**** PayZen MODE TEST")
		fmt.Println("***********************************")
	}

	if this.Debug {
		fmt.Println("**********************************************")
		fmt.Println("%v", string(data.Bytes()))
		fmt.Println("**********************************************")
	}

	r, err := http.Post(url, "text/xml", data)

	if err != nil {
		return result, errors.New(fmt.Sprintf("http.Post: %v", err.Error()))
	}

	response, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return result, errors.New(fmt.Sprintf("ioutil.ReadAll: %v", err.Error()))		
	}

	result.Response = string(response)

	fmt.Println("***********************************")
	fmt.Println("***** PAYZEN START RESPONSE ******")	
	fmt.Println("***** STATUS CODE: %v", r.StatusCode)
	if this.Debug {
		fmt.Println("***** RESPONSE: %v", result.Response)		
	}
	fmt.Println("***** PAYZEN END RESPONSE *********")
	fmt.Println("***********************************")

	switch r.StatusCode {
		case 200:


			err := xml.Unmarshal(response, soapResponse)

			if err != nil {				
				return result, errors.New(fmt.Sprintf("MakeRequest - Unmarshal Response: %v", err.Error()))
			}

			//result.ResponseObject = soapResponse
			
			/*
			if this.Debug {
				// print xml response content
				content, err := xml.MarshalIndent(soapResponse, "  ", "    ")

				if err != nil {
					fmt.Println("** error on get response to debug log. MakeRequest - MarshalIndent: %v", err)					
				} else {				
					fmt.Println("******************************************************************")
					fmt.Println("********************* XML REPONSE CONTENT ************************")
					fmt.Println(string(content))			
					fmt.Println("********************* XML REPONSE CONTENT ************************")
					fmt.Println("******************************************************************")			
				}
			}*/

			commonResponse := soapResponse.Body.GetCommonResponse()

			responseCodeStr := commonResponse.ResponseCode
			responseCode, _ := strconv.Atoi(responseCodeStr)

			responseCodeDetail := commonResponse.ResponseCodeDetail
			transactionStatusLabel := commonResponse.TransactionStatusLabel

			result.TransactionStatusLabel = transactionStatusLabel
			result.ResponseCode = responseCodeStr
			result.ResponseCodeDetail = responseCodeDetail

			if responseCode != 0 {
				result.Error = true
				result.Message = fmt.Sprintf("Code: %v, Details: ", responseCodeStr, responseCodeDetail)
				return result, nil
			}

			if !this.ValidateResponse(soapResponse.Header) {
				result.Error = true
				result.Message = "Header not is valid"
				return result, nil
			}


			return result, nil

		default:
			return result, errors.New(fmt.Sprintf("PayZen API error - Code %v, Status: %v", r.StatusCode, r.Status))
	}
}

func (this *PayZenToolBox) MakeFormRequest(payload map[string]string) (*api.PaymentResult, error) {

	result := api.NewPaymentResult()

	form := url.Values{}

	for k, v := range payload {
		form.Add(k, v)
	}

	if this.Debug {
		fmt.Println("Form Data = %v", form)
	}

	r, err := http.PostForm(PAYZEN_FORM_URL, form)
	r.Header.Set("Accept-Language", "pt-BR")
	r.Header.Set("Content-Language", "pt-BR")

	if err != nil {
		return result, errors.New(fmt.Sprintf("http.Post: %v", err.Error()))
	}

	response, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return result, errors.New(fmt.Sprintf("ioutil.ReadAll: %v", err.Error()))		
	}

	payzenBoletoItau := payload["vads_payment_cards"] == string(api.BoletoOnlineItauBoleto) 
	payzenBoletoBradesco := payload["vads_payment_cards"] == string(api.BoletoOnlineBradescoBoleto)
	
	result.Response = string(response)

	fmt.Println("***********************************")
	fmt.Println("***** PAYZEN START RESPONSE ******")	
	fmt.Println("***** STATUS CODE: %v", r.StatusCode)
	if this.Debug {
		fmt.Println("***** RESPONSE: %v", result.Response)		
	}
	fmt.Println("***** PAYZEN END RESPONSE *********")
	fmt.Println("***********************************")


	d1 := []byte(result.Response)

	if payzenBoletoItau {		
		if err := ioutil.WriteFile(fmt.Sprintf("/tmp/payzen_boleto_result_%v_%v.html", payload["vads_order_id"], util.DateNow().Format("20060102150405")), d1, 0644); err != nil {
			fmt.Println("Error on save file %v", err)
		}
	}

	switch r.StatusCode {
		case 200:

			if payzenBoletoItau {

			  itauCode, err := getItauCode(result.Response)

				if err != nil {
					return result, errors.New(fmt.Sprintf("getItauCode - %v", err.Error()))		
				}		  

			  if itauCode == "" {
							  	
			  	boletoError, err := getBoletoError(result.Response)

					if err != nil {
						return result, errors.New(fmt.Sprintf("getBoletoError - %v", err.Error()))		
					}		  

			  	if boletoError == "" {
			  		return result, errors.New("Erro ao processar retorno PayZen: Nenhum resultado esperado foi encontrado.")		
			  	} else {
			  		return result, errors.New(boletoError)		
			  	}

			  } else {
					result.BoletoUrl = fmt.Sprintf("https://shopline.itau.com.br/shopline/impressao.aspx?DC=%v", itauCode)
			  }

			} else if payzenBoletoBradesco {
				

				if support.IsValidHtml(result.Response) {

					//fmt.Println("is valid html %v", len(d1))

					if err := ioutil.WriteFile(fmt.Sprintf("/tmp/payzen_boleto_result_%v_%v.html", payload["vads_order_id"], util.DateNow().Format("20060102150405")), d1, 0644); err != nil {
						fmt.Println("Error on save file %v", err)
					}

			  	boletoError, err := getBoletoError(result.Response)

					if err != nil {
						return result, errors.New(fmt.Sprintf("getBoletoError - %v", err.Error()))		
					}		  

			  	if boletoError == "" {
			  		return result, errors.New("Erro ao processar retorno PayZen: Nenhum resultado esperado foi encontrado.")		
			  	} else {
			  		return result, errors.New(boletoError)		
			  	}

					
				} else{					
					result.BoletoOutputContent = d1
					//fmt.Println("is not valid html")
				}

		  } 

		  result.Response = ""
			return result, nil

		default:
			return result, errors.New(fmt.Sprintf("PayZen API error - Code %v, Status: %v", r.StatusCode, r.Status))
	}
}

func getItauCode(html string) (string, error){

	doc, err := support.HtmlParse(html)

	if err != nil {
		return "", err
	}

	redirectForm := support.HtmlParseFindByName(doc, "redirectForm")


	if redirectForm == nil {
		fmt.Println("HtmlParse: redirectForm is nil")
		return "", nil
	}

	dc := support.HtmlParseFindByName(redirectForm, "DC")

	if redirectForm == nil {
		fmt.Println("HtmlParse: dc is nil")
		return "", nil
	}

	return support.GetNodeAttrValue(dc, "value"), nil	
}

func getBoletoError(html string) (string, error){

 doc, err := support.HtmlParse(html)

  if err != nil {
		return "", errors.New(fmt.Sprintf("HtmlParse: cannot parse %v", err.Error()))		
	}

	paymentResult := support.HtmlParseFindById(doc, "paymentResult")

	if paymentResult == nil {
		fmt.Println("HtmlParse: paymentResult is nil")
		return "", nil
	}	
	
	span := support.HtmlParseFindByType(paymentResult, "span")

	if span == nil {
		fmt.Println("HtmlParse: span is nil")
		return "", nil
	}	

	divParent := support.HtmlParseFindByType(span, "div")

	if divParent == nil {
		fmt.Println("HtmlParse: div1 is nil")
		return "", nil
	}	

	divChild := support.HtmlParseFindByType(divParent, "div")
	
	if divChild == nil {

		if divParent != nil {
			return divParent.FirstChild.Data, nil
		}

		fmt.Println("HtmlParse: div2 is nil")
		return "", nil
	}	

	return divChild.FirstChild.Data, nil

}


func switchStatus(transactionStatusLabel string) (*api.PaymentResult, error) {

	result := new(api.PaymentResult)

	fmt.Sprintf("switch transaction status: %v", transactionStatusLabel)

	switch transactionStatusLabel {
		case "INITIAL":
			result.TransactionStatus = api.Initial
			break
		case "NOT_CREATED":
			result.Error = true
			result.Message = "Não foi possível criar a transação"
			result.TransactionStatus = api.NotCreated
			break
		case "AUTHORISED":
			result.TransactionStatus = api.Authorised
			break
		case "AUTHORISED_TO_VALIDATE":					
			result.TransactionStatus = api.AuthorisedToValidate
			break
		case "WAITING_AUTHORISATION":
			result.TransactionStatus = api.WaitingAuthorisation
			break
		case "WAITING_AUTHORISATION_TO_VALIDATE":
			result.TransactionStatus = api.WaitingAuthorisationToValidate
			break
		case "REFUSED":
			result.Error = true
			result.Message = "A transação foi recusada, verifique os dados cartão"
			result.TransactionStatus = api.Refused
			break
		case "CAPTURED":			
			result.TransactionStatus = api.Captured
			break
		case "CANCELLED":			
			result.TransactionStatus = api.Cancelled
			break
		case "EXPIRED":		
			result.TransactionStatus = api.Expired
			break
		case "UNDER_VERIFICATION":
			result.TransactionStatus = api.UnderVerification
			break				
		default:
			//result.TransactionStatusLabel = transactionStatusLabel
			fmt.Println("Problemas na transação. Status %v não reconhecido.", transactionStatusLabel)
			return result, errors.New(fmt.Sprintf("Problemas na trasanção. Status %v não reconhecido.", transactionStatusLabel))				
	}		

	return result, nil
}

func formatAmount(amount float64) string {

	if amount == 0.0 {
		return "0"
	}

	ac := accounting.Accounting{Symbol: "", Precision: 2, Thousand: "", Decimal: ""}
	text := strings.Replace(ac.FormatMoney(amount), ",", "", -1)
	text = strings.Replace(text, ".", "", -1)	
	return text
}

func parseAmount(amount string) float64 {
	v, _ := strconv.ParseFloat(amount, 64)
	return v / 100.0
}