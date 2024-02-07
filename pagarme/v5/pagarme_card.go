package v5

import (
	"encoding/json"
	"fmt"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/v2/either"
	"github.com/mobilemindtec/go-utils/validator/cnpj"
	"github.com/mobilemindtec/go-utils/validator/cpf"
	"reflect"
)

type PagarmeCard struct {
	Pagarme
}

func NewPagarmeCard(lang string, auth *Authentication) *PagarmeCard {
	p := &PagarmeCard{}
	p.Pagarme.init(lang, auth)
	return p
}

func (this *PagarmeCard) Create(customerId string, card CardPtr) *either.Either[*ErrorResponse, CardPtr] {

	if len(customerId) == 0 {
		return either.Left[*ErrorResponse, CardPtr](
			NewErrorResponse("customer id is required"))
	}

	if !this.validate(card) {
		return either.Left[*ErrorResponse, CardPtr](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.ValidationErrors))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Card)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/customers/%v/cards", customerId)

	return either.
		MapIf(
			this.post(card, uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) CardPtr {
				return e.UnwrapRight().Content.(CardPtr)
			})
}

func (this *PagarmeCard) Get(customerId string, cardId string) *either.Either[*ErrorResponse, CardPtr] {

	if len(customerId) == 0 || len(cardId) == 0 {
		return either.Left[*ErrorResponse, CardPtr](
			NewErrorResponse("customer id and card id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Card)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/customers/%v/cards/%v", customerId, cardId)

	return either.
		MapIf(
			this.get(uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) CardPtr {
				return e.UnwrapRight().Content.(CardPtr)
			})
}

func (this *PagarmeCard) List(customerId string) *either.Either[*ErrorResponse, Cards] {

	if len(customerId) == 0 {
		return either.Left[*ErrorResponse, Cards](
			NewErrorResponse("customer id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Content[Cards])
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/customers/%v/cards", customerId)

	return either.
		MapIf(
			this.get(uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) Cards {
				return e.UnwrapRight().Content.(*Content[Cards]).Data
			})
}

func (this *PagarmeCard) Update(customerId string, card CardPtr) *either.Either[*ErrorResponse, CardPtr] {

	if len(customerId) == 0 {
		return either.Left[*ErrorResponse, CardPtr](
			NewErrorResponse("customer id is required"))
	}

	if len(card.Id) == 0 {
		return either.Left[*ErrorResponse, CardPtr](
			NewErrorResponse("card id is required"))
	}

	if !this.validate(card) {
		return either.Left[*ErrorResponse, CardPtr](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.ValidationErrors))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Card)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/customers/%v/cards/%v", customerId, card.Id)

	return either.
		MapIf(
			this.put(card, uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) CardPtr {
				return e.UnwrapRight().Content.(CardPtr)
			})
}

func (this *PagarmeCard) Delete(customerId string, cardId string) *either.Either[*ErrorResponse, bool] {

	if len(customerId) == 0 || len(cardId) == 0 {
		return either.Left[*ErrorResponse, bool](
			NewErrorResponse("customer id and card id is required"))
	}

	uri := fmt.Sprintf("/customers/%v/cards/%v", customerId, cardId)

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

func (this *PagarmeCard) Renew(customerId string, cardId string) *either.Either[*ErrorResponse, bool] {

	if len(customerId) == 0 || len(cardId) == 0 {
		return either.Left[*ErrorResponse, bool](
			NewErrorResponse("customer id and card id is required"))
	}

	uri := fmt.Sprintf("/customers/%v/cards/%v/renew", customerId, cardId)

	return either.
		MapIf(
			this.post(nil, uri),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) bool {
				return true
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
