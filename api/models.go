package api

import (
	"fmt"
	"time"
	"github.com/beego/beego/v2/core/logs"
)

type Gateway string
type AsaasMode int64
type PayZenApiVersion string
type PagarmeApiVersion string

const (
	GatewayNone    Gateway = ""
	GatewayPagarme Gateway = "Pagarme"
	GatewayPayZen  Gateway = "PayZen"
	GatewayAsaas   Gateway = "Asaas"
	GatewayPicPay  Gateway = "PicPay"

	// Asaas mode
	AsaasModeProd AsaasMode = 1
	AsaasModeTest AsaasMode = 2

	// API version
	PagarmeApi5     PagarmeApiVersion = "5"
	PayZenApiRestV4 PayZenApiVersion  = "RESTFul.v4"

	// Api mode
	ApiModeTest       = "TEST"
	ApiModeProduction = "PRODUCTION"

	// scheme
	SchemeBoleto = "BOLETO"
)

type PaymentStatus int

const (
	PaymentInitial PaymentStatus = 1 + iota
	PaymentWaitingPayment
	PaymentPaid
	PaymentRefused
	PaymentCancelled
	PaymentRefound
	PaymentExpired
	PaymentChargeback
	PaymentOther
	PaymentSuccess
	PaymentError
	PaymentCreated
)

type PaymentStatusLabel string

const (
	PaymentInitialLabel        PaymentStatusLabel = "initial"
	PaymentWaitingPaymentLabel PaymentStatusLabel = "waiting_payment"
	PaymentPaidLabel           PaymentStatusLabel = "paid"
	PaymentRefusedLabel        PaymentStatusLabel = "refused"
	PaymentCancelledLabel      PaymentStatusLabel = "cancelled"
	PaymentRefoundLabel        PaymentStatusLabel = "refound"
	PaymentExpiredLabel        PaymentStatusLabel = "expired"
	PaymentChargebackLabel     PaymentStatusLabel = "changeback"
	PaymentOtherLabel          PaymentStatusLabel = "other"
	PaymentSuccessLabel        PaymentStatusLabel = "success"
	PaymentErrorLabel          PaymentStatusLabel = "error"
	PaymentCreatedLabel        PaymentStatusLabel = "created"
)

type PaymentEvent string

const (

	/*

		Asaas Events

		PAYMENT_CREATED - Geração de nova cobrança.
		PAYMENT_UPDATED - Alteração no vencimento ou valor de cobrança existente.
		PAYMENT_CONFIRMED - Cobrança confirmada (pagamento efetuado, porém o saldo ainda não foi disponibilizado).
		PAYMENT_RECEIVED - Cobrança recebida.
		PAYMENT_OVERDUE - Cobrança vencida.
		PAYMENT_DELETED - Cobrança removida.
		PAYMENT_RESTORED - Cobrança restaurada.
		PAYMENT_REFUNDED - Cobrança estornada.
		PAYMENT_RECEIVED_IN_CASH_UNDONE - Recebimento em dinheiro desfeito.
		PAYMENT_CHARGEBACK_REQUESTED - Recebido chargeback.
		PAYMENT_CHARGEBACK_DISPUTE - Em disputa de chargeback (caso sejam apresentados documentos para contestação).
		PAYMENT_AWAITING_CHARGEBACK_REVERSAL - Disputa vencida, aguardando repasse da adquirente.
		PAYMENT_DUNNING_RECEIVED - Recebimento de recuperação.
		PAYMENT_DUNNING_REQUESTED - Requisição de recuperação.
	*/

	PaymentEventCreated                    PaymentEvent = "PAYMENT_CREATED"
	PaymentEventUpdated                    PaymentEvent = "PAYMENT_UPDATED"
	PaymentEventConfirmed                  PaymentEvent = "PAYMENT_CONFIRMED"
	PaymentEventReceived                   PaymentEvent = "PAYMENT_RECEIVED"
	PaymentEventOverdue                    PaymentEvent = "PAYMENT_OVERDUE"
	PaymentEventDeleted                    PaymentEvent = "PAYMENT_DELETED"
	PaymentEventRestored                   PaymentEvent = "PAYMENT_RESTORED"
	PaymentEventRefunded                   PaymentEvent = "PAYMENT_REFUNDED"
	PaymentEventReceivedInCashUndone       PaymentEvent = "PAYMENT_RECEIVED_IN_CASH_UNDONE"
	PaymentEventChargebackRequested        PaymentEvent = "PAYMENT_CHARGEBACK_REQUESTED"
	PaymentEventChargebackDispute          PaymentEvent = "PAYMENT_CHARGEBACK_DISPUTE"
	PaymentEventAwaitingChargebackReversal PaymentEvent = "PAYMENT_AWAITING_CHARGEBACK_REVERSAL"
	PaymentEventDunningReceived            PaymentEvent = "PAYMENT_DUNNING_RECEIVED"
	PaymentEventDunningRequested           PaymentEvent = "PAYMENT_DUNNING_REQUESTED"
	PaymentEventPaymentStatusChanged       PaymentEvent = "PAYMENT_STATUS_CHANGED"
	PaymentEventSubscriptionStatusChanged  PaymentEvent = "SUBSCRIPTION_STATUS_CHANGED"
	PaymentEventRecipientStatusChanged     PaymentEvent = "RECIPIENT_STATUS_CHANGED"
	PaymentEventTransactionCreated         PaymentEvent = "TRANSACTION_CREATED"
	PaymentEventPicPayStatusChanged        PaymentEvent = "PICPAY_STATUS_CHANGED"
	PaymentEventNotFound                   PaymentEvent = "EVENT_NOT_FOUND"

	PaymentEventOrderChanged PaymentEvent = "ORDER_CHANGED"

	//
)

