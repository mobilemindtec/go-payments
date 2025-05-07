package v5

type ServiceRefererName string

type PagarmeApi struct {
	Card         *PagarmeCard
	Customer     *PagarmeCustomer
	Order        *PagarmeOrder
	Plan         *PagarmePlan
	Subscription *PagarmeSubscription
	Invoice      *PagarmeInvoice
	Charge       *PagarmeCharge
	Recipient    *PagarmeRecipient
	ServiceRefererName ServiceRefererName
}

func NewPagarmeApi(lang string, auth *Authentication, serviceRefererName ServiceRefererName) *PagarmeApi {
	card := &PagarmeCard{}
	card.Pagarme.init(lang, auth, serviceRefererName)

	customer := &PagarmeCustomer{}
	customer.Pagarme.init(lang, auth, serviceRefererName)

	order := &PagarmeOrder{}
	order.Pagarme.init(lang, auth, serviceRefererName)

	plan := &PagarmePlan{}
	plan.Pagarme.init(lang, auth, serviceRefererName)

	subscription := &PagarmeSubscription{}
	subscription.Pagarme.init(lang, auth, serviceRefererName)

	invoice := &PagarmeInvoice{}
	invoice.Pagarme.init(lang, auth, serviceRefererName)

	charge := &PagarmeCharge{}
	charge.Pagarme.init(lang, auth, serviceRefererName)

	recipient := &PagarmeRecipient{}
	recipient.Pagarme.init(lang, auth, serviceRefererName)

	return &PagarmeApi{
		Card:         card,
		Customer:     customer,
		Order:        order,
		Plan:         plan,
		Subscription: subscription,
		Invoice:      invoice,
		Charge:       charge,
		Recipient:    recipient,
		ServiceRefererName: serviceRefererName,
	}
}

func (this *PagarmeApi) DebugOn() {
	this.Card.DebugOn()
	this.Customer.DebugOn()
	this.Order.DebugOn()
	this.Plan.DebugOn()
	this.Subscription.DebugOn()
}
