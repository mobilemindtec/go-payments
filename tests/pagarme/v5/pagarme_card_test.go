package v5

import (
	gopayments "github.com/mobilemindtec/go-payments/tests"
	"testing"

	_ "github.com/mobilemindtec/go-payments/api"
	pagarme "github.com/mobilemindtec/go-payments/pagarme/v5"
	"github.com/stretchr/testify/assert"
)

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5CardCreate
func TestPagarmev5CardCreate(t *testing.T) {

	Pagarme := pagarme.NewPagarmeCard("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	card := fillCreditCard(pagarme.NewCreditCard()).Card
	customerId := "cus_zBknRWwFWfK4RQxa"
	result := Pagarme.Create(customerId, card)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty card response")
	if result.IsRight() {
		assert.NotEmptyf(t, result.UnwrapRight().Data.Id, "empty card id")
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5CardList
func TestPagarmev5CardList(t *testing.T) {

	Pagarme := pagarme.NewPagarmeCard("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	customerId := "cus_zBknRWwFWfK4RQxa"
	result := Pagarme.List(customerId)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty card response")
	if result.IsRight() {
		assert.True(t, result.Right().SliceNonEmpty())
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5CardGet
func TestPagarmev5CardGet(t *testing.T) {

	Pagarme := pagarme.NewPagarmeCard("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	customerId := "cus_zBknRWwFWfK4RQxa"
	cardId := "card_ZrBONNuP4TzROQ2j"
	result := Pagarme.Get(customerId, cardId)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty card response")
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5CardUpdate
func TestPagarmev5CardUpdate(t *testing.T) {

	Pagarme := pagarme.NewPagarmeCard("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	customerId := "cus_zBknRWwFWfK4RQxa"
	cardId := "card_ZrBONNuP4TzROQ2j"
	result := Pagarme.Get(customerId, cardId)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty card response")

	result.
		Right().
		Foreach(func(card pagarme.SuccessCard) {
			creditCard := fillCreditCard(pagarme.NewCreditCard())
			creditCard.Card.Id = card.Data.Id
			result := Pagarme.Update(customerId, creditCard.Card)
			assert.False(t, result.IsLeft())
			assert.Truef(t, result.Right().NonEmpty(), "empty card response")
		})
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5CardRenew
func TestPagarmev5CardRenew(t *testing.T) {

	Pagarme := pagarme.NewPagarmeCard("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	customerId := "cus_zBknRWwFWfK4RQxa"
	cardId := "card_ZrBONNuP4TzROQ2j"
	result := Pagarme.Renew(customerId, cardId)

	assert.False(t, result.IsLeft())

}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5CardDelete
func TestPagarmev5CardDelete(t *testing.T) {

	Pagarme := pagarme.NewPagarmeCard("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	customerId := "cus_zBknRWwFWfK4RQxa"
	cardId := "card_ZrBONNuP4TzROQ2j"
	result := Pagarme.Delete(customerId, cardId)

	assert.False(t, result.IsLeft())

}
