package payzen

import (
	"encoding/xml"
	_"fmt"
)


type SOAPEnvelope struct {
	XMLName xml.Name    `xml:"soap:Envelope"`
	XmlNSSoap string `xml:"xmlns:soap,attr"`
	XmlNS string `xml:"xmlns:v5,attr"`
	Header  *SOAPHeader `xml:",omitempty"`
	Body    *SOAPBody    `xml:",omitempty"`
}

func NewSOAPEnvelope() *SOAPEnvelope {
	soap := new(SOAPEnvelope)
	soap.XmlNSSoap = "http://www.w3.org/2003/05/soap-envelope"
	soap.XmlNS = "http://v5.ws.vads.lyra.com/"
	soap.Header = NewSOAPHeader()
	soap.Body = new(SOAPBody)
	return soap
}

type SOAPHeader struct {
	XMLName xml.Name    `xml:"soap:Header"`
	XmlNSSoap string `xml:"xmlns:soapHeader,attr"`
	ShopId string `xml:"soapHeader:shopId"`
	RequestId string `xml:"soapHeader:requestId"`
	Timestamp string `xml:"soapHeader:timestamp"`
	Mode string `xml:"soapHeader:mode"`
	AuthToken string `xml:"soapHeader:authToken"`
}

func NewSOAPHeader() *SOAPHeader {
	header := new(SOAPHeader)
	header.XmlNSSoap = "http://v5.ws.vads.lyra.com/Header"
	return header
}	

type SOAPBody struct {
	XMLName xml.Name    `xml:"soap:Body"`
	OperationRequest interface{} `xml:""`
}

type SOAPCommonRequest struct {
	XMLName xml.Name  `xml:"commonRequest"`
	PaymentSource string `xml:"paymentSource"`
	SubmissionDate string `xml:"submissionDate"`
	ContractNumber string `xml:"contractNumber"`
}

type SOAPThreeDSRequest struct {
	XMLName xml.Name    `xml:"threeDSRequest"`
	Mode string `xml:"mode"`
	Disabled string `xml:"DISABLED"`
}

type SOAPPaymentRequest struct {
	XMLName xml.Name    `xml:"paymentRequest"`
	Amount string `xml:"amount,omitempty"`
	Currency string `xml:"currency,omitempty"`	
	/*
		Data da captura solicitada apresentada no formato ISO 8601 definido pelo W3C.

		Exemplo: 2016-07-16T19:20:00Z.

		Este parâmetro é utilizado para efetuar um pagamento a prazo.

		Se o número de dias entre a data de captura e a data do dia vigente for superior ao tempo de validade da autorização, uma autorização de 1 BRL será realizada o mesmo dia da transação. Permite verificar a validade do cartão.

		A autorização para o valor total será efetuada:

		processo padrão: o dia da data de captura no banco desejada,
		processo com autorização antecipada: em função do meio de pagamento selecionado, a D- quantidade de dias correspondente ao tempo de validade de uma autorização antes da data de captura no banco desejada.
		Se você quiser receber a notificação de resultado desta solicitação de autorização, você deve configurar a regra de notificação URL de notificação em autorização por Batch no Back Office (Configuração > Regras de notificações).

		Observação : Se o prazo anterior à captura for superior a 365 dias na solicitação de pagamento, ele será automaticamente redefinido a 365 dias.
	*/
	ExpectedCaptureDate string `xml:"expectedCaptureDate,omitempty"`
	PaymentOptionCode int `xml:"paymentOptionCode,omitempty"` // numero de parcelas
	ManualValidation string `xml:"manualValidation,omitempty"` // 0 = automatico, 1 manual

	
	//TransactionId string `xml:"transactionId,omitempty"`
	//FirstInstallmentDelay string `xml:"firstInstallmentDelay,omitempty"`
}

type SOAPExtraInfo struct {
	Key string `xml:"key"`
	Value string `xml:"value"`
}

