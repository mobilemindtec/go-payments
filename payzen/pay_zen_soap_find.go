package payzen

import (
	"encoding/xml"
	_"fmt"
)


/************************************************************/
/** request findPaymentsResponse */
/************************************************************/

type SOAPFindPayments struct {
	XMLName xml.Name    `xml:"v5:findPayments"`
	QueryRequest *SOAPQueryRequest `xml:""`
}

func NewSOAPFindPayments() *SOAPFindPayments {
	FindPayments := new(SOAPFindPayments)
	FindPayments.QueryRequest = new(SOAPQueryRequest)
	return FindPayments
}


type SOAPFindPaymentsResponse struct {
	XMLName xml.Name `xml:"findPaymentsResponse"`
	FindPaymentsResult *SOAPFindPaymentsResult `xml:""`
}

type SOAPFindPaymentsResult struct {
	XMLName xml.Name `xml:"findPaymentsResult"`
	RequestId string `xml:"requestId"`
	CommonResponse *SOAPCommonResponse `xml:"commonResponse"`
	OrderResponse *SOAPOrderResponse `xml:"orderResponse"`
	TransactionItem []*SOAPTransactionItem `xml:"transactionItem"`
}

type SOAPTransactionItem struct {
	XMLName xml.Name `xml:"transactionItem"`
	TransactionUuid  string `xml:"transactionUuid"`
	TransactionStatusLabel string `xml:"transactionStatusLabel"`
	Amount string `xml:"amount"`
	Currency string `xml:"currency"`
	ExpectedCaptureDate string `xml:"expectedCaptureDate"`
}

func NewSOAPFindPaymentsResponse() *SOAPFindPaymentsResponse{
	FindPaymentsResponse := new(SOAPFindPaymentsResponse)
	FindPaymentsResponse.FindPaymentsResult = new(SOAPFindPaymentsResult)
	//FindPaymentsResponse.FindPaymentsResult.TransactionItem = []*SOAPTransactionItem{}
	return FindPaymentsResponse
}


/************************************************************/
/** request getPaymentDetails */
/************************************************************/

type SOAPGetPaymentDetails struct {
	XMLName xml.Name    `xml:"v5:getPaymentDetails"`
	QueryRequest *SOAPQueryRequest `xml:""`
	ExtendedResponseRequest *SOAPExtendedResponseRequest `xml:""`
}

type SOAPExtendedResponseRequest struct {
	XMLName xml.Name    `xml:"extendedResponseRequest"`
	IsNsuRequested int `xml:""`
}

func NewSOAPGetPaymentDetails() *SOAPGetPaymentDetails {
	GetPaymentDetails := new(SOAPGetPaymentDetails)
	GetPaymentDetails.QueryRequest = new(SOAPQueryRequest)
	return GetPaymentDetails
}

func NewSOAPGetPaymentDetailsWithNsu() *SOAPGetPaymentDetails {
	GetPaymentDetails := new(SOAPGetPaymentDetails)
	GetPaymentDetails.QueryRequest = new(SOAPQueryRequest)
	GetPaymentDetails.ExtendedResponseRequest = new(SOAPExtendedResponseRequest)
	GetPaymentDetails.ExtendedResponseRequest.IsNsuRequested = 1
	return GetPaymentDetails
}


type SOAPGetPaymentDetailsResponse struct {
	XMLName xml.Name `xml:"getPaymentDetailsResponse"`
	GetPaymentDetailsResult *SOAPGetPaymentDetailsResult `xml:""`
}

type SOAPGetPaymentDetailsResult struct {
	XMLName xml.Name `xml:"getPaymentDetailsResult"`
	CommonResponse *SOAPCommonResponse `xml:""`
	PaymentResponse *SOAPPaymentResponse `xml:""`
	OrderResponse *SOAPOrderResponse `xml:""`
	CardResponse *SOAPCardResponse `xml:""`
	AuthorizationResponse *SOAPAuthorizationResponse `xml:""`
	CaptureResponse *SOAPCaptureResponse `xml:""`
	CustomerResponse *SOAPCustomerResponse `xml:""`
	MarkResponse *SOAPMarkResponse `xml:""`
	ThreeDSResponse *SOAPThreeDSResponse `xml:""`
	ExtraResponse *SOAPExtraResponse `xml:""`
	FraudManagementResponse *SOAPFraudManagementResponse `xml:""`
}

func NewSOAPGetPaymentDetailsResponse() *SOAPGetPaymentDetailsResponse {
	GetPaymentDetailsResponse := new(SOAPGetPaymentDetailsResponse)
	GetPaymentDetailsResponse.GetPaymentDetailsResult = new(SOAPGetPaymentDetailsResult)
	return GetPaymentDetailsResponse
}
