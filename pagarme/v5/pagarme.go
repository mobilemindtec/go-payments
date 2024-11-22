package v5

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mobilemindtec/go-utils/v2/either"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/beego/i18n"

	"github.com/beego/beego/v2/core/logs"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/v2/optional"
)

const (
	PAGARME_URL = "https://api.pagar.me/core/v5"
)

type ResponseParser func(data []byte, response *Response) error

type Authentication struct {
	SecretKey string
	PublicKey string
}

func NewAuthentication(secretKey string, publicKey string) *Authentication {
	return &Authentication{secretKey, publicKey}
}

func (this *Authentication) Basic() string {
	val := fmt.Sprintf("%v:", this.SecretKey)
	return base64.StdEncoding.EncodeToString([]byte(val))
}

type Pagarme struct {
	Auth                  *Authentication
	Lang                  string
	EntityValidator       *validator.EntityValidator
	EntityValidatorResult *validator.EntityValidatorResult
	ValidationErrors      map[string]string
	HasValidationError    bool
	Debug                 bool
}

func (this *Pagarme) DebugOn()  *Pagarme {
	this.SetDebug(true)
	return this
}

func (this *Pagarme) SetDebug(debug bool) {
	this.Debug = debug
}

func NewPagarme(lang string, auth *Authentication) *Pagarme {
	return (&Pagarme{}).init(lang, auth)
}

func (this *Pagarme) init(lang string, auth *Authentication) *Pagarme {
	this.Lang = lang
	this.Auth = auth
	this.EntityValidator = validator.NewEntityValidator(this.Lang, "Pagarme")
	this.EntityValidatorResult = new(validator.EntityValidatorResult)
	this.EntityValidatorResult.Errors = map[string]string{}
	return this
}

func (this *Pagarme) validationsToMapOfStringSlice() map[string][]string{
	results := make(map[string][]string)
	for k, v := range this.ValidationErrors{
		results[k] = []string{v}
	}

	return results
}

// HTTP
func (this *Pagarme) get(
	action string, parsers ...ResponseParser) *either.Either[error, *Response] {
	return this.request("GET", action, nil, tryParser(parsers))
}

func (this *Pagarme) delete(action string, payloads ...interface{}) *either.Either[error, *Response] {
	return this.request("DELETE", action, tryPayload(payloads), nil)
}

func (this *Pagarme) patch(action string, payload interface{}, parsers ...ResponseParser) *either.Either[error, *Response] {
	return this.request("PATCH", action, payload, tryParser(parsers))
}

func (this *Pagarme) post(
	action string, data interface{}, parsers ...ResponseParser) *either.Either[error, *Response] {
	return this.request("POST", action, data, tryParser(parsers))
}

func (this *Pagarme) put(
	action string, data interface{}, parsers ...ResponseParser) *either.Either[error, *Response] {
	return this.request("PUT", action, data, tryParser(parsers))
}

