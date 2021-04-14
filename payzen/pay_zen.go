package payzen

import (
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/i18n"
	"github.com/mobilemindtec/go-payments/pagarme"
	"github.com/mobilemindtec/go-payments/pickpay"
	"strconv"
	"strings"
	"errors"
	"time"
	"fmt"
)

type PayZenCard struct {
	Number string `valid:""`
	Scheme string `valid:""`
	ExpiryMonth string `valid:""`
	ExpiryYear string `valid:""`
	CardSecurityCode string `valid:""`
	CardHolderBirthDay string `valid:""`
	CardHolderName string `valid:""`
	Token string`valid:""`

	BoletoOnline BoletoOnlineTipo `valid:""`	
	BoletoOnlineDaysDalay int // Obs.: Não suportado para boletos online Bradesco.

	BoletoOnlineTexto string 
	BoletoOnlineTexto2 string
	BoletoOnlineTexto3 string

	BoletoInstructions string `valid:""` 
	BoletoExpirationDate string // 2006-01-02
}

type PayZenCustomer struct {
	FirstName string `valid:"Required"`
	LastName string `valid:""`
	PhoneNumber string `valid:""`
	Email string `valid:"Required"`
	StreetNumber string `valid:""`
	Address string `valid:""`
	District string `valid:""`
	ZipCode string `valid:""`
	City string `valid:""`
	State string `valid:""`
	Country string `valid:"Required"`
	IdentityCode string `valid:"Required"`
}

type PayZenSubscription struct{

	OrderId string `valid:"Required"`
	SubscriptionId string `valid:"Required"`
	Description string `valid:"Required"`

	// valor da recorrência
	Amount float64 `valid:"Required"`

	// valor inicial da recorrência
	InitialAmount float64 `valid:""`
	// quantas vezes o valor inicial deve ser cobrado
	InitialAmountNumber float64 `valid:""`

	// data de inicio da cobrança
	EffectDate time.Time `valid:"Required"`
	// cobrar no último dia do mês
	LastDayOfMonth bool
	// quantidade de cobranças
	Count int64 `valid:"Required;"`

	MonthDay int64

	FrequencyByDay int64

	Rule string

	Token string `valid:"Required;"`

	Account *PayZenAccount `valid:"Required"`
}

type PayZenPayment struct{
	OrderId string `valid:""`
	Installments int `valid:""`
	Amount float64 `valid:""`
	Customer *PayZenCustomer `valid:""`
	Card *PayZenCard `valid:""`
	Account *PayZenAccount `valid:"Required"`
	TokenOperation bool
	TransactionUuid string
	VadsTransId string
	ValidationType PayZenPaymentValidationType	

	// pagarme
	PostbackUrl string `valid:""` 
	SoftDescriptor string `valid:""` 	
	Metadata map[string]string `valid:""` 
	SaveBoletoAtPath string

	//PickPay
  CallbackUrl string `json:""`
  ReturnUrl string `json:""`
  Plugin string `json:""`
  AdditionalInfo map[string]interface{} `json:""`
  ExpiresAt time.Time `json:""`
}

type PayZenCapturePayment struct {
	TransactionUuids string `valid:"Required"`
	Commission float64
	Account *PayZenAccount `valid:"Required"`
}

type PayZenPaymentToken struct {
	Token string `valid:"Required"`
	Account *PayZenAccount `valid:"Required"`
}

type PayZenPaymentFind struct {
	TransactionId string
	TransactionUuid string
	OrderId string
	SubscriptionId string
	Token string
	AuthorizationId string
	Account *PayZenAccount `valid:"Required"`
}

func NewPayZenSubscription(shopId string, mode string, cert string) *PayZenSubscription {
	subscription := new(PayZenSubscription)
	subscription.Account = &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
	return subscription
}


func NewPayZenPaymentFind(shopId string, mode string, cert string) *PayZenPaymentFind {
	find := new(PayZenPaymentFind)
	find.Account = &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
	return find
}

func NewEmptyPayZenPaymentFind() *PayZenPaymentFind {
	find := new(PayZenPaymentFind)
	find.Account = &PayZenAccount{ }
	return find
}

func NewPayZenPaymentToken(shopId string, mode string, cert string) *PayZenPaymentToken {
	tokenPayment := new(PayZenPaymentToken)
	tokenPayment.Account = &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
	return tokenPayment
}

func NewPayZenCapturePayment(shopId string, mode string, cert string) *PayZenCapturePayment {
	capturePayment := new(PayZenCapturePayment)
	capturePayment.Account = &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
	return capturePayment
}

