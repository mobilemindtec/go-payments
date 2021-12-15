package asaas

import (
  "github.com/mobilemindtec/go-utils/beego/validator"	
  "github.com/mobilemindtec/go-payments/api"	
	"encoding/json"
	"errors"
	"fmt"
)




type WebhookData struct {
	Event api.PaymentEvent `json:"event" valid:"Required"`
	Response *Response `json:"payment" valid:"Required"`
	Raw string `json:"raw" valid:"Required"`
} 

func NewWebhookData() *WebhookData{
	return &WebhookData{}
}

type Webhook struct {
	Debug bool
  EntityValidator *validator.EntityValidator  
  ValidationErrors map[string]string
  HasValidationError bool	
}

func NewWebhook(lang string) *Webhook {
	entityValidator := validator.NewEntityValidator(lang, "Asaas")
	return &Webhook{ EntityValidator: entityValidator }
}

func NewDefaultWebhook() *Webhook {
	entityValidator := validator.NewEntityValidator("pt-BR", "Asaas")
	return &Webhook{ EntityValidator: entityValidator }
}

func (this *Webhook) SetDebug()  {
	this.Debug = true
}

func (this *Webhook) Parse(body []byte) (*WebhookData, error) {
	data := NewWebhookData()

  if this.Debug {
    fmt.Println("************************************************")
    fmt.Println("**** Asaas.Webhook: ", string(body))
    fmt.Println("************************************************")
  }

	err := json.Unmarshal(body, data)	

	if data.Response != nil {
		data.Response.BuildStatus()
	}

	data.Raw = string(body)

	if err != nil {
	  entityValidatorResult, _ := this.EntityValidator.IsValid(data, nil)  

	  if entityValidatorResult.HasError {
	    this.HasValidationError = true
	    this.ValidationErrors = this.EntityValidator.GetValidationErrors(entityValidatorResult)
	    return nil, errors.New("Validation error")
	  }		
	}

	return data, err
}
