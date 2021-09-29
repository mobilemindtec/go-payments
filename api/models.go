package api

import (
	"time"

)

type Gateway string

const (
  GatewayNone Gateway = ""
  GatewayPagarme Gateway = "Pagarme"
  GatewayPayZen Gateway = "PayZen"
  GatewayAsaas Gateway = "Asaas"
  GatewayPicPay Gateway = "PicPay"
)

type PaymentStatus int

const (
	PaymentInitial PaymentStatus = 1 + iota
	PaymentWaitingPayment
	PaymentPaid
	PaymentRefused
	PaymentCancelled
	PaymentRefound
	PaymentExpired
	PaymentChargeback
	PaymentOther
	PaymentSuccess
	PaymentError
)

type PaymentStatusLabel string

const (
	PaymentInitialLabel PaymentStatusLabel = "initial"
	PaymentWaitingPaymentLabel  PaymentStatusLabel = "waiting_payment"
	PaymentPaidLabel PaymentStatusLabel = "paid"
	PaymentRefusedLabel  PaymentStatusLabel = "refused"
	PaymentCancelledLabel  PaymentStatusLabel = "cancelled"
	PaymentRefoundLabel  PaymentStatusLabel = "refound"
	PaymentExpiredLabel  PaymentStatusLabel = "expired"
	PaymentChargebackLabel PaymentStatusLabel = "changeback"
	PaymentOtherLabel PaymentStatusLabel = "other"
	PaymentSuccessLabel PaymentStatusLabel = "success"
	PaymentErrorLabel  PaymentStatusLabel = "error"
)


/*
INITIAL Em andamento. Status temporário. O status definitivo será retornado assim que a sincronização for realizada.
NOT_CREATED A transação não foi criada e portanto não será exibida no Back Office.
AUTHORISED Captura em andamento. A transação foi aceita e será capturada automaticamente no banco na data prevista.
AUTHORISED_TO_VALIDATE Para ser aprovado. A transação, criada em validação manual, foi autorizada. O vendedor deve validar manualmente a captura no banco. A transação pode ser aprovada enquanto a data de captura não for vencida. Se esta data estiver vencida, então o pagamento tem o status Expirado (status definitivo).
WAITING_AUTHORISATION Autorização em andamento. A data de captura solicitada é superior à data de fim de validade da solicitação de autorização. Uma autorização de 1 BRL foi efetuada e aceita pelo banco emissor. A solicitado de autorização e a captura no banco serão acionadas automaticamente.
WAITING_AUTHORISATION_TO_VALIDATE Para ser aprovado e autorizado. A data de captura solicitada é superior à data de fim de validade da solicitação de autorização. Uma autorização de 1 BRL foi efetuada e aceita pelo banco emissor. A solicitação de autorização será automaticamente efetuada a D-1 antes da data de captura no banco. O pagamento poderá ser aceito ou recusado. Captura automática no banco.
REFUSED Recusada. A transação foi recusada.	
CAPTURED A transação foi capturada no banco.
CANCELLED Cancelada. A transação foi cancelada pelo vendedor.
EXPIRED Expirada. A date de captura foi atingida mas o vendedor não validou a transação.
UNDER_VERIFICATION (Específico a PayPal) Verificação por PayPal em andamento. Este valor significa que Paypal segura a transação por causa de uma suspeita de fraude. O pagamento fica então na aba Pagamento em andamento.

*/
type TransactionStatus int

const (
	Initial TransactionStatus = 1 + iota
	NotCreated
	Authorised
	AuthorisedToValidate
	WaitingAuthorisation
	WaitingAuthorisationToValidate
	Refused
	Captured
	Cancelled
	Expired
	UnderVerification
	PartiallyAuthorised
	Refunded //pickpay,pagarme, asaas
	Created //pickpay, pagarme
	Chargeback //pickpay,pagarme,asaas
	WaitingPayment //pagarme, asaas
	PendingRefund //pagarme, asaas
	Analyzing //pagarme
	PendingReview //pagarme
	ReceivedInCash // asaas
	Other
	Canceled
	Success
	Error
)

type PagarmeStatus int

