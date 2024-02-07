package v5

import (
	gopayments "github.com/mobilemindtec/go-payments/tests"
	"testing"

	_ "github.com/mobilemindtec/go-payments/api"
	pagarme "github.com/mobilemindtec/go-payments/pagarme/v5"
	"github.com/stretchr/testify/assert"
)

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5PlanCreate
func TestPagarmev5PlanCreate(t *testing.T) {

	Pagarme := pagarme.NewPagarmePlan("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	plan := pagarme.
		NewPlan("Plan test").
		AddPaymentMethod(pagarme.MethodCreditCard, pagarme.MethodBoleto).
		AddPlanItem(pagarme.NewPlanItem("Item test", 1, 100)).
		SetIntervalRule(pagarme.Month, 1) // mensal

	plan.StatementDescriptor = "MMIND"

	result := Pagarme.Create(plan)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty plan response")
	if result.IsRight() {
		assert.NotEmptyf(t, result.UnwrapRight().Id, "empty plan id")
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5PlanEdit
func TestPagarmev5PlanEdit(t *testing.T) {

	Pagarme := pagarme.NewPagarmePlan("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	planId := "plan_6wl3pk2HrxcrjdY7"

	result := Pagarme.Get(planId)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty plan response")
	if result.IsRight() {
		assert.NotEmptyf(t, result.UnwrapRight().Id, "empty plan id")
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5PlanList
func TestPagarmev5PlanList(t *testing.T) {

	Pagarme := pagarme.NewPagarmePlan("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	result := Pagarme.List(pagarme.NewPlanQuery())

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty plan response")

	result.
		Right().
		Foreach(func(plans pagarme.Plans) {

			l := len(plans)

			assert.True(t, l > 0)

			if l > 0 {
				assert.True(t, len(plans[0].Items) > 0)
			}
		})
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5PlanUpdate
func TestPagarmev5PlanUpdate(t *testing.T) {

	Pagarme := pagarme.NewPagarmePlan("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	plan := pagarme.
		NewPlan("Plan test").
		AddPaymentMethod(pagarme.MethodCreditCard, pagarme.MethodBoleto).
		AddPlanItem(pagarme.NewPlanItem("Item test", 1, 100)).
		SetIntervalRule(pagarme.Month, 1) // mensal

	plan.StatementDescriptor = "MMIND"
	plan.Id = "plan_6wl3pk2HrxcrjdY7"
	plan.Status = pagarme.Active

	result := Pagarme.Update(plan)

	assert.False(t, result.IsLeft())
	assert.Truef(t, result.Right().NonEmpty(), "empty plan response")
	if result.IsRight() {
		assert.NotEmptyf(t, result.UnwrapRight().Id, "empty plan id")
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests/pagarme/v5 -run TestPagarmev5PlanDelete
func TestPagarmev5PlanDelete(t *testing.T) {

	Pagarme := pagarme.NewPagarmePlan("pt-BR", pagarme.NewAuthentication(gopayments.SecretKey, gopayments.PublicKey))
	Pagarme.DebugOn()

	planId := "plan_XvmYa6rFoRs4Gq6Z"

	result := Pagarme.Delete(planId)

	assert.False(t, result.IsLeft())
}
