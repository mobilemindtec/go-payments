package payzen

import (
	"encoding/xml"
	_"fmt"
)

/*

Descrição da regra da assinatura.

O valor esperado neste atributo é um string de dígitos conforme à especificação iCalendar, ou Internet Calendar, apresentado na RFC5545 (ver http://tools.ietf.org/html/rfc5545).

Por razões técnicas. não é possível definir períodos de assinatura inferiores a um dia.

As palavras chaves SECONDLY" / "MINUTELY" / "HOURLY não são levados em conta.

Exemplos:
Para definir parcelas de pagamento que ocorrem o último dia de cada mês, durante 12 meses, a regra se escreve:

RRULE:FREQ=MONTHLY;BYMONTHDAY=28,29,30,31;BYSETPOS=-1;COUNT=12

Esta regra significa que se o mês corrente não contém um dia 31, então o motor levará em conta o dia 30. Se o dia 30 não existe, então ele levará em conta o dia 29, e assim por diante até o dia 28.

Outra versão desta regra é: RRULE:FREQ=MONTHLY;COUNT=5;BYMONTHDAY=-1

Para definir parcelas de pagamento que ocorrem o dia 10 de cada mês, durante 12 meses, a regra de assinatura se escreve da seguinte forma: RRULE:FREQ=MONTHLY;COUNT=12;BYMONTHDAY=10
Para definir parcelas de pagamento que ocorrem todo trimestre, até o 31/12/2016: RRULE:FREQ=YEARLY;BYMONTHDAY=1;BYMONTH=1,4,7,10;UNTIL=20161231
As parcelas ocorrerão todo dia 1° de janeiro, abril, julho e outubro. A quantidade total deles depende da data de início da assinatura .

Para maiores detalhes e exemplos, você pode consultar o site http://recurrance.sourceforge.net/.

https://payzen.io/pt-BR//webservices-payment/implementation-webservices-v5/efetuar-operacoes-especificas-aos-pagamento-por-codigo.html

*/


/*

	create subscription

*/

type SOAPCreateSubscription struct {
	XMLName xml.Name    `xml:"v5:createSubscription"`
	CommonRequest *SOAPCommonRequest `xml:"commonRequest"`
	OrderRequest *SOAPOrderRequest `xml:"orderRequest"`
	SubscriptionRequest *SOAPSubscriptionRequest `xml:"subscriptionRequest"`
	CardRequest *SOAPCardRequest `xml:"cardRequest"`
}

type SOAPSubscriptionRequest struct {
	XMLName xml.Name    `xml:"subscriptionRequest"`

	EffectDate string `xml:"effectDate"` // 2016-07-16T19:20Z
	Amount string `xml:"amount"`
	Currency string `xml:"currency"`
	InitialAmount string `xml:"initialAmount"`
	InitialAmountNumber string `xml:"initialAmountNumber"`
	Rrule string `xml:"rrule"`
	SubscriptionId string `xml:"subscriptionId"`
	Description string `xml:"description"`		
}

func NewSOAPCreateSubscription() *SOAPCreateSubscription {
	soap := new(SOAPCreateSubscription)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.OrderRequest = new(SOAPOrderRequest)
	soap.CardRequest = new(SOAPCardRequest)
	soap.SubscriptionRequest = new(SOAPSubscriptionRequest)
	return soap
}

type SOAPCreateSubscriptionResponse struct {
	XMLName xml.Name    `xml:"createSubscriptionResponse"`
	CreateSubscriptionResult *SOAPCreateSubscriptionResult `xml:""`
}

type SOAPCreateSubscriptionResult struct {
	XMLName xml.Name    `xml:"createSubscriptionResult"`
	CommonResponse *SOAPCommonResponse `xml:""`
	SubscriptionResponse *SOAPSubscriptionResponse `xml:""`
}

type SOAPSubscriptionResponse struct {
	XMLName xml.Name    `xml:"subscriptionResponse"`
	SubscriptionId string `xml:"subscriptionId"`

	EffectDate string `xml:"effectDate"` // 2016-07-16T19:20Z
	InitialAmountNumber int64 `xml:"initialAmountNumber"`
	Rrule string `xml:"rrule"`	
	Description string `xml:"description"`
	
	PastPaymentsNumber int64 `xml:"pastPaymentNumber"`
	TotalPaymentsNumber int64 `xml:"totalPaymentNumber"`
	CancelDate string `xml:"cancelDate"`	
}

