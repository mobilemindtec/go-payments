package gopayments

import (
	"fmt"
	"github.com/mobilemindtec/go-payments/api"
	"github.com/mobilemindtec/go-payments/picpay"
	"testing"
	"time"
)

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPicPayCreateTransaction
func TestPicPayCreateTransaction(t *testing.T) {

	Picpay := picpay.NewPicPay("pt-BR", Token, SallerToken)
	Picpay.Debug = true

	request := picpay.NewPicPayTransactionRequest()

	request.Buyer.FirstName = "Ricardo"
	request.Buyer.LastName = "Bocchi"
	request.Buyer.Document = "83361855004"
	request.Buyer.Email = "ricardobocchi@gmail.com"
	request.Buyer.Phone = "+5554999767081"

	request.ReferenceId = GenUUID()
	request.CallbackUrl = fmt.Sprintf("https://portal.appmobloja.com.br/gateway/picpay/postback/%v", request.ReferenceId)
	request.ReturnUrl = fmt.Sprintf("https://portal.appmobloja.com.br/gateway/picpay/success/%v", request.ReferenceId)
	request.Value = "5"
	//request.Plugin =
	//request.AdditionalInfo =
	request.ExpiresAt = time.Now().Add(time.Duration(time.Hour * 48))

	result, err := Picpay.CreateTransaction(request)

	if err != nil {
		t.Errorf("Erro ao criar transacao: %v", err)
	} else {

		t.Log(fmt.Sprintf("result = %v", result))
	}

	CacheClient.Set("ReferenceId", request.ReferenceId, 0)
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPicPayCheckStatus
func TestPicPayCheckStatus(t *testing.T) {

	Picpay := picpay.NewPicPay("pt-BR", Token, SallerToken)
	Picpay.Debug = true

	referenceId, _ := CacheClient.Get("ReferenceId").Result()
	result, err := Picpay.CheckStatus(referenceId)

	if err != nil {
		t.Errorf("Erro ao verificar status: %v", err)
	} else {

		if result.Transaction.StatusText != "created" {
			t.Errorf("Status esperado: created, encontrado %v", result.Transaction.StatusText)
			return
		}

		if result.Transaction.PicPayStatus != api.PicPayCreated {
			t.Errorf("Status esperado: PicPayCreated, encontrado %v", result.Transaction.PicPayStatus)
			return
		}

		t.Log(fmt.Sprintf("result = %v", result))
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPicPayCheckCancel
func TestPicPayCheckCancel(t *testing.T) {

	Picpay := picpay.NewPicPay("pt-BR", Token, SallerToken)
	Picpay.Debug = true

	referenceId, _ := CacheClient.Get("ReferenceId").Result()
	result, err := Picpay.Cancel(referenceId, "")

	if err != nil {
		t.Errorf("Erro ao verificar status: %v", err)
	} else {

		if result.Transaction.StatusText != "created" {
			t.Errorf("Status esperado: created, encontrado %v", result.Transaction.StatusText)
			return
		}

		if result.Transaction.PicPayStatus != api.PicPayCreated {
			t.Errorf("Status esperado: PicPayCreated, encontrado %v", result.Transaction.PicPayStatus)
			return
		}

		t.Log(fmt.Sprintf("result = %v", result))
	}

}
