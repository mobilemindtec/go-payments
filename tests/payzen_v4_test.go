package gopayments

import (
  "github.com/mobilemindtec/go-payments/payzen/v4"
  "github.com/mobilemindtec/go-payments/api"
  _"encoding/json"
  _"net/url"
  _"strings"
	"testing"
  "time"
	"fmt"
)


func createPayment() *v4.Payment {
  payment := v4.NewPayment(15.00)
  payment.Card.Number = "4970100000000007"
  payment.Card.Brand = "VISA"
  payment.Card.ExpiryMonth = 12
  payment.Card.ExpiryYear = 2025
  payment.Card.SecurityCode = "123"
  payment.Card.CardHolderName = "Ricardo Bocchi"
  payment.Card.InstallmentNumber = 1

  payment.OrderId = GenUUID()
  payment.Customer.Email = "ricardo@mobilemind.com.br"

  payment.Customer.BillingDetails.Address = "Rua Vitória"
  payment.Customer.BillingDetails.Address2 = "Ed Barzenski, Sala 8"
  payment.Customer.BillingDetails.StreetNumber = "255"
  payment.Customer.BillingDetails.ZipCode = "95700540"
  payment.Customer.BillingDetails.CellPhoneNumber = "54999767081"
  payment.Customer.BillingDetails.PhoneNumber = "5430553222"
  payment.Customer.BillingDetails.City = "Bento Goncalves"
  payment.Customer.BillingDetails.State = "RS"
  payment.Customer.BillingDetails.District = "Botafogo"
  payment.Customer.BillingDetails.FirstName = "Ricardo"
  payment.Customer.BillingDetails.LastName = "Bocchi"
  payment.Customer.BillingDetails.IdentityCode = "83361855004"  
  return payment
}


// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPayZenV4PaymentCreate
func TestPayZenV4PaymentCreate(t *testing.T) {

  payment := createPayment() 
  payment.IpnTargetUrl = "https://mobilemind.free.beeceptor.com/webhook/payzen"
  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  result, err := payzen.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v - %v", result.Message, result.Errors)
    return
  }
  
  trans := result.GetTransaction()

  if trans.TransactionStatus != v4.AUTHORISED {
    t.Errorf("Status expected: %v, returned: %v", v4.AUTHORISED, trans.TransactionStatus)
    return
  }

  if trans.PaymentStatus != v4.PAID {
    t.Errorf("PaymentStatus expected: %v, returned: %v", v4.PAID, trans.PaymentStatus)
  }


  client.Set("OrderId", payment.OrderId, 0)
  client.Set("TransactionUuid", trans.Uuid, 0)
}

// go test -v github.com/mobilemindtec/go-payments/tests -run TestPayZenV4PaymentCancelOrRefund -v
func TestPayZenV4PaymentCancelOrRefund(t *testing.T) {

  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  uuid, _ := client.Get("TransactionUuid").Result()
  result, err := payzen.PaymentCancelOrRefund(uuid, 10.00)

  if err != nil {
    t.Errorf("Erro cancel or refund: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao cancel or refund: %v - %v", result.Message, result.Errors)
    return
  }

  if !result.EmptyResponseSuccess {
    t.Errorf("empty success response expected")
    return    
  }   
}

