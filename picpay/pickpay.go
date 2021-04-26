package picpay

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

type PicPayStatus int64

const (
  PicPayCreated PicPayStatus = 1 + iota
  PicPayExpired
  PicPayAnalysis
  PicPayPaid
  PicPayCompleted
  PicPayRefunded
  PicPayChargeback
)

const (
  PicPayApiUrl = "https://appws.picpay.com"
)

type PicPayQrCode struct {
  Content string `json:"content"`
  Base64 string `json:"base64"`
}

type PicPayTransaction struct {
  ReferenceId string `json:"referenceId"`
  PaymentUrl string `json:"paymentUrl"`
  ExpiresAt time.Time `json:"expiresAt"`
  QrCode *PicPayQrCode `json:"qrcode"`
  Message string `json:"message"`
  StatusText string `json:"status"`
  PicPayStatus PicPayStatus `json:"picpayStatus"`
  AuthorizationId string `json:"authorizationId"`
  CancellationId string `json:"cancellationId"`
}

type PicPayResult struct {
  Transaction *PicPayTransaction `json:"transaction"`
  ValidationErrors map[string]string `json:"error"`
  Error bool `json:"error"`  
  Message string `json:"message"`
  Response string
  Request string
}

func (this *PicPayTransaction) HasError() bool {
  return len(this.Message) > 0
}

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

type PicPayBuyer struct {
  FirstName string `json:"firstName" valid:"Required"`
  LastName string `json:"lastName" valid:"Required"`
  Document string `json:"document" valid:"Required"`
  Email string `json:"email" valid:"Required"`
  Phone string `json:"phone" valid:"Required"`    
}

type PicPayTransactionRequest struct {
  ReferenceId string `json:"referenceId" valid:"Required"`
  CallbackUrl string `json:"callbackUrl" valid:"Required"`
  ReturnUrl string `json:"returnUrl" valid:"Required"`
  Value string `json:"value" valid:"Required"`
  Plugin string `json:"plugin" valid:""`
  AdditionalInfo map[string]interface{} `json:"additionalInfo" valid:"Required"`
  Buyer *PicPayBuyer `json:"buyer" valid:"Required"`
  ExpiresAt time.Time `json:"" valid:"Required"`
  ExpiresAtFormatted string `json:"expiresAt" valid:""`
  AuthorizationId string `json:"authorizationId,omitempty" valid:""`
}

func NewPicPayTransactionRequest() *PicPayTransactionRequest {
  request := new(PicPayTransactionRequest)
  request.Buyer = new(PicPayBuyer)
  request.AdditionalInfo = make(map[string]interface{})
  return request
}

func NewPicPay(lang string, token string, sallerToken string) *PicPay {
  entityValidator := validator.NewEntityValidator(lang, "PicPay")
	return &PicPay{ EntityValidator: entityValidator, PicPayToken: token, PicPaySallerToken: sallerToken }
}

func (this *PicPay) CreateTransaction(request *PicPayTransactionRequest) (*PicPayResult, error) {

  if this.Debug {
    fmt.Println("PicPay CreateTransaction")
  }

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

  jsonData, err := json.Marshal(request)

  if err != nil {
    fmt.Println("error json.Marshal ", err.Error())    
    return result, err
  }

  postData := bytes.NewBuffer(jsonData)

  result.Request = string(jsonData)

	method := "POST"

  client := &http.Client {}
  req, err := http.NewRequest(method, fmt.Sprintf("%v/ecommerce/public/payments", PicPayApiUrl), postData)

  if err != nil {
    fmt.Println("err = %v", err)
  	return nil, errors.New(fmt.Sprintf("error on http.NewRequest: %v", err))
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

  tresult := new(PicPayTransaction)
  err = json.Unmarshal(body, tresult)



  if err != nil {
    fmt.Println("err = %v", err)
    fmt.Println(result.Response)    
    return nil, errors.New(fmt.Sprintf("error on json.Unmarshal: %v", err))
  }

  if res.StatusCode != 200 {
    result.Error = true
    result.Message = fmt.Sprintf("PicPay error. Status: %v, Detalhes: %v", res.StatusCode, tresult.Message)
    return result, errors.New(result.Message) 
  }

  if len(tresult.Message) > 0 {
    result.Error = true
    result.Message = fmt.Sprintf("%v", tresult.Message)
    return result, errors.New(result.Message)  
  }


  tresult.StatusText = "created"
  tresult.PicPayStatus = PicPayCreated

  result.Transaction = tresult

  return result, nil
}

func (this *PicPay) CheckStatus(referenceId string) (*PicPayResult, error) {

  if this.Debug {
    fmt.Println("PicPay CheckStatus")
  }

  result := new(PicPayResult)

  method := "GET"

  client := &http.Client {}
  req, err := http.NewRequest(method, fmt.Sprintf("%v/ecommerce/public/payments/%v/status", PicPayApiUrl, referenceId), nil)

  if err != nil {
    fmt.Println("err = %v", err)
    return nil, errors.New(fmt.Sprintf("error on http.NewRequest: %v", err))
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

  tresult := new(PicPayTransaction)
  err = json.Unmarshal(body, tresult)



  if err != nil {
    fmt.Println("err = %v", err)
    fmt.Println(result.Response)    
    return nil, errors.New(fmt.Sprintf("error on json.Unmarshal: %v", err))
  }

  if res.StatusCode != 200 {
    result.Error = true
    result.Message = fmt.Sprintf("PicPay error. Status: %v, Detalhes: %v", res.StatusCode, tresult.Message)
    return result, errors.New(result.Message) 
  }

  if len(tresult.Message) > 0 {
    result.Error = true
    result.Message = fmt.Sprintf("%v", tresult.Message)
    return result, errors.New(result.Message)  
  }

  switch tresult.StatusText {
    case "created":
      tresult.PicPayStatus = PicPayCreated
      break
    case "expired":
      tresult.PicPayStatus = PicPayExpired
      break
    case "analysis":
      tresult.PicPayStatus = PicPayAnalysis
      break
    case "paid":
      tresult.PicPayStatus = PicPayPaid
      break
    case "completed":
      tresult.PicPayStatus = PicPayCompleted
      break
    case "refunded":
      tresult.PicPayStatus = PicPayRefunded
      break
    case "chargeback":
      tresult.PicPayStatus = PicPayChargeback
      break    
    default:
      fmt.Println("PicPay: status %v not found", tresult.StatusText)
  }

  result.Transaction = tresult

  return result, nil
}

func (this *PicPay) Cancel(referenceId string, authorizationId string) (*PicPayResult, error) {

  if this.Debug {
    fmt.Println("PicPay Cancel")
  }

  result := new(PicPayResult)

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
  req, err := http.NewRequest(method, fmt.Sprintf("%v/ecommerce/public/payments/%v/cancellations", PicPayApiUrl, referenceId), postData)

  if err != nil {
    fmt.Println("err = %v", err)
    return nil, errors.New(fmt.Sprintf("error on http.NewRequest: %v", err))
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

  tresult := new(PicPayTransaction)
  err = json.Unmarshal(body, tresult)



  if err != nil {
    fmt.Println("err = %v", err)
    fmt.Println(result.Response)    
    return nil, errors.New(fmt.Sprintf("error on json.Unmarshal: %v", err))
  }

  if res.StatusCode != 200 {
    result.Error = true
    result.Message = fmt.Sprintf("PicPay error. Status: %v, Detalhes: %v", res.StatusCode, tresult.Message)
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

func (this *PicPay) onValid(request *PicPayTransactionRequest) bool {
  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(request, func (validator *validation.Validation) {
    
 
  })
 
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
