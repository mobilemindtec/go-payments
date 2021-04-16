package pickpay

import (
  "github.com/mobilemindtec/go-utils/beego/validator"  
  "github.com/beego/beego/v2/core/validation"
  "github.com/beego/i18n"  
  "encoding/json"
  "io/ioutil"
  "net/http"
  "errors"
  "bytes"
  "time"
  "fmt"
)


/*

curl --location --request GET 'https://api.chat24.io/v1/help/transports'

*/
/*
  "created": registro criado
  "expired": prazo para pagamento expirado
  "analysis": pago e em processo de análise anti-fraude
  "paid": pago
  "completed": pago e saldo disponível
  "refunded": pago e devolvido
  "chargeback": pago e com chargeback
*/ 

type PickPayStatus int64

const (
  PickPayCreated PickPayStatus = 1 + iota
  PickPayExpired
  PickPayAnalysis
  PickPayPaid
  PickPayCompleted
  PickPayRefunded
  PickPayChargeback
)

const (
  PickPayApiUrl = "https://appws.picpay.com"
)

type PickPayQrCode struct {
  Content string `json:"content"`
  Base64 string `json:"base64"`
}

type PickPayTransaction struct {
  ReferenceId string `json:"referenceId"`
  PaymentUrl string `json:"paymentUrl"`
  ExpiresAt time.Time `json:"expiresAt"`
  QrCode *PickPayQrCode `json:"qrcode"`
  Message string `json:"message"`
  StatusText string `json:"status"`
  PickPayStatus PickPayStatus `json:"pickpayStatus"`
  AuthorizationId string `json:"authorizationId"`
  CancellationId string `json:"cancellationId"`
}

type PickPayResult struct {
  Transaction *PickPayTransaction `json:"transaction"`
  ValidationErrors map[string]string `json:"error"`
  Error bool `json:"error"`  
  Message string `json:"message"`
  Response string
  Request string
}

func (this *PickPayTransaction) HasError() bool {
  return len(this.Message) > 0
}

type PickPay struct {
  PickPayToken string
  PickPaySallerToken string
  Debug bool
  EntityValidator *validator.EntityValidator
  EntityValidatorResult *validator.EntityValidatorResult  
  Lang string  
  ValidationErrors map[string]string
  HasValidationError bool  
}

type PickPayBuyer struct {
  FirstName string `json:"firstName" valid:"Required"`
  LastName string `json:"lastName" valid:"Required"`
  Document string `json:"document" valid:"Required"`
  Email string `json:"email" valid:"Required"`
  Phone string `json:"phone" valid:"Required"`    
}

type PickPayTransactionRequest struct {
  ReferenceId string `json:"referenceId" valid:"Required"`
  CallbackUrl string `json:"callbackUrl" valid:"Required"`
  ReturnUrl string `json:"returnUrl" valid:"Required"`
  Value string `json:"value" valid:"Required"`
  Plugin string `json:"plugin" valid:""`
  AdditionalInfo map[string]interface{} `json:"additionalInfo" valid:"Required"`
  Buyer *PickPayBuyer `json:"buyer" valid:"Required"`
  ExpiresAt time.Time `json:"" valid:"Required"`
  ExpiresAtFormatted string `json:"expiresAt" valid:""`
  AuthorizationId string `json:"authorizationId,omitempty" valid:""`
}

func NewPickPayTransactionRequest() *PickPayTransactionRequest {
  request := new(PickPayTransactionRequest)
  request.Buyer = new(PickPayBuyer)
  request.AdditionalInfo = make(map[string]interface{})
  return request
}

func NewPickPay(lang string, token string, sallerToken string) *PickPay {
  entityValidator := validator.NewEntityValidator(lang, "PickPay")
	return &PickPay{ EntityValidator: entityValidator, PickPayToken: token, PickPaySallerToken: sallerToken }
}

func (this *PickPay) CreateTransaction(request *PickPayTransactionRequest) (*PickPayResult, error) {

  if this.Debug {
    fmt.Println("PickPay CreateTransaction")
  }

  result := new(PickPayResult)

  this.EntityValidatorResult = new(validator.EntityValidatorResult)
  this.EntityValidatorResult.Errors = map[string]string{}

  if !this.onValid(request) {
    result.Error = true
    result.ValidationErrors = this.EntityValidatorResult.Errors      
    return result, errors.New(this.getMessage("PickPay.ValidationError"))
  }

  if !request.ExpiresAt.IsZero() {
    request.ExpiresAtFormatted = request.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z")
  }

  jsonData, err := json.Marshal(request)

  if err != nil {
    fmt.Println("error json.Marshal ", err.Error())    
    return result, err
  }

  postData := bytes.NewBuffer(jsonData)

  result.Request = string(jsonData)

	method := "POST"

  client := &http.Client {}
  req, err := http.NewRequest(method, fmt.Sprintf("%v/ecommerce/public/payments", PickPayApiUrl), postData)

  if err != nil {
    fmt.Println("err = %v", err)
  	return nil, errors.New(fmt.Sprintf("error on http.NewRequest: %v", err))
  }

  req.Header.Add("x-picpay-token", this.PickPayToken)  
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
    fmt.Println("****************** PickPay Response ******************")
    fmt.Println(result.Response)
    fmt.Println("****************** PickPay Response ******************")
  }

  tresult := new(PickPayTransaction)
  err = json.Unmarshal(body, tresult)



  if err != nil {
    fmt.Println("err = %v", err)
    fmt.Println(result.Response)    
    return nil, errors.New(fmt.Sprintf("error on json.Unmarshal: %v", err))
  }

  if res.StatusCode != 200 {
    result.Error = true
    result.Message = fmt.Sprintf("PickPay error. Status: %v, Detalhes: %v", res.StatusCode, tresult.Message)
    return result, errors.New(result.Message) 
  }

  if len(tresult.Message) > 0 {
    result.Error = true
    result.Message = fmt.Sprintf("%v", tresult.Message)
    return result, errors.New(result.Message)  
  }


  tresult.StatusText = "created"
  tresult.PickPayStatus = PickPayCreated

  result.Transaction = tresult

  return result, nil
}

