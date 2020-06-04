package gopayments

import (
  "github.com/mobilemindtec/go-payments/payzen"
  "github.com/mobilemindtec/go-utils/app/util"
	"github.com/satori/go.uuid"
	"github.com/go-redis/redis"
	"testing"
	"time"
	"fmt"
	"os"
)

const(
  ShopId = "33015842"
  Mode = "TEST"
  Cert = "6820023378838080"

  KeyCardTransactionId = "CardTransactionId"
  KeyCardTransactionUuid = "CardTransactionUuid"
  KeyCardOrderId = "CardOrderId"

  KeyBoletoTransactionId = "BoletoTransactionId"
  KeyBoletoTransactionUuid = "BoletoTransactionUuid"
  KeyBoletoOrderId = "BoletoOrderId"

  KeyCardTokenCancelled = "CardTokenCancelled"

  KeyCardTokenActive = "CardTokenActive"

  KeySubscriptionId = "SubscriptionId"
)

var (
	client *redis.Client
)


func fillCard(card *payzen.PayZenCard) {
  card.Number = "4970100000000007"
  card.Scheme = "VISA"
  card.ExpiryMonth = "12"
  card.ExpiryYear = "2020"
  card.CardSecurityCode = "235"
}

func fillCustomer(customer *payzen.PayZenCustomer) {
  customer.FirstName = "Tony"
  customer.LastName = "Montana"
  customer.PhoneNumber = "54999999999"
  customer.Email = "ricardobocchi@gmail.com"
  customer.StreetNumber = "255"
  customer.Address = "Rua Vitoria"
  customer.District = "Botafogo"
  customer.ZipCode = "95700540"
  customer.City = "Bento Goncalves"
  customer.State = "RS"
  customer.Country = payzen.CountryBR
  customer.IdentityCode = "833.618.550-04"
}

func genUUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func setup(){
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func shutdown() {

}

func TestMain(m *testing.M) {
  setup()
  code := m.Run()
  shutdown()
  os.Exit(code)
}

func TestPayZenPaymentCreateCartao(t *testing.T) {

  time.Sleep(3 * time.Second)

	PayZen := payzen.NewPayZen("pt-BR")
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)


  payment.OrderId = genUUID()
  payment.Installments = 1
  payment.Amount = 10.0


	fillCard(payment.Card)
  fillCustomer(payment.Customer)

  result, err := PayZen.PaymentCreate(payment)

  if err != nil {
  	t.Errorf("Erro ao criar autorização: %v", err)
  	return
  }

  if result.Error {
  	t.Errorf("Erro ao criar autorização: %v", result.Message)
  }else{

  	if len(result.TransactionId) == 0 {
  		t.Errorf("Erro ao criar autorização: %v", "TransactionId não informada")
  	} else if len(result.TransactionUuid) == 0 {
  		t.Errorf("Erro ao criar autorização: %v", "TransactionUuid não informada")
  	} else {
  		client.Set(KeyCardTransactionUuid, result.TransactionUuid, 0)
  		client.Set(KeyCardTransactionId, result.TransactionId, 0)
  		client.Set(KeyCardOrderId, payment.OrderId, 0)
  	}
	}

}

func TestPayZenPaymentCreateBoleto(t *testing.T) {

  time.Sleep(1 * time.Second)

	PayZen := payzen.NewPayZen("pt-BR")
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)


  payment.OrderId = genUUID()
  payment.Installments = 1
  payment.Amount = 10.0


	payment.Card.Scheme = payzen.SchemeBoleto
  fillCustomer(payment.Customer)

  result, err := PayZen.PaymentCreate(payment)

  if err != nil {
  	t.Errorf("Erro ao criar autorização: %v", err)
  	return
  }

  if result.Error {
  	t.Errorf("Erro ao criar autorização: %v", result.Message)
  }else{

  	if len(result.TransactionId) == 0 {
  		t.Errorf("Erro ao criar autorização: %v", "TransactionId não informada")
  	} else if len(result.TransactionUuid) == 0 {
  		t.Errorf("Erro ao criar autorização: %v", "TransactionUuid não informada")
  	} else {
  		client.Set(KeyBoletoTransactionUuid, result.TransactionUuid, 0)
  		client.Set(KeyBoletoTransactionId, result.TransactionId, 0)
  		client.Set(KeyBoletoOrderId, payment.OrderId, 0)
  	}

	}
}