type TransactionStatus int

const (
	Initial TransactionStatus = 1 + iota
	NotCreated
	Authorised
	AuthorisedToValidate
	WaitingAuthorisation
	WaitingAuthorisationToValidate
	Refused
	Captured
	Canceled
	Expired
	UnderVerification
	PartiallyAuthorised
	Refunded       //pickpay,pagarme, asaas
	Created        //pickpay, pagarme
	Chargeback     //pickpay,pagarme,asaas
	WaitingPayment //pagarme, asaas
	PendingRefund  //pagarme, asaas
	Analyzing      //pagarme
	PendingReview  //pagarme
	ReceivedInCash // asaas
	Other
	Canceled__invalid
	Success
	Error
)

// used in payzen
type SubscriptionCycle string

const (
	SubscriptionCycleNone SubscriptionCycle = ""
	Daily                 SubscriptionCycle = "DAILY"        // diario
	Weekly                SubscriptionCycle = "WEEKLY"       // semanal
	Biweekly              SubscriptionCycle = "BIWEEKLY"     // quinzenal
	Monthly               SubscriptionCycle = "MONTHLY"      // mensal
	Quarterly             SubscriptionCycle = "QUARTERLY"    // trimestral
	Semiannually          SubscriptionCycle = "SEMIANNUALLY" // semestral
	Yearly                SubscriptionCycle = "YEARLY"       // anual
)

type PaymentType string

const (
	PaymentTypeNone       PaymentType = ""
	PaymentTypeCreditCard PaymentType = "credit_card"
	PaymentTypeDebitCard  PaymentType = "debit_card"
	PaymentTypeBoleto     PaymentType = "boleto"
	PaymentTypePix        PaymentType = "pix"
	PaymentTypePicPay     PaymentType = "picpay"
	PaymentTypeTransfer   PaymentType = "transfer"
	PaymentTypeDeposit    PaymentType = "deposit"
	PaymentTypeUndefined  PaymentType = "undefined"
)

type BankAccountType string

const (
	BankAccountTypeNone BankAccountType = ""
	ContaCorrente       BankAccountType = "CONTA_CORRENTE"
	ContaPoupanca       BankAccountType = "CONTA_POUPANCA"
)

type TransferStatus string

const (
	TransferPending        TransferStatus = "PENDING"
	TransferBankProcessing TransferStatus = "BANK_PROCESSING"
	TransferDone           TransferStatus = "DONE"
	TransferCancelled      TransferStatus = "CANCELLED"
	TransferFailed         TransferStatus = "FAILED"
)

type Filter struct {
	Limit      int64  `jsonp:""`
	Offset     int64  `jsonp:""`
	StartDate  string `jsonp:""`
	FinishDate string `jsonp:""`

	DateCreated string `jsonp:""` // transfer filter asaas

	Status        string `jsonp:""` // pagarme
	BankAccountId string `jsonp:""` //pagarme

	RecebedorId string `jsonp:""`
}

func NewFilter() *Filter {
	return &Filter{}
}

type Bank struct {
	Code string `jsonp:""`
}

func NewBank() *Bank {
	return &Bank{}
}