/*

processing	Transação está processo de autorização.
authorized	Transação foi autorizada. Cliente possui saldo na conta e este valor foi reservado para futura captura, que deve acontecer em até 5 dias para transações criadas com api_key. Caso não seja capturada, a autorização é cancelada automaticamente pelo banco emissor, e o status da transação permanece authorized.
paid	Transação paga. Foi autorizada e capturada com sucesso, e para boleto, significa que nossa API já identificou o pagamento de seu cliente.
refunded	Transação estornada completamente.
waiting_payment	Transação aguardando pagamento (status válido para boleto bancário).
pending_refund	Transação do tipo boleto e que está aguardando para confirmação do estorno solicitado.
refused	Transação recusada, não autorizada.
chargedback	Transação sofreu chargeback. Mais em nossa central de ajuda
analyzing	Transação encaminhada para a análise manual feita por um especialista em prevenção a fraude.
pending_review	Transação pendente de revisão manual por parte do lojista. Uma transação ficar

*/

const (
  PagarmeProcessing PagarmeStatus = 1 + iota
  PagarmeAuthorized         
  PagarmePaid 
  PagarmeRefunded
  PagarmeWaitingPayment 
  PagarmePendingRefund 
  PagarmeRefused
	PagarmeChargedback
	PagarmeAnalyzing
	PagarmePendingReview
	PagarmeSuccess
	PagarmeEmpty
	PagarmeCancelled
	PagarmeError
)


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
  PicPayCancelled
)

// ASAAS STATUS

type AsaasStatus int64

const (
	AsaasPending AsaasStatus = iota + 1 //- Aguardando pagamento
	AsaasReceived //- Recebida (saldo já creditado na conta)
	AsaasConfirmed //- Pagamento confirmado (saldo ainda não creditado)
	AsaasOverdue //- Vencida
	AsaasRefunded //- Estornada
	AsaasReceivedInCash //- Recebida em dinheiro (não gera saldo na conta)
	AsaasRefundRequested //- Estorno Solicitado
	AsaasChargebackRequested //- Recebido chargeback
	AsaasChargebackDispute //- Em disputa de chargeback (caso sejam apresentados documentos para contestação)
	AsaasAwaitingChargebackReversal //- Disputa vencida, aguardando repasse da adquirente
	AsaasDunningRequested //- Em processo de recuperação
	AsaasDunningReceived //- Recuperada
	AsaasAwaitingRiskAnalysis //- Pagamento em análise
	AsaasActive // subscription
	AsaasExpired // subscription
	AsaasDeleted
	AsaasSuccess
)

type PayZenPaymentValidationType int

const (
	Automatica PayZenPaymentValidationType = 0 + iota
	Manual 
)

type BoletoOnlineTipo string

const (
	BoletoOnline BoletoOnlineTipo = "BOLETO"
	BoletoOnlineItauIb BoletoOnlineTipo = "ITAU_IB"
	BoletoOnlineItauBoleto BoletoOnlineTipo = "ITAU_BOLETO"
	BoletoOnlineBradescoBoleto BoletoOnlineTipo = "BRADESCO_BOLETO"
	BoletoOnlineNenhum BoletoOnlineTipo = "NENHUM"
)

type PayZenApiVersion string

const (
	PayZenApiSOAP PayZenApiVersion = "SOAP"
	PayZenApiRestV4 PayZenApiVersion = "RESTFul.v4"
)

type SubscriptionCycle string

const (
	SubscriptionCycleNone SubscriptionCycle = ""
	Weekly SubscriptionCycle = "WEEKLY" // semanal
	Biweekly SubscriptionCycle = "BIWEEKLY" // quinzenal
	Monthly SubscriptionCycle = "MONTHLY" // mensal
	Quarterly SubscriptionCycle = "QUARTERLY" // trimestral
	Semiannually SubscriptionCycle = "SEMIANNUALLY" // semestral
	Yearly SubscriptionCycle = "YEARLY" // anual
)

//type PaymentMethod string // Pagarme
//
//const(
//	PaymentMethodBoleto PaymentMethod = "boleto"
//	PaymentMethodCreditCard PaymentMethod = "credit_card"
//	PaymentMethodPix PaymentMethod = "pix"
//)

type PaymentType string

