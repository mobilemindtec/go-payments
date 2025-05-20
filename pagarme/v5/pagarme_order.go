package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/v2/either"
	"reflect"
)

type SuccessOrder = *Success[OrderPtr]
type SuccessOrders = *Success[Orders]

type PagarmeOrder struct {
	Pagarme
}

func NewPagarmeOrder(lang string, auth *Authentication, serviceRefererName ServiceRefererName) *PagarmeOrder {
	p := &PagarmeOrder{}
	p.Pagarme.init(lang, auth, serviceRefererName)
	return p
}

func (this *PagarmeOrder) Create(order OrderPtr) *either.Either[*ErrorResponse, SuccessOrder] {

	if !this.onValidOrder(order) {
		return either.Left[*ErrorResponse, SuccessOrder](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.validationsToMapOfStringSlice()))
	}

	return either.
		MapIf(
			this.post("/orders", order, createParser[Order]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessOrder {
				return NewSuccess[OrderPtr](e.UnwrapRight())
			})
}

func (this *PagarmeOrder) Get(orderId string) *either.Either[*ErrorResponse, SuccessOrder] {

	if empty, left := checkEmpty[SuccessOrder]("order id", orderId); empty {
		return left
	}

	uri := fmt.Sprintf("/orders/%v", orderId)

	return either.
		MapIf(
			this.get(uri, createParser[Order]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessOrder {
				return NewSuccess[OrderPtr](e.UnwrapRight())
			})
}

func (this *PagarmeOrder) List(query *OrderQuery) *either.Either[*ErrorResponse, SuccessOrders] {

	uri := fmt.Sprintf("/orders/?%v", query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, createParserContent[Orders]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessOrders {
				return NewSuccessSlice[Orders](e.UnwrapRight())
			})
}

func (this *Pagarme) onValidOrder(order *Order) bool {

	this.EntityValidator.AddValidationForType(
		reflect.TypeOf(order), func(entity interface{}, validator *validator.Validation) {
			p := entity.(*Order)

			if p.Payments == nil || len(p.Payments) == 0 {
				validator.SetError("Payments", "Payments array is required")
			}


			if len(p.Items) == 0 {
				validator.SetError("Items", "Items is required")
			}

			for _, it := range p.Items {
				if it.Amount <= 0 {
					validator.SetError("Item", "Item Amount must be bigger than zero")
				}
			}

		})

	
	this.EntityValidator.AddEntity(order)

	if order.Payments != nil {
		for _, it := range order.Payments {

			this.EntityValidator.AddEntity(it)

			this.EntityValidator.AddValidationForType(
				reflect.TypeOf(it), func(entity interface{}, validator *validator.Validation) {

					p := entity.(*Payment)

					if p.Amount <= 0 {
						validator.SetError("Amount", "Amount is required")
					}

					switch p.PaymentMethod {
					case MethodCreditCard:
						if p.CreditCard == nil {
							validator.SetError("Payment", "CreditCard object is required")
						}
					case MethodBoleto:
						if p.Boleto == nil {
							validator.SetError("Payment", "Boleto object is required")
						}
					case MethodPix:
						if p.Pix == nil {
							validator.SetError("Payment", "Pix object is required")
						}
					default:
						validator.SetError("Payment", "PaymentMethod is required")
					}
				})

			switch it.PaymentMethod {
			case MethodCreditCard:
				if it.CreditCard != nil {

					this.EntityValidator.AddEntity(it.CreditCard)

					if len(it.CreditCard.CardId) == 0 && len(it.CreditCard.CardToken) == 0 {


						if it.CreditCard.Card != nil {

							this.EntityValidator.AddEntity(it.CreditCard.Card)
							this.EntityValidator.AddEntity(it.CreditCard.Card.BillingAddress)
							this.EntityValidator.AddValidationForType(
								reflect.TypeOf(it.CreditCard.Card),
								cardValidator(true, false))
						}


						this.EntityValidator.AddValidationForType(
							reflect.TypeOf(it.CreditCard), func(entity interface{}, validator *validator.Validation) {

								if it.CreditCard.Card == nil {
									validator.SetError("Card", "Card is required")
								}
							})
					}



				}
			case MethodBoleto:
				if it.Boleto != nil {

					this.EntityValidator.AddEntity(it.Boleto)

					this.EntityValidator.AddValidationForType(
						reflect.TypeOf(it.Boleto), func(entity interface{}, validator *validator.Validation) {

						})
				}
			case MethodPix:
				if it.Pix != nil {

					this.EntityValidator.AddEntity(it.Pix)

					this.EntityValidator.AddValidationForType(
						reflect.TypeOf(it.Pix), func(entity interface{}, validator *validator.Validation) {

						})
				}
			}
		}
	}

	return this.processValidator()
}
