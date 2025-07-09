package asaas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/i18n"
	"github.com/mobilemindtec/go-payments/api"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	_ "time"
)

/*

curl --location --request GET 'https://api.chat24.io/v1/help/transports'

*/

type ResultProcessor func(data []byte, response *Response) error

const (
	AsaasProdApiUrl  = "https://api.asaas.com/v3"
	AsaasTestApiUrl  = "https://sandbox.asaas.com/api/v3"
	AccesTokenHeader = "access_token"
)

type Asaas struct {
	Debug                 bool
	EntityValidator       *validator.EntityValidator
	EntityValidatorResult *validator.EntityValidatorResult
	Lang                  string
	ValidationErrors      map[string]string
	HasValidationError    bool
	Mode                  api.AsaasMode
	AccessToken           string
}

func NewAsaas(lang string, accessToken string, mode api.AsaasMode) *Asaas {
	entityValidator := validator.NewEntityValidator(lang, "Asaas")
	entityValidatorResult := new(validator.EntityValidatorResult)
	entityValidatorResult.Errors = map[string]string{}
	return &Asaas{EntityValidator: entityValidator, Mode: mode, AccessToken: accessToken, EntityValidatorResult: entityValidatorResult}
}

func (this *Asaas) getApiUrl() string {
	if this.Mode == api.AsaasModeTest {
		return AsaasTestApiUrl
	}
	return AsaasProdApiUrl
}

func (this *Asaas) CustomerCreate(customer *Customer) (*Response, error) {

	this.Log("Call CustomerCreate")

	if !this.onValid(customer) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.post(customer, "customers", nil)
}

func (this *Asaas) CustomerUpdate(customer *Customer) (*Response, error) {

	this.Log("Call CustomerUpdate")

	if !this.onValid(customer) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.post(customer, fmt.Sprintf("customers/%v", customer.Id), nil)
}

func (this *Asaas) CustomerFindByKey(key string, value string) (*Response, error) {
	return this.CustomerFind(map[string]string{key: value})
}

func (this *Asaas) CustomerFind(filter map[string]string) (*Response, error) {

	this.Log("Call CustomerFind")

	resultProcessor := func(data []byte, response *Response) error {
		return json.Unmarshal(data, response.CustomerResults)
	}

	url := fmt.Sprintf("customers%v", this.urlQuery(filter))

	return this.get(url, resultProcessor)
}

