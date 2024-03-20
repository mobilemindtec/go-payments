package v1

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/leekchan/accounting"
	"github.com/mobilemindtec/go-payments/api"
)

type BoletoRule string

const (
	BoletoStrictExpirationDate BoletoRule = "strict_expiration_date" // restringe o pagamento para até a data de vencimento e apenas o valor exato do documento
	BoletoNoStrict             BoletoRule = "no_strict"              // permite pagamento após o vencimento e valores diferentes do impresso
)

type SubscriptionCycle string

const (
	SubscriptionCycleNone SubscriptionCycle = ""
	Weekly                SubscriptionCycle = "WEEKLY"       // semanal
	Biweekly              SubscriptionCycle = "BIWEEKLY"     // quinzenal
	Monthly               SubscriptionCycle = "MONTHLY"      // mensal
	Quarterly             SubscriptionCycle = "QUARTERLY"    // trimestral
	Semiannually          SubscriptionCycle = "SEMIANNUALLY" // semanal
	Yearly                SubscriptionCycle = "YEARLY"       // anual
)

type Filter struct {
	Count     int64 // max 1000
	Page      int64
	Status    string
	StartDate string
	EndDate   string
	ApiKey    string

	BankAccountId string
}

func NewFilter() *Filter {
	return &Filter{}
}

func (this *Filter) SetStartDate(date time.Time) {
	this.StartDate = date.Format("2006-01-02")
}

func (this *Filter) SetEndDate(date time.Time) {
	this.EndDate = date.Format("2006-01-02")
}

func (this *Filter) ToMap() map[string]string {
	data := make(map[string]string)

	if this.Count > 0 {
		data["count"] = fmt.Sprintf("%v", this.Count)
	}

	if this.Page > 0 {
		data["page"] = fmt.Sprintf("%v", this.Page)
	}

	if len(this.StartDate) > 0 {
		data["start_date"] = this.StartDate
	}

	if len(this.EndDate) > 0 {
		data["end_date"] = this.EndDate
	}

	if len(this.Status) > 0 {
		data["status"] = this.Status
	}

	if len(this.ApiKey) > 0 {
		data["api_key"] = this.ApiKey
	}

	if len(this.BankAccountId) > 0 {
		data["bank_account_id"] = this.BankAccountId
	}

	return data
}

type ResponseErrorItem struct {
	Message       string `json:"message"`
	ParameterName string `json:"parameter_name"`
	Type          string `json:"type"`
}

type ResponseError struct {
	Method string               `json:"method"`
	Url    string               `json:"url"`
	Errors []*ResponseErrorItem `json:"errors"`
}

func NewResponseError() *ResponseError {
	return &ResponseError{Errors: []*ResponseErrorItem{}}
}

type Address struct {
	Neighborhood string ` json:"neighborhood" valid:"Required" `
	Street       string `json:"street" valid:"Required" `
	StreetNumber string `json:"street_number" valid:"Required" `
	ZipCode      string `json:"zipcode" valid:"Required" `
	City         string `json:"city" valid:"Required" `
	State        string `json:"state" valid:"Required" `
}

type Phone struct {
	Ddd    string `json:"ddd" valid:"Required;MaxSize(2)" `
	Number string `json:"number" valid:"Required;MaxSize(9);MinSize(9)" `
}

type CustomerDocument struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

type Customer struct {
	DocumentNumber string   `json:"document_number" valid:"Required;MinSize(11);MaxSize(14)"`
	Email          string   `json:"email" valid:"Required;Email"`
	Name           string   `json:"name" valid:"Required"`
	Address        *Address `json:"address" valid:"Required"`
	Phone          *Phone   `json:"phone" valid:"Required"`
	ApiKey         string   `json:"api_key,omitempty" valid:""`
	Id             int64    `json:"id,omitempty"` //
}

func NewCustomer() *Customer {
	entity := new(Customer)
	return entity
}

type BoletoFine struct {
	Days       int64 `json:"days"`
	Amount     int64 `json:"amount"`
	Percentage int64 `json:"percentage"`
}

func NewBoletoFine() *BoletoFine {
	return &BoletoFine{}
}

