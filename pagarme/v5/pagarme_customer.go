package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/v2/either"
	"net/url"
)

type SuccessCustomer = *Success[CustomerPtr]
type SuccessCustomers = *Success[Customers]

type PagarmeCustomer struct {
	Pagarme
}

func NewPagarmeCustomer(lang string, auth *Authentication) *PagarmeCustomer {
	p := &PagarmeCustomer{}
	p.Pagarme.init(lang, auth)
	return p
}

func (this *PagarmeCustomer) Create(customer CustomerPtr) *either.Either[*ErrorResponse, SuccessCustomer] {

	if !this.onValidCustomer(customer) {
		return either.Left[*ErrorResponse, SuccessCustomer](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.validationsToMapOfStringSlice()))
	}

	return either.
		MapIf(
			this.post("/customers", customer, createParser[Customer]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCustomer {
				return NewSuccess[CustomerPtr](e.UnwrapRight())
			})
}

func (this *PagarmeCustomer) Update(customer CustomerPtr) *either.Either[*ErrorResponse, SuccessCustomer] {

	if !this.onValidCustomer(customer) {
		return either.Left[*ErrorResponse, SuccessCustomer](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.validationsToMapOfStringSlice()))
	}

	if empty, left := checkEmpty[SuccessCustomer]("customer id", customer.Id); empty {
		return left
	}

	uri := fmt.Sprintf("/customers/%v", customer.Id)

	return either.
		MapIf(
			this.put(uri, customer, createParser[Customer]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCustomer {
				return NewSuccess[CustomerPtr](e.UnwrapRight())
			})
}

func (this *PagarmeCustomer) Get(customerId string) *either.Either[*ErrorResponse, SuccessCustomer] {

	if empty, left := checkEmpty[SuccessCustomer]("customer id", customerId); empty {
		return left
	}

	uri := fmt.Sprintf("/customers/%v", customerId)

	return either.
		MapIf(
			this.get(uri, createParser[Customer]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCustomer {
				return NewSuccess[CustomerPtr](e.UnwrapRight())
			})
}

func (this *PagarmeCustomer) List(query *CustomerQuery) *either.Either[*ErrorResponse, SuccessCustomers] {

	uri := fmt.Sprintf("/customers/?%v", url.QueryEscape(query.UrlQuery()))

	return either.
		MapIf(
			this.get(uri, createParserContent[SuccessCustomers]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCustomers {
				return NewSuccessSlice[Customers](e.UnwrapRight())
			})
}