type SOAPOrderRequest struct {
	XMLName xml.Name    `xml:"orderRequest"`
	OrderId string `xml:"orderId"`
	// objecto customizado chave valor..
	ExtInfo []*SOAPExtraInfo `xml:"extInfo"`	
}

type SOAPCardRequest struct {
	XMLName xml.Name    `xml:"cardRequest"`
	Number string `xml:"number,omitempty"`
	Scheme string `xml:"scheme,omitempty"`
	ExpiryMonth string `xml:"expiryMonth,omitempty"`
	ExpiryYear string `xml:"expiryYear,omitempty"`
	CardSecurityCode string `xml:"cardSecurityCode,omitempty"`
	CardHolderBirthDay string `xml:"cardHolderBirthDay,omitempty"`
	CardHolderName string `xml:"cardHolderName,omitempty"`	
	PaymentToken string `xml:"paymentToken,omitempty"`
}

type SOAPQueryRequest struct {
	XMLName xml.Name    `xml:"queryRequest"`
	OrderId string `xml:"orderId,omitempty"`
	Uuid string `xml:"uuid,omitempty"`
	PaymentToken string `xml:"paymentToken,omitempty"`
}

type SOAPBillingDetails struct {
	XMLName xml.Name    `xml:"billingDetails"`
	Reference string `xml:"reference,omitempty"`
	Title string `xml:"title,omitempty"`
	Type string `xml:"type,omitempty"` // PRIVATE or COMPANY
	FirstName string `xml:"firstName,omitempty"`
	LastName string `xml:"lastName,omitempty"`
	PhoneNumber string `xml:"phoneNumber,omitempty"`
	Email string `xml:"email,omitempty"`
	StreetNumber string `xml:"streetNumber,omitempty"`
	Address string `xml:"address,omitempty"`
	Address2 string `xml:"address2,omitempty"`
	District string `xml:"district,omitempty"`
	ZipCode string `xml:"zipCode,omitempty"`
	City string `xml:"city,omitempty"`
	State string `xml:"state,omitempty"`
	Country string `xml:"country,omitempty"`
	Language string `xml:"language,omitempty"`
	CellPhoneNumber string `xml:"cellPhoneNumber,omitempty"`
	IdentityCode string `xml:"identityCode,omitempty"`

}

type SOAPShippingDetails struct {
	XMLName xml.Name    `xml:"shippingDetails"`
	Type string `xml:"type"` // PRIVATE or COMPANY
	Title string `xml:"title"`
	FirstName string `xml:"firstName"`
	LastName string `xml:"lastName"`
	PhoneNumber string `xml:"phoneNumber"`
	Email string `xml:"email"`
	StreetNumber string `xml:"streetNumber"`
	Address string `xml:"address"`
	Address2 string `xml:"address2"`
	District string `xml:"district"`
	ZipCode string `xml:"zipCode"`
	City string `xml:"city"`
	State string `xml:"state"`
	Country string `xml:"country"`
	Language string `xml:"language"`
	CellPhoneNumber string `xml:"cellPhoneNumber"`
	DeliveryCompanyName string `xml:"deliveryCompanyName"`
	ShippingSpeed string `xml:"shippingSpeed"`
	ShippingMethod string `xml:"shippingMethod"`
	LegalName string `xml:"legalName"`
	IdentityCode string `xml:"identityCode"`	
}

type SOAPTechRequest struct {
	XMLName xml.Name    `xml:"techRequest"`
	BrowserUserAgent string `xml:"browserUserAgent"`
	BrowserAccept string `xml:"browserAccept"`
	
}

type SOAPExtraDetails struct {
	XMLName xml.Name    `xml:"extraDetails"`
	IpAddress string `xml:"ipAddress"`
	FingerPrintId string `xml:"fingerPrintId"`
	
}

type SOAPCustomerRequest struct {
	XMLName xml.Name    `xml:"customerRequest"`
	BillingDetails *SOAPBillingDetails `xml:""`
	ShippingDetails *SOAPShippingDetails `xml:"-"`
	ExtraDetails *SOAPExtraDetails `xml:"-"`
}


