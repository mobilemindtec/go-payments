package picpay

import (
  "github.com/mobilemindtec/go-utils/beego/validator" 
  "github.com/mobilemindtec/go-utils/support"  
  "github.com/mobilemindtec/go-payments/api"  
  "errors"
)


type WebhookData struct {
  ReferenceId string `json:"referenceId" valid:"Required"`
  AuthorizationId string `json:"authorizationId"`
  Raw string `json:"row" valid:"Required"`
  Response *PicPayResult
}

func NewWebhookData() *WebhookData {
  return &WebhookData{}
}

type Webhook struct {
  JsonParser *support.JsonParser
  SallerToken string
  Debug bool

  EntityValidator *validator.EntityValidator  
  ValidationErrors map[string]string
  HasValidationError bool
}

func NewWebhook(lang string, sallerToken string) *Webhook {
  entityValidator := validator.NewEntityValidator(lang, "PicPay")
  return &Webhook{ 
    SallerToken: sallerToken, 
    JsonParser:  new(support.JsonParser), 
    EntityValidator: entityValidator, 
  }
}

func NewDefaultWebhook() *Webhook{
  entityValidator := validator.NewEntityValidator("pt-BR", "PicPay")
  return &Webhook{
    JsonParser:  new(support.JsonParser), 
    EntityValidator: entityValidator, 
  }
}

func (this *Webhook) SetDebug() {
  this.Debug = true
}

func (this *Webhook) Parse(body []byte) (*WebhookData, error) {

	jsonMap, err := this.JsonParser.JsonBytesToMap(body)

	if err != nil {
		return nil, err
	}

  data := NewWebhookData()

  data.ReferenceId = this.JsonParser.GetJsonString(jsonMap, "referenceId")
  data.AuthorizationId = this.JsonParser.GetJsonString(jsonMap, "authorizationId")
  data.Raw = string(body)

  entityValidatorResult, _ := this.EntityValidator.IsValid(data, nil)  

  if entityValidatorResult.HasError {
    this.HasValidationError = true
    this.ValidationErrors = this.EntityValidator.GetValidationErrors(entityValidatorResult)
    return nil, errors.New("Validation error")
  }

  data.Response = new(PicPayResult)
  data.Response.Response = string(body)
  data.Response.Transaction = new(PicPayTransaction)
  data.Response.Transaction.ReferenceId = data.ReferenceId
  data.Response.Transaction.AuthorizationId = data.AuthorizationId
  data.Response.Transaction.PicPayStatus = api.PicPayCreated
  data.Response.Transaction.StatusText = "new status received"

  return data, nil  
}