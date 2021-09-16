package pagarme

import (
  "github.com/mobilemindtec/go-utils/beego/validator" 
  beego "github.com/beego/beego/v2/server/web"
  "github.com/mobilemindtec/go-utils/support"  
  "encoding/json"
  "encoding/hex"
  "crypto/sha1"
  "crypto/hmac" 
  "strings"
  "errors"
  "fmt"

)

type WebhookEvent string

const (
  EventTransactionStatusChanged WebhookEvent = "transaction_status_changed"
  EventSubscriptionStatusChanged WebhookEvent = "subscription_status_changed"
  EventRecipientStatusChanged WebhookEvent = "recipient_status_changed"
  EventTransactionCreated WebhookEvent = "transaction_created"
)

type ObjectType string

const (
  ObjectTransaction ObjectType = "transaction"
  ObjectSubscription ObjectType = "subscription"
  ObjectRecipient ObjectType = "recipient"
)

type WebhookData struct {
  Id string `json:"id" valid:"Required"`
  Uuid string  `json:"uuid" valid:"Required"`
  Fingerprint string `json:"fingerprint" valid:"Required"`
  Event WebhookEvent  `json:"event" valid:"Required"`
  OldStatus string `json:"old_status" valid:"Required"`
  DesiredStatus string `json:"desired_status" valid:"Required"`
  CurrentStatus string `json:"current_status" valid:"Required"`
  Object ObjectType `json:"object" valid:"Required"`
  Raw string `json:"raw" valid:"Required"`
}

func NewWebhookData() *WebhookData {
  return &WebhookData{}
}

type Webhook struct {
  Controller *beego.Controller  
  JsonParser *support.JsonParser
  ApiKey string
  Debug bool

  EntityValidator *validator.EntityValidator  
  ValidationErrors map[string]string
  HasValidationError bool
}

func NewWebhook(lang string, apiKey string, controller *beego.Controller) *Webhook {
  entityValidator := validator.NewEntityValidator(lang, "Pagarme")
  return &Webhook{ 
    ApiKey: apiKey, 
    JsonParser:  new(support.JsonParser), 
    EntityValidator: entityValidator, 
  }
}

func (this *Webhook) SetDebug() {
  this.Debug = true
}

func (this *Webhook) IsValid() bool {
  signature := this.Controller.Ctx.Request.Header.Get("X-Hub-Signature")
  return CheckPostbackSignature(this.ApiKey, signature, this.Controller.Ctx.Input.RequestBody)
}

func (this *Webhook) GetData() (*WebhookData, error) {
  body := this.Controller.Ctx.Input.RequestBody
  return this.Parse(body)
}

func (this *Webhook) Parse(body []byte) (*WebhookData, error) {

  jsonMap := this.JsonParser.FormToJson(this.Controller.Ctx)
  data := NewWebhookData()

  if this.Debug {
    fmt.Println("************************************************")
    fmt.Println("**** Pagarme.Webhook: ", jsonMap)
    fmt.Println("************************************************")
  }

  data.Id = this.Controller.GetString("id")
  data.Fingerprint = this.Controller.GetString("fingerprint")
  data.Event = WebhookEvent(this.Controller.GetString("event"))
  data.OldStatus = this.Controller.GetString("old_status")
  data.DesiredStatus = this.Controller.GetString("desired_status")
  data.CurrentStatus = this.Controller.GetString("current_status")
  data.Object = ObjectType(this.Controller.GetString("object"))
  data.Uuid = this.Controller.Ctx.Input.Param(":uuid")

  jsonString, err := json.Marshal(jsonMap)

  if err != nil {
    return nil, err
  }

  data.Raw = string(jsonString)

  entityValidatorResult, _ := this.EntityValidator.IsValid(data, nil)  

  if entityValidatorResult.HasError {
    this.HasValidationError = true
    this.ValidationErrors = this.EntityValidator.GetValidationErrors(entityValidatorResult)
    return nil, errors.New("Validation error")
  }

  return data, nil

}

func CheckPostbackSignature(apiKey string, hubSignature string, requestBody []byte) bool {

  pagarmeSignature := hubSignature

  if !strings.Contains(pagarmeSignature, "="){
    fmt.Println("************************************************")
    fmt.Println("** Pagarme Signature not has =")
    fmt.Println("************************************************")        
    return false
  }

  finalSignature := strings.Split(pagarmeSignature, "=")[1]

  mac := hmac.New(sha1.New, []byte(apiKey))
  mac.Write(requestBody)
  rawBodyMAC := mac.Sum(nil)
  computedSignature := hex.EncodeToString(rawBodyMAC)


  if !hmac.Equal([]byte(finalSignature), []byte(computedSignature)){
    fmt.Println("************************************************")
    fmt.Println("**  X-Hub-Signature = ", pagarmeSignature)
    fmt.Println("**  final signature = ", finalSignature)
    fmt.Println("**  computed signature = ", computedSignature)
    fmt.Println("** Inválid Pagarme Signature: Expected: ", string(finalSignature), " Received: %v", string(computedSignature))
    fmt.Println("************************************************")    
    return false
  }
    

  //if this.Debug {
  //  fmt.Println("************************************************")
  //  fmt.Println("**  X-Hub-Signature = ", pagarmeSignature)
  //  fmt.Println("**  final signature = ", finalSignature)
  //  fmt.Println("**  computed signature = ", computedSignature)    
  //  fmt.Println("************************************************")
  //}

  return true
}
