package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/v2/either"
	"github.com/mobilemindtec/go-utils/v2/maps"
	"reflect"
)

type PagarmeRecipient struct {
	Pagarme
}

func NewPagarmeRecipient(lang string, auth *Authentication) *PagarmeRecipient {
	p := &PagarmeRecipient{}
	p.Pagarme.init(lang, auth)
	return p
}

func (this *PagarmeRecipient) Create(recipient *Recipient) *either.Either[*ErrorResponse, RecipientPtr] {

	if !this.validate(recipient) {
		return either.Left[*ErrorResponse, RecipientPtr](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.ValidationErrors))
	}

	return either.
		MapIf(
			this.post("/recipients", recipient, createParser[Recipient]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) RecipientPtr {
				return e.UnwrapRight().Content.(RecipientPtr)
			})
}

func (this *PagarmeRecipient) Update(id string, recipient *RecipientUpdate) *either.Either[*ErrorResponse, RecipientPtr] {

	if empty, left := checkEmpty[RecipientPtr]("recipiente id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v", id)

	return either.
		MapIf(
			this.put(uri, recipient, createParser[Recipient]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) RecipientPtr {
				return e.UnwrapRight().Content.(RecipientPtr)
			})
}

func (this *PagarmeRecipient) Gt(id string) *either.Either[*ErrorResponse, RecipientPtr] {

	if empty, left := checkEmpty[RecipientPtr]("recipiente id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v", id)

	return either.
		MapIf(
			this.get(uri, createParser[Recipient]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) RecipientPtr {
				return e.UnwrapRight().Content.(RecipientPtr)
			})
}

func (this *PagarmeRecipient) List() *either.Either[*ErrorResponse, Recipients] {

	return either.
		MapIf(
			this.get("/recipients", createParserContent[Recipients]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) Recipients {
				return e.UnwrapRight().Content.(Recipients)
			})
}

func (this *PagarmeRecipient) UpdateTransferSettings(id string, settings *TransferSettings) *either.Either[*ErrorResponse, bool] {

	if empty, left := checkEmpty[bool]("recipiente id", id); empty {
		return left
	}

	switch settings.TransferInterval {
	case Daily:
		if settings.TransferDay <= 0 {
			return either.Left[*ErrorResponse, bool](
				NewErrorResponse("transfer day is required"))

		}
	}

	uri := fmt.Sprintf("/recipients/%v/transfer-settings", id)

	return either.
		MapIf(
			this.patch(uri, settings),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) bool {
				return true
			})
}

func (this *PagarmeRecipient) UpdateBankAccount(id string, account BankAccount) *either.Either[*ErrorResponse, bool] {

	if empty, left := checkEmpty[bool]("recipiente id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v/default-bank-account", id)

	return either.
		MapIf(
			this.patch(uri, account),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) bool {
				return true
			})
}

func (this *PagarmeRecipient) BalancePtr(recipientId string) *either.Either[*ErrorResponse, BalancePtr] {

	if empty, left := checkEmpty[BalancePtr]("recipiente id", recipientId); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v/balance", recipientId)

	return either.
		MapIf(
			this.get(uri, createParser[Balance]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) BalancePtr {
				return e.UnwrapRight().Content.(BalancePtr)
			})
}

func (this *PagarmeRecipient) CreateTransfer(recipientId string, amount int64) *either.Either[*ErrorResponse, TransferPtr] {

	if empty, left := checkEmpty[TransferPtr]("recipiente id", recipientId); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v/withdrawals", recipientId)

	payload := maps.JSON("amount", amount)

	return either.
		MapIf(
			this.post(uri, payload, createParser[Transfer]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) TransferPtr {
				return e.UnwrapRight().Content.(TransferPtr)
			})
}

func (this *PagarmeRecipient) ListTransfers(recipientId string, query *TransferQuery) *either.Either[*ErrorResponse, Transfers] {

	if empty, left := checkEmpty[Transfers]("recipiente id", recipientId); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v/withdrawals?%v", recipientId, query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, createParserContent[Transfers]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) Transfers {
				return e.UnwrapRight().Content.(*Content[Transfers]).Data
			})
}

func (this *PagarmeRecipient) GetTransfer(recipientId string, transferId string) *either.Either[*ErrorResponse, TransferPtr] {

	if empty, left := checkEmpty[TransferPtr]("recipiente id and transfer id", recipientId, transferId); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v/withdrawals/%v", recipientId, transferId)

	return either.
		MapIf(
			this.get(uri, createParser[Transfer]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) TransferPtr {
				return e.UnwrapRight().Content.(TransferPtr)
			})
}

func (this *PagarmeRecipient) validate(recipient *Recipient) bool {
	this.EntityValidator.AddEntity(recipient)
	this.EntityValidator.AddEntity(recipient.TransferSettings)
	this.EntityValidator.AddEntity(recipient.DefaultBankAccount)
	this.EntityValidator.AddValidationForType(reflect.TypeOf(recipient.TransferSettings), transferSettingsValidator)
	return this.processValidator()
}

func transferSettingsValidator(entity interface{}, validator *validator.Validation) {
	settings := entity.(*TransferSettings)

	switch settings.TransferInterval {
	case Daily:
		if settings.TransferDay <= 0 {
			validator.SetError("TransferDay", "TransferDay is required")

		}
	}
}
