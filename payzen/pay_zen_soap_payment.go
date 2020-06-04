package payzen

import (
	"encoding/xml"
)

/* create payment */

type SOAPCreatePayment struct {
	XMLName xml.Name    `xml:"v5:createPayment"`
	CommonRequest *SOAPCommonRequest `xml:""`
	ThreeDSRequest *SOAPThreeDSRequest `xml:""`
	PaymentRequest *SOAPPaymentRequest `xml:""`
	OrderRequest *SOAPOrderRequest `xml:""`
	CardRequest *SOAPCardRequest `xml:""`
	CustomerRequest *SOAPCustomerRequest `xml:""`
	TechRequest *SOAPTechRequest `xml:""`
}

func NewSOAPCreatePayment() *SOAPCreatePayment {
	soap := new(SOAPCreatePayment)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.ThreeDSRequest = new(SOAPThreeDSRequest)
	soap.PaymentRequest = new(SOAPPaymentRequest)
	soap.OrderRequest = new(SOAPOrderRequest)
	soap.CardRequest = new(SOAPCardRequest)
	soap.CustomerRequest = new(SOAPCustomerRequest)
	soap.CustomerRequest.BillingDetails = new(SOAPBillingDetails)
	soap.CustomerRequest.ExtraDetails = new(SOAPExtraDetails)
	return soap
}

type SOAPCreatePaymentResponse struct {
	XMLName xml.Name    `xml:"createPaymentResponse"`
	CreatePaymentResult *SOAPCreatePaymentResult `xml:""`	
}

type SOAPCreatePaymentResult struct {
	XMLName xml.Name    `xml:"createPaymentResult"`
	RequestId string `xml:"requestId"`
	CommonResponse	*SOAPCommonResponse `xml:""`
	PaymentResponse *SOAPPaymentResponse `xml:""`
	OrderResponse *SOAPOrderResponse `xml:""`
	CardResponse *SOAPCardResponse `xml:""`
	AuthorizationResponse *SOAPAuthorizationResponse `xml:""`
	CaptureResponse *SOAPCaptureResponse `xml:""`
	CustomerResponse *SOAPCustomerResponse `xml:""`
	MarkResponse *SOAPMarkResponse `xml:""`
	ExtraResponse *SOAPExtraResponse `xml:""`
	ThreeDSResponse *SOAPThreeDSResponse `xml:",omitempty"`	
	FraudManagementResponse *SOAPFraudManagementResponse `xml:",omitempty"`		
}

func NewSOAPCreatePaymentResponse() *SOAPCreatePaymentResponse {
	soap := new(SOAPCreatePaymentResponse)
	soap.CreatePaymentResult = new(SOAPCreatePaymentResult)
	soap.CreatePaymentResult.CommonResponse = new(SOAPCommonResponse)
	soap.CreatePaymentResult.PaymentResponse = new(SOAPPaymentResponse)
	soap.CreatePaymentResult.OrderResponse = new(SOAPOrderResponse)
	return soap
}

/* cancel payment */

type SOAPCancelPayment struct {
	XMLName xml.Name    `xml:"v5:cancelPayment"`
	CommonRequest *SOAPCommonRequest `xml:""`
	QueryRequest *SOAPQueryRequest `xml:""`
}

func NewSOAPCancelPayment() *SOAPCancelPayment {
	soap := new(SOAPCancelPayment)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.QueryRequest = new(SOAPQueryRequest)
	return soap
}

type SOAPCancelPaymentResponse struct {
	XMLName xml.Name `xml:"cancelPaymentResponse"`
	CancelPaymentResult *SOAPCancelPaymentResult `xml:""`
}


type SOAPCancelPaymentResult struct {
	XMLName xml.Name `xml:"cancelPaymentResult"`
	CommonResponse *SOAPCommonResponse `xml:""`
}

func NewSOAPCancelPaymentResponse() *SOAPCancelPaymentResponse {
	soap := new(SOAPCancelPaymentResponse)
	soap.CancelPaymentResult = new(SOAPCancelPaymentResult)
	return soap
}


/* update payment */

type SOAPUpdatePayment struct {
	XMLName xml.Name    `xml:"v5:updatePayment"`
	CommonRequest *SOAPCommonRequest `xml:""`
	QueryRequest *SOAPQueryRequest `xml:""`
	PaymentRequest *SOAPPaymentRequest `xml:""`	
}