func TestPayZenPaymentCreateBoletoOnline(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)

  PayZen.OnDebug()


  payment.OrderId = genUUID()
  payment.Installments = 1
  payment.Amount = 10.0

  //client.Set(KeyBoletoOrderId, payment.OrderId, 0)



  payment.VadsTransId = "000004" // deve ser um valor númerico de 6 digitos que não pode repetir no mesmo dia
  payment.Card.Scheme = payzen.SchemeBoleto
  payment.Card.BoletoOnline = payzen.BoletoOnlineItauBoleto
  payment.Card.BoletoOnlineDaysDalay = 3

  fillCustomer(payment.Customer)

  result, err := PayZen.PaymentCreateBoletoOnline(payment)

  if err != nil {
    t.Errorf("Erro ao criar autorização: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar autorização: %v", result.Message)
  }else{

    //t.Errorf("URL Boleto: %v", result.BoletoUrl)

    if len(result.TransactionId) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionId não informada")
    } else if len(result.TransactionUuid) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionUuid não informada")
    } else {
      client.Set(KeyBoletoTransactionUuid, result.TransactionUuid, 0)
      client.Set(KeyBoletoTransactionId, result.TransactionId, 0)
      client.Set(KeyBoletoOrderId, payment.OrderId, 0)
    }

  }
}

func TestPayZenCaptureBoleto(t *testing.T) {

  time.Sleep(1 * time.Second)

	PayZen := payzen.NewPayZen("pt-BR")
	//PayZen.OnDebug()
  //account := payzen.NewPayZenAccount(ShopId, Mode, Cert)

  BoletoTransactionUuid, err := client.Get(KeyBoletoTransactionUuid).Result()

  if err != nil {
  	t.Errorf("Erro ao recuperar boleto transaction uuid: %v", err)
  	return
  }

  fmt.Printf("BoletoTransactionUuid %v\n", BoletoTransactionUuid)

  capture := payzen.NewPayZenCapturePayment(ShopId, Mode, Cert)
  capture.TransactionUuids = "000001"//BoletoTransactionUuid
  result, err := PayZen.PaymentCapture(capture)

  if err != nil {
  	t.Errorf("Erro ao criar autorização: %v", err)
  }

  if result.Error {
  	t.Errorf("Erro ao criar autorização: %v", result.Message)
  }

}

func TestPayZenCaptureCartao(t *testing.T) {

  time.Sleep(1 * time.Second)

	PayZen := payzen.NewPayZen("pt-BR")
	//PayZen.OnDebug()
  //account := payzen.NewPayZenAccount(ShopId, Mode, Cert)

  CardTransactionUuid, err := client.Get(KeyCardTransactionUuid).Result()

  if err != nil {
  	t.Errorf("Erro ao recuperar card transaction uuid: %v", err)
  	return
  }

  fmt.Printf("CardTransactionUuid %v\n", CardTransactionUuid)

  capture := payzen.NewPayZenCapturePayment(ShopId, Mode, Cert)
  capture.TransactionUuids = CardTransactionUuid
  result, err := PayZen.PaymentCapture(capture)

  if err != nil {
  	t.Errorf("Erro ao criar autorização: %v", err)
  }

  if result.Error {
  	t.Errorf("Erro ao criar autorização: %v", result.Message)
  }

}

func TestPayZenCreateTokenActive(t *testing.T) {

  time.Sleep(1 * time.Second)

	PayZen := payzen.NewPayZen("pt-BR")
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)

	fillCard(payment.Card)
  fillCustomer(payment.Customer)

  result, err := PayZen.PaymentTokenCreate(payment)

  if err != nil {
  	t.Errorf("Erro ao criar token: %v", err)
  	return
  }

  if result.Error {
  	t.Errorf("Erro ao criar token: %v", result.Message)
  }else{

  	if len(result.TokenInfo.Token) == 0 {
  		t.Errorf("Erro ao criar token: %v", "Token não informado")
  	}else{
  		client.Set(KeyCardTokenActive, result.TokenInfo.Token, 0)
  	}
	}

}

