package v4

import (
	"fmt"
	"github.com/leekchan/accounting"
	"github.com/mobilemindtec/go-payments/api"
	"strconv"
	"strings"
	"time"
)

const (
	Currency               = "BRL"
	TimeZoneDateTimeLayout = "2006-01-02T15:04:05-07:00"
)

type TransactionStatus int64
type PaymentStatus int64
type ApiMode int64

const (
	ACCEPTED TransactionStatus = 1 + iota
	AUTHORISED
	AUTHORISED_TO_VALIDATE
	CANCELLED
	CAPTURED
	EXPIRED
	PARTIALLY_AUTHORISED
	REFUSED
	UNDER_VERIFICATION
	WAITING_AUTHORISATION
	WAITING_AUTHORISATION_TO_VALIDATE
	ERROR
	EMPTY
	UNKNOW
)

const (
	PAID PaymentStatus = 1 + iota
	UNPAID
	RUNNING
	PARTIALLY_PAID
)

const (
	Test ApiMode = 1 + iota
	Prod
)

func ConvertAmount(amount float64) int64 {

	if amount == 0.0 {
		return 0
	}

	ac := accounting.Accounting{Symbol: "", Precision: 2, Thousand: "", Decimal: ""}
	text := strings.Replace(ac.FormatMoney(amount), ",", "", -1)
	text = strings.Replace(text, ".", "", -1)
	val, _ := strconv.ParseInt(text, 10, 64)
	return val
}