type BankAccount struct {
	Bank        *Bank  `jsonp:""`
	AccountName string `jsonp:""`
	OwnerName   string `jsonp:""`
	//Data de nascimento do proprietário da conta.
	//Somente quando a conta bancária não pertencer ao mesmo CPF ou CNPJ da conta Asaas.
	OwnerBirthDate  string          `jsonp:""`
	CpfCnpj         string          `jsonp:""`
	Agency          string          `jsonp:""`
	Account         string          `jsonp:""`
	AccountDigit    string          `jsonp:""`
	BankAccountType BankAccountType `jsonp:""`
}

func NewBankAccount() *BankAccount {
	return &BankAccount{}
}

type Transfer struct {
	Amount        float64      `jsonp:""`
	BankAccountId int64        `jsonp:""`
	BankAccount   *BankAccount `jsonp:""`
	RecebedorId   string       `jsonp:""`
}

func NewTransfer() *Transfer {
	return &Transfer{}
}

type BoletoFine struct {
	Days       int64   `jsonp:""`
	Amount     float64 `jsonp:""`
	Percentage float64 `jsonp:""`
}

type BoletoInterest struct {
	Days       int64   `jsonp:""`
	Amount     float64 `jsonp:""`
	Percentage float64 `jsonp:""`
}

type BoletoDiscount struct {
	Days       int64   `jsonp:""`
	Amount     float64 `jsonp:""`
	Percentage float64 `jsonp:""`
}

type Plan struct {
	Id              string        `jsonp:""`
	Amount          float64       `jsonp:""`
	Days            int64         `jsonp:""` // Prazo, em dias, para cobrança das parcelas
	Name            string        `jsonp:""`
	TrialDays       int64         `jsonp:""`
	PaymentMethods  []PaymentType `jsonp:""`
	Charges         int64         `jsonp:""`
	InvoiceReminder int64         `jsonp:""`
	Installments    int64         `jsonp:""`

	IntervalRule  SubscriptionCycle `jsonp:""` // usado no pagarme v5
	IntervalCount int64             `jsonp:""` // usado no pagarme v5
}

func NewPlan() *Plan {
	return &Plan{PaymentMethods: []PaymentType{}}
}

type Card struct {
	Id             string `valid:"" jsonp:""`
	Number         string `valid:"" jsonp:""`
	Scheme         string `valid:"" jsonp:"brand"`
	ExpiryMonth    string `valid:"" jsonp:"expiry_month"`
	ExpiryYear     string `valid:"" jsonp:"expiry_year"`
	SecurityCode   string `valid:"" jsonp:"cvv"`
	HolderName     string `valid:"" jsonp:"holder_name"`
	HolderDocument string `valid:"" jsonp:"holder_document"`
	Token          string `valid:"" jsonp:""`
}

type Customer struct {
	Id                   string `jsonp:""`
	FirstName            string `valid:"Required" jsonp:""`
	LastName             string `valid:"" jsonp:""`
	PhoneNumber          string `valid:"" jsonp:""`
	CellPhoneNumber      string `valid:"" jsonp:""`
	Address2             string `valid:"" jsonp:""`
	Email                string `valid:"Required" jsonp:""`
	StreetNumber         string `valid:"" jsonp:""`
	Address              string `valid:"" jsonp:""`
	District             string `valid:"" jsonp:""`
	ZipCode              string `valid:"" jsonp:""`
	City                 string `valid:"" jsonp:""`
	State                string `valid:"" jsonp:""`
	Country              string `valid:"Required" jsonp:""`
	Document             string `valid:"Required" jsonp:"document"`
	ExternalReference    string `jsonp:""`
	NotificationDisabled bool   `jsonp:""`
	IpAddress            string `jsonp:""`
}

func NewCustomer() *Customer {
	return &Customer{}
}

func (this *Customer) IsCreated() bool {
	return len(this.Id) > 0
}