func TestPayZenCreateTokenCancelled(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)

  fillCard(payment.Card)
  fillCustomer(payment.Customer)

  result, err := PayZen.PaymentTokenCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v", result.Message)
  }else{

    if len(result.TokenInfo.Token) == 0 {
      t.Errorf("Erro ao criar token: %v", "Token não informado")
    }else{
      client.Set(KeyCardTokenCancelled, result.TokenInfo.Token, 0)
    }
  }

}

func TestPayZenUpdateToken(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  //PayZen.OnDebug()
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)

  fillCard(payment.Card)
  fillCustomer(payment.Customer)

  CardToken, err := client.Get(KeyCardTokenActive).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card token: %v", err)
    return
  }


  payment.Card.Token = CardToken

  result, err := PayZen.PaymentTokenUpdate(payment)

  if err != nil {
    t.Errorf("Erro ao criar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar token: %v", result.Message)
  }
}

func TestPayZenCancelToken(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")

  CardToken, err := client.Get(KeyCardTokenCancelled).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card token: %v", err)
    return
  }


  paymentToken := payzen.NewPayZenPaymentToken(ShopId, Mode, Cert)
  paymentToken.Token = CardToken

  result, err := PayZen.PaymentTokenCancel(paymentToken)

  if err != nil {
    t.Errorf("Erro ao cancelar token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao cancelar token: %v", result.Message)
  }
}

func TestPayZenGetDetailsTokenCancelled(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")

  CardToken, err := client.Get(KeyCardTokenCancelled).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card token: %v", err)
    return
  }


  paymentToken := payzen.NewPayZenPaymentToken(ShopId, Mode, Cert)
  paymentToken.Token = CardToken

  result, err := PayZen.PaymentTokenGetDetails(paymentToken)

  if err != nil {
    t.Errorf("Erro ao recuperar informações do token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do token: %v", result.Message)
  } else {

    fmt.Printf("**********************************************")
    fmt.Printf("Token Card Number %v\n", result.TokenInfo.Number)
    fmt.Printf("Token Card Brand %v\n", result.TokenInfo.Brand)
    fmt.Printf("Token CreationDate %v\n", result.TokenInfo.CreationDate)
    fmt.Printf("Token CancellationDate %v\n", result.TokenInfo.CancellationDate)
    fmt.Printf("**********************************************")

    if len(result.TokenInfo.Number) == 0 || len(result.TokenInfo.Brand) == 0 {
      t.Errorf("Algumas informações do token não estão presentes")
      return
    }

    if !result.TokenInfo.Cancelled {
      t.Errorf("O token não está cancelado")
      return
    }

    if result.TokenInfo.Active {
      t.Errorf("O token está ativo")
      return
    }

  }
}

func TestPayZenGetDetailsTokenActive(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")

  CardToken, err := client.Get(KeyCardTokenActive).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card token: %v", err)
    return
  }


  paymentToken := payzen.NewPayZenPaymentToken(ShopId, Mode, Cert)
  paymentToken.Token = CardToken

  result, err := PayZen.PaymentTokenGetDetails(paymentToken)

  if err != nil {
    t.Errorf("Erro ao recuperar informações do token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do token: %v", result.Message)
  } else {

    fmt.Printf("**********************************************")
    fmt.Printf("Token Card Number %v\n", result.TokenInfo.Number)
    fmt.Printf("Token Card Brand %v\n", result.TokenInfo.Brand)
    fmt.Printf("Token CreationDate %v\n", result.TokenInfo.CreationDate)
    fmt.Printf("Token CancellationDate %v\n", result.TokenInfo.CancellationDate)
    fmt.Printf("**********************************************")

    if len(result.TokenInfo.Number) == 0 || len(result.TokenInfo.Brand) == 0 {
      t.Errorf("Algumas informações do token não estão presentes")
      return
    }

    if result.TokenInfo.Cancelled {
      t.Errorf("O token está cancelado")
      return
    }

    if !result.TokenInfo.Active {
      t.Errorf("O token não está ativo")
      return
    }

  }
}

