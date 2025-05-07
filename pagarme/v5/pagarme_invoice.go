package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/v2/either"
)

type SuccessInvoice = *Success[InvoicePtr]
type SuccessInvoices = *Success[Invoices]

type PagarmeInvoice struct {
	Pagarme
}

func NewPagarmeInvoice(lang string, auth *Authentication, serviceRefererName ServiceRefererName) *PagarmeInvoice {
	p := &PagarmeInvoice{}
	p.Pagarme.init(lang, auth, serviceRefererName)
	return p
}

func (this *PagarmeInvoice) Get(id string) *either.Either[*ErrorResponse, SuccessInvoice] {

	if empty, left := checkEmpty[SuccessInvoice]("invoice id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/customers/%v/invoices/%v", id)

	return either.
		MapIf(
			this.get(uri, createParser[Invoice]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessInvoice {
				return NewSuccess[InvoicePtr](e.UnwrapRight())
			})
}

func (this *PagarmeInvoice) List(query *InvoiceQuery) *either.Either[*ErrorResponse, SuccessInvoices] {
	uri := fmt.Sprintf("/invoices?%v", query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, createParserContent[Invoices]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessInvoices {
				return NewSuccessSlice[Invoices](e.UnwrapRight())
			})
}

func (this *PagarmeInvoice) Cancel(id string) *either.Either[*ErrorResponse, SuccessBool] {
	
	if len(id) == 0 {
		return either.Left[*ErrorResponse, SuccessBool](
			NewErrorResponse("invoice id is required"))
	}

	uri := fmt.Sprintf("/invoices/%v", id)

	return either.
		MapIf(
			this.delete(uri,nil),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessBool {
				return NewSuccessWithValue[bool](e.UnwrapRight(), true)
			})
}