const (
	PaymentTypeNone PaymentType = ""
	PaymentTypeCreditCard PaymentType = "credit_card"	
	PaymentTypeDebitCard PaymentType = "debit_card"	
	PaymentTypeBoleto PaymentType = "boleto"	
	PaymentTypePix PaymentType = "pix"	
	PaymentTypePicPay PaymentType = "picpay"	
)

type BankAccountType string

const (
	BankAccountTypeNone BankAccountType = ""
	ContaCorrente BankAccountType = "CONTA_CORRENTE"
	ContaPoupanca BankAccountType = "CONTA_POUPANCA"
)

type TransferStatus string 

const (
	TransferPending TransferStatus	 = "PENDING"
	TransferBankProcessing TransferStatus	 = "BANK_PROCESSING"
	TransferDone TransferStatus	 = "DONE"
	TransferCancelled TransferStatus	 = "CANCELLED"
	TransferFailed TransferStatus	 = "FAILED"

)

type AsaasMode int64

const (
	AsaasModeProd AsaasMode = iota + 1
	AsaasModeTest
)

type Filter struct {
	Limit int64 `jsonp:""`
	Offset int64 `jsonp:""`
	StartDate string `jsonp:""`
	FinishDate string`jsonp:""`
	
	DateCreated string `jsonp:""` // transfer filter asaas

	Status string `jsonp:""` // pagarme
	BankAccountId string `jsonp:""`	//pagarme

	RecebedorId string `jsonp:""`
}

func NewFilter() *Filter{
	return &Filter{}
}

type Bank struct {
	Code string `jsonp:""`
}

func NewBank() *Bank{
	return &Bank{}
}

type BankAccount struct {
	Bank *Bank `jsonp:""`
	AccountName string `jsonp:""`
	OwnerName string `jsonp:""`
	//Data de nascimento do proprietário da conta. 
	//Somente quando a conta bancária não pertencer ao mesmo CPF ou CNPJ da conta Asaas.
	OwnerBirthDate string `jsonp:""`
	CpfCnpj string `jsonp:""`
	Agency string `jsonp:""`
	Account string `jsonp:""`
	AccountDigit string `jsonp:""`
	BankAccountType BankAccountType `jsonp:""`
}

func NewBankAccount() *BankAccount{
	return &BankAccount{}
}

type Transfer struct {
	Amount float64 `jsonp:""`
	BankAccountId int64 `jsonp:""`
	BankAccount *BankAccount `jsonp:""`
}

func NewTransfer() *Transfer{
	return &Transfer{}
}


type BoletoFine struct {
	Days int64 `jsonp:""`
	Amount float64 `jsonp:""`
	Percentage float64 `jsonp:""`
}

func NewBoletoFine() *BoletoFine {
	return &BoletoFine{}
}

type BoletoInterest struct {
	Days int64 `jsonp:""`
	Amount float64 `jsonp:""`
	Percentage float64 `jsonp:""`
}

func NewBoletoInterest() *BoletoInterest {
	return &BoletoInterest{}
}

type BoletoDiscount struct {
	Days int64 `jsonp:""`
	Amount float64 `jsonp:""`
	Percentage float64 `jsonp:""`
}

func NewBoletoDiscount() *BoletoFine {
	return &BoletoFine{}
}

type Plan struct {
	Id string `jsonp:""`
	Amount int64 `jsonp:""`
	Days int64 `jsonp:""` // Prazo, em dias, para cobrança das parcelas
	Name string `jsonp:""`
	TrialDays int64 `jsonp:""`
	PaymentMethods []PaymentType `jsonp:""`
	Charges int64 `jsonp:""`
	InvoiceReminder int64 `jsonp:""` 
}

func NewPlan() *Plan {
	return &Plan{ PaymentMethods: []PaymentType{} }
}

type PayZenAccount struct {
	ShopId string `valid:"Required"`
	Mode string `valid:"Required"`
	Cert string `valid:"Required"`
}

type Card struct {
	Number string `valid:"" jsonp:""`
	Scheme string `valid:"" jsonp:"brand"`
	ExpiryMonth string `valid:"" jsonp:"expiry_month"`
	ExpiryYear string `valid:"" jsonp:"expiry_year"`
	CardSecurityCode string `valid:"" jsonp:"cvv"`
	CardHolderBirthDay string `valid:""`
	CardHolderName string `valid:"" jsonp:"holder_name"`
	Token string`valid:"" jsonp:""`
}