func NewPayZenPayment(shopId string, mode string, cert string) *PayZenPayment {
	payment := new(PayZenPayment)
	payment.Account = &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
	payment.Customer = new(PayZenCustomer)
	payment.Card = new(PayZenCard)
	payment.ValidationType = Automatica
	return payment
}

func NewEmptyPayZenPayment() *PayZenPayment {
	payment := new(PayZenPayment)
	payment.Customer = new(PayZenCustomer)
	payment.Card = new(PayZenCard)
	return payment
}

func NewPayZenAccount(shopId string, mode string, cert string) *PayZenAccount {
	return &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
}

type PayZenTokenResult struct {
	Token string
	Number string
	Brand string
	CreationDate time.Time
	CancellationDate time.Time
	Cancelled bool
	Active bool
	NotFound bool
}

type PayZenSubscriptionResult struct {
	SubscriptionId string

	PastPaymentsNumber int64
	TotalPaymentsNumber int64
	EffectDate time.Time
	CancelDate time.Time
	InitialAmountNumber int64
	Rule string
	Description string

	Active bool
	Cancelled bool
	Started bool

	Token string
}


type PayZenTransactionItemResult struct {
	TransactionUuid string
	TransactionId string
	TransactionStatus PayZenTransactionStatus
	PagarmeStatus pagarme.PagarmeStatus
	PickPayStatus pickpay.PickPayStatus
	TransactionStatusLabel string
	ExternalTransactionId string
	Amount float64
	ExpectedCaptureDate time.Time
	CreationDate time.Time
}

type PayZenResult struct {
	Error bool
	Message string
	Request string
	Response string

	//RequestObject interface{}
	//ResponseObject interface{}

	PagarmeStatus pagarme.PagarmeStatus
	PickPayStatus pickpay.PickPayStatus
	TransactionStatus PayZenTransactionStatus

	TransactionStatusLabel string

	TransactionId string
	TransactionUuid string
	//ExternalTransactionId string

  ResponseCode string
  ResponseCodeDetail string

	BoletoUrl string
	BoletoNumber string

	SubscriptionId string

	TokenInfo *PayZenTokenResult

	PaymentNotFound bool
	SubscriptionInvalid bool
	SubscriptionIdNotFound bool
	PaymentTokenNotFound bool

	Transactions []*PayZenTransactionItemResult

	SubscriptionInfo *PayZenSubscriptionResult

	ValidationErrors map[string]string

	Platform string
	Nsu string	
	BoletoOutputContent []byte	`json:"-"`
	BoletoFileName string `json:"-"`

	QrCode string
	QrCodeUrl string
	PaymentUrl string
  AuthorizationId string
  CancellationId string
}

func NewPayZenResult() *PayZenResult {
	result := new(PayZenResult)
	result.TokenInfo = new(PayZenTokenResult)
	result.Transactions = []*PayZenTransactionItemResult{}
	result.SubscriptionInfo = new(PayZenSubscriptionResult)
	result.ValidationErrors = make(map[string]string)
	result.Platform = "PayZen"
	return result
}

func NewPayZenResultWithErrors(errors map[string]string) *PayZenResult {
	result := NewPayZenResult()
	result.Error = true
	result.ValidationErrors = errors
	return result
}

type PayZen struct {

  Lang string

	ToolBox *PayZenToolBox
  EntityValidator *validator.EntityValidator
  EntityValidatorResult *validator.EntityValidatorResult

  ValidationErrors map[string]string
  HasValidationError bool

  Debug bool
}

func NewPayZen(lang string) *PayZen {
	entityValidator := validator.NewEntityValidator(lang, "PayZen")
	return &PayZen{ Lang: lang, EntityValidator: entityValidator }
}

func (this *PayZen) OnDebug() {
	this.Debug = true
}

/* pay operations */