func (this *Asaas) CustomerGet(id string) (*Response, error) {

	this.Log("Call CustomerGet")

	if len(id) == 0 {
		this.SetValidationError("id", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	resultProcessor := func(data []byte, response *Response) error {
		customer := new(Customer)
		response.CustomerResults.Data = append(response.CustomerResults.Data, customer)
		response.CustomerResults.TotalCount = 1
		return json.Unmarshal(data, customer)
	}

	return this.get(fmt.Sprintf("customers/%v", id), resultProcessor)
}

func (this *Asaas) PaymentCreate(payment *Payment) (*Response, error) {

	this.Log("Call PaymentCreate")

	if !this.onValidPayment(payment) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.post(payment, "payments", nil)
}

// /
// / Somente cobranças aguardando pagamento ou vencidas podem ser removidas.
// /
func (this *Asaas) PaymentCancel(id string) (*Response, error) {

	this.Log("Call PaymentCancel")

	if len(id) == 0 {
		this.SetValidationError("id", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.delete(fmt.Sprintf("payments/%v", id))
}

// /
// / É possível estornar cobranças via cartão de crédito recebidas ou confirmadas.
// / Ao fazer isto o saldo correspondente é debitado de sua conta no Asaas e a cobrança
// / cancelada no cartão do seu cliente.
// / O cancelamento pode levar até 10 dias úteis para aparecer na fatura de seu cliente.
// /
func (this *Asaas) PaymentRefund(id string) (*Response, error) {

	this.Log("Call PaymentRefund")

	if len(id) == 0 {
		this.SetValidationError("id", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.post(nil, fmt.Sprintf("payments/%v/refund", id), nil)
}

func (this *Asaas) PaymentReceiveInCash(payment *PaymentInCash) (*Response, error) {

	this.Log("Call PaymentInCash")

	if !this.onValid(payment) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.post(payment, fmt.Sprintf("payments/%v/receiveInCash", payment.Id), nil)
}

func (this *Asaas) PaymentGet(id string) (*Response, error) {

	this.Log("Call PaymentGet")

	if len(id) == 0 {
		this.SetValidationError("id", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.get(fmt.Sprintf("payments/%v", id), nil)
}

func (this *Asaas) PaymentFindByKey(key string, value string) (*Response, error) {
	return this.PaymentFind(map[string]string{key: value})
}

func (this *Asaas) PaymentFind(filter map[string]string) (*Response, error) {

	this.Log("Call PaymentFind")

	resultProcessor := func(data []byte, response *Response) error {
		err := json.Unmarshal(data, response.PaymentResults)

		for _, it := range response.PaymentResults.Data {
			it.BuildStatus()
		}

		return err
	}

	url := fmt.Sprintf("payments%v", this.urlQuery(filter))

	return this.get(url, resultProcessor)
}

func (this *Asaas) Payments(filter *DefaultFilter) (*Response, error) {

	this.Log("Call Payments")

	resultProcessor := func(data []byte, response *Response) error {
		err := json.Unmarshal(data, response.PaymentResults)

		for _, it := range response.PaymentResults.Data {
			it.BuildStatus()
		}

		return err
	}

	url := fmt.Sprintf("payments%v", this.urlQuery(filter.ToMap()))

	return this.get(url, resultProcessor)
}

func (this *Asaas) PaymentGetPixQrCode(id string) (*Response, error) {

	this.Log("Call PaymentGetPixQrCode")

	if len(id) == 0 {
		this.SetValidationError("id", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.get(fmt.Sprintf("payments/%v/pixQrCode", id), nil)
}

func (this *Asaas) InstallmentGet(installmentId string) (*Response, error) {

	this.Log("Call PaymentGetInstallment")

	if len(installmentId) == 0 {
		this.SetValidationError("installmentId", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.get(fmt.Sprintf("installments/%v", installmentId), nil)
}

// /
// / busca todas as parcelas de uma parcelamento
// /
func (this *Asaas) InstallmentsGet(installmentId string) (*Response, error) {

	this.Log("Call PaymentInstallments")

	if len(installmentId) == 0 {
		this.SetValidationError("installmentId", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	resultProcessor := func(data []byte, response *Response) error {
		err := json.Unmarshal(data, response.PaymentResults)

		for _, it := range response.PaymentResults.Data {
			it.BuildStatus()
		}

		return err
	}

	return this.get(fmt.Sprintf("payments?installment=%v", installmentId), resultProcessor)
}

// /
// / Somente cobranças aguardando pagamento ou vencidas podem ser removidas.
// /
func (this *Asaas) InstallmentCancel(installmentId string) (*Response, error) {

	this.Log("Call InstallmentCancel")

	if len(installmentId) == 0 {
		this.SetValidationError("installmentId", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.delete(fmt.Sprintf("installments/%v", installmentId))
}

// /
// / É possível estornar cobranças via cartão de crédito recebidas ou confirmadas.
// / Ao fazer isto o saldo correspondente é debitado de sua conta no Asaas e a cobrança
// / cancelada no cartão do seu cliente.
// / O cancelamento pode levar até 10 dias úteis para aparecer na fatura de seu cliente.
// /
func (this *Asaas) InstallmentRefund(installmentId string) (*Response, error) {

	this.Log("Call InstallmentRefund")

	if len(installmentId) == 0 {
		this.SetValidationError("installmentId", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.post(nil, fmt.Sprintf("installments/%v/refund", installmentId), nil)
}

func (this *Asaas) TokenCreate(tokenRequest *TokenRequest) (*Response, error) {

	this.Log("Call TokenCreate")

	if !this.onValidToquenRequest(tokenRequest) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	resultProcessor := func(data []byte, response *Response) error {
		response.Card = new(CardResponse)
		err := json.Unmarshal(data, response.Card)
		return err
	}

	result, err := this.post(tokenRequest, "creditCard/tokenizeCreditCard", resultProcessor)

	if err == nil && !result.Error {
		result.Status = api.AsaasSuccess
	}

	return result, err
}

func (this *Asaas) SubscriptionCreate(payment *Payment) (*Response, error) {

	this.Log("Call SubscriptionCreate")

	if !this.onValidPayment(payment) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.post(payment, "subscriptions", nil)
}

func (this *Asaas) SubscriptionUpdate(payment *Payment) (*Response, error) {

	this.Log("Call SubscriptionUpdate")

	if len(payment.Id) == 0 {
		this.SetValidationError("Id", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	//if !this.onValidPayment(payment) {
	//  return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	//}

	return this.post(payment, fmt.Sprintf("subscriptions/%v", payment.Id), nil)
}

/*
Sobre a atualização, você consegue atualizar o cartão de uma assinatura por exemplo:
Na assinatura, quando o cliente paga com cartão de crédito, o cartão dele é automaticamente cadastrado para ser usado na recorrência.
Caso o cliente queira informar outro cartão, você precisará recuperar uma cobrança dessa assinatura que ainda não tenha sido paga (confirmada).
Para listar as cobranças de uma assinatura e os status, use as instruções desse trecho de nosso manual -> https://asaasv3.docs.apiary.io/#reference/0/assinaturas/listar-cobrancas-de-uma-assinatura
Após recuperar o ID da cobrança, você precisará fazer uma chamada adicional no seguinte endpoint, passando o ID da cobrança no lugar do {id_cobranca}:
/api/v3/payments/{id_cobrança}/payWithCreditCard"    ["POST"]

	{
	   "creditCard":{
	      "holderName":"marcelo h almeida",
	      "number":"5162306219378829",
	      "expiryMonth":"05",
	      "expiryYear":"2021",
	      "ccv":"318"
	   },
	   "creditCardHolderInfo":{
	      "name":"Marcelo Henrique Almeida",
	      "email":"marcelo.almeida@gmail.com",
	      "cpfCnpj":"24971563792",
	      "postalCode":"89223-005",
	      "addressNumber":"277",
	      "addressComplement":null,
	      "phone":"4738010919",
	      "mobilePhone":"47998781877"
	   }
	}

Ou, pode ser enviado apenas o Token caso você tenha tokenização habilitada e possua o token do cartão a ser utilizado, dessa forma:

	{
	    "creditCardToken": "461f086a-e2ff-426e-b1ab-22f9118a07e8"
	}
*/
func (this *Asaas) SubscriptionUpdateCardToken(payment *Payment) (*Response, error) {

	this.Log("Call SubscriptionUpdate")

	if len(payment.Id) == 0 {
		this.SetValidationError("Id", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	data := make(map[string]interface{})

	if len(payment.CardToken) > 0 {
		data["creditCardToken"] = payment.CardToken
	} else {
		if !this.onValidCard(payment) {
			return nil, errors.New(this.getMessage("Asaas.ValidationError"))
		}

		data["creditCard"] = payment.Card
		data["creditCardHolderInfo"] = payment.CardHolderInfo
	}

	return this.post(data, fmt.Sprintf("payments/%v/payWithCreditCard", payment.Id), nil)
}

func (this *Asaas) SubscriptionCancel(subscriptionId string) (*Response, error) {

	this.Log("Call PaymentCancel")

	if len(subscriptionId) == 0 {
		this.SetValidationError("subscriptionId", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.delete(fmt.Sprintf("subscriptions/%v", subscriptionId))
}

func (this *Asaas) SubscriptionGet(subscriptionId string) (*Response, error) {

	this.Log("Call SubscriptionGet")

	if len(subscriptionId) == 0 {
		this.SetValidationError("subscriptionId", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.get(fmt.Sprintf("subscriptions/%v", subscriptionId), nil)
}

func (this *Asaas) SubscriptionPaymentsGet(subscriptionId string) (*Response, error) {

	this.Log("Call SubscriptionPaymentsGet")

	if len(subscriptionId) == 0 {
		this.SetValidationError("subscriptionId", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	resultProcessor := func(data []byte, response *Response) error {
		err := json.Unmarshal(data, response.PaymentResults)

		for _, it := range response.PaymentResults.Data {
			it.BuildStatus()
		}

		return err
	}

	return this.get(fmt.Sprintf("subscriptions/%v/payments", subscriptionId), resultProcessor)
}

func (this *Asaas) SubscriptionFindByKey(key string, value string) (*Response, error) {
	return this.SubscriptionFind(map[string]string{key: value})
}

func (this *Asaas) SubscriptionFind(filter map[string]string) (*Response, error) {

	this.Log("Call SubscriptionFind")

	resultProcessor := func(data []byte, response *Response) error {
		err := json.Unmarshal(data, response.PaymentResults)

		for _, it := range response.PaymentResults.Data {
			it.BuildStatus()
		}

		return err
	}

	url := fmt.Sprintf("subscriptions%v", this.urlQuery(filter))

	return this.get(url, resultProcessor)
}

func (this *Asaas) PaymentLinkCreate(payment *Payment) (*Response, error) {

	this.Log("Call PaymentLinkCreate")

	if !this.onValidPayment(payment) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.post(payment, "paymentLinks", nil)
}

func (this *Asaas) PaymentLinkCancel(peymentLinkId string) (*Response, error) {

	this.Log("Call PaymentLinkCancel")

	if len(peymentLinkId) == 0 {
		this.SetValidationError("peymentLinkId", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.delete(fmt.Sprintf("paymentLinks/%v", peymentLinkId))
}

func (this *Asaas) PaymentLinkGet(peymentLinkId string) (*Response, error) {

	this.Log("Call PaymentLinkGet")

	if len(peymentLinkId) == 0 {
		this.SetValidationError("peymentLinkId", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.get(fmt.Sprintf("paymentLinks/%v", peymentLinkId), nil)
}

func (this *Asaas) PaymentLinkFindByKey(key string, value string) (*Response, error) {
	return this.PaymentLinkFind(map[string]string{key: value})
}

func (this *Asaas) PaymentLinkFind(filter map[string]string) (*Response, error) {

	this.Log("Call PaymentLinkFind")

	resultProcessor := func(data []byte, response *Response) error {
		err := json.Unmarshal(data, response.PaymentResults)

		for _, it := range response.PaymentResults.Data {
			it.BuildStatus()
		}

		return err
	}

	url := fmt.Sprintf("paymentLinks%v", this.urlQuery(filter))

	return this.get(url, resultProcessor)
}

func (this *Asaas) FinancialTransactionsList(filter *DefaultFilter) (*Response, error) {

	this.Log("Call FinancialTransactionsList")

	resultProcessor := func(data []byte, response *Response) error {
		return json.Unmarshal(data, response.FinancialTransactionResults)
	}

	url := fmt.Sprintf("financialTransactions%v", this.urlQuery(filter.ToMap()))
	return this.get(url, resultProcessor)
}

func (this *Asaas) CurrentBalance() (*Response, error) {
	this.Log("Call CurrentBalance")
	return this.get("finance/getCurrentBalance", nil)
}

func (this *Asaas) TransferCreate(transfer *Transfer) (*Response, error) {
	this.Log("Call TransferCreate")

	if !this.onValidTransfer(transfer) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.post(transfer, "transfers", nil)
}

func (this *Asaas) TransferGet(id string) (*Response, error) {
	this.Log("Call TransferCreate")

	if len(id) == 0 {
		this.SetValidationError("id", "is required")
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	resultProcessor := func(data []byte, response *Response) error {

		transferResult := new(TransferResult)

		err := json.Unmarshal(data, transferResult)

		if err != nil {
			return err
		}

		response.TransferResults.Data = append(response.TransferResults.Data, transferResult)
		return nil
	}

	return this.get(fmt.Sprintf("transfers/%v", id), resultProcessor)
}

func (this *Asaas) TransferList(filter *DefaultFilter) (*Response, error) {
	this.Log("Call TransferList")

	resultProcessor := func(data []byte, response *Response) error {
		return json.Unmarshal(data, response.TransferResults)
	}

	url := fmt.Sprintf("transfers%v", this.urlQuery(filter.ToMap()))
	return this.get(url, resultProcessor)
}

func (this *Asaas) AccountCreate(account *Account) (*Response, error) {
	this.Log("Call AccountCreate")

	if !this.onValidAccount(account, true) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	resultProcessor := func(data []byte, response *Response) error {
		accountResult := NewAccount(nil)
		err := json.Unmarshal(data, accountResult)

		if err == nil {
			response.AccountResults.Data = append(response.AccountResults.Data, accountResult)
		}

		return err
	}

	return this.post(account, "accounts", resultProcessor)
}

func (this *Asaas) BankAccountMainCreateOrUpdate(bankAccount *BankAccountSimple) (*Response, error) {
	this.Log("Call BankAccountCreate")

	if !this.onValidBankAccount(bankAccount) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	return this.post(bankAccount, "bankAccounts/mainAccount")
}

func (this *Asaas) AccountUpdate(account *Account) (*Response, error) {
	this.Log("Call AccountUpdate")

	if !this.onValidAccount(account, false) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	resultProcessor := func(data []byte, response *Response) error {
		accountResult := NewAccount(nil)
		err := json.Unmarshal(data, accountResult)

		if err == nil {
			response.AccountResults.Data = append(response.AccountResults.Data, accountResult)
		}

		return err
	}

	return this.post(account, "myAccount/commercialInfo", resultProcessor)
}

func (this *Asaas) AccountGet() (*Response, error) {
	this.Log("Call AccountGet")

	resultProcessor := func(data []byte, response *Response) error {
		accountResult := NewAccount(nil)
		err := json.Unmarshal(data, accountResult)

		if err == nil {
			response.AccountResults.Data = append(response.AccountResults.Data, accountResult)
		}

		return err
	}

	return this.get("myAccount/commercialInfo", resultProcessor)
}

func (this *Asaas) AccountDocuments() (*Response, error) {
	this.Log("Call AccountDocuments")

	resultProcessor := func(data []byte, response *Response) error {
		docs := new(Documents)
		err := json.Unmarshal(data, docs)

		if err == nil {
			response.Documents = docs
		}

		return err
	}

	return this.get("myAccount/documents", resultProcessor)
}

func (this *Asaas) AccountDocumentSend(doc *Document) (*Response, error) {
	this.Log("Call AccountDocumentSend")

	resultProcessor := func(data []byte, response *Response) error {
		docResp := new(DocumentResponse)
		err := json.Unmarshal(data, docResp)

		if err == nil {
			response.DocumentResponse = docResp
		}

		return err
	}

	data := map[string]interface{}{}
	data["documentFile"] = doc.DocumentFile
	data["type"] = string(doc.Type)

	return this.postMultiPart(data, fmt.Sprintf("myAccount/documents/%v", doc.Id), resultProcessor)
}

func (this *Asaas) AccountStatus() (*Response, error) {
	this.Log("Call AccountStatus")

	resultProcessor := func(data []byte, response *Response) error {
		status := &AccountStatus{}
		err := json.Unmarshal(data, status)

		if err == nil {
			response.AccountStatus = status
		}

		return err
	}

	return this.get("myAccount/status", resultProcessor)
}

func (this *Asaas) AccountList() (*Response, error) {
	this.Log("Call AccountList")

	resultProcessor := func(data []byte, response *Response) error {
		return json.Unmarshal(data, response.AccountResults)
	}

	return this.get("accounts", resultProcessor)
}

func (this *Asaas) WebhookCreateOrChange(webhook *WebhookObject) (*Response, error) {
	this.Log("Call WebhookCreateOrChange")

	if !this.onValid(webhook) {
		return nil, errors.New(this.getMessage("Asaas.ValidationError"))
	}

	resultProcessor := func(data []byte, response *Response) error {
		response.Webhook = NewWebhookObject()
		return json.Unmarshal(data, response.Webhook)
	}

	uri := "webhook"

	switch webhook.Type {
	case WebhookInvoice:
		uri = fmt.Sprintf("%v/invoice", uri)
		break
	case WebhookTransfer:
		uri = fmt.Sprintf("%v/transfer", uri)
		break
	case WebhookBill:
		uri = fmt.Sprintf("%v/bill", uri)
		break
	case WebhookAnticipation:
		uri = fmt.Sprintf("%v/anticipation", uri)
		break
	case WebhookMobilePhoneRecharge:
		uri = fmt.Sprintf("%v/mobilePhoneRecharge", uri)
		break
	case WebhookAccountStatus:
		uri = fmt.Sprintf("%v/accountStatus", uri)
		break
	case WebhookPayment:
		// default
		break
	default:
		return nil, errors.New("unknown webhook type")
	}

	return this.post(webhook, uri, resultProcessor)
}

func (this *Asaas) WebhookStatus(webhookType WebhookType) (*Response, error) {
	this.Log("Call WebhookStatus")

	resultProcessor := func(data []byte, response *Response) error {
		response.Webhook = NewWebhookObject()
		return json.Unmarshal(data, response.Webhook)
	}

	uri := "webhook"

	switch webhookType {
	case WebhookInvoice:
		uri = fmt.Sprintf("%v/invoice", uri)
		break
	case WebhookTransfer:
		uri = fmt.Sprintf("%v/transfer", uri)
		break
	case WebhookBill:
		uri = fmt.Sprintf("%v/bill", uri)
		break
	case WebhookAnticipation:
		uri = fmt.Sprintf("%v/anticipation", uri)
		break
	case WebhookMobilePhoneRecharge:
		uri = fmt.Sprintf("%v/mobilePhoneRecharge", uri)
		break
	case WebhookAccountStatus:
		uri = fmt.Sprintf("%v/accountStatus", uri)
		break
	case WebhookPayment:
		// default
		break
	default:
		return nil, errors.New("unknown webhook type")
	}

	return this.get(uri, resultProcessor)
}

func (this *Asaas) Wallets() (*Response, error) {
	this.Log("Call Wallets")

	resultProcessor := func(data []byte, response *Response) error {
		return json.Unmarshal(data, response.WalletResults)
	}

	return this.get("wallets", resultProcessor)
}

func (this *Asaas) get(action string, resultProcessor ResultProcessor) (*Response, error) {
	return this.request(nil, action, "GET", resultProcessor, false)
}

func (this *Asaas) delete(action string) (*Response, error) {
	return this.request(nil, action, "DELETE", nil, false)
}

func (this *Asaas) post(data interface{}, action string, resultProcessor ...ResultProcessor) (*Response, error) {

	var processor ResultProcessor

	if len(resultProcessor) > 0 {
		processor = resultProcessor[0]
	}

	return this.request(data, action, "POST", processor, false)
}

func (this *Asaas) postMultiPart(data map[string]interface{}, action string, resultProcessor ...ResultProcessor) (*Response, error) {

	var processor ResultProcessor

	if len(resultProcessor) > 0 {
		processor = resultProcessor[0]
	}

	return this.request(data, action, "POST", processor, true)
}

func (this *Asaas) put(data interface{}, action string, resultProcessor ResultProcessor) (*Response, error) {
	return this.request(data, action, "PUT", resultProcessor, false)
}

func (this *Asaas) request(
	data interface{}, action string, method string, resultProcessor ResultProcessor, multiPartFormData bool) (*Response, error) {

	result := NewResponse()

	var req *http.Request
	var err error

	client := new(http.Client)
	apiUrl := fmt.Sprintf("%v/%v", this.getApiUrl(), action)

	this.Log("URL %v, METHOD = %v", apiUrl, method)

	if method == "POST" && data != nil {

		if !multiPartFormData {
			payload, err := json.Marshal(data)

			if err != nil {
				logs.Debug("error json.Marshal ", err.Error())
				return result, err
			}

			postData := bytes.NewBuffer(payload)

			result.Request = string(payload)

			if this.Debug {
				logs.Debug("****************** Asaas Request ******************")
				pettry, _ := json.MarshalIndent(data, "", "  ")
				logs.Debug(string(pettry))
				logs.Debug("****************** Asaas Request ******************")
			}

			req, err = http.NewRequest(method, apiUrl, postData)
		} else {
			var formData bytes.Buffer
			writer := multipart.NewWriter(&formData)
			m := data.(map[string]interface{})
			for k, v := range m {

				if x, ok := v.(io.Closer); ok {
					defer x.Close()
				}

				if val, ok := v.(string); ok {
					fw, err := writer.CreateFormField(k)

					if err != nil {
						return result, err
					}

					if _, err := fw.Write([]byte(val)); err != nil {
						return result, err
					}
				} else if bits, ok := v.([]byte); ok {
					fw, err := writer.CreateFormField(k)

					if err != nil {
						return result, err
					}

					if _, err := fw.Write(bits); err != nil {
						return result, err
					}
				} else if reader, ok := v.(io.Reader); ok {
					fw, err := writer.CreateFormField(k)

					if err != nil {
						return result, err
					}

					if _, err := io.Copy(fw, reader); err != nil {
						return result, err
					}
				} else if fd, ok := v.(*os.File); ok {

					fw, err := writer.CreateFormFile(k, fd.Name())

					if err != nil {
						return result, err
					}

					if _, err := io.Copy(fw, fd); err != nil {
						return result, err
					}

				} else {
					return result, fmt.Errorf("invalid field type: %v", reflect.TypeOf(v))
				}
			}
			if err := writer.Close(); err != nil {
				return result, err
			}

			/*if this.Debug {
				logs.Debug("****************** Asaas Request ******************")
				logs.Debug(string(formData.Bytes()))
				logs.Debug("****************** Asaas Request ******************")
			}*/

			req, err = http.NewRequest(method, apiUrl, &formData)
		}

	} else {
		req, err = http.NewRequest(method, apiUrl, nil)
	}

	if err != nil {
		logs.Debug("err = ", err)
		return nil, errors.New(fmt.Sprintf("error on http.NewRequest: %v", err))
	}

	if multiPartFormData {
		req.Header.Add("Content-Type", "multipart/form-data")
	} else {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add(AccesTokenHeader, this.AccessToken)

	res, err := client.Do(req)

	if err != nil {
		logs.Debug("err = ", err)
		return nil, errors.New(fmt.Sprintf("error on client.Do: %v", err))
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		logs.Debug("err = ", err)
		return nil, errors.New(fmt.Sprintf("error on ioutil.ReadAll: %v", err))
	}

	result.Response = string(body)

	if this.Debug {
		logs.Debug("****************** Asaas Response ", res.StatusCode, " ******************")
		logs.Debug(result.Response)
		logs.Debug("****************** Asaas Response ******************")
	}

	if res.StatusCode == 200 || res.StatusCode == 400 {
		if res.StatusCode == 200 && resultProcessor != nil {
			if err := resultProcessor(body, result); err != nil {
				logs.Debug("err =", err)
				return nil, errors.New(fmt.Sprintf("error on resultProcessor: %v", err))
			}
		} else {

			err = json.Unmarshal(body, result)

			if err != nil {
				logs.Debug("err = ", err)
				return nil, errors.New(fmt.Sprintf("error on json.Unmarshal: %v", err))
			}

		}
	}

	if res.StatusCode == 400 {
		result.Error = true

		if result.ErrorsCount() > 0 {

			for _, erro := range result.Errors {
				result.Message = fmt.Sprintf("%v%v, ", result.Message, erro.Description)
			}

		} else {
			result.Message = fmt.Sprintf("Asaas validation errror")
		}

		return result, nil
	}

	if res.StatusCode != 200 {
		result.Error = true
		result.Message = fmt.Sprintf("Assas error. Status: %v", res.StatusCode)
		return result, errors.New(result.Message)
	}

	result.Error = result.HasError()
	result.BuildStatus()

	return result, nil
}

func (this *Asaas) onValid(entity interface{}) bool {
	this.EntityValidatorResult, _ = this.EntityValidator.IsValid(entity, nil)

	if this.EntityValidatorResult.HasError {
		this.onValidationErrors()
		return false
	}

	return true
}

func (this *Asaas) onValidBankAccount(bankAccount *BankAccountSimple) bool {
	this.EntityValidatorResult, _ = this.EntityValidator.ValidSimple(bankAccount)

	if this.EntityValidatorResult.HasError {
		this.onValidationErrors()
		return false
	}

	return true
}
func (this *Asaas) onValidAccount(account *Account, validateBankAccount bool) bool {

	items := []interface{}{
		account,
		account.BankAccount,
	}

	//if account.BankAccount != nil {
	//  items = append(items, account.BankAccount.Bank)
	//}

	this.EntityValidatorResult, _ = this.EntityValidator.ValidMult(items, func(validator *validation.Validation) {

		if account.Webhooks == nil || len(account.Webhooks) == 0 {
			validator.SetError("Webhooks", this.getMessage("Asaas.required"))
		}

		if validateBankAccount {
			if account.BankAccount == nil {
				validator.SetError("BankAccount", this.getMessage("Asaas.required"))
				validator.SetError("Bank", this.getMessage("Asaas.required"))
			}

			if account.BankAccount != nil && len(account.BankAccount.Bank) == 0 {
				validator.SetError("Bank", this.getMessage("Asaas.required"))
			}

			if len(account.LoginEmail) == 0 {
				validator.SetError("LoginEmail", this.getMessage("Asaas.required"))
			}

		} else {
			if len(account.PersonType) == 0 {
				validator.SetError("PersonType", this.getMessage("Asaas.required"))
			}
		}

	})

	if this.EntityValidatorResult.HasError {
		this.onValidationErrors()
		return false
	}

	return true
}

func (this *Asaas) onValidTransfer(transfer *Transfer) bool {

	items := []interface{}{
		transfer,
		transfer.BankAccount,
	}

	if transfer.BankAccount != nil {
		items = append(items, transfer.BankAccount.Bank)
	}

	this.EntityValidatorResult, _ = this.EntityValidator.ValidMult(items, func(validator *validation.Validation) {

		if transfer.BankAccount == nil {
			validator.SetError("BankAccount", this.getMessage("Asaas.required"))
			validator.SetError("Bank", this.getMessage("Asaas.required"))
		}

		if transfer.BankAccount != nil && transfer.BankAccount.Bank == nil {
			validator.SetError("Bank", this.getMessage("Asaas.required"))
		}

	})

	if this.EntityValidatorResult.HasError {
		this.onValidationErrors()
		return false
	}

	return true
}

func (this *Asaas) onValidToquenRequest(tokenRequest *TokenRequest) bool {
	this.EntityValidatorResult, _ = this.EntityValidator.Valid(tokenRequest, nil)

	if this.EntityValidatorResult.HasError {
		this.onValidationErrors()
		return false
	}

	return true
}

func (this *Asaas) onValidCard(payment *Payment) bool {

	items := []interface{}{
		payment.Card,
		payment.CardHolderInfo,
	}

	this.EntityValidatorResult, _ = this.EntityValidator.ValidMult(items, func(validator *validation.Validation) {

		if payment.Card == nil {
			validator.SetError("Card", this.getMessage("Asaas.required"))
		}

		if payment.CardHolderInfo == nil {
			validator.SetError("CardHolderInfo", this.getMessage("Asaas.required"))
		}

	})

	if this.EntityValidatorResult.HasError {
		this.onValidationErrors()
		return false
	}

	return true
}

func (this *Asaas) onValidPayment(payment *Payment) bool {

	validCard := payment.BillingType == BillingCreditCard && len(payment.CardToken) == 0

	items := []interface{}{
		payment,
	}

	if payment.Splits != nil {
		for _, it := range payment.Splits {
			items = append(items, it)
		}
	}

	if validCard {
		items = append(items, payment.Card)
		items = append(items, payment.CardHolderInfo)
	}

	this.EntityValidatorResult, _ = this.EntityValidator.ValidMult(items, func(validator *validation.Validation) {

		if payment.PaymentType == PaymentLink {

			if payment.ChargeType == ChargeTypeNone {
				validator.SetError("ChargeType", this.getMessage("Asaas.required"))
			}

			if len(payment.Name) == 0 {
				validator.SetError("Name", this.getMessage("Asaas.required"))
			}

			if payment.ChargeType == Recurrent {
				if payment.SubscriptionCycle == api.SubscriptionCycleNone {
					validator.SetError("SubscriptionCycle", this.getMessage("Asaas.required"))
				}
			} else if payment.ChargeType == Installment {
				if payment.MaxInstallmentCount <= 0 {
					validator.SetError("MaxInstallmentCount", this.getMessage("Asaas.required"))
				}
			}

			if payment.BillingType == BillingBoleto || payment.BillingType == BillingUndefined {
				if payment.DueDateLimitDays <= 0 {
					validator.SetError("DueDateLimitDays", this.getMessage("Asaas.required"))
				}
			}

		} else {

			if len(payment.Customer) == 0 {
				validator.SetError("Customer", this.getMessage("Asaas.required"))
			}

			if len(payment.ExternalReference) == 0 {
				validator.SetError("ExternalReference", this.getMessage("Asaas.required"))
			}

			if payment.PaymentType == PaymentDefault {

				if len(payment.DueDate) == 0 {
					validator.SetError("DueDate", this.getMessage("Asaas.required"))
				}

			} else {

				if len(payment.NextDueDate) == 0 {
					validator.SetError("NextDueDate", this.getMessage("Asaas.required"))
				}

				if payment.SubscriptionCycle == api.SubscriptionCycleNone {
					validator.SetError("SubscriptionCycle", this.getMessage("Asaas.required"))
				}

			}

			if payment.BillingType == BillingCreditCard {

				if validCard {
					if len(payment.RemoteIp) == 0 {
						validator.SetError("RemoteIp", this.getMessage("Asaas.required"))
					}
				} else {
					if payment.Card != nil {
						validator.SetError("Card", this.getMessage("Asaas.shouldNil"))
					}

					if payment.CardHolderInfo != nil {
						validator.SetError("CardHolderInfo", this.getMessage("Asaas.shouldNil"))
					}
				}
			}

			if payment.Splits != nil {
				for i, it := range payment.Splits {
					if it.FixedValue <= 0 && it.PercentualValue <= 0 {
						validator.SetError(fmt.Sprintf("Split.%v", i), "Set fixed valur or percentual value")
					}
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

func (this *Asaas) getMessage(key string, args ...interface{}) string {
	return i18n.Tr(this.Lang, key, args)
}

func (this *Asaas) onValidationErrors() {
	this.HasValidationError = true
	this.ValidationErrors = this.EntityValidator.GetValidationErrors(this.EntityValidatorResult)
}

func (this *Asaas) SetValidationError(key string, value string) {
	this.HasValidationError = true
	if this.ValidationErrors == nil {
		this.ValidationErrors = make(map[string]string)
	}
	this.ValidationErrors[key] = value
}

func (this *Asaas) Log(message string, args ...interface{}) {
	if this.Debug {
		logs.Debug("Assas: ", fmt.Sprintf(message, args...))
	}
}

func (this *Asaas) urlQuery(filter map[string]string) string {
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
