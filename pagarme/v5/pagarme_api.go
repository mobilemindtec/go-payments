package v5

type PagarmeApi struct {
	Card         *PagarmeCard
	Customer     *PagarmeCustomer
	Order        *PagarmeOrder
	Plan         *PagarmePlan
	Subscription *PagarmeSubscription
	Invoice      *PagarmeInvoice
}

func NewPagarmeApi(lang string, auth *Authentication) *PagarmeApi {
	card := &PagarmeCard{}
	card.Pagarme.init(lang, auth)

	customer := &PagarmeCustomer{}
	customer.Pagarme.init(lang, auth)

	order := &PagarmeOrder{}
	order.Pagarme.init(lang, auth)

	plan := &PagarmePlan{}
	plan.Pagarme.init(lang, auth)

	subscription := &PagarmeSubscription{}
	subscription.Pagarme.init(lang, auth)

	invoice := &PagarmeInvoice{}
	invoice.Pagarme.init(lang, auth)

	return &PagarmeApi{
		Card:         card,
		Customer:     customer,
		Order:        order,
		Plan:         plan,
		Subscription: subscription,
		Invoice:      invoice,
	}
}

func (this *PagarmeApi) DebugOn() {
	this.Card.DebugOn()
	this.Customer.DebugOn()
	this.Order.DebugOn()
	this.Plan.DebugOn()
	this.Subscription.DebugOn()
}