func (this *Pagarme) request(
	method string, action string, data interface{}, parser ResponseParser) *either.Either[error, *Response] {

	response := NewResponse()

	var req *http.Request
	var err error

	client := new(http.Client)
	apiUrl := fmt.Sprintf("%v%v", PAGARME_URL, action)

	logs.Debug("URL %v, METHOD = %v", apiUrl, method)

	if data != nil {

		payload, err := json.Marshal(data)

		if err != nil {
			logs.Debug("error json.Marshal ", err.Error())
			return either.Left[error, *Response](err)
		}

		postData := bytes.NewBuffer(payload)

		response.RawRequest = string(payload)

		if this.Debug {
			logs.Debug("****************** Pagarme Request ******************")
			logs.Debug(response.RawRequest)
			logs.Debug("****************** Pagarme Request ******************")
		}

		req, err = http.NewRequest(method, apiUrl, postData)

	} else {
		req, err = http.NewRequest(method, apiUrl, nil)
	}

	if err != nil {
		logs.Debug("http.NewRequest err = ", err)
		return either.Left[error, *Response](err)
	}

	if this.Debug {
		logs.Debug("Secret %v", this.Auth.SecretKey)
		logs.Debug("Authorization: Basic %v", this.Auth.Basic())
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", this.Auth.Basic()))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("client.Do err = %v", err)
		return either.Left[error, *Response](err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("ioutil.ReadAll err = %v", err)
		return either.Left[error, *Response](err)
	}

	response.RawResponse = string(body)
	response.StatusCode = res.StatusCode

	if this.Debug {
		fmt.Println("****************** Pagarme Response ******************")
		fmt.Println("STATUS CODE ", res.StatusCode)
		fmt.Println(response.RawResponse)
		fmt.Println("****************** Pagarme Response ******************")
	}

	switch res.StatusCode {

	case 200:

		if parser != nil {
			if err := parser(body, response); err != nil {
				fmt.Println("parser err = %v", err)
				return either.Left[error, *Response](err)
			}
		} else {
			err = json.Unmarshal(body, response)
			if err != nil {
				fmt.Println("json.Unmarshal err = %v", err)
				return either.Left[error, *Response](err)
			}
		}

		return either.Right[error, *Response](response)

	default:


		err = json.Unmarshal(body, response.Error)
		if err != nil {
			return either.Left[error, *Response](
				errors.New(fmt.Sprintf("Pagarme error. Status: %v", res.StatusCode)))
		} else {
			response.Error.StatusCode = int64(res.StatusCode)
			return either.Left[error, *Response](response.Error)
		}

	}
}

// UTILS

func (this *Pagarme) onValidCustomer(customer CustomerPtr) bool {

	this.EntityValidator.AddValidationForType(
		reflect.TypeOf(customer), func(entity interface{}, validator *validator.Validation) {
			//p := entity.(CustomerPtr)

		})

	this.EntityValidator.AddEntity(customer)

	if customer.Address == nil {
		this.EntityValidator.AddEntity(customer.Address)
	}

	return this.processValidator()
}

func (this *Pagarme) onValidEntity(entity interface{}) bool {
	this.EntityValidatorResult, _ = this.EntityValidator.IsValid(entity, nil)

	if this.EntityValidatorResult.HasError {
		this.onValidationErrors()
		return false
	}

	return true
}

func (this *Pagarme) getMessage(key string, args ...interface{}) string {
	return i18n.Tr(this.Lang, key, args)
}

func (this *Pagarme) onValidationErrors() {
	this.HasValidationError = true
	this.ValidationErrors = this.EntityValidator.GetValidationErrors(this.EntityValidatorResult)
}

func (this *Pagarme) SetValidationError(key string, value string) {
	this.HasValidationError = true
	if this.ValidationErrors == nil {
		this.ValidationErrors = make(map[string]string)
	}
	this.ValidationErrors[key] = value
}

func (this *Pagarme) processValidator() bool {
	val := this.EntityValidator.Validate()

	switch val.(type) {
	case *optional.Fail:
		fail := val.(*optional.Fail).Item
		if errs, ok := fail.(map[string]string); ok {
			this.ValidationErrors = errs
			this.HasValidationError = true
		}
		return false
	default:
		return true
	}
}

func (this *Pagarme) Log(message string, args ...interface{}) {
	if this.Debug {
		logs.Debug("Pagarme: ", fmt.Sprintf(message, args...))
	}
}

func (this *Pagarme) urlQuery(filter map[string]string) string {
	url := ""
	if filter != nil && len(filter) > 0 {
		url = fmt.Sprintf("%v?", url)

		for k, v := range filter {
			url = fmt.Sprintf("%v%v=%v", url, k, v)
			url = fmt.Sprintf("%v&", url)
		}
	}

	return url
}

func tryPayload(payloads []interface{}) interface{} {
	var payload interface{}
	if len(payloads) > 0 {
		payload = payloads[0]
	}
	return payload
}
func tryParser(parsers []ResponseParser) ResponseParser {
	var parser ResponseParser
	if len(parsers) > 0 {
		parser = parsers[0]
	}
	return parser
}

func createParser[T any]() func([]byte, *Response) error {
	return func(data []byte, response *Response) error {
		response.Content = new(T)
		return json.Unmarshal(data, response.Content)
	}
}

func unwrapError(err error) *ErrorResponse {

	if er, ok := err.(*ErrorResponse); ok {
		return er
	}

	return NewErrorResponse(fmt.Sprintf("%v", err.Error()))
}

func createParserContent[T any]() func([]byte, *Response) error {
	return func(data []byte, response *Response) error {
		response.Content = new(Content[T])
		return json.Unmarshal(data, response.Content)
	}
}

func checkEmpty[T any](label string, values ...string) (bool, *either.Either[*ErrorResponse, T]) {
	for _, value := range values {
		if len(value) == 0 {
			return true, either.Left[*ErrorResponse, T](
				NewErrorResponse(fmt.Sprintf("%v is required", label)))
		}
	}
	return false, nil
}