type SOAPResponseEnvelop struct {
	XMLName xml.Name    `xml:"Envelope"`
	//XmlNSSoap string `xml:"xmlns:soap,attr"`	
	Header *SOAPResponseHeader `xml:""`
	Body * SOAPResponseBody `xml:""`
}

func NewSOAPResponseEnvelop() *SOAPResponseEnvelop {
	soap := new(SOAPResponseEnvelop)
	//soap.XmlNSSoap = "http://www.w3.org/2003/05/soap-envelope"
	soap.Header = NewSOAPResponseHeader()
	soap.Body = new(SOAPResponseBody)
	return soap
}

type SOAPResponseHeader struct {
	XMLName xml.Name    `xml:"Header"`
	///XmlNSSoap string `xml:"xmlns:env,attr"`
	ShopId string `xml:"shopId"`
	RequestId string `xml:"requestId"`
	Timestamp string `xml:"timestamp"`
	Mode string `xml:"mode"`
	AuthToken string `xml:"authToken"`		
}

func NewSOAPResponseHeader() *SOAPResponseHeader {
	header := new(SOAPResponseHeader)
	//header.XmlNSSoap = "http://www.w3.org/2003/05/soap-envelope"
	return header
}

type SOAPResponseBody struct {
	XMLName xml.Name    `xml:"Body"`
	
	Return *SOAPReturn `xml:",omitempty"`
	
	CreatePaymentResponse *SOAPCreatePaymentResponse `xml:",omitempty"`
	CancelPaymentResponse *SOAPCancelPaymentResponse `xml:",omitempty"`
	UpdatePaymentResponse *SOAPUpdatePaymentResponse `xml:",omitempty"`
	DuplicatePaymentResponse *SOAPDuplicatePaymentResponse `xml:",omitempty"`

	CapturePaymentResponse *SOAPCapturePaymentResponse `xml:",omitempty"`
	ValidatePaymentResponse *SOAPValidatePaymentResponse `xml:",omitempty"`
	RefundPaymentResponse *SOAPRefundPaymentResponse `xml:",omitempty"`

	CreateTokenResponse *SOAPCreateTokenResponse `xml:",omitempty"`
	UpdateTokenResponse *SOAPUpdateTokenResponse `xml:",omitempty"`
	CancelTokenResponse *SOAPCancelTokenResponse `xml:",omitempty"`
	ReactiveTokenResponse *SOAPReactiveTokenResponse `xml:",omitempty"`
	GetTokenDetailsResponse *SOAPGetTokenDetailsResponse `xml:",omitempty"`
	
	
	FindPaymentsResponse *SOAPFindPaymentsResponse `xml:",omitempty"`
	GetPaymentDetailsResponse *SOAPGetPaymentDetailsResponse `xml:",omitempty"`

	CreateSubscriptionResponse *SOAPCreateSubscriptionResponse `xml:",omitempty"`
	GetSubscriptionDetailsResponse *SOAPGetSubscriptionDetailsResponse `xml:",omitempty"`
	CancelSubscriptionResponse *SOAPCancelSubscriptionResponse `xml:",omitempty"`
	UpdateSubscriptionResponse *SOAPUpdateSubscriptionResponse `xml:",omitempty"`
}