func TestPayZenReactiveToken(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")

  CardToken, err := client.Get(KeyCardTokenCancelled).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card token: %v", err)
    return
  }


  // busca status para ver se está inativo
  paymentToken := payzen.NewPayZenPaymentToken(ShopId, Mode, Cert)
  paymentToken.Token = CardToken

  result, err := PayZen.PaymentTokenGetDetails(paymentToken)

  if !result.TokenInfo.Cancelled {
    t.Errorf("O token não está cancelado")
    return
  }

  // reativa token
  result, err = PayZen.PaymentTokenReactive(paymentToken)

  if err != nil {
     t.Errorf("Erro ao reativer token: %v", err)
     return
  }

  if result.Error {
    t.Errorf("Erro ao reativer token: %v", result.Message)
    return
  }

    // busca status para ver se está ativo
  paymentToken = payzen.NewPayZenPaymentToken(ShopId, Mode, Cert)
  paymentToken.Token = CardToken

  result, err = PayZen.PaymentTokenGetDetails(paymentToken)

  if !result.TokenInfo.Active {
    t.Errorf("O token não está cancelado")
    return
  }


}

func TestPayZenGetDetailsTokenNotFound(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  //PayZen.OnDebug()


  paymentToken := payzen.NewPayZenPaymentToken(ShopId, Mode, Cert)
  paymentToken.Token = "3123213233213"

  result, _ := PayZen.PaymentTokenGetDetails(paymentToken)

  if !result.TokenInfo.NotFound {
    t.Errorf("token found: %v", result.Message)
    return
  }
}

func TestPayZenFindPayment(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  PayZen.OnDebug()

  CardOrderId, err := client.Get(KeyCardOrderId).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card token: %v", err)
    return
  }

  CardTransactionUuid, err := client.Get(KeyCardTransactionUuid).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card TransactionUuid: %v", err)
    return
  }

  paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
  paymentFind.OrderId = CardOrderId //"ecf02704-f155-40b0-bee3-3477d752da9d" //CardOrderId

  result, err := PayZen.PaymentFind(paymentFind)

  if result.PaymentNotFound {
    t.Errorf("O pagamento não foi encontrado")
  }

  if err != nil {
    t.Errorf("Erro ao recuperar informações do token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do token: %v", result.Message)
    return
  }

  if len(result.Transactions) != 1 {
    t.Errorf("Transaction count expected: %v, returned: %v", 1, len(result.Transactions))
    return
  }

  transaction := result.Transactions[0]

  if transaction.TransactionUuid != CardTransactionUuid {
    t.Errorf("TransactionUuid expected: %v, returned: %v", CardTransactionUuid, transaction.TransactionUuid)
    return
  }

  if transaction.Amount != 10.0 {
    t.Errorf("Transaction amount expected: %v, returned: %v", 10.0, transaction.Amount)
    return
  }

}

func TestPayZenFindPaymentBoletoOnline(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  PayZen.OnDebug()

  BoletoOrderId, err := client.Get(KeyBoletoOrderId).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card token: %v", err)
    return
  }

  
  BoletoOrderId = "6d5cef2b-c27d-4905-af85-86e29ca3c0fb"
  
  paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
  paymentFind.OrderId = BoletoOrderId //"205e125f-3bac-4210-b80e-ebec704a5845" //CardOrderId

  result, err := PayZen.PaymentFind(paymentFind)

  if result.PaymentNotFound {
    t.Errorf("O pagamento não foi encontrado")
  }

  if err != nil {
    t.Errorf("Erro ao recuperar informações do token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do token: %v", result.Message)
    return
  }

  if len(result.Transactions) != 1 {
    t.Errorf("Transaction count expected: %v, returned: %v", 1, len(result.Transactions))
    return
  }

  transaction := result.Transactions[0]

  if transaction.Amount != 10.0 {
    t.Errorf("Transaction amount expected: %v, returned: %v", 10.0, transaction.Amount)
    return
  }

}

func TestPayZenFindPaymentNotFound(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")


  paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
  paymentFind.OrderId = "31323213312"

  result, _ := PayZen.PaymentFind(paymentFind)

  if !result.PaymentNotFound {
    t.Errorf("O pagamento foi encontrado")
  }

}

