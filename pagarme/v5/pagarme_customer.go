package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/v2/either"
	"net/url"
)

type PagarmeCustomer struct {
	Pagarme
}

func NewPagarmeCustomer(lang string, auth *Authentication) *PagarmeCustomer {
	p := &PagarmeCustomer{}
	p.Pagarme.init(lang, auth)
	return p
}

func (this *PagarmeCustomer) Create(customer CustomerPtr) *either.Either[*ErrorResponse, CustomerPtr] {

	if !this.onValidCustomer(customer) {
		return either.Left[*ErrorResponse, CustomerPtr](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.ValidationErrors))
	}

	return either.
		MapIf(
			this.post("/customers", customer, createParser[Customer]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) CustomerPtr {
				return e.UnwrapRight().Content.(CustomerPtr)
			})
}

func (this *PagarmeCustomer) Update(customer CustomerPtr) *either.Either[*ErrorResponse, CustomerPtr] {

	if !this.onValidCustomer(customer) {
		return either.Left[*ErrorResponse, CustomerPtr](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.ValidationErrors))
	}

	if empty, left := checkEmpty[CustomerPtr]("customer id", customer.Id); empty {
		return left
	}

	uri := fmt.Sprintf("/customers/%v", customer.Id)

	return either.
		MapIf(
			this.put(uri, customer, createParser[Customer]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) CustomerPtr {
				return e.UnwrapRight().Content.(CustomerPtr)
			})
}

func (this *PagarmeCustomer) Get(customerId string) *either.Either[*ErrorResponse, CustomerPtr] {

	if empty, left := checkEmpty[CustomerPtr]("customer id", customerId); empty {
		return left
	}

	uri := fmt.Sprintf("/customers/%v", customerId)

	return either.
		MapIf(
			this.get(uri, createParser[Customer]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) CustomerPtr {
				return e.UnwrapRight().Content.(CustomerPtr)
			})
}

func (this *PagarmeCustomer) List(query *CustomerQuery) *either.Either[*ErrorResponse, Customers] {

	uri := fmt.Sprintf("/customers/?%v", url.QueryEscape(query.UrlQuery()))

	return either.
		MapIf(
			this.get(uri, createParserContent[Customers]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) Customers {
				return e.UnwrapRight().Content.(*Content[Customers]).Data
			})
}