type Subscription struct {
	OrderId        string `valid:"Required" jsonp:""`
	SubscriptionId string `valid:"Required" jsonp:""`
	TransactionId  string `valid:"" jsonp:""`
	Description    string `valid:"Required" jsonp:""`

	// Pagarme v5
	SubscriptionItemId string `valid:"" jsonp:""`

	// valor da recorrência
	Amount float64 `valid:"Required" jsonp:""`

	Installments int64 `jsonp:""`

	// valor inicial da recorrência
	InitialAmount float64 `valid:"" jsonp:""`
	// quantas vezes o valor inicial deve ser cobrado
	InitialAmountNumber int64 `valid:"" jsonp:""`

	// data de inicio da cobrança
	EffectDate time.Time `valid:"Required" jsonp:""`
	// cobrar no último dia do mês

	// pagarme, mês, anos, etc..
	IntervalCount int64 `jsonp:""`

	// Quantidade de cobranças
	Count int64 `jsonp:""`

	Cycle                   SubscriptionCycle `jsonp:""`
	PaymentAtLastDayOfMonth bool              `jsonp:""`
	PaymentAtDayOfMonth     int64             `jsonp:""`
	SoftDescriptor          string            `valid:"" jsonp:""`

	// regra de recorrência do payzen
	Rule string `jsonp:""`

	Token string `valid:"" jsonp:""`

	BoletoFine     *BoletoFine            `jsonp:""`
	BoletoInterest *BoletoInterest        `jsonp:""`
	BoletoDiscount *BoletoDiscount        `jsonp:""`
	Card           *Card                  `valid:"" jsonp:""`
	PostbackUrl    string                 `jsonp:""`
	WebhookUrl     string                 `jsonp:""`
	PaymentType    PaymentType            `jsonp:""` // pagarme/asaas
	PlanId         string                 `jsonp:""` // pagarme
	Customer       *Customer              `jsonp:""`
	AdditionalInfo map[string]interface{} `jsonp:""`

	UpdatePendingPayments bool

	Order *Order `jsonp:""`
}

func NewSubscription() *Subscription {
	subscription := new(Subscription)
	return subscription
}

type Payment struct {
	OrderId      string    `valid:"" jsonp:""`
	Installments int64     `valid:"" jsonp:""`
	Amount       float64   `valid:"" jsonp:""`
	Customer     *Customer `valid:"" jsonp:""`
	Card         *Card     `valid:"" jsonp:""`

	//TokenOperation bool
	TransactionId string `jsonp:""`

	PaymentType PaymentType `jsonp:""`

	// usado no pagarme v1 e payzen para sinalizar onde a resposta deve ser enviada
	// na ausencia desse valor, o payments vai entrar com uma URL conhecida que processará
	// o postback e depois enviará a informação para a URL fornecida em WebhookUrl
	PostbackUrl    string `valid:"" jsonp:""`
	SoftDescriptor string `valid:"" jsonp:""`
	//Metadata map[string]interface{} `valid:"" jsonp:"pagarme_metadata"`
	SaveBoletoAtPath string
	BoletoFine       *BoletoFine     `jsonp:""`
	BoletoInterest   *BoletoInterest `jsonp:""`
	BoletoDiscount   *BoletoDiscount `jsonp:""`

	// usado para envio de postback payment -> consumidor api (4gym,mobloja,etc..)
	WebhookUrl string `jsonp:""`

	//PicPay
	ReturnUrl      string                 `json:"" jsonp:"picpay_return_url"`
	Plugin         string                 `json:"" jsonp:"picpay_plugin"`
	AdditionalInfo map[string]interface{} `json:"" jsonp:""`

	ExpiresAt time.Time `json:"" jsonp:"expires_at"`

	BoletoInstructions string `jsonp:"boleto_instructions"`

	Order *Order `jsonp:""`
}

func NewPayment() *Payment {
	payment := new(Payment)
	payment.Customer = new(Customer)
	payment.Card = new(Card)
	payment.Installments = 1
	payment.AdditionalInfo = make(map[string]interface{})
	payment.Order = NewOrder()
	return payment
}

type PaymentToken struct {
	Token string `valid:"Required" jsonp:""`
}

func NewPaymentToken() *PaymentToken {
	tokenPayment := new(PaymentToken)
	return tokenPayment
}

type PaymentFind struct {
	TransactionId string `jsonp:""`
	//TransactionUuid string  `jsonp:""`
	OrderId         string  `jsonp:""`
	SubscriptionId  string  `jsonp:""`
	Amount          float64 `valid:"" jsonp:""`
	Token           string  `jsonp:""`
	AuthorizationId string  `jsonp:""`

	// Consulta
	//ChargeId   string `jsonp:""`
	ChargeCode string `jsonp:""`
	Size       int    `jsonp:""` // results size
	Page       int    `jsonp:""` // results page

	PaymentType PaymentType `jsonp:""` // picpay, pix

	CustomerDocument string `jsonp:""`

	CustomerExternalReference string `jsonp:""`
}

func NewPaymentFind() *PaymentFind {
	find := new(PaymentFind)
	return find
}

type OrderItem struct {
	Id          int64   `jsonp:""`
	Description string  `jsonp:""`
	Type        string  `jsonp:""`
	Reference   string  `jsonp:""`
	Quantity    int64   `jsonp:""`
	Amount      float64 `jsonp:""`
}

func NewOrderItem() *OrderItem {
	return &OrderItem{}
}