func TestPayZenGetDetailsPayment(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")

  //PayZen.OnDebug()

  CardTransactionUuid, err := client.Get(KeyCardTransactionUuid).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card TransactionUuid: %v", err)
    return
  }

  paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
  paymentFind.TransactionUuid = CardTransactionUuid

  result, err := PayZen.PaymentGetDetails(paymentFind)

  if result.PaymentNotFound {
    t.Errorf("O pagamento não foi encontrado")
  }

  if err != nil {
    t.Errorf("Erro ao recuperar informações do token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do token: %v", result.Message)
    return
  }

  if len(result.Transactions) != 1 {
    t.Errorf("Transaction count expected: %v, returned: %v", 1, len(result.Transactions))
    return
  }

  transaction := result.Transactions[0]

  if transaction.TransactionUuid != CardTransactionUuid {
    t.Errorf("TransactionUuid expected: %v, returned: %v", CardTransactionUuid, transaction.TransactionUuid)
    return
  }

  if transaction.Amount != 10.0 {
    t.Errorf("Transaction amount expected: %v, returned: %v", 10.0, transaction.Amount)
    return
  }

}

func TestPayZenGetDetailsPaymentWithNsu(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")

  //PayZen.OnDebug()

  CardTransactionUuid, err := client.Get(KeyCardTransactionUuid).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card TransactionUuid: %v", err)
    return
  }

  paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
  paymentFind.TransactionUuid = CardTransactionUuid

  result, err := PayZen.PaymentGetDetailsWithNsu(paymentFind)

  if result.PaymentNotFound {
    t.Errorf("O pagamento não foi encontrado")
  }

  if err != nil {
    t.Errorf("Erro ao recuperar informações do token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do token: %v", result.Message)
    return
  }

  if len(result.Transactions) != 1 {
    t.Errorf("Transaction count expected: %v, returned: %v", 1, len(result.Transactions))
    return
  }

  transaction := result.Transactions[0]

  if transaction.TransactionUuid != CardTransactionUuid {
    t.Errorf("TransactionUuid expected: %v, returned: %v", CardTransactionUuid, transaction.TransactionUuid)
    return
  }

  if len(transaction.ExternalTransactionId) == 0 {
    t.Errorf("ExternalTransactionId expected")
    return
  }

  if transaction.Amount != 10.0 {
    t.Errorf("Transaction amount expected: %v, returned: %v", 10.0, transaction.Amount)
    return
  }

}

func TestPayZenCreateSubscription(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  //PayZen.OnDebug()

  CardTokenActive, err := client.Get(KeyCardTokenActive).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card CardTokenActive: %v", err)
    return
  }

  subscription := payzen.NewPayZenSubscription(ShopId, Mode, Cert)

  subscription.OrderId = genUUID()
  subscription.SubscriptionId = genUUID()
  subscription.Description = "Recorrência diária"

  // valor da recorrência
  subscription.Amount = 10
  // valor inicial da recorrência
  subscription.InitialAmount = 0
  // quantas vezes o valor inicial deve ser cobrado
  subscription.InitialAmountNumber = 0
  // data de inicio da cobrança
  subscription.EffectDate = util.DateNow()
  // cobrar no último dia do mês
  subscription.LastDayOfMonth = false
  // quantidade de cobranças
  subscription.Count = 50
  subscription.MonthDay = 10
  subscription.FrequencyByDay = 1
  subscription.Rule = ""
  subscription.Token = CardTokenActive

  result, err := PayZen.PaymentCreateSubscription(subscription)

  if err != nil {
    t.Errorf("Erro ao recuperar informações do token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do token: %v", result.Message)
    return
  }

  if result.SubscriptionInfo.SubscriptionId != subscription.SubscriptionId {
    t.Errorf("SubscriptionId expected: %v, returned: %v", result.SubscriptionInfo.SubscriptionId, subscription.SubscriptionId)
    return
  }

  client.Set(KeySubscriptionId, result.SubscriptionInfo.SubscriptionId, 0)

}