func (this *SOAPResponseBody) GetCommonResponse() *SOAPCommonResponse {

	// payments
	if this.CreatePaymentResponse != nil {
		return this.CreatePaymentResponse.CreatePaymentResult.CommonResponse
	}

	if this.CancelPaymentResponse != nil {
		return this.CancelPaymentResponse.CancelPaymentResult.CommonResponse
	} 

	if this.UpdatePaymentResponse != nil {
		return this.UpdatePaymentResponse.UpdatePaymentResult.CommonResponse
	}

	if this.DuplicatePaymentResponse != nil {
		return this.DuplicatePaymentResponse.DuplicatePaymentResult.CommonResponse
	}

	// payment opts

	if this.CapturePaymentResponse != nil {
		return this.CapturePaymentResponse.CapturePaymentResult.CommonResponse
	}

	if this.ValidatePaymentResponse != nil {
		return this.ValidatePaymentResponse.ValidatePaymentResult.CommonResponse	
	}

	if this.RefundPaymentResponse != nil {
		return this.RefundPaymentResponse.RefundPaymentResult.CommonResponse
	}

	// token
	if this.CreateTokenResponse != nil {
		return this.CreateTokenResponse.CreateTokenResult.CommonResponse
	}

	if this.UpdateTokenResponse != nil {
		return this.UpdateTokenResponse.UpdateTokenResult.CommonResponse
	}

	if this.CancelTokenResponse != nil {
		return this.CancelTokenResponse.CancelTokenResult.CommonResponse
	}

	if this.ReactiveTokenResponse != nil {
		return this.ReactiveTokenResponse.ReactiveTokenResult.CommonResponse
	}

	if this.GetTokenDetailsResponse != nil {
		return this.GetTokenDetailsResponse.GetTokenDetailsResult.CommonResponse
	}

	// find 	
	if this.FindPaymentsResponse != nil {
		return this.FindPaymentsResponse.FindPaymentsResult.CommonResponse
	}

	if this.GetPaymentDetailsResponse != nil {
		return this.GetPaymentDetailsResponse.GetPaymentDetailsResult.CommonResponse
	}

	// subscription
	if this.CreateSubscriptionResponse != nil {
		return this.CreateSubscriptionResponse.CreateSubscriptionResult.CommonResponse
	}

	if this.GetSubscriptionDetailsResponse != nil {
		return this.GetSubscriptionDetailsResponse.GetSubscriptionDetailsResult.CommonResponse
	}

	if this.CancelSubscriptionResponse != nil {
		return this.CancelSubscriptionResponse.CancelSubscriptionResult.CommonResponse
	}

	if this.UpdateSubscriptionResponse != nil {
		return this.UpdateSubscriptionResponse.UpdateSubscriptionResult.CommonResponse
	}

	return nil
}

type SOAPReturn struct {
	XMLName xml.Name    `xml:"return"`
	CommonResponse	*SOAPCommonResponse `xml:""`
}

type SOAPCommonResponse struct {
	XMLName xml.Name    `xml:"commonResponse"`
	ResponseCode string `xml:"responseCode"`
	ResponseCodeDetail string `xml:"responseCodeDetail"`
	TransactionStatusLabel string `xml:"transactionStatusLabel"`
	ShopId string `xml:"shopId"`
	PaymentSource string `xml:"paymentSource"` // EC, MOTO,  CC, OTHER -> EC pedido eletrônico
	SubmissionDate string `xml:"submissionDate"`
	ContractNumber string `xml:"contractNumber"`
	PaymentToken string `xml:"paymentToken"`
}



type SOAPCardResponse struct {
	XMLName xml.Name    `xml:"cardResponse"`
	Number string `xml:"number"`
	Scheme string `xml:"scheme"`
	Brand string `xml:"brand"`
	Country string `xml:"country"`
	ProductCode string `xml:"productCode"`
	BankCode string `xml:"bankCode"`
	ExpiryMonth string `xml:"expiryMonth"`
	ExpiryYear string `xml:"expiryYear"`
}

type SOAPAuthorizationResponse struct{
	XMLName xml.Name    `xml:"authorizationResponse"`
	Mode string `xml:"mode"`
	Amount string `xml:"amount"`
	Currency string `xml:"currency"`
	Date string `xml:"date"`
	Number string `xml:"number"`
	Result string `xml:"result"`
}

type SOAPCaptureResponse struct {
	XMLName xml.Name    `xml:"captureResponse"`
	Date string `xml:"date"`
	Number string `xml:"number"`
	ReconciliationStatus string `xml:"reconciliationStatus"`
	RefundAmount string `xml:"refundAmount"`
	RefundCurrency string `xml:"refundCurrency"`
	Chargeback string `xml:"chargeback"`
}