func (this *PickPay) CheckStatus(referenceId string) (*PickPayResult, error) {

  if this.Debug {
    fmt.Println("PickPay CheckStatus")
  }

  result := new(PickPayResult)

  method := "GET"

  client := &http.Client {}
  req, err := http.NewRequest(method, fmt.Sprintf("%v/ecommerce/public/payments/%v/status", PickPayApiUrl, referenceId), nil)

  if err != nil {
    fmt.Println("err = %v", err)
    return nil, errors.New(fmt.Sprintf("error on http.NewRequest: %v", err))
  }

  req.Header.Add("x-picpay-token", this.PickPayToken)  
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
    fmt.Println("****************** PickPay Response ******************")
    fmt.Println(result.Response)
    fmt.Println("****************** PickPay Response ******************")
  }

  tresult := new(PickPayTransaction)
  err = json.Unmarshal(body, tresult)



  if err != nil {
    fmt.Println("err = %v", err)
    fmt.Println(result.Response)    
    return nil, errors.New(fmt.Sprintf("error on json.Unmarshal: %v", err))
  }

  if res.StatusCode != 200 {
    result.Error = true
    result.Message = fmt.Sprintf("PickPay error. Status: %v, Detalhes: %v", res.StatusCode, tresult.Message)
    return result, errors.New(result.Message) 
  }

  if len(tresult.Message) > 0 {
    result.Error = true
    result.Message = fmt.Sprintf("%v", tresult.Message)
    return result, errors.New(result.Message)  
  }

  switch tresult.StatusText {
    case "created":
      tresult.PickPayStatus = PickPayCreated
      break
    case "expired":
      tresult.PickPayStatus = PickPayExpired
      break
    case "analysis":
      tresult.PickPayStatus = PickPayAnalysis
      break
    case "paid":
      tresult.PickPayStatus = PickPayPaid
      break
    case "completed":
      tresult.PickPayStatus = PickPayCompleted
      break
    case "refunded":
      tresult.PickPayStatus = PickPayRefunded
      break
    case "chargeback":
      tresult.PickPayStatus = PickPayChargeback
      break    
    default:
      fmt.Println("PickPay: status %v not found", tresult.StatusText)
  }

  result.Transaction = tresult

  return result, nil
}

func (this *PickPay) Cancel(referenceId string, authorizationId string) (*PickPayResult, error) {

  if this.Debug {
    fmt.Println("PickPay Cancel")
  }

  result := new(PickPayResult)

  payload := map[string]string{}

  if len(authorizationId) > 0 {
    payload["authorizationId"] = authorizationId
  } 

  jsonData, err := json.Marshal(payload)

  if err != nil {
    fmt.Println("error json.Marshal ", err.Error())    
    return result, err
  }

  postData := bytes.NewBuffer(jsonData)

  method := "POST"

  client := &http.Client {}
  req, err := http.NewRequest(method, fmt.Sprintf("%v/ecommerce/public/payments/%v/cancellations", PickPayApiUrl, referenceId), postData)

  if err != nil {
    fmt.Println("err = %v", err)
    return nil, errors.New(fmt.Sprintf("error on http.NewRequest: %v", err))
  }

  req.Header.Add("x-picpay-token", this.PickPayToken)  
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
    fmt.Println("****************** PickPay Response ******************")
    fmt.Println(result.Response)
    fmt.Println("****************** PickPay Response ******************")
  }

  tresult := new(PickPayTransaction)
  err = json.Unmarshal(body, tresult)



  if err != nil {
    fmt.Println("err = %v", err)
    fmt.Println(result.Response)    
    return nil, errors.New(fmt.Sprintf("error on json.Unmarshal: %v", err))
  }

  if res.StatusCode != 200 {
    result.Error = true
    result.Message = fmt.Sprintf("PickPay error. Status: %v, Detalhes: %v", res.StatusCode, tresult.Message)
    return result, errors.New(result.Message) 
  }

  if len(tresult.Message) > 0 {
    result.Error = true
    result.Message = fmt.Sprintf("%v", tresult.Message)
    return result, errors.New(result.Message)  
  }  

  result.Transaction = tresult

  return result, nil
}

func (this *PickPay) onValid(request *PickPayTransactionRequest) bool {
  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(request, func (validator *validation.Validation) {
    
 
  })
 
   if this.EntityValidatorResult.HasError {
    this.onValidationErrors()
    return false
  }

  return true
}

func (this *PickPay) getMessage(key string, args ...interface{}) string{
  return i18n.Tr(this.Lang, key, args)
}

func (this *PickPay) onValidationErrors(){
  this.HasValidationError = true
  data := make(map[interface{}]interface{})
  this.EntityValidator.CopyErrorsToView(this.EntityValidatorResult, data)
  this.ValidationErrors = data["errors"].(map[string]string)
}