type BoletoInterest struct {
	Days       int64 `json:"days"`
	Amount     int64 `json:"amount"`
	Percentage int64 `json:"percentage"`
}

func NewBoletoInterest() *BoletoInterest {
	return &BoletoInterest{}
}

/*

	Exemplo de uma transação de R$ 100, onde 99 vai para o cliente e 1 real vai para a mobile mind

	cliente := new(SplitRule)
	cliente.Liable = true
	cliente.ChargeProcessingFee = true
	//cliente.Percentage = 100 // apenas no caso de percentual
	cliente.ChargeRemainderFee = true
	cliente.RecipientId = id do recebedor no pagarme
	cliente.Amount = 99 // 99 reais

	mobilemind := new(SplitRule)
	mobilemind.Liable = false
	mobilemind.ChargeProcessingFee = false
	//mobilemind.Percentage = 0 // apenas no caso de percentual
	mobilemind.ChargeRemainderFee = false
	mobilemind.RecipientId = id do recebedor mobile mind no pagarme
	mobilemind.Amount = 1 // 1 real


*/

type SplitRule struct {
	Liable              bool   `json:"liable"`                // Se o recebedor é responsável ou não pelo chargeback. Default true para todos os recebedores da transação.
	ChargeProcessingFee bool   `json:"charge_processing_fee"` // Se o recebedor será cobrado das taxas da criação da transação. Default true para todos os recebedores da transação.
	Percentage          int64  `json:"percentage"`            // Qual a porcentagem que o recebedor receberá. Deve estar entre 0 e 100. Se amount já está preenchido, não é obrigatório
	Amount              int64  `json:"amount"`                // Qual o valor da transação o recebedor receberá. Se percentage já está preenchido, não é obrigatório
	ChargeRemainderFee  bool   `json:"charge_remainder_fee"`  //Se o recebedor deverá pagar os eventuais restos das taxas, calculadas em porcentagem. Sendo que o default vai para o primeiro recebedor definido na regra.
	RecipientId         string `json:"recipient_id"`          // Id do recebedor
}

type Plano struct {
	Id             int64             `json:"id,omitempty"`
	Amount         int64             `json:"amount" valid:"Required"`
	Days           int64             `json:"days,omitempty"` // Prazo, em dias, para cobrança das parcelas
	Name           string            `json:"name" valid:"Required"`
	TrialDays      int64             `json:"trial_days"`
	PaymentMethods []api.PaymentType `json:"payment_methods"`
	//Número de cobranças que poderão ser feitas nesse plano.
	// Ex: Plano cobrado 1x por ano, válido por no máximo 3 anos.
	// Nesse caso, nossos parâmetros serão: days = 365, installments = 1,
	// charges=2 (cartão de crédito) ou charges=3 (boleto).
	// OBS: No caso de cartão de crédito, a cobrança feita na ativação da assinatura não é considerada.
	// OBS: null irá cobrar o usuário indefinidamente, ou até o plano ser cancelado
	Charges int64 `json:"charges,omitempty"`
	//Número de parcelas entre cada cobrança (charges).
	// Ex: Plano anual, válido por 2 anos, sendo que cada transação será dividida em 12 vezes.
	// Nesse caso, nossos parâmetros serão: days = 365, installments = 12, charges=2 (cartão de crédito)
	// ou charges=3 (boleto). OBS: Boleto sempre terá installments = 1
	// (parcelamento de cada item da recorência, padrão 1x)
	Installments int64 `json:"installments"`
	// Define em até quantos dias antes o cliente será avisado sobre o vencimento do boleto.
	InvoiceReminder int64  `json:"invoice_reminder,omitempty"`
	ApiKey          string `json:"api_key" valid:"Required"`
}

func NewPlano(name string, amount float64) *Plano {
	return &Plano{
		Amount:       FormatAmount(amount),
		Name:         name,
		Installments: 1,
	}
}