type Customer struct {
	Id string `jsonp:""`
	FirstName string `valid:"Required" jsonp:""`
	LastName string `valid:"" jsonp:""`
	PhoneNumber string `valid:"" jsonp:""`
	CellPhoneNumber string `valid:"" jsonp:""`
	Address2 string `valid:"" jsonp:""`
	Email string `valid:"Required" jsonp:""`
	StreetNumber string `valid:"" jsonp:""`
	Address string `valid:"" jsonp:""`
	District string `valid:"" jsonp:""`
	ZipCode string `valid:"" jsonp:""`
	City string `valid:"" jsonp:""`
	State string `valid:"" jsonp:""`
	Country string `valid:"Required" jsonp:""`
	IdentityCode string `valid:"Required" jsonp:"document"`
	ExternalReference string `jsonp:""`
}

func NewCustomer() *Customer{
	return &Customer{}
}

type Subscription struct{

	OrderId string `valid:"Required" jsonp:""`
	SubscriptionId string `valid:"Required" jsonp:""`
	TransactionId string `valid:"" jsonp:""`
	Description string `valid:"Required" jsonp:""`

	// valor da recorrência
	Amount float64 `valid:"Required" jsonp:""`

	// valor inicial da recorrência
	InitialAmount float64 `valid:"" jsonp:""`
	// quantas vezes o valor inicial deve ser cobrado
	InitialAmountNumber int64 `valid:"" jsonp:""`

	// data de inicio da cobrança
	EffectDate time.Time `valid:"Required" jsonp:""`
	// cobrar no último dia do mês

	Count int64 `jsonp:""`
	Cycle SubscriptionCycle `jsonp:""`
	PaymentAtLastDayOfMonth bool `jsonp:""`
	PaymentAtDayOfMonth int64 `jsonp:""`

	Rule string `jsonp:""`

	
	Account *PayZenAccount `valid:"Required"`
	Token string `valid:"" jsonp:""`

	BoletoFine *BoletoFine `jsonp:""`
	BoletoInterest *BoletoInterest `jsonp:""`
	BoletoDiscount *BoletoDiscount `jsonp:""`
	Card *Card `valid:"" jsonp:""`
	PostbackUrl string `jsonp:""`
	WebhookUrl string `jsonp:""`
	PaymentType PaymentType `jsonp:""` // pagarme/asaas
	PlanId string `jsonp:""` // pagarme
	Customer *Customer `jsonp:""`
	AdditionalInfo map[string]interface{} `jsonp:""`
}

type Payment struct{
	OrderId string `valid:"" jsonp:""`
	Installments int64 `valid:"" jsonp:""`
	Amount float64 `valid:"" jsonp:""`
	Customer *Customer `valid:"" jsonp:""`
	Card *Card `valid:"" jsonp:""`
	Account *PayZenAccount `valid:"Required"`
	TokenOperation bool 
	TransactionUuid string `jsonp:""`
	VadsTransId string`jsonp:"payzen_vads_trans_id"`
	ValidationType PayZenPaymentValidationType	`jsonp:"payzen_validation_type"`
	
	PaymentType PaymentType `jsonp:""`

	// pagarme
	PostbackUrl string `valid:"" jsonp:""` 
	SoftDescriptor string `valid:"" jsonp:""` 	
	//Metadata map[string]interface{} `valid:"" jsonp:"pagarme_metadata"` 
	SaveBoletoAtPath string
	BoletoFine *BoletoFine `jsonp:""`
	BoletoInterest *BoletoInterest `jsonp:""`
	BoletoDiscount *BoletoDiscount `jsonp:""`
	WebhookUrl string `jsonp:""`

	//PicPay
  ReturnUrl string `json:"" jsonp:"picpay_return_url"`
  Plugin string `json:"" jsonp:"" picpay_plugin`
  AdditionalInfo map[string]interface{} `json:"" jsonp:""`
  
  ExpiresAt time.Time `json:"" jsonp:"expires_at"`

	BoletoOnline BoletoOnlineTipo `valid:"" `	
	//dalay para pagemento do boleto (válido apenas para itaú. O bradesco deve ser configurado na plataforma)
	BoletoOnlineDaysDalay int `jsonp:""`// Obs.: Não suportado para boletos online Bradesco.

	BoletoOnlineTexto string  `jsonp:"payzen_boleto_text"`
	BoletoOnlineTexto2 string `jsonp:"payzen_boleto_tex2"`
	BoletoOnlineTexto3 string `jsonp:"payzen_boleto_text3"`

	BoletoInstructions string `jsonp:"boleto_instructions"`

}

