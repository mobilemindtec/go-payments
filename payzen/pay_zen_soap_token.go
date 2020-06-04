package payzen

import (
	"encoding/xml"
	_"fmt"
)

/*
	ractive token
*/

type SOAPReactiveToken struct {
	XMLName xml.Name    `xml:"v5:reactivateToken"`
	QueryRequest *SOAPQueryRequest `xml:"queryRequest"`
}

func NewSOAPReactiveToken() *SOAPReactiveToken {
	soap := new(SOAPReactiveToken)
	soap.QueryRequest = new(SOAPQueryRequest)
	return soap	
}

type SOAPReactiveTokenResponse struct {
	XMLName xml.Name    `xml:"reactiveTokenResponse"`
	ReactiveTokenResult *SOAPReactiveTokenResult `xml:""`
}

type SOAPReactiveTokenResult struct {
	XMLName xml.Name    `xml:"reactiveTokenResult"`	
	CommonResponse *SOAPCommonResponse `xml:"commonResponse"`
	AuthorizationResponse *SOAPAuthorizationResponse `xml:"authorizationResponse"`
}

func NewSOAPReactiveTokenResponse() *SOAPReactiveTokenResponse {
	soap := new(SOAPReactiveTokenResponse)
	soap.ReactiveTokenResult = new(SOAPReactiveTokenResult)
	soap.ReactiveTokenResult.CommonResponse = new(SOAPCommonResponse)
	return soap
}

/*
	update token
*/

type SOAPUpdateToken struct {
	XMLName xml.Name    `xml:"v5:updateToken"`
	CommonRequest *SOAPCommonRequest `xml:"commonRequest"`
	QueryRequest *SOAPQueryRequest `xml:"queryRequest"`
	CardRequest *SOAPCardRequest `xml:"cardRequest"`
	CustomerRequest *SOAPCustomerRequest `xml:"customerRequest"`
}

func NewSOAPUpdateToken() *SOAPUpdateToken {
	soap := new(SOAPUpdateToken)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.CardRequest = new(SOAPCardRequest)
	soap.QueryRequest = new(SOAPQueryRequest)
	soap.CustomerRequest = new(SOAPCustomerRequest)
	soap.CustomerRequest.BillingDetails = new(SOAPBillingDetails)
	soap.CustomerRequest.ExtraDetails = new(SOAPExtraDetails)
	return soap	
}

type SOAPUpdateTokenResponse struct {
	XMLName xml.Name    `xml:"updateTokenResponse"`
	UpdateTokenResult *SOAPUpdateTokenResult `xml:""`
}

type SOAPUpdateTokenResult struct {
	XMLName xml.Name    `xml:"updateTokenResult"`	
	CommonResponse *SOAPCommonResponse `xml:"commonResponse"`
	AuthorizationResponse *SOAPAuthorizationResponse `xml:"authorizationResponse"`
}

func NewSOAPUpdateTokenResponse() *SOAPUpdateTokenResponse {
	soap := new(SOAPUpdateTokenResponse)
	soap.UpdateTokenResult = new(SOAPUpdateTokenResult)
	soap.UpdateTokenResult.CommonResponse = new(SOAPCommonResponse)
	return soap
}

/*
	create token
*/

type SOAPCreateToken struct {
	XMLName xml.Name    `xml:"v5:createToken"`
	CommonRequest *SOAPCommonRequest `xml:"commonRequest"`
	CardRequest *SOAPCardRequest `xml:"cardRequest"`
	CustomerRequest *SOAPCustomerRequest `xml:"customerRequest"`
}

func NewSOAPCreateToken() *SOAPCreateToken {
	soap := new(SOAPCreateToken)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.CardRequest = new(SOAPCardRequest)
	soap.CustomerRequest = new(SOAPCustomerRequest)
	soap.CustomerRequest.BillingDetails = new(SOAPBillingDetails)
	soap.CustomerRequest.ExtraDetails = new(SOAPExtraDetails)
	return soap	
}

type SOAPCreateTokenResponse struct {
	XMLName xml.Name    `xml:"createTokenResponse"`
	CreateTokenResult *SOAPCreateTokenResult `xml:""`
}

type SOAPCreateTokenResult struct {
	XMLName xml.Name    `xml:"createTokenResult"`
	RequestId string `xml:"requestId"`
	CommonResponse *SOAPCommonResponse `xml:"commonResponse"`
	AuthorizationResponse *SOAPAuthorizationResponse `xml:"authorizationResponse"`
}

func NewSOAPCreateTokenResponse() *SOAPCreateTokenResponse {
	soap := new(SOAPCreateTokenResponse)
	soap.CreateTokenResult = new(SOAPCreateTokenResult)
	soap.CreateTokenResult.CommonResponse = new(SOAPCommonResponse)
	return soap
}

/*
	cancel token
*/

type SOAPCancelToken struct {
	XMLName xml.Name    `xml:"v5:cancelToken"`
	CommonRequest *SOAPCommonRequest `xml:"commonRequest"`
	QueryRequest *SOAPQueryRequest `xml:"queryRequest"`
}

func NewSOAPCancelToken() *SOAPCancelToken {
	soap := new(SOAPCancelToken)
	soap.CommonRequest = new(SOAPCommonRequest)
	soap.QueryRequest = new(SOAPQueryRequest)
	return soap	
}

type SOAPCancelTokenResponse struct {
	XMLName xml.Name    `xml:"cancelTokenResponse"`
	CancelTokenResult *SOAPCancelTokenResult `xml:""`
}

type SOAPCancelTokenResult struct {
	XMLName xml.Name    `xml:"cancelTokenResult"`	
	CommonResponse *SOAPCommonResponse `xml:"commonResponse"`
}

func NewSOAPCancelTokenResponse() *SOAPCancelTokenResponse {
	soap := new(SOAPCancelTokenResponse)
	soap.CancelTokenResult = new(SOAPCancelTokenResult)
	return soap
}


/*
	token details
*/

type SOAPGetTokenDetails struct {
	XMLName xml.Name    `xml:"v5:getTokenDetails"`	
	QueryRequest *SOAPQueryRequest `xml:"queryRequest"`
}

func NewSOAPGetTokenDetails() *SOAPGetTokenDetails {
	soap := new(SOAPGetTokenDetails)	
	soap.QueryRequest = new(SOAPQueryRequest)
	return soap	
}

type SOAPGetTokenDetailsResponse struct {
	XMLName xml.Name    `xml:"getTokenDetailsResponse"`
	GetTokenDetailsResult *SOAPGetTokenDetailsResult `xml:""`
}

type SOAPGetTokenDetailsResult struct {
	XMLName xml.Name    `xml:"getTokenDetailsResult"`	
	CommonResponse *SOAPCommonResponse `xml:""`	
	CardResponse *SOAPCardResponse `xml:""`
	CustomerResponse *SOAPCustomerResponse `xml:""`
	AuthorizationResponse *SOAPAuthorizationResponse `xml:""`
	TokenResponse *SOAPTokenResponse `xml:""`
}

type SOAPTokenResponse struct {
	XMLName xml.Name    `xml:"tokenResponse"`	
	CreationDate string `xml:"creationDate"`
	CancellationDate string `xml:"cancellationDate"`
}

func NewSOAPGetTokenDetailsResponse() *SOAPGetTokenDetailsResponse {
	soap := new(SOAPGetTokenDetailsResponse)
	soap.GetTokenDetailsResult = new(SOAPGetTokenDetailsResult)
	return soap
}
