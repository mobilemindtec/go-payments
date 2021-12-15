package v4


import (
  "github.com/mobilemindtec/go-utils/beego/validator" 
  "github.com/mobilemindtec/go-utils/support" 
  "crypto/hmac"
  "crypto/sha256" 
  //"encoding/base64"
  "encoding/json"
  "encoding/hex"
  _ "crypto/sha1"
  _ "crypto/hmac" 
  _ "strings"
  "net/url"
  "strings"
  "errors"
  _ "fmt"
)




type WebhookData struct {
	Answer *Answer `json:"kr-answer" valid:"Required"`
  Response *PayZenResult
 	KrAnswerType string `json:"kr-answer-type" valid:"Required"`
  KrAnswer string  
 	KrHash string `json:"kr-hash" valid:"Required"`
 	KrHashAlgorithm string `json:"kr-hash-algorithm" valid:"Required"`
 	KrHashKey string `json:"kr-hash-key" valid:"Required"`
  Raw string `json:"row" valid:"Required"`
}

func NewWebhookData() *WebhookData {
  return &WebhookData{ Answer: new(Answer) }
}

type Webhook struct {
  JsonParser *support.JsonParser
  Debug bool

  EntityValidator *validator.EntityValidator  
  ValidationErrors map[string]string
  HasValidationError bool
}

func NewWebhook(lang string, sallerToken string) *Webhook {
  entityValidator := validator.NewEntityValidator(lang, "PayZen")
  return &Webhook{ 
    JsonParser:  new(support.JsonParser), 
    EntityValidator: entityValidator, 
  }
}

func NewDefaultWebhook() *Webhook {
  entityValidator := validator.NewEntityValidator("pt-BR", "PayZen")
  return &Webhook{ 
    JsonParser:  new(support.JsonParser), 
    EntityValidator: entityValidator, 
  }
}

func (this *Webhook) SetDebug() {
  this.Debug = true
}

func (this *Webhook) Parse(body []byte) (*WebhookData, error) {


  formData := make(map[string]interface{})
  urlQuery, err := url.QueryUnescape(string(body))

  if err != nil {
    return nil, err
  }

  data := NewWebhookData()  

  splited := strings.Split(urlQuery, "&")


  for _, value := range splited {

    vals := strings.Split(value, "=")

    switch vals[0] {
      case "kr-hash-key":
        data.KrHashKey = vals[1]
        break
      case "kr-answer-type":
        data.KrAnswerType = vals[1]
        break
      case "kr-hash-algorithm":
        data.KrHashAlgorithm = vals[1]
        break
      case "kr-hash":
        data.KrHash = vals[1]
        break
      case "kr-answer":
        if err := json.Unmarshal([]byte(vals[1]), data.Answer); err != nil {
          return nil, err
        }
        data.KrAnswer = vals[1]
        formData[vals[0]] = data.Answer
        break
      default:
        formData[vals[0]]  = vals[1]
        break
    }


  }

  jsonString, err := json.Marshal(formData)

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

  response := new(PaymentResponse)
  response.Answer = data.Answer
  data.Response = NewPayZenResultWithResponse(response)
  return data, nil  
}

func GenerateSignatureFromBody(cert string, krAnswer string) string{
  
  mac := hmac.New(sha256.New, []byte(cert))
  mac.Write([]byte(krAnswer))
  raw := mac.Sum(nil)

  base64Content := hex.EncodeToString(raw)
  //return string(raw)

  return base64Content
  
}