func TestPayZenGetDetailsSubscription(t *testing.T) {

  time.Sleep(1 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")

  PayZen.OnDebug()

  SubscriptionId, err := client.Get(KeySubscriptionId).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar SubscriptionId: %v", err)
    return
  }

  CardTokenActive, err := client.Get(KeyCardTokenActive).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card CardTokenActive: %v", err)
    return
  }


  paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
  paymentFind.SubscriptionId = SubscriptionId
  paymentFind.Token = CardTokenActive

  result, err := PayZen.PaymentGetDetailsSubscription(paymentFind)

  if err != nil {
    t.Errorf("Erro ao recuperar informações da Subscription: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do Subscription: %v", result.Message)
    return
  }

  if result.SubscriptionInfo.SubscriptionId != SubscriptionId {
    t.Errorf("SubscriptionId expected: %v, returned: %v", SubscriptionId, result.SubscriptionInfo.SubscriptionId)
  }

  if result.SubscriptionInfo.TotalPaymentsNumber != 3 {
    t.Errorf("TotalPaymentsNumber expected: %v, returned: %v", 3, result.SubscriptionInfo.TotalPaymentsNumber)
  }

  /*
  if result.SubscriptionInfo.PastPaymentsNumber != 1 {
    t.Errorf("PastPaymentsNumber expected: %v, returned: %v", 3, result.SubscriptionInfo.TotalPaymentsNumber)
  }

  if !result.SubscriptionInfo.Started {
    t.Errorf("Subscription not started")
  }

  if !result.SubscriptionInfo.Active {
    t.Errorf("Subscription is not active")
  }
  */

  if result.SubscriptionInfo.Cancelled {
    t.Errorf("Subscription is cancelled")
  }
}

func TestPayZenUpdateSubscription(t *testing.T) {

  time.Sleep(3 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  PayZen.OnDebug()

  CardTokenActive, err := client.Get(KeyCardTokenActive).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card CardTokenActive: %v", err)
    return
  }

  SubscriptionId, err := client.Get(KeySubscriptionId).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card SubscriptionId: %v", err)
    return
  }

  OrderId, err := client.Get(KeyCardOrderId).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card OrderId: %v", err)
    return
  }

  subscription := payzen.NewPayZenSubscription(ShopId, Mode, Cert)

  subscription.OrderId = OrderId
  subscription.SubscriptionId = SubscriptionId
  subscription.Description = "Subscription UPdate"

  // valor da recorrência
  subscription.Amount = 9
  // valor inicial da recorrência
  subscription.InitialAmount = 0
  // quantas vezes o valor inicial deve ser cobrado
  subscription.InitialAmountNumber = 0
  // data de inicio da cobrança
  subscription.EffectDate = util.DateNow()
  // cobrar no último dia do mês
  subscription.LastDayOfMonth = false
  // quantidade de cobranças
  subscription.Count = 3
  subscription.MonthDay = 10
  subscription.Rule = ""
  subscription.Token = CardTokenActive

  result, err := PayZen.PaymentUpdateSubscription(subscription)

  if err != nil {
    t.Errorf("Erro ao recuperar informações do token: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do token: %v", result.Message)
    return
  }

  if result.SubscriptionInfo.SubscriptionId != subscription.SubscriptionId {
    t.Errorf("SubscriptionId expected: %v, returned: %v", result.SubscriptionInfo.SubscriptionId, subscription.SubscriptionId)
    return
  }
}

func TestPayZenCancelSubscription(t *testing.T) {

  time.Sleep(3 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")

  PayZen.OnDebug()

  SubscriptionId, err := client.Get(KeySubscriptionId).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar SubscriptionId: %v", err)
    return
  }

  CardTokenActive, err := client.Get(KeyCardTokenActive).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card CardTokenActive: %v", err)
    return
  }


  paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
  paymentFind.SubscriptionId = SubscriptionId
  paymentFind.Token = CardTokenActive

  result, err := PayZen.PaymentCancelSubscription(paymentFind)

  if err != nil {
    t.Errorf("Erro ao recuperar informações da Subscription: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do Subscription: %v", result.Message)
    return
  }

  time.Sleep(3 * time.Second)

  result, err = PayZen.PaymentGetDetailsSubscription(paymentFind)

  if err != nil {
    t.Errorf("Erro ao recuperar informações da Subscription: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do Subscription: %v", result.Message)
    return
  }
  if !result.SubscriptionInfo.Cancelled {
    t.Errorf("Subscription is not cancelled")
  }
}