func (this *Plano) SetCycle(cycle SubscriptionCycle, charges int64, invoiceReminder int64) {

	this.Charges = charges
	this.InvoiceReminder = invoiceReminder

	switch cycle {
	case Weekly: // semanal
		this.Days = 7
		break
	case Biweekly: // quinzenal
		this.Days = 14
		break
	case Monthly: // mensal
		this.Days = 30
		break
	case Quarterly: // trimestral
		this.Days = 90
		break
	case Semiannually: // semestral
		this.Days = 180
		break
	case Yearly: // anual
		this.Days = 365
		break
	}
}

type PixAdditionalFields struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewPixAdditionalFields(name string, value string) *PixAdditionalFields {
	return &PixAdditionalFields{Name: name, Value: value}
}

type Payment struct {
	Amount        int64           `json:"amount" valid:"Required"`
	ReferenceKey  string          `json:"reference_key,omitempty" valid:""`
	ApiKey        string          `json:"api_key" valid:"Required"`
	Installments  int64           `json:"Installments,omitempty" valid:"Required"` // parcelas
	Customer      *Customer       `json:"customer" valid:"Required"`
	PaymentMethod api.PaymentType `json:"payment_method" valid:"Required"`

	PostbackUrl    string                 `json:"postback_url,omitempty" valid:""`
	SoftDescriptor string                 `json:"soft_descriptor,omitempty" valid:"Required"` // nome que aparece na fatura do cliente
	Metadata       map[string]interface{} `json:"metadata,omitempty"`

	Capture bool `json:"capture"`
	Async   bool `json:"async"`

	BoletoExpirationDate string `json:"boleto_expiration_date,omitempty" valid:""`
	BoletoInstructions   string `json:"boleto_instructions,omitempty" valid:""`

	CardId             string `json:"card_id,omitempty" valid:""`
	CardHolderName     string `json:"card_holder_name,omitempty" valid:""`
	CardExpirationDate string `json:"card_expiration_date,omitempty" valid:""`
	CardNumber         string `json:"card_number,omitempty" valid:""`
	CardCvv            string `json:"card_cvv,omitempty" valid:""`
	CardHash           string `json:"card_hash,omitempty" valid:""`

	SplitRules     []*SplitRule    `json:"split_rules,omitempty"`
	BoletoFine     *BoletoFine     `json:"boleto_fine"`
	BoletoInterest *BoletoInterest `json:"boleto_interest"`
	BoletoRule     []BoletoRule    `json:"boleto_rules"`

	PixAdditionalFields []*PixAdditionalFields `json:"pix_additional_fields,omitempty"`
	PixExpirationDate   string                 `json:"pix_expiration_date,omitempty"`
}

func NewPaymentWithCard(amount float64) *Payment {
	return &Payment{Amount: FormatAmount(amount), Installments: 1, PaymentMethod: api.PaymentTypeCreditCard, Customer: new(Customer)}
}

func NewPaymentWithBoleto(amount float64) *Payment {
	return &Payment{PaymentMethod: api.PaymentTypeBoleto, Amount: FormatAmount(amount), Customer: new(Customer), Installments: 1}
}

func NewPaymentWithPix(amount float64) *Payment {
	return &Payment{PaymentMethod: api.PaymentTypePix, Amount: FormatAmount(amount), PixAdditionalFields: []*PixAdditionalFields{}, Customer: new(Customer), Installments: 1}
}

func (this *Payment) SetPixExpirationDate(date time.Time) {
	this.PixExpirationDate = date.Format("2006-01-02")
}

func (this *Payment) AddPixAdditionalFields(name string, value string) {

	if this.PixAdditionalFields == nil {
		this.PixAdditionalFields = []*PixAdditionalFields{}
	}

	this.PixAdditionalFields = append(this.PixAdditionalFields, NewPixAdditionalFields(name, value))
}

func (this *Payment) AddMetadata(name string, value string) {

	if this.Metadata == nil {
		this.Metadata = make(map[string]interface{})
	}

	this.Metadata[name] = value
}

