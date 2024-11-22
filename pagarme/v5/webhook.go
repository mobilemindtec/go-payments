package v5

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/mobilemindtec/go-utils/support"
	"strings"
)

type WebhookData struct {
	Id         string       `json:"id" valid:"Required"`
	Event      WebhookEvent `json:"event" valid:"Required"`
	Status     string       `json:"status" valid:""`
	Raw        string       `json:"raw" valid:"Required"`
	PayloadMap map[string]interface{}
}

func NewWebhookData() *WebhookData {
	return &WebhookData{}
}

func (this *WebhookData) IsOrder() bool {
	return strings.HasPrefix(string(this.Event), "order.")
}

func (this *WebhookData) IsSubscription() bool {
	return strings.HasPrefix(string(this.Event), "subscription.")
}

func (this *WebhookData) IsCharge() bool {
	return strings.HasPrefix(string(this.Event), "charge.")
}

type Webhook struct {
	JsonParser *support.JsonParser
	Debug      bool

	EntityValidator    *validator.EntityValidator
	ValidationErrors   map[string]string
	HasValidationError bool
}

func NewWebhook(lang string) *Webhook {
	entityValidator := validator.NewEntityValidator(lang, "Pagarme")
	return &Webhook{
		JsonParser:      new(support.JsonParser),
		EntityValidator: entityValidator,
	}
}

func NewDefaultWebhook() *Webhook {
	entityValidator := validator.NewEntityValidator("pt-BR", "Pagarme")
	return &Webhook{
		JsonParser:      new(support.JsonParser),
		EntityValidator: entityValidator,
	}
}

func (this *Webhook) SetDebug() {
	this.Debug = true
}

func ParseWebhookObject(body []byte, entity interface{}) error {
	return json.Unmarshal(body, entity)
}

func (this *Webhook) Parse(body []byte) (*WebhookData, error) {

	jsonMap, err := this.JsonParser.JsonBytesToMap(body)

	if err != nil {
		return nil, err
	}

	data := NewWebhookData()

	if this.Debug {
		fmt.Println("************************************************")
		fmt.Println("**** Pagarmev5.Webhook: ", jsonMap)
		fmt.Println("************************************************")
	}

	payload := this.JsonParser.GetJsonObject(jsonMap, "data")
	data.Id = this.JsonParser.GetJsonString(payload, "id")
	data.Status = this.JsonParser.GetJsonString(payload, "status")
	data.PayloadMap = jsonMap
	data.Event = WebhookEvent(this.JsonParser.GetJsonString(jsonMap, "type"))
	data.Event = WebhookEvent(this.JsonParser.GetJsonString(data.PayloadMap, "event"))
	data.Raw = string(body)

	entityValidatorResult, _ := this.EntityValidator.IsValid(data, nil)

	if entityValidatorResult.HasError {
		this.HasValidationError = true
		this.ValidationErrors = this.EntityValidator.GetValidationErrors(entityValidatorResult)
		return nil, errors.New("validation error")
	}

	return data, nil

}
