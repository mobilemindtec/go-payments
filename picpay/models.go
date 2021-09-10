package picpay

import(
  "github.com/mobilemindtec/go-payments/api"
  "strings"
  "time"
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
  PicPayStatus api.PicPayStatus `json:"picpayStatus"`
  AuthorizationId string `json:"authorizationId"`
  CancellationId string `json:"cancellationId"`
}

func (this *PicPayTransaction) GetPayZenSOAPStatus() api.PayZenTransactionStatus {


  switch this.PicPayStatus {
    case api.PicPayCreated:
      return api.Created
    case api.PicPayExpired:
      return api.Expired
    case api.PicPayAnalysis:
      return api.UnderVerification
    case api.PicPayPaid:
      return api.Authorised
    case api.PicPayCompleted:
      return api.Authorised
    case api.PicPayRefunded:
      return api.Refunded
    case api.PicPayChargeback:
      return api.Chargeback
    default:

      if len(strings.TrimSpace(this.CancellationId)) > 0 {
        this.PicPayStatus = api.PicPayCancelled
        return api.Cancelled
      } else {
        return api.Error
      }
  } 
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