type PaymentCapture struct {
	TransactionUuids string `valid:"Required" jsonp:""`
	Commission float64 
	Account *PayZenAccount `valid:"Required"`
}

type PaymentToken struct {
	Token string `valid:"Required" jsonp:""`
	Account *PayZenAccount `valid:"Required"`
}

type PaymentFind struct {
	//TransactionId string `jsonp:""`
	TransactionUuid string `jsonp:""`
	OrderId string `jsonp:""`
	SubscriptionId string `jsonp:""`
	Amount float64 `valid:"" jsonp:""`
	Token string `jsonp:""`
	AuthorizationId string `jsonp:""`
	Account *PayZenAccount `valid:"Required"`
	
	PaymentType PaymentType `jsonp:""` // picpay, pix
}

func NewSubscriptionWithShopId(shopId string, mode string, cert string) *Subscription {
	subscription := new(Subscription)
	subscription.Account = &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
	return subscription
}

func NewSubscription() *Subscription {
	subscription := new(Subscription)
	return subscription
}

func NewPaymentFindWithShopId(shopId string, mode string, cert string) *PaymentFind {
	find := new(PaymentFind)
	find.Account = &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
	return find
}

func NewPaymentFind() *PaymentFind {
	find := new(PaymentFind)
	find.Account = &PayZenAccount{ }
	return find
}

func NewPaymentTokenWithShopId(shopId string, mode string, cert string) *PaymentToken {
	tokenPayment := new(PaymentToken)
	tokenPayment.Account = &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
	return tokenPayment
}

func NewPaymentToken() *PaymentToken {
	tokenPayment := new(PaymentToken)
	tokenPayment.Account = &PayZenAccount{  }
	return tokenPayment	
}

func NewPaymentCaptureWithShopId(shopId string, mode string, cert string) *PaymentCapture {
	capturePayment := new(PaymentCapture)
	capturePayment.Account = &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
	return capturePayment
}

func NewPaymentCapture() *PaymentCapture {
	capturePayment := new(PaymentCapture)
	capturePayment.Account = &PayZenAccount{  }
	return capturePayment
}

func NewPaymentWithShopId(shopId string, mode string, cert string) *Payment {
	payment := new(Payment)
	payment.Account = &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
	payment.Customer = new(Customer)
	payment.Card = new(Card)
	payment.Installments = 1
	payment.ValidationType = Automatica
	return payment
}

func NewPayment() *Payment {
	payment := new(Payment)
	payment.Customer = new(Customer)
	payment.Card = new(Card)
	payment.Installments = 1
	return payment
}

func NewPayZenAccount(shopId string, mode string, cert string) *PayZenAccount {
	return &PayZenAccount{ ShopId: shopId, Mode: mode, Cert: cert }
}

type TokenResult struct {
	Token string `jsonp:""`
	Number string `jsonp:""`
	Brand string `jsonp:""`
	CreationDate time.Time 
	CancellationDate time.Time
	Cancelled bool `jsonp:""`
	Active bool `jsonp:""`
	NotFound bool `jsonp:""`
}

type SubscriptionResult struct {
	SubscriptionId string `jsonp:""`

	PastPaymentsNumber int64 `jsonp:""`
	TotalPaymentsNumber int64 `jsonp:""`
	EffectDate time.Time `jsonp:""`
	CancelDate time.Time `jsonp:""`
	InitialAmountNumber int64 `jsonp:""`
	InitialAmount int64 `jsonp:""`
	Rule string `jsonp:""`
	Description string `jsonp:""`

	Active bool `jsonp:""`
	Cancelled bool `jsonp:""`
	Started bool `jsonp:""`

	Token string `jsonp:""`
}