type Subscription struct {
	Id                 int64           `json:"id,omitempty"`
	PlanoId            int64           `json:"plan_id" valid:"Required"`
	ReferenceKey       string          `json:"reference_key,omitempty" valid:""`
	PaymentMethod      api.PaymentType `json:"payment_method" valid:"Required"`
	ApiKey             string          `json:"api_key" valid:"Required"`
	CardId             string          `json:"card_id,omitempty" valid:""`
	CardHolderName     string          `json:"card_holder_name,omitempty" valid:""`
	CardExpirationDate string          `json:"card_expiration_date,omitempty" valid:""`
	CardNumber         string          `json:"card_number,omitempty" valid:""`
	CardCvv            string          `json:"card_cvv,omitempty" valid:""`
	CardHash           string          `json:"card_hash,omitempty" valid:""`

	PostbackUrl string                 `json:"postback_url,omitempty" valid:""`
	Customer    *Customer              `json:"customer" valid:"Required"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`

	SplitRules     []*SplitRule    `json:"split_rules,omitempty"`
	BoletoFine     *BoletoFine     `json:"boleto_fine,omitempty"`
	BoletoInterest *BoletoInterest `json:"boleto_interest,omitempty"`
}

func NewSubscription(planoId int64) *Subscription {
	return &Subscription{PlanoId: planoId, Customer: new(Customer)}
}

func NewSubscriptionWithCard(planoId int64) *Subscription {
	return &Subscription{PlanoId: planoId, Customer: new(Customer), PaymentMethod: api.PaymentTypeCreditCard}
}

func NewSubscriptionWithBoleto(planoId int64) *Subscription {
	return &Subscription{PlanoId: planoId, Customer: new(Customer), PaymentMethod: api.PaymentTypeBoleto}
}

type CaptureData struct {
	ApiKey        string                 `json:"api_key" valid:"Required"`
	TransactionId string                 `json:"-" valid:"Required"` // id or token
	Amount        int64                  `json:"amount" valid:"Required"`
	SplitRules    []*SplitRule           `json:"split_rules,omitempty"`
	Metadata      map[string]interface{} `json:"metadata"`
}

func NewCaptureData(transactionId string, amount float64) *CaptureData {
	return &CaptureData{TransactionId: transactionId, Amount: FormatAmount(amount)}
}

type Card struct {
	Id             string `json:"card_id,omitempty" valid:""`
	HolderName     string `json:"card_holder_name" valid:""`
	ExpirationDate string `json:"card_expiration_date" valid:""`
	Number         string `json:"card_number" valid:""`
	Cvv            string `json:"card_cvv,omitempty" valid:""`
	CustomerId     string `json:"customer_id,omitempty" valid:""`
	Hash           string `json:"card_hash,omitempty" valid:""`
	ApiKey         string `json:"api_key,omitempty" valid:""`
}

type CardHashKey struct {
	Id        int64  `json:"id"`
	PublicKey string `json:"public_key"`
	Hash      string
}

type CardResult struct {
	Object         string    `json:"card"`
	Id             string    `json:"id"`
	DateCreated    time.Time `json:"date_created"`
	DateUpdated    time.Time `json:"date_updated"`
	Brand          string    `json:"brand"`
	HolderName     string    `json:"holder_name"`
	FirstDigits    string    `json:"first_digits"`
	LastDigits     string    `json:"last_digits"`
	Country        string    `json:"country"`
	Fingerprint    string    `json:"fingerprint"`
	Customer       string    `json:"customer"`
	ExpirationDate string    `json:"expiration_date"`
	Valid          bool      `json:"valid"`
}

type Movement struct {
	Object           string          `json:"object"`             //: "balance_operation",
	Id               int64           `json:"id"`                 //: 4852,
	Status           string          `json:"status"`             //: "available",
	BalanceAmount    int64           `json:"balance_amount"`     //: 2920013,
	BalanceOldAmount int64           `json:"balance_old_amount"` //: 2910128,
	Type             string          `json:"type"`               //: "payable",
	Amount           int64           `json:"amount"`             //: 10000,
	Fee              int64           `json:"fee"`                //: 115,
	DateCreated      string          `json:"date_created"`       //: "2015-03-06T18:44:42.000Z",
	MovementObject   *MovementObject `json:"movement_object"`    //
}

