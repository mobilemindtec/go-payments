package asaas

import (
  "github.com/mobilemindtec/go-utils/beego/validator"	
	beego "github.com/beego/beego/v2/server/web"
	"encoding/json"
	"errors"
	"fmt"
)


/*
PAYMENT_CREATED - Geração de nova cobrança.
PAYMENT_UPDATED - Alteração no vencimento ou valor de cobrança existente.
PAYMENT_CONFIRMED - Cobrança confirmada (pagamento efetuado, porém o saldo ainda não foi disponibilizado).
PAYMENT_RECEIVED - Cobrança recebida.
PAYMENT_OVERDUE - Cobrança vencida.
PAYMENT_DELETED - Cobrança removida.
PAYMENT_RESTORED - Cobrança restaurada.
PAYMENT_REFUNDED - Cobrança estornada.
PAYMENT_RECEIVED_IN_CASH_UNDONE - Recebimento em dinheiro desfeito.
PAYMENT_CHARGEBACK_REQUESTED - Recebido chargeback.
PAYMENT_CHARGEBACK_DISPUTE - Em disputa de chargeback (caso sejam apresentados documentos para contestação).
PAYMENT_AWAITING_CHARGEBACK_REVERSAL - Disputa vencida, aguardando repasse da adquirente.
PAYMENT_DUNNING_RECEIVED - Recebimento de recuperação.
PAYMENT_DUNNING_REQUESTED - Requisição de recuperação.
*/
type WebhookEvent string

const (
	EventPaymentCreated WebhookEvent = "PAYMENT_CREATED"
	EventPaymentUpdated WebhookEvent = "PAYMENT_UPDATED"
	EventPaymentConfirmed WebhookEvent = "PAYMENT_CONFIRMED"
	EventPaymentReceived WebhookEvent = "PAYMENT_RECEIVED"
	EventPaymentOverdue WebhookEvent = "PAYMENT_OVERDUE"
	EventPaymentDeleted WebhookEvent = "PAYMENT_DELETED"
	EventPaymentRestored WebhookEvent = "PAYMENT_RESTORED"
	EventPaymentRefunded WebhookEvent = "PAYMENT_REFUNDED"
	EventPaymentReceivedInCashUndone WebhookEvent = "PAYMENT_RECEIVED_IN_CASH_UNDONE"
	EventPaymentChargebackRequested WebhookEvent = "PAYMENT_CHARGEBACK_REQUESTED"
	EventPaymentChargebackDispute WebhookEvent = "PAYMENT_CHARGEBACK_DISPUTE"
	EventPaymentAwaitingChargebackReversal WebhookEvent = "PAYMENT_AWAITING_CHARGEBACK_REVERSAL"
	EventPaymentDunningReceived WebhookEvent = "PAYMENT_DUNNING_RECEIVED"
	EventPaymentDunningRequested WebhookEvent = "PAYMENT_DUNNING_REQUESTED"	
)

type WebhookData struct {
	Event WebhookEvent `json:"event" valid:"Required"`
	Response *Response `json:"payment" valid:"Required"`
	Raw string `json:"raw" valid:"Required"`
	Uuid string `json:"uuid" valid:""`
} 

func NewWebhookData() *WebhookData{
	return &WebhookData{}
}

type Webhook struct {
	AccessToken string
	Controller *beego.Controller	
	Debug bool
  EntityValidator *validator.EntityValidator  
  ValidationErrors map[string]string
  HasValidationError bool	
}

func NewWebhook(lang string, accessToken string, controller *beego.Controller) *Webhook {
	entityValidator := validator.NewEntityValidator(lang, "Asaas")
	return &Webhook{ AccessToken: accessToken, EntityValidator: entityValidator }
}

func (this *Webhook) SetDebug()  {
	this.Debug = true
}

func (this *Webhook) IsValid() bool {
	token := this.Controller.Ctx.Request.Header.Get("asaas-access-token")
	return token == this.AccessToken
}

func (this *Webhook) GetData() (*WebhookData, error) {
	body := this.Controller.Ctx.Input.RequestBody
	return this.Parse(body)
}

func (this *Webhook) Parse(body []byte) (*WebhookData, error) {
	data := NewWebhookData()

  if this.Debug {
    fmt.Println("************************************************")
    fmt.Println("**** Asaas.Webhook: ", string(body))
    fmt.Println("************************************************")
  }

	err := json.Unmarshal(body, data)	

	if data.Response != nil {
		data.Response.BuildStatus()
	}

	data.Raw = string(body)

	if err != nil {
	  entityValidatorResult, _ := this.EntityValidator.IsValid(data, nil)  

	  if entityValidatorResult.HasError {
	    this.HasValidationError = true
	    this.ValidationErrors = this.EntityValidator.GetValidationErrors(entityValidatorResult)
	    return nil, errors.New("Validation error")
	  }		
	}

	return data, err
}
