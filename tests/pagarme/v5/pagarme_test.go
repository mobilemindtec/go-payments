package v5

import (
	_ "github.com/mobilemindtec/go-payments/api"
	pagarme "github.com/mobilemindtec/go-payments/pagarme/v5"
	"github.com/mobilemindtec/go-payments/tests"
)

func fillCreditCard(creditCard *pagarme.CreditCard) *pagarme.CreditCard {
	//creditCard := new(pagarme.CreditCard)
	creditCard.StatementDescriptor = "MMIND"
	creditCard.Card.Number = "4901720080344448"
	creditCard.Card.HolderName = "Aardvark Silva"
	creditCard.Card.HolderDocument = "83361855004"
	creditCard.Card.ExpMonth = 12
	creditCard.Card.ExpYear = 2028
	creditCard.Card.Cvv = "314"
	creditCard.Card.Brand = "VISA"

	fillBillingAddress(creditCard.Card.BillingAddress)

	return creditCard
}

func fillBillingAddress(billingAddress *pagarme.BillingAddress) *pagarme.BillingAddress {
	billingAddress.City = "Bento Gonçalves"
	billingAddress.State = "RS"
	billingAddress.ZipCode = "95700000"
	billingAddress.Line1 = "Rua Vitória, 255"
	return billingAddress
}

func newCustomer() *pagarme.Customer {
	customer := pagarme.NewCustomer()
	customer.Name = "Ricardo Bocchi"
	customer.Email = "ricardobocchi@gmail.com"
	customer.Address.City = "Bento Gonçalves"
	customer.Address.State = "RS"
	customer.Address.ZipCode = "95700000"
	customer.Address.Line1 = "Rua Vitória, 255"
	customer.Birthdate = "1986-11-27"
	customer.Gender = pagarme.Male
	customer.Code = gopayments.GenUUID()
	customer.Document = "83361855004"
	customer.Phones.MobilePhone = pagarme.NewPhone("55", "54", "999767081")
	return customer
}

func newOrder() *pagarme.Order {
	order := pagarme.NewOrder()
	order.Code = gopayments.GenUUID()
	order.Customer = newCustomer()
	order.Ip = "127.0.0.1"
	order.
		AddPayment(1000, pagarme.MethodCreditCard).
		WithCreditCard(func(creditCard *pagarme.CreditCard) {
			fillCreditCard(creditCard)
		})

	item := order.AddItem()
	item.Amount = 1000
	item.Code = gopayments.GenUUID()
	item.Description = "test item"
	return order
}
