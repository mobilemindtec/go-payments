package api

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
