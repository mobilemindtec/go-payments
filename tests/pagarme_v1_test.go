package gopayments

import (
	"fmt"
	"testing"
	"time"

	"github.com/mobilemindtec/go-payments/api"
	pagarme "github.com/mobilemindtec/go-payments/pagarme/v1"
)

func pagarmeFillCard(card *pagarme.Card) {
	card.Number = "4901720080344448"
	card.HolderName = "Aardvark Silva"
	card.ExpirationDate = "1225"
	card.Cvv = "314"
}

func pagarmeFillPayments(payment *pagarme.Payment) {
	customer := new(pagarme.Customer)

	pagarmefillCustomer(customer)

	//payment.Amount = 10
	payment.Installments = 1
	payment.Customer = customer
	payment.PaymentMethod = api.PaymentTypeCreditCard
	//payment.PostbackUrl = "https://mobilemind.com.br"
	payment.SoftDescriptor = "Mobile Mind"
	//payment.Metadata = ""
	payment.Capture = true
	//payment.BoletoExpirationDate
	//payment.BoletoInstructions
	//payment.CardId
	payment.CardHolderName = "Ricardo Bocchi"
	payment.CardExpirationDate = "0925"
	payment.CardNumber = "4056769270964567"
	payment.CardCvv = "123"
	//payment.CardHash
	//payment.SplitRules
}

