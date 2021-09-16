package gopayments

import (
  "github.com/mobilemindtec/go-payments/asaas"
  "github.com/mobilemindtec/go-payments/api"
	"testing"
  "time"
	"fmt"
)

var (  
  AsaasApiMode = api.AsaasModeTest
)


func fillAsaasCard(payment *asaas.Payment) {
  payment.Card.HolderName = "Ricardo Bocchi"
  payment.Card.Number = "4916561358240742" //"4916561358240741" cartão de erro
  payment.Card.ExpiryMonth = "12"
  payment.Card.ExpiryYear = "2025"
  payment.Card.SecurityCode = "123"

  payment.CardHolderInfo.Name = "Ricardo Bocchi"
  payment.CardHolderInfo.Email = "ricardo@mobilemind.com.br"
  payment.CardHolderInfo.CpfCnpj = "83361855004"
  payment.CardHolderInfo.PostalCode = "95700540"
  payment.CardHolderInfo.AddressNumber = "255"
  payment.CardHolderInfo.Phone = "54999767081"

  payment.RemoteIp = GetHostIp()    
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCustomerCreate
func TestAsaasCustomerCreate(t *testing.T) {

  customer := new(asaas.Customer)
  customer.Name = "Ricardo Bocchi"
  customer.Email = "ricardobocchi@mobilemind.com.br"
  customer.CpfCnpj = "83361855004"
  customer.MobilePhone = "54999767081"
  customer.ExternalReference = "12345"
  customer.NotificationDisabled = true

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  result, err := Asaas.CustomerCreate(customer)

  if err != nil {
    t.Errorf("Erro ao criar Customer: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar Customer: %v", result.Message)
    return
  }

  if len(result.Id) == 0 {
    t.Errorf("Customer criado sem ID")
    return    
  }

  client.Set("ClientId", result.Id, 0)

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCustomerFind
func TestAsaasCustomerFind(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("ClientId").Result()

  result, err := Asaas.CustomerFindByKey("externalReference", "12345")

  if err != nil {
    t.Errorf("Erro ao criar Customer: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar Customer: %v", result.Message)
    return
  }

  if result.CustomerResults.TotalCount != 1 {
    t.Errorf(fmt.Sprintf("Customer count expected %v  returned %v", 1, result.CustomerResults.TotalCount))
    return        
  }

  customer := result.CustomerResults.First()

  

  if customer.Id != id {
    t.Errorf(fmt.Sprintf("Customer id expected %v  returned %v", id, customer.Id))
    return    
  }

  if customer.Name != "Ricardos Bocchi" {
    t.Errorf(fmt.Sprintf("Customer name expected %v  returned %v", "Ricardos Bocchi", customer.Name))
    return    
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCustomerGet
func TestAsaasCustomerGet(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("ClientId").Result()

  result, err := Asaas.CustomerGet(id)

  if err != nil {
    t.Errorf("Erro ao criar Customer: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar Customer: %v", result.Message)
    return
  }

  if !result.CustomerResults.HasData() {
    t.Errorf(fmt.Sprintf("Customer expected, but is null"))
    return        
  }

  customer := result.CustomerResults.First()
  

  if customer.Id != id {
    t.Errorf(fmt.Sprintf("Customer id expected %v  returned %v", id, customer.Id))
    return    
  }

  if customer.Name != "Ricardos Bocchi" {
    t.Errorf(fmt.Sprintf("Customer name expected %v  returned %v", "Ricardos Bocchi", customer.Name))
    return    
  }
}


// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCustomerUpdate
func TestAsaasCustomerUpdate(t *testing.T) {

  id, _ := client.Get("ClientId").Result()
  customer := new(asaas.Customer)
  customer.Id = id
  customer.Name = "Ricardos Bocchi"
  customer.Email = "ricardobocchi@mobilemind.com.br"
  customer.CpfCnpj = "83361855004"
  customer.MobilePhone = "54999767081"
  customer.ExternalReference = "12345"
  customer.NotificationDisabled = true

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  result, err := Asaas.CustomerUpdate(customer)

  if err != nil {
    t.Errorf("Erro ao criar Customer: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar Customer: %v", result.Message)
    return
  }

  if len(result.Id) == 0 {
    t.Errorf("Customer criado sem ID")
    return    
  }

  client.Set("ClientId", result.Id, 0)

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCreatePaymentBoleto
func TestAsaasCreatePaymentBoleto(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  customerId, _ := client.Get("ClientId").Result()
  orderId := GenUUID()

  payment := asaas.NewPaymentWithBoleto(customerId, orderId, time.Now().AddDate(0, 0, 3), 10)

  result, err := Asaas.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar payment: %v", result.Message)
    return
  }

  if result.Status != api.AsaasPending {
    t.Errorf("Status expected: %v, Received %v", api.AsaasPending, result.Status)
    return
  }


  client.Set("OrderId", orderId, 0)
  client.Set("PaymentId", result.Id, 0)

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCreatePaymentPix
func TestAsaasCreatePaymentPix(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  customerId, _ := client.Get("ClientId").Result()
  orderId := GenUUID()

  payment := asaas.NewPaymentWithPix(customerId, orderId, time.Now().AddDate(0, 0, 3), 10)

  result, err := Asaas.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar payment: %v", result.Message)
    return
  }

  if result.Status != api.AsaasPending {
    t.Errorf("Status expected: %v, Received %v", api.AsaasPending, result.Status)
    return
  }


  client.Set("OrderId", orderId, 0)
  client.Set("PaymentId", result.Id, 0)

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCreatePaymentCreditCard
func TestAsaasCreatePaymentCreditCard(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  customerId, _ := client.Get("ClientId").Result()
  orderId := GenUUID()

  payment := asaas.NewPaymentWithCard(customerId, orderId, 10)
  
  fillAsaasCard(payment)

  result, err := Asaas.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar payment: %v, %v", result.Message, result.ErrorsToMap())
    return
  }

  if result.Status != api.AsaasConfirmed {
    t.Errorf("Status expected: %v, Received %v", api.AsaasConfirmed, result.Status)
    return
  }


  client.Set("OrderId", orderId, 0)
  client.Set("PaymentId", result.Id, 0)

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCreatePaymentCreditCardParcelado
func TestAsaasCreatePaymentCreditCardParcelado(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  customerId, _ := client.Get("ClientId").Result()
  orderId := GenUUID()

  payment := asaas.NewPaymenInstallmenttWithCard(customerId, orderId, 500, 10)
  
  fillAsaasCard(payment)

  result, err := Asaas.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar payment: %v, %v", result.Message, result.ErrorsToMap())
    return
  }

  if result.Status != api.AsaasConfirmed {
    t.Errorf("Status expected: %v, Received %v", api.AsaasConfirmed, result.Status)
    return
  }


  client.Set("OrderId", orderId, 0)
  client.Set("PaymentId", result.Id, 0)
  client.Set("InstallmentId", result.Installment, 0)

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCreatePaymentBoletoParcelado
func TestAsaasCreatePaymentBoletoParcelado(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  customerId, _ := client.Get("ClientId").Result()
  orderId := GenUUID()

  payment := asaas.NewPaymenInstallmenttWithBoleto(customerId, orderId, 500, 10)
  
  result, err := Asaas.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar payment: %v, %v", result.Message, result.ErrorsToMap())
    return
  }

  if result.Status != api.AsaasPending {
    t.Errorf("Status expected: %v, Received %v", api.AsaasConfirmed, result.Status)
    return
  }

  client.Set("OrderId", orderId, 0)
  client.Set("PaymentId", result.Id, 0)
  client.Set("InstallmentId", result.Installment, 0)
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasPaymentCancel
func TestAsaasPaymentCancel(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("PaymentId").Result()

  result, err := Asaas.PaymentCancel(id)

  if err != nil {
    t.Errorf("Erro ao canceler payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao canceler payment: %v", result.Message)
    return
  }
  
  if !result.Deleted {
    t.Errorf("Delete true expected, but is not")
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasPaymentRefund
func TestAsaasPaymentRefund(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("PaymentId").Result()

  result, err := Asaas.PaymentRefund(id)

  if err != nil {
    t.Errorf("Erro ao devolver payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao devolver payment: %v", result.Message)
    return
  }
  
  if result.Status != api.AsaasRefunded {
    t.Errorf("Status expected: %v, Received %v", api.AsaasRefunded, result.Status)
    return
  }  
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasPaymentGet
func TestAsaasPaymentGet(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("PaymentId").Result()

  result, err := Asaas.PaymentGet(id)

  if err != nil {
    t.Errorf("Erro ao buscar payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar payment: %v", result.Message)
    return
  }
  
  if result.Status != api.AsaasConfirmed {
    t.Errorf("Status expected: %v, Received %v", api.AsaasConfirmed, result.Status)
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasPaymentGetPixQrCode
func TestAsaasPaymentGetPixQrCode(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("PaymentId").Result()

  result, err := Asaas.PaymentGetPixQrCode(id)

  if err != nil {
    t.Errorf("Erro ao buscar payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar payment: %v", result.Message)
    return
  }
  
  if len(result.EncodedImage) > 0{
    t.Errorf("EncodedImage expected, but is empty")
    return
  }

  if result.Status != api.AsaasReceivedInCash {
    t.Errorf("Status expected: %v, Received %v", api.AsaasReceivedInCash, result.Status)
    return
  }

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasPaymentReceiveInCash
func TestAsaasPaymentReceiveInCash(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("PaymentId").Result()

  payment := asaas.NewPaymentInCash(id, time.Now(), 10)

  result, err := Asaas.PaymentReceiveInCash(payment)

  if err != nil {
    t.Errorf("Erro ao buscar payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar payment: %v", result.Message)
    return
  }

}


// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasPaymentFind
func TestAsaasPaymentFind(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("ClientId").Result()

  result, err := Asaas.PaymentFindByKey("customer", id)

  if err != nil {
    t.Errorf("Erro ao pesquisar payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao pesquisar payment: %v, %v", result.Message, result.ErrorsToMap())
    return
  }

  if !result.PaymentResults.HasData() {
    t.Errorf("Expected payments, but not have")
    return    
  }

  first := result.PaymentResults.First()
  
  if first.Status != api.AsaasConfirmed {
    t.Errorf("Status expected: %v, Received %v", api.AsaasConfirmed, first.Status)
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasPaymentGetInstalments
func TestAsaasPaymentGetInstalments(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("InstallmentId").Result()

  result, err := Asaas.InstallmentsGet(id)

  if err != nil {
    t.Errorf("Erro ao buscar payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar payment: %v", result.Message)
    return
  }
  
  if !result.PaymentResults.HasData() {
    t.Errorf("payments expected, but not have")
    return
  }

  if result.PaymentResults.Count() != 10 {
    t.Errorf("Installments count Expected %v, Received %v", 10, result.PaymentResults.Count())
    return
  }

  if result.PaymentResults.First().Status != api.AsaasConfirmed {
    t.Errorf("First Status Expected %v, Received %v", api.AsaasConfirmed, result.PaymentResults.First().Status)
    return
  }

  if result.PaymentResults.Last().Status != api.AsaasConfirmed {
    t.Errorf("Last Status Expected %v, Received %v", api.AsaasConfirmed, result.PaymentResults.First().Status)
    return
  }

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasInstallmentCancel
func TestAsaasInstallmentCancel(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("InstallmentId").Result()

  result, err := Asaas.InstallmentCancel(id)

  if err != nil {
    t.Errorf("Erro ao canceler parcelamento: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao canceler parcelamento: %v", result.Message)
    return
  }
  
  if !result.Deleted {
    t.Errorf("Delete true expected, but is not")
    return
  }  
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasInstallmentRefund
func TestAsaasInstallmentRefund(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("InstallmentId").Result()

  result, err := Asaas.InstallmentRefund(id)

  if err != nil {
    t.Errorf("Erro ao devolver payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao devolver payment: %v", result.Message)
    return
  }
  
  if result.Status != api.AsaasRefunded {
    t.Errorf("Status expected: %v, Received %v", api.AsaasRefunded, result.Status)
    return
  }  
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaaTokenCreate
func TestAsaaTokenCreate(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  customerId, _ := client.Get("ClientId").Result()

  payment := asaas.NewSubscriptionWithCardToken("")
  fillAsaasCard(payment)

  tokenRequest := asaas.NewTokenRequest(customerId)
  tokenRequest.CreditCardCcv = payment.Card.SecurityCode
  tokenRequest.CreditCardHolderName = payment.Card.HolderName
  tokenRequest.CreditCardExpiryMonth = payment.Card.ExpiryMonth
  tokenRequest.CreditCardNumber = payment.Card.Number
  tokenRequest.CreditCardExpiryYear = payment.Card.ExpiryYear

  result, err := Asaas.TokenCreate(tokenRequest)

  if err != nil {
    t.Errorf("Erro ao update subscription token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao update subscription token: %v, %v", result.Message, result.ErrorsToMap())
    return
  }

  if result.Status != api.AsaasSuccess {
    t.Errorf("Status expected: %v, Received %v", api.AsaasSuccess, result.Status)
    return
  }

  if len(result.Card.Token) == 0 {
    t.Errorf("token expected, but not have")
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCreateSubscriptionBoleto
func TestAsaasCreateSubscriptionBoleto(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  customerId, _ := client.Get("ClientId").Result()
  orderId := GenUUID()

  payment := asaas.NewSubscriptionWithBoleto(customerId, orderId, api.Monthly, time.Now(), 100)
  //payment.SubscriptionCycle
  //payment.NextDueDate
  //payment.EndDate
  payment.MaxPayments = 12
  
  result, err := Asaas.SubscriptionCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar subscription: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar subscription: %v, %v", result.Message, result.ErrorsToMap())
    return
  }

  if result.Status != api.AsaasActive {
    t.Errorf("Status expected: %v, Received %v", api.AsaasActive, result.Status)
    return
  }


  client.Set("OrderId", orderId, 0)
  client.Set("PaymentId", result.Id, 0)
  client.Set("SubscriptionId", result.Id, 0)

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCreateSubscriptionCard
func TestAsaasCreateSubscriptionCard(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  customerId, _ := client.Get("ClientId").Result()
  orderId := GenUUID()

  payment := asaas.NewSubscriptionWithCard(customerId, orderId, api.Monthly, time.Now(), 100)
  //payment.SubscriptionCycle
  //payment.NextDueDate
  //payment.EndDate
  payment.MaxPayments = 12
  fillAsaasCard(payment)
  
  result, err := Asaas.SubscriptionCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar subscription: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar subscription: %v, %v", result.Message, result.ErrorsToMap())
    return
  }

  if result.Status != api.AsaasActive {
    t.Errorf("Status expected: %v, Received %v", api.AsaasActive, result.Status)
    return
  }


  client.Set("OrderId", orderId, 0)
  client.Set("PaymentId", result.Id, 0)
  client.Set("SubscriptionId", result.Id, 0)

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaaSubscriptionUpdate
func TestAsaaSubscriptionUpdate(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  customerId, _ := client.Get("ClientId").Result()
  orderId, _ := client.Get("OrderId").Result()
  id, _ := client.Get("SubscriptionId").Result()

  payment := asaas.NewSubscriptionWithCard(customerId, orderId, api.Monthly, time.Now(), 110)
  payment.Id = id
  //payment.SubscriptionCycle
  //payment.NextDueDate
  //payment.EndDate
  payment.MaxPayments = 12
  fillAsaasCard(payment)
  
  result, err := Asaas.SubscriptionUpdate(payment)

  if err != nil {
    t.Errorf("Erro ao criar subscription: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar subscription: %v, %v", result.Message, result.ErrorsToMap())
    return
  }

  if result.Status != api.AsaasActive {
    t.Errorf("Status expected: %v, Received %v", api.AsaasActive, result.Status)
    return
  }

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaaSubscriptionUpdateCardToken
func TestAsaaSubscriptionUpdateCardToken(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  //paymentId, _ := client.Get("PaymentId").Result()

  payment := asaas.NewSubscriptionWithCardToken("pay_7551240616502368")
  fillAsaasCard(payment)
  
  result, err := Asaas.SubscriptionUpdateCardToken(payment)

  if err != nil {
    t.Errorf("Erro ao update subscription token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao update subscription token: %v, %v", result.Message, result.ErrorsToMap())
    return
  }

  if result.Status != api.AsaasActive {
    t.Errorf("Status expected: %v, Received %v", api.AsaasActive, result.Status)
    return
  }
}


// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasSubscriptionGet
func TestAsaasSubscriptionGet(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("SubscriptionId").Result()

  result, err := Asaas.SubscriptionGet(id)

  if err != nil {
    t.Errorf("Erro ao buscar subscription: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar subscription: %v", result.Message)
    return
  }
  
  if result.Status != api.AsaasActive {
    t.Errorf("Status expected: %v, Received %v", api.AsaasActive, result.Status)
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasSubscriptionPaymentsGet
func TestAsaasSubscriptionPaymentsGet(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("SubscriptionId").Result()

  result, err := Asaas.SubscriptionPaymentsGet(id)

  if err != nil {
    t.Errorf("Erro ao buscar subscription: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar subscription: %v", result.Message)
    return
  }
  
  if !result.PaymentResults.HasData() {
    t.Errorf("payments expected, but not have")
    return
  }

  t.Errorf("Count = %v", result.PaymentResults.Count())

  if result.PaymentResults.First().Status != api.AsaasPending {
    t.Errorf("Status expected: %v, Received %v", api.AsaasPending, result.PaymentResults.First().Status)
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasSubscriptionCancel
func TestAsaasSubscriptionCancel(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("SubscriptionId").Result()

  result, err := Asaas.SubscriptionCancel(id)

  if err != nil {
    t.Errorf("Erro ao canceler payment: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao canceler payment: %v", result.Message)
    return
  }
  
  if !result.Deleted {
    t.Errorf("Delete true expected, but is not")
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasPaymentLinkCreate
func TestAsaasPaymentLinkCreate(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  payment := asaas.NewPaymentLink(100, asaas.Detached, 10)
  payment.Name = "Payment test"
  //payment.DueDateLimitDays
  //payment.MaxInstallmentCount
  payment.SetEndDate(time.Now().AddDate(0, 0, 5))
  
  result, err := Asaas.PaymentLinkCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar payment link: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar payment link: %v, %v", result.Message, result.ErrorsToMap())
    return
  }

  if !result.Active && !result.Deleted {
    t.Errorf("active expected, but is not")
    return
  }

  client.Set("PaymentLinkId", result.Id, 0)
}


// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasPaymentLinkCancel
func TestAsaasPaymentLinkCancel(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("PaymentLinkId").Result()

  result, err := Asaas.PaymentLinkCancel(id)

  if err != nil {
    t.Errorf("Erro ao canceler payment link: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao canceler payment link: %v", result.Message)
    return
  }
  
  if !result.Deleted {
    t.Errorf("Delete true expected, but is not")
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasPaymentLinkGet
func TestAsaasPaymentLinkGet(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  id, _ := client.Get("PaymentLinkId").Result()

  result, err := Asaas.PaymentLinkGet(id)

  if err != nil {
    t.Errorf("Erro ao buscar payment link: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar payment link: %v", result.Message)
    return
  }
  

  if result.Status != api.AsaasPending {
    t.Errorf("Status expected: %v, Received %v", api.AsaasPending, result.Status)
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasFinancialTransactionsList
func TestAsaasFinancialTransactionsList(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true


  filter := asaas.NewDefaultFilter()
  filter.SetStartDate(time.Now().AddDate(0, 0, -10))
  filter.SetFinishDate(time.Now())
  filter.Limit = 25

  result, err := Asaas.FinancialTransactionsList(filter)

  if err != nil {
    t.Errorf("Erro ao buscar financial transactions: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar financial transactions: %v", result.Message)
    return
  }
  

  if !result.FinancialTransactionResults.HasData() {
    t.Errorf("data is expected, but not have")
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCurrentBalance
func TestAsaasCurrentBalance(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true


  result, err := Asaas.CurrentBalance()

  if err != nil {
    t.Errorf("Erro ao buscar financial transactions: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar financial transactions: %v", result.Message)
    return
  }
  

  if result.TotalBalance == 0 {
    t.Errorf("TotalBalance is expected, but not have")
    return
  }
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasTransferCreate
func TestAsaasTransferCreate(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  bank := asaas.NewBank("104")
  bankAccount := asaas.NewBankAccount(bank, api.ContaCorrente)
  bankAccount.AccountName = "Mobile Mind - Caixa Corrente"
  bankAccount.OwnerName = "Mobile Mind Empresa de Tecnologia LTDA"
  //bankAccount.OwnerBirthDate = "1995-04-12"
  bankAccount.CpfCnpj = "15095430000101"
  bankAccount.Agency = "3060"
  bankAccount.Account = "1128"
  bankAccount.AccountDigit = "2"

  transfer := asaas.NewTransfer(bankAccount, 20)

  result, err := Asaas.TransferCreate(transfer)

  if err != nil {
    t.Errorf("Erro ao criar transaferencia: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar transaferencia: %v", result.Message)
    return
  }
  
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasTransferList
func TestAsaasTransferList(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true


  filter := asaas.NewDefaultFilter()
  filter.SetDateCreated(time.Now())

  result, err := Asaas.TransferList(filter)

  if err != nil {
    t.Errorf("Erro ao buscar transaferencias: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar transaferencias: %v", result.Message)
    return
  }
  
  if !result.TransferResults.HasData() {
    t.Errorf("data is expected, but not have")
    return
  }

  if !result.TransferResults.First().Authorized {
    t.Errorf("authorized is expected, but is not")
    return
  }

  if result.TransferResults.First().Status != api.TransferPending {
    t.Errorf("status pending is expected, but is not")
    return
  }
  
}



// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasWebhook
func TestAsaasWebhook(t *testing.T) {
  jsonData := []byte(`
    {
      "event": "PAYMENT_RECEIVED",
      "payment": {
        "object": "payment",
        "id": "pay_080225913252",
        "dateCreated": "2017-03-10",
        "customer": "cus_G7Dvo4iphUNk",
        "subscription": "sub_VXJBYgP2u0eO",
        "installment": "ins_000000001031",
        "paymentLink": "123517639363",
        "dueDate": "2017-03-15",
        "value": 100.00,
        "netValue": 94.51,
        "billingType": "CREDIT_CARD",
        "status": "RECEIVED",
        "description": "Pedido 056984",
        "externalReference": "056984",
        "confirmedDate": "2017-03-15",
        "originalValue": null,
        "interestValue": null,
        "originalDueDate": "2017-06-10",
        "paymentDate": null,
        "clientPaymentDate": null,
        "invoiceUrl": "https://www.asaas.com/i/080225913252",
        "bankSlipUrl": null,
        "invoiceNumber": "00005101",
        "deleted": false,
        "creditCard": {
          "creditCardNumber": "8829",
          "creditCardBrand": "MASTERCARD",
          "creditCardToken": "a75a1d98-c52d-4a6b-a413-71e00b193c99"
        }
      }
    }
  `)

  webhook := asaas.NewWebhook("pt-BR", "123", nil)

  result, err := webhook.Parse(jsonData)

  if err != nil {
    t.Errorf("Erro ao criar webhook: %v", err)
    return
  }

  if result.Response.Status != api.AsaasReceived {
    t.Errorf("Status expected: %v, Received %v", api.AsaasReceived, result.Response.Status)
    return
  }

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasAccountCreate
func TestAsaasAccountCreate(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  bank := asaas.NewBank("104")
  bankAccount := asaas.NewBankAccount(bank, api.ContaCorrente)
  bankAccount.AccountName = "Mobile Mind - Caixa Corrente"
  bankAccount.OwnerName = "Mobile Mind Empresa de Tecnologia LTDA"
  //bankAccount.OwnerBirthDate = "1995-04-12"
  bankAccount.CpfCnpj = "15095430000101"
  bankAccount.Agency = "3060"
  bankAccount.Account = "1128"
  bankAccount.AccountDigit = "2"

  account := asaas.NewAccount(bankAccount)
  account.Name = "Rede Inova ST"
  account.Email = "rmedrogaria5@hotmail.com"
  account.LoginEmail = "rmedrogaria5@hotmail.com"
  account.CpfCnpj = "34735958000142"
  account.CompanyType = asaas.LIMITED
  account.Phone = "27999874613"
  account.MobilePhone = "27999874613"
  account.Address = "Angelo Pretti"
  account.AddressNumber = "12"
  account.Complement = ""
  account.Province = "Centro"
  account.PostalCode = "29650000"

  result, err := Asaas.AccountCreate(account)

  if err != nil {
    t.Errorf("Erro ao criar conta: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar conta: %v", result.Message)
    return
  }

  if !result.AccountResults.HasData() {
    t.Errorf("data is expected")
    return
  }  

  if len(result.AccountResults.First().WalletId) == 0 {
    t.Errorf("WalletId não pode ser vazio")
    return
  }  

  if len(result.AccountResults.First().ApiKey) == 0 {
    t.Errorf("ApiKey não pode ser vazio")
    return
  }  

}


// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasAccountList
func TestAsaasAccountList(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  result, err := Asaas.AccountList()

  if err != nil {
    t.Errorf("Erro ao buscar transaferencias: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar transaferencias: %v", result.Message)
    return
  }
  
  if !result.AccountResults.HasData() {
    t.Errorf("AccountResults is expected, but not have")
    return
  }

  if len(result.AccountResults.First().WalletId) == 0 {
    t.Errorf("WalletId não pode ser vazio")
    return
  }   
  
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasCreateOrChangeWebhook
func TestAsaasCreateOrChangeWebhook(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  webhookData := asaas.NewWebhookObject()

  webhookData.Url = "https://pay.mobilemind.com.br/gateway/asaas/tenant-uuid"
  webhookData.Email = "ricardo@mobilemind.com.br"
  webhookData.Enabled = true
  webhookData.Interrupted = false
  webhookData.ApiVersion = 3
  webhookData.AuthToken = "5tLxsL6uoN"  

  result, err := Asaas.WebhookCreateOrChange(webhookData)

  if err != nil {
    t.Errorf("Erro ao buscar transaferencias: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar transaferencias: %v", result.Message)
    return
  }
  
  
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestAsaasWebhookStatus
func TestAsaasWebhookStatus(t *testing.T) {

  Asaas := asaas.NewAsaas("pt-BR", AsaasAccessToken, AsaasApiMode)
  Asaas.Debug = true

  result, err := Asaas.WebhookStatus()

  if err != nil {
    t.Errorf("Erro ao buscar transaferencias: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao buscar transaferencias: %v", result.Message)
    return
  }
  
  if len(result.Webhook.Url) == 0 {
    t.Errorf("webhook is required")
    return
  }
  
}