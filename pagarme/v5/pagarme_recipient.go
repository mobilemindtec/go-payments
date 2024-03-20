package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/v2/either"
	"github.com/mobilemindtec/go-utils/v2/maps"
	"reflect"
)

type SuccessRecipient = *Success[RecipientPtr]
type SuccessRecipients = *Success[Recipients]

type SuccessBalance = *Success[BalancePtr]
type SuccessBalanceOperations = *Success[BalanceOperations]

type SuccessTransfer = *Success[TransferPtr]
type SuccessTransfers = *Success[Transfers]

type PagarmeRecipient struct {
	Pagarme
}

func NewPagarmeRecipient(lang string, auth *Authentication) *PagarmeRecipient {
	p := &PagarmeRecipient{}
	p.Pagarme.init(lang, auth)
	return p
}

func (this *PagarmeRecipient) Create(recipient *Recipient) *either.Either[*ErrorResponse, SuccessRecipient] {

	if !this.validate(recipient) {
		return either.Left[*ErrorResponse, SuccessRecipient](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.ValidationErrors))
	}

	return either.
		MapIf(
			this.post("/recipients", recipient, createParser[Recipient]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SuccessRecipient {
				return NewSuccess[RecipientPtr](e.UnwrapRight())
			})
}

func (this *PagarmeRecipient) Update(id string, recipient *RecipientUpdate) *either.Either[*ErrorResponse, SuccessRecipient] {

	if empty, left := checkEmpty[SuccessRecipient]("recipiente id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v", id)

	return either.
		MapIf(
			this.put(uri, recipient, createParser[Recipient]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SuccessRecipient {
				return NewSuccess[RecipientPtr](e.UnwrapRight())
			})
}

func (this *PagarmeRecipient) Gt(id string) *either.Either[*ErrorResponse, SuccessRecipient] {

	if empty, left := checkEmpty[SuccessRecipient]("recipiente id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v", id)

	return either.
		MapIf(
			this.get(uri, createParser[Recipient]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SuccessRecipient {
				return NewSuccess[RecipientPtr](e.UnwrapRight())
			})
}

func (this *PagarmeRecipient) List() *either.Either[*ErrorResponse, SuccessRecipients] {

	return either.
		MapIf(
			this.get("/recipients", createParserContent[SuccessRecipients]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SuccessRecipients {
				return NewSuccessSlice[Recipients](e.UnwrapRight())
			})
}

func (this *PagarmeRecipient) UpdateTransferSettings(id string, settings *TransferSettings) *either.Either[*ErrorResponse, SuccessBool] {

	if empty, left := checkEmpty[SuccessBool]("recipiente id", id); empty {
		return left
	}

	switch settings.TransferInterval {
	case Daily:
		if settings.TransferDay <= 0 {
			return either.Left[*ErrorResponse, SuccessBool](
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
			func(e *either.Either[error, *Response]) SuccessBool {
				return NewSuccessWithValue(e.UnwrapRight(), true)
			})
}

func (this *PagarmeRecipient) UpdateBankAccount(id string, account BankAccount) *either.Either[*ErrorResponse, SuccessBool] {

	if empty, left := checkEmpty[SuccessBool]("recipiente id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v/default-bank-account", id)

	return either.
		MapIf(
			this.patch(uri, account),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SuccessBool {
				return NewSuccessWithValue(e.UnwrapRight(), true)
			})
}

func (this *PagarmeRecipient) Balance(recipientId string) *either.Either[*ErrorResponse, SuccessBalance] {

	if empty, left := checkEmpty[SuccessBalance]("recipiente id", recipientId); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v/balance", recipientId)

	return either.
		MapIf(
			this.get(uri, createParser[Balance]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SuccessBalance {
				return NewSuccess[BalancePtr](e.UnwrapRight())
			})
}

func (this *PagarmeRecipient) BalanceOperations(recipientId string, query *BalanceQuery) *either.Either[*ErrorResponse, SuccessBalanceOperations] {

	if empty, left := checkEmpty[SuccessBalanceOperations]("recipiente id", recipientId); empty {
		return left
	}

	uri := fmt.Sprintf("/balance/operations?%b", query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, createParserContent[BalanceOperations]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SuccessBalanceOperations {
				return NewSuccessSlice[BalanceOperations](e.UnwrapRight())
			})
}

func (this *PagarmeRecipient) CreateTransfer(recipientId string, amount int64) *either.Either[*ErrorResponse, SuccessTransfer] {

	if empty, left := checkEmpty[SuccessTransfer]("recipiente id", recipientId); empty {
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
			func(e *either.Either[error, *Response]) SuccessTransfer {
				return NewSuccess[TransferPtr](e.UnwrapRight())
			})
}

func (this *PagarmeRecipient) ListTransfers(recipientId string, query *TransferQuery) *either.Either[*ErrorResponse, SuccessTransfers] {

	if empty, left := checkEmpty[SuccessTransfers]("recipiente id", recipientId); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v/withdrawals?%v", recipientId, query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, createParserContent[Transfers]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SuccessTransfers {
				return NewSuccessSlice[Transfers](e.UnwrapRight())
			})
}

func (this *PagarmeRecipient) GetTransfer(recipientId string, transferId string) *either.Either[*ErrorResponse, SuccessTransfer] {

	if empty, left := checkEmpty[SuccessTransfer]("recipiente id and transfer id", recipientId, transferId); empty {
		return left
	}

	uri := fmt.Sprintf("/recipients/%v/withdrawals/%v", recipientId, transferId)

	return either.
		MapIf(
			this.get(uri, createParser[Transfer]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) SuccessTransfer {
				return NewSuccess[TransferPtr](e.UnwrapRight())
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
