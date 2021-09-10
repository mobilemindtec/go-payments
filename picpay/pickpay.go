package picpay

import (
  "github.com/mobilemindtec/go-utils/beego/validator"  
  _ "github.com/beego/beego/v2/core/validation"
  "github.com/mobilemindtec/go-payments/api"
  "github.com/beego/i18n"  
  "encoding/json"
  "io/ioutil"
  "net/http"
  "errors"
  "bytes"
  "fmt"
)


/*

curl --location --request GET 'https://api.chat24.io/v1/help/transports'

*/

const (
  PicPayApiUrl = "https://appws.picpay.com"
)

type PicPay struct {
  PicPayToken string
  PicPaySallerToken string
  Debug bool
  EntityValidator *validator.EntityValidator
  EntityValidatorResult *validator.EntityValidatorResult  
  Lang string  
  ValidationErrors map[string]string
  HasValidationError bool  
}


func NewPicPay(lang string, token string, sallerToken string) *PicPay {
  entityValidator := validator.NewEntityValidator(lang, "PicPay")
	return &PicPay{ EntityValidator: entityValidator, PicPayToken: token, PicPaySallerToken: sallerToken }
}

func (this *PicPay) CreateTransaction(request *PicPayTransactionRequest) (*PicPayResult, error) {

  if this.Debug {
    fmt.Println("PicPay CreateTransaction")
  }

  var err error
  result := new(PicPayResult)

  this.EntityValidatorResult = new(validator.EntityValidatorResult)
  this.EntityValidatorResult.Errors = map[string]string{}

  if !this.onValid(request) {
    result.Error = true
    result.ValidationErrors = this.EntityValidatorResult.Errors      
    return result, errors.New(this.getMessage("PicPay.ValidationError"))
  }

  if !request.ExpiresAt.IsZero() {
    request.ExpiresAtFormatted = request.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z")
  }


  result, err = this.request(request, "ecommerce/public/payments")

  if err != nil {
    return result, err
  }

  result.Transaction.StatusText = "created"
  result.Transaction.PicPayStatus = api.PicPayCreated

  return result, nil
}

func (this *PicPay) CheckStatus(referenceId string) (*PicPayResult, error) {

  if this.Debug {
    fmt.Println("PicPay CheckStatus")
  }

  if len(referenceId) == 0 {
    this.SetValidationError("referenceId", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }

  return this.request(nil, fmt.Sprintf("ecommerce/public/payments/%v/status", referenceId))
}

func (this *PicPay) Cancel(referenceId string, authorizationId string) (*PicPayResult, error) {

  if this.Debug {
    fmt.Println("PicPay Cancel")
  }

  if len(referenceId) == 0 {
    this.SetValidationError("referenceId", "is required")
    return nil, errors.New(this.getMessage("PayZen.ValidationError"))       
  }  

  payload := map[string]string{}

  if len(authorizationId) > 0 {
    payload["authorizationId"] = authorizationId
  } 

  return this.request(payload, fmt.Sprintf("ecommerce/public/payments/%v/cancellations", referenceId))
}

func (this *PicPay) request(data interface{}, action string) (*PicPayResult, error) {

  result := new(PicPayResult)

  client := new(http.Client)
  apiUrl := fmt.Sprintf("%v/%v", PicPayApiUrl, action)

  method := "GET"
  var postData *bytes.Buffer = nil

  if data != nil {
    method = "POST"

    jsonData, err := json.Marshal(data)

    if err != nil {
      fmt.Println("error json.Marshal ", err.Error())    
      return nil, err
    }

    postData = bytes.NewBuffer(jsonData)

    result.Request = string(jsonData)

    if this.Debug {
      fmt.Println("****************** PicPay Request ******************")
      fmt.Println(result.Request)
      fmt.Println("****************** PicPay Request ******************")
    }

  } else {
    result.Request = "http get"
  }

  this.Log("URL %v, METHOD = %v", apiUrl, method)

  var req *http.Request 
  var reqError error

  if method == "GET" {
    req, reqError = http.NewRequest(method, apiUrl, nil)
  } else {
    req, reqError = http.NewRequest(method, apiUrl, postData)
  }

  if reqError != nil {
    fmt.Println("err = %v", reqError)
    return nil, errors.New(fmt.Sprintf("error on http.NewRequest: %v", reqError))
  }

  req.Header.Add("x-picpay-token", this.PicPayToken)  
  req.Header.Add("Content-Type", "application/json")

  res, err := client.Do(req)

  if err != nil {
    fmt.Println("err = %v", err)
    return nil, errors.New(fmt.Sprintf("error on client.Do: %v", err))
  }

  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)

  if err != nil {
    fmt.Println("err = %v", err)
    return nil, errors.New(fmt.Sprintf("error on ioutil.ReadAll: %v", err))
  }

  result.Response = string(body)

  if this.Debug {
    fmt.Println("****************** PicPay Response ******************")
    fmt.Println(result.Response)
    fmt.Println("****************** PicPay Response ******************")
  }

  transaction := new(PicPayTransaction)
  err = json.Unmarshal(body, transaction)



  if err != nil {
    fmt.Println("err = %v", err)
    fmt.Println(result.Response)    
    return nil, errors.New(fmt.Sprintf("error on json.Unmarshal: %v", err))
  }

  if res.StatusCode != 200 {
    result.Error = true
    result.Message = fmt.Sprintf("PicPay error. Status: %v, Detalhes: %v", res.StatusCode, transaction.Message)
    return result, errors.New(result.Message) 
  }

  if len(transaction.Message) > 0 {
    result.Error = true
    result.Message = fmt.Sprintf("%v", transaction.Message)
    return result, errors.New(result.Message)  
  }

  switch transaction.StatusText {
    case "created":
      transaction.PicPayStatus = api.PicPayCreated
      break
    case "expired":
      transaction.PicPayStatus = api.PicPayExpired
      break
    case "analysis":
      transaction.PicPayStatus = api.PicPayAnalysis
      break
    case "paid":
      transaction.PicPayStatus = api.PicPayPaid
      break
    case "completed":
      transaction.PicPayStatus = api.PicPayCompleted
      break
    case "refunded":
      transaction.PicPayStatus = api.PicPayRefunded
      break
    case "chargeback":
      transaction.PicPayStatus = api.PicPayChargeback
      break    
    default:
      fmt.Println("PicPay: status %v not found", transaction.StatusText)
  }

  result.Transaction = transaction  

  return result, nil
}

func (this *PicPay) onValid(request *PicPayTransactionRequest) bool {
  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(request, nil)
 
  if this.EntityValidatorResult.HasError {
    this.onValidationErrors()
    return false
  }

  return true
}

func (this *PicPay) getMessage(key string, args ...interface{}) string{
  return i18n.Tr(this.Lang, key, args)
}

func (this *PicPay) onValidationErrors(){
  this.HasValidationError = true
  data := make(map[interface{}]interface{})
  this.EntityValidator.CopyErrorsToView(this.EntityValidatorResult, data)
  this.ValidationErrors = data["errors"].(map[string]string)
}


func (this *PicPay) Log(message string, args ...interface{}) {
  if this.Debug {
    fmt.Println("PicPay: ", fmt.Sprintf(message, args...))
  }
}

func (this *PicPay) SetValidationError(key string, value string){
  this.HasValidationError = true
  if this.ValidationErrors == nil {
    this.ValidationErrors = make(map[string]string)
  }
  this.ValidationErrors[key]= value
}