type MovementObject struct {
	Object        string `json:"object"`         //"payable",
	Id            int64  `json:"id"`             //1294,
	Status        string `json:"status"`         //"paid",
	Amount        int64  `json:"amount"`         //10000,
	Fee           int64  `json:"fee"`            //115,
	Installment   int64  `json:"installment"`    //1,
	TransactionId int64  `json:"transaction_id"` //185507,
	PaymentDate   string `json:"payment_date"`   //"2015-03-06T03:00:00.000Z",
	DateCreated   string `json:"date_created"`   //"2015-03-06T18:44:42.000Z"
}

type BalanceValue struct {
	Amount int64 `json:"amount"`
}

type BankAccount struct {
	Object             string `json:"object"`               //"bank_account",
	Id                 int64  `json:"id"`                   //17346045,
	BankCode           string `json:"bank_code"`            //"000",
	Agencia            string `json:"agencia"`              //"00000",
	AgenciaDv          string `json:"agencia_dv"`           //"2",
	Conta              string `json:"conta"`                //"00000",
	ContaDv            string `json:"conta_dv"`             //"00",
	Type               string `json:"type"`                 //"conta_corrente",
	DocumentType       string `json:"document_type"`        //"cpf",
	DocumentNumber     string `json:"document_number"`      //"03602396681",
	LegalName          string `json:"legal_name"`           //"nome2",
	ChargeTransferFees bool   `json:"charge_transfer_fees"` //true,
	DateCreated        string `json:"date_created"`         //"2016-12-27T22:08:10.536Z"
}

func NewBankAccount(id int64) *BankAccount {
	return &BankAccount{Id: id}
}

type Transfer struct {
	Amount        int64  `json:"amount"`
	BankAccountId int64  `json:"bank_account_id"`
	ApiKey        string `json:"api_key"`
}

func NewTransfer(amount float64, BankAccount *BankAccount) *Transfer {
	return &Transfer{
		Amount:        FormatAmount(amount),
		BankAccountId: BankAccount.Id,
	}
}

type TransferResult struct {
	Object               string `json:"object"`                 // "transfer",
	Id                   int64  `json:"id"`                     // 65485,
	Amount               int64  `json:"amount"`                 // 100,
	Type                 string `json:"type"`                   // "ted",
	Status               string `json:"status"`                 // "pending_transfer",
	SourceType           string `json:"source_type"`            // "recipient",
	SourceId             string `json:"source_id"`              // "re_cix7pxz6f02ppcv6dn4ckcrcc",
	TargetType           string `json:"target_type"`            // "bank_account",
	TargetId             string `json:"target_id"`              // "17346045",
	Fee                  int64  `json:"fee"`                    // 367,
	FundingDate          string `json:"funding_date"`           // null,
	FundingEstimatedDate string `json:"funding_estimated_date"` // "2017-02-18T02:00:00.000Z",
	TransactionId        string `json:"transaction_id"`         // null,
	DateCreated          string `json:"date_created"`           // "2017-02-17T16:24:20.933Z",
	BankAccount          string `json:"bank_account"`           //
	ApiKey               string `json:"api_key,omitempty" valid:""`
}