// go test -v github.com/mobilemindtec/go-payments/tests -run TestPayZenV4PaymentCapture -v
func TestPayZenV4PaymentCapture(t *testing.T) {

  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  uuid, _ := client.Get("TransactionUuid").Result()
  result, err := payzen.PaymentCapture(uuid)

  if err != nil {
    t.Errorf("Erro ao capture: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao capture: %v - %v", result.Message, result.Errors)
    return
  }

  if !result.EmptyResponseSuccess {
    t.Errorf("empty success response expected")
    return    
  }

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPayZenV4TokenCreate
func TestPayZenV4TokenCreate(t *testing.T) {

  payment := createPayment() 
  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  result, err := payzen.TokenCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v - %v", result.Message, result.Errors)
    return
  }


  t.Log(fmt.Sprintf("result = %v", result.GetResponse()))

  trans := result.GetTransaction()

  client.Set("Token", trans.PaymentMethodToken, 0)
  client.Set("TransactionUuid", trans.Uuid, 0)
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPayZenV4TokenUpdate
func TestPayZenV4TokenUpdate(t *testing.T) {

  token, _ := client.Get("Token").Result()

  payment := createPayment() 
  payment.PaymentMethodToken = token
  payment.Card.ExpiryYear = 2028
  payment.Card.SecurityCode = "321"
  payment.Card.CardHolderName = "Ricardo Jao Bocchi"
  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  result, err := payzen.TokenUpdate(payment)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v - %v", result.Message, result.Errors)
    return
  }

  if result.IsCancelled() {
    t.Errorf("token NOT cancelled expected, but is not")
    return
  }  

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPayZenV4TokenGet
func TestPayZenV4TokenGet(t *testing.T) {

  token, _ := client.Get("Token").Result()

  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  result, err := payzen.TokenGet(token)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v - %v", result.Message, result.Errors)
    return
  }

  if !result.IsCancelled() {
    t.Errorf("token cancelled expected, but is not")
    return
  }

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPayZenV4TokenCancel
func TestPayZenV4TokenCancel(t *testing.T) {

  token, _ := client.Get("Token").Result()

  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  result, err := payzen.TokenCancel(token)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v - %v", result.Message, result.Errors)
    return
  }

  if !result.EmptyResponseSuccess {
    t.Errorf("empty success response expected")
    return    
  }  

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPayZenV4SubscriptionCreate
func TestPayZenV4SubscriptionCreate(t *testing.T) {

  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  orderId := GenUUID()
  token, _ := client.Get("Token").Result()
  subscription := v4.NewSubscription(orderId, 100, token, time.Now())

  subscription.SetRule(api.Monthly, 12, true, 0)

  result, err := payzen.SubscriptionCreate(subscription)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v - %v", result.Message, result.Errors)
    return
  }

  if !result.EmptyResponseSuccess {
    t.Errorf("empty success response expected")
    return    
  }  

  resp := result.GetResponse()

  if len(resp.SubscriptionId) == 0 {
    t.Errorf("SubscriptionId expected, but is empty")
    return    
  }  

  client.Set("SubscriptionId", resp.SubscriptionId, 0)
  client.Set("OrderId", orderId, 0)

}


// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPayZenV4SubscriptionUpdate
func TestPayZenV4SubscriptionUpdate(t *testing.T) {

  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  orderId, _ := client.Get("OrderId").Result()
  subscriptionId, _ := client.Get("SubscriptionId").Result()
  token, _ := client.Get("Token").Result()
  subscription := v4.NewSubscription(orderId, 110, token, time.Now().AddDate(0, 0, 1))
  subscription.SubscriptionId = subscriptionId
  subscription.SetRule(api.Monthly, 0, false, 15)

  result, err := payzen.SubscriptionUpdate(subscription)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v - %v", result.Message, result.Errors)
    return
  }

  if !result.EmptyResponseSuccess {
    t.Errorf("empty success response expected")
    return    
  }  


}

// go test -v github.com/mobilemindtec/go-payments/tests -run TestPayZenV4SubscriptionGet
func TestPayZenV4SubscriptionGet(t *testing.T) {

  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  subscriptionId, _ := client.Get("SubscriptionId").Result()
  token, _ := client.Get("Token").Result()
  result, err := payzen.SubscriptionGet(subscriptionId, token)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v - %v", result.Message, result.Errors)
    return
  }

  resp := result.GetResponse()

  if len(resp.SubscriptionId) == 0 {
    t.Errorf("SubscriptionId expected, but is empty")
    return    
  }  

  if resp.IsCancelled() {
    t.Errorf("cancelled is not expected, but is")
    return    
  }  
}

// go test -v github.com/mobilemindtec/go-payments/tests -run TestPayZenV4SubscriptionCancel
func TestPayZenV4SubscriptionCancel(t *testing.T) {

  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  subscriptionId, _ := client.Get("SubscriptionId").Result()
  token, _ := client.Get("Token").Result()
  result, err := payzen.SubscriptionCancel(subscriptionId, token)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v - %v", result.Message, result.Errors)
    return
  } 

  if !result.EmptyResponseSuccess {
    t.Errorf("empty success response expected")
    return    
  } 

}

