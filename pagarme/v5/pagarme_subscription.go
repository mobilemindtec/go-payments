package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/v2/either"
	"github.com/mobilemindtec/go-utils/v2/maps"
	"reflect"
)

type SuccessSubscription = *Success[SubscriptionPtr]
type SuccessSubscriptions = *Success[Subscriptions]

type SuccessSubscriptionItem = *Success[SubscriptionItemPtr]
type SuccessSubscriptionItems = *Success[SubscriptionItems]

type CancelPendingInvoices bool
type CardId string

const (
	CancelPendingInvoicesYes CancelPendingInvoices = true
	CancelPendingInvoicesNo  CancelPendingInvoices = false
)

type PagarmeSubscription struct {
	Pagarme
}

func NewPagarmeSubscription(lang string, auth *Authentication, serviceRefererName ServiceRefererName) *PagarmeSubscription {
	p := &PagarmeSubscription{}
	p.Pagarme.init(lang, auth, serviceRefererName)
	return p
}

func (this *PagarmeSubscription) Create(subscription SubscriptionPtr) *either.Either[*ErrorResponse, SuccessSubscription] {

	if !this.validate(subscription) {
		return either.Left[*ErrorResponse, SuccessSubscription](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.validationsToMapOfStringSlice()))
	}

	return either.
		MapIf(
			this.post("/subscriptions", subscription, createParser[Subscription]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessSubscription {
				return NewSuccess[SubscriptionPtr](e.UnwrapRight())
			})
}

func (this *PagarmeSubscription) Get(id string) *either.Either[*ErrorResponse, SuccessSubscription] {

	if empty, left := checkEmpty[SuccessSubscription]("subscription id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/subscriptions/%v", id)

	return either.
		MapIf(
			this.get(uri, createParser[Subscription]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessSubscription {
				return NewSuccess[SubscriptionPtr](e.UnwrapRight())
			})
}

func (this *PagarmeSubscription) List(query *SubscriptionQuery) *either.Either[*ErrorResponse, SuccessSubscriptions] {

	uri := fmt.Sprintf("/subscriptions/?%v", query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, createParserContent[Subscriptions]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessSubscriptions {
				return NewSuccessSlice[Subscriptions](e.UnwrapRight())
			})
}

func (this *PagarmeSubscription) ListItems(id string) *either.Either[*ErrorResponse, SuccessSubscriptionItems] {

	uri := fmt.Sprintf("/subscriptions/%v/items", id)

	return either.
		MapIf(
			this.get(uri, createParserContent[SubscriptionItems]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessSubscriptionItems {
				return NewSuccess[SubscriptionItems](e.UnwrapRight())
			})
}

func (this *PagarmeSubscription) Cancel(id string, cancelPendingInvoices CancelPendingInvoices) *either.Either[*ErrorResponse, SuccessBool] {

	if empty, left := checkEmpty[SuccessBool]("subscription id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/subscriptions/%v", id)
	payload := maps.JSON("cancel_pending_invoices", cancelPendingInvoices)

	return either.
		MapIf(
			this.delete(uri, payload),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessBool {
				return NewSuccessWithValue[bool](e.UnwrapRight(), true)
			})
}

func (this *PagarmeSubscription) UpdateCard(updateData *SubscriptionUpdate) *either.Either[*ErrorResponse, SuccessBool] {

	if empty, left := checkEmpty[SuccessBool]("subscription id", updateData.Id); empty {
		return left
	}

	if len(updateData.CardId) == 0 && updateData.Card == nil {
		return either.Left[*ErrorResponse, SuccessBool](
			NewErrorResponse("card id or card is required"))
	}

	uri := fmt.Sprintf("/subscriptions/%v/card", updateData.Id)

	return either.
		MapIf(
			this.patch(uri, updateData),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessBool {
				return NewSuccessWithValue[bool](e.UnwrapRight(), true)
			})
}

func (this *PagarmeSubscription) UpdatePaymentMethod(updateData *SubscriptionUpdate) *either.Either[*ErrorResponse, SuccessBool] {

	if empty, left := checkEmpty[SuccessBool]("subscription id", updateData.Id); empty {
		return left
	}

	if len(updateData.PaymentMethod) == 0 {
		return either.Left[*ErrorResponse, SuccessBool](
			NewErrorResponse("payment_method is required"))
	}

	if (len(updateData.CardId) == 0 && len(updateData.CardToken) == 0) || (len(updateData.CardId) > 0 && len(updateData.CardToken) > 0) {
		return either.Left[*ErrorResponse, SuccessBool](
			NewErrorResponse("card is or card token is required"))
	}

	uri := fmt.Sprintf("/subscriptions/%v/payment-method", updateData.Id)

	return either.
		MapIf(
			this.patch(uri, updateData),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessBool {
				return NewSuccessWithValue[bool](e.UnwrapRight(), true)
			})
}

func (this *PagarmeSubscription) UpdateItem(subscriptionId string, itemId string, item SubscriptionItemPtr) *either.Either[*ErrorResponse, SuccessSubscriptionItem] {

	if empty, left := checkEmpty[SuccessSubscriptionItem]("subscription id and subscription item id", subscriptionId, itemId); empty {
		return left
	}

	uri := fmt.Sprintf("/subscriptions/%v/items/%v", subscriptionId, itemId)

	return either.
		MapIf(
			this.put(uri, item, createParser[SubscriptionItem]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessSubscriptionItem {
				return NewSuccess[SubscriptionItemPtr](e.UnwrapRight())
			})
}

func (this *PagarmeSubscription) validate(subscription *Subscription) bool {
	this.EntityValidator.AddEntity(subscription)
	this.EntityValidator.AddValidationForType(reflect.TypeOf(subscription), this.subscriptionValidator)
	return this.processValidator()
}


func (this *PagarmeSubscription) subscriptionValidator(entity interface{}, validator *validator.Validation) {

	s := entity.(*Subscription)

	if s.BillingType == ExactDay {
		if s.BillingDay <= 0 {
			validator.SetError("BillingDay", "BillingDay is required to BillingType equal ExactDay")
		}
	}

	if len(s.Items) == 0 {
		validator.SetError("Items", "Items is required")
	}

	for _, it := range s.Items {
		if it.PricingScheme == nil || it.PricingScheme.Price <= 0 {
			validator.SetError("PricingScheme", "Item PricingScheme must be bigger than zero")
		}
	}

	if s.Customer == nil && len(s.CustomerId) == 0 {
		validator.SetError("Customer", "Customer or CustomerId is required")
	}

	if len(s.CardId) == 0 && len(s.CardToken) == 0 {

		if s.Card == nil {
			validator.SetError("Card", "Card is required")
		} else {

			this.EntityValidator.AddEntity(s.Card)
			this.EntityValidator.AddEntity(s.Card.BillingAddress)
			this.EntityValidator.AddValidationForType(
				reflect.TypeOf(s.Card),
				cardValidator(ValidateCardCreate))
			
		}
		
		
	}
	
	if s.Card != nil {
		
	}
}