type TransactionItemResult struct {
	TransactionUuid string `jsonp:""`
	TransactionId string `jsonp:""`
	AuthorizationId string `jsonp:""`
	CancellationId string `jsonp:""`
	
	PagarmeStatus PagarmeStatus `jsonp:"pagarme_status"`
	PicPayStatus PicPayStatus `jsonp:"picpay_status"`
	TransactionStatus TransactionStatus `jsonp:"payzen_status"`
	PayZenV4Status string `jsonp:"payzen_v4_status"`
	AsaasStatus AsaasStatus `jsonp:"asaas_status"`

	TransactionStatusLabel string `jsonp:""`
	ExternalTransactionId string `jsonp:""`
	Nsu string `jsonp:""`
	Amount float64 `jsonp:""`
	ExpectedCaptureDate time.Time `jsonp:""`
	CreationDate time.Time `jsonp:""`
	Status PaymentStatus `jsonp:""`
	StatusLabel PaymentStatusLabel `jsonp:""`	

	Platform Gateway `jsonp:""`
}

func NewTransactionItemResult(platform Gateway) *TransactionItemResult {
	return &TransactionItemResult{ Platform: platform }
}

func (this *TransactionItemResult) isPagarme() bool {
	return this.Platform == GatewayPagarme
}

func (this *TransactionItemResult) isPayZen() bool {
	return this.Platform == GatewayPayZen
}

func (this *TransactionItemResult) isPicPay() bool {
	return this.Platform == GatewayPicPay
}

func (this *TransactionItemResult) isAsaas() bool {
	return this.Platform == GatewayAsaas
}

func (this *TransactionItemResult) BuildStatus() {

	switch this.TransactionStatus {
		case Initial:
			this.Status = PaymentInitial
			this.StatusLabel = PaymentInitialLabel
			break
		case NotCreated:
			this.Status = PaymentInitial
			this.StatusLabel = PaymentInitialLabel
			break
		case Authorised:
			if this.isPayZen() {			
				this.Status = PaymentPaid
				this.StatusLabel = PaymentPaidLabel
			} else {
				this.Status = PaymentWaitingPayment
				this.StatusLabel = PaymentWaitingPaymentLabel
			}
			break
		case AuthorisedToValidate:
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
			break
		case WaitingAuthorisation:
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
			break
		case WaitingAuthorisationToValidate:
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
			break
		case Refused:
			this.Status = PaymentRefused
			this.StatusLabel = PaymentRefusedLabel
			break
		case Captured:
			this.Status = PaymentPaid
			this.StatusLabel = PaymentPaidLabel
			break
		case Cancelled:
			this.Status = PaymentCancelled
			this.StatusLabel = PaymentCancelledLabel
			break			
		case Expired:
			this.Status = PaymentExpired
			this.StatusLabel = PaymentCancelledLabel
			break			
		case UnderVerification:
			this.Status = PaymentOther
			this.StatusLabel = PaymentOtherLabel
			break						
		case PartiallyAuthorised:
			this.Status = PaymentOther
			this.StatusLabel = PaymentOtherLabel
			break						
		case Refunded:
			this.Status = PaymentRefound
			this.StatusLabel = PaymentRefoundLabel
			break						
		case Created:
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
			break			
		case Chargeback:
			this.Status = PaymentChargeback
			this.StatusLabel = PaymentChargebackLabel
			break			
		case WaitingPayment:
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
			break			
		case PendingRefund:
			this.Status = PaymentOther
			this.StatusLabel = PaymentOtherLabel
			break			
		case Analyzing:
			this.Status = PaymentOther
			this.StatusLabel = PaymentOtherLabel
			break			
		case PendingReview:
			this.Status = PaymentOther
			this.StatusLabel = PaymentOtherLabel
			break						
		case Success:
			this.Status = PaymentSuccess
			this.StatusLabel = PaymentSuccessLabel
			break						
		default:
			this.Status = PaymentError
			this.StatusLabel = PaymentErrorLabel
			break						
	}	
}

type PaymentResultValidationError struct {
	Error bool `jsonp:""`
	Message string `jsonp:""`
	ValidationErrors map[string]string `jsonp:""`
}

func (this *PaymentResultValidationError) AddError(key string, value string) *PaymentResultValidationError{
	if this.ValidationErrors == nil {
		this.ValidationErrors = make(map[string]string)
	}

	this.ValidationErrors[key] = value
	return this
}

type Balance struct {
	WaitingFunds float64 `jsonp:""`
	Available float64 `jsonp:""`
	Transferred float64 `jsonp:""`  	
}

