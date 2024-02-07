package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/v2/either"
)

type PagarmeInvoice struct {
	Pagarme
}

func NewPagarmeInvoice(lang string, auth *Authentication) *PagarmeInvoice {
	p := &PagarmeInvoice{}
	p.Pagarme.init(lang, auth)
	return p
}

func (this *PagarmeInvoice) Get(id string) *either.Either[*ErrorResponse, InvoicePtr] {

	if empty, left := checkEmpty[InvoicePtr]("invoice id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/customers/%v/invoices/%v", id)

	return either.
		MapIf(
			this.get(uri, createParser[Invoice]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) InvoicePtr {
				return e.UnwrapRight().Content.(InvoicePtr)
			})
}

func (this *PagarmeInvoice) List(query *InvoiceQuery) *either.Either[*ErrorResponse, Invoices] {
	uri := fmt.Sprintf("/invoices?%v", query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, createParserContent[Invoices]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) Invoices {
				return e.UnwrapRight().Content.(*Content[Invoices]).Data
			})
}

func (this *PagarmeInvoice) Cancel(id string) *either.Either[*ErrorResponse, bool] {

	if len(id) == 0 {
		return either.Left[*ErrorResponse, bool](
			NewErrorResponse("invoice id is required"))
	}

	uri := fmt.Sprintf("/invoices/%v", id)

	return either.
		MapIf(
			this.delete(uri),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) bool {
				return true
			})
}