type Response struct {
	Object               string `json:"object"`
	StatusText           string `json:"status"` // processing, authorized, paid, refunded, waiting_payment, pending_refund, refused
	OldStatusText        string
	DesiredStatusText    string
	RefuseReason         string      `json:"refuse_reason"` // acquirer, antifraud, internal_error, no_acquirer, acquirer_timeout
	StatusReason         string      `json:"status_reason"` // acquirer, antifraud, internal_error, no_acquirer, acquirer_timeout
	AcquirerName         string      `json:"acquirer_name"` // tone, cielo, rede
	AcquirerId           string      `json:"acquirer_id"`
	AcquirerResponseCode string      `json:"acquirer_response_code"`
	AuthorizationCode    string      `json:"authorization_code"`
	SoftDescriptor       string      `json:"soft_descriptor"`
	Tid                  interface{} `json:"tid"`
	Nsu                  int64       `json:"nsu"`
	DateCreated          string      `json:"date_created"`
	DateUpdated          string      `json:"date_updated"`
	Amount               int64       `json:"amount"`
	AuthorizedAmount     int64       `json:"authorized_amount"`
	PaidAmount           int64       `json:"paid_amount"`
	RefundedAmount       int64       `json:"refunded_amount"`
	Installments         int64       `json:"installments"`
	Id                   int64       `json:"id"`
	Cost                 float64     `json:"cost"`
	CardHolderName       string      `json:"card_holder_name"`
	CardLastDigits       string      `json:"card_last_digits"`
	CardFirstDigits      string      `json:"card_first_digits"`
	CardBrand            string      `json:"card_brand"`
	CardPinMode          string      `json:"card_pin_mode"`
	PostbackUrl          string      `json:"postback_url"`
	PaymentMethod        string      `json:"payment_method"`
	CaptureMethod        string      `json:"capture_method"`
	AntifraudScore       string      `json:"antifraud_score"`
	BoletoUrl            string      `json:"boleto_url"`
	BoletoBarcode        string      `json:"boleto_barcode"`
	BoletoExpirationDate string      `json:"boleto_expiration_date"`
	Referer              string      `json:"referer"`
	Ip                   string      `json:"ip"`
	ReferenceKey         string      `json:"reference_key"`
	ManageUrl            string      `json:"manage_url"`

	PixQrCode           string                 `json:"pix_qr_code"`
	PixExpirationDate   string                 `json:"pix_expiration_date"`
	PixAdditionalFields []*PixAdditionalFields `json:"pix_additional_fields"`

	Errors         []*ResponseErrorItem `json:"errors"`
	ResponseValues map[string]interface{}
	Response       string
	Request        string
	Message        string
	Error          bool

	Status api.PagarmeStatus

	CardResult         *CardResult `json:"card"`
	CardHashKey        *CardHashKey
	Plano              *Plano                 `json:"plan"`
	CurrentTransaction *Response              `json:"current_transaction"`
	CurrentPeriodStart string                 `json:"current_period_start"`
	CurrentPeriodEnd   string                 `json:"current_period_end"`
	Charges            int64                  `json:"charges,omitempty"`
	Customer           *Customer              `json:"customer"`
	Metadata           map[string]interface{} `json:"metadata"`
	Fine               *BoletoFine            `json:"fine"`
	Interest           *BoletoInterest        `json:"interest"`

	Transactions []*Response

	WaitingFunds *BalanceValue `json:"waiting_funds"`
	Available    *BalanceValue `json:"available"`
	Transferred  *BalanceValue `json:"transferred"`

	Movements []*Movement

	TransferResult  *TransferResult
	TransferResults []*TransferResult
}

func NewResponse() *Response {
	return &Response{
		CardResult:      new(CardResult),
		CardHashKey:     new(CardHashKey),
		Plano:           new(Plano),
		Errors:          []*ResponseErrorItem{},
		ResponseValues:  make(map[string]interface{}),
		Transactions:    []*Response{},
		WaitingFunds:    new(BalanceValue),
		Available:       new(BalanceValue),
		Transferred:     new(BalanceValue),
		Movements:       []*Movement{},
		TransferResult:  new(TransferResult),
		TransferResults: []*TransferResult{},
	}
}

func (this *Response) HasTransactions() bool {
	return this.Transactions != nil && len(this.Transactions) > 0
}

func (this *Response) TransactionsCount() int {
	if this.HasTransactions() {
		return len(this.Transactions)
	}
	return 0
}

func (this *Response) FirstTransaction() *Response {
	if this.HasTransactions() {
		return this.Transactions[0]
	}
	return nil
}

func (this *Response) LastTransaction() *Response {
	if this.HasTransactions() {
		return this.Transactions[len(this.Transactions)-1]
	}
	return nil
}

func (this *Response) HasMovements() bool {
	return this.Movements != nil && len(this.Movements) > 0
}

func (this *Response) MovementsCount() int {
	if this.HasMovements() {
		return len(this.Movements)
	}
	return 0
}

func (this *Response) FirstMovement() *Movement {
	if this.HasMovements() {
		return this.Movements[0]
	}
	return nil
}

