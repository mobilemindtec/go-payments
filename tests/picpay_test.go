package gopayments

import (
  "github.com/mobilemindtec/go-payments/picpay"
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
  picpayData := jsonParser.GetJsonObject(clienteData, "picpay")
  
  Token = jsonParser.GetJsonString(picpayData, "token")
  SallerToken = jsonParser.GetJsonString(picpayData, "sallerToken")

  fmt.Printf("init picpay token = %v, sallerToken = %v", Token, SallerToken)
}

//go test  github.com/mobilemindtec/go-payments/tests -run TestPicPayCreateTransaction

func TestPicPayCreateTransaction(t *testing.T) {
		
	Picpay := picpay.NewPicPay("pt-BR", Token, SallerToken)
	Picpay.Debug = true

	request := picpay.NewPicPayTransactionRequest()

	request.Buyer.FirstName = "Ricardo"
	request.Buyer.LastName = "Bocchi"
	request.Buyer.Document = "83361855004"
	request.Buyer.Email = "ricardobocchi@gmail.com"
	request.Buyer.Phone = "+5554999767081"

	request.ReferenceId = "000001"
	request.CallbackUrl = fmt.Sprintf("https://portal.appmobloja.com.br/gateway/picpay/postback/%v", request.ReferenceId)
	request.ReturnUrl = fmt.Sprintf("https://portal.appmobloja.com.br/gateway/picpay/success/%v", request.ReferenceId)
	request.Value = "5"
	//request.Plugin = 
	//request.AdditionalInfo = 
	request.ExpiresAt	 = time.Now().Add(time.Duration(time.Hour * 48))
	

	result, err := Picpay.CreateTransaction(request)
	

  if err != nil {
  	t.Errorf("Erro ao criar transacao: %v", err)
  }else{

  	t.Log(fmt.Sprintf("result = %v", result))
  }
}

func TestPicPayCheckStatus(t *testing.T) {
		
	Picpay := picpay.NewPicPay("pt-BR", Token, SallerToken)
	Picpay.Debug = true

	result, err := Picpay.CheckStatus("000001")
	

  if err != nil {
  	t.Errorf("Erro ao verificar status: %v", err)
  }else{

  	if result.Transaction.StatusText != "created" {
  		t.Errorf("Status esperado: created, encontrado %v", result.Transaction.StatusText)
  		return
  	}

  	if result.Transaction.PicPayStatus != picpay.PicPayCreated {
  		t.Errorf("Status esperado: PicPayCreated, encontrado %v", result.Transaction.PicPayStatus)
  		return
  	}

  	t.Log(fmt.Sprintf("result = %v", result))
  }
}

func TestPicPayCheckCancel(t *testing.T) {
		
	Picpay := picpay.NewPicPay("pt-BR", Token, SallerToken)
	Picpay.Debug = true

	result, err := Picpay.Cancel("000001", "")
	

  if err != nil {
  	t.Errorf("Erro ao verificar status: %v", err)
  }else{

  	if result.Transaction.StatusText != "created" {
  		t.Errorf("Status esperado: created, encontrado %v", result.Transaction.StatusText)
  		return
  	}

  	if result.Transaction.PicPayStatus != picpay.PicPayCreated {
  		t.Errorf("Status esperado: PicPayCreated, encontrado %v", result.Transaction.PicPayStatus)
  		return
  	}

  	t.Log(fmt.Sprintf("result = %v", result))
  }


}