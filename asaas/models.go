package asaas

import (
	"fmt"
	"github.com/mobilemindtec/go-payments/api"
	"io"
	"time"
)

type BillingType string

const (
	BillingBoleto BillingType = "BOLETO"
	// ao gerar um pagamento CREDIT_CARD sem informar os dados do carão,
	// será retornado um link de pagamento onde será possível pagar
	// com catão de crédito ou débito
	BillingCreditCard BillingType = "CREDIT_CARD"
	BillingUndefined  BillingType = "UNDEFINED" // Perguntar ao Cliente

	BillingDebitCard BillingType = "DEBIT_CARD" // webhook
	BillingPix       BillingType = "PIX"        // webhook
	BillingTransfer  BillingType = "TRANSFER"   // webhook
	BillingDeposit   BillingType = "DEPOSIT"    // webhook

)

type DiscountType string

const (
	DiscountFixed      DiscountType = "FIXED"
	DiscountPercentage DiscountType = "PERCENTAGE"
)

//type SubscriptionCycle string
//
//const (
//	SubscriptionCycleNone SubscriptionCycle = ""
//	Weekly SubscriptionCycle = "WEEKLY" // semanal
//	Biweekly SubscriptionCycle = "BIWEEKLY" // quinzenal
//	Monthly SubscriptionCycle = "MONTHLY" // mensal
//	Quarterly SubscriptionCycle = "QUARTERLY" // trimestral
//	Semiannually SubscriptionCycle = "SEMIANNUALLY" // semestral
//	Yearly SubscriptionCycle = "YEARLY" // anual
//)

type ChargeType string

const (
	ChargeTypeNone ChargeType = ""
	Detached       ChargeType = "DETACHED"    // avulsa
	Recurrent      ChargeType = "RECURRENT"   // assinatura
	Installment    ChargeType = "INSTALLMENT" // parcelado
)

type PaymentType int64

const (
	PaymentDefault PaymentType = iota + 1
	PaymentSubscription
	PaymentLink
)

type CompanyType string

const (
	MEI         CompanyType = "MEI"
	LIMITED     CompanyType = "LIMITED"
	INDIVIDUAL  CompanyType = "INDIVIDUAL"
	ASSOCIATION CompanyType = "ASSOCIATION"
)

type PersonType string

const (
	FISICA   PersonType = "FISICA"
	JURIDICA PersonType = "JURIDICA"
)

type WebhookType string

const (
	WebhookPayment             WebhookType = "PAYMENT"
	WebhookInvoice             WebhookType = "INVOICE"
	WebhookTransfer            WebhookType = "TRANSFER"
	WebhookBill                WebhookType = "BILL"
	WebhookAnticipation        WebhookType = "RECEIVABLE_ANTICIPATION"
	WebhookMobilePhoneRecharge WebhookType = "MOBILE_PHONE_RECHARGE"
	WebhookAccountStatus       WebhookType = "ACCOUNT_STATUS"
)

type DocumentType string

const (
	IDENTIFICATION           DocumentType = "IDENTIFICATION"
	SOCIAL_CONTRACT          DocumentType = "SOCIAL_CONTRACT"
	ENTREPRENEUR_REQUIREMENT DocumentType = "ENTREPRENEUR_REQUIREMENT"
	MINUTES_OF_ELECTION      DocumentType = "MINUTES_OF_ELECTION"
	CUSTOM                   DocumentType = "CUSTOM"
)

type DocumentStatus string

const (
	NOT_SENT DocumentStatus = "NOT_SENT"
	PENDING  DocumentStatus = "PENDING"
	APPROVED DocumentStatus = "APPROVED"
	REJECTED DocumentStatus = "REJECTED"
)

type WebhookObject struct {
	Url         string      `json:"url" valid:"Required"`
	Email       string      `json:"email" valid:"Required"`
	Interrupted bool        `json:"interrupted" valid:"Required"`
	Enabled     bool        `json:"enabled" valid:"Required"`
	ApiVersion  int64       `json:"apiVersion" valid:"Required"`
	AuthToken   string      `json:"authToken" valid:"Required"`
	Type        WebhookType `json:"type" valid:"Required"`
}

func NewWebhookObject() *WebhookObject {
	return &WebhookObject{}
}

type Bank struct {
	Code string `json:"code" valid:"Required"`
}

func NewBank(code string) *Bank {
	return &Bank{Code: code}
}

type BankAccount struct {
	Bank        *Bank  `json:"bank" valid:"Required"`
	AccountName string `json:"accountName"` // Nome da conta bancária
	OwnerName   string `json:"ownerName" valid:"Required"`
	//Data de nascimento do proprietário da conta.
	//Somente quando a conta bancária não pertencer ao mesmo CPF ou CNPJ da conta Asaas.
	OwnerBirthDate  string              `json:"ownerBirthDate"`
	CpfCnpj         string              `json:"cpfCnpj" valid:"Required"`
	Agency          string              `json:"agency" valid:"Required"`
	Account         string              `json:"account" valid:"Required"`
	AccountDigit    string              `json:"accountDigit" valid:"Required"`
	BankAccountType api.BankAccountType `json:"bankAccountType" valid:"Required"`
}