func NewSOAPUpdatePayment() *SOAPUpdatePayment {
	soap := new(SOAPUpdatePayment)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.QueryRequest = new(SOAPQueryRequest)
	soap.PaymentRequest = new(SOAPPaymentRequest)
	return soap
}

type SOAPUpdatePaymentResponse struct {
	XMLName xml.Name    `xml:"updatePaymentResponse"`
	UpdatePaymentResult *SOAPUpdatePaymentResult `xml:""`	
}

type SOAPUpdatePaymentResult struct {
	XMLName xml.Name    `xml:"updatePaymentResult"`	
	CommonResponse	*SOAPCommonResponse `xml:""`
	PaymentResponse *SOAPPaymentResponse `xml:""`
	OrderResponse *SOAPOrderResponse `xml:""`
	CardResponse *SOAPCardResponse `xml:""`
	AuthorizationResponse *SOAPAuthorizationResponse `xml:""`
	CaptureResponse *SOAPCaptureResponse `xml:""`
	CustomerResponse *SOAPCustomerResponse `xml:""`
	MarkResponse *SOAPMarkResponse `xml:""`
	ExtraResponse *SOAPExtraResponse `xml:""`
	ThreeDSResponse *SOAPThreeDSResponse `xml:",omitempty"`	
	FraudManagementResponse *SOAPFraudManagementResponse `xml:",omitempty"`	
}

func NewSOAPUpdatePaymentResponse() *SOAPUpdatePaymentResponse {
	soap := new(SOAPUpdatePaymentResponse)
	soap.UpdatePaymentResult = new(SOAPUpdatePaymentResult)
	soap.UpdatePaymentResult.CommonResponse = new(SOAPCommonResponse)
	soap.UpdatePaymentResult.PaymentResponse = new(SOAPPaymentResponse)
	soap.UpdatePaymentResult.OrderResponse = new(SOAPOrderResponse)
	return soap
}



/*
	duplicate payment
*/

type SOAPDuplicatePayment struct {
	XMLName xml.Name    `xml:"v5:duplicatePayment"`
	CommonRequest *SOAPCommonRequest `xml:""`
	PaymentRequest *SOAPPaymentRequest `xml:""`	
	QueryRequest *SOAPQueryRequest `xml:""`
	OrderRequest *SOAPOrderRequest `xml:"`
}

func NewSOAPDuplicatePayment() *SOAPDuplicatePayment {
	soap := new(SOAPDuplicatePayment)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.QueryRequest = new(SOAPQueryRequest)
	soap.PaymentRequest = new(SOAPPaymentRequest)
	soap.OrderRequest = new(SOAPOrderRequest)
	return soap
}

type SOAPDuplicatePaymentResponse struct {
	XMLName xml.Name    `xml:"duplicatePaymentResponse"`
	DuplicatePaymentResult *SOAPDuplicatePaymentResult `xml:""`	
}

type SOAPDuplicatePaymentResult struct {
	XMLName xml.Name    `xml:"duplicatePaymentResult"`	
	CommonResponse	*SOAPCommonResponse `xml:""`
	PaymentResponse *SOAPPaymentResponse `xml:""`
	OrderResponse *SOAPOrderResponse `xml:""`
	CardResponse *SOAPCardResponse `xml:""`
	AuthorizationResponse *SOAPAuthorizationResponse `xml:""`
	CaptureResponse *SOAPCaptureResponse `xml:""`
	CustomerResponse *SOAPCustomerResponse `xml:""`
	MarkResponse *SOAPMarkResponse `xml:""`
	ExtraResponse *SOAPExtraResponse `xml:""`
	ThreeDSResponse *SOAPThreeDSResponse `xml:",omitempty"`	
	FraudManagementResponse *SOAPFraudManagementResponse `xml:",omitempty"`	
}

func NewSOAPDuplicatePaymentResponse() *SOAPDuplicatePaymentResponse {
	soap := new(SOAPDuplicatePaymentResponse)
	soap.DuplicatePaymentResult = new(SOAPDuplicatePaymentResult)
	soap.DuplicatePaymentResult.CommonResponse = new(SOAPCommonResponse)
	soap.DuplicatePaymentResult.PaymentResponse = new(SOAPPaymentResponse)
	soap.DuplicatePaymentResult.OrderResponse = new(SOAPOrderResponse)
	return soap
}
