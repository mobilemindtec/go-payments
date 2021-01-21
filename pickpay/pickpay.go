package pickpay

import (
  "encoding/json"
  "fmt"
  "net/http"
  "io/ioutil"
  "errors"
  "bytes"
  "time"
)


/*

curl --location --request GET 'https://api.chat24.io/v1/help/transports'

*/

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
}

func (this *PickPayTransaction) HasError() bool {
  return len(this.Message) > 0
}

type PickPay struct {
  PickPayToken string
  PickPaySallerToken string
  Debug bool
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
  ExpiresAt string `json:"expiresAt" valid:"Required"`
}

func NewPickPayTransactionRequest() *PickPayTransactionRequest {
  request := new(PickPayTransactionRequest)
  request.Buyer = new(PickPayBuyer)
  request.AdditionalInfo = make(map[string]interface{})
  return request
}

func NewPickPay(token string, sallerToken string) *PickPay {
	return &PickPay{ PickPayToken: token, PickPaySallerToken: sallerToken }
}

func (this *PickPay) CreateTransaction(request *PickPayTransactionRequest) (*PickPayTransaction, error) {

  if this.Debug {
    fmt.Println("PickPay CreateTransaction")
  }

  jsonData, err := json.Marshal(request)

  if err != nil {
    fmt.Println("error json.Marshal ", err.Error())
    return nil, err
  }

  postData := bytes.NewBuffer(jsonData)

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

  if this.Debug {
    fmt.Println("****************** PickPay Response ******************")
    fmt.Println(string(body))
    fmt.Println("****************** PickPay Response ******************")
  }

  result := new(PickPayTransaction)
  err = json.Unmarshal(body, result)

  if err != nil {
    fmt.Println("err = %v", err)
    fmt.Println(string(body))
    return nil, errors.New(fmt.Sprintf("error on json.Unmarshal: %v", err))
  }

  if res.StatusCode != 200 {
    return result, errors.New(fmt.Sprintf("PickPay error. Status: %v, Detalhes: %v", res.StatusCode, result.Message)) 
  }

  if len(result.Message) > 0 {
    return result, errors.New(fmt.Sprintf("%v", result.Message))  
  }

  return result, nil
}

