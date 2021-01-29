package gopayments

import (
  "github.com/mobilemindtec/go-payments/pickpay"
  _ "github.com/mobilemindtec/go-utils/app/util"
  "github.com/mobilemindtec/go-utils/support"
	_ "github.com/satori/go.uuid"
	_ "github.com/go-redis/redis"
	"testing"
	"time"
	"fmt"
	_ "os"
  "encoding/json"
  "io/ioutil"	
)

var (

	Token = ""
	SallerToken = ""

)

func init(){
  file, err := ioutil.ReadFile("../certs.json")
  if err != nil {
      fmt.Printf("error on open file ../certs.json: %v\n", err)
      return
  }

  data := make(map[string]interface{})
  
  err = json.Unmarshal(file, &data)
  if err != nil {
      fmt.Printf("JSON error: %v\n", err)
      return
  }  

  jsonParser := new(support.JsonParser)

  clienteData := jsonParser.GetJsonObject(data, "mobilemind")
  pickpayData := jsonParser.GetJsonObject(clienteData, "pickpay")
  
  Token = jsonParser.GetJsonString(pickpayData, "token")
  SallerToken = jsonParser.GetJsonString(pickpayData, "sallerToken")

  fmt.Printf("init pickpay token = %v, sallerToken = %v", Token, SallerToken)
}

//go test  github.com/mobilemindtec/go-payments/tests -run TestPickPayCreateTransaction

func TestPickPayCreateTransaction(t *testing.T) {
		
	Pickpay := pickpay.NewPickPay("pt-BR", Token, SallerToken)
	Pickpay.Debug = true

	request := pickpay.NewPickPayTransactionRequest()

	request.Buyer.FirstName = "Ricardo"
	request.Buyer.LastName = "Bocchi"
	request.Buyer.Document = "83361855004"
	request.Buyer.Email = "ricardobocchi@gmail.com"
	request.Buyer.Phone = "+5554999767081"

	request.ReferenceId = "000001"
	request.CallbackUrl = fmt.Sprintf("https://portal.appmobloja.com.br/gateway/pickpay/postback/%v", request.ReferenceId)
	request.ReturnUrl = fmt.Sprintf("https://portal.appmobloja.com.br/gateway/pickpay/success/%v", request.ReferenceId)
	request.Value = "5"
	//request.Plugin = 
	//request.AdditionalInfo = 
	request.ExpiresAt	 = time.Now().Add(time.Duration(time.Hour * 48))
	

	result, err := Pickpay.CreateTransaction(request)
	

  if err != nil {
  	t.Errorf("Erro ao criar transacao: %v", err)
  }else{

  	t.Log(fmt.Sprintf("result = %v", result))
  }
}

func TestPickPayCheckStatus(t *testing.T) {
		
	Pickpay := pickpay.NewPickPay("pt-BR", Token, SallerToken)
	Pickpay.Debug = true

	result, err := Pickpay.CheckStatus("000001")
	

  if err != nil {
  	t.Errorf("Erro ao verificar status: %v", err)
  }else{

  	if result.Transaction.StatusText != "created" {
  		t.Errorf("Status esperado: created, encontrado %v", result.Transaction.StatusText)
  		return
  	}

  	if result.Transaction.PickPayStatus != pickpay.PickPayCreated {
  		t.Errorf("Status esperado: PickPayCreated, encontrado %v", result.Transaction.PickPayStatus)
  		return
  	}

  	t.Log(fmt.Sprintf("result = %v", result))
  }
}

func TestPickPayCheckCancel(t *testing.T) {
		
	Pickpay := pickpay.NewPickPay("pt-BR", Token, SallerToken)
	Pickpay.Debug = true

	result, err := Pickpay.Cancel("000001", "")
	

  if err != nil {
  	t.Errorf("Erro ao verificar status: %v", err)
  }else{

  	if result.Transaction.StatusText != "created" {
  		t.Errorf("Status esperado: created, encontrado %v", result.Transaction.StatusText)
  		return
  	}

  	if result.Transaction.PickPayStatus != pickpay.PickPayCreated {
  		t.Errorf("Status esperado: PickPayCreated, encontrado %v", result.Transaction.PickPayStatus)
  		return
  	}

  	t.Log(fmt.Sprintf("result = %v", result))
  }


}