package v5

import (
	"fmt"
	"log"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type CustomerType string
type DocumentType string
type Gender string
type PaymentMethod string
type BankCode string
type BoletoType string
type CardType string
type Status string
type Currency string
type Interval string
type BillingType string
type SchemeType string
type OperationType string
type InvoiceStatus string
type ChargeStatus string

const (
	BRL        Currency = "BRL"
	DateLayout          = "2006-01-02"
)

const (
	Percentage OperationType = "percentage"
	Flat       OperationType = "flat"
)
const (
	Unit    SchemeType = "unit"
	Package SchemeType = "package"
	Volume  SchemeType = "volume"
	Tier    SchemeType = "tier"
)

const (
	PrePaid  BillingType = "prepaid"
	PostPaid BillingType = "postpaid"
	ExactDay BillingType = "exact_day"
)

const (
	Day   Interval = "day"
	Week  Interval = "week"
	Month Interval = "month"
	Year  Interval = "year"
)

const (
	BrasilCode = "BR"
)

const (
	Individual CustomerType = "individual"
	Company    CustomerType = "company"
)
const (
	CPF      DocumentType = "CPF"
	CNPJ     DocumentType = "CNPJ"
	PASSPORT DocumentType = "PASSPORT"
)
const (
	Male   Gender = "male"
	Female Gender = "female"
)
const (
	MethodCreditCard PaymentMethod = "credit_card"
	MethodBoleto     PaymentMethod = "boleto"
	MethodPix        PaymentMethod = "pix"
)
const (
	AuthAndCapture OperationType = "auth_and_capture"
	AuthOnly       OperationType = "auth_only"
	PreAuth        OperationType = "pre_auth"
)

const (
	DM  BoletoType = "DM"
	BDP BoletoType = "BDP"
)

const (
	BancoBrasil BankCode = "001"
	Santander   BankCode = "033"
	Bradesco    BankCode = "237"
	Itau        BankCode = "341"
	Citibank    BankCode = "745"
	Caixa       BankCode = "104"
)

const (
	CardTypeCredit  CardType = "credit"
	CardTypeVoucher CardType = "voucher"
)

const (
	Active   Status = "active"
	Deleted  Status = "deleted"
	Expired  Status = "expired"
	Inactive Status = "inactive"
)

const (
	Pending   InvoiceStatus = "pending"
	Paid      InvoiceStatus = "paid"
	Canceled  InvoiceStatus = "canceled"
	Scheduled InvoiceStatus = "scheduled"
	Failed    InvoiceStatus = "failed"
)

const (
	ChargePending    ChargeStatus = "pending"
	ChargePaid       ChargeStatus = "paid"
	ChargeCanceled   ChargeStatus = "canceled"
	ChargeProcessing ChargeStatus = "processing"
	ChargeFailed     ChargeStatus = "failed"
	ChargeOverpaid   ChargeStatus = "overpaid"
	ChargeUnderpaid  ChargeStatus = "underpaid"
)

type Order struct {
	Code             string       `json:"code" valid:"Required;MaxSize(64)"`
	Customer         *Customer    `json:"customer"`
	CustomerId       string       `json:"customer_id,omitempty"`
	Items            []*OrderItem `json:"items"`
	Payments         []*Payment   `json:"payments"`
	Closed           bool         `json:"closed"` // Informa se o pedido será criado aberto ou fechado.
	Ip               string       `json:"ip" valid:"Required"`
	SessionId        string       `json:"session_id"`
	AntifraudEnabled bool         `json:"antifraud_enabled"`

	Id        string      `json:"id,omitempty"`
	Status    string      `json:"status,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
	UpdatedAt string      `json:"updated_at,omitempty"`
	Charges   []*Charge   `json:"charges,omitempty"`
	Checkouts []*Checkout `json:"checkouts,omitempty"`
}

type OrderPtr = *Order
type Orders = []OrderPtr

func NewOrder() *Order {
	return &Order{Payments: []*Payment{}, Items: []*OrderItem{}}
}

func (this *Order) AddItem() *OrderItem {
	item := &OrderItem{Quantity: 1}
	this.Items = append(this.Items, item)
	return item
}

func (this *Order) AddPayment(amount int64, method PaymentMethod) *Order {
	this.Payments = append(this.Payments, NewPayment(amount, method))
	return this
}

func (this *Order) WithBoleto(cb func(*Boleto)) *Order {
	failIf(len(this.Payments) == 0, "Payments must be greater than zero")
	payment := this.Payments[len(this.Payments)-1]
	failIf(payment.Boleto == nil, "Boleto must be not nil")
	cb(payment.Boleto)
	return this
}

func (this *Order) WithCreditCard(cb func(*CreditCard)) *Order {
	failIf(len(this.Payments) == 0, "Payments must be greater than zero")
	payment := this.Payments[len(this.Payments)-1]
	failIf(payment.CreditCard == nil, "CreditCard must be not nil")
	cb(payment.CreditCard)
	return this
}

func (this *Order) WithPix(cb func(*Pix)) *Order {
	failIf(len(this.Payments) == 0, "Payments must be greater than zero")
	payment := this.Payments[len(this.Payments)-1]
	failIf(payment.Pix == nil, "CreditCard must be not nil")
	cb(payment.Pix)
	return this
}

type OrderItem struct {
	Id          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	Quantity    int64  `json:"quantity"`
	Code        string `json:"code" valid:"Required"` // // código na plataforma
	Status      string `json:"status,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type Customer struct {
	Id         string            `json:"id,omitempty"`
	Name       string            `json:"name" valid:"Required;MaxSize(64)"`
	Type       CustomerType      `json:"type" valid:"Required"`
	Email      string            `json:"email" valid:"Required;Email"`
	Code       string            `json:"code" valid:"MaxSize(52)"` // código na plataforma
	Document   string            `json:"document" valid:"MaxSize(50)"`
	Gender     Gender            `json:"gender"`
	Delinquent bool              `json:"delinquent"`
	Address    *Address          `json:"address"`
	Phones     *Phones           `json:"phones"`
	Birthdate  string            `json:"birthdate"`
	Metadata   map[string]string `json:"metadata"`
	CreatedAt  string            `json:"created_at,omitempty"`
	UpdatedAt  string            `json:"updated_at,omitempty"`
}

func NewCustomer() *Customer {
	return &Customer{Address: NewAddress(), Phones: new(Phones), Type: Individual}
}

type CustomerPtr = *Customer
type Customers = []CustomerPtr

type Address struct {
	Id        string `json:"id,omitempty"`
	Country   string `json:"country" valid:"Required;MaxSize(2)"`
	State     string `json:"state" valid:"Required;MaxSize(2)"`
	City      string `json:"city" valid:"Required;MaxSize(2)"`
	ZipCode   string `json:"zip_code" valid:"Required;MaxSize(8)"`
	Line1     string `json:"line_1" valid:"Required"`
	Line2     string `json:"line_2" valid:""`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func NewAddress() *Address {
	return &Address{Country: BrasilCode}
}

type Phones struct {
	HomePhone   *Phone `json:"home_phone"`
	MobilePhone *Phone `json:"mobile_phone"`
}

type Phone struct {
	CountryCode string `json:"country_code"`
	AreaCode    string `json:"area_code"`
	Number      string `json:"number"`
}

func NewPhone(countryCode string, areaCode string, number string) *Phone {
	return &Phone{CountryCode: countryCode, AreaCode: areaCode, Number: number}
}

type Payment struct {
	PaymentMethod PaymentMethod `json:"payment_method"`
	CreditCard    *CreditCard   `json:"credit_card"`
	Boleto        *Boleto       `json:"boleto"`
	Pix           *Pix          `json:"pix"`
	Amount        int64         `json:"amount" valid:"Required"`
}

func NewPayment(amount int64, method PaymentMethod) *Payment {

	switch method {
	case MethodCreditCard:
		return &Payment{PaymentMethod: method, Amount: amount, CreditCard: NewCreditCard()}
	case MethodBoleto:
		return &Payment{PaymentMethod: method, Amount: amount, Boleto: NewBoleto()}
	case MethodPix:
		return &Payment{PaymentMethod: method, Amount: amount, Pix: NewPix()}
	default:
		fail("payment method %v not found", method)
		logs.Debug("payment method %v not found", method)
		return nil
	}

}

type CreditCard struct {
	OperationType       OperationType `json:"operation_type" valid:"Required"`
	Installments        int64         `json:"installments"`
	StatementDescriptor string        `json:"statement_descriptor" valid:"Required;MaxSize(13)"`
	Card                *Card         `json:"card"`
}

func NewCreditCard() *CreditCard {
	return &CreditCard{OperationType: AuthAndCapture, Installments: 1, Card: NewCard()}
}

type Card struct {
	Number           string          `json:"number,omitempty" valid:""`
	HolderName       string          `json:"holder_name,omitempty" valid:""`
	HolderDocument   string          `json:"holder_document,omitempty" valid:""`
	ExpMonth         int64           `json:"exp_month,omitempty" valid:""` // Value between 1 and 12 (included)
	ExpYear          int64           `json:"exp_year,omitempty" valid:""`  // Formatos yy ou yyyy. Ex: 23 ou 2023
	Cvv              string          `json:"cvv,omitempty" valid:""`
	Brand            string          `json:"brand,omitempty" valid:""`
	Label            string          `json:"label,omitempty" valid:""`
	CardId           string          `json:"card_id,omitempty"`
	CardToken        string          `json:"card_token,omitempty"`
	BillingAddressId string          `json:"billing_address_id,omitempty"`
	BillingAddress   *BillingAddress `json:"billing_address"`

	Id             string            `json:"id,omitempty"`
	FirstSixDigits string            `json:"first_six_digits,omitempty"`
	LastFourDigits string            `json:"last_four_digits,omitempty"`
	Status         Status            `json:"status,omitempty"`
	PrivateLabel   bool              `json:"private_label,omitempty"`
	Type           CardType          `json:"type,omitempty"`
	CreatedAt      string            `json:"created_at,omitempty"`
	UpdatedAt      string            `json:"updated_at,omitempty"`
	DeletedAt      string            `json:"deleted_at,omitempty"`
	Customer       *Customer         `json:"customer,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type CardPtr = *Card
type Cards = []CardPtr

func NewCard() *Card {
	return &Card{BillingAddress: &BillingAddress{Country: BrasilCode}}
}

type BillingAddress struct {
	Country string `json:"country" valid:"Required;MaxSize(2)"`
	State   string `json:"state" valid:"Required;MaxSize(2)"`
	City    string `json:"city" valid:"Required;MaxSize(2)"`
	ZipCode string `json:"zip_code" valid:"Required;MaxSize(8)"`
	Line1   string `json:"line_1" valid:"Required"`
	Line2   string `json:"line_2" valid:""`
}

type Boleto struct {
	Bank           BankCode   `json:"bank" valid:"Required"`
	Instructions   string     `json:"instructions" valid:"MaxSize(256)"`
	DueAt          string     `json:"due_at"`
	NossoNumero    string     `json:"nosso_numero"`
	Type           BoletoType `json:"type"`
	DocumentNumber string     `json:"document_number" valid:"MaxSize(16)"` // Identificador do boleto
	Interest       *Interest  `json:"interest"`
	Fine           *Fine      `json:"fine"`
}

func NewBoleto() *Boleto {
	return &Boleto{Interest: new(Interest), Fine: new(Fine)}
}

type Interest struct {
	Days   int64         `json:"days" valid:"Required"`
	Type   OperationType `json:"type" valid:"Required"`
	Amount string        `json:"amount" valid:"Required"` // Valor em porcentagem ou em centavos da taxa de juros que será cobrada ao mês.
}

type Fine struct {
	Days   int64         `json:"days" valid:"Required"`
	Type   OperationType `json:"type" valid:"Required"`
	Amount string        `json:"amount" valid:"Required"` // Valor em porcentagem ou em centavos da taxa de juros que será cobrada ao mês.
}

type Pix struct {
	ExpiresIn             int64                  `json:"expires_in" ` //Data de expiração do Pix em segundos.
	ExpiresAt             string                 `json:"expires_at"`  // Data de expiração do Pix [Formato: YYYY-MM-DDThh:mm:ss] UTC
	AdditionalInformation *AdditionalInformation `json:"additional_information"`
}

func NewPix() *Pix {
	return &Pix{AdditionalInformation: new(AdditionalInformation)}
}

type Plan struct {
	Id             string          `json:"id"`
	Name           string          `json:"name" valid:"Required;MaxSize(64)"`
	Description    string          `json:"description" valid:""`
	Shippable      bool            `json:"shippable"`
	PaymentMethods []PaymentMethod `json:"payment_methods"`
	// Opções de parcelamento disponíveis para assinaturas criadas a partir do plano.
	Installments        []int64  `json:"installments,omitempty"`
	MinimumPrice        int64    `json:"int64,omitempty"`
	StatementDescriptor string   `json:"statement_descriptor" valid:"MaxSize(13)"`
	Currency            Currency `json:"currency"`
	Interval            Interval `json:"interval"`
	/**
	Número de intervalos de acordo com a propriedade interval entre cada cobrança da assinatura.
	Ex.: plano mensal = interval_count (1) e interval (month)
	plano trimestral = interval_count (3) e interval (month)
	plano semestral = interval_count (6) e interval (month)
	*/
	IntervalCount   int64       `json:"interval_count"`
	TrialPeriodDays int64       `json:"trial_period_days,omitempty"`
	BillingType     BillingType `json:"billing_type" valid:"Required"`
	/**
	Dias disponíveis para cobrança das assinaturas criadas a partir do plano.
	Deve ser maior ou igual a 1 e menor ou igual a 28.
	Obrigatório, caso o billing_type seja igual a exact_day.
	*/
	BillingDays   []int64           `json:"billing_days,omitempty"`
	Items         []*PlanItem       `json:"items"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	PricingScheme *PricingScheme    `json:"pricing_scheme,omitempty"`
	/**
	Quantidade para o pricing_scheme.
	Obrigatório quando o pricing_scheme.scheme_type for igual a unit.
	*/
	Quantity int64  `json:"quantity"`
	Status   Status `json:"status,omitempty"`
}

type PlanPtr = *Plan
type Plans = []PlanPtr

func (this *Plan) AddPaymentMethod(methods ...PaymentMethod) *Plan {
	for _, method := range methods {
		this.PaymentMethods = append(this.PaymentMethods, method)
	}
	return this
}

func (this *Plan) SetPricingScheme(scheme *PricingScheme) *Plan {
	this.PricingScheme = scheme
	return this
}

func (this *Plan) AddPlanItem(items ...*PlanItem) *Plan {
	for _, item := range items {
		this.Items = append(this.Items, item)
	}
	return this
}

func (this *Plan) SetIntervalRule(interval Interval, count int64) *Plan {
	this.Interval = interval
	this.IntervalCount = count
	return this
}

func NewPlan(name string) *Plan {
	return &Plan{
		Name:           name,
		Currency:       BRL,
		Installments:   []int64{},
		Items:          []*PlanItem{},
		BillingType:    PrePaid,
		PaymentMethods: []PaymentMethod{}}
}

type PlanItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Quantidade de itens.
	Quantity int64 `json:"quantity"`
	/**
	Número de ciclos durante o qual o item será cobrado.
	Ex: Um item com cycles = 1 representa que um item será cobrado apenas uma vez.
	Caso não seja informado, o item será cobrado até que seja desativado.
	*/
	Cycles        int64          `json:"cycles,omitempty"`
	Interval      string         `json:"interval,omitempty"`
	CreatedAt     string         `json:"created_at,omitempty"`
	UpdatedAt     string         `json:"updated_at,omitempty"`
	DeletedAt     string         `json:"deleted_at,omitempty"`
	PricingScheme *PricingScheme `json:"pricing_scheme,omitempty"`
	Status        Status         `json:"status,omitempty"`
	Plan          *Plan          `json:"plan"`
}

func NewPlanItem(name string, quantity int64, price int64) *PlanItem {
	return &PlanItem{Name: name, Quantity: quantity, PricingScheme: NewPricingScheme(price)}
}

type PricingScheme struct {
	Price        int64      `json:"price"`
	MininumPrice int64      `json:"mininum_price,omitempty"`
	SchemeType   SchemeType `json:"scheme_type,omitempty"`
}

func NewPricingScheme(price int64) *PricingScheme {
	return &PricingScheme{Price: price}
}

type AdditionalInformation struct {
	Name  string `json:"name" valid:"Required"`
	Value string `json:"value" valid:"Required"`
}

type Subscription struct {
	Id                  string              `json:"id,omitempty"`
	Code                string              `json:"code,omitempty" valid:"MaxSize(52)"`
	PlanId              string              `json:"plan_id,omitempty"`
	PaymentMethod       PaymentMethod       `json:"payment_method" valid:"Required"`
	Currency            Currency            `json:"currency"`
	StartAt             string              `json:"start_at" valid:"Required"`
	Interval            Interval            `json:"interval" valid:"Required"`
	MinimumPrice        int64               `json:"minimum_price" valid:""`
	IntervalCount       int64               `json:"interval_count" valid:"Required"`
	BillingType         BillingType         `json:"billing_type" valid:"Required"`
	BillingDay          int64               `json:"billing_day,omitempty" valid:""`
	Installments        int64               `json:"installments" valid:""`
	StatementDescriptor string              `json:"statement_descriptor" valid:"Required;MaxSize(13)"`
	Customer            *Customer           `json:"customer,omitempty"`
	CustomerId          string              `json:"customer_id,omitempty"`
	Discounts           []*Discount         `json:"discounts,omitempty"`
	Increments          []*Increment        `json:"increments,omitempty"`
	Items               []*SubscriptionItem `json:"items"`
	Card                *Card               `json:"card"`
	CardId              string              `json:"card_id"`
	CardToken           string              `json:"card_token"`
	BoletoDueDays       int64               `json:"boleto_due_days"`
	Metadata            map[string]string   `json:"metadata"`
	Description         string              `json:"description"`
	PricingScheme       *PricingScheme      `json:"pricing_scheme,omitempty"`
	Quantity            int64               `json:"quantity"`
	Boleto              *Boleto             `json:"boleto,omitempty"`
}

func NewSubscription() *Subscription {
	return &Subscription{
		PaymentMethod: MethodCreditCard,
		Currency:      BRL,
		Installments:  1,
		BillingType:   PrePaid,
	}
}

func (this *Subscription) WithCustomer(customer CustomerPtr) *Subscription {
	this.Customer = customer
	return this
}

func (this *Subscription) WithCard(card CardPtr) *Subscription {
	this.Card = card
	return this
}

func (this *Subscription) WithPlanId(planId string) *Subscription {
	this.PlanId = planId
	return this
}

func (this *Subscription) SetStartAt(date time.Time) *Subscription {
	this.StartAt = date.Format("2006-01-02")
	return this
}

func (this *Subscription) SetIntervalRule(interval Interval, count int64) *Subscription {
	this.Interval = interval
	this.IntervalCount = count
	return this
}

func (this *Subscription) AddItem(items ...*SubscriptionItem) *Subscription {
	for _, item := range items {
		this.Items = append(this.Items, item)
	}
	return this
}

type SubscriptionPtr = *Subscription
type Subscriptions = []SubscriptionPtr

type SubscriptionItem struct {
	Id            string         `json:"id,omitempty"`
	Description   string         `json:"description,omitempty"`
	Cycles        int64          `json:"cycles,omitempty"`
	Quantity      int64          `json:"quantity"`
	Status        Status         `json:"status,omitempty"`
	CreatedAt     string         `json:"created_at,omitempty"`
	UpdatedAt     string         `json:"updated_at,omitempty"`
	DeletedAt     string         `json:"deleted_at,omitempty"`
	Discounts     []*Discount    `json:"discounts,omitempty"`
	Increments    []*Increment   `json:"increments,omitempty"`
	PricingScheme *PricingScheme `json:"pricing_scheme,omitempty"`
	Name          string         `json:"name"`
}

type SubscriptionItemPtr = *SubscriptionItem
type SubscriptionItems = []SubscriptionItemPtr

func NewSubscriptionItem(description string, quantity int64) *SubscriptionItem {
	return &SubscriptionItem{
		Description: description,
		Quantity:    quantity}
}

type SubscriptionUpdate struct {
	Id            string
	Card          *Card         `json:"card,omitempty"`
	CardId        string        `json:"card_id,omitempty"`
	CardToken     string        `json:"card_token,omitempty"`
	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
}

func NewSubscriptionUpdate(id string) *SubscriptionUpdate {
	return &SubscriptionUpdate{Id: id}
}

func (this *SubscriptionItem) SetPricingScheme(scheme *PricingScheme) *SubscriptionItem {
	this.PricingScheme = scheme
	return this
}

type Discount struct {
	Value         int64         `json":"value"`
	Cycles        int64         `json":"cycles"`
	IncrementType OperationType `json:"increment_type,omitempty"`
}

type Increment struct {
	Value         int64         `json":"value"`
	Cycles        int64         `json":"cycles"`
	IncrementType OperationType `json:"increment_type,omitempty"`
}

type Cycle struct {
	Id        string `json:"id"`
	BillingAt string `json:"billing_at"`
	Cycle     int64  `json:"cycle"`
	StartAt   string `json:"start_at"`
	EndAt     string `json:"end_at"`
	Duration  int64  `json:"duration"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Status    string `json:"status"`
}

type InvoiceItem struct {
	Name        string `json:"name"`
	Amount      int64  `json:"amount"`
	Quantity    int64  `json:"quantity"`
	Description string `json:"description"`
}

type Invoice struct {
	Id              string         `json:"id,omitempty"`
	Url             string         `json:"url,omitempty"`
	Code            string         `json:"code,omitempty"`
	Amount          int64          `json:"amount"`
	PaymentMethod   PaymentMethod  `json:"payment_method"`
	Installments    int64          `json:"installments"`
	Status          InvoiceStatus  `json:"status"`
	BillingAt       string         `json:"billing_at"`
	SeenAt          string         `json:"seen_at"`
	CreatedAt       string         `json:"created_at"`
	CanceledAt      string         `json:"canceled_at"`
	TotalDiscounts  int64          `json:"total_discounts"`
	TotalIncrements int64          `json:"total_increments"`
	Items           []*InvoiceItem `json:"items"`
	Customer        *Customer      `json:"customer"`
	Subscription    *Subscription  `json:"subscription"`
	Cycle           *Cycle         `json:"cycle"`
	Charge          *Charge        `json:"charge"`
}

type InvoicePtr = *Invoice
type Invoices = []InvoicePtr

type ErrorResponse struct {
	Message    string            `json:"message"`
	Errors     map[string]string `json:"errors"`
	StatusCode int64
}

func (this *ErrorResponse) String() string {
	return fmt.Sprint("%v - %v", this.Message, this.Errors)
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Message: msg}
}

func NewErrorResponseWithErrors(msg string, errors map[string]string) *ErrorResponse {
	return &ErrorResponse{Message: msg, Errors: errors}
}

type Paging struct {
	Total int `json:"total"`
}

type Content[T any] struct {
	Paging *Paging `json:"paging"`
	Data   T       `json:"data"`
}

type Response struct {
	Error *ErrorResponse

	Content interface{}

	RawResponse string
	RawRequest  string
	StatusCode  int
}

func NewResponse() *Response {
	return &Response{Error: new(ErrorResponse)}
}

func (this *Response) HasError() bool {
	if len(this.Error.Message) > 0 {
		return true
	}
	if this.Error.Errors != nil {
		return len(this.Error.Errors) > 0
	}
	return false
}

type Charge struct {
	Id              string           `json:"id"`
	Code            string           `json:"code"`
	GatewayId       string           `json:"gateway_id"`
	Amount          int64            `json:"amount"`
	PaidAmount      int64            `json:"paid_amount"`
	Status          ChargeStatus     `json:"status"`
	Currency        Currency         `json:"currency"`
	PaymentMethod   PaymentMethod    `json:"payment_method"`
	DueAt           string           `json:"due_at"`
	PaidAt          string           `json:"paid_at"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
	Customer        *Customer        `json:"customer"`
	LastTransaction *LastTransaction `json:"last_transaction"`
}

type ChargePtr = *Charge
type Charges = []ChargePtr

type ChargeUpdate struct {
	UpdateSubscription bool   `json:"update_subscription,omitempty"`
	CardId             string `json:"card_id,omitempty"`
	CardToken          string `json:"card_token,omitempty"`
	Card               *Card  `json:"card,omitempty"`
}

type Checkout struct {
	// Define a estrutura de acordo com o JSON dos checkouts, se necessário.
}

type LastTransaction struct {
	Id                  string             `json:"id"`
	TransactionType     string             `json:"transaction_type"`
	GatewayId           string             `json:"gateway_id"`
	Amount              int64              `json:"amount"`
	Status              string             `json:"status"`
	Success             bool               `json:"success"`
	Installments        int64              `json:"installments"`
	StatementDescriptor string             `json:"statement_descriptor"`
	AcquirerName        string             `json:"acquirer_name"`
	AcquirerTID         string             `json:"acquirer_tid"`
	AcquirerNSU         string             `json:"acquirer_nsu"`
	AcquirerAuthCode    string             `json:"acquirer_auth_code"`
	AcquirerMessage     string             `json:"acquirer_message"`
	AcquirerReturnCode  string             `json:"acquirer_return_code"`
	OperationType       string             `json:"operation_type"`
	Card                *Card              `json:"card"`
	FundingSource       string             `json:"funding_source"`
	CreatedAt           string             `json:"created_at"`
	UpdatedAt           string             `json:"updated_at"`
	GatewayResponse     *GatewayResponse   `json:"gateway_response"`
	AntifraudResponse   *AntifraudResponse `json:"antifraud_response"`
	Metadata            map[string]string  `json:"metadata"`
}

type GatewayResponse struct {
	Code   string   `json:"code"`
	Errors []string `json:"errors"`
}

type AntifraudResponse struct {
	Status       string `json:"status"`
	Score        string `json:"score"`
	ProviderName string `json:"provider_name"`
}

func failIf(ok bool, msg string, args ...interface{}) {
	if ok {
		fail(msg, args...)
	}
}

func fail(msg string, args ...interface{}) {
	log.Fatalf("error " + fmt.Sprintf(msg, args...))
}