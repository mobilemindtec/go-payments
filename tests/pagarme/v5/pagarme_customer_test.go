package v5

import (
	"fmt"
	pagarme "github.com/mobilemindtec/go-payments/pagarme/v5"
	gopayments "github.com/mobilemindtec/go-payments/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5CustomerCreate
func TestPagarmev5CustomerCreate(t *testing.T) {

	Pagarme := pagarme.NewPagarmeCustomer("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	customer := newCustomer()

	result := Pagarme.Create(customer)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty customer response")

	result.
		Right().
		Foreach(func(c *pagarme.Customer) {
			assert.NotEmptyf(t, result.UnwrapRight().Id, "empty customer id")
			gopayments.CacheClient.Set("CustomerId", c.Id, 0)
		})
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5CustomerGet
func TestPagarmev5CustomerGet(t *testing.T) {

	Pagarme := pagarme.NewPagarmeCustomer("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	customerId := "cus_zBknRWwFWfK4RQxa"
	result := Pagarme.Get(customerId)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty customer response")

	result.
		Right().
		Foreach(func(c pagarme.CustomerPtr) {
			fmt.Print(c.Id)
		})
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5CustomerList
func TestPagarmev5CustomerList(t *testing.T) {

	Pagarme := pagarme.NewPagarmeCustomer("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	result := Pagarme.List(pagarme.NewCustomerQuery())

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty customer response")

	result.
		Right().
		ListForeach(func(c pagarme.CustomerPtr) {
			fmt.Print(c.Id)
		})
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5CustomerUpdate
func TestPagarmev5CustomerUpdate(t *testing.T) {

	Pagarme := pagarme.NewPagarmeCustomer("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	customer := newCustomer()
	customerId := "cus_zBknRWwFWfK4RQxa"
	customer.Id = customerId
	result := Pagarme.Update(customer)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty customer response")

	result.
		Right().
		Foreach(func(c pagarme.CustomerPtr) {
			fmt.Print(c.Id)
		})
}