func NewBankAccount(bank *Bank, bankAccountType api.BankAccountType) *BankAccount {
	return &BankAccount{Bank: bank, BankAccountType: bankAccountType}
}

type BankAccountSimple struct {
	Bank        string `json:"bank" valid:"Required"`
	AccountName string `json:"accountName"` // Nome da conta bancária
	Name        string `json:"name" valid:"Required"`
	//Data de nascimento do proprietário da conta.
	//Somente quando a conta bancária não pertencer ao mesmo CPF ou CNPJ da conta Asaas.
	OwnerBirthDate  string              `json:"ownerBirthDate"`
	CpfCnpj         string              `json:"cpfCnpj" valid:"Required"`
	Agency          string              `json:"agency" valid:"Required"`
	Account         string              `json:"account" valid:"Required"`
	AccountDigit    string              `json:"accountDigit" valid:"Required"`
	BankAccountType api.BankAccountType `json:"bankAccountType" valid:"Required"`
}

func NewBankAccountSimple(bank string, bankAccountType api.BankAccountType) *BankAccountSimple {
	return &BankAccountSimple{Bank: bank, BankAccountType: bankAccountType}
}

type Account struct {
	Name          string             `json:"name" valid:"Required"`
	Email         string             `json:"email" valid:"Required"`
	LoginEmail    string             `json:"loginEmail,omitempty" valid:""`
	CpfCnpj       string             `json:"cpfCnpj" valid:"Required"`
	CompanyType   CompanyType        `json:"companyType"`
	Phone         string             `json:"phone" valid:"Required"`
	MobilePhone   string             `json:"mobilePhone" valid:"Required"`
	Address       string             `json:"address" valid:"Required"`
	AddressNumber string             `json:"addressNumber" valid:"Required"`
	Complement    string             `json:"complement"`
	Province      string             `json:"province" valid:"Required"`   // bairro
	PostalCode    string             `json:"postalCode" valid:"Required"` // bairro
	BankAccount   *BankAccountSimple `json:"bankAccount,omitempty" valid:""`
	// result
	WalletId string `json:"walletId,omitempty"`
	ApiKey   string `json:"apiKey,omitempty"`
	//City string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
	//PersonType string `json:"personType,omitempty"`

	CompanyName string `json:"companyName,omitempty"`

	PersonType PersonType `json:"personType,omitempty"`

	//DenialReason string `json:"denialReason,omitempty"`

	Webhooks []*WebhookObject `json:"webhooks"`
}

func (this *Account) AddWebhook(objects ...*WebhookObject) *Account {
	for _, object := range objects {
		this.Webhooks = append(this.Webhooks, object)
	}
	return this
}

func NewAccount(bankAccount *BankAccountSimple) *Account {
	return &Account{BankAccount: bankAccount}
}

type AccountStatus struct {
	Id              string `json:"id"`
	CommercialInfo  string `json:"commercialInfo"`
	Documentation   string `json:"documentation"`
	BankAccountInfo string `json:"bankAccountInfo"`
	General         string `json:"general"`
}

type AccountStatusEvent struct {
	Event  string         `json:"event"`
	Status *AccountStatus `json:"accountStatus"`
}

type Document struct {
	Id           string       `json:"-"`
	DocumentFile io.Reader    `json:"documentFile"`
	Type         DocumentType `json:"type"`
}

type DocumentResponse struct {
	Id     string         `json:"id"`
	Status DocumentStatus `json:"status"`
}

type DocumentResponsible struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type DocumentEntry struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type DocumentDescription struct {
	Id            string               `json:"id"`
	Status        string               `json:"status"`
	Type          string               `json:"type"`
	Title         string               `json:"title"`
	Description   string               `json:"description"`
	Responsible   *DocumentResponsible `json:"responsible"`
	OnboardingUrl string               `json:"onboardingUrl"`
	Documents     []*DocumentEntry     `json:"documents"`
}

type Documents struct {
	RejectReasons string                 `json:"rejectReasons"`
	Data          []*DocumentDescription `json:"data"`
}

type Transfer struct {
	Value       float64      `json:"value" valid:"Required"`
	BankAccount *BankAccount `json:"bankAccount" valid:"Required"`
}

func NewTransfer(bankAccount *BankAccount, value float64) *Transfer {
	return &Transfer{Value: value, BankAccount: bankAccount}
}

type TransferEventData struct {
	Object                string       `json:"object"`
	Id                    string       `json:"id"`
	DateCreated           string       `json:"dateCreated"`
	Status                string       `json:"status"`
	EffectiveDate         string       `json:"effectiveDate"`
	EndToEndIdentifier    string       `json:"endToEndIdentifier"`
	Type                  string       `json:"type"`
	Value                 int64        `json:"value"`
	NetValue              int64        `json:"netValue"`
	TransferFee           int64        `json:"transferFee"`
	ScheduleDate          string       `json:"scheduleDate"`
	Authorized            bool         `json:"authorized"`
	FailReason            string       `json:"failReason"`
	TransactionReceiptUrl string       `json:"transactionReceiptUrl"`
	OperationType         string       `json:"operationType"`
	Description           string       `json:"description"`
	BankAccount           *BankAccount `json:"bankAccount"`
}