type OrderDeliveryAddress struct {
	StreetNumber  string `jsonp:""`
	Address       string `jsonp:""`
	Address2      string `jsonp:""`
	District      string `jsonp:""`
	ZipCode       string `jsonp:""`
	City          string `jsonp:""`
	State         string `jsonp:""`
	Country       string `jsonp:""`
	ReclaimInShop bool   `jsonp:""`
}

func NewOrderDeliveryAddress() *OrderDeliveryAddress {
	return &OrderDeliveryAddress{}
}

type Order struct {
	Id              int64                 `jsonp:""`
	DeliveryCost    float64               `jsonp:""`
	FirstName       string                `jsonp:""`
	LastName        string                `jsonp:""`
	PhoneNumber     string                `jsonp:""`
	DeliveryAddress *OrderDeliveryAddress `jsonp:""`
	Items           []*OrderItem          `jsonp:""`
}

func NewOrder() *Order {
	return &Order{DeliveryAddress: NewOrderDeliveryAddress(), Items: []*OrderItem{}}
}

type TokenInfo struct {
	Id               string `jsonp:""`
	Token            string `jsonp:""`
	Number           string `jsonp:""`
	Brand            string `jsonp:""`
	CreationDate     time.Time
	CancellationDate time.Time
	Cancelled        bool   `jsonp:""`
	Active           bool   `jsonp:""`
	NotFound         bool   `jsonp:""`
	FirstSixDigits   string `jsonp:""`
	LastFourDigits   string `jsonp:""`
}

type SubscriptionResult struct {
	SubscriptionId      string    `jsonp:""`
	SubscriptionItemId  string    `jsonp:""`
	PastPaymentsNumber  int64     `jsonp:""`
	TotalPaymentsNumber int64     `jsonp:""`
	EffectDate          time.Time `jsonp:""`
	CancelDate          time.Time `jsonp:""`
	InitialAmountNumber int64     `jsonp:""`
	InitialAmount       int64     `jsonp:""`
	Rule                string    `jsonp:""`
	Description         string    `jsonp:""`

	Active    bool `jsonp:""`
	Cancelled bool `jsonp:""`
	Started   bool `jsonp:""`

	Token string `jsonp:""`
}

type TransactionItemResult struct {
	//TransactionUuid   string            `jsonp:""`
	TransactionId     string            `jsonp:""`
	AuthorizationId   string            `jsonp:""`
	CancellationId    string            `jsonp:""`
	PagarmeV5Status   PagarmeV5Status   `jsonp:"pagarme_v5_status"`
	PicPayStatus      PicPayStatus      `jsonp:"picpay_status"`
	TransactionStatus TransactionStatus `jsonp:"payzen_status"`
	PayZenV4Status    string            `jsonp:"payzen_v4_status"`
	AsaasStatus       AsaasStatus       `jsonp:"asaas_status"`

	TransactionStatusLabel string `jsonp:""`
	//ExternalTransactionId  string             `jsonp:""`
	Nsu     string    `jsonp:""`
	Amount  float64   `jsonp:""`
	DueDate time.Time `jsonp:""`
	//ExpectedCaptureDate    time.Time          `jsonp:""`
	CreationDate           time.Time          `jsonp:""`
	Status      PaymentStatus      `jsonp:""`
	StatusLabel PaymentStatusLabel `jsonp:""`

	StatusText string `jsonp:""`

	Platform Gateway `jsonp:""`
}

func NewTransactionItemResult(platform Gateway) *TransactionItemResult {
	return &TransactionItemResult{Platform: platform}
}

func (this *TransactionItemResult) isPagarme() bool {
	return this.Platform == GatewayPagarme
}

func (this *TransactionItemResult) isPayZen() bool {
	return this.Platform == GatewayPayZen
}

func (this *TransactionItemResult) isPicPay() bool {
	return this.Platform == GatewayPicPay
}

func (this *TransactionItemResult) isAsaas() bool {
	return this.Platform == GatewayAsaas
}

func (this *TransactionItemResult) IsPaid() bool {
	return this.TransactionStatus == Authorised || this.TransactionStatus == Captured
}

func (this *TransactionItemResult) IsCancelled() bool {
	return this.TransactionStatus == Canceled
}

func (this *TransactionItemResult) IsRefound() bool {
	return this.TransactionStatus == Refunded
}

func (this *TransactionItemResult) IsRefused() bool {
	return this.TransactionStatus == Refused
}

