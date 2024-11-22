package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/v2/either"
	"github.com/mobilemindtec/go-utils/validator/cnpj"
	"github.com/mobilemindtec/go-utils/validator/cpf"
	"reflect"
)

type SuccessCard = *Success[CardPtr]
type SuccessCards = *Success[Cards]

type SuccessBool = *Success[bool]

type PagarmeCard struct {
	Pagarme
}

func NewPagarmeCard(lang string, auth *Authentication) *PagarmeCard {
	p := &PagarmeCard{}
	p.Pagarme.init(lang, auth)
	return p
}

func (this *PagarmeCard) Create(customerId string, card CardPtr) *either.Either[*ErrorResponse, SuccessCard] {

	if empty, left := checkEmpty[SuccessCard](customerId, "customer id"); empty {
		return left
	}

	if !this.validate(card) {
		return either.Left[*ErrorResponse, SuccessCard](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.validationsToMapOfStringSlice()))
	}

	uri := fmt.Sprintf("/customers/%v/cards", customerId)

	return either.
		MapIf(
			this.post(uri, card, createParser[Card]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCard {
				return NewSuccess[CardPtr](e.UnwrapRight())
			})
}

func (this *PagarmeCard) Get(customerId string, cardId string) *either.Either[*ErrorResponse, SuccessCard] {

	if empty, left := checkEmpty[SuccessCard]("customer id and card id", customerId, cardId); empty {
		return left
	}

	uri := fmt.Sprintf("/customers/%v/cards/%v", customerId, cardId)

	return either.
		MapIf(
			this.get(uri, createParser[Card]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCard {
				return NewSuccess[CardPtr](e.UnwrapRight())
			})
}

func (this *PagarmeCard) List(customerId string) *either.Either[*ErrorResponse, SuccessCards] {

	if empty, left := checkEmpty[SuccessCards]("customer id", customerId); empty {
		return left
	}

	uri := fmt.Sprintf("/customers/%v/cards", customerId)

	return either.
		MapIf(
			this.get(uri, createParserContent[Cards]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCards {
				return NewSuccessSlice[Cards](e.UnwrapRight())
			})
}

func (this *PagarmeCard) Update(customerId string, card CardPtr) *either.Either[*ErrorResponse, SuccessCard] {

	if empty, left := checkEmpty[SuccessCard]("customer id and card id", customerId, card.Id); empty {
		return left
	}

	if !this.validate(card) {
		return either.Left[*ErrorResponse, SuccessCard](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.validationsToMapOfStringSlice()))
	}

	uri := fmt.Sprintf("/customers/%v/cards/%v", customerId, card.Id)

	return either.
		MapIf(
			this.put(uri, card, createParser[Card]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessCard {
				return NewSuccess[CardPtr](e.UnwrapRight())
			})
}

func (this *PagarmeCard) Delete(customerId string, cardId string) *either.Either[*ErrorResponse, SuccessBool] {

	if empty, left := checkEmpty[SuccessBool]("customer id and card id", customerId, cardId); empty {
		return left
	}

	uri := fmt.Sprintf("/customers/%v/cards/%v", customerId, cardId)

	return either.
		MapIf(
			this.delete(uri),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessBool {
				return NewSuccessWithValue[bool](e.UnwrapRight(), true)
			})
}

func (this *PagarmeCard) Renew(customerId string, cardId string) *either.Either[*ErrorResponse, SuccessBool] {

	if empty, left := checkEmpty[SuccessBool]("customer id and card id", customerId, cardId); empty {
		return left
	}

	uri := fmt.Sprintf("/customers/%v/cards/%v/renew", customerId, cardId)

	return either.
		MapIf(
			this.post(uri, nil),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessBool {
				return NewSuccessWithValue[bool](e.UnwrapRight(), true)
			})
}

func (this *PagarmeCard) validate(card *Card) bool {
	this.EntityValidator.AddEntity(card)
	this.EntityValidator.AddValidationForType(reflect.TypeOf(card), cardValidator)
	return this.processValidator()
}

func cardValidator(entity interface{}, validator *validator.Validation) {
	c := entity.(*Card)

	isUpdate := len(c.CardId) > 0
	isCreate := len(c.CardToken) == 0 && len(c.CardId) == 0

	if len(c.Number) == 0 && len(c.CardToken) == 0 && len(c.CardId) == 0 {
		validator.SetError("Card", "CardToken, CardId or Number is required")
	}

	if isUpdate || isCreate {
		if len(c.Brand) == 0 {
			validator.SetError("Brand", "Brand is required")
		}
		if !cpf.Validate(c.HolderDocument) && !cnpj.Validate(c.HolderDocument) {
			validator.SetError("HolderDocument", "HolderDocument is required CPF or CNPJ")
		}
		if len(c.HolderName) == 0 {
			validator.SetError("HolderName", "HolderName is required")
		}
		if len(c.Number) < 13 || len(c.Number) > 19 {
			validator.SetError("Number", "Number is required size between 13 and 19")
		}
		if len(c.Cvv) < 3 || len(c.Cvv) > 4 {
			validator.SetError("Cvv", "CVV is required size between 3 and 4")
		}
		if c.ExpMonth < 1 || c.ExpMonth > 12 {
			validator.SetError("ExpMonth", "ExpMonth is required size between 1 and 12")
		}
		if c.ExpYear < 1900 {
			validator.SetError("Brand", "ExpYear is required size greater 1900")
		}
	}

	if len(c.BillingAddressId) == 0 && c.BillingAddress == nil {
		validator.SetError("Card", "BillingAddress or BillingAddressId is required")
	}

}
