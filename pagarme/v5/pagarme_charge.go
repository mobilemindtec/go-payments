package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/v2/either"
	"github.com/mobilemindtec/go-utils/v2/maps"
	"time"
)

type SuccessCharge = *Success[ChargePtr]
type SuccessCharges = *Success[Charges]

type PagarmeCharge struct {
	Pagarme
}

func NewPagarmeCharge(lang string, auth *Authentication, serviceRefererName ServiceRefererName) *PagarmeCharge {
	p := &PagarmeCharge{}
	p.Pagarme.init(lang, auth, serviceRefererName)
	return p
}

func (this *PagarmeCharge) Capture(id string, code string) *either.Either[*ErrorResponse, SuccessCharge] {

	if empty, left := checkEmpty[SuccessCharge]("charge id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/charges/%v/capture", id)

	payload := maps.JSON("code", code)

	return either.
		MapIf(
			this.post(uri, payload, createParser[Charge]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCharge {
				return NewSuccess[ChargePtr](e.UnwrapRight())
			})
}

func (this *PagarmeCharge) Get(chargeId string) *either.Either[*ErrorResponse, SuccessCharge] {

	if empty, left := checkEmpty[SuccessCharge]("charge id", chargeId); empty {
		return left
	}

	uri := fmt.Sprintf("/charges/%v", chargeId)

	return either.
		MapIf(
			this.get(uri, createParser[Charge]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCharge {
				return NewSuccess[ChargePtr](e.UnwrapRight())
			})
}

func (this *PagarmeCharge) UpdateCard(chargeId string, updateData ChargeUpdate) *either.Either[*ErrorResponse, SuccessCharge] {

	if empty, left := checkEmpty[SuccessCharge]("charge id", chargeId); empty {
		return left
	}

	if (len(updateData.CardId) == 0 && len(updateData.CardToken) == 0 && updateData.Card == nil) ||
		(len(updateData.CardId) > 0 && len(updateData.CardToken) > 0 && updateData.Card != nil) {
		return either.Left[*ErrorResponse, SuccessCharge](
			NewErrorResponse("card id, card token or card is required"))
	}

	uri := fmt.Sprintf("/charges/%v/card", chargeId)

	return either.
		MapIf(
			this.patch(uri, updateData, createParser[Charge]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCharge {
				return NewSuccess[ChargePtr](e.UnwrapRight())
			})
}

func (this *PagarmeCharge) UpdateDueDate(chargeId string, dueDate time.Time) *either.Either[*ErrorResponse, SuccessCharge] {

	if empty, left := checkEmpty[SuccessCharge]("charge id", chargeId); empty {
		return left
	}

	if dueDate.IsZero() {
		return either.Left[*ErrorResponse, SuccessCharge](
			NewErrorResponse("dueDate id is required"))
	}

	payload := maps.JSON("due_date", dueDate.Format(DateLayout))
	uri := fmt.Sprintf("/charges/%v/due-date", chargeId)

	return either.
		MapIf(
			this.patch(uri, payload, createParser[Charge]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCharge {
				return NewSuccess[ChargePtr](e.UnwrapRight())
			})
}

func (this *PagarmeCharge) Cancel(chargeId string) *either.Either[*ErrorResponse, SuccessCharge] {

	if empty, left := checkEmpty[SuccessCharge]("charge id", chargeId); empty {
		return left
	}

	uri := fmt.Sprintf("/charges/%v", chargeId)

	return either.
		MapIf(
			this.delete(uri, nil, createParser[Charge]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCharge {
				return NewSuccess[ChargePtr](e.UnwrapRight())
			})
}

func (this *PagarmeCharge) ConfirmPayment(chargeId string, opts ...*ChargePaymentConfirm) *either.Either[*ErrorResponse, SuccessCharge] {

	if empty, left := checkEmpty[SuccessCharge]("charge id", chargeId); empty {
		return left
	}

	uri := fmt.Sprintf("/charges/%v/confirm-payment", chargeId)

	var payload *ChargePaymentConfirm

	if len(opts) > 0 {
		payload = opts[0]
	}

	return either.
		MapIf(
			this.post(uri, payload, createParser[Charge]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCharge {
				return NewSuccess[ChargePtr](e.UnwrapRight())
			})
}

func (this *PagarmeCharge) Retry(chargeId string) *either.Either[*ErrorResponse, SuccessCharge] {

	if empty, left := checkEmpty[SuccessCharge]("charge id", chargeId); empty {
		return left
	}

	uri := fmt.Sprintf("/charges/%v/retry", chargeId)

	payload := maps.JSON()

	return either.
		MapIf(
			this.post(uri, payload, createParser[Charge]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCharge {
				return NewSuccess[ChargePtr](e.UnwrapRight())
			})
}

func (this *PagarmeCharge) List(query *ChargeQuery) *either.Either[*ErrorResponse, SuccessCharges] {

	uri := fmt.Sprintf("/charges?%v", query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, createParserContent[Charges]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCharges {

				return NewSuccessSlice[Charges](e.UnwrapRight())
			})
}

func (this *PagarmeCharge) ListPending(orderId string, size int) *either.Either[*ErrorResponse, SuccessCharges] {
	query := NewChargeQuery()
	query.OrderId = orderId
	query.Status = ChargePending
	query.Size = size
	return this.List(query)
}
