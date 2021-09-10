package v4


import (
  "github.com/mobilemindtec/go-utils/beego/validator" 
  beego "github.com/beego/beego/v2/server/web"
  "github.com/mobilemindtec/go-utils/support"  
  "encoding/json"
  _ "encoding/hex"
  _ "crypto/sha1"
  _ "crypto/hmac" 
  _ "strings"
  "errors"
  _ "fmt"
)




type WebhookData struct {
	Answer *Answer `json:"kr-answer" valid:"Required"`
 	KrAnswerType string `json:"kr-answer-type" valid:"Required"`
 	KrHash string `json:"kr-hash" valid:"Required"`
 	KrHashAlgorithm string `json:"kr-hash-algorithm" valid:"Required"`
 	KrHashKey string `json:"kr-hash-key" valid:"Required"`
  Raw string `json:"row" valid:"Required"`
  Uuid string `json:"uuid" valid:"Required"`
}

func NewWebhookData() *WebhookData {
  return &WebhookData{ Answer: new(Answer) }
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
  entityValidator := validator.NewEntityValidator(lang, "PayZen")
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

	jsonMap := this.JsonParser.FormToJson(this.Controller.Ctx)


  data := NewWebhookData()
  data.Uuid = this.Controller.Ctx.Input.Param(":uuid")

  jsonString, err := json.Marshal(jsonMap)

  if err != nil {
    return nil, err
  }

  data.Raw = string(jsonString)

  krAnswer := this.JsonParser.GetJsonString(jsonMap, "kr-answer")

  if len(krAnswer) == 0 {
  	return nil, errors.New("empty kr-answer")
  }

  err = json.Unmarshal([]byte(krAnswer), data.Answer)    

  if err != nil {
  	return nil, err
  }

  entityValidatorResult, _ := this.EntityValidator.IsValid(data, nil)  

  if entityValidatorResult.HasError {
    this.HasValidationError = true
    this.ValidationErrors = this.EntityValidator.GetValidationErrors(entityValidatorResult)
    return nil, errors.New("Validation error")
  }

  return data, nil  
}