type TransferEvent struct {
	Event    string             `json:"event"`
	Transfer *TransferEventData `json:"transfer"`
}

type DefaultFilter struct {
	Limit      int64
	Offset     int64
	StartDate  string
	FinishDate string

	DateCreated     string              // transfer filter
	BankAccountType api.BankAccountType // transfer filter
}

func NewDefaultFilter() *DefaultFilter {
	return &DefaultFilter{}
}

func (this *DefaultFilter) SetStartDate(date time.Time) {
	this.StartDate = DateFormat(date)
}

func (this *DefaultFilter) SetFinishDate(date time.Time) {
	this.FinishDate = DateFormat(date)
}

func (this *DefaultFilter) SetDateCreated(date time.Time) {
	this.DateCreated = DateFormat(date)
}

func (this *DefaultFilter) ToMap() map[string]string {
	filterMap := make(map[string]string)
	if len(this.StartDate) > 0 {
		filterMap["startDate"] = this.StartDate
	}
	if len(this.FinishDate) > 0 {
		filterMap["finishDate"] = this.FinishDate
	}
	if this.Limit > 0 {
		filterMap["limit"] = fmt.Sprintf("%v", this.Limit)
	}
	if this.Offset > 0 {
		filterMap["offset"] = fmt.Sprintf("%v", this.Offset)
	}
	if len(this.DateCreated) > 0 {
		filterMap["dateCreated"] = this.DateCreated
	}
	if this.BankAccountType != api.BankAccountTypeNone {
		filterMap["type"] = string(this.BankAccountType)
	}
	return filterMap
}

type Customer struct {
	Id                   string `json:"id,omitempty"`
	Name                 string `json:"name", valid:"Required"`
	CpfCnpj              string `json:"cpfCnpj" valid:"Required"`
	Email                string `json:"email"`
	MobilePhone          string `json:"mobilePhone"`
	Phone                string `json:"phone"`
	NotificationDisabled bool   `json:"notificationDisabled"` // true para desabilitar o envio de notificações de cobrança
	ExternalReference    string `json:"externalReference" valid:"Required"`
	AdditionalEmails     string `json:"additionalEmails"` // Emails adicionais para envio de notificações de cobrança separados por ","

	Address       string `json:"address"`
	AddressNumber string `json:"addressNumber"`
	Province      string `json:"province"`
	PostalCode    string `json:"postalCode"`
	City          int `json:"city"`
	State         string `json:"state"`
}

type Customers = []*Customer

func NewCustomer(name string, cpfCnpj string) *Customer {
	return &Customer{Name: name, CpfCnpj: cpfCnpj}
}

// /
// / Informações de desconto
// /
type Discount struct {
	Value            float64      `json:"value"`            // Valor percentual ou fixo de desconto a ser aplicado sobre o valor da cobrança
	DueDateLimitDays int64        `json:"dueDateLimitDays"` // Dias antes do vencimento para aplicar desconto. Ex: 0 = até o vencimento, 1 = até um dia antes, 2 = até dois dias antes, e assim por diante
	Type             DiscountType `json:"type"`
}

func NewDiscount(value float64, dueDateLimitDays int64, discountType DiscountType) *Discount {
	return &Discount{Value: value, DueDateLimitDays: dueDateLimitDays, Type: discountType}
}

// /
// / Informações de juros para pagamento após o vencimento
// /
type Interest struct {
	Value float64 `json:"value"` // Percentual de juros ao mês sobre o valor da cobrança para pagamento após o vencimento
}

func NewInterest(value float64) *Interest {
	return &Interest{Value: value}
}

// /
// / Informações de multa para pagamento após o vencimento
// /
type Fine struct {
	Value         float64 `json:"value"`         // Percentual de multa sobre o valor da cobrança para pagamento após o vencimento
	PostalService bool    `json:"postalService"` // Define se a cobrança será enviada via Correios
}

func NewFine(value float64) *Fine {
	return &Fine{Value: value}
}

type TokenRequest struct {
	CreditCardCcv         string `json:"creditCardCcv" valid:"Required"`
	CreditCardHolderName  string `json:"creditCardHolderName" valid:"Required"`
	CreditCardExpiryMonth string `json:"creditCardExpiryMonth" valid:"Required"`
	CreditCardNumber      string `json:"creditCardNumber" valid:"Required"`
	CreditCardExpiryYear  string `json:"creditCardExpiryYear" valid:"Required"`
	Customer              string `json:"customer" valid:"Required"`
}

func NewTokenRequest(customer string) *TokenRequest {
	return &TokenRequest{Customer: customer}
}