// go test -v github.com/mobilemindtec/go-payments/tests -run TestPayZenV4GetTransaction
func TestPayZenV4GetTransaction(t *testing.T) {

  payzen := v4.NewPayZen("pt-BR", ApiMode, Authentication)
  payzen.SetDebug()

  uuid, _ := client.Get("TransactionUuid").Result()
  result, err := payzen.TransactionGet(uuid)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v - %v", result.Message, result.Errors)
    return
  }
}

// go test -v github.com/mobilemindtec/go-payments/tests -run TestPayZenV4NotificacaoFormData
func TestPayZenV4NotificacaoFormData(t *testing.T) {
  urlQuery := []byte(`kr-hash-key=password&kr-hash-algorithm=sha256_hmac&kr-answer=%7B%22shopId%22%3A%2231187067%22%2C%22orderCycle%22%3A%22CLOSED%22%2C%22orderStatus%22%3A%22PAID%22%2C%22serverDate%22%3A%222021-08-03T21%3A34%3A24%2B00%3A00%22%2C%22orderDetails%22%3A%7B%22orderTotalAmount%22%3A1500%2C%22orderEffectiveAmount%22%3A1500%2C%22orderCurrency%22%3A%22BRL%22%2C%22mode%22%3A%22TEST%22%2C%22orderId%22%3A%22a218c526-222f-4ab3-b055-248669e30b34%22%2C%22_type%22%3A%22V4%2FOrderDetails%22%7D%2C%22customer%22%3A%7B%22billingDetails%22%3A%7B%22address%22%3A%22Rua+Vit%C3%B3ria%22%2C%22category%22%3A%22PRIVATE%22%2C%22cellPhoneNumber%22%3A%2254999767081%22%2C%22city%22%3A%22Bento+Goncalves%22%2C%22country%22%3A%22BR%22%2C%22district%22%3A%22Botafogo%22%2C%22firstName%22%3A%22Ricardo%22%2C%22identityCode%22%3A%2283361855004%22%2C%22language%22%3A%22PT%22%2C%22lastName%22%3A%22Bocchi%22%2C%22phoneNumber%22%3A%225430553222%22%2C%22state%22%3A%22RS%22%2C%22streetNumber%22%3A%22255%22%2C%22title%22%3Anull%2C%22zipCode%22%3A%2295700540%22%2C%22legalName%22%3Anull%2C%22_type%22%3A%22V4%2FCustomer%2FBillingDetails%22%7D%2C%22email%22%3A%22ricardo%40mobilemind.com.br%22%2C%22reference%22%3Anull%2C%22shippingDetails%22%3A%7B%22address%22%3Anull%2C%22address2%22%3Anull%2C%22category%22%3Anull%2C%22city%22%3Anull%2C%22country%22%3Anull%2C%22deliveryCompanyName%22%3Anull%2C%22district%22%3Anull%2C%22firstName%22%3Anull%2C%22identityCode%22%3Anull%2C%22lastName%22%3Anull%2C%22legalName%22%3Anull%2C%22phoneNumber%22%3Anull%2C%22shippingMethod%22%3Anull%2C%22shippingSpeed%22%3Anull%2C%22state%22%3Anull%2C%22streetNumber%22%3Anull%2C%22zipCode%22%3Anull%2C%22_type%22%3A%22V4%2FCustomer%2FShippingDetails%22%7D%2C%22extraDetails%22%3A%7B%22browserAccept%22%3Anull%2C%22fingerPrintId%22%3Anull%2C%22ipAddress%22%3A%22138.36.81.242%22%2C%22browserUserAgent%22%3A%22Go-http-client%2F1.1%22%2C%22_type%22%3A%22V4%2FCustomer%2FExtraDetails%22%7D%2C%22shoppingCart%22%3A%7B%22insuranceAmount%22%3Anull%2C%22shippingAmount%22%3Anull%2C%22taxAmount%22%3Anull%2C%22cartItemInfo%22%3Anull%2C%22_type%22%3A%22V4%2FCustomer%2FShoppingCart%22%7D%2C%22_type%22%3A%22V4%2FCustomer%2FCustomer%22%7D%2C%22transactions%22%3A%5B%7B%22shopId%22%3A%2231187067%22%2C%22uuid%22%3A%2237fe11c80f0646b6911f72ffe16f60e2%22%2C%22amount%22%3A1500%2C%22currency%22%3A%22BRL%22%2C%22paymentMethodType%22%3A%22CARD%22%2C%22paymentMethodToken%22%3Anull%2C%22status%22%3A%22PAID%22%2C%22detailedStatus%22%3A%22AUTHORISED%22%2C%22operationType%22%3A%22DEBIT%22%2C%22effectiveStrongAuthentication%22%3A%22DISABLED%22%2C%22creationDate%22%3A%222021-08-03T18%3A25%3A56%2B00%3A00%22%2C%22errorCode%22%3Anull%2C%22errorMessage%22%3Anull%2C%22detailedErrorCode%22%3Anull%2C%22detailedErrorMessage%22%3Anull%2C%22metadata%22%3Anull%2C%22transactionDetails%22%3A%7B%22liabilityShift%22%3A%22NO%22%2C%22effectiveAmount%22%3A1500%2C%22effectiveCurrency%22%3A%22BRL%22%2C%22creationContext%22%3A%22CHARGE%22%2C%22cardDetails%22%3A%7B%22paymentSource%22%3A%22EC%22%2C%22manualValidation%22%3A%22NO%22%2C%22expectedCaptureDate%22%3A%222021-08-03T18%3A25%3A56%2B00%3A00%22%2C%22effectiveBrand%22%3A%22VISA%22%2C%22pan%22%3A%22497010XXXXXX0007%22%2C%22expiryMonth%22%3A12%2C%22expiryYear%22%3A2025%2C%22country%22%3A%22BR%22%2C%22issuerCode%22%3Anull%2C%22issuerName%22%3Anull%2C%22effectiveProductCode%22%3A%22F%22%2C%22legacyTransId%22%3A%22943225%22%2C%22legacyTransDate%22%3A%222021-08-03T18%3A25%3A56%2B00%3A00%22%2C%22paymentMethodSource%22%3A%22NEW%22%2C%22authorizationResponse%22%3A%7B%22amount%22%3A1500%2C%22currency%22%3A%22BRL%22%2C%22authorizationDate%22%3A%222021-08-03T18%3A25%3A56%2B00%3A00%22%2C%22authorizationNumber%22%3A%22009894%22%2C%22authorizationResult%22%3A%220%22%2C%22authorizationMode%22%3A%22FULL%22%2C%22_type%22%3A%22V4%2FPaymentMethod%2FDetails%2FCards%2FCardAuthorizationResponse%22%7D%2C%22captureResponse%22%3A%7B%22refundAmount%22%3Anull%2C%22refundCurrency%22%3Anull%2C%22captureDate%22%3Anull%2C%22captureFileNumber%22%3Anull%2C%22effectiveRefundAmount%22%3Anull%2C%22effectiveRefundCurrency%22%3Anull%2C%22_type%22%3A%22V4%2FPaymentMethod%2FDetails%2FCards%2FCardCaptureResponse%22%7D%2C%22threeDSResponse%22%3A%7B%22authenticationResultData%22%3A%7B%22transactionCondition%22%3Anull%2C%22enrolled%22%3Anull%2C%22status%22%3Anull%2C%22eci%22%3Anull%2C%22xid%22%3Anull%2C%22cavvAlgorithm%22%3Anull%2C%22cavv%22%3Anull%2C%22signValid%22%3Anull%2C%22brand%22%3Anull%2C%22_type%22%3A%22V4%2FPaymentMethod%2FDetails%2FCards%2FCardAuthenticationResponse%22%7D%2C%22_type%22%3A%22V4%2FPaymentMethod%2FDetails%2FCards%2FThreeDSResponse%22%7D%2C%22authenticationResponse%22%3Anull%2C%22installmentNumber%22%3A1%2C%22installmentCode%22%3A%221%22%2C%22markAuthorizationResponse%22%3A%7B%22amount%22%3Anull%2C%22currency%22%3Anull%2C%22authorizationDate%22%3Anull%2C%22authorizationNumber%22%3Anull%2C%22authorizationResult%22%3Anull%2C%22_type%22%3A%22V4%2FPaymentMethod%2FDetails%2FCards%2FMarkAuthorizationResponse%22%7D%2C%22cardHolderName%22%3Anull%2C%22identityDocumentNumber%22%3Anull%2C%22identityDocumentType%22%3Anull%2C%22_type%22%3A%22V4%2FPaymentMethod%2FDetails%2FCardDetails%22%7D%2C%22acquirerDetails%22%3Anull%2C%22fraudManagement%22%3A%7B%22riskControl%22%3A%5B%5D%2C%22riskAnalysis%22%3A%5B%5D%2C%22riskAssessments%22%3Anull%2C%22_type%22%3A%22V4%2FPaymentMethod%2FDetails%2FFraudManagement%22%7D%2C%22subscriptionDetails%22%3A%7B%22subscriptionId%22%3Anull%2C%22_type%22%3A%22V4%2FPaymentMethod%2FDetails%2FSubscriptionDetails%22%7D%2C%22parentTransactionUuid%22%3Anull%2C%22mid%22%3A%2280151051%22%2C%22sequenceNumber%22%3A1%2C%22taxAmount%22%3Anull%2C%22preTaxAmount%22%3Anull%2C%22taxRate%22%3Anull%2C%22externalTransactionId%22%3A%2201709790205590289818%22%2C%22dcc%22%3Anull%2C%22nsu%22%3A%22015736736298%22%2C%22tid%22%3Anull%2C%22acquirerNetwork%22%3A%22REDE%22%2C%22taxRefundAmount%22%3Anull%2C%22userInfo%22%3A%22API+REST%22%2C%22paymentMethodTokenPreviouslyRegistered%22%3Anull%2C%22occurrenceType%22%3A%22UNITAIRE%22%2C%22_type%22%3A%22V4%2FTransactionDetails%22%7D%2C%22_type%22%3A%22V4%2FPaymentTransaction%22%7D%5D%2C%22subMerchantDetails%22%3Anull%2C%22_type%22%3A%22V4%2FPayment%22%7D&kr-answer-type=V4%2FPayment&kr-hash=b40b62ec15baae5fa57ce914704f46c1edfa766cb1dbf0f9a5a2b766156c5c63`)


  webhook := v4.NewDefaultWebhook()
  data, err := webhook.Parse(urlQuery)

  if err != nil {
    t.Errorf("Parse error: %v", err)  
    return
  }

  signeture := v4.GenerateSignatureFromBody(Authentication.PasswordTest, data.KrAnswer)

  if signeture != data.KrHash {
    t.Errorf("Signature calculed = %v,  received = %v", signeture, data.KrHash)  
  }

  //t.Errorf("%v", data)

  /*
  formData := make(map[string]interface{})
  
  webhookData := v4.NewWebhookData()
  urlQuery, _ = url.QueryUnescape(urlQuery)
  t.Errorf(urlQuery)
  splited := strings.Split(urlQuery, "&")

  krAnswer := ""
  krHash := ""

  for _, value := range splited {
    vals := strings.Split(value, "=")

    if vals[0] == "kr-hash" {
      krHash = vals[1]
    }

    if vals[0] == "kr-answer" {
      krAnswer = vals[1]
      //answer := make(map[string]interface{})
      json.Unmarshal([]byte(vals[1]), webhookData.Answer)    
      formData[vals[0]] = webhookData.Answer
    } else{


      formData[vals[0]]  = vals[1]
    }
  }

  signeture := v4.GenerateSignatureFromBody(Authentication.PasswordTest, krAnswer)

  //if signeture != krHash {
    t.Errorf("Signature calculed = %v,  received = %v", signeture, krHash)  
  //}

  jsonData, _ := json.MarshalIndent(formData, "", " ")
  t.Errorf(string(jsonData))
  */
}