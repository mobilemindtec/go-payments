package v5

import (
	"github.com/mobilemindtec/go-utils/json"
	"github.com/mobilemindtec/go-utils/v2/maps"
	"time"
)

type CustomerQuery struct {
	Name     string `jsonp:"name,omitempty"`
	Document string `jsonp:"document,omitempty"`
	Email    string `jsonp:"email,omitempty"`
	Code     string `jsonp:"code,omitempty"`
	Gender   Gender `jsonp:"gender,omitempty"`
	Page     int    `jsonp:"page,omitempty"`
	Size     int    `jsonp:"size,omitempty"`
}

func NewCustomerQuery() *CustomerQuery {
	return &CustomerQuery{}
}

func (this *CustomerQuery) WithDocument(doc string) *CustomerQuery {
	this.Document = doc
	return this
}

func (this *CustomerQuery) UrlQuery() string {
	m, _ := json.EncodeAsMap(this)
	return maps.ToUrlQuery(m)
}

type PlanQuery struct {
	Name         string    `jsonp:"name,omitempty"`
	Status       Status    `jsonp:"status,omitempty"`
	CreatedSince time.Time `jsonp:"created_since,date,omitempty"`
	CreatedUntil time.Time `jsonp:"created_until,date,omitempty"`
	Page         int       `jsonp:"page,omitempty"`
	Size         int       `jsonp:"size,omitempty"`
}

func (this *PlanQuery) UrlQuery() string {
	m, _ := json.EncodeAsMap(this)
	return maps.ToUrlQuery(m)
}

func NewPlanQuery() *PlanQuery {
	return &PlanQuery{}
}

type OrderQuery struct {
	Code         string    `jsonp:"code,omitempty"`
	CustomerId   string    `jsonp:"customer_id,omitempty"`
	Status       Status    `jsonp:"status,omitempty"`
	CreatedSince time.Time `jsonp:"created_since,date,omitempty"`
	CreatedUntil time.Time `jsonp:"created_until,date,omitempty"`
	Page         int       `jsonp:"page,omitempty"`
	Size         int       `jsonp:"size,omitempty"`
}

func (this *OrderQuery) UrlQuery() string {
	m, _ := json.EncodeAsMap(this)
	return maps.ToUrlQuery(m)
}

func NewOrderQuery() *OrderQuery {
	return &OrderQuery{}
}

type SubscriptionQuery struct {
	Status           Status      `jsonp:"status,omitempty"`
	Code             string      `jsonp:"code,omitempty"`
	BillingType      BillingType `jsonp:"billing_type,omitempty"`
	CustomerId       string      `jsonp:"customer_id,omitempty"`
	PlanId           string      `jsonp:"plan_id,omitempty"`
	CardId           string      `jsonp:"card_id,omitempty"`
	NextBillingSince time.Time   `jsonp:"next_billing_since,date,omitempty"`
	NextBillingUntil time.Time   `jsonp:"next_billing_until,date,omitempty"`
	CreatedSince     time.Time   `jsonp:"created_since,date,omitempty"`
	CreatedUntil     time.Time   `jsonp:"created_until,date,omitempty"`
	Page             int         `jsonp:"page,omitempty"`
	Size             int         `jsonp:"size,omitempty"`
}

func (this *SubscriptionQuery) UrlQuery() string {
	m, _ := json.EncodeAsMap(this)
	return maps.ToUrlQuery(m)
}

func NewSubscriptionQuery() *SubscriptionQuery {
	return &SubscriptionQuery{}
}

type InvoiceQuery struct {
	Status         InvoiceStatus `jsonp:"status,omitempty"`
	CustomerId     string        `jsonp:"customer_id,omitempty"`
	SubscriptionId string        `jsonp:"subscription_id,omitempty"`
	DueSince       time.Time     `jsonp:"due_since,date,omitempty"`
	DueUntil       time.Time     `jsonp:"due_until,date,omitempty"`
	CreatedSince   time.Time     `jsonp:"created_since,date,omitempty"`
	CreatedUntil   time.Time     `jsonp:"created_until,date,omitempty"`
	Page           int           `jsonp:"page,omitempty"`
	Size           int           `jsonp:"size,omitempty"`
}

func (this *InvoiceQuery) UrlQuery() string {
	m, _ := json.EncodeAsMap(this)
	return maps.ToUrlQuery(m)
}

func NewInvoiceQuery() *InvoiceQuery {
	return &InvoiceQuery{}
}

type ChargeQuery struct {
	Code          string        `jsonp:"code,omitempty"`
	Status        ChargeStatus  `jsonp:"status,omitempty"`
	PaymentMethod PaymentMethod `jsonp:"payment_method,omitempty"`
	CustomerId    string        `jsonp:"customer_id,omitempty"`
	OrderId       string        `jsonp:"order_id,omitempty"`
	CreatedSince  time.Time     `jsonp:"created_since,date,omitempty"`
	CreatedUntil  time.Time     `jsonp:"created_until,date,omitempty"`
	Page          int           `jsonp:"page,omitempty"`
	Size          int           `jsonp:"size,omitempty"`
}

func (this *ChargeQuery) UrlQuery() string {
	m, _ := json.EncodeAsMap(this)
	return maps.ToUrlQuery(m)
}

func NewChargeQuery() *ChargeQuery {
	return &ChargeQuery{}
}

type TransferQuery struct {
	Status       ChargeStatus `jsonp:"status,omitempty"`
	CreatedSince time.Time    `jsonp:"created_since,date,omitempty"`
	CreatedUntil time.Time    `jsonp:"created_until,date,omitempty"`
	Page         int          `jsonp:"page,omitempty"`
	Size         int          `jsonp:"size,omitempty"`
}

func (this *TransferQuery) UrlQuery() string {
	m, _ := json.EncodeAsMap(this)
	return maps.ToUrlQuery(m)
}

func NewTransferQuery() *TransferQuery {
	return &TransferQuery{}
}

type BalanceQuery struct {
	Status       ChargeStatus `jsonp:"status,omitempty"`
	CreatedSince time.Time    `jsonp:"created_since,date,omitempty"`
	CreatedUntil time.Time    `jsonp:"created_until,date,omitempty"`
	RecipientId  string       `jsonp:"recipient_id"`
	Page         int          `jsonp:"page,omitempty"`
	Size         int          `jsonp:"size,omitempty"`
}

func (this *BalanceQuery) UrlQuery() string {
	m, _ := json.EncodeAsMap(this)
	return maps.ToUrlQuery(m)
}

func NewBalanceQuery() *BalanceQuery {
	return &BalanceQuery{}
}