func NewBalance() *Balance{
	return &Balance{}
}

type Movement struct {
	Id string `jsonp:""`
	Amount float64 `jsonp:""`
	Type string `jsonp:""`
	Date string `jsonp:""`  	
	Status string `jsonp:""`
	TransactionId string `jsonp:""`	
}

func NewMovement() *Movement{
	return &Movement{}
}

type TransferResult struct {
  Id string `jsonp:""` 
  DateCreated string `jsonp:""` 
  Status TransferStatus `jsonp:""` 
  EffectiveDate string `jsonp:""` 
  Type string `jsonp:""` 
  Value float64 `jsonp:""` 
  NetValue float64 `jsonp:""` 
  TransferFee float64 `jsonp:""` 
  ScheduleDate string `jsonp:""` 
  Authorized bool `jsonp:""` 
  TransactionReceiptUrl string `jsonp:""` 
  TransactionId string `jsonp:""`
}

func NewTransferResult() *TransferResult {
	return &TransferResult{}
}

type PaymentResult struct {
	Error bool `jsonp:""`
	Message string `jsonp:""`
	Request string `jsonp:""`
	Response string `jsonp:""`

	CaptureRequest string 
	CaptureResponse string

	PagarmeStatus PagarmeStatus `jsonp:"pagarme_status"`
	PicPayStatus PicPayStatus `jsonp:"picpay_status"`
	TransactionStatus TransactionStatus `jsonp:"status"`
	PayZenV4Status string `jsonp:"payzen_v4_status"`
	AsaasStatus AsaasStatus `jsonp:"asaas_status"`
	
	Status PaymentStatus `jsonp:""`
	StatusLabel PaymentStatusLabel `jsonp:""`

	//v4

	TransactionStatusLabel string `jsonp:""`

	TransactionId string `jsonp:""`
	TransactionUuid string `jsonp:""`
	//ExternalTransactionId string

  ResponseCode string `jsonp:""`
  ResponseCodeDetail string `jsonp:""`

  //v4
  ErroCode string `jsonp:""`
  ErroCodeDetail string `jsonp:""`

	BoletoUrl string `jsonp:""`
	BoletoNumber string `jsonp:""`

	SubscriptionId string `jsonp:""`

	InstallmentId string `jsonp:""`
	InstallmentCount int64 `jsonp:""`
	Amount float64 `jsonp:""`

	TokenInfo *TokenResult `jsonp:""`

	PaymentNotFound bool `jsonp:""`
	SubscriptionInvalid bool `jsonp:""`
	SubscriptionIdNotFound bool `jsonp:""`
	PaymentTokenNotFound bool `jsonp:""`

	Transactions []*TransactionItemResult `jsonp:""`

	SubscriptionInfo *SubscriptionResult `jsonp:""`

	ValidationErrors map[string]string `jsonp:""`

	Plan *Plan `jsonp:""`

	Platform Gateway `jsonp:""`
	Nsu string	`jsonp:""`
	BoletoOutputContent []byte	`json:"-"`
	BoletoFileName string `json:"-"`

	QrCode string `jsonp:"qrcode"`
	QrCodeUrl string `jsonp:"qrcode_url"`
	QrPayload string `jsonp:"qrcode_payload"`
	QrExpirationDate string `jsonp:"qrcode_expiration_date"`
	PaymentUrl string `jsonp:"payment_url"`

	InvoiceUrl string `jsonp:""` // url da fatura
	PaymentLink string `jsonp:""` // identificador do link de pagamento

  AuthorizationId string `jsonp:"picpay_authorization_id"`
  CancellationId string `jsonp:"picpay_cancellation_id"`
  ReferenceId string `jsonp:"picpay_reference_id"`

  IsPagarme bool `jsonp:"is_pagarme"`
  IsPicPay bool `jsonp:"is_picpay"`
  IsPayZen bool `jsonp:"is_payzen"`
  IsAsaas bool `jsonp:"is_asaas"`

  Balance *Balance `jsonp:""`
  Movements []*Movement `jsonp:""`
  Transfers []*TransferResult `jsonp:""`

  Customer *Customer `jsonp:""`
}

func (this *PaymentResult) isPagarme() bool {
	return this.Platform == GatewayPagarme
}