func pagarmefillCustomer(customer *pagarme.Customer) {
	customer.Email = "ricardobocchi@gmail.com"
	customer.Name = "Ricardo Bocchi"
	customer.DocumentNumber = "83361855004"

	customer.Phone = new(pagarme.Phone)
	customer.Phone.Ddd = "054"
	customer.Phone.Number = "999767081"

	customer.Address = new(pagarme.Address)
	customer.Address.Neighborhood = "Botafogo"
	customer.Address.Street = "Vitoria"
	customer.Address.StreetNumber = "255"
	customer.Address.ZipCode = "95700540"
	customer.Address.City = "Bento Goncalves"
	customer.Address.State = "RS"
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeCreateCardHashKey
func TestPagarmeCreateCardHashKey(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	result, err := Pagarme.GetCardHashKey()

	if err != nil {
		t.Errorf("Erro ao criar card hash key: %v", err)
		return
	}

	if len(result.PublicKey) == 0 {
		t.Errorf("PublicKey is expected")
		return
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeEncryptCard
func TestPagarmeEncryptCard(t *testing.T) {

	Card := new(pagarme.Card)
	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	pagarmeFillCard(Card)

	result, err := Pagarme.EncryptCard(Card)

	if err != nil {
		t.Errorf("Erro ao encrypt card: %v", err)
		return
	}

	if len(result.Hash) == 0 {
		t.Errorf("card hash is expected")
		return
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeTokenCreate
func TestPagarmeTokenCreate(t *testing.T) {

	Card := new(pagarme.Card)
	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()
	pagarmeFillCard(Card)

	result, err := Pagarme.TokenCreate(Card)

	if err != nil {
		t.Errorf("Erro ao criar card cartão: %v", err)
		return
	}

	if len(result.CardResult.Id) == 0 {
		t.Errorf("card token is expected")
		return
	}

	client.Set("CardId", result.CardResult.Id, 0)
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmePaymentCreateWithCard
func TestPagarmePaymentCreateWithCard(t *testing.T) {

	payment := pagarme.NewPaymentWithCard(1)
	payment.PostbackUrl = "https://mobilemind.free.beeceptor.com/webhook/pagarme"

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	pagarmeFillPayments(payment)

	result, err := Pagarme.PaymentCreate(payment)

	if err != nil {
		t.Errorf("Erro ao create card payment: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Id == 0 {
			t.Errorf("Id cant be empty")
			return
		}

		if result.Status == api.PagarmeAuthorized {

			captureData := pagarme.NewCaptureData(fmt.Sprintf("%v", result.Id), 1)

			result, err := Pagarme.PaymentCapture(captureData)

			if err != nil {
				t.Errorf("Erro ao capture data: %v", err)
				return
			}

			if result.Status != api.PagarmeAuthorized {
				t.Errorf("status expected %v, returned %v", result.Status, api.PagarmeAuthorized)
				return
			}

		}

		client.Set("TransactionId", result.Id, 0)

	}

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmePaymentCreateWithBoleto
func TestPagarmePaymentCreateWithBoleto(t *testing.T) {

	payment := pagarme.NewPaymentWithBoleto(1)
	payment.PostbackUrl = "https://mobilemind.free.beeceptor.com/webhook/pagarme"

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	pagarmeFillPayments(payment)

	result, err := Pagarme.PaymentCreate(payment)

	if err != nil {
		t.Errorf("Erro ao create card payment: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Id == 0 {
			t.Errorf("Id cant be empty")
			return
		}

		if result.Status == api.PagarmeAuthorized {

			captureData := pagarme.NewCaptureData(fmt.Sprintf("%v", result.Id), 1)

			result, err := Pagarme.PaymentCapture(captureData)

			if err != nil {
				t.Errorf("Erro ao capture data: %v", err)
				return
			}

			if result.Status != api.PagarmeWaitingPayment {
				t.Errorf("status expected %v, returned %v", result.Status, api.PagarmeWaitingPayment)
				return
			}

		}

		client.Set("TransactionId", result.Id, 0)

	}

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmePaymentCreateWithPix
func TestPagarmePaymentCreateWithPix(t *testing.T) {

	payment := pagarme.NewPaymentWithPix(1)
	payment.SetPixExpirationDate(time.Now().AddDate(0, 0, 3))
	payment.AddPixAdditionalFields("Mobile Mind", "Test")
	payment.AddPixAdditionalFields("Mobile Mind 2", "Test 2")
	payment.SoftDescriptor = "Mobile Mind"
	payment.PostbackUrl = "https://mobilemind.free.beeceptor.com/webhook/pagarme"

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	pagarmefillCustomer(payment.Customer)

	result, err := Pagarme.PaymentCreate(payment)

	if err != nil {
		t.Errorf("Erro ao create card payment: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Id == 0 {
			t.Errorf("Id cant be empty")
			return
		}

		if result.Status == api.PagarmeAuthorized {

			captureData := pagarme.NewCaptureData(fmt.Sprintf("%v", result.Id), 1)

			result, err := Pagarme.PaymentCapture(captureData)

			if err != nil {
				t.Errorf("Erro ao capture data: %v", err)
				return
			}

			if result.Status != api.PagarmeWaitingPayment {
				t.Errorf("status expected %v, returned %v", result.Status, api.PagarmeAuthorized)
				return
			}

		}

		client.Set("TransactionId", result.Id, 0)

	}

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmePaymentStatus
func TestPagarmePaymentStatus(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	id, _ := client.Get("TransactionId").Int64()
	result, err := Pagarme.PaymentGet(fmt.Sprintf("%v", id))

	if err != nil {
		t.Errorf("Erro ao create card payment: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Status != api.PagarmeRefunded {
			t.Errorf("status expected %v, returned %v", result.Status, api.PagarmeRefunded)
			return
		}

	}

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmePaymentCancel
func TestPagarmePaymentCancel(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	id, _ := client.Get("TransactionId").Int64()

	result, err := Pagarme.PaymentRefund(fmt.Sprintf("%v", id), 10)

	if err != nil {
		t.Errorf("Erro ao create card payment: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Id == 0 {
			t.Errorf("Id cant be empty")
			return
		}

		if result.Status == api.PagarmeAuthorized {

			captureData := pagarme.NewCaptureData(fmt.Sprintf("%v", result.Id), 1)

			result, err := Pagarme.PaymentCapture(captureData)

			if err != nil {
				t.Errorf("Erro ao capture data: %v", err)
				return
			}

			if result.Status != api.PagarmeAuthorized {
				t.Errorf("status expected %v, returned %v", result.Status, api.PagarmeAuthorized)
				return
			}

		}

		client.Set("TransactionId", result.Id, 0)

	}

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmePlanoCreate
func TestPagarmePlanoCreate(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	plano := pagarme.NewPlano(fmt.Sprintf("My plan %v", time.Now().Unix()), 120)
	plano.SetCycle(pagarme.Monthly, 0, 5)

	result, err := Pagarme.PlanoCreate(plano)

	if err != nil {
		t.Errorf("Erro ao create plano: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Plano.Id == 0 {
			t.Errorf("plano id is expected")
			return
		}

		client.Set("PlanoId", result.Plano.Id, 0)

	}

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmePlanoUpdate
func TestPagarmePlanoUpdate(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	id, _ := client.Get("PlanoId").Int64()

	plano := pagarme.NewPlano(fmt.Sprintf("My plan %v", time.Now().Unix()), 110)
	plano.Id = id
	//plano.SetCycle(pagarme.Monthly, 0, 5)

	result, err := Pagarme.PlanoUpdate(plano)

	if err != nil {
		t.Errorf("Erro ao create plano: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Plano.Id == 0 {
			t.Errorf("plano id is expected")
			return
		}

	}

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmePlanoGet
func TestPagarmePlanoGet(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	id, _ := client.Get("PlanoId").Int64()

	result, err := Pagarme.PlanoGet(id)

	if err != nil {
		t.Errorf("Erro ao create plano: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Plano.Id == 0 {
			t.Errorf("plano id is expected")
			return
		}

	}

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeSubscriptionCreate
func TestPagarmeSubscriptionCreate(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	planId, _ := client.Get("PlanoId").Int64()
	cardId, _ := client.Get("CardId").Result()
	subscription := pagarme.NewSubscriptionWithCard(planId)
	subscription.CardId = cardId
	subscription.PostbackUrl = "https://mobilemind.free.beeceptor.com/webhook/pagarme"

	pagarmefillCustomer(subscription.Customer)

	result, err := Pagarme.SubscriptionCreate(subscription)

	if err != nil {
		t.Errorf("Erro ao create subscription: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Id == 0 {
			t.Errorf("Subscription id is expected")
			return
		}

		client.Set("SubscriptionId", result.Id, 0)

	}

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeSubscriptionUpdate
func TestPagarmeSubscriptionUpdate(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	planId, _ := client.Get("PlanoId").Int64()
	cardId, _ := client.Get("CardId").Result()
	subscriptionId, _ := client.Get("SubscriptionId").Int64()
	subscription := pagarme.NewSubscriptionWithCard(planId)
	subscription.CardId = cardId
	subscription.Id = subscriptionId
	subscription.PostbackUrl = "https://mobilemind.free.beeceptor.com/webhook/pagarme"

	pagarmefillCustomer(subscription.Customer)

	result, err := Pagarme.SubscriptionUpdate(subscription)

	if err != nil {
		t.Errorf("Erro ao create subscription: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Id == 0 {
			t.Errorf("Subscription id is expected")
			return
		}

	}

}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeSubscriptionGet
func TestPagarmeSubscriptionGet(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	id, _ := client.Get("SubscriptionId").Int64()

	result, err := Pagarme.SubscriptionGet(fmt.Sprintf("%v", id))

	if err != nil {
		t.Errorf("Erro ao get subscription: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if result.Plano.Id == 0 {
			t.Errorf("plano id is expected")
			return
		}
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeSubscriptionTransactionsGet
func TestPagarmeSubscriptionTransactionsGet(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	id, _ := client.Get("SubscriptionId").Int64()

	result, err := Pagarme.SubscriptionTransactionsGet(fmt.Sprintf("%v", id))

	if err != nil {
		t.Errorf("Erro ao list subscription transactions: %v", err)
	} else {
		//t.Log(fmt.Sprintf("result = %v", customer.Id))

		if !result.HasTransactions() {
			t.Errorf("transactions is expected")
			return
		}

		for _, x := range result.Transactions {
			t.Errorf("%v", x.StatusText)
		}

		if result.TransactionsCount() != 2 {
			t.Errorf("transactions count expected 2, but returned %v", result.TransactionsCount())
			return
		}

		if result.FirstTransaction().Status == api.PagarmePaid {
			t.Errorf("transactions paid expected")
			return
		}

		if result.FirstTransaction().Status == api.PagarmePaid {
			t.Errorf("transactions paid expected")
			return
		}

	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeSubscriptionSkip
func TestPagarmeSubscriptionSkip(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	id, _ := client.Get("SubscriptionId").Int64()

	_, err := Pagarme.SubscriptionSkip(fmt.Sprintf("%v", id), 1)

	if err != nil {
		t.Errorf("Erro ao get subscription: %v", err)
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeSubscriptionCancel
func TestPagarmeSubscriptionCancel(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	id, _ := client.Get("SubscriptionId").Int64()

	result, err := Pagarme.SubscriptionCancel(fmt.Sprintf("%v", id))

	if err != nil {
		t.Errorf("Erro ao get subscription: %v", err)
	}

	if result.Status == api.PagarmeCancelled {
		t.Errorf("transactions paid expected")
		return
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeCurrentBalance
func TestPagarmeCurrentBalance(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	result, err := Pagarme.CurrentBalance("re_ciiahjw06003a546eedfngbv8")

	if err != nil {
		t.Errorf("Erro ao get subscription: %v", err)
	}

	if result.Status == api.PagarmeCancelled {
		t.Errorf("transactions paid expected")
		return
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeMovements
func TestPagarmeMovements(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	filter := pagarme.NewFilter()
	result, err := Pagarme.Movements(filter)

	if err != nil {
		t.Errorf("Erro ao get subscription: %v", err)
		return
	}

	if !result.HasMovements() {
		t.Errorf("movement expected, but not has")
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeTransferList
func TestPagarmeTransferList(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	filter := pagarme.NewFilter()
	result, err := Pagarme.TransferList(filter)

	if err != nil {
		t.Errorf("Erro ao get subscription: %v", err)
		return
	}

	if !result.HasMovements() {
		t.Errorf("movement expected, but not has")
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeTransferCreate
func TestPagarmeTransferCreate(t *testing.T) {

	Pagarme := pagarme.NewPagarme("pt-BR", ApiKey, CryptoKey)
	Pagarme.SetDebug()

	bankAccount := pagarme.NewBankAccount(2020)
	tranfer := pagarme.NewTransfer(100, bankAccount)
	result, err := Pagarme.TransferCreate(tranfer)

	if err != nil {
		t.Errorf("Erro ao get subscription: %v", err)
		return
	}

	if !result.HasMovements() {
		t.Errorf("movement expected, but not has")
	}
}

// go test -v  github.com/mobilemindtec/go-payments/tests -run TestPagarmeWebhooTest
func TestPagarmeWebhooTest(t *testing.T) {

	requestBody := []byte(`{
    "object": "postback",
    "status": "pending_retry",
    "model": "transaction",
    "model_id": "1662527",
    "headers": "{\"Content-Type\":\"application/x-www-form-urlencoded\",\"X-PagarMe-Event\":\"transaction_status_changed\",\"X-Hub-Signature\":\"sha1=0c62a0b489e1138ef39ae71dece45be1c0e97c1e\",\"User-Agent\":\"PagarMe-Hookshot/1.0\"}",
    "payload": "id=1662527&fingerprint=a67597c98a493cc8b2c62ab018a553c19747e8a5&event=transaction_status_changed&old_status=waiting_payment&desired_status=paid&current_status=paid&object=transaction&transaction%5Bobject%5D=transaction&transaction%5Bstatus%5D=paid&transaction%5Brefuse_reason%5D=&transaction%5Bstatus_reason%5D=acquirer&transaction%5Bacquirer_response_code%5D=&transaction%5Bacquirer_name%5D=pagarme&transaction%5Bacquirer_id%5D=56f9d019decf72cc70055d58&transaction%5Bauthorization_code%5D=&transaction%5Bsoft_descriptor%5D=&transaction%5Btid%5D=1662527&transaction%5Bnsu%5D=1662527&transaction%5Bdate_created%5D=2017-06-28T17%3A36%3A52.808Z&transaction%5Bdate_updated%5D=2017-06-28T17%3A37%3A25.949Z&transaction%5Bamount%5D=15000&transaction%5Bauthorized_amount%5D=15000&transaction%5Bpaid_amount%5D=0&transaction%5Brefunded_amount%5D=0&transaction%5Binstallments%5D=1&transaction%5Bid%5D=1662527&transaction%5Bcost%5D=380&transaction%5Bcard_holder_name%5D=&transaction%5Bcard_last_digits%5D=&transaction%5Bcard_first_digits%5D=&transaction%5Bcard_brand%5D=&transaction%5Bcard_pin_mode%5D=&transaction%5Bpostback_url%5D=https%3A%2F%2Frequestb.in%2F10m5xva1&transaction%5Bpayment_method%5D=boleto&transaction%5Bcapture_method%5D=ecommerce&transaction%5Bantifraud_score%5D=&transaction%5Bboleto_url%5D=https%3A%2F%2Fpagar.me&transaction%5Bboleto_barcode%5D=1234%205678&transaction%5Bboleto_expiration_date%5D=2017-07-03T03%3A00%3A00.000Z&transaction%5Breferer%5D=api_key&transaction%5Bip%5D=177.63.194.231&transaction%5Bsubscription_id%5D=&transaction%5Bphone%5D=&transaction%5Baddress%5D=&transaction%5Bcustomer%5D=&transaction%5Bbilling%5D=&transaction%5Bshipping%5D=&transaction%5Bcard%5D=&transaction%5Bsplit_rules%5D=",
    "request_url": "https://requestb.in/10m5xva1",
    "retries": 0,
    "next_retry": null,
    "deliveries": [
      {
        "object": "postback_delivery",
        "status": "success",
        "status_reason": "http_status_code",
        "status_code": "200",
        "response_time": 228,
        "response_headers": "{\"date\":\"Wed, 28 Jun 2017 17:37:26 GMT\",\"content-type\":\"text/html; charset=utf-8\",\"transfer-encoding\":\"chunked\",\"connection\":\"close\",\"set-cookie\":[\"__cfduid=dd5481bfa0252320257fb1d3de05f19961498671446; expires=Thu, 28-Jun-18 17:37:26 GMT; path=/; domain=.requestb.in; HttpOnly\"],\"sponsored-by\":\"https://www.runscope.com\",\"via\":\"1.1 vegur\",\"strict-transport-security\":\"max-age=15552000\",\"x-content-type-options\":\"nosniff\",\"server\":\"cloudflare-nginx\",\"cf-ray\":\"376269f9ef4c0ed9-EWR\"}",
        "response_body": "ok",
        "date_created": "2017-06-28T17:37:26.033Z",
        "date_updated": "2017-06-28T17:37:26.266Z",
        "id": "pd_cj4haa8lt14slps730qybvjij"
      }
    ],
    "date_created": "2017-06-28T17:37:26.008Z",
    "date_updated": "2017-06-28T18:11:42.999Z",
    "signature": "sha1=0c62a0b489e1138ef39ae71dece45be1c0e97c1e",
    "id": "po_cj4haa8l4131bpi73glgzbnpp"
  }`)

	webkook := pagarme.NewDefaultWebhook()

	data, err := webkook.Parse(requestBody)

	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	if data.Response.Id != 1662527 {
		t.Errorf("Error is not pased: %v", data.Response.Id)
	}

}
