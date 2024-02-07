package v5

import (
	pagarme "github.com/mobilemindtec/go-payments/pagarme/v5"
	gopayments "github.com/mobilemindtec/go-payments/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmev5OrderCreate
func TestPagarmev5OrderCreate(t *testing.T) {

	Pagarme := pagarme.NewPagarmeOrder("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	order := newOrder()
	result := Pagarme.Create(order)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty order response")
	if result.IsRight() {
		assert.NotEmptyf(t, result.UnwrapRight().Id, "empty order id")
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmev5OrderList
func TestPagarmev5OrderList(t *testing.T) {

	Pagarme := pagarme.NewPagarmeOrder("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	result := Pagarme.List(pagarme.NewOrderQuery())

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty order response")

	result.
		Right().
		Foreach(func(orders pagarme.Orders) {

			l := len(orders)

			assert.True(t, l > 0)

			if l > 0 {
				assert.True(t, len(orders[0].Items) > 0)
			}
		})
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmev5OrderGet
func TestPagarmev5OrderGet(t *testing.T) {

	Pagarme := pagarme.NewPagarmeOrder("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	orderId := "or_xD4GJmgiLeCyndkn"
	result := Pagarme.Get(orderId)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty order response")
}