func (this *TransactionItemResult) BuildStatus() {

	switch this.TransactionStatus {
	case Initial:
		this.Status = PaymentInitial
		this.StatusLabel = PaymentInitialLabel
		break
	case NotCreated:
		this.Status = PaymentInitial
		this.StatusLabel = PaymentInitialLabel
		break
	case Authorised:
		if this.isPayZen() {
			this.Status = PaymentPaid
			this.StatusLabel = PaymentPaidLabel
		} else {
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
		}
		break
	case AuthorisedToValidate:
		this.Status = PaymentWaitingPayment
		this.StatusLabel = PaymentWaitingPaymentLabel
		break
	case WaitingAuthorisation:
		this.Status = PaymentWaitingPayment
		this.StatusLabel = PaymentWaitingPaymentLabel
		break
	case WaitingAuthorisationToValidate:
		this.Status = PaymentWaitingPayment
		this.StatusLabel = PaymentWaitingPaymentLabel
		break
	case Refused:
		this.Status = PaymentRefused
		this.StatusLabel = PaymentRefusedLabel
		break
	case Captured:
		this.Status = PaymentPaid
		this.StatusLabel = PaymentPaidLabel
		break
	case Canceled:
		this.Status = PaymentCancelled
		this.StatusLabel = PaymentCancelledLabel
		break
	case Expired:
		this.Status = PaymentExpired
		this.StatusLabel = PaymentExpiredLabel
		break
	case UnderVerification:
		this.Status = PaymentOther
		this.StatusLabel = PaymentOtherLabel
		break
	case PartiallyAuthorised:
		this.Status = PaymentOther
		this.StatusLabel = PaymentOtherLabel
		break
	case Refunded:
		this.Status = PaymentRefound
		this.StatusLabel = PaymentRefoundLabel
		break
	case Created:
		this.Status = PaymentWaitingPayment
		this.StatusLabel = PaymentWaitingPaymentLabel
		break
	case Chargeback:
		this.Status = PaymentChargeback
		this.StatusLabel = PaymentChargebackLabel
		break
	case WaitingPayment:
		this.Status = PaymentWaitingPayment
		this.StatusLabel = PaymentWaitingPaymentLabel
		break
	case PendingRefund:
		this.Status = PaymentOther
		this.StatusLabel = PaymentOtherLabel
		break
	case Analyzing:
		this.Status = PaymentOther
		this.StatusLabel = PaymentOtherLabel
		break
	case PendingReview:
		this.Status = PaymentOther
		this.StatusLabel = PaymentOtherLabel
		break
	case Success:
		this.Status = PaymentSuccess
		this.StatusLabel = PaymentSuccessLabel
		break
	default:
		this.Status = PaymentError
		this.StatusLabel = PaymentErrorLabel
		break
	}
}

type PaymentResultValidationError struct {
	Error            bool              `jsonp:""`
	Message          string            `jsonp:""`
	ValidationErrors map[string]string `jsonp:""`
}

func (this *PaymentResultValidationError) AddError(key string, value string) *PaymentResultValidationError {
	if this.ValidationErrors == nil {
		this.ValidationErrors = make(map[string]string)
	}

	this.ValidationErrors[key] = value
	return this
}

func (this *PaymentResultValidationError) String() string {
	msg := this.Message
	if this.ValidationErrors != nil && len(this.ValidationErrors) > 0 {
		msg = fmt.Sprintf("%v\n%v", msg, this.ValidationErrors)
	}
	return msg
}

type Balance struct {
	WaitingFunds float64 `jsonp:""`
	Available    float64 `jsonp:""`
	Transferred  float64 `jsonp:""`
}

func NewBalance() *Balance {
	return &Balance{}
}

type Movement struct {
	Id            string  `jsonp:""`
	Amount        float64 `jsonp:""`
	Type          string  `jsonp:""`
	Date          string  `jsonp:""`
	Status        string  `jsonp:""`
	TransactionId string  `jsonp:""`
}

func NewMovement() *Movement {
	return &Movement{}
}

type TransferResult struct {
	Id                    string         `jsonp:""`
	DateCreated           string         `jsonp:""`
	Status                TransferStatus `jsonp:""`
	EffectiveDate         string         `jsonp:""`
	Type                  string         `jsonp:""`
	Value                 float64        `jsonp:""`
	NetValue              float64        `jsonp:""`
	TransferFee           float64        `jsonp:""`
	ScheduleDate          string         `jsonp:""`
	Authorized            bool           `jsonp:""`
	TransactionReceiptUrl string         `jsonp:""`
	TransactionId         string         `jsonp:""`
}