func TestPayZenGetDetailsCancelledSubscription(t *testing.T) {

  time.Sleep(3 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")

  PayZen.OnDebug()

  SubscriptionId, err := client.Get(KeySubscriptionId).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar SubscriptionId: %v", err)
    return
  }

  CardTokenActive, err := client.Get(KeyCardTokenActive).Result()

  if err != nil {
    t.Errorf("Erro ao recuperar card CardTokenActive: %v", err)
    return
  }

  paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
  paymentFind.SubscriptionId = SubscriptionId
  paymentFind.Token = CardTokenActive

  result, err := PayZen.PaymentGetDetailsSubscription(paymentFind)

  if err != nil {
    t.Errorf("Erro ao recuperar informações da Subscription: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao recuperar informações do Subscription: %v", result.Message)
    return
  }

  if !result.SubscriptionInfo.Cancelled {
    t.Errorf("Subscription is not cancelled")
  }
}


func TestPayZenPaymentCancel(t *testing.T) {

  time.Sleep(3 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  PayZen.OnDebug()
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)


  payment.OrderId = genUUID()
  payment.Installments = 1
  payment.Amount = 10.0


  fillCard(payment.Card)
  fillCustomer(payment.Customer)

  result, err := PayZen.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar autorização: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar autorização: %v", result.Message)
  }else{

    if len(result.TransactionId) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionId não informada")
      return
    } else if len(result.TransactionUuid) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionUuid não informada")
      return
    }

    time.Sleep(3 * time.Second)

    paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
    paymentFind.TransactionUuid = result.TransactionUuid

    result, err = PayZen.PaymentCancel(paymentFind)

    if err != nil {
      t.Errorf("Erro ao cancelar autorização: %v", err)
      return
    }

    if result.Error {
      t.Errorf("Erro ao cancelar autorização: %v", result.Message)
      return
    }

    if result.TransactionStatus != payzen.Cancelled{
      t.Errorf("autorização não foi cancelada: %v", result.Message)
    }

  }

}

func TestPayZenPaymentUpdate(t *testing.T) {

  time.Sleep(3 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  //PayZen.OnDebug()
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)


  payment.OrderId = genUUID()
  payment.Installments = 1
  payment.Amount = 10.0


  //fillCard(payment.Card)
  //fillCustomer(payment.Customer)

  result, err := PayZen.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar autorização: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar autorização: %v", result.Message)
  }else{

    if len(result.TransactionId) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionId não informada")
      return
    } else if len(result.TransactionUuid) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionUuid não informada")
      return
    }

    time.Sleep(3 * time.Second)

    transactionUuid := result.TransactionUuid
    payment.Amount = 9.0
    payment.TransactionUuid = transactionUuid

    result, err := PayZen.PaymentUpdate(payment)

    if err != nil {
      t.Errorf("Erro ao criar autorização: %v", err)
      return
    }

    if result.Error {
      t.Errorf("Erro ao criar autorização: %v", result.Message)
      return
    }

    time.Sleep(3 * time.Second)

    paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
    paymentFind.OrderId = payment.OrderId

    result, err = PayZen.PaymentFind(paymentFind)

    if err != nil {
      t.Errorf("Erro ao buscar autorização: %v", err)
      return
    }

    if result.Error {
      t.Errorf("Erro ao buscar autorização: %v", result.Message)
      return
    }

    trans := result.Transactions[0]

    if trans.Amount != payment.Amount {
      t.Errorf("Amount expected: %v, returned: %v", payment.Amount, trans.Amount)
      return
    }

  }

}