type Card struct {
	HolderName   string `json:"holderName" valid:"Required"`
	Number       string `json:"number" valid:"Required"`
	ExpiryMonth  string `json:"expiryMonth" valid:"Required"` // 06
	ExpiryYear   string `json:"expiryYear" valid:"Required"`  // 2019
	SecurityCode string `json:"ccv" valid:"Required"`
}

type CardHolderInfo struct {
	Name              string `json:"name" valid:"Required"`
	Email             string `json:"email" valid:"Required"`
	CpfCnpj           string `json:"cpfCnpj" valid:"Required"`
	PostalCode        string `json:"postalCode" valid:"Required"`
	AddressNumber     string `json:"addressNumber" valid:"Required"`
	AddressComplement string `json:"addressComplement"`
	Phone             string `json:"phone" valid:"Required"`
	MobilePhone       string `json:"mobilePhone"`
}

type Split struct {
	WalletId        string  `json:"walletId" valid:"Required"` // Identificador da carteira (retornado no momento da criação da conta)
	FixedValue      float64 `json:"fixedValue"`                // Valor fixo a ser transferido para a conta quando a cobrança for recebida
	PercentualValue float64 `json:"percentualValue"`
}

/*
Token

Ao realizar uma primeira transação para o cliente com cartão de crédito,
a resposta do Asaas lhe devolverá dentro do objeto creditCard, o atributo creditCardToken.

# Parcelamento

Para criar uma cobrança parcelada, ao invés de enviar o parâmetro value,
envie installmentCount e installmentValue,
que representam o número de parcelas e o valor da cada parcela respectivamente.
*/
type Payment struct {
	BillingType       BillingType     `json:"billingType" valid:"Required"`
	Value             float64         `json:"value" valid:"Required"`     // Valor da cobrança
	DueDate           string          `json:"dueDate,omitempty" valid:""` // date do vencimento da cobrança
	Description       string          `json:"description,omitempty"`
	ExternalReference string          `json:"externalReference,omitempty"` // Campo livre para busca
	InstallmentCount  int64           `json:"installmentCount,omitempty"`  // Número de parcelas (somente no caso de cobrança parcelada)
	InstallmentValue  int64           `json:"installmentValue,omitempty"`  // Valor de cada parcela (somente no caso de cobrança parcelada)
	TotalValue        float64         `json:"totalValue,omitempty"`        // valor total para parcelamento
	Discount          *Discount       `json:"discount,omitempty"`
	Interest          *Interest       `json:"interest,omitempty"`
	Fine              *Fine           `json:"file,omitempty"`
	PostalService     bool            `json:"postalService,omitempty"`
	Customer          string          `json:"customer"`
	Card              *Card           `json:"creditCard,omitempty"`           // obrigatório compra cartão
	CardHolderInfo    *CardHolderInfo `json:"creditCardHolderInfo,omitempty"` // obrigatório compra cartão
	CardToken         string          `json:"creditCardToken,omitempty"`      // obrigatório compra cartão
	RemoteIp          string          `json:"remoteIp,omitempty"`             // obrigatório compra cartão
	Splits            []*Split        `json:"split,omitempty"`

	PaymentType PaymentType `json:"-"`

	// subscription
	SubscriptionCycle     api.SubscriptionCycle `json:"cycle,omitempty"`
	NextDueDate           string                `json:"nextDueDate,omitempty"` // Vencimento da primeira mensalidade
	EndDate               string                `json:"endDate,omitempty"`     // Data limite para vencimento das mensalidades - assuntura e link de pagamento
	MaxPayments           int64                 `json:"maxPayments,omitempty"` //Número máximo de mensalidades a serem geradas para esta assinatura
	UpdatePendingPayments bool                  `json:"updatePendingPayments"` // true para atualizar mensalidades já existentes com o novo valor ou forma de pagamento

	// link de pagamento

	Name string `json:"name,omitempty"`
	// Caso seja possível o pagamento via boleto bancário, define a quantidade de dias úteis que
	// o seu cliente poderá pagar o boleto após gerado
	DueDateLimitDays int64      `json:"dueDateLimitDays,omitempty"`
	ChargeType       ChargeType `json:"chargeType,omitempty"`
	//Quantidade máxima de parcelas que seu cliente poderá parcelar o valor do link de pagamentos caso a
	//forma de cobrança selecionado seja Parcelamento. Caso não informado o valor padrão será de 1 parcela
	MaxInstallmentCount int64 `json:"maxInstallmentCount,omitempty"`
	// EndDate string `json:"endDate,omitempty"` Data de encerramento, a partir desta data o seu link de pagamentos será desativado automaticamente

	Id string `json:"-"`
}

func (this *Payment) SetDueDate(date time.Time) {
	this.DueDate = DateFormat(date)
}

func (this *Payment) SetEndDate(date time.Time) {
	this.EndDate = DateFormat(date)
}

func NewPayment() *Payment {
	return &Payment{}
}

