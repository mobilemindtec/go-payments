package pagarme

import (
  "github.com/mobilemindtec/go-utils/beego/validator" 
  "github.com/mobilemindtec/go-utils/support"  
  "encoding/json"
  "encoding/hex"
  "crypto/sha1"
  "crypto/hmac" 
  "strings"
  _ "net/url"
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
  Fingerprint string `json:"fingerprint" valid:"Required"`
  Event WebhookEvent  `json:"event" valid:"Required"`
  OldStatus string `json:"old_status" valid:"Required"`
  DesiredStatus string `json:"desired_status" valid:"Required"`
  CurrentStatus string `json:"current_status" valid:"Required"`
  Object ObjectType `json:"object" valid:"Required"`
  Raw string `json:"raw" valid:"Required"`
  Response *Response
  Payload string
  PayloadMap map[string]interface{}
  Signature string
}

func NewWebhookData() *WebhookData {
  return &WebhookData{}
}

type Webhook struct {
  JsonParser *support.JsonParser
  Debug bool

  EntityValidator *validator.EntityValidator  
  ValidationErrors map[string]string
  HasValidationError bool
}

func NewWebhook(lang string) *Webhook {
  entityValidator := validator.NewEntityValidator(lang, "Pagarme")
  return &Webhook{ 
    JsonParser:  new(support.JsonParser), 
    EntityValidator: entityValidator, 
  }
}

func NewDefaultWebhook() *Webhook {
  entityValidator := validator.NewEntityValidator("pt-BR", "Pagarme")
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

  if this.Debug {
    fmt.Println("************************************************")
    fmt.Println("**** Pagarme.Webhook: ", jsonMap)
    fmt.Println("************************************************")
  }
  
  data.Object = ObjectType(this.JsonParser.GetJsonString(jsonMap, "object"))
  data.Payload = this.JsonParser.GetJsonString(jsonMap, "payload")
  
  data.PayloadMap = make(map[string]interface{})
  splited := strings.Split(data.Payload, "&")

  for _, it := range splited {
    sp := strings.Split(it, "=")
    if len(sp) == 2 {
      //fmt.Println(sp)
      data.PayloadMap[sp[0]] = sp[1]
    }
  }

  data.Id = this.JsonParser.GetJsonString(data.PayloadMap, "id")
  data.Fingerprint = this.JsonParser.GetJsonString(data.PayloadMap, "fingerprint")
  data.Event = WebhookEvent(this.JsonParser.GetJsonString(data.PayloadMap, "event"))
  data.OldStatus = this.JsonParser.GetJsonString(data.PayloadMap, "old_status")
  data.DesiredStatus = this.JsonParser.GetJsonString(data.PayloadMap, "desired_status")
  data.CurrentStatus = this.JsonParser.GetJsonString(data.PayloadMap, "current_status")
  data.Signature = this.JsonParser.GetJsonString(data.PayloadMap, "signature")

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

  data.Response = new(Response)
  data.Response.Object = this.JsonParser.GetJsonString(data.PayloadMap, "object")
  data.Response.Id = this.JsonParser.GetJsonInt64(data.PayloadMap, "id")
  data.Response.StatusText = data.CurrentStatus
  data.Response.OldStatusText = data.OldStatus
  data.Response.DesiredStatusText = data.DesiredStatus

  //query, err := url.ParseQuery(payload)
  //if err != nil {
  //  return nil, err
  //}

  //trasactionMap := make(map[string]interface{})
  //for k, v := range query {
  //  if strings.Contains(k, "transaction[") {
  //    kk := strings.Replace(strings.Replace(k, "transaction[", "", -1), "]", "", -1)
  //    trasactionMap[kk] = v[0]
  //    //fmt.Println(kk, " = ", v[0])
  //  }
  //}



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
