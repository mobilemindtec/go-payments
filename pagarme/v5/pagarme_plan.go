package v5

import (
	"fmt"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/lists"
	"github.com/mobilemindtec/go-utils/v2/either"
	"reflect"
)

type SuccessPlan = *Success[PlanPtr]
type SuccessPlans = *Success[Plans]

type PagarmePlan struct {
	Pagarme
}

func NewPagarmePlan(lang string, auth *Authentication, serviceRefererName ServiceRefererName) *PagarmePlan {
	p := &PagarmePlan{}
	p.Pagarme.init(lang, auth, serviceRefererName)
	return p
}

func (this *PagarmePlan) Create(plan *Plan) *either.Either[*ErrorResponse, SuccessPlan] {

	if !this.validate(plan) {
		return either.Left[*ErrorResponse, SuccessPlan](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.validationsToMapOfStringSlice()))
	}

	return either.
		MapIf(
			this.post("/plans", plan, createParser[Plan]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessPlan {
				return NewSuccess[PlanPtr](e.UnwrapRight())
			})
}

func (this *PagarmePlan) Get(id string) *either.Either[*ErrorResponse, SuccessPlan] {

	if empty, left := checkEmpty[SuccessPlan]("plan id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/plans/%v", id)

	return either.
		MapIf(
			this.get(uri, createParser[Plan]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessPlan {
				return NewSuccess[PlanPtr](e.UnwrapRight())
			})
}

func (this *PagarmePlan) List(query *PlanQuery) *either.Either[*ErrorResponse, SuccessPlans] {

	uri := fmt.Sprintf("/plans/?%v", query.UrlQuery())

	return either.
		MapIf(
			this.get(uri, createParserContent[Plan]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessPlans {
				return NewSuccessSlice[Plans](e.UnwrapRight())
			})
}

func (this *PagarmePlan) Update(plan PlanPtr) *either.Either[*ErrorResponse, SuccessPlan] {

	if empty, left := checkEmpty[SuccessPlan]("plan id", plan.Id); empty {
		return left
	}

	if !this.validate(plan) {
		return either.Left[*ErrorResponse, SuccessPlan](
			NewErrorResponseWithErrors(this.getMessage("Pagarme.ValidationError"), this.validationsToMapOfStringSlice()))
	}

	uri := fmt.Sprintf("/plans/%v", plan.Id)

	return either.
		MapIf(
			this.put(uri, plan, createParser[Plan]()),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessPlan {
				return NewSuccess[PlanPtr](e.UnwrapRight())
			})
}

func (this *PagarmePlan) Delete(id string) *either.Either[*ErrorResponse, SuccessBool] {

	if empty, left := checkEmpty[SuccessBool]("plan id", id); empty {
		return left
	}

	uri := fmt.Sprintf("/plans/%v", id)

	return either.
		MapIf(
			this.delete(uri, nil),
			func(e *either.Either[error, *Response]) *ErrorResponse {
				return unwrapError(e.UnwrapLeft())
			},
			func(e *either.Either[error, *Response]) SuccessBool {
				return NewSuccessWithValue(e.UnwrapRight(), true)
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