func NewPaymentWithCard(customerId string, orderId string, value float64) *Payment {
	return &Payment{
		Card:              new(Card),
		CardHolderInfo:    new(CardHolderInfo),
		BillingType:       BillingCreditCard,
		Value:             value,
		Customer:          customerId,
		ExternalReference: orderId,
		DueDate:           DateFormat(time.Now()),
		PaymentType:       PaymentDefault,
	}
}

func NewPaymenInstallmenttWithCard(customerId string, orderId string, value float64, installmentCount int64) *Payment {
	return &Payment{
		Card:              new(Card),
		CardHolderInfo:    new(CardHolderInfo),
		BillingType:       BillingCreditCard,
		TotalValue:        value,
		InstallmentCount:  installmentCount,
		Customer:          customerId,
		ExternalReference: orderId,
		DueDate:           DateFormat(time.Now()),
		PaymentType:       PaymentDefault,
	}
}

func NewPaymenInstallmenttWithBoleto(customerId string, orderId string, value float64, installmentCount int64) *Payment {
	return &Payment{
		BillingType:       BillingBoleto,
		TotalValue:        value,
		InstallmentCount:  installmentCount,
		Customer:          customerId,
		ExternalReference: orderId,
		DueDate:           DateFormat(time.Now()),
		PaymentType:       PaymentDefault,
	}
}

func NewPaymentWithPix(customerId string, orderId string, dueDate time.Time, value float64) *Payment {
	return &Payment{
		BillingType:       BillingPix,
		Value:             value,
		Customer:          customerId,
		ExternalReference: orderId,
		DueDate:           DateFormat(dueDate),
		PaymentType:       PaymentDefault,
	}
}

func NewPaymentWithBoleto(customerId string, orderId string, dueDate time.Time, value float64) *Payment {
	return &Payment{
		BillingType:       BillingBoleto,
		Value:             value,
		Customer:          customerId,
		ExternalReference: orderId,
		DueDate:           DateFormat(dueDate),
		PaymentType:       PaymentDefault,
	}
}

func NewSubscription() *Payment {
	return &Payment{}
}

func NewSubscriptionWithBoleto(customerId string, orderId string, cycle api.SubscriptionCycle, nextDueDate time.Time, value float64) *Payment {
	return &Payment{
		BillingType:       BillingBoleto,
		Value:             value,
		Customer:          customerId,
		ExternalReference: orderId,
		NextDueDate:       DateFormat(nextDueDate),
		PaymentType:       PaymentSubscription,
		SubscriptionCycle: cycle,
	}
}

func NewSubscriptionWithCard(customerId string, orderId string, cycle api.SubscriptionCycle, nextDueDate time.Time, value float64) *Payment {
	return &Payment{
		Card:              new(Card),
		CardHolderInfo:    new(CardHolderInfo),
		BillingType:       BillingCreditCard,
		Value:             value,
		Customer:          customerId,
		ExternalReference: orderId,
		NextDueDate:       DateFormat(nextDueDate),
		PaymentType:       PaymentSubscription,
		SubscriptionCycle: cycle,
	}
}

func NewSubscriptionWithCardToken(paymentId string) *Payment {
	return &Payment{
		Id:             paymentId,
		Card:           new(Card),
		CardHolderInfo: new(CardHolderInfo),
	}
}

func NewPaymentLink(value float64, chargeType ChargeType, dueDateLimitDays int64) *Payment {
	return &Payment{
		Value:            value,
		ChargeType:       chargeType,
		PaymentType:      PaymentLink,
		DueDateLimitDays: dueDateLimitDays,
		BillingType:      BillingUndefined,
	}
}

func NewPaymentLinkWithBoleto(value float64, chargeType ChargeType, dueDateLimitDays int64) *Payment {
	return &Payment{
		Value:            value,
		ChargeType:       chargeType,
		PaymentType:      PaymentLink,
		DueDateLimitDays: dueDateLimitDays,
		BillingType:      BillingBoleto,
	}
}

func NewPaymentLinkWithCard(value float64, chargeType ChargeType) *Payment {
	return &Payment{
		Value:       value,
		ChargeType:  chargeType,
		PaymentType: PaymentLink,
		BillingType: BillingBoleto,
	}
}

type PaymentInCash struct {
	Id             string  `valid:"Required"`
	PaymentDate    string  `json:"paymentDate" valid:"Required"`
	NotifyCustomer bool    `json:"notifyCustomer"`
	Value          float64 `json:"value" valid:"Required"` // Valor da cobrança
}

func NewPaymentInCash(id string, date time.Time, value float64) *PaymentInCash {
	return &PaymentInCash{Id: id, Value: value, PaymentDate: DateFormat(date)}
}

func (this *PaymentInCash) SetPaymentDate(date time.Time) {
	this.PaymentDate = DateFormat(date)
}

type CardResponse struct {
	Number string `json:"creditCardNumber"` // Últimos 4 dígitos do cartão utilizado
	Brand  string `json:"creditCardBrand"`  // Bandeira do cartão utilizado
	Token  string `json:"creditCardToken"`  // Token do cartão de crédito que poderá ser enviado nas próximas transações sem a necessidade de informar novamente os dados de cartão e do titular.
}

