package v5

import (
	"encoding/json"
	"fmt"
	"github.com/mobilemindtec/go-utils/v2/either"
	"github.com/mobilemindtec/go-utils/v2/maps"
	"time"
)

type PagarmeCharge struct {
	Pagarme
}

func NewPagarmeCharge(lang string, auth *Authentication) *PagarmeCharge {
	p := &PagarmeCharge{}
	p.Pagarme.init(lang, auth)
	return p
}

func (this *PagarmeCharge) Capture(id string, code string) *either.Either[*ErrorResponse, ChargePtr] {

	if len(id) == 0 {
		return either.Left[*ErrorResponse, ChargePtr](
			NewErrorResponse("charge id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Charge)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/charges/%v/capture", id)

	payload := maps.JSON("code", code)

	return either.
		MapIf(
			this.post(payload, uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) ChargePtr {
				return e.UnwrapRight().Content.(ChargePtr)
			})
}

func (this *PagarmeCharge) Get(id string) *either.Either[*ErrorResponse, ChargePtr] {

	if len(id) == 0 {
		return either.Left[*ErrorResponse, ChargePtr](
			NewErrorResponse("charge id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Charge)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/charges/%v", id)

	return either.
		MapIf(
			this.get(uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) ChargePtr {
				return e.UnwrapRight().Content.(ChargePtr)
			})
}

func (this *PagarmeCharge) UpdateCard(id string, updateData ChargeUpdate) *either.Either[*ErrorResponse, ChargePtr] {

	if len(id) == 0 {
		return either.Left[*ErrorResponse, ChargePtr](
			NewErrorResponse("charge id is required"))
	}

	if (len(updateData.CardId) == 0 && len(updateData.CardToken) == 0 && updateData.Card == nil) ||
		(len(updateData.CardId) > 0 && len(updateData.CardToken) > 0 && updateData.Card != nil) {
		return either.Left[*ErrorResponse, ChargePtr](
			NewErrorResponse("card id, card token or card is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Charge)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/charges/%v/card", id)

	return either.
		MapIf(
			this.patch(uri, updateData, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) ChargePtr {
				return e.UnwrapRight().Content.(ChargePtr)
			})
}

func (this *PagarmeCharge) UpdateDueDate(id string, dueDate time.Time) *either.Either[*ErrorResponse, ChargePtr] {

	if len(id) == 0 {
		return either.Left[*ErrorResponse, ChargePtr](
			NewErrorResponse("charge id is required"))
	}

	if dueDate.IsZero() {
		return either.Left[*ErrorResponse, ChargePtr](
			NewErrorResponse("dueDate id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Charge)
		return json.Unmarshal(data, response.Content)
	}

	payload := maps.JSON("due_date", dueDate.Format(DateLayout))
	uri := fmt.Sprintf("/charges/%v/due-date", id)

	return either.
		MapIf(
			this.patch(uri, payload, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) ChargePtr {
				return e.UnwrapRight().Content.(ChargePtr)
			})
}

func (this *PagarmeCharge) Cancel(id string) *either.Either[*ErrorResponse, ChargePtr] {

	if len(id) == 0 {
		return either.Left[*ErrorResponse, ChargePtr](
			NewErrorResponse("charge id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Charge)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/charges/%v", id)

	return either.
		MapIf(
			this.delete(uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) ChargePtr {
				return e.UnwrapRight().Content.(ChargePtr)
			})
}

func (this *PagarmeCharge) ConfirmPayment(id string, code string, description string) *either.Either[*ErrorResponse, ChargePtr] {

	if len(id) == 0 {
		return either.Left[*ErrorResponse, ChargePtr](
			NewErrorResponse("charge id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Charge)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/charges/%v/confirm-payment", id)

	payload := maps.JSON("code", code, "description", description)

	return either.
		MapIf(
			this.post(payload, uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) ChargePtr {
				return e.UnwrapRight().Content.(ChargePtr)
			})
}

func (this *PagarmeCharge) Retry(id string) *either.Either[*ErrorResponse, ChargePtr] {

	if len(id) == 0 {
		return either.Left[*ErrorResponse, ChargePtr](
			NewErrorResponse("charge id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Charge)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/charges/%v/retry", id)

	payload := maps.JSON()

	return either.
		MapIf(
			this.post(payload, uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) ChargePtr {
				return e.UnwrapRight().Content.(ChargePtr)
			})
}

func (this *PagarmeCharge) List(query *ChargeQuery) *either.Either[*ErrorResponse, ChargePtr] {

	parser := func(data []byte, response *Response) error {
		response.Content = new(Content[Charges])
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/charges?%v", query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) ChargePtr {
				return e.UnwrapRight().Content.(ChargePtr)
			})
}
