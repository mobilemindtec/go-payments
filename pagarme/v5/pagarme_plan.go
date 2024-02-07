package v5

import (
	"encoding/json"
	"fmt"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/lists"
	"github.com/mobilemindtec/go-utils/v2/either"
	"net/url"
	"reflect"
)

type PagarmePlan struct {
	Pagarme
}

func NewPagarmePlan(lang string, auth *Authentication) *PagarmePlan {
	p := &PagarmePlan{}
	p.Pagarme.init(lang, auth)
	return p
}

func (this *PagarmePlan) Create(plan *Plan) *either.Either[*ErrorResponse, PlanPtr] {

	if !this.validate(plan) {
		return either.Left[*ErrorResponse, PlanPtr](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.ValidationErrors))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Plan)
		return json.Unmarshal(data, response.Content)
	}

	return either.
		MapIf(
			this.post(plan, "/plans", parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) PlanPtr {
				return e.UnwrapRight().Content.(PlanPtr)
			})
}

func (this *PagarmePlan) Get(id string) *either.Either[*ErrorResponse, PlanPtr] {

	if len(id) == 0 {
		return either.Left[*ErrorResponse, PlanPtr](
			NewErrorResponse("plan id is required"))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Plan)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/plans/%v", id)

	return either.
		MapIf(
			this.get(uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) PlanPtr {
				return e.UnwrapRight().Content.(PlanPtr)
			})
}

func (this *PagarmePlan) List(query *PlanQuery) *either.Either[*ErrorResponse, Plans] {

	parser := func(data []byte, response *Response) error {
		response.Content = new(Content[Plans])
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/plans/?%v", url.QueryEscape(query.UrlQuery()))

	return either.
		MapIf(
			this.get(uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) Plans {
				return e.UnwrapRight().Content.(*Content[Plans]).Data
			})
}

func (this *PagarmePlan) Update(plan PlanPtr) *either.Either[*ErrorResponse, PlanPtr] {

	if len(plan.Id) == 0 {
		return either.Left[*ErrorResponse, PlanPtr](
			NewErrorResponse("plan id is required"))
	}

	if !this.validate(plan) {
		return either.Left[*ErrorResponse, PlanPtr](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.ValidationErrors))
	}

	parser := func(data []byte, response *Response) error {
		response.Content = new(Plan)
		return json.Unmarshal(data, response.Content)
	}

	uri := fmt.Sprintf("/plans/%v", plan.Id)

	return either.
		MapIf(
			this.put(plan, uri, parser),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return NewErrorResponse(fmt.Sprintf("%v", e.UnwrapLeft()))
			},
			func(e *either.Either[error, *Response]) PlanPtr {
				return e.UnwrapRight().Content.(PlanPtr)
			})
}

func (this *PagarmePlan) Delete(id string) *either.Either[*ErrorResponse, bool] {

	if len(id) == 0 {
		return either.Left[*ErrorResponse, bool](
			NewErrorResponse("customer id and card id is required"))
	}

	uri := fmt.Sprintf("/plans/%v", id)

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

func (this *PagarmePlan) validate(plan *Plan) bool {
	this.EntityValidator.AddEntity(plan)

	items := lists.Map(plan.Items, func(i interface{}) interface{} { return i })
	this.EntityValidator.AddEntities(items...)

	if plan.PricingScheme != nil {
		this.EntityValidator.AddEntity(plan.PricingScheme)
	}

	this.EntityValidator.AddValidationForType(reflect.TypeOf(plan), planValidator)
	return this.processValidator()
}

func planValidator(entity interface{}, validator *validator.Validation) {
	c := entity.(*Plan)

	// validator.SetError("Brand", "Brand is required")
	if len(c.PaymentMethods) == 0 && c.PricingScheme == nil {
		validator.SetError("PaymentMethods", "PaymentMethods or PricingScheme is required")
	}

	if len(c.Items) == 0 {
		validator.SetError("Items", "Items is required")
	}

	if len(c.Id) > 0 {
		if len(c.Status) == 0 {
			validator.SetError("Status", "Status is required")
		}
	}
}