type ResponseError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type CustomerResults struct {
	Object     string      `json:"object"`
	HasMore    bool        `json:"hasMore"`
	TotalCount int64       `json:"totalCount"`
	Limit      int64       `json:"limit"`
	Offset     int64       `json:"offset"`
	Data       []*Customer `json:"data"`
}

func (this *CustomerResults) HasData() bool {
	return len(this.Data) > 0
}

func (this *CustomerResults) Count() int {
	return len(this.Data)
}

func (this *CustomerResults) First() *Customer {
	if this.HasData() {
		return this.Data[0]
	}
	return nil
}

func (this *CustomerResults) Last() *Customer {
	if this.HasData() {
		return this.Data[len(this.Data)-1]
	}
	return nil
}

type Wallet struct {
	Id string `json:"id"`
}

type WalletResults struct {
	Object     string    `json:"object"`
	HasMore    bool      `json:"hasMore"`
	TotalCount int64     `json:"totalCount"`
	Limit      int64     `json:"limit"`
	Offset     int64     `json:"offset"`
	Data       []*Wallet `json:"data"`
}

func (this *WalletResults) HasData() bool {
	return len(this.Data) > 0
}

func (this *WalletResults) Count() int {
	return len(this.Data)
}

func (this *WalletResults) First() *Wallet {
	if this.HasData() {
		return this.Data[0]
	}
	return nil
}

func (this *WalletResults) Last() *Wallet {
	if this.HasData() {
		return this.Data[len(this.Data)-1]
	}
	return nil
}

type TransferResult struct {
	Id                    string             `json:"id"`                    // "777eb7c8-b1a2-4356-8fd8-a1b0644b5282",
	Object                string             `json:"object"`                // "transfer",
	DateCreated           string             `json:"dateCreated"`           // "2019-05-02",
	Status                api.TransferStatus `json:"status"`                // "PENDING",
	EffectiveDate         string             `json:"effectiveDate"`         // null,
	Type                  string             `json:"type"`                  // "BANK_ACCOUNT",
	Value                 float64            `json:"value"`                 // 1000,
	NetValue              float64            `json:"netValue"`              // 1000,
	TransferFee           float64            `json:"transferFee"`           // 0, // Taxa de transferência
	ScheduleDate          string             `json:"scheduleDate"`          // "2019-05-02",
	Authorized            bool               `json:"authorized"`            // true,
	BankAccount           *BankAccount       `json:"bankAccount"`           //
	TransactionReceiptUrl string             `json:"transactionReceiptUrl"` // null
}

type TransferResults struct {
	Object     string            `json:"object"`
	HasMore    bool              `json:"hasMore"`
	TotalCount int64             `json:"totalCount"`
	Limit      int64             `json:"limit"`
	Offset     int64             `json:"offset"`
	Data       []*TransferResult `json:"data"`
}

func (this *TransferResults) HasData() bool {
	return len(this.Data) > 0
}

func (this *TransferResults) Count() int {
	return len(this.Data)
}

func (this *TransferResults) First() *TransferResult {
	if this.HasData() {
		return this.Data[0]
	}
	return nil
}

func (this *TransferResults) Last() *TransferResult {
	if this.HasData() {
		return this.Data[len(this.Data)-1]
	}
	return nil
}

