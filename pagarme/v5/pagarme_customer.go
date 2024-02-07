package v5

import (
	"encoding/json"
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

	parser := func(data []byte, response *Response) error {
		response.Content = new(Customer)
		return json.Unmarshal(data, response.Content)
	}

	return either.
		MapIf(
			this.post(customer, "/customers", parser),
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

	if len(customer.Id) == 0 {
		return either.Left[*ErrorResponse, CustomerPtr](
			NewErrorResponse("id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Customer)
		return json.Unmarshal(data, response.Content)
	}

	return either.
		MapIf(
			this.put(customer, fmt.Sprintf("/customers/%v", customer.Id), parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) CustomerPtr {
				return e.UnwrapRight().Content.(CustomerPtr)
			})
}

func (this *PagarmeCustomer) Get(customerId string) *either.Either[*ErrorResponse, CustomerPtr] {

	if len(customerId) == 0 {
		return either.Left[*ErrorResponse, CustomerPtr](
			NewErrorResponse("id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Customer)
		return json.Unmarshal(data, response.Content)
	}

	return either.
		MapIf(
			this.get(fmt.Sprintf("/customers/%v", customerId), parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) CustomerPtr {
				return e.UnwrapRight().Content.(CustomerPtr)
			})
}

func (this *PagarmeCustomer) List(query *CustomerQuery) *either.Either[*ErrorResponse, Customers] {

	parser := func(data []byte, response *Response) error {
		response.Content = new(Content[Customers])
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/customers/?%v", url.QueryEscape(query.UrlQuery()))

	return either.
		MapIf(
			this.get(uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) Customers {
				return e.UnwrapRight().Content.(*Content[Customers]).Data
			})
}
