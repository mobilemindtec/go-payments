package asaas

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mobilemindtec/go-payments/api"
	"github.com/mobilemindtec/go-utils/beego/validator"
)

type EventType string

const (
	// Transference
	//Geração de nova transferência.
	TransferCreated EventType = "TRANSFER_CREATED"
	//Transferência pendente de execução.
	TransferPending EventType = "TRANSFER_PENDING"
	//Transferência em processamento bancário.
	TransferInBankProcessing EventType = "TRANSFER_IN_BANK_PROCESSING"
	//Transferência bloqueada.
	TransferBlocked EventType = "TRANSFER_BLOCKED"
	//Transferência realizada.
	TransferDone EventType = "TRANSFER_DONE"
	//Transferência falhou.
	TransferFailed EventType = "TRANSFER_FAILED"
	//Transferência cancelada.
	TransferCancelled EventType = "TRANSFER_CANCELLED"

	// Payment
	//Geração de nova cobrança.
	PaymentCreated EventType = "PAYMENT_CREATED"
	//Pagamento em cartão aguardando aprovação pela análise manual de risco.
	PaymentAwaitingRiskAnalysis EventType = "PAYMENT_AWAITING_RISK_ANALYSIS"
	//Pagamento em cartão aprovado pela análise manual de risco.
	PaymentApprovedByRiskAnalysis EventType = "PAYMENT_APPROVED_BY_RISK_ANALYSIS"
	//Pagamento em cartão reprovado pela análise manual de risco.
	PaymentReprovedByRiskAnalysis EventType = "PAYMENT_REPROVED_BY_RISK_ANALYSIS"
	//Pagamento em cartão que foi autorizado e precisa ser capturado.
	PaymentAuthorized EventType = "PAYMENT_AUTHORIZED"
	//Alteração no vencimento ou valor de cobrança existente.
	PaymentUpdated EventType = "PAYMENT_UPDATED"
	//Cobrança confirmada (pagamento efetuado, porém o saldo ainda não foi disponibilizado).
	PaymentConfirmed EventType = "PAYMENT_CONFIRMED"
	//Cobrança recebida.
	PaymentReceived EventType = "PAYMENT_RECEIVED"
	//Falha no pagamento de cartão de crédito
	PaymentCreditCardCaptureRefused EventType = "PAYMENT_CREDIT_CARD_CAPTURE_REFUSED"
	//Cobrança antecipada.
	PaymentAnticipated EventType = "PAYMENT_ANTICIPATED"
	//Cobrança vencida.
	PaymentOverdue EventType = "PAYMENT_OVERDUE"
	//Cobrança removida.
	PaymentDeleted EventType = "PAYMENT_DELETED"
	//Cobrança restaurada.
	PaymentRestored EventType = "PAYMENT_RESTORED"
	//Cobrança estornada.
	PaymentRefunded EventType = "PAYMENT_REFUNDED"
	//Estorno em processamento (liquidação já está agendada, cobrança será estornada após executar a liquidação).
	PaymentRefundInProgress EventType = "PAYMENT_REFUND_IN_PROGRESS"
	//Recebimento em dinheiro desfeito.
	PaymentReceivedInCashUndone EventType = "PAYMENT_RECEIVED_IN_CASH_UNDONE"
	//Recebido chargeback.
	PaymentChargebackRequested EventType = "PAYMENT_CHARGEBACK_REQUESTED"
	//Em disputa de chargeback (caso sejam apresentados documentos para contestação).
	PaymentChargebackDispute EventType = "PAYMENT_CHARGEBACK_DISPUTE"
	//Disputa vencida, aguardando repasse da adquirente.
	PaymentAwaitingChargebackReversal EventType = "PAYMENT_AWAITING_CHARGEBACK_REVERSAL"
	//Recebimento de negativação.
	PaymentDunningReceived EventType = "PAYMENT_DUNNING_RECEIVED"
	//Requisição de negativação.
	PaymentDunningRequested EventType = "PAYMENT_DUNNING_REQUESTED"
	//Boleto da cobrança visualizado pelo cliente.
	PaymentBankSlipViewed EventType = "PAYMENT_BANK_SLIP_VIEWED"
	//Fatura da cobrança visualizada pelo cliente.
	PaymentCheckoutViewed EventType = "PAYMENT_CHECKOUT_VIEWED"

	//Account Status
	//Conta bancária aprovada
	AccountStatusBankAccountInfoApproved EventType = "ACCOUNT_STATUS_BANK_ACCOUNT_INFO_APPROVED"
	//Conta bancária está em análise
	AccountStatusBankAccountInfoAwaitingApproval EventType = "ACCOUNT_STATUS_BANK_ACCOUNT_INFO_AWAITING_APPROVAL"
	//Conta bancária voltou para pendente
	AccountStatusBankAccountInfoPending EventType = "ACCOUNT_STATUS_BANK_ACCOUNT_INFO_PENDING"
	//Conta bancária reprovada
	AccountStatusBankAccountInfoRejected EventType = "ACCOUNT_STATUS_BANK_ACCOUNT_INFO_REJECTED"
	//Informações comerciais aprovada
	AccountStatusCommercialInfoApproved EventType = "ACCOUNT_STATUS_COMMERCIAL_INFO_APPROVED"
	//Informações comerciais em análise
	AccountStatusCommercialInfoAwaitingApproval EventType = "ACCOUNT_STATUS_COMMERCIAL_INFO_AWAITING_APPROVAL"
	//Informações comerciais voltou para pendente
	AccountStatusCommercialInfoPending EventType = "ACCOUNT_STATUS_COMMERCIAL_INFO_PENDING"
	//Informações comerciais reprovada
	AccountStatusCommercialInfoRejected EventType = "ACCOUNT_STATUS_COMMERCIAL_INFO_REJECTED"
	//Documentos aprovados
	AccountStatusDocumentApprovedEventType = "ACCOUNT_STATUS_DOCUMENT_APPROVED"
	//Documentos em análise
	AccountStatusDocumentAwaitingApproval EventType = "ACCOUNT_STATUS_DOCUMENT_AWAITING_APPROVAL"
	//Documentos voltaram para pendente
	AccountStatusDocumentPendingEventType = "ACCOUNT_STATUS_DOCUMENT_PENDING"
	//Documentos reprovados
	AccountStatusDocumentRejectedEventType = "ACCOUNT_STATUS_DOCUMENT_REJECTED"
	//Conta aprovada
	AccountStatusGeneralApprovalApproved EventType = "ACCOUNT_STATUS_GENERAL_APPROVAL_APPROVED"
	//Conta em análise
	AccountStatusGeneralApprovalAwaitingApproval EventType = "ACCOUNT_STATUS_GENERAL_APPROVAL_AWAITING_APPROVAL"
	//Conta voltou para pendente
	AccountStatusGeneralApprovalPending EventType = "ACCOUNT_STATUS_GENERAL_APPROVAL_PENDING"
	//Conta reprovada
	AccountStatusGeneralApprovalRejected EventType = "ACCOUNT_STATUS_GENERAL_APPROVAL_REJECTED"
)

type WebhookData struct {
	Event    api.PaymentEvent `json:"event" valid:"Required"`
	Response *Response        `json:"payment" valid:"Required"`
	Raw      string           `json:"raw" valid:"Required"`
}

func NewWebhookData() *WebhookData {
	return &WebhookData{}
}

type Webhook struct {
	Debug              bool
	EntityValidator    *validator.EntityValidator
	ValidationErrors   map[string]string
	HasValidationError bool
}

func NewWebhook(lang string) *Webhook {
	entityValidator := validator.NewEntityValidator(lang, "Asaas")
	return &Webhook{EntityValidator: entityValidator}
}

func NewDefaultWebhook() *Webhook {
	entityValidator := validator.NewEntityValidator("pt-BR", "Asaas")
	return &Webhook{EntityValidator: entityValidator}
}

func (this *Webhook) SetDebug() {
	this.Debug = true
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
