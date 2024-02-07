package v5

import (
	gopayments "github.com/mobilemindtec/go-payments/tests"
	"testing"
	"time"

	_ "github.com/mobilemindtec/go-payments/api"
	pagarme "github.com/mobilemindtec/go-payments/pagarme/v5"
	"github.com/stretchr/testify/assert"
)

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5SubscriptionCreateWithoutPlan
func TestPagarmev5SubscriptionCreateWithoutPlan(t *testing.T) {

	Pagarme := pagarme.
		NewPagarmeSubscription("pt-BR",
			pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))

	Pagarme.DebugOn()

	customer := newCustomer()
	card := fillCreditCard(pagarme.NewCreditCard()).Card

	sub := pagarme.
		NewSubscription().
		WithCustomer(customer).
		WithCard(card).
		SetStartAt(time.Now().Add(24*5*time.Hour)).
		SetIntervalRule(pagarme.Month, 1).
		AddItem(
			pagarme.
				NewSubscriptionItem("Item test", 1).
				SetPricingScheme(pagarme.NewPricingScheme(100)))
	sub.StatementDescriptor = "MMIND"

	result := Pagarme.Create(sub)

	assert.False(t, result.IsLeft())
	assert.True(t, result.Right().NonEmpty())
	if result.IsRight() {
		assert.NotEmpty(t, result.UnwrapRight().Id)
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5SubscriptionCreateFromPlan
func TestPagarmev5SubscriptionCreateFromPlan(t *testing.T) {

	PagarmeSubs := pagarme.
		NewPagarmeSubscription("pt-BR",
			pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))

	PagarmePlan := pagarme.
		NewPagarmePlan("pt-BR",
			pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))

	PagarmeSubs.DebugOn()

	plan := pagarme.
		NewPlan("Plan test").
		AddPaymentMethod(pagarme.MethodCreditCard, pagarme.MethodBoleto).
		AddPlanItem(pagarme.NewPlanItem("Item test", 1, 100)).
		SetIntervalRule(pagarme.Month, 1) // mensal

	presult := PagarmePlan.Create(plan)

	assert.True(t, presult.IsRight())

	pid := presult.Right().Get().Id

	customer := newCustomer()
	card := fillCreditCard(pagarme.NewCreditCard()).Card

	sub := pagarme.
		NewSubscription().
		WithCustomer(customer).
		WithCard(card).
		WithPlanId(pid).
		SetStartAt(time.Now().Add(24*5*time.Hour)).
		SetIntervalRule(pagarme.Month, 1).
		AddItem(
			pagarme.
				NewSubscriptionItem("Item test", 1).
				SetPricingScheme(pagarme.NewPricingScheme(100)))
	sub.StatementDescriptor = "MMIND"

	result := PagarmeSubs.Create(sub)

	assert.False(t, result.IsLeft())
	assert.True(t, result.Right().NonEmpty())
	if result.IsRight() {
		assert.NotEmpty(t, result.UnwrapRight().Id)
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5SubscriptionGet
func TestPagarmev5SubscriptionGet(t *testing.T) {

	Pagarme := pagarme.
		NewPagarmeSubscription("pt-BR",
			pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))

	Pagarme.DebugOn()

	result := Pagarme.Get("sub_pG6KjZ0iOivgNRw2")

	assert.False(t, result.IsLeft())
	assert.True(t, result.Right().NonEmpty())
	if result.IsRight() {
		assert.NotEmpty(t, result.UnwrapRight().Id)
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5SubscriptionList
func TestPagarmev5SubscriptionList(t *testing.T) {

	Pagarme := pagarme.
		NewPagarmeSubscription("pt-BR",
			pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))

	Pagarme.DebugOn()

	q := pagarme.NewSubscriptionQuery()
	q.Size = 10
	q.Page = 1
	q.CreatedSince = time.Now()
	q.CustomerId = "123"
	result := Pagarme.List(q)

	assert.False(t, result.IsLeft())
	assert.True(t, result.Right().NonEmpty())
	if result.IsRight() {
		assert.True(t, len(result.UnwrapRight().Data) > 0)
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5SubscriptionItemsList
func TestPagarmev5SubscriptionItemsList(t *testing.T) {

	Pagarme := pagarme.
		NewPagarmeSubscription("pt-BR",
			pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))

	Pagarme.DebugOn()

	result := Pagarme.ListItems("sub_jw4qrlaUqUwMQWmO")

	assert.False(t, result.IsLeft())
	assert.True(t, result.Right().NonEmpty())
	if result.IsRight() {
		assert.True(t, len(result.UnwrapRight().Data) > 0)
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5SubscriptionCancel
func TestPagarmev5SubscriptionCancel(t *testing.T) {

	Pagarme := pagarme.
		NewPagarmeSubscription("pt-BR",
			pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))

	Pagarme.DebugOn()

	result := Pagarme.Cancel("sub_jw4qrlaUqUwMQWmO", pagarme.CancelPendingInvoicesYes)

	assert.False(t, result.IsLeft())
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5SubscriptionEditWithCardId
func TestPagarmev5SubscriptionEditWithCardId(t *testing.T) {

	Pagarme := pagarme.
		NewPagarmeApi("pt-BR",
			pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))

	Pagarme.DebugOn()

	card := fillCreditCard(pagarme.NewCreditCard()).Card
	updateData := pagarme.NewSubscriptionUpdate("sub_pG6KjZ0iOivgNRw2")
	updateData.Card = card
	result := Pagarme.Subscription.UpdateCard(updateData)

	assert.False(t, result.IsLeft())
}