func (this *PayZen) PaymentCreate(payment *PayZenPayment) (*PayZenResult, error) {
	if !this.onValid(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

  if payment.Card.Scheme == SchemeBoleto && payment.Account.Mode == PayZenModeProduction {
    switch payment.Card.BoletoOnline {
      case BoletoOnline, BoletoOnlineItauIb, BoletoOnlineItauBoleto, BoletoOnlineBradescoBoleto:
        return this.PaymentCreateBoletoOnline(payment)
    }
  }

	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CreatePayment(payment)
}


func (this *PayZen) PaymentCreateBoletoOnline(payment *PayZenPayment) (*PayZenResult, error) {
	if !this.onValid(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}
	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CreatePaymentBoletoOnline(payment)
}

func (this *PayZen) PaymentUpdate(payment *PayZenPayment) (*PayZenResult, error) {

	valid := this.onValidOther(payment, func (validator *validation.Validation) {
		if len(payment.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.rquired"))
		}
		if payment.Amount <= 0 {
			validator.SetError("Amount", this.getMessage("Pagarme.rquired"))
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}


	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.UpdatePayment(payment)
}

// Para cancelar: Initial, Authorised,
func (this *PayZen) PaymentCancel(paymentFind *PayZenPaymentFind) (*PayZenResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.rquired"))
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentFind.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CancelPayment(paymentFind.TransactionUuid)
}

func (this *PayZen) PaymentCapture(capturePayment *PayZenCapturePayment) (*PayZenResult, error) {

	if !this.onValidOther(capturePayment, nil) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(capturePayment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CapturePayment(capturePayment)
}

func (this *PayZen) PaymentValidate(paymentFind *PayZenPaymentFind) (*PayZenResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.rquired"))
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentFind.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.ValidatePayment(paymentFind.TransactionUuid)
}

func (this *PayZen) PaymentDuplicate(payment *PayZenPayment) (*PayZenResult, error) {

	valid := this.onValidOther(payment, func (validator *validation.Validation) {
		if len(payment.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.rquired"))
		}
		if len(payment.OrderId) == 0 {
			validator.SetError("OrderId", this.getMessage("Pagarme.rquired"))
		}
		if payment.Amount <= 0 {
			validator.SetError("Amount", this.getMessage("Pagarme.rquired"))
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}


	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.DuplicatePayment(payment)
}


func (this *PayZen) PaymentRefund(payment *PayZenPayment) (*PayZenResult, error) {

	valid := this.onValidOther(payment, func (validator *validation.Validation) {
		if len(payment.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.rquired"))
		}
		if payment.Amount <= 0 {
			validator.SetError("Amount", this.getMessage("Pagarme.rquired"))
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}


	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.RefundPayment(payment)
}


/* token operations */

func (this *PayZen) PaymentTokenCreate(payment *PayZenPayment) (*PayZenResult, error) {
	payment.TokenOperation = true

	if !this.onValid(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CreatePaymentToken(payment)
}

func (this *PayZen) PaymentTokenUpdate(payment *PayZenPayment) (*PayZenResult, error) {
	payment.TokenOperation = true

	if !this.onValid(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	valid := this.onValidOther(payment, func (validator *validation.Validation) {
		if len(payment.Card.Token) == 0 {
			validator.SetError("Token", this.getMessage("Pagarme.rquired"))
			return
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}


	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.UpdatePaymentToken(payment)
}

func (this *PayZen) PaymentTokenCancel(paymentToken *PayZenPaymentToken) (*PayZenResult, error) {

	if !this.onValidOther(paymentToken, nil) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentToken.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CancelPaymentToken(paymentToken.Token)
}

func (this *PayZen) PaymentTokenReactive(paymentToken *PayZenPaymentToken) (*PayZenResult, error) {

	if !this.onValidOther(paymentToken, nil) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentToken.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.ReactivePaymentToken(paymentToken.Token)
}

func (this *PayZen) PaymentTokenGetDetails(paymentToken *PayZenPaymentToken) (*PayZenResult, error) {

	if !this.onValidOther(paymentToken, nil) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentToken.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.GetDetailsPaymentToken(paymentToken.Token)
}


/* find payment operations */
/*
	Com esse método é possível buscar todos os pagamentos relacionados a uma recorrência. As transações retornam no atributo Transactions do resultado
*/
func (this *PayZen) PaymentFind(paymentFind *PayZenPaymentFind) (*PayZenResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.OrderId) == 0 {
			validator.SetError("OrderId", this.getMessage("Pagarme.rquired"))
			return
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}


	this.ToolBox = NewPayZenToolBox(paymentFind.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.FindPayment(paymentFind.OrderId)
}

func (this *PayZen) PaymentGetDetails(paymentFind *PayZenPaymentFind) (*PayZenResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.rquired"))
			return
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}


	this.ToolBox = NewPayZenToolBox(paymentFind.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.GetPaymentDetails(paymentFind.TransactionUuid)
}

func (this *PayZen) PaymentGetDetailsWithNsu(paymentFind *PayZenPaymentFind) (*PayZenResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.rquired"))
			return
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}


	this.ToolBox = NewPayZenToolBox(paymentFind.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.GetPaymentDetailsWithNsu(paymentFind.TransactionUuid)
}

func (this *PayZen) PaymentCreateSubscription(subscription *PayZenSubscription) (*PayZenResult, error) {

	if !this.onValidSubscription(subscription) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(subscription.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CreateSubscription(subscription)
}

func (this *PayZen) PaymentGetDetailsSubscription(paymentFind *PayZenPaymentFind) (*PayZenResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.SubscriptionId) == 0 {
			validator.SetError("SubscriptionId", this.getMessage("Pagarme.rquired"))
			return
		}
		if len(paymentFind.Token) == 0 {
			validator.SetError("Token", this.getMessage("Pagarme.rquired"))
			return
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentFind.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.GetSubscriptionDetails(paymentFind.SubscriptionId, paymentFind.Token)
}

func (this *PayZen) PaymentCancelSubscription(paymentFind *PayZenPaymentFind) (*PayZenResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.SubscriptionId) == 0 {
			validator.SetError("SubscriptionId", this.getMessage("Pagarme.rquired"))
			return
		}
		if len(paymentFind.Token) == 0 {
			validator.SetError("Token", this.getMessage("Pagarme.rquired"))
			return
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentFind.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CancelSubscription(paymentFind.SubscriptionId, paymentFind.Token)
}

func (this *PayZen) PaymentUpdateSubscription(subscription *PayZenSubscription) (*PayZenResult, error) {

	if !this.onValidSubscription(subscription) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	valid := this.onValidOther(subscription, func (validator *validation.Validation) {
		if len(subscription.SubscriptionId) == 0 {
			validator.SetError("SubscriptionId", this.getMessage("Pagarme.rquired"))
			return
		}
		if len(subscription.Token) == 0 {
			validator.SetError("Token", this.getMessage("Pagarme.rquired"))
			return
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}


	this.ToolBox = NewPayZenToolBox(subscription.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.UpdateSubscription(subscription)
}

func (this *PayZen) CanCancelPayment(status PayZenTransactionStatus) bool {

	switch status {
		case Initial:
		case Authorised:
		case AuthorisedToValidate:
		case WaitingAuthorisation:
		case WaitingAuthorisationToValidate:
			return true
	}
	return false

}

func (this *PayZen) CanUpdatePayment(status PayZenTransactionStatus) bool {

	switch status {
		case Initial:
		case Authorised:
		case AuthorisedToValidate:
		case WaitingAuthorisation:
		case WaitingAuthorisationToValidate:
			return true
	}
	return false

}

func (this *PayZen) CanDuplicatePayment(status PayZenTransactionStatus) bool {

	switch status {
		case Captured:
			return true
	}
	return false

}

func (this *PayZen) onValidOther(object interface{}, action func (*validation.Validation)) bool {
	this.EntityValidatorResult, _ = this.EntityValidator.IsValid(object, action)

  if this.EntityValidatorResult.HasError {
  	this.onValidationErrors()
  	return false
  }

	return true
}

func (this *PayZen) onValidSubscription(subscription *PayZenSubscription) bool {

  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(subscription, func (validator *validation.Validation) {

  	if len(subscription.Rule) == 0 {

  		if !subscription.LastDayOfMonth && subscription.MonthDay == 0 {
  			validator.SetError("Rule", "you need set Rule or LastDayOfMonth or MonthDay")
  		}

  		if subscription.Count <= 0 {
  			validator.SetError("Count", this.getMessage("Pagarme.rquired"))
  		}

  		if subscription.EffectDate.IsZero() {
  			validator.SetError("EffectDate", this.getMessage("Pagarme.rquired"))
  		}

  		if subscription.InitialAmountNumber > 0 && subscription.InitialAmount == 0.0 {
  			validator.SetError("InitialAmount", this.getMessage("Pagarme.rquired"))
  		}

  		if subscription.InitialAmountNumber == 0 && subscription.InitialAmount > 0.0 {
  			validator.SetError("InitialAmountNumber", this.getMessage("Pagarme.rquired"))
  		}

  		if subscription.Amount <= 0 {
  			validator.SetError("Amount", this.getMessage("Pagarme.rquired"))
  		}

  	}

  })

  if this.EntityValidatorResult.HasError {
  	this.onValidationErrors()
  	return false
  }

  return true
}

func (this *PayZen) onValid(payment *PayZenPayment) bool {

  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(payment, func (validator *validation.Validation) {

  	if !payment.TokenOperation {

			if len(strings.TrimSpace(payment.OrderId)) == 0 {
				validator.SetError(this.getMessage("Pagarme.OrderId"), this.getMessage("Pagarme.rquired"))
			}

			if payment.Installments <= 0 {
				validator.SetError(this.getMessage("Pagarme.Installments"), this.getMessage("Pagarme.rquired"))
			}

			if payment.Amount <= 0 {
				validator.SetError(this.getMessage("Pagarme.Amount"), this.getMessage("Pagarme.rquired"))
			}

			if payment.Card == nil {
				validator.SetError(this.getMessage("Pagarme.Card"), this.getMessage("Pagarme.rquired"))
			}

			if payment.Customer == nil {
				validator.SetError(this.getMessage("Pagarme.Customer"), this.getMessage("Pagarme.rquired"))
			} else {
			
				if payment.Card.Scheme == SchemeBoleto {
					if len(strings.TrimSpace(payment.Customer.IdentityCode)) == 0 {
						validator.SetError(this.getMessage("Pagarme.IdentityCode"), this.getMessage("Pagarme.rquired"))
					}
				}

			}

  	}

  	if payment.Card.Scheme == SchemeBoleto {
	    switch payment.Card.BoletoOnline {

	      case BoletoOnline:
	      	validator.SetError(this.getMessage("Pagarme.BoletoOnline"), fmt.Sprintf("Boleto On-Line %v não implementado", payment.Card.BoletoOnline))
	        break

	      case BoletoOnlineItauIb:
	      	validator.SetError(this.getMessage("Pagarme.BoletoOnline"), fmt.Sprintf("Boleto On-Line %v não implementado", payment.Card.BoletoOnline))
	        break

	      case BoletoOnlineItauBoleto:
					if _, err := strconv.Atoi(payment.VadsTransId); err != nil {
						validator.SetError(this.getMessage("Pagarme.VadsTransId"), "vads_trans_id deve ser um valor númerico de 6 digitos que não pode repetir no mesmo dia.")
					}	        
	        break

	      case BoletoOnlineBradescoBoleto:
					if _, err := strconv.Atoi(payment.VadsTransId); err != nil {
						validator.SetError(this.getMessage("Pagarme.VadsTransId"), "vads_trans_id deve ser um valor númerico de 6 digitos que não pode repetir no mesmo dia.")
					}	        
	      	//validator.SetError(this.getMessage("Pagarme.BoletoOnline"), fmt.Sprintf("Boleto On-Line %v não implementado", payment.Card.BoletoOnline))
	        break
	        
	    }  		
  	}

  })

  if this.EntityValidatorResult.HasError {
  	this.onValidationErrors()
  	return false
  }

  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(payment.Card, func (validator *validation.Validation) {

  	if len(strings.TrimSpace(payment.Card.Token)) == 0 && payment.Card.Scheme != SchemeBoleto {

			if len(strings.TrimSpace(payment.Card.Number)) == 0 {
				validator.SetError(this.getMessage("Pagarme.Number"), this.getMessage("Pagarme.rquired"))
			}

			if len(strings.TrimSpace(payment.Card.Scheme)) == 0 {
				validator.SetError(this.getMessage("Pagarme.Scheme"), this.getMessage("Pagarme.rquired"))
			}

			if len(strings.TrimSpace(payment.Card.ExpiryMonth)) == 0 {
				validator.SetError(this.getMessage("Pagarme.ExpiryMonth"), this.getMessage("Pagarme.rquired"))
			}

			if len(strings.TrimSpace(payment.Card.ExpiryYear)) == 0 {
				validator.SetError(this.getMessage("Pagarme.ExpiryYear"), this.getMessage("Pagarme.rquired"))
			}

			/*
			if len(strings.TrimSpace(payment.Card.Number)) == 0 {
				validator.SetError("Number", this.getMessage("Pagarme.rquired"))
			}*/

			if len(strings.TrimSpace(payment.Card.CardSecurityCode)) == 0 {
				validator.SetError(this.getMessage("Pagarme.CardSecurityCode"), this.getMessage("Pagarme.rquired"))
			}

			if len(strings.TrimSpace(payment.Card.CardHolderName)) == 0 {
				validator.SetError(this.getMessage("Pagarme.CardHolderName"), this.getMessage("Pagarme.rquired"))
			}
  	}

  })

  if this.EntityValidatorResult.HasError {
  	this.onValidationErrors()
  	return false
  }

  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(payment.Customer, func (validator *validation.Validation) {

  })

  if this.EntityValidatorResult.HasError {
  	this.onValidationErrors()
  	return false
  }

  return true
}

func (this *PayZen) onValidationErrors(){
	this.HasValidationError = true
	data := make(map[interface{}]interface{})
  this.EntityValidator.CopyErrorsToView(this.EntityValidatorResult, data)
  this.ValidationErrors = data["errors"].(map[string]string)
}

func (this *PayZen) getMessage(key string, args ...interface{}) string{
  return i18n.Tr(this.Lang, key, args)
}
