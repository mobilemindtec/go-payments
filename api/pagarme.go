package api

type PagarmeV5Status string

const (
	PagarmeV5None PagarmeV5Status = "none"

	// Card
	// 	Authorized pending capture
	PagarmeV5AuthorizedPendingCapture PagarmeV5Status = "authorized_pending_capture"
	// 	Not allowed
	PagarmeV5NotAuthorized PagarmeV5Status = "not_authorized"
	// 	Captured
	PagarmeV5Captured PagarmeV5Status = "captured"
	// 	Partially Captured
	PagarmeV5PartialCapture PagarmeV5Status = "partial_capture"
	// 	Waiting for capture
	PagarmeV5WaitingCapture PagarmeV5Status = "waiting_capture"
	// 	Reversed
	PagarmeV5Refunded PagarmeV5Status = "refunded"
	// 	Canceled
	PagarmeV5Voided PagarmeV5Status = "voided"
	// 	Partially reversed
	PagarmeV5PartialRefunded PagarmeV5Status = "partial_refunded"
	// 	Partially canceled
	PagarmeV5PartialVoid PagarmeV5Status = "partial_void"
	// 	Cancellation error
	PagarmeV5ErrorOnVoiding PagarmeV5Status = "error_on_voiding"
	// 	Error in refund
	PagarmeV5ErrorOnRefunding PagarmeV5Status = "error_on_refunding"
	// 	Awaiting cancellation
	PagarmeV5WaitingCancellation PagarmeV5Status = "waiting_cancellation"
	// 	With error
	PagarmeV5WithError PagarmeV5Status = "with_error"
	// 	Failure
	PagarmeV5Failed PagarmeV5Status = "failed"

	// Pix
	//Aguardando pagamento
	PagarmeV5WaitingPayment PagarmeV5Status = "waiting_payment"
	//Paid out
	PagarmeV5Paid PagarmeV5Status = "paid"
	//Aguardando estorno
	PagarmeV5PendingRefund PagarmeV5Status = "pending_refund"
	//Estornado
	//Refunded PagarmeV5Status = "refunded"
	//With error
	//WithError PagarmeV5Status = "with_error"
	//Failure
	//Failed PagarmeV5Status = "failed"

	// Boleto
	//Generated
	PagarmeV5Generated PagarmeV5Status = "generated"
	//Home
	PagarmeV5Viewed PagarmeV5Status = "viewed"
	//Payment to less
	PagarmeV5Underpaid PagarmeV5Status = "underpaid"
	//Highest payment
	PagarmeV5Overpaid PagarmeV5Status = "overpaid"
	//Paid out
	//PagarmeV5Paid PagarmeV5Status = "paid"
	//Canceled
	//Voided PagarmeV5Status = "voided"
	//With error
	//WithError PagarmeV5Status = "with_error"
	//Failure
	//Failed PagarmeV5Status = "failed"
	//Boleto ainda está em etapa de criação
	PagarmeV5Processing PagarmeV5Status = "processing"
)