func (this *PaymentResult) isPayZen() bool {
	return this.Platform == GatewayPayZen
}

func (this *PaymentResult) isPicPay() bool {
	return this.Platform == GatewayPicPay
}

func (this *PaymentResult) isAsaas() bool {
	return this.Platform == GatewayAsaas
}

func (this *PaymentResult) BuildStatus() {

	this.IsPayZen = this.isPayZen()
	this.IsPicPay = this.isPicPay()
	this.IsPagarme = this.isPagarme()
	this.IsPagarme = this.isAsaas()

	switch this.TransactionStatus {
		case Initial:
			this.Status = PaymentInitial
			this.StatusLabel = PaymentInitialLabel
			break
		case NotCreated:
			this.Status = PaymentInitial
			this.StatusLabel = PaymentInitialLabel
			break
		case Authorised:
			if this.IsPayZen {			
				this.Status = PaymentPaid
				this.StatusLabel = PaymentPaidLabel
			} else {
				this.Status = PaymentWaitingPayment
				this.StatusLabel = PaymentWaitingPaymentLabel
			}
			break
		case AuthorisedToValidate:
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
			break
		case WaitingAuthorisation:
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
			break
		case WaitingAuthorisationToValidate:
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
			break
		case Refused:
			this.Status = PaymentRefused
			this.StatusLabel = PaymentRefusedLabel
			break
		case Captured:
			this.Status = PaymentPaid
			this.StatusLabel = PaymentPaidLabel
			break
		case Cancelled:
			this.Status = PaymentCancelled
			this.StatusLabel = PaymentCancelledLabel
			break			
		case Expired:
			this.Status = PaymentExpired
			this.StatusLabel = PaymentCancelledLabel
			break			
		case UnderVerification:
			this.Status = PaymentOther
			this.StatusLabel = PaymentOtherLabel
			break						
		case PartiallyAuthorised:
			this.Status = PaymentOther
			this.StatusLabel = PaymentOtherLabel
			break						
		case Refunded:
			this.Status = PaymentRefound
			this.StatusLabel = PaymentRefoundLabel
			break						
		case Created:
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
			break			
		case Chargeback:
			this.Status = PaymentChargeback
			this.StatusLabel = PaymentChargebackLabel
			break			
		case WaitingPayment:
			this.Status = PaymentWaitingPayment
			this.StatusLabel = PaymentWaitingPaymentLabel
			break			
		case PendingRefund:
			this.Status = PaymentOther
			this.StatusLabel = PaymentOtherLabel
			break			
		case Analyzing:
			this.Status = PaymentOther
			this.StatusLabel = PaymentOtherLabel
			break			
		case PendingReview:
			this.Status = PaymentOther
			this.StatusLabel = PaymentOtherLabel
			break						
		case Success:
			this.Status = PaymentSuccess
			this.StatusLabel = PaymentSuccessLabel
			break						
		default:
			this.Status = PaymentError
			this.StatusLabel = PaymentErrorLabel
			break						
	}	
}

func NewPaymentResult() *PaymentResult {
	result := new(PaymentResult)
	result.TokenInfo = new(TokenResult)
	result.Transactions = []*TransactionItemResult{}
	result.SubscriptionInfo = new(SubscriptionResult)
	result.ValidationErrors = make(map[string]string)
	result.Movements = []*Movement{}
	result.Transfers = []*TransferResult{}
	result.Platform = GatewayPayZen
	return result
}

func NewPaymentResultWithErrors(errors map[string]string) *PaymentResult {
	result := NewPaymentResult()
	result.Error = true
	result.ValidationErrors = errors
	return result
}

func NewPaymentResultValidationError(errors map[string]string) *PaymentResultValidationError {
	result := new(PaymentResultValidationError)
	result.Error = true
	result.ValidationErrors = errors
	return result
}

func NewPaymentResultValidationErrorWithErrorKey(message string, key string, value string) *PaymentResultValidationError {
	result := new(PaymentResultValidationError)
	result.Error = true	
	return result.AddError(key, value)
}

func NewPaymentResultValidationErrorWithMessage(errors map[string]string, message string) *PaymentResultValidationError {
	result := new(PaymentResultValidationError)
	result.Error = true
	result.Message = message
	result.ValidationErrors = errors
	return result
}