func NewTransferResult() *TransferResult {
	return &TransferResult{}
}

type PaymentResult struct {
	Error    bool   `jsonp:""`
	Message  string `jsonp:""`
	Request  string `jsonp:""`
	Response string `jsonp:""`

	CaptureRequest  string
	CaptureResponse string

	PagarmeV5Status   PagarmeV5Status   `jsonp:"pagarme_v5_status"`
	PicPayStatus      PicPayStatus      `jsonp:"picpay_status"`
	TransactionStatus TransactionStatus `jsonp:"payzen_status"`
	PayZenV4Status    string            `jsonp:"payzen_v4_status"`
	AsaasStatus       AsaasStatus       `jsonp:"asaas_status"`

	Status      PaymentStatus      `jsonp:""`
	StatusLabel PaymentStatusLabel `jsonp:""`

	NotificationId int64 `jsonp:""`
	OperationId int64 `jsonp:""`

	//v4

	TransactionStatusLabel string `jsonp:""`

	TransactionId string `jsonp:""`
	OrderId string `jsonp:""`

	ResponseCode       string `jsonp:""`
	ResponseCodeDetail string `jsonp:""`

	//v4
	ErroCode       string `jsonp:""`
	ErroCodeDetail string `jsonp:""`

	BoletoUrl    string `jsonp:""`
	BoletoNumber string `jsonp:""`

	PaymentType  PaymentType  `jsonp:""`
	PaymentEvent PaymentEvent `jsonp:""`

	//SubscriptionId string `jsonp:""`

	InstallmentId    string  `jsonp:""`
	InstallmentCount int64   `jsonp:""`
	Amount           float64 `jsonp:""`

	TokenInfo *TokenInfo `jsonp:""`

	Transactions []*TransactionItemResult `jsonp:""`

	SubscriptionInfo *SubscriptionResult `jsonp:""`

	ValidationErrors map[string]string `jsonp:""`

	Plan *Plan `jsonp:""`

	Platform            Gateway   `jsonp:""`
	Nsu                 string    `jsonp:""`
	BoletoOutputContent []byte    `json:"-"`
	BoletoFileName      string    `json:"-"`
	DueAt               time.Time `jsonp:""`

	QrCode           string `jsonp:"qrcode"`
	QrCodeUrl        string `jsonp:"qrcode_url"`
	QrPayload        string `jsonp:"qrcode_payload"`
	QrExpirationDate string `jsonp:"qrcode_expiration_date"`
	Barcode string `jsonp:"barcode"`
	PaymentUrl       string `jsonp:"payment_url"`

	InvoiceUrl  string `jsonp:""` // url da fatura
	PaymentLink string `jsonp:""` // identificador do link de pagamento

	AuthorizationId string `jsonp:"picpay_authorization_id"`
	CancellationId  string `jsonp:"picpay_cancellation_id"`
	ReferenceId     string `jsonp:"picpay_reference_id"`

	IsPagarme bool `jsonp:"is_pagarme"`
	IsPicPay  bool `jsonp:"is_picpay"`
	IsPayZen  bool `jsonp:"is_payzen"`
	IsAsaas   bool `jsonp:"is_asaas"`

	Balance   *Balance          `jsonp:""`
	Movements []*Movement       `jsonp:""`
	Transfers []*TransferResult `jsonp:""`

	Customer *Customer `jsonp:""`

	OverridePaymentStatusUrl string `json:"override_payment_status_url"`


	Tag interface{}
}

func (this *PaymentResult) isPagarme() bool {
	return this.Platform == GatewayPagarme
}

func (this *PaymentResult) isPayZen() bool {
	return this.Platform == GatewayPayZen
}

func (this *PaymentResult) isPicPay() bool {
	return this.Platform == GatewayPicPay
}

func (this *PaymentResult) isAsaas() bool {
	return this.Platform == GatewayAsaas
}

func (this *PaymentResult) IsInstallment() bool {
	return len(this.InstallmentId) > 0
}

func (this *PaymentResult) IsSubscription() bool {
	return len(this.SubscriptionInfo.SubscriptionId) > 0
}


func (this *PaymentResult) BuildStatus() {

	this.IsPayZen = this.isPayZen()
	this.IsPicPay = this.isPicPay()
	this.IsPagarme = this.isPagarme()
	this.IsAsaas = this.isAsaas()

	logs.Info("this.TransactionStatus!! %v", this.TransactionStatus)

	switch this.TransactionStatus {
	case Initial:
		this.Status = PaymentInitial
		this.StatusLabel = PaymentInitialLabel
		break
	case NotCreated:
		this.Status = PaymentInitial
		this.StatusLabel = PaymentInitialLabel
		break
	case Authorised:

		logs.Info("Authorized!! %v", this.IsPagarme)

		if this.IsPayZen || this.IsPagarme {
			this.Status = PaymentPaid
			this.StatusLabel = PaymentPaidLabel
		} else {
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
		}
		break
	case AuthorisedToValidate:
		this.Status = PaymentWaitingPayment
		this.StatusLabel = PaymentWaitingPaymentLabel
		break
	case WaitingAuthorisation:
		this.Status = PaymentWaitingPayment
		this.StatusLabel = PaymentWaitingPaymentLabel
		break
	case WaitingAuthorisationToValidate:
		this.Status = PaymentWaitingPayment
		this.StatusLabel = PaymentWaitingPaymentLabel
		break
	case Refused:
		this.Status = PaymentRefused
		this.StatusLabel = PaymentRefusedLabel
		break
	case Captured:
		this.Status = PaymentPaid
		this.StatusLabel = PaymentPaidLabel
		break
	case Canceled:
		this.Status = PaymentCancelled
		this.StatusLabel = PaymentCancelledLabel
		break
	case Expired:
		this.Status = PaymentExpired
		this.StatusLabel = PaymentExpiredLabel
		break
	case UnderVerification:
		this.Status = PaymentOther
		this.StatusLabel = PaymentOtherLabel
		break
	case PartiallyAuthorised:
		this.Status = PaymentOther
		this.StatusLabel = PaymentOtherLabel
		break
	case Refunded:
		this.Status = PaymentRefound
		this.StatusLabel = PaymentRefoundLabel
		break
	case Created:
		this.Status = PaymentWaitingPayment
		this.StatusLabel = PaymentWaitingPaymentLabel
		break
	case Chargeback:
		this.Status = PaymentChargeback
		this.StatusLabel = PaymentChargebackLabel
		break
	case WaitingPayment:
		this.Status = PaymentWaitingPayment
		this.StatusLabel = PaymentWaitingPaymentLabel
		break
	case PendingRefund:
		this.Status = PaymentOther
		this.StatusLabel = PaymentOtherLabel
		break
	case Analyzing:
		this.Status = PaymentOther
		this.StatusLabel = PaymentOtherLabel
		break
	case PendingReview:
		this.Status = PaymentOther
		this.StatusLabel = PaymentOtherLabel
		break
	case Success:
		this.Status = PaymentSuccess
		this.StatusLabel = PaymentSuccessLabel
		break
	default:
		this.Status = PaymentError
		this.StatusLabel = PaymentErrorLabel
		break
	}
}

func (this *PaymentResult) WithTransactionStatus(status TransactionStatus) *PaymentResult {
	this.TransactionStatus = status
	return this
}

func (this *PaymentResult) WithSuccess() *PaymentResult {
	this.TransactionStatus = Success
	return this
}

func (this *PaymentResult) IsCancelled() bool {
	return this.Status == PaymentCancelled
}

func (this *PaymentResult) IsRefused() bool {
	return this.Status == PaymentRefused
}

func (this *PaymentResult) IsRefound() bool {
	return this.Status == PaymentRefound
}

func (this *PaymentResult) IsExpired() bool {
	return this.Status == PaymentExpired
}

func (this *PaymentResult) IsChargebacked() bool {
	return this.Status == PaymentChargeback
}



func NewPaymentResult() *PaymentResult {
	result := new(PaymentResult)
	result.TokenInfo = new(TokenInfo)
	result.Transactions = []*TransactionItemResult{}
	result.SubscriptionInfo = new(SubscriptionResult)
	result.ValidationErrors = make(map[string]string)
	result.Customer = new(Customer)
	result.Movements = []*Movement{}
	result.Transfers = []*TransferResult{}
	result.Platform = GatewayPayZen
	return result
}

func NewPaymentResultValidationErrorWithMessage(errors map[string]string, message string) *PaymentResultValidationError {
	result := new(PaymentResultValidationError)
	result.Error = true
	result.Message = message
	result.ValidationErrors = errors
	return result
}

func NewPaymentResultValidationErrorWithErrorKey(message string, key string, value string) *PaymentResultValidationError {
	result := new(PaymentResultValidationError)
	result.Error = true
	return result.AddError(key, value)
}

func NewPaymentResultWithError(format string, args ...interface{}) *PaymentResultValidationError {
	result := new(PaymentResultValidationError)
	result.Error = true
	result.Message = fmt.Sprintf(format, args...)
	return result
}
