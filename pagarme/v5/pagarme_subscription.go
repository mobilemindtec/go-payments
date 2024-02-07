package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/v2/either"
	"github.com/mobilemindtec/go-utils/v2/maps"
	"reflect"
)

type CancelPendingInvoices bool
type CardId string

const (
	CancelPendingInvoicesYes CancelPendingInvoices = true
	CancelPendingInvoicesNo  CancelPendingInvoices = false
)

type PagarmeSubscription struct {
	Pagarme
}

func NewPagarmeSubscription(lang string, auth *Authentication) *PagarmeSubscription {
	p := &PagarmeSubscription{}
	p.Pagarme.init(lang, auth)
	return p
}

func (this *PagarmeSubscription) Create(subscription SubscriptionPtr) *either.Either[*ErrorResponse, SubscriptionPtr] {

	if !this.validate(subscription) {
		return either.Left[*ErrorResponse, SubscriptionPtr](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.ValidationErrors))
	}

	return either.
		MapIf(
			this.post("/subscriptions", subscription, createParser[Subscription]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SubscriptionPtr {
				return e.UnwrapRight().Content.(SubscriptionPtr)
			})
}

func (this *PagarmeSubscription) Get(id string) *either.Either[*ErrorResponse, SubscriptionPtr] {

	if empty, left := checkEmpty[SubscriptionPtr]("subscription id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/subscriptions/%v", id)

	return either.
		MapIf(
			this.get(uri, createParser[Subscription]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SubscriptionPtr {
				return e.UnwrapRight().Content.(SubscriptionPtr)
			})
}

func (this *PagarmeSubscription) List(query *SubscriptionQuery) *either.Either[*ErrorResponse, *Content[Subscriptions]] {

	uri := fmt.Sprintf("/subscriptions/?%v", query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, createParserContent[Subscriptions]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) *Content[Subscriptions] {
				return e.UnwrapRight().Content.(*Content[Subscriptions])
			})
}

func (this *PagarmeSubscription) ListItems(id string) *either.Either[*ErrorResponse, *Content[SubscriptionItems]] {

	uri := fmt.Sprintf("/subscriptions/%v/items", id)

	return either.
		MapIf(
			this.get(uri, createParserContent[SubscriptionItems]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) *Content[SubscriptionItems] {
				return e.UnwrapRight().Content.(*Content[SubscriptionItems])
			})
}

func (this *PagarmeSubscription) Cancel(id string, cancelPendingInvoices CancelPendingInvoices) *either.Either[*ErrorResponse, bool] {

	if empty, left := checkEmpty[bool]("subscription id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/subscriptions/%v", id)
	payload := maps.JSON("cancel_pending_invoices", cancelPendingInvoices)

	return either.
		MapIf(
			this.delete(uri, payload),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) bool {
				return true
			})
}

func (this *PagarmeSubscription) UpdateCard(updateData *SubscriptionUpdate) *either.Either[*ErrorResponse, bool] {

	if empty, left := checkEmpty[bool]("subscription id", updateData.Id); empty {
		return left
	}

	if (len(updateData.CardId) == 0 && updateData.Card == nil) || (len(updateData.CardId) > 0 && updateData.Card != nil) {
		return either.Left[*ErrorResponse, bool](
			NewErrorResponse("card id or card is required"))
	}

	uri := fmt.Sprintf("/subscriptions/%v/card", updateData.Id)

	return either.
		MapIf(
			this.patch(uri, updateData),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) bool {
				return true
			})
}

func (this *PagarmeSubscription) UpdatePaymentMethod(updateData *SubscriptionUpdate) *either.Either[*ErrorResponse, bool] {

	if empty, left := checkEmpty[bool]("subscription id", updateData.Id); empty {
		return left
	}

	if len(updateData.PaymentMethod) == 0 {
		return either.Left[*ErrorResponse, bool](
			NewErrorResponse("payment_method is required"))
	}

	if (len(updateData.CardId) == 0 && len(updateData.CardToken) == 0) || (len(updateData.CardId) > 0 && len(updateData.CardToken) > 0) {
		return either.Left[*ErrorResponse, bool](
			NewErrorResponse("card is or card token is required"))
	}

	uri := fmt.Sprintf("/subscriptions/%v/payment-method", updateData.Id)

	return either.
		MapIf(
			this.patch(uri, updateData),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) bool {
				return true
			})
}

func (this *PagarmeSubscription) UpdateItem(subscriptionId string, itemId string, item SubscriptionItemPtr) *either.Either[*ErrorResponse, SubscriptionItemPtr] {

	if empty, left := checkEmpty[SubscriptionItemPtr]("subscription id and subscription item id", subscriptionId, itemId); empty {
		return left
	}

	uri := fmt.Sprintf("/subscriptions/%v/items/%v", subscriptionId, itemId)

	return either.
		MapIf(
			this.put(uri, item, createParser[SubscriptionPtr]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SubscriptionItemPtr {
				return e.UnwrapRight().Content.(SubscriptionItemPtr)
			})
}

func (this *PagarmeSubscription) validate(subscription *Subscription) bool {
	this.EntityValidator.AddEntity(subscription)
	this.EntityValidator.AddValidationForType(reflect.TypeOf(subscription), subscriptionValidator)
	return this.processValidator()
}

func subscriptionValidator(entity interface{}, validator *validator.Validation) {

	s := entity.(*Subscription)

	if s.BillingType == ExactDay {
		if s.BillingDay <= 0 {
			validator.SetError("BillingDay", "BillingDay is required to BillingType equal ExactDay")
		}
	}

	if s.Customer == nil && len(s.CustomerId) == 0 {
		validator.SetError("Customer", "Customer or CustomerId is required")
	}

	if s.Card == nil && len(s.CardId) == 0 && len(s.CardToken) == 0 {
		validator.SetError("Customer", "Card, CardId or CardToken is required")
	}
}