type FinancialTransaction struct {
	Object      string  `json:"object"`
	Id          string  `json:"id"`
	Value       float64 `json:"value"`
	Balance     float64 `json:"balance"`
	Type        string  `json:"type"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
}

type FinancialTransactionResults struct {
	Object     string                  `json:"object"`
	HasMore    bool                    `json:"hasMore"`
	TotalCount int64                   `json:"totalCount"`
	Limit      int64                   `json:"limit"`
	Offset     int64                   `json:"offset"`
	Data       []*FinancialTransaction `json:"data"`
}

func (this *FinancialTransactionResults) HasData() bool {
	return len(this.Data) > 0
}

func (this *FinancialTransactionResults) Count() int {
	return len(this.Data)
}

func (this *FinancialTransactionResults) First() *FinancialTransaction {
	if this.HasData() {
		return this.Data[0]
	}
	return nil
}

func (this *FinancialTransactionResults) Last() *FinancialTransaction {
	if this.HasData() {
		return this.Data[len(this.Data)-1]
	}
	return nil
}

type AccountResults struct {
	Object     string     `json:"object"`
	HasMore    bool       `json:"hasMore"`
	TotalCount int64      `json:"totalCount"`
	Limit      int64      `json:"limit"`
	Offset     int64      `json:"offset"`
	Data       []*Account `json:"data"`
}

func (this *AccountResults) HasData() bool {
	return len(this.Data) > 0
}

func (this *AccountResults) Count() int {
	return len(this.Data)
}

func (this *AccountResults) First() *Account {
	if this.HasData() {
		return this.Data[0]
	}
	return nil
}

func (this *AccountResults) Last() *Account {
	if this.HasData() {
		return this.Data[len(this.Data)-1]
	}
	return nil
}

type PaymentResults struct {
	Object     string      `json:"object"`
	HasMore    bool        `json:"hasMore"`
	TotalCount int64       `json:"totalCount"`
	Limit      int64       `json:"limit"`
	Offset     int64       `json:"offset"`
	Data       []*Response `json:"data"`
}

func (this *PaymentResults) HasData() bool {
	return len(this.Data) > 0
}

func (this *PaymentResults) Count() int {
	return len(this.Data)
}

func (this *PaymentResults) First() *Response {
	if this.HasData() {
		return this.Data[0]
	}
	return nil
}

func (this *PaymentResults) Last() *Response {
	if this.HasData() {
		return this.Data[len(this.Data)-1]
	}
	return nil
}

type Response struct {
	ReturnedId     interface{} `json:"id"`
	Id             string      // "pay_4440248962351893",
	PaymentLink    string      `json:"paymentLink"` // Identificador único do link de pagamentos ao qual a cobrança pertence
	StatusText     string      `json:"status"`      // "PENDING",
	Status         api.AsaasStatus
	InvoiceUrl     string `json:"invoiceUrl"`   // URL da fatura,
	BankSlipUrl    string `json:"bankSlipUrl"`  // URL para download do boleto
	SubscriptionId string `json:"subscription"` // Identificador único da assinatura (quando cobrança recorrente)

	Object                 string      `json:"object"`                  // "payment",
	DateCreated            string      `json:"dateCreated"`             // "2021-07-21",
	Customer               string      `json:"customer"`                // "cus_000004699156",
	Value                  float64     `json:"value"`                   // 10,
	NetValue               float64     `json:"netValue"`                // Valor líquido da cobrança após desconto da tarifa do Asaas
	OriginalValue          float64     `json:"originalValue,omitempty"` // Valor original da cobrança (preenchido quando paga com juros e multa)
	InterestValue          float64     `json:"interestValue,omitempty"` // Valor calculado de juros e multa que deve ser pago após o vencimento da cobrança
	Description            string      `json:"description"`             // "",
	BillingType            BillingType `json:"billingType"`             // "BOLETO",
	DueDate                string      `json:"dueDate"`                 // Data de vencimento da cobrança
	OriginalDueDate        string      `json:"originalDueDate"`         // "2023-07-07",
	PaymentDate            string      `json:"paymentDate"`             // Data de liquidação da cobrança no Asaas
	ClientPaymentDate      string      `json:"clientPaymentDate"`       // Data em que o cliente efetuou o pagamento do boleto
	InvoiceNumber          string      `json:"invoiceNumber"`           // Número da fatura
	ExternalReference      string      `json:"externalReference"`       // "4d8ccb10-6c6c-4cd3-8514-450434e4c323",
	Deleted                bool        `json:"deleted"`                 // false,
	Anticipated            bool        `json:"anticipated"`             // Define se a cobrança foi antecipada ou está em processo de antecipação
	CreditDate             string      `json:"creditDate"`              // null,
	EstimatedCreditDate    string      `json:"estimatedCreditDate"`     // null,
	LastInvoiceViewedDate  string      `json:"lastInvoiceViewedDate"`   // null,
	LastBankSlipViewedDate string      `json:"lastBankSlipViewedDate"`  // null,
	ConfirmedDate          string      `json:"confirmedDate"`           // Data de confirmação da cobrança (Somente para cartão de crédito)
	Discount               *Discount   `json:"discount"`
	Fine                   *Fine       `json:"fine"`
	Interest               *Interest   `json:"interest"`
	PostalService          bool        `json:"postalService"` // false
	Installment            string      `json:"installment"`   // Identificador único do parcelamento (quando cobrança parcelada)

	InstallmentCount int64 `json:"installmentCount"`

	Card *CardResponse `json:"creditCard"`

	CustomerResults             *CustomerResults
	PaymentResults              *PaymentResults
	FinancialTransactionResults *FinancialTransactionResults
	BankAccount                 *BankAccount `json:"bankAccount,omitempty"`
	TransferResults             *TransferResults
	AccountResults              *AccountResults
	AccountStatus               *AccountStatus
	Webhook                     *WebhookObject
	WalletResults               *WalletResults
	Documents                   *Documents
	DocumentResponse            *DocumentResponse

	EncodedImage   string `json:"encodedImage"`
	Payload        string `json:"payload"`
	ExpirationDate string `json:"expirationDate"`

	// assinatura
	MaxPayments       int64                 `json:"maxPayments,omitempty"`
	SubscriptionCycle api.SubscriptionCycle `json:"cycle,omitempty"`
	NextDueDate       string                `json:"nextDueDate"`
	EndDate           string                `json:"endDate"`

	// link de pagamento
	Name                string     `json:"name"`
	ChargeType          ChargeType `json:"chargeType"`
	Url                 string     `json:"url"`
	Active              bool       `json:"active,omitempty"`
	MaxInstallmentCount int64      `json:"maxInstallmentCount,omitempty"`
	ViewCount           int64      `json:"viewCount,omitempty"`
	DueDateLimitDays    int64      `json:"dueDateLimitDays,omitempty"`

	TotalBalance float64 `json:"totalBalance,omitempty"`

	Request  string
	Response string
	Errors   []*ResponseError `json:"errors"`
	Message  string
	Error    bool
}

func NewResponse() *Response {
	return &Response{
		CustomerResults:             &CustomerResults{Data: []*Customer{}},
		PaymentResults:              &PaymentResults{Data: []*Response{}},
		FinancialTransactionResults: &FinancialTransactionResults{Data: []*FinancialTransaction{}},
		TransferResults:             &TransferResults{Data: []*TransferResult{}},
		AccountResults:              &AccountResults{Data: []*Account{}},
		WalletResults:               &WalletResults{Data: []*Wallet{}},
	}
}

func (this *Response) HasError() bool {
	return this.Errors != nil && len(this.Errors) > 0
}

func (this *Response) BuildStatus() {

	this.Id = fmt.Sprintf("%v", this.ReturnedId)

	if this.Deleted {
		this.Status = api.AsaasDeleted
		return
	}

	switch this.StatusText {
	case "PENDING":
		this.Status = api.AsaasPending
		break
	case "RECEIVED":
		this.Status = api.AsaasReceived
		break
	case "CONFIRMED":
		this.Status = api.AsaasConfirmed
		break
	case "OVERDUE":
		this.Status = api.AsaasOverdue
		break
	case "REFUNDED":
		this.Status = api.AsaasRefunded
		break
	case "RECEIVED_IN_CASH":
		this.Status = api.AsaasReceivedInCash
		break
	case "REFUND_REQUESTED":
		this.Status = api.AsaasRefundRequested
		break
	case "CHARGEBACK_REQUESTED":
		this.Status = api.AsaasChargebackRequested
		break
	case "CHARGEBACK_DISPUTE":
		this.Status = api.AsaasChargebackDispute
		break
	case "AWAITING_CHARGEBACK_REVERSAL":
		this.Status = api.AsaasAwaitingChargebackReversal
		break
	case "DUNNING_REQUESTED":
		this.Status = api.AsaasDunningRequested
		break
	case "DUNNING_RECEIVED":
		this.Status = api.AsaasDunningReceived
		break
	case "AWAITING_RISKANALYSIS":
		this.Status = api.AsaasAwaitingRiskAnalysis
		break
	case "ACTIVE":
		this.Status = api.AsaasActive
		break
	case "EXPIRED":
		this.Status = api.AsaasExpired
		break
	}
}

func (this *Response) GetPayZenSOAPStatus() api.TransactionStatus {

	this.BuildStatus()

	switch this.Status {
	case api.AsaasPending:
		return api.WaitingPayment
	case api.AsaasReceived:
		return api.Other // status de quando o valor fica disponível na conta para saque
	case api.AsaasConfirmed:
		return api.Captured
	case api.AsaasOverdue:
		return api.Expired
	case api.AsaasRefunded:
		return api.Refunded
	case api.AsaasReceivedInCash:
		return api.Captured
	case api.AsaasRefundRequested:
		return api.PendingRefund
	case api.AsaasChargebackRequested:
		return api.Chargeback
	case api.AsaasChargebackDispute:
		return api.Chargeback
	case api.AsaasAwaitingChargebackReversal:
		return api.Chargeback
	case api.AsaasDunningRequested:
		return api.Other
	case api.AsaasDunningReceived:
		return api.Other
	case api.AsaasAwaitingRiskAnalysis:
		return api.Other
	case api.AsaasActive:
		return api.Authorised
	case api.AsaasExpired:
		return api.Expired
	case api.AsaasDeleted:
		return api.Canceled
	case api.AsaasSuccess:
		return api.Success
	default:
		return api.Error
	}
}

func (this *Response) GetPaymentType() api.PaymentType {
	switch this.BillingType {
	case BillingBoleto:
		return api.PaymentTypeBoleto
	case BillingCreditCard:
		return api.PaymentTypeCreditCard
	case BillingDebitCard:
		return api.PaymentTypeDebitCard
	case BillingPix:
		return api.PaymentTypePix
	case BillingTransfer:
		return api.PaymentTypeTransfer
	case BillingDeposit:
		return api.PaymentTypeDeposit
		//case BillingUndefined:
	default:
		return api.PaymentTypeUndefined
	}
}

func (this *Response) ErrorsToMap() map[string]string {
	errors := make(map[string]string)

	if this.Errors != nil {
		for _, it := range this.Errors {
			errors[it.Code] = it.Description
		}
	}

	return errors
}

func (this *Response) ErrorsCount() int {

	if this.Errors != nil {
		return len(this.Errors)
	}

	return 0
}

func (this *Response) FirstError() string {
	if len(this.Errors) > 0 {
		return this.Errors[0].Description
	}
	return ""
}

func DateFormat(date time.Time) string {
	return date.Format("2006-01-02")
}