func TestPayZenPaymentDuplicate(t *testing.T) {

  time.Sleep(3 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  PayZen.OnDebug()
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)

  orderId := genUUID()

  payment.OrderId = orderId
  payment.Installments = 1
  payment.Amount = 10.0


  fillCard(payment.Card)
  fillCustomer(payment.Customer)

  result, err := PayZen.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar autorização: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar autorização: %v", result.Message)
  }else{

    if len(result.TransactionId) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionId não informada")
      return
    } else if len(result.TransactionUuid) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionUuid não informada")
      return
    }

    transactionUuid := result.TransactionUuid

    time.Sleep(3 * time.Second)

    capture := payzen.NewPayZenCapturePayment(ShopId, Mode, Cert)
    capture.TransactionUuids = transactionUuid
    result, err := PayZen.PaymentCapture(capture)

    if err != nil {
      t.Errorf("Erro ao criar autorização: %v", err)
    }

    if result.Error {
      t.Errorf("Erro ao criar autorização: %v", result.Message)
    }

    time.Sleep(3 * time.Second)

    payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)
    payment.TransactionUuid = transactionUuid // UUID Transação BackOffice
    payment.OrderId = orderId // Referência do pedido BackOffice
    payment.Amount = 10.0

    result, err = PayZen.PaymentDuplicate(payment)

    if err != nil {
      t.Errorf("Erro ao criar autorização 2: %v", err)
      return
    }

    if result.Error {
      t.Errorf("Erro ao criar autorização 2: %v", result.Message)
      return
    }

    time.Sleep(3 * time.Second)

    paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
    paymentFind.OrderId = payment.OrderId

    result, err = PayZen.PaymentFind(paymentFind)

    if err != nil {
      t.Errorf("Erro ao buscar autorização: %v", err)
      return
    }

    if result.Error {
      t.Errorf("Erro ao buscar autorização: %v", result.Message)
      return
    }

    if len(result.Transactions) != 2 {
      t.Errorf("Transactions expected: %v, returned: %v", 2, len(result.Transactions))
      return
    }

  }

}

func TestPayZenPaymentValidate(t *testing.T) {

  time.Sleep(3 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  PayZen.OnDebug()
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)

  orderId := genUUID()

  payment.OrderId = orderId
  payment.Installments = 1
  payment.Amount = 10.0
  payment.ValidationType = payzen.Manual


  fillCard(payment.Card)
  fillCustomer(payment.Customer)

  result, err := PayZen.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar autorização: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar autorização: %v", result.Message)
  }else{

    if len(result.TransactionId) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionId não informada")
      return
    } else if len(result.TransactionUuid) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionUuid não informada")
      return
    }

    fmt.Println("** OrderId", orderId)
    fmt.Println("** TransactionUuid", result.TransactionUuid)

    time.Sleep(3 * time.Second)

    paymentFind := payzen.NewPayZenPaymentFind(ShopId, Mode, Cert)
    paymentFind.TransactionUuid = result.TransactionUuid

    result, err := PayZen.PaymentValidate(paymentFind)

    if err != nil {
      t.Errorf("Erro ao criar autorização 2: %v", err)
      return
    }

    if result.Error {
      t.Errorf("Erro ao criar autorização 2: %v", result.Message)
      return
    }

  }

}

func TestPayZenPaymentRefund(t *testing.T) {

  time.Sleep(3 * time.Second)

  PayZen := payzen.NewPayZen("pt-BR")
  PayZen.OnDebug()
  payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)

  orderId := genUUID()

  payment.OrderId = orderId
  payment.Installments = 1
  payment.Amount = 10.0


  fillCard(payment.Card)
  fillCustomer(payment.Customer)

  result, err := PayZen.PaymentCreate(payment)

  if err != nil {
    t.Errorf("Erro ao criar autorização: %v", err)
    return
  }

  if result.Error {
    t.Errorf("Erro ao criar autorização: %v", result.Message)
  }else{

    if len(result.TransactionId) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionId não informada")
      return
    } else if len(result.TransactionUuid) == 0 {
      t.Errorf("Erro ao criar autorização: %v", "TransactionUuid não informada")
      return
    }

    transactionUuid := result.TransactionUuid

    time.Sleep(3 * time.Second)

    capture := payzen.NewPayZenCapturePayment(ShopId, Mode, Cert)
    capture.TransactionUuids = transactionUuid
    result, err := PayZen.PaymentCapture(capture)

    if err != nil {
      t.Errorf("Erro ao criar autorização: %v", err)
    }

    if result.Error {
      t.Errorf("Erro ao criar autorização: %v", result.Message)
    }

    time.Sleep(3 * time.Second)

    payment := payzen.NewPayZenPayment(ShopId, Mode, Cert)
    payment.TransactionUuid = transactionUuid // UUID Transação BackOffice
    payment.Amount = 5.0

    result, err = PayZen.PaymentRefund(payment)

    if err != nil {
      t.Errorf("Erro ao criar refund 2: %v", err)
      return
    }

    if result.Error {
      t.Errorf("Erro ao criar refund 2: %v", result.Message)
      return
    }


  }

}
