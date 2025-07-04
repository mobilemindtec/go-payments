package v5

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/mobilemindtec/go-payments/api"
	"github.com/mobilemindtec/go-utils/v2/optional"
	"log"
	"strings"
	"time"
)

// Quantidade de vezes para obrar na recorrência
type CycleCount int64

// /Ex .: monthly plan = interval_count (1) and interval month
// /quarterly plan = interval_count (3) and interval month
// /Semi-annual plan = interval_count (6) and interval month
type IntervalCount int64
type Quantity int64

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
type TransferInterval string
type BankAccountType string

type WebhookEvent string
type WebhookStatus string

const (
	BRL            Currency = "BRL"
	DateLayout              = "2006-01-02"
	DateTimeLayout          = "2006-01-02T15:04:05"
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
	Canceled Status = "canceled"
)

const (
	InvoicePending   InvoiceStatus = "pending"
	InvoicePaid      InvoiceStatus = "paid"
	InvoiceCanceled  InvoiceStatus = "canceled"
	InvoiceScheduled InvoiceStatus = "scheduled"
	InvoiceFailed    InvoiceStatus = "failed"
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

const (
	Daily   TransferInterval = "daily"
	Weekly  TransferInterval = "weekly"
	Monthly TransferInterval = "monthly"
)

const (
	Checking BankAccountType = "checking" // corrente
	Savings  BankAccountType = "savings"  // poupança
)

const (
	WebhookPending WebhookStatus = "pending"
	WebhookSent    WebhookStatus = "sent"
	WebhookFailed  WebhookStatus = "failed"
)

const (
	EventCustomerCreated         WebhookEvent = "customer.created"          //	Occurs whenever a customer is created.
	EventCustomerUpdated         WebhookEvent = "customer.updated"          //	Occurs whenever a customer is updated.
	EventCardCreated             WebhookEvent = "card.created"              //	Occurs whenever a card is created.
	EventCardUpdated             WebhookEvent = "card.updated"              //	Occurs whenever a card is updated.
	EventCardDeleted             WebhookEvent = "card.deleted"              //	Occurs whenever a card is deleted.
	EventAddressCreated          WebhookEvent = "address.created"           //	Occurs whenever an address is created.
	EventAddressUpdated          WebhookEvent = "address.updated"           //	Occurs whenever an address is updated.
	EventAddressDeleted          WebhookEvent = "address.deleted"           //	Occurs whenever an address is deleted.
	EventCardExpired             WebhookEvent = "card.expired"              //	Occurs whenever a card expires by the expiration date.
	EventPlanCreated             WebhookEvent = "plan.created"              //	Occurs whenever a plan is created.
	EventPlanUpdated             WebhookEvent = "plan.updated"              //	Occurs whenever a plan is updated.
	EventPlanDeleted             WebhookEvent = "plan.deleted"              //	Occurs whenever a plan is deleted.
	EventPlanItemCreated         WebhookEvent = "plan_item.created"         //	Occurs whenever a plan item is created.
	EventPlanItemUpdated         WebhookEvent = "plan_item.updated"         //	Occurs whenever a plan item is updated.
	EventPlanItemDeleted         WebhookEvent = "plan_item.deleted"         //	Occurs whenever a plan item is deleted.
	EventSubscriptionCreated     WebhookEvent = "subscription.created"      //	Occurs whenever a subscription is created.
	EventSubscriptionCanceled    WebhookEvent = "subscription.canceled"     //	Occurs whenever the subscription is canceled.
	EventSubscriptionItemCreated WebhookEvent = "subscription_item.created" //	Occurs whenever a subscription item is created.
	EventSubscriptionItemUpdated WebhookEvent = "subscription_item.updated" //	Occurs whenever a subscription item is updated.
	EventSubscriptionItemDeleted WebhookEvent = "subscription_item.deleted" //	Occurs whenever a subscription item is deleted.
	EventDiscountCreated         WebhookEvent = "discount.created"          //	Occurs whenever a discount is created.
	EventDiscountDeleted         WebhookEvent = "discount.deleted"          //	Occurs whenever a discount is deleted.
	EventIncrementCreated        WebhookEvent = "increment.created"         //	Occurs whenever an increment is created.
	EventIncrementDeleted        WebhookEvent = "increment.deleted"         //	Occurs whenever an increment is deleted.
	EventOrderPaid               WebhookEvent = "order.paid"                //	Occurs whenever an order is paid.
	EventOrderPaymentFailed      WebhookEvent = "order.payment_failed"      //	Occurs whenever payment for an order fails.
	EventOrderCreated            WebhookEvent = "order.created"             //	Occurs whenever a order is created.
	EventOrderCanceled           WebhookEvent = "order.canceled"            //	Occurs whenever an order is canceled.
	EventOrderItemCreated        WebhookEvent = "order_item.created"        //	Occurs whenever an order item is created.
	EventOrderItemUpdated        WebhookEvent = "order_item.updated"        //	Occurs whenever an order item is updated.
	EventOrderItemDeleted        WebhookEvent = "order_item.deleted"        //	Occurs whenever an order item is deleted.
	EventOrderClosed             WebhookEvent = "order.closed"              //	Occurs whenever a request is closed.
	EventOrderUpdated            WebhookEvent = "order.updated"             //	Occurs whenever a order is updated.
	EventInvoiceCreated          WebhookEvent = "invoice.created"           //	Occurs whenever an invoice is created.
	EventInvoiceUpdated          WebhookEvent = "invoice.updated"           //	Occurs whenever an invoice is updated.
	EventInvoicePaid             WebhookEvent = "invoice.paid"              //	Occurs whenever an invoice is paid.
	EventInvoicePaymentFailed    WebhookEvent = "invoice.payment_failed"    //	Occurs when an invoice payment fails.
	EventInvoiceCanceled         WebhookEvent = "invoice.canceled"          //	Occurs whenever an invoice is canceled
	EventChargeCreated           WebhookEvent = "charge.created"            //	Occurs whenever a charge is created.
	EventChargeUpdated           WebhookEvent = "charge.updated"            //	Occurs when a charge is updated.
	EventChargePaid              WebhookEvent = "charge.paid"               //	Occurs whenever a charge is paid.
	EventChargePaymentFailed     WebhookEvent = "charge.payment_failed"     //	Occurs when a charge for a charge fails.
	EventChargeRefunded          WebhookEvent = "charge.refunded"           //	Occurs whenever a charge is reversed.
	EventChargePending           WebhookEvent = "charge.pending"            //	Occurs whenever a charge is pending.
	EventChargeProcessing        WebhookEvent = "charge.processing"         //	Occurs whenever a charge is still being processed.
	EventChargeUnderpaid         WebhookEvent = "charge.underpaid"          //	Occurs whenever a charge has been underpaid.
	EventChargeOverpaid          WebhookEvent = "charge.overpaid"           //	Occurs whenever a charge has been overpaid.
	EventChargePartialCanceled   WebhookEvent = "charge.partial_canceled"   //	Occurs when a charge has been partially canceled.
	EventChargeChargedback       WebhookEvent = "charge.chargedback"        //	Ocorre sempre que uma cobrança sofre chargeback.
	EventUsageCreated            WebhookEvent = "usage.created"             //	Occurs whenever the usage of an item in the period is created.
	EventUsageDeleted            WebhookEvent = "usage.deleted"             //	Occurs whenever the usage of an item in the period is deleted.
	EventRecipientCreated        WebhookEvent = "recipient.created"         //	Occurs whenever a recipient is created.
	EventRecipientDeleted        WebhookEvent = "recipient.deleted"         //	Occurs whenever a recipient is deleted.
	EventRecipientUpdated        WebhookEvent = "recipient.updated"         //	Occurs whenever a recipient is updated.
	EventBankAccountCreated      WebhookEvent = "bank_account.created"      //	Occurs whenever a bank account is created.
	EventBankAccountUpdated      WebhookEvent = "bank_account.updated"      //	Occurs whenever a bank account is updated.
	EventBankAccountDeleted      WebhookEvent = "bank_account.deleted"      //	Occurs whenever a bank account is deleted.
	EventSellerCreated           WebhookEvent = "seller.created"            //	Occurs whenever a salesperson is created.
	EventSellerUpdated           WebhookEvent = "seller.updated"            //	Occurs whenever a seller is edited.
	EventSellerDeleted           WebhookEvent = "seller.deleted"            //	Occurs whenever a seller is deleted.
	EventTransferPending         WebhookEvent = "transfer.pending"          //	Occurs whenever a transfer is pending.
	EventTransferCreated         WebhookEvent = "transfer.created"          //	Occurs whenever a transfer is created.
	EventTransferProcessing      WebhookEvent = "transfer.processing"       //	Occurs whenever a transfer is in process.
	EventTransferPaid            WebhookEvent = "transfer.paid"             //	It occurs whenever a transfer is paid.
	EventTransferCanceled        WebhookEvent = "transfer.canceled"         //	Occurs whenever a transfer is canceled.
	EventTransferFailed          WebhookEvent = "transfer.failed"           //	Occurs whenever a transfer fails.
	EventCheckoutCreated         WebhookEvent = "checkout.created"          //	Occurs when a checkout is created.
	EventCheckoutCanceled        WebhookEvent = "checkout.canceled"         //	Occurs when a checkout is canceled.
	EventCheckoutClosed          WebhookEvent = "checkout.closed"           //	Ocorre quando um checkout é fechado.
	EventChargeAntifraudApproved WebhookEvent = "charge.antifraud_approved" //	Occurs when an anti-fraud order is approved.
	EventChargeAntifraudReproved WebhookEvent = "charge.antifraud_reproved" //	It occurs when an anti-fraud order is disapproved.
	EventChargeAntifraudManual   WebhookEvent = "charge.antifraud_manual"   //	Occurs when an order in anti-fraud is marked for manual analysis.
	EventChargeAntifraudPending  WebhookEvent = "charge.antifraud_pending"  //	It occurs when an order is pending to be sent for analysis by the anti-fraud service.
)

type WebhookObject[T any] struct {
	Id        string       `json:"id"`
	Account   *Account     `json:"account"`
	Type      WebhookEvent `json:"type"`
	CreatedAt string       `json:"created_at"`
	Data      T            `json:"data"`
}

type Order struct {
	Code             string       `json:"code" valid:"Required;MaxSize(64)"`
	Customer         *Customer    `json:"customer,omitempty"`
	CustomerId       string       `json:"customer_id,omitempty"`
	Items            []*OrderItem `json:"items"`
	Payments         []*Payment   `json:"payments"`
	Closed           bool         `json:"closed"` // Informa se o pedido será criado aberto ou fechado.
	Ip               string       `json:"ip" valid:"Required"`
	SessionId        string       `json:"session_id"`
	AntifraudEnabled bool         `json:"antifraud_enabled"`

	Id     string `json:"id,omitempty"`
	Amount int64  `json:"amount,omitempty"`

	Status    string      `json:"status,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
	UpdatedAt string      `json:"updated_at,omitempty"`
	Charges   []*Charge   `json:"charges,omitempty"`
	Checkouts []*Checkout `json:"checkouts,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func (this *Order) IsCreditCard() bool {
	if len(this.Payments) > 0 {
		return this.Payments[0].IsCreditCard()
	}
	return false
}

func (this *Order) HasCardIdOrToken() bool {
	if this.IsCreditCard() {
		return this.Payments[0].HasCardIdOrToken()
	}
	return false
}

func (this *Order) HasCardToken() bool {
	if this.IsCreditCard() {
		return this.Payments[0].HasCardToken()
	}
	return false
}

func (this *Order) GetCard() CardPtr {
	if len(this.Payments) > 0 && this.Payments[0].IsCreditCard() {
		return this.Payments[0].CreditCard.Card
	}
	return nil
}

// SetCardToken set card token and clean card sensitive information
func (this *Order) SetCardToken(token string) {
	if len(this.Payments) > 0 && this.Payments[0].IsCreditCard() {
		this.Payments[0].CreditCard.Card.Clean()
		this.Payments[0].CreditCard.CardToken = token
	}
}

func (this *Order) HasCardId() bool {
	if this.IsCreditCard() {
		return this.Payments[0].HasCardId()
	}
	return false
}

func (this *Order) GetLastCharge() *optional.Optional[ChargePtr] {
	if len(this.Charges) > 0 {
		return optional.Of[ChargePtr](this.Charges[len(this.Charges)-1])
	}
	return optional.OfNone[ChargePtr]()
}

func (this *Order) GetLastTransaction() *optional.Optional[LastTransactionPtr] {
	if len(this.Charges) > 0 {
		return optional.Just(this.Charges[len(this.Charges)-1].LastTransaction)
	}
	return optional.OfNone[LastTransactionPtr]()
}

func (this *Order) GetStatus() api.PagarmeV5Status {
	return optional.FlatMap[LastTransactionPtr, api.PagarmeV5Status](
		this.GetLastTransaction(),
		func(transaction LastTransactionPtr) *optional.Optional[api.PagarmeV5Status] {
			return optional.Of[api.PagarmeV5Status](transaction.Status)
		},
	).GetOr(api.PagarmeV5None)
}

func (this *Order) GetTranscationId() string {
	return optional.FlatMap[LastTransactionPtr, string](
		this.GetLastTransaction(),
		func(transaction LastTransactionPtr) *optional.Optional[string] {
			return optional.Of[string](transaction.Id)
		},
	).GetOr("")
}

func (this *Order) GetChargeId() string {



	return optional.FlatMap[ChargePtr, string](
		this.GetLastCharge(),
		func(charge ChargePtr) *optional.Optional[string] {
			return optional.Of[string](charge.Id)
		},
	).GetOr("")
}

func (this *Order) GetPayZenSOAPStatus() api.TransactionStatus {
	return optional.FlatMap[LastTransactionPtr, api.TransactionStatus](
		this.GetLastTransaction(),
		func(transaction LastTransactionPtr) *optional.Optional[api.TransactionStatus] {
			return optional.Of[api.TransactionStatus](transaction.GetPayZenSOAPStatus())
		},
	).GetOr(api.NotCreated)
}

type OrderPtr = *Order
type Orders = []OrderPtr

func NewOrder() *Order {
	return &Order{Payments: []*Payment{}, Items: []*OrderItem{}, Metadata: make(map[string]string)}
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

func (this *Order) GetPayment() *Payment {
	if len(this.Payments) > 0 {
		return this.Payments[0]
	}
	return nil
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
	Email      string            `json:"email" valid:"Required"`
	Code       string            `json:"code" valid:"MaxSize(52)"` // código na plataforma
	Document   string            `json:"document" valid:"MaxSize(50)"`
	DocumentType DocumentType `json:"document_type"`
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
	City      string `json:"city" valid:"Required;MaxSize(64)"`
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

func (this *Phone) ToString() string {
	return fmt.Sprintf("%v%v",this.AreaCode, this.Number)
}

type Payment struct {
	PaymentMethod PaymentMethod `json:"payment_method"`
	CreditCard    *CreditCard   `json:"credit_card,omitempty"`
	Boleto        *Boleto       `json:"boleto,omitempty"`
	Pix           *Pix          `json:"pix,omitempty"`
	Amount        int64         `json:"amount" valid:"Required"`
}

func (this *Payment) IsCreditCard() bool {
	return this.PaymentMethod == MethodCreditCard
}

func (this *Payment) HasCardIdOrToken() bool {
	return this.IsCreditCard() &&
		(this.HasCardId() || this.HasCardToken())
}

func (this *Payment) HasCardId() bool {
	return this.IsCreditCard() && len(this.CreditCard.CardId) > 0
}

func (this *Payment) HasCardToken() bool {
	return this.IsCreditCard() && len(this.CreditCard.CardToken) > 0
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
	Card                *Card         `json:"card,omitempty"`
	CardId string `json:"card_id,omitempty"`
	CardToken string `json:"card_token,omitempty"`
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
	//CardToken        string          `json:"card_token,omitempty"`
	BillingAddressId string          `json:"billing_address_id,omitempty"`
	BillingAddress   *BillingAddress `json:"billing_address,omitempty"`

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

func (this *Card) Copy() *Card {
	return &Card{
		Number: this.Number,
		HolderName: this.HolderName,
		HolderDocument: this.HolderDocument,
		ExpMonth: this.ExpMonth,
		ExpYear: this.ExpYear,
		Cvv: this.Cvv,
		Brand: this.Brand,
	}
}

// Clean clean card sensitive information
func (this *Card) Clean(){
	this.Number = ""
	this.HolderName = ""
	this.HolderDocument = ""
	this.ExpMonth = 0
	this.ExpYear = 0
	this.Cvv = ""
	this.Brand = ""
}

type BillingAddress struct {
	Country string `json:"country" valid:"Required;MaxSize(2)"`
	State   string `json:"state" valid:"Required;MaxSize(2)"`
	City    string `json:"city" valid:"Required;MaxSize(64)"`
	ZipCode string `json:"zip_code" valid:"Required;MaxSize(8)"`
	Line1   string `json:"line_1" valid:"Required"`
	Line2   string `json:"line_2" valid:""`
}

type Boleto struct {
	Bank           BankCode   `json:"bank,omitempty" valid:""`
	Instructions   string     `json:"instructions" valid:"MaxSize(256)"`
	DueAt          string     `json:"due_at"`
	NossoNumero    string     `json:"nosso_numero,omitempty"`
	Type           BoletoType `json:"type,omitempty"`
	DocumentNumber string     `json:"document_number,omitempty" valid:"MaxSize(16)"` // Identificador do boleto
	Interest       *Interest  `json:"interest,omitempty"`
	Fine           *Fine      `json:"fine,omitempty"`
}

func NewBoleto() *Boleto {
	return &Boleto{}
}

type Interest struct {
	Days   int64         `json:"days" valid:"Required"`
	Type   OperationType `json:"type" valid:"Required"`
	Amount int64         `json:"amount" valid:"Required"` // Valor em porcentagem ou em centavos da taxa de juros que será cobrada ao mês.
}

type Fine struct {
	Days   int64         `json:"days" valid:"Required"`
	Type   OperationType `json:"type" valid:"Required"`
	Amount int64         `json:"amount" valid:"Required"` // Valor em porcentagem ou em centavos da taxa de juros que será cobrada ao mês.
}

type Pix struct {
	ExpiresIn             int64                  `json:"expires_in,omitempty" ` //Data de expiração do Pix em segundos.
	ExpiresAt             string                 `json:"expires_at,omitempty"`  // Data de expiração do Pix [Formato: YYYY-MM-DDThh:mm:ss] UTC
	AdditionalInformation *AdditionalInformation `json:"additional_information,omitempty"`
}

func NewPix() *Pix {
	return &Pix{}
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
	IntervalCount   IntervalCount `json:"interval_count"`
	TrialPeriodDays int64         `json:"trial_period_days,omitempty"`
	BillingType     BillingType   `json:"billing_type" valid:"Required"`
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
	Quantity Quantity `json:"quantity"`
	Status   Status   `json:"status,omitempty"`
}

func (this *Plan) GetIntervalRule() api.SubscriptionCycle {
	switch this.Interval {
	case Day:
		return api.Daily
	case Week:
		return api.Weekly
	case Month:
		return api.Monthly
	case Year:
		return api.Yearly
	default:
		return api.SubscriptionCycleNone
	}
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

func (this *Plan) SetIntervalRule(interval Interval, count IntervalCount) *Plan {
	this.Interval = interval
	this.IntervalCount = count
	return this
}

func (this *Plan) GetAmount() int64 {
	var amount int64
	for _, it := range this.Items {
		amount += it.PricingScheme.Price * int64(it.Quantity)
	}
	return amount
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
	Quantity Quantity `json:"quantity"`
	//Indicates how many times the item will be charged.
	/**
	Número de ciclos durante o qual o item será cobrado.
	Ex: Um item com cycles = 1 representa que um item será cobrado apenas uma vez.
	Caso não seja informado, o item será cobrado até que seja desativado.
	*/
	Cycles        CycleCount     `json:"cycles,omitempty"`
	Interval      Interval       `json:"interval,omitempty"`
	CreatedAt     string         `json:"created_at,omitempty"`
	UpdatedAt     string         `json:"updated_at,omitempty"`
	DeletedAt     string         `json:"deleted_at,omitempty"`
	PricingScheme *PricingScheme `json:"pricing_scheme,omitempty"`
	Status        Status         `json:"status,omitempty"`
	Plan          *Plan          `json:"plan"`
}

func NewPlanItem(name string, quantity Quantity, cycles CycleCount, price int64) *PlanItem {
	return &PlanItem{
		Name:          name,
		Quantity:      quantity,
		PricingScheme: NewPricingScheme(price),
		Cycles:        cycles,
	}
}

type PricingScheme struct {
	Price        int64      `json:"price"`
	MininumPrice int64      `json:"mininum_price,omitempty"`
	SchemeType   SchemeType `json:"scheme_type,omitempty"`
}

func NewPricingScheme(price int64) *PricingScheme {
	return &PricingScheme{Price: price, SchemeType: Unit}
}

type AdditionalInformation struct {
	Name  string `json:"name" valid:"Required"`
	Value string `json:"value" valid:"Required"`
}

type Subscription struct {
	Id                  string              `json:"id,omitempty"`
	Code                string              `json:"code,omitempty" valid:"MaxSize(52)"`
	PlanId              string              `json:"plan_id,omitempty"`
	Status              Status              `json:"status,omitempty"`
	PaymentMethod       PaymentMethod       `json:"payment_method" valid:"Required"`
	Currency            Currency            `json:"currency"`
	StartAt             string              `json:"start_at" valid:"Required"`
	Interval            Interval            `json:"interval" valid:"Required"`
	MinimumPrice        int64               `json:"minimum_price" valid:""`
	IntervalCount       IntervalCount       `json:"interval_count" valid:"Required"`
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
	Quantity            Quantity            `json:"quantity"`
	Boleto              *Boleto             `json:"boleto,omitempty"`
}

func NewSubscription() *Subscription {
	return &Subscription{
		PaymentMethod: MethodCreditCard,
		Currency:      BRL,
		Installments:  1,
		BillingType:   PrePaid,
		Metadata: make(map[string]string),
	}
}

// SetCardToken set card token and clean card sensitive information
func (this *Subscription) SetCardToken(token string) {
	this.Card.Clean()
	this.CardToken = token
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

func (this *Subscription) SetIntervalRule(interval Interval, intervalCount IntervalCount) *Subscription {
	this.Interval = interval
	this.IntervalCount = intervalCount
	return this
}

func (this *Subscription) SetPricingScheme(scheme *PricingScheme) *Subscription {
	this.PricingScheme = scheme
	return this
}

func (this *Subscription) AddItem(items ...*SubscriptionItem) *Subscription {
	for _, item := range items {
		this.Items = append(this.Items, item)
	}
	return this
}

func (this *Subscription) HasCardIdOrToken() bool {
	return this.PaymentMethod == MethodCreditCard &&
		(len(this.CardId) > 0 || len(this.CardToken) > 0)
}


type SubscriptionPtr = *Subscription
type Subscriptions = []SubscriptionPtr

type SubscriptionItem struct {
	Id            string         `json:"id,omitempty"`
	Description   string         `json:"description,omitempty"`
	Cycles        CycleCount     `json:"cycles,omitempty"` // determina quantas vezes vai ser cobrado. se zero, a recorrência é indeterminada.
	Quantity      Quantity       `json:"quantity"`
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

func NewSubscriptionItem(description string, quantity Quantity, cycles CycleCount) *SubscriptionItem {
	return &SubscriptionItem{
		Description: description,
		Quantity:    quantity,
		Cycles:      cycles,
	}
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

// The object invoice representa os documentos gerados automaticamente ao final
// de cada ciclo de uma assinatura, discriminando todos os valores referentes à
// assinatura, como itens e descontos, para realização da cobrança do assinante.
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
	DueAt           string         `json:"due_at"`
	SeenAt          string         `json:"seen_at"`
	CreatedAt       string         `json:"created_at"`
	CanceledAt      string         `json:"canceled_at"`
	TotalDiscounts  int64          `json:"total_discounts"`
	TotalIncrements int64          `json:"total_increments"`
	SubscriptionId  string         `json:"subscriptionId"`
	Items           []*InvoiceItem `json:"items"`
	Customer        *Customer      `json:"customer"`
	Subscription    *Subscription  `json:"subscription"`
	Cycle           *Cycle         `json:"cycle"`
	Charge          *Charge        `json:"charge"`
}

func (this *Invoice) ToPaymentStatus() api.PaymentStatus {
	switch this.Status {
	case InvoicePending, InvoiceScheduled:
		return api.PaymentWaitingPayment
	case InvoicePaid:
		return api.PaymentPaid
	case InvoiceCanceled:
		return api.PaymentCancelled
	case InvoiceFailed:
		return api.PaymentRefused
	default:
		return api.PaymentError
	}
}

type InvoicePtr = *Invoice
type Invoices = []InvoicePtr

type CardTokenResponse struct {
	Id string `json:"id"`
	Type string `json:"type"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
	Card *Card `json:"card"`
}

type  CardTokenResponsePtr = *CardTokenResponse

type ErrorResponse struct {
	Message    string              `json:"message"`
	Errors     map[string][]string `json:"errors"`
	StatusCode int64
}

func (this *ErrorResponse) String() string {
	return fmt.Sprint("State: %v: %v - %v", this.StatusCode, this.Message, this.Errors)
}

func (this *ErrorResponse) ToMapOfString() map[string]string {
	errors := make(map[string]string)
	if this.Errors != nil {
		for k, v := range this.Errors {
			errors[k] = strings.Join(v, ", ")
		}
	}
	return errors
}

func (this *ErrorResponse) Error() string {
	return this.String()
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Message: msg}
}

func NewErrorResponseWithErrors(msg string, errors map[string][]string) *ErrorResponse {
	return &ErrorResponse{Message: msg, Errors: errors}
}

type Paging struct {
	Total int `json:"total"`
}

type Content[T any] struct {
	Paging *Paging `json:"paging"`
	Data   T       `json:"data"`
}

type Success[T any] struct {
	Data        T
	RawResponse string
	RawRequest  string
}

func NewSuccess[T any](response *Response) *Success[T] {
	return &Success[T]{
		Data:        response.Content.(T),
		RawRequest:  response.RawRequest,
		RawResponse: response.RawResponse,
	}
}

func NewSuccessWithValue[T any](response *Response, val T) *Success[T] {
	return &Success[T]{
		Data:        val,
		RawRequest:  response.RawRequest,
		RawResponse: response.RawResponse,
	}
}

func NewSuccessSlice[T any](response *Response) *Success[T] {
	return &Success[T]{
		Data:        response.Content.(*Content[T]).Data,
		RawRequest:  response.RawRequest,
		RawResponse: response.RawResponse,
	}
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

func (this *Response) HasErrors() bool {
	if this.Error.Errors != nil {
		return len(this.Error.Errors) > 0
	}
	return false
}

func (this *Response) GetErrros() map[string][]string {
	if this.Error.Errors != nil {
		return this.Error.Errors
	}
	return make(map[string][]string)
}

func (this *Response) GetMessage() string {
	if len(this.Error.Message) > 0 {
		return this.Error.Message
	}
	return ""
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
	Invoice         *Invoice         `json:"invoice"`
	Customer        *Customer        `json:"customer"`
	LastTransaction *LastTransaction `json:"last_transaction"`
}

func (this *Charge) ToPaymentStatus() api.PaymentStatus {
	switch this.Status {
	case ChargePending:
		return api.PaymentWaitingPayment
	case ChargePaid:
		return api.PaymentPaid
	case ChargeCanceled:
		return api.PaymentCreated
	case ChargeProcessing:
		return api.PaymentWaitingPayment
	case ChargeFailed:
		return api.PaymentRefused
	case ChargeOverpaid, ChargeUnderpaid:
		return api.PaymentOther
	default:
		return api.PaymentOther

	}
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
	Id                  string              `json:"id"`
	TransactionType     string              `json:"transaction_type"`
	GatewayId           string              `json:"gateway_id"`
	Amount              int64               `json:"amount"`
	PaidAmount              int64               `json:"paid_amount"`
	Status              api.PagarmeV5Status `json:"status"`
	Success             bool                `json:"success"`
	Installments        int64               `json:"installments"`
	StatementDescriptor string              `json:"statement_descriptor"`
	AcquirerName        string              `json:"acquirer_name"`
	AcquirerTID         string              `json:"acquirer_tid"`
	AcquirerNSU         string              `json:"acquirer_nsu"`
	AcquirerAuthCode    string              `json:"acquirer_auth_code"`
	AcquirerMessage     string              `json:"acquirer_message"`
	AcquirerReturnCode  string              `json:"acquirer_return_code"`
	OperationType       string              `json:"operation_type"`
	Card                *Card               `json:"card"`
	FundingSource       string              `json:"funding_source"`
	CreatedAt           string              `json:"created_at"`
	UpdatedAt           string              `json:"updated_at"`
	DueAt           string              `json:"due_at"`
	GatewayResponse     *GatewayResponse    `json:"gateway_response"`
	AntifraudResponse   *AntifraudResponse  `json:"antifraud_response"`
	Metadata            map[string]string   `json:"metadata"`

	Url         string `json:"url,omitempty"`
	Pdf         string `json:"pdf,omitempty"`
	Line        string `json:"line,omitempty"`
	Barcode     string `json:"barcode,omitempty"`
	QrCode      string `json:"qr_code,omitempty"`
	QrCodeUrl   string `json:"qr_code_url,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"`
	NossoNumero string `json:"nosso_numero,omitempty"`
}

func (this *LastTransaction) GetPayZenSOAPStatus() api.TransactionStatus {
	switch this.Status {
	case api.PagarmeV5Generated:
		return api.Created
	case api.PagarmeV5AuthorizedPendingCapture, api.PagarmeV5WaitingCapture:
		return api.Authorised
	case api.PagarmeV5Captured, api.PagarmeV5PartialCapture:
		return api.Captured
	case api.PagarmeV5Refunded, api.PagarmeV5PartialRefunded:
		return api.Refunded
	case api.PagarmeV5WaitingPayment:
		return api.WaitingPayment
	case api.PagarmeV5PendingRefund:
		return api.PendingRefund
	case api.PagarmeV5NotAuthorized:
		return api.Refused
	case api.PagarmeV5Voided, api.PagarmeV5PartialVoid:
		return api.Canceled
	case api.PagarmeV5Paid:
		return api.Authorised
	case api.PagarmeV5None:
		return api.NotCreated
	case api.PagarmeV5Chargedback:
		return api.Chargeback
	default:
		return api.Error
	}
}

type LastTransactionPtr = *LastTransaction

type TransferSettings struct {
	TransferEnabled  bool             `json:"transfer_enabled"`
	TransferInterval TransferInterval `json:"transfer_interval"`
	TransferDay      int64            `json:"transfer_day,omitempty"`
}

type BankAccount struct {
	HolderName        string          `json:"holder_name" valid:"Required;MaxSize(30)"`
	Bank              string          `json:"bank" valid:"Required;MaxSize(3)"`
	BranchNumber      string          `json:"branch_number" valid:"Required;MaxSize(4)"`
	BranchCheckDigit  string          `json:"branch_check_digit" valid:"Required;MaxSize(1)"`
	AccountNumber     string          `json:"account_number" valid:"Required;MaxSize(13)"`
	AccountCheckDigit string          `json:"account_check_digit" valid:"Required;MaxSize(2)"`
	HolderType        CustomerType    `json:"holder_type" valid:"Required"`
	HolderDocument    string          `json:"holder_document" valid:"Required"`
	Type              BankAccountType `json:"type" valid:"Required"`
}

// The object recipient represents a receiver,
// which will receive part of the sale made.
type Recipient struct {
	Id                 string            `json:"id,omitempty"`
	Name               string            `json:"name" valid:"Required;MaxSize(128)"`
	Email              string            `json:"email" valid:"Required;MaxSize(64)"`
	Description        string            `json:"description" valid:"Required;MaxSize(256)"`
	Document           string            `json:"document" valid:"Required;MaxSize(16)"`
	Type               CustomerType      `json:"type" valid:"Required"`
	Code               string            `json:"code" valid:"Required"`
	TransferSettings   *TransferSettings `json:"transfer_settings" valid:"Required"`
	DefaultBankAccount BankAccount       `json:"default_bank_account" valid:"Required"`
	Metadata           map[string]string `json:"metadata"`
}

type RecipientPtr = *Recipient
type Recipients = []RecipientPtr

type RecipientUpdate struct {
	Email       string       `json:"email" valid:"Required;MaxSize(64)"`
	Description string       `json:"description" valid:"Required;MaxSize(256)"`
	Type        CustomerType `json:"type" valid:"Required"`
}

type Balance struct {
	Currency           string `json:"currency"`
	AvailableAmount    int64  `json:"available_amount"`
	WaitingFundsAmount int64  `json:"waiting_funds_amount"`
	TransferredAmount  int64  `json:"transferred_amount"`
}

type BalancePtr = *Balance

type BalanceOperation struct {
	Id             string `json:"id"`
	Status         string `json:"status"`
	BalanceAmount  int64  `json:"balance_amount"`
	Type           string `json:"type"`
	Amount         int64  `json:"amount"`
	Fee            int64  `json:"fee"`
	CreatedAt      string `json:"created_at"`
	MovementObject string `json:"movement_object"`
}

type BalanceOperationPtr = *BalanceOperation
type BalanceOperations = []BalanceOperationPtr

type MovementObject struct {
	Fee               int64  `json:"fee"`
	AnticipationFee   int64  `json:"anticipation_fee"`
	FraudCoverageFee  int64  `json:"fraud_coverage_fee"`
	RecipientId       string `json:"recipient_id"`
	OriginatorModel   string `json:"originator_model"`
	OriginatorModelId string `json:"originator_model_id"`
	PaymentDate       string `json:"payment_date"`
	PaymentMethod     string `json:"payment_method"`
	Object            string `json:"object"`
	Id                string `json:"id"`
	Status            string `json:"status"`
	Amount            int64  `json:"amount"`
	CreatedAt         string `json:"created_at"`
	Type              string `json:"type"`
	GatewayId         string `json:"gateway_id"`
}

type Transfer struct {
	Id          string       `json:"id"`
	Amount      int64        `json:"amount"`
	Status      string       `json:"status"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
	BankAccount *BankAccount `json:"bank_account"`
	Recipient   *Recipient   `json:"recipient"`
}

type TransferPtr = *Transfer
type Transfers = []TransferPtr

type GatewayResponse struct {
	Code   string              `json:"code"`
	Errors []map[string]string `json:"errors"`
}

type AntifraudResponse struct {
	Status       string `json:"status"`
	Score        string `json:"score"`
	ProviderName string `json:"provider_name"`
}

type Account struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ChargePaymentConfirm struct {
	ChargeCode  string `json:"charge_code"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
}

func failIf(ok bool, msg string, args ...interface{}) {
	if ok {
		fail(msg, args...)
	}
}

func fail(msg string, args ...interface{}) {
	log.Fatalf("error " + fmt.Sprintf(msg, args...))
}