func NewSOAPCreateSubscriptionResponse() *SOAPCreateSubscriptionResponse {
	soap := new(SOAPCreateSubscriptionResponse)
	soap.CreateSubscriptionResult = new(SOAPCreateSubscriptionResult)
	soap.CreateSubscriptionResult.SubscriptionResponse = new(SOAPSubscriptionResponse)
	return soap
}


/*

	get subscription details

*/

type SOAPSubscriptionQueryRequest struct {
	XMLName xml.Name    `xml:"queryRequest"`
	PaymentToken string `xml:"paymentToken"`
	SubscriptionId string `xml:"subscriptionId"`
}

type SOAPGetSubscriptionDetails struct {
	XMLName xml.Name    `xml:"v5:getSubscriptionDetails"`
	QueryRequest *SOAPSubscriptionQueryRequest `xml:""`
}	

func NewSOAPGetSubscriptionDetails() *SOAPGetSubscriptionDetails {
	soap := new(SOAPGetSubscriptionDetails)
	soap.QueryRequest = new(SOAPSubscriptionQueryRequest)
	return soap
}

type SOAPGetSubscriptionDetailsResponse struct {
	XMLName xml.Name `xml:"getSubscriptionDetailsResponse"`
	GetSubscriptionDetailsResult *SOAPGetSubscriptionDetailsResult `xml:""`
}

type SOAPGetSubscriptionDetailsResult struct {
	XMLName xml.Name `xml:"getSubscriptionDetailsResult"`	
	CommonResponse *SOAPCommonResponse `xml:""`
	OrderResponse *SOAPOrderResponse `xml:""`
	SubscriptionResponse *SOAPSubscriptionResponse `xml:""`
} 

func NewSOAPGetSubscriptionDetailsResponse() *SOAPGetSubscriptionDetailsResponse {
	soap := new(SOAPGetSubscriptionDetailsResponse)
	soap.GetSubscriptionDetailsResult = new(SOAPGetSubscriptionDetailsResult)
	return soap
}


/* cancel subscription */



type SOAPCancelSubscription struct {
	XMLName xml.Name    `xml:"v5:cancelSubscription"`
	QueryRequest *SOAPSubscriptionQueryRequest `xml:""`
	CommonRequest *SOAPCommonRequest `xml:"commonRequest"`
}	

func NewSOAPCancelSubscription() *SOAPCancelSubscription{
	soap := new(SOAPCancelSubscription)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.QueryRequest = new(SOAPSubscriptionQueryRequest)
	return soap
}

type SOAPCancelSubscriptionResponse struct {
	XMLName xml.Name    `xml:"cancelSubscriptionResponse"`
	CancelSubscriptionResult *SOAPCancelSubscriptionResult `xml:""`
}

type SOAPCancelSubscriptionResult struct {
	XMLName xml.Name    `xml:"cancelSubscriptionResult"`
	CommonResponse *SOAPCommonResponse `xml:""`
}	

func NewSOAPCancelSubscriptionResponse() *SOAPCancelSubscriptionResponse{
	soap := new(SOAPCancelSubscriptionResponse)
	soap.CancelSubscriptionResult = new(SOAPCancelSubscriptionResult)
	return soap
}


/* update subscription */

type SOAPUpdateSubscription struct {
	XMLName xml.Name    `xml:"v5:updateSubscription"`
	CommonRequest *SOAPCommonRequest `xml:"commonRequest"`
	QueryRequest *SOAPSubscriptionQueryRequest `xml:""`	
	SubscriptionRequest *SOAPSubscriptionRequest `xml:"subscriptionRequest"`
	//PaymentRequest *SOAPPaymentRequest `xml:"paymentRequest"`
}

func NewSOAPUpdateSubscription() *SOAPUpdateSubscription {
	soap := new(SOAPUpdateSubscription)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.QueryRequest = new(SOAPSubscriptionQueryRequest)
	//soap.PaymentRequest = new(SOAPPaymentRequest)
	soap.SubscriptionRequest = new(SOAPSubscriptionRequest)
	return soap
}

type SOAPUpdateSubscriptionResponse struct {
	XMLName xml.Name    `xml:"updateSubscriptionResponse"`
	UpdateSubscriptionResult *SOAPUpdateSubscriptionResult `xml:""`
}

type SOAPUpdateSubscriptionResult struct {
	XMLName xml.Name    `xml:"updateSubscriptionResult"`
	CommonResponse *SOAPCommonResponse `xml:""`
}


func NewSOAPUpdateSubscriptionResponse() *SOAPUpdateSubscriptionResponse {
	soap := new(SOAPUpdateSubscriptionResponse)
	soap.UpdateSubscriptionResult = new(SOAPUpdateSubscriptionResult)
	return soap
}



