package gopayments

import (
  "github.com/mobilemindtec/go-payments/pagarme"
  _ "github.com/mobilemindtec/go-utils/app/util"
	_ "github.com/satori/go.uuid"
	_ "github.com/go-redis/redis"
	"testing"
	_ "time"
	"fmt"
	_ "os"
)


const (

	ApiKey = ""
	CryptoKey = ""

)



func pagarmeFillCard(card *pagarme.PagarmeCard) {
  card.Number = "4901720080344448"
  card.HolderName = "Aardvark Silva"
  card.ExpirationDate = "1213"  
  card.Cvv = "314"
}

func pagarmeFillPayments(payment *pagarme.PagarmePayment) {
  customer := new(pagarme.PagarmeCustomer)

  pagarmefillCustomer(customer)

  //payment.Amount = 10  
  payment.Installments = 1
  payment.Customer = customer
  payment.PaymentMethod = pagarme.PAYMENT_METHOD_CREDIT_CARD
  //payment.PostbackUrl = "https://mobilemind.com.br"
  payment.SoftDescriptor = "Mobile Mind"
  //payment.Metadata = ""
  payment.Capture = true
  //payment.BoletoExpirationDate
  //payment.BoletoInstructions
  //payment.CardId
  payment.CardHolderName = "Ricardo Bocchi"
  payment.CardExpirationDate = "0921"
  payment.CardNumber = "4024007140405134"
  payment.CardCvv = "680"
  //payment.CardHash
  //payment.SplitRules
}



func pagarmefillCustomer(customer *pagarme.PagarmeCustomer) {
  customer.Email = "ricardobocchi@gmail.com"
  customer.Name = "Ricardo Bocchi"
  customer.DocumentNumber = "83361855004"

  customer.Phone = new(pagarme.PagarmePhone)
  customer.Phone.Ddd = "054"
  customer.Phone.Number = "999767081"

  customer.Address = new(pagarme.PagarmeAddress)
  customer.Address.Neighborhood = "Botafogo"
  customer.Address.Street = "Vitoria"
  customer.Address.StreetNumber = "255"
  customer.Address.ZipCode = "95700540"
  customer.Address.City = "Bento Goncalves"
  customer.Address.State = "RS"
}

func TestPagarmeGetCardHashKey(t *testing.T) {
		
	Pagarme := pagarme.NewPagarmeServiceWithCert("pt-BR", ApiKey, CryptoKey)
	result, err :=  Pagarme.GetCardHashKey()

  if err != nil {
  	t.Errorf("Erro ao criar card hash key: %v", err)
  }else{
  	t.Log(fmt.Sprintf("result = %v", result))
  }
}

func TestPagarmeEncryptCard(t *testing.T) {
	
	Card := new(pagarme.PagarmeCard)
	Pagarme := pagarme.NewPagarmeServiceWithCert("pt-BR", ApiKey, CryptoKey)

	pagarmeFillCard(Card)

	result, err :=  Pagarme.EncryptCard(Card)

  if err != nil {
  	t.Errorf("Erro ao encrypt card: %v", err)
  }else{
  	t.Log(fmt.Sprintf("result = %v", result))  	
  }
}

func TestPagarmeCreateCard(t *testing.T) {
	
	Card := new(pagarme.PagarmeCard)
	Pagarme := pagarme.NewPagarmeServiceWithCert("pt-BR", ApiKey, CryptoKey)

	pagarmeFillCard(Card)

  Card.CustomerId = "123456"

  Pagarme.CreateCard(Card)

}

func TestPagarmeCreateCustomer(t *testing.T) {
  
  customer := pagarme.NewPagarmeCustomer()
  Pagarme := pagarme.NewPagarmeServiceWithCert("pt-BR", ApiKey, CryptoKey)

  pagarmefillCustomer(customer)

  err := Pagarme.CreateCustomer(customer)

  if err != nil {
    t.Errorf("Erro ao create customer: %v", err)
  }else{
    t.Log(fmt.Sprintf("result = %v", customer.Id))   
  }

}

func TestPagarmeCreatePaymentCard(t *testing.T) {
  
  payment := pagarme.NewPagarmePaymentCard(1)
  Pagarme := pagarme.NewPagarmeServiceWithCert("pt-BR", ApiKey, CryptoKey)

  pagarmeFillPayments(payment)

  result, err := Pagarme.CreatePayment(payment)

  if err != nil {
    t.Errorf("Erro ao create card payment: %v", err)
  }else{
    //t.Log(fmt.Sprintf("result = %v", customer.Id))   


    if result.TransactionId == 0 {
      t.Errorf("TransactionId cant be empty")
      return
    }

    captureData := pagarme.NewCaptureData(fmt.Sprintf("%v"), result.TransactionId, 1)

    captureResult, err := Pagarme.Capture(captureData)

    if err != nil {
      t.Errorf("Erro ao capture data: %v", err)
    }else{
      t.Log(captureResult)
    }

  }

}

