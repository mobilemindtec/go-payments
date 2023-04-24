package payzen

import (
	"github.com/mobilemindtec/go-utils/beego/validator"
	"github.com/beego/beego/v2/core/validation"
	"github.com/mobilemindtec/go-payments/api"
	"github.com/beego/i18n"
	"strconv"
	"strings"
	"errors"
	"fmt"
)



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

func (this *PayZen) PaymentCreate(payment *api.Payment) (*api.PaymentResult, error) {
	if !this.onValid(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

  if payment.Card.Scheme == SchemeBoleto && payment.Account.Mode == PayZenModeProduction {
    switch payment.BoletoOnline {
      case api.BoletoOnline, api.BoletoOnlineItauIb, api.BoletoOnlineItauBoleto, api.BoletoOnlineBradescoBoleto:
        return this.PaymentCreateBoletoOnline(payment)
    }
  }

	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CreatePayment(payment)
}


func (this *PayZen) PaymentCreateBoletoOnline(payment *api.Payment) (*api.PaymentResult, error) {
	if !this.onValid(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}
	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CreatePaymentBoletoOnline(payment)
}

func (this *PayZen) PaymentUpdate(payment *api.Payment) (*api.PaymentResult, error) {

	valid := this.onValidOther(payment, func (validator *validation.Validation) {
		if len(payment.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.required"))
		}
		if payment.Amount <= 0 {
			validator.SetError("Amount", this.getMessage("Pagarme.required"))
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
func (this *PayZen) PaymentCancel(paymentFind *api.PaymentFind) (*api.PaymentResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.required"))
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentFind.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CancelPayment(paymentFind.TransactionUuid)
}

func (this *PayZen) PaymentCapture(capturePayment *api.PaymentCapture) (*api.PaymentResult, error) {

	if !this.onValidOther(capturePayment, nil) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(capturePayment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CapturePayment(capturePayment)
}

func (this *PayZen) PaymentValidate(paymentFind *api.PaymentFind) (*api.PaymentResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.required"))
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentFind.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.ValidatePayment(paymentFind.TransactionUuid)
}

func (this *PayZen) PaymentDuplicate(payment *api.Payment) (*api.PaymentResult, error) {

	valid := this.onValidOther(payment, func (validator *validation.Validation) {
		if len(payment.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.required"))
		}
		if len(payment.OrderId) == 0 {
			validator.SetError("OrderId", this.getMessage("Pagarme.required"))
		}
		if payment.Amount <= 0 {
			validator.SetError("Amount", this.getMessage("Pagarme.required"))
		}
	})

	if !valid {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}


	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.DuplicatePayment(payment)
}


func (this *PayZen) PaymentRefund(payment *api.Payment) (*api.PaymentResult, error) {

	valid := this.onValidOther(payment, func (validator *validation.Validation) {
		if len(payment.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.required"))
		}
		if payment.Amount <= 0 {
			validator.SetError("Amount", this.getMessage("Pagarme.required"))
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

func (this *PayZen) PaymentTokenCreate(payment *api.Payment) (*api.PaymentResult, error) {
	payment.TokenOperation = true

	if !this.onValid(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(payment.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CreatePaymentToken(payment)
}

func (this *PayZen) PaymentTokenUpdate(payment *api.Payment) (*api.PaymentResult, error) {
	payment.TokenOperation = true

	if !this.onValid(payment) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	valid := this.onValidOther(payment, func (validator *validation.Validation) {
		if len(payment.Card.Token) == 0 {
			validator.SetError("Token", this.getMessage("Pagarme.required"))
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

func (this *PayZen) PaymentTokenCancel(paymentToken *api.PaymentToken) (*api.PaymentResult, error) {

	if !this.onValidOther(paymentToken, nil) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentToken.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CancelPaymentToken(paymentToken.Token)
}

func (this *PayZen) PaymentTokenReactive(paymentToken *api.PaymentToken) (*api.PaymentResult, error) {

	if !this.onValidOther(paymentToken, nil) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(paymentToken.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.ReactivePaymentToken(paymentToken.Token)
}

func (this *PayZen) PaymentTokenGetDetails(paymentToken *api.PaymentToken) (*api.PaymentResult, error) {

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
func (this *PayZen) PaymentFind(paymentFind *api.PaymentFind) (*api.PaymentResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.OrderId) == 0 {
			validator.SetError("OrderId", this.getMessage("Pagarme.required"))
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

func (this *PayZen) PaymentGetDetails(paymentFind *api.PaymentFind) (*api.PaymentResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.required"))
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

func (this *PayZen) PaymentGetDetailsWithNsu(paymentFind *api.PaymentFind) (*api.PaymentResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.TransactionUuid) == 0 {
			validator.SetError("TransactionUuid", this.getMessage("Pagarme.required"))
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

func (this *PayZen) PaymentCreateSubscription(subscription *api.Subscription) (*api.PaymentResult, error) {

	if !this.onValidSubscription(subscription) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	this.ToolBox = NewPayZenToolBox(subscription.Account)
	this.ToolBox.Debug = this.Debug
	return this.ToolBox.CreateSubscription(subscription)
}

func (this *PayZen) PaymentGetDetailsSubscription(paymentFind *api.PaymentFind) (*api.PaymentResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.SubscriptionId) == 0 {
			validator.SetError("SubscriptionId", this.getMessage("Pagarme.required"))
			return
		}
		if len(paymentFind.Token) == 0 {
			validator.SetError("Token", this.getMessage("Pagarme.required"))
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

func (this *PayZen) PaymentCancelSubscription(paymentFind *api.PaymentFind) (*api.PaymentResult, error) {

	valid := this.onValidOther(paymentFind, func (validator *validation.Validation) {
		if len(paymentFind.SubscriptionId) == 0 {
			validator.SetError("SubscriptionId", this.getMessage("Pagarme.required"))
			return
		}
		if len(paymentFind.Token) == 0 {
			validator.SetError("Token", this.getMessage("Pagarme.required"))
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

func (this *PayZen) PaymentUpdateSubscription(subscription *api.Subscription) (*api.PaymentResult, error) {

	if !this.onValidSubscription(subscription) {
		return nil, errors.New(this.getMessage("PayZen.ValidationError"))
	}

	valid := this.onValidOther(subscription, func (validator *validation.Validation) {
		if len(subscription.SubscriptionId) == 0 {
			validator.SetError("SubscriptionId", this.getMessage("Pagarme.required"))
			return
		}
		if len(subscription.Token) == 0 {
			validator.SetError("Token", this.getMessage("Pagarme.required"))
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

func (this *PayZen) CanCancelPayment(status api.TransactionStatus) bool {

	switch status {
		case api.Initial:
		case api.Authorised:
		case api.AuthorisedToValidate:
		case api.WaitingAuthorisation:
		case api.WaitingAuthorisationToValidate:
			return true
	}
	return false

}

func (this *PayZen) CanUpdatePayment(status api.TransactionStatus) bool {

	switch status {
		case api.Initial:
		case api.Authorised:
		case api.AuthorisedToValidate:
		case api.WaitingAuthorisation:
		case api.WaitingAuthorisationToValidate:
			return true
	}
	return false

}

func (this *PayZen) CanDuplicatePayment(status api.TransactionStatus) bool {

	switch status {
		case api.Captured:
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

func (this *PayZen) onValidSubscription(subscription *api.Subscription) bool {

  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(subscription, func (validator *validation.Validation) {

  	if len(subscription.Rule) == 0 {

  		if subscription.EffectDate.IsZero() {
  			validator.SetError("EffectDate", this.getMessage("Pagarme.required"))
  		}

  		if subscription.InitialAmountNumber > 0 && subscription.InitialAmount == 0.0 {
  			validator.SetError("InitialAmount", this.getMessage("Pagarme.required"))
  		}

  		if subscription.InitialAmountNumber == 0 && subscription.InitialAmount > 0.0 {
  			validator.SetError("InitialAmountNumber", this.getMessage("Pagarme.required"))
  		}

  		if subscription.Amount <= 0 {
  			validator.SetError("Amount", this.getMessage("Pagarme.required"))
  		}

  	}

  })

  if this.EntityValidatorResult.HasError {
  	this.onValidationErrors()
  	return false
  }

  return true
}

func (this *PayZen) onValid(payment *api.Payment) bool {

  this.EntityValidatorResult, _ = this.EntityValidator.IsValid(payment, func (validator *validation.Validation) {

  	if !payment.TokenOperation {

			if len(strings.TrimSpace(payment.OrderId)) == 0 {
				validator.SetError(this.getMessage("Pagarme.OrderId"), this.getMessage("Pagarme.required"))
			}

			if payment.Installments <= 0 {
				validator.SetError(this.getMessage("Pagarme.Installments"), this.getMessage("Pagarme.required"))
			}

			if payment.Amount <= 0 {
				validator.SetError(this.getMessage("Pagarme.Amount"), this.getMessage("Pagarme.required"))
			}

			if payment.Card == nil {
				validator.SetError(this.getMessage("Pagarme.Card"), this.getMessage("Pagarme.required"))
			}

			if payment.Customer == nil {
				validator.SetError(this.getMessage("Pagarme.Customer"), this.getMessage("Pagarme.required"))
			} else {
			
				if payment.Card.Scheme == SchemeBoleto {
					if len(strings.TrimSpace(payment.Customer.IdentityCode)) == 0 {
						validator.SetError(this.getMessage("Pagarme.IdentityCode"), this.getMessage("Pagarme.required"))
					}
				}

			}

  	}

  	if payment.Card.Scheme == SchemeBoleto {
	    switch payment.BoletoOnline {

	      case api.BoletoOnline:
	      	validator.SetError(this.getMessage("Pagarme.BoletoOnline"), fmt.Sprintf("Boleto On-Line %v não implementado", payment.BoletoOnline))
	        break

	      case api.BoletoOnlineItauIb:
	      	validator.SetError(this.getMessage("Pagarme.BoletoOnline"), fmt.Sprintf("Boleto On-Line %v não implementado", payment.BoletoOnline))
	        break

	      case api.BoletoOnlineItauBoleto:
					if _, err := strconv.Atoi(payment.VadsTransId); err != nil {
						validator.SetError(this.getMessage("Pagarme.VadsTransId"), "vads_trans_id deve ser um valor númerico de 6 digitos que não pode repetir no mesmo dia.")
					}	        
	        break

	      case api.BoletoOnlineBradescoBoleto:
					if _, err := strconv.Atoi(payment.VadsTransId); err != nil {
						validator.SetError(this.getMessage("Pagarme.VadsTransId"), "vads_trans_id deve ser um valor númerico de 6 digitos que não pode repetir no mesmo dia.")
					}	        
	      	//validator.SetError(this.getMessage("Pagarme.BoletoOnline"), fmt.Sprintf("Boleto On-Line %v não implementado", payment.BoletoOnline))
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
				validator.SetError(this.getMessage("Pagarme.Number"), this.getMessage("Pagarme.required"))
			}

			if len(strings.TrimSpace(payment.Card.Scheme)) == 0 {
				validator.SetError(this.getMessage("Pagarme.Scheme"), this.getMessage("Pagarme.required"))
			}

			if len(strings.TrimSpace(payment.Card.ExpiryMonth)) == 0 {
				validator.SetError(this.getMessage("Pagarme.ExpiryMonth"), this.getMessage("Pagarme.required"))
			}

			if len(strings.TrimSpace(payment.Card.ExpiryYear)) == 0 {
				validator.SetError(this.getMessage("Pagarme.ExpiryYear"), this.getMessage("Pagarme.required"))
			}

			/*
			if len(strings.TrimSpace(payment.Card.Number)) == 0 {
				validator.SetError("Number", this.getMessage("Pagarme.required"))
			}*/

			if len(strings.TrimSpace(payment.Card.CardSecurityCode)) == 0 {
				validator.SetError(this.getMessage("Pagarme.CardSecurityCode"), this.getMessage("Pagarme.required"))
			}

			if len(strings.TrimSpace(payment.Card.CardHolderName)) == 0 {
				validator.SetError(this.getMessage("Pagarme.CardHolderName"), this.getMessage("Pagarme.required"))
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