func (this *Response) LastMovement() *Movement {
	if this.HasMovements() {
		return this.Movements[len(this.Movements)-1]
	}
	return nil
}

func (this *Response) ToMap() map[string]string {

	errorMap := make(map[string]string)

	if this.Errors != nil {
		for _, it := range this.Errors {
			errorMap[fmt.Sprintf("%v, %v", it.ParameterName, it.Type)] = it.Message
		}
	}

	return errorMap

}

func (this *Response) HasError() bool {
	if this.Errors != nil {
		return len(this.Errors) > 0
	}
	return false
}

func (this *Response) ErrorsCount() int {
	if this.Errors != nil {
		return len(this.Errors)
	}
	return 0
}

func (this *Response) FirstError() string {
	if this.Errors != nil {
		for _, it := range this.Errors {
			return fmt.Sprintf("%v: %v", it.ParameterName, it.Message)
		}
	}
	return ""
}

func (this *Response) GetPayZenSOAPStatus() api.TransactionStatus {
	switch this.Status {
	case api.PagarmeProcessing:
		return api.Created
	case api.PagarmeAuthorized:
		return api.Authorised
	case api.PagarmePaid:
		return api.Captured
	case api.PagarmeRefunded:
		return api.Refunded
	case api.PagarmeWaitingPayment:
		return api.WaitingPayment
	case api.PagarmeUnpaid:
		return api.WaitingPayment
	case api.PagarmePendingRefund:
		return api.PendingRefund
	case api.PagarmeRefused:
		return api.Refused
	case api.PagarmeChargedback:
		return api.Chargeback
	case api.PagarmeAnalyzing:
		return api.Analyzing
	case api.PagarmeCancelled:
		return api.Cancelled
	case api.PagarmePendingReview:
		return api.Other
	case api.PagarmeSuccess:
		return api.Success
	default:
		return api.Error
	}
}

func (this *Response) BuildStatus() {
	switch this.StatusText {
	case "processing":
		this.Status = api.PagarmeProcessing
		break
	case "authorized":
		this.Status = api.PagarmeAuthorized
		break
	case "paid":
		this.Status = api.PagarmePaid
		break
	case "unpaid":
		this.Status = api.PagarmeUnpaid
		break
	case "refunded":
		this.Status = api.PagarmeRefunded
		break
	case "canceled":
		this.Status = api.PagarmeCancelled
		break
	case "waiting_payment":
		this.Status = api.PagarmeWaitingPayment
		break
	case "pending_refund":
		this.Status = api.PagarmeRefunded
		break
	case "chargedback":
		this.Status = api.PagarmeChargedback
		break
	case "analyzing":
		this.Status = api.PagarmeAnalyzing
		break
	case "pending_review":
		this.Status = api.PagarmePendingReview
		break
	case "refused":
		this.Message = "A transação foi recusada, verifique os dados cartão"
		this.Error = true
		this.Status = api.PagarmeRefused
		break
	case "":
		this.Status = api.PagarmeEmpty
		break
	default:
		this.Message = fmt.Sprintf("Status desconhecido: %v", this.StatusText)
		this.Error = true
		this.Status = api.PagarmeRefused
	}
}

func FormatAmount(amount float64) int64 {
	ac := accounting.Accounting{Symbol: "", Precision: 2, Thousand: "", Decimal: ""}
	text := strings.Replace(ac.FormatMoney(amount), ",", "", -1)
	text = strings.Replace(text, ".", "", -1)
	val, _ := strconv.Atoi(text)
	return int64(val)
}

func FormatAmountToFloat(amount int64) float64 {

	unformateed := accounting.UnformatNumber(fmt.Sprintf("%v", amount), 2, "BRL")
	val, _ := strconv.ParseFloat(unformateed, 64)
	return val
	/*text := fmt.Sprintf("%v", amount)
	if len(text) >= 2 {
		a := text[:len(text)-2]
		b := text[len(text)-2:]

		if len(strings.TrimSpace(a)) == 0 {
			a = "0"
		}

		val, _ := strconv.ParseFloat(fmt.Sprintf("%v.%v", a, b), 64)
		return val
	} else {
		return 0
	}*/

}