type BillingDetails struct {
	Address         string `json:"address" valid:""`
	Address2        string `json:"address2,omitempty"`
	StreetNumber    string `json:"streetNumber" valid:""`
	ZipCode         string `json:"zipCode" valid:""`
	CellPhoneNumber string `json:"cellPhoneNumber,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	City            string `json:"city" valid:""`
	State           string `json:"state" valid:""`
	Country         string `json:"country" valid:""`
	District        string `json:"district" valid:""`
	FirstName       string `json:"firstName" valid:"Required"`
	LastName        string `json:"lastName" valid:"Required"`
	IdentityCode    string `json:"identityCode" valid:"Required"` // documento
	Category        string `json:"category"`                      // PRIVATE, COMPANY
	LegalName       string `json:"legalName,omitempty"`           // razão social
	Language        string `json:"language"`                      // PT
}

func NewBillingDetails() *BillingDetails {
	return &BillingDetails{Language: "PT", Category: "PRIVATE", Country: "BR"}
}

type ShippingDetails struct {
	Category            string `json:"category,omitempty"` // PRIVATE or COMPANY
	FirstName           string `json:"firstName" valid:"Required"`
	LastName            string `json:"lastName" valid:"Required"`
	PhoneNumber         string `json:"phoneNumber,omitempty"`
	StreetNumber        string `json:"streetNumber" valid:""`
	Address             string `json:"address" valid:"Required"`
	Address2            string `json:"address2,omitempty"`
	District            string `json:"district" valid:"Required"`
	ZipCode             string `json:"zipCode" valid:""`
	City                string `json:"city" valid:"Required"`
	State               string `json:"state" valid:"Required"`
	Country             string `json:"country" valid:"Required"`
	DeliveryCompanyName string `json:"deliveryCompanyName,omitempty"` // PRIVATE or COMPANY
	ShippingSpeed       string `json:"shippingSpeed,omitempty"`       // STANDARD EXPRESS PRIORITY
	ShippingMethod      string `json:"shippingMethod,omitempty"`      // RECLAIM_IN_SHOP (retirada na loja) VERIFIED_ADDRESS (entrega no endereço)
	LegalName           string `json:"legalName,omitempty"`
	IdentityCode        string `json:"identityCode,omitempty"`
}

func NewShippingDetails() *ShippingDetails {
	return &ShippingDetails{}
}

type CartItemInfo struct {
	ProductLabel string `json:"productLabel,omitempty"`

	/*
		FOOD_AND_GROCERY 	Produtos alimentares e de mercearia
		AUTOMOTIVE 	Automóvel / Moto
		ENTERTAINMENT 	Lazer / Cultura
		HOME_AND_GARDEN 	Casa e jardim
		HOME_APPLIANCE 	Equipamentos para a casa
		AUCTION_AND_GROUP_BUYING 	Leilões e compras em grupo
		FLOWERS_AND_GIFTS 	Flores e presentes
		COMPUTER_AND_SOFTWARE 	Computadores e softwares
		HEALTH_AND_BEAUTY 	Saúde e beleza
		SERVICE_FOR_INDIVIDUAL 	Serviços para pessoa física
		SERVICE_FOR_BUSINESS 	Serviços para pessoa jurídica
		SPORTS 	Esportes
		CLOTHING_AND_ACCESSORIES 	Roupas e acessórios
		TRAVEL 	Viagem
		HOME_AUDIO_PHOTO_VIDEO 	Som, imagem e vídeo
		TELEPHONY 	Telefonia
	*/
	ProductType string `json:"productType,omitempty"`

	ProductRef    string `json:"productRef,omitempty"`
	ProductQty    string `json:"productQty,omitempty"`
	ProductAmount string `json:"productAmount,omitempty"`
	ProductVat    string `json:"productVat,omitempty"` // imposto
}

func NewCartItemInfo() *CartItemInfo {
	return &CartItemInfo{}
}

type ShoppingCart struct {
	InsuranceAmount string          `json:"insuranceAmount,omitempty"` // seguro
	ShippingAmount  string          `json:"shippingAmount,omitempty"`  // despesas
	TaxAmount       string          `json:"taxAmount,omitempty"`       // impostos
	CartItemInfo    []*CartItemInfo `json:"cartItemInfo"`
}

func NewShoppingCart() *ShoppingCart {
	return &ShoppingCart{CartItemInfo: []*CartItemInfo{}}
}

type Customer struct {
	Email           string           `json:"email" valid:"Required"`
	IpAddress       string           `json:"ipAddress,omitempty"`
	Reference       string           `json:"reference,omitempty"`
	BillingDetails  *BillingDetails  `json:"billingDetails"`
	ShippingDetails *ShippingDetails `json:"shippingDetails"`
	ShoppingCart    *ShoppingCart    `json:"shoppingCart"`
}

func NewCustomer() *Customer {
	return &Customer{BillingDetails: NewBillingDetails(), ShoppingCart: NewShoppingCart(), ShippingDetails: NewShippingDetails()}
}

type CardOptions struct {
	ManualValidation      string `json:"manualValidation"` // NO
	CaptureDelay          int64  `json:"captureDelay"`
	FirstInstallmentDelay int64  `json:"firstInstallmentDelay"` // Número de meses adiados a serem aplicados à primeira parcela para um pagamento parcelado.
	InstallmentNumber     int64  `json:"installmentNumber"`     // parcelas
	Retry                 int64  `json:"retry"`                 // default 3
}

type TransactionOptions struct {
	CardOptions *CardOptions `json:"cardOptions"`
}

type Card struct {
	PaymentMethodType     string `json:"paymentMethodType" valid:"Required"` // CARD
	Number                string `json:"pan" valid:"Required"`               // card number
	ExpiryMonth           int64  `json:"expiryMonth" valid:"Required"`
	ExpiryYear            int64  `json:"expiryYear" valid:"Required"`
	SecurityCode          string `json:"securityCode" valid:"Required"`
	Brand                 string `json:"brand" valid:"Required"`
	CardHolderName        string `json:"cardHolderName" valid:"Required"`
	FirstInstallmentDelay int64  `json:"firstInstallmentDelay,omitempty"`
	InstallmentNumber     int64  `json:"installmentNumber"`
	PaymentMethodToken    string `json:"paymentMethodToken,omitempty"`
}

type Device struct {
	DeviceType     string `json:"deviceType" valid:"Required"`
	AcceptHeader   string `json:"acceptHeader" valid:"Required"`
	ColorDepth     string `json:"colorDepth" valid:"Required"`
	JavaEnabled    bool   `json:"javaEnabled" valid:"Required"`
	Language       string `json:"language" valid:"Required"`
	ScreenHeight   int64  `json:"screenHeight" valid:"Required"`
	ScreenWidth    int64  `json:"screenWidth" valid:"Required"`
	TimeZoneOffset int64  `json:"timeZoneOffset" valid:""`
	UserAgent      string `json:"userAgent" valid:"Required"`
	Ip             string `json:"ip,omitempty"`
}

func NewDevice() *Device {
	return &Device{
		DeviceType:     "BROWSER",
		AcceptHeader:   "application/json",
		ColorDepth:     "32",
		JavaEnabled:    true,
		Language:       "PT",
		ScreenHeight:   768,
		ScreenWidth:    1024,
		TimeZoneOffset: 0,
		UserAgent:      "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101",
	}
}

func NewCard() *Card {
	return &Card{PaymentMethodType: "CARD"}
}

type Payment struct {
	PaymentOrderId     string              `json:"paymentOrderId,omitempty"`
	OrderId            string              `json:"orderId" valid:"Required"`
	Amount             int64               `json:"amount" valid:"Required"`
	Currency           string              `json:"currency" valid:"Required"`
	IpnTargetUrl       string              `json:"ipnTargetUrl,omitempty"`       // URL de notificação estantânea
	PaymentMethodToken string              `json:"paymentMethodToken,omitempty"` // topken de um cartão
	Customer           *Customer           `json:"customer"`
	TransactionOptions *TransactionOptions `json:"transactionOptions,omitempty"`
	FormAction         string              `json:"formAction" valid:"Required"` // PAYMENT
	PaymentForms       []*Card             `json:"paymentForms"`
	Device             *Device             `json:"device"`
	Metadata           map[string]string   `json:"metadata,omitempty"`
	FingerPrintId      string              `json:"fingerPrintId,omitempty"`
	Card               *Card               `json:"-"`
}

func (this *Payment) SetAmount(amount float64) {
	this.Amount = ConvertAmount(amount)
}

func NewPayment(amount float64) *Payment {
	card := NewCard()
	return &Payment{
		FormAction:   "PAYMENT",
		Amount:       ConvertAmount(amount),
		Currency:     Currency,
		Customer:     NewCustomer(),
		Card:         card,
		PaymentForms: []*Card{card},
		Device:       NewDevice(),
		Metadata:     map[string]string{},
	}
}

type Subscription struct {
	SubscriptionId     string `json:"subscriptionId,omitempty"`
	OrderId            string `json:"orderId" valid:"Required"`
	Amount             int64  `json:"amount" valid:"Required"`
	Currency           string `json:"currency" valid:"Required"`
	PaymentMethodToken string `json:"paymentMethodToken" valid:"Required"`
	//subscription
	Comment             string            `json:"comment"`
	Description         string            `json:"description"`
	EffectDate          string            `json:"effectDate" valid:"Required"`   // date de inicio da recorrencia
	InitialAmount       int64             `json:"initialAmount,omitempty"`       // Valor das primeiras parcelas. O valor deve ser um número inteiro positivo (ex: 1234 para R$ 12,34).
	InitialAmountNumber int64             `json:"initialAmountNumber,omitempty"` // Quantidade de parcelas às quais aplicar o valor definido em initialAmount.
	Metadata            map[string]string `json:"metadata,omitempty"`
	FingerPrintId       string            `json:"fingerPrintId,omitempty"`
	Rrule               string            `json:"rrule" valid:"Required"`

	TransactionOptions *TransactionOptions `json:"transactionOptions,omitempty"`

	Count                   int64                 `json:"-"`
	Cycle                   api.SubscriptionCycle `json:"-"`
	PaymentAtLastDayOfMonth bool                  `json:"-"`
	PaymentAtDayOfMonth     int64                 `json:"-"`
}

func (this *Subscription) SetRule(cycle api.SubscriptionCycle, count int64, lastDayOfMonth bool, dayOfMonth int64) {
	this.Cycle = cycle
	this.Count = count
	this.PaymentAtLastDayOfMonth = lastDayOfMonth
	this.PaymentAtDayOfMonth = dayOfMonth
}

func (this *Subscription) SetAmount(amount float64) {
	this.Amount = ConvertAmount(amount)
}

func (this *Subscription) SetInitialAmount(initialAmount float64) {
	this.InitialAmount = ConvertAmount(initialAmount)
}

func (this *Subscription) BuildRule() string {

	if len(this.Rrule) == 0 {

		rules := []string{}

		if this.Count > 0 {
			rules = append(rules, fmt.Sprintf("COUNT=%v", this.Count))
		}

		switch this.Cycle {
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
			return "cycle is required"
		}

		if this.PaymentAtLastDayOfMonth && this.PaymentAtDayOfMonth > 0 {
			return "use PaymentAtLastDayOfMonth or PaymentAtDayOfMonth"
		}

		if this.PaymentAtDayOfMonth > 0 {
			rules = append(rules, fmt.Sprintf("BYMONTHDAY=%v", this.PaymentAtDayOfMonth))
		}

		if this.PaymentAtLastDayOfMonth {
			rules = append(rules, "BYMONTHDAY=28,29,30,31")
			rules = append(rules, "BYSETPOS=-1")
		}

		if len(rules) == 0 {
			return "is required"
		}

		this.Rrule = "RRULE"

		for i, it := range rules {
			if i == 0 {
				this.Rrule = fmt.Sprintf("%v:%v", this.Rrule, it)
			} else {
				this.Rrule = fmt.Sprintf("%v;%v", this.Rrule, it)
			}
		}

	}

	return ""
}

func NewSubscription(orderId string, amount float64, token string, effectDate time.Time) *Subscription {
	return &Subscription{
		EffectDate:         effectDate.Format(TimeZoneDateTimeLayout),
		OrderId:            orderId,
		Amount:             ConvertAmount(amount),
		Currency:           Currency,
		PaymentMethodToken: token,
	}
}

type OrderDetails struct {
	OrderEffectiveAmount int64     `json:"orderEffectiveAmount"`
	OrderTotalAmount     int64     `json:"orderTotalAmount"`
	Mode                 string    `json:"mode"` // TEST,PRODUCTION
	OrderCurrency        string    `json:"orderCurrency"`
	OrderId              string    `json:"orderId"`
	Customer             *Customer `json:"customer"`
}

type SubscriptionDetails struct {
	SubscriptionId string `json:"subscriptionId"`
}

type RiskAnalysi struct {
	ResultCode string `json:"resultCode"`
	Status string `json:"status"`
	RequestId string `json:"requestId"`
	Score string `json:"score"`
}

type FraudManagement struct {
	RiskAnalysis []*RiskAnalysi `json:"riskAnalysis"`
}

type TransactionDetails struct {
	ExternalTransactionId string               `json:"externalTransactionId"` // nsu
	Nsu                   string               `json:"nsu"`                   // nsu
	SubscriptionDetails   *SubscriptionDetails `json:"subscriptionDetails"`
	FraudManagement *FraudManagement `json:"fraudManagement"`
}

type Transaction struct {
	Amount               int64               `json:"amount"`
	creationDate         string              `json:"creationDate"`
	Currency             string              `json:"currency"`
	DetailedErrorCode    string              `json:"detailedErrorCode"`    // erro da adquirente
	DetailedErrorMessage string              `json:"detailedErrorMessage"` // erro da adquirente
	ErrorCode            string              `json:"errorCode"`            // erro payzen
	ErrorMessage         string              `json:"errorMessage"`         // erro payzen
	DetailedStatus       string              `json:"detailedStatus"`       //
	OperationType        string              `json:"operationType"`        // DEBIT,CREDIT,VERIFICATION
	Uuid                 string              `json:"uuid"`
	PaymentMethodToken   string              `json:"paymentMethodToken"`
	PaymentMethodType    string              `json:"paymentMethodType"` // CARD
	ShopId               string              `json:"shopId"`
	Status               string              `json:"status"`
	TransactionDetails   *TransactionDetails `json:"transactionDetails"`
	TransactionStatus    TransactionStatus
	PaymentStatus        PaymentStatus
	ResponseCode         int64 `json:"responseCode,omitempty"`
}

func (this *Transaction) HasError() bool {
	return this.TransactionStatus == ERROR
}

func (this *Transaction) BuildStatus() {

	switch this.DetailedStatus {
	case "ACCEPTED":
		this.TransactionStatus = ACCEPTED
		break
	case "AUTHORISED":
		this.TransactionStatus = AUTHORISED
		break
	case "AUTHORISED_TO_VALIDATE":
		this.TransactionStatus = AUTHORISED_TO_VALIDATE
		break
	case "CANCELLED":
		this.TransactionStatus = CANCELLED
		break
	case "CAPTURED":
		this.TransactionStatus = CAPTURED
		break
	case "EXPIRED":
		this.TransactionStatus = EXPIRED
		break
	case "PARTIALLY_AUTHORISED":
		this.TransactionStatus = PARTIALLY_AUTHORISED
		break
	case "REFUSED":
		this.TransactionStatus = REFUSED
		break
	case "UNDER_VERIFICATION":
		this.TransactionStatus = UNDER_VERIFICATION
		break
	case "WAITING_AUTHORISATION":
		this.TransactionStatus = WAITING_AUTHORISATION
		break
	case "WAITING_AUTHORISATION_TO_VALIDATE":
		this.TransactionStatus = WAITING_AUTHORISATION_TO_VALIDATE
		break
	case "ERROR":
		this.TransactionStatus = ERROR
		break
	case "":

		if len(this.ErrorCode) > 0 {
			this.TransactionStatus = ERROR
		} else {
			this.TransactionStatus = EMPTY
		}
		break
	default:
		this.TransactionStatus = UNKNOW
		break
	}

	switch this.Status {
	case "PAID":
		this.PaymentStatus = PAID
		break
	case "UNPAID":
		this.PaymentStatus = UNPAID
		break
	case "RUNNING":
		this.PaymentStatus = RUNNING
		break
	case "PARTIALLY_PAID":
		this.PaymentStatus = PARTIALLY_PAID
		break
	}
}

func (this *Transaction) GetSOAPStatus() api.TransactionStatus {

	this.BuildStatus()

	switch this.TransactionStatus {
	case ACCEPTED:
		return api.Success
	case AUTHORISED:
		return api.Authorised
	case AUTHORISED_TO_VALIDATE:
		return api.AuthorisedToValidate
	case CANCELLED:
		return api.Cancelled
	case CAPTURED:
		return api.Captured
	case EXPIRED:
		return api.Expired
	case PARTIALLY_AUTHORISED:
		return api.PartiallyAuthorised
	case REFUSED:
		return api.Refused
	case UNDER_VERIFICATION:
		return api.UnderVerification
	case WAITING_AUTHORISATION:
		return api.UnderVerification
	case WAITING_AUTHORISATION_TO_VALIDATE:
		return api.WaitingAuthorisationToValidate
	case ERROR:
		return api.Error
	case UNKNOW:
		return api.Error
	case EMPTY:
		return api.Success
	default:
		return api.NotCreated
	}
}

type Answer struct {
	OrderStatus          string         `json:"orderStatus"` // PAID, UNPAID, RUNNING, PARTIALLY_PAID
	OrderCycle           string         `json:"orderCycle"`  // OPEN, CLOSED
	ShopId               string         `json:"shopId"`
	OrderDetails         *OrderDetails  `json:"orderDetails"`
	Transactions         []*Transaction `json:"transactions"`
	TransactionDetails   *Transaction   `json:"transactionDetails"`
	DetailedErrorCode    string         `json:"detailedErrorCode"`    // erro da adquirente
	DetailedErrorMessage string         `json:"detailedErrorMessage"` // erro da adquirente
	ErrorCode            string         `json:"errorCode"`            // erro payzen
	ErrorMessage         string         `json:"errorMessage"`         // erro payzen
	CancellationDate     string         `json:"cancellationDate"`     // token cancelado
	ResponseCode         int64          `json:"responseCode,omitempty"`

	PaymentMethodToken string `json:"paymentMethodToken"`

	SubscriptionId string `json:"subscriptionId"`

	OrderId             string `json:"orderId"`                       //: "55d979c6-e3c8-45b1-bc63-a5d3f4c1969b",
	Description         string `json:"description"`                   //: "",
	Rrule               string `json:"rrule"`                         //: "RRULE:COUNT=12;FREQ=MONTHLY;BYMONTHDAY=28,29,30,31;BYSETPOS=-1",
	EffectDate          string `json:"effectDate"`                    //: "2021-07-30T18:46:12+00:00",
	CancelDate          string `json:"cancelDate"`                    //: null,
	InitialAmount       int64  `json:"initialAmount,omitempty"`       //: null,
	InitialAmountNumber int64  `json:"initialAmountNumber,omitempty"` //: null,
	PastPaymentsNumber  int64  `json:"pastPaymentsNumber,omitempty"`  //: 0,
	TotalPaymentsNumber int64  `json:"totalPaymentsNumber,omitempty"` //: 12,
	Metadata            string `json:"metadata"`                      //: null,
}

func (this *Answer) IsCancelled() bool {
	return len(this.CancellationDate) > 0 || len(this.CancelDate) > 0
}

type PaymentResponse struct {
	Answer              *Answer `json:"answer"`
	ApplicationProvider string  `json:"applicationProvider"` //"PAYZEN",
	ApplicationVersion  string  `json:"applicationVersion"`  //"5.26.0",
	Metadata            string  `json:"metadata"`            //null,
	Mode                string  `json:"mode"`                //"TEST",
	ServerDate          string  `json:"serverDate"`          //"2021-07-29T22:15:42+00:00",
	ServerUrl           string  `json:"serverUrl"`           //"https://api.payzen.com.br",
	Status              string  `json:"status"`              //"ERROR",
	Ticket              string  `json:"ticket"`              //"a4609adcc22e429ca6cd34526d48cccf",
	Version             string  `json:"version"`             //"V4",
	WebService          string  `json:"webService"`          //"PCI/Charge/CreateToken"
	Response            string
	Request             string
}

type TransationResponse struct {
	Answer              *Transaction `json:"answer"`
	ApplicationProvider string       `json:"applicationProvider"` //"PAYZEN",
	ApplicationVersion  string       `json:"applicationVersion"`  //"5.26.0",
	Metadata            string       `json:"metadata"`            //null,
	Mode                string       `json:"mode"`                //"TEST",
	ServerDate          string       `json:"serverDate"`          //"2021-07-29T22:15:42+00:00",
	ServerUrl           string       `json:"serverUrl"`           //"https://api.payzen.com.br",
	Status              string       `json:"status"`              //"ERROR",
	Ticket              string       `json:"ticket"`              //"a4609adcc22e429ca6cd34526d48cccf",
	Version             string       `json:"version"`             //"V4",
	WebService          string       `json:"webService"`          //"PCI/Charge/CreateToken"
	Response            string
	Request             string
}

type PayZenResult struct {
	Response             *PaymentResponse
	Transaction          *TransationResponse
	Errors               map[string]string
	Message              string
	Error                bool
	EmptyResponseSuccess bool
}

func (this *PayZenResult) IsCancelled() bool {
	if this.IsResponse() && this.Response.Answer != nil {
		return this.Response.Answer.IsCancelled()
	}
	return false
}

func (this *PayZenResult) BuildResult() {

	if this.Transaction != nil && this.Transaction.Answer != nil {

		this.Transaction.Answer.BuildStatus()

		if len(this.Transaction.Answer.ErrorCode) > 0 {
			this.Errors["ErrorCode"] = this.Transaction.Answer.ErrorCode
		}
		if len(this.Transaction.Answer.ErrorMessage) > 0 {
			this.Errors["ErrorMessage"] = this.Transaction.Answer.ErrorMessage
		}
		if len(this.Transaction.Answer.DetailedErrorCode) > 0 {
			this.Errors["DetailedErrorCode"] = this.Transaction.Answer.DetailedErrorCode
		}
		if len(this.Transaction.Answer.DetailedErrorMessage) > 0 {
			this.Errors["DetailedErrorMessage"] = this.Transaction.Answer.DetailedErrorMessage
			this.Message = this.Transaction.Answer.DetailedErrorMessage
		}

		this.Error = len(this.Errors) > 0

		if this.Transaction.Status == "ERROR" {
			this.Error = true
		}

		if !this.Error {
			if this.Transaction.Answer.TransactionStatus == EMPTY && this.Transaction.Answer.ResponseCode == 0 {
				this.EmptyResponseSuccess = true
			}
		}

	} else if this.Response != nil {

		if this.Response.Status == "ERROR" {
			this.Error = true
		}

		if len(this.Response.Answer.ErrorCode) > 0 {
			this.Errors["ErrorCode"] = this.Response.Answer.ErrorCode
		}
		if len(this.Response.Answer.ErrorMessage) > 0 {
			this.Errors["ErrorMessage"] = this.Response.Answer.ErrorMessage
		}
		if len(this.Response.Answer.DetailedErrorCode) > 0 {
			this.Errors["DetailedErrorCode"] = this.Response.Answer.DetailedErrorCode
		}
		if len(this.Response.Answer.DetailedErrorMessage) > 0 {
			this.Errors["DetailedErrorMessage"] = this.Response.Answer.DetailedErrorMessage
			this.Message = this.Response.Answer.DetailedErrorMessage
		}

		this.Error = len(this.Errors) > 0

		if this.Response.Status == "ERROR" {
			this.Error = true
		}

		haveTrans := false

		if this.Response.Answer.TransactionDetails != nil {
			this.Response.Answer.TransactionDetails.BuildStatus()
			haveTrans = true
		}

		if this.Response.Answer.Transactions != nil {
			for _, it := range this.Response.Answer.Transactions {
				it.BuildStatus()
				haveTrans = true
			}
		}

		if !this.Error {
			if this.Response.Answer.ResponseCode == 0 || !haveTrans {
				this.EmptyResponseSuccess = true
			}
		}
	}
}

func NewPayZenResultWithResponse(response *PaymentResponse) *PayZenResult {
	result := &PayZenResult{Response: response, Errors: make(map[string]string)}
	result.BuildResult()
	return result
}

func NewPayZenResultWithTransaction(transaction *TransationResponse) *PayZenResult {
	result := &PayZenResult{Transaction: transaction, Errors: make(map[string]string)}
	result.BuildResult()
	return result
}

func (this *PayZenResult) IsResponse() bool {
	return this.Response != nil
}

func (this *PayZenResult) IsTransaction() bool {
	return this.Transaction != nil
}

func (this *PayZenResult) GetResponse() *Answer {
	if this.IsResponse() {
		return this.Response.Answer
	}

	return nil
}

func (this *PayZenResult) GetTransaction() *Transaction {

	if this.IsTransaction() {
		return this.Transaction.Answer
	}

	if this.GetResponse().Transactions != nil && len(this.GetResponse().Transactions) > 0 {
		return this.GetResponse().Transactions[0]
	}

	return nil
}

func (this *PayZenResult) GetTransactions() []*Transaction {

	if this.IsTransaction() {
		return []*Transaction{this.Transaction.Answer}
	}

	if this.GetResponse().Transactions != nil && len(this.GetResponse().Transactions) > 0 {
		return this.GetResponse().Transactions
	}

	return nil
}

func (this *PayZenResult) GetHttpResponse() string {
	if this.IsResponse() {
		return this.Response.Response
	}

	return this.Transaction.Response
}

func (this *PayZenResult) GetHttpRequest() string {
	if this.IsResponse() {
		return this.Response.Request
	}

	return this.Transaction.Request
}

func (this *PayZenResult) SetHttpResponse(httpResponse string) {
	if this.IsResponse() {
		this.Response.Response = httpResponse
	} else {
		this.Transaction.Response = httpResponse
	}
}

func (this *PayZenResult) SetHttpRequest(httpRequest string) {
	if this.IsResponse() {
		this.Response.Request = httpRequest
	} else {
		this.Transaction.Request = httpRequest
	}
}
