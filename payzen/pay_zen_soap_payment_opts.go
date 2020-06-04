package payzen

import (
	"encoding/xml"
)


/************************************************************/
/** request capturePayment */
/************************************************************/

type SOAPCapturePayment struct { 
	XMLName xml.Name    `xml:"v5:capturePayment"`
	SettlementRequest *SOAPSettlementRequest `xml:""`
}

type SOAPSettlementRequest struct {
	XMLName xml.Name    `xml:"settlementRequest"`
	TransactionUuids string `xml:"transactionUuids"`
	Commission string `xml:"commission"`
	Date string `xml:"date"`
}

func NewSOAPCapturePayment() *SOAPCapturePayment {
	soap := new(SOAPCapturePayment)
	soap.SettlementRequest = new(SOAPSettlementRequest)
	return soap
}

type SOAPCapturePaymentResponse struct {
	XMLName xml.Name    `xml:"capturePaymentResponse"`
	CapturePaymentResult *SOAPCapturePaymentResult `xml:""`
}

type SOAPCapturePaymentResult struct {
	XMLName xml.Name    `xml:"capturePaymentResult"`
	CommonResponse *SOAPCommonResponse `xml:"commonResponse"`
}

func NewSOAPCreateCaptureResponse() *SOAPCapturePaymentResponse {
	soap := new(SOAPCapturePaymentResponse)
	soap.CapturePaymentResult = new(SOAPCapturePaymentResult)
	return soap
}

/************************************************************/
/** validate payment */
/************************************************************/

type SOAPValidatePayment struct { 
	XMLName xml.Name    `xml:"v5:validatePayment"`
	CommonRequest *SOAPCommonRequest `xml:""`
	QueryRequest *SOAPQueryRequest `xml:""`
}


func NewSOAPValidatePayment() *SOAPValidatePayment {
	soap := new(SOAPValidatePayment)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.QueryRequest = new(SOAPQueryRequest)
	return soap
}

type SOAPValidatePaymentResponse struct {
	XMLName xml.Name    `xml:"validatePaymentResponse"`
	ValidatePaymentResult *SOAPValidatePaymentResult `xml:""`
}

type SOAPValidatePaymentResult struct {
	XMLName xml.Name    `xml:"validatePaymentResult"`
	CommonResponse *SOAPCommonResponse `xml:"commonResponse"`
}

func NewSOAPValidatePaymentResponse() *SOAPValidatePaymentResponse {
	soap := new(SOAPValidatePaymentResponse)
	soap.ValidatePaymentResult = new(SOAPValidatePaymentResult)
	return soap
}

/************************************************************/
/** refund payment */
/************************************************************/

type SOAPRefundPayment struct { 
	XMLName xml.Name    `xml:"v5:refundPayment"`
	CommonRequest *SOAPCommonRequest `xml:""`
	PaymentRequest *SOAPPaymentRequest `xml:""`
	QueryRequest *SOAPQueryRequest `xml:""`
}


func NewSOAPRefundPayment() *SOAPRefundPayment {
	soap := new(SOAPRefundPayment)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.QueryRequest = new(SOAPQueryRequest)
	soap.PaymentRequest = new(SOAPPaymentRequest)
	return soap
}

type SOAPRefundPaymentResponse struct {
	XMLName xml.Name    `xml:"refundPaymentResponse"`
	RefundPaymentResult *SOAPRefundPaymentResult `xml:""`
}

type SOAPRefundPaymentResult struct {
	XMLName xml.Name    `xml:"refundPaymentResult"`
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

func NewSOAPRefundPaymentResponse() *SOAPRefundPaymentResponse {
	soap := new(SOAPRefundPaymentResponse)
	soap.RefundPaymentResult = new(SOAPRefundPaymentResult)
	soap.RefundPaymentResult.CommonResponse = new(SOAPCommonResponse)
	soap.RefundPaymentResult.PaymentResponse = new(SOAPPaymentResponse)
	soap.RefundPaymentResult.OrderResponse = new(SOAPOrderResponse)
	return soap
}