type SOAPCustomerResponse struct {
	XMLName xml.Name    `xml:"customerResponse"`
	BillingDetails *SOAPBillingDetails `xml:""`
	ShippingDetails *SOAPShippingDetails `xml:""`
	ExtraDetails *SOAPExtraDetails `xml:""`
}

type SOAPMarkResponse struct {
	XMLName xml.Name    `xml:"markResponse"`
	Amount  string `xml:"amount"`
	Currency  string `xml:"currency"`
	Date  string `xml:"date"`
	Number  string `xml:"number"`
	Result  string `xml:"result"`
}

type SOAPThreeDSResponse struct {
	XMLName xml.Name    `xml:"threeDSResponse"`
	AuthenticationRequestData  SOAPAuthenticationRequestData `xml:""`
	AuthenticationResultData  SOAPAuthenticationResultData `xml:""`
}

type SOAPAuthenticationRequestData struct {
	XMLName xml.Name    `xml:"authenticationRequestData"`
	ThreeDSAcctId string `xml:"threeDSAcctId"`
	ThreeDSAcsUrl string `xml:"threeDSAcsUrl"`
	ThreeDSBrand string `xml:"threeDSBrand"`
	ThreeDSEncodedPareq string `xml:"threeDSEncodedPareq"`
	ThreeDSEnrolled string `xml:"threeDSEnrolled"`
	ThreeDSRequestId string `xml:"threeDSRequestId"`
}

type SOAPAuthenticationResultData struct {
	XMLName xml.Name    `xml:"authenticationResultData"`
	TransactionCondition string `xml:"transactionCondition"`
	Enrolled string `xml:"enrolled"`
	Status string `xml:"status"`
	Eci string `xml:"eci"`
	Xid string `xml:"xid"`
	Cavv string `xml:"cavv"`
	Brand string `xml:"brand"`
}


type SOAPRiskControl struct {
	XMLName xml.Name `xml:"riskControl"`
	Name string `xml:"name"`
	Result string `xml:"result"`
}

type SOAPRiskAnalysis struct {
	XMLName xml.Name `xml:"riskAnalysis"`
	Score string `xml:"score"`
	ResultCode string `xml:"resultCode"`
	Status string `xml:"status"`
	RequestId string `xml:"requestId"`	
	ExtInfo []*SOAPExtraInfo `xml:"extInfo"`

}

//https://payzen.io/pt-BR/webservices-payment/implementation-webservices-v5/fraudmanagementresponse.html
type SOAPFraudManagementResponse struct {
	XMLName xml.Name `xml:"fraudManagementResponse"`
	RiskControl *SOAPRiskControl `xml:""`
	RiskAnalysis *SOAPRiskAnalysis `xml:""`
	RiskAssessment string `xml:"riskAssessment"`
}

type SOAPExtraResponse struct {
	XMLName xml.Name    `xml:"extraResponse"`
	PaymentOptionCode  string `xml:"paymentOptionCode"`
	PaymentOptionOccNumb  string `xml:"paymentOptionOccNumb"`
}

type SOAPPaymentResponse struct {
	XMLName xml.Name    `xml:"paymentResponse"`
	TransactionId string `xml:"transactionId"`
	TransactionUuid string `xml:"transactionUuid"`
	Amount string `xml:"amount"`
	Currency string `xml:"currency"`
	ExpectedCaptureDate string `xml:"expectedCaptureDate"`
	PaymentType string `xml:"paymentType"` // SINGLE
	PaymentError string `xml:"paymentError"`

	EffectiveAmount string `xml:"effectiveAmount"`
	EffectiveCurrency string `xml:"effectiveCurrency"`
	OperationType string `xml:"operationType"` // 0 debito, 1 reembolso
	CreationDate string `xml:"creationDate"`
	ExternalTransactionId string `xml:"externalTransactionId"`
	SequenceNumber string `xml:"sequenceNumber"`

}

type SOAPOrderResponse struct {
	XMLName xml.Name    `xml:"orderResponse"`
	OrderId string `xml:"orderId"`
	ExtInfo []*SOAPExtraInfo `xml:"extInfo"`
}