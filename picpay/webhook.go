package picpay

import (
  "github.com/mobilemindtec/go-utils/beego/validator" 
  beego "github.com/beego/beego/v2/server/web"
  "github.com/mobilemindtec/go-utils/support"  
  "errors"
)


type WebhookData struct {
  ReferenceId string `json:"referenceId" valid:"Required"`
  AuthorizationId string `json:"authorizationId"`
  Raw string `json:"row" valid:"Required"`
  Uuid string `json:"uuid" valid:"Required"`
}

func NewWebhookData() *WebhookData {
  return &WebhookData{}
}

type Webhook struct {
  Controller *beego.Controller  
  JsonParser *support.JsonParser
  SallerToken string
  Debug bool

  EntityValidator *validator.EntityValidator  
  ValidationErrors map[string]string
  HasValidationError bool
}

func NewWebhook(lang string, sallerToken string, controller *beego.Controller) *Webhook {
  entityValidator := validator.NewEntityValidator(lang, "PicPay")
  return &Webhook{ 
    SallerToken: sallerToken, 
    JsonParser:  new(support.JsonParser), 
    EntityValidator: entityValidator, 
  }
}

func (this *Webhook) SetDebug() {
  this.Debug = true
}

func (this *Webhook) IsValid() bool {
  token := this.Controller.Ctx.Request.Header.Get("x-seller-token")
  return this.SallerToken == token
}

func (this *Webhook) GetData() (*WebhookData, error) {
  body := this.Controller.Ctx.Input.RequestBody
  return this.Parse(body)
}

func (this *Webhook) Parse(body []byte) (*WebhookData, error) {

	jsonMap, err := this.JsonParser.JsonToMap(this.Controller.Ctx)

	if err != nil {
		return nil, err
	}

  data := NewWebhookData()

  data.ReferenceId = this.JsonParser.GetJsonString(jsonMap, "referenceId")
  data.AuthorizationId = this.JsonParser.GetJsonString(jsonMap, "authorizationId")
  data.Raw = string(body)
  data.Uuid = this.Controller.Ctx.Input.Param(":uuid")

  entityValidatorResult, _ := this.EntityValidator.IsValid(data, nil)  

  if entityValidatorResult.HasError {
    this.HasValidationError = true
    this.ValidationErrors = this.EntityValidator.GetValidationErrors(entityValidatorResult)
    return nil, errors.New("Validation error")
  }

  return data, nil  
}