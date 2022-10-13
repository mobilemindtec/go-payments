package gopayments

import (
  "github.com/mobilemindtec/go-payments/payzen/v4"
  "github.com/mobilemindtec/go-utils/support"
  "github.com/satori/go.uuid"
  "github.com/go-redis/redis"
  "encoding/json"
  "io/ioutil"
  "testing"
  "net"
  "fmt"
  "os"
)

var (
  client *redis.Client
    
  //payzen v4
  Authentication *v4.Authentication
  ApiMode = v4.Test  

  // pagarme
  ApiKey = ""
  CryptoKey = ""
  

  // payzen soap
  Mode = "TEST"
  ShopId = ""
  Cert = ""   

  // picpay
  Token = ""
  SallerToken = ""   

  // asaas
  AsaasAccessToken = ""
)

func init(){
  file, err := ioutil.ReadFile("../certs.json")
  if err != nil {
      panic(fmt.Sprintf("error on open file ../certs.json: %v\n", err))
      return
  }

  data := make(map[string]interface{})
  
  err = json.Unmarshal(file, &data)
  if err != nil {
      panic(fmt.Sprintf("error load configs json: %v\n", err))
      return
  }  

  jsonParser := new(support.JsonParser)

  // payzen v4
  mobilemindObj := jsonParser.GetJsonObject(data, "mobilemind")
  payzenObj := jsonParser.GetJsonObject(mobilemindObj, "payzen")
  v4Data := jsonParser.GetJsonObject(payzenObj, "v4")
  username := jsonParser.GetJsonString(v4Data, "username")
  v4Pwd := jsonParser.GetJsonObject(v4Data, "password")
  passwordTest := jsonParser.GetJsonString(v4Pwd, "test")
  passwordProd := jsonParser.GetJsonString(v4Pwd, "prod")

  mode, _ := payzenObj["mode"].(string)

  if mode == "PRODUCTION" {
    ApiMode = v4.Prod
  }

  Authentication = v4.NewAuthentication(username, passwordProd, passwordTest)

  // payzen soap  
  certObj := jsonParser.GetJsonObject(payzenObj, "cert")
  Mode, _ = payzenObj["mode"].(string)
  ShopId, _ = payzenObj["shop_id"].(string)
  Cert, _ = certObj["test"].(string)  

  // pagarme
  pagarmeObj := jsonParser.GetJsonObject(mobilemindObj, "pagarme")
  ApiKey = jsonParser.GetJsonString(pagarmeObj, "api_key")
  CryptoKey = jsonParser.GetJsonString(pagarmeObj, "crypto_key")

  // picpay
  picpayObj := jsonParser.GetJsonObject(mobilemindObj, "picpay")
  Token = jsonParser.GetJsonString(picpayObj, "token")
  SallerToken = jsonParser.GetJsonString(picpayObj, "sallerToken")

  // asaas

  asaasObj := jsonParser.GetJsonObject(mobilemindObj, "asaas")
  AsaasAccessToken = jsonParser.GetJsonString(asaasObj, "api_key")

  fmt.Printf("init picpay token = %v, sallerToken = %v", Token, SallerToken)
  fmt.Printf("init payzen data: Mode = %v, ShopId = %v, Cert = %v", Mode, ShopId, Cert)
  fmt.Printf("init pagarme toApiKey = %v, CryptoKey = %v", Token, CryptoKey)  
  fmt.Printf("init payzen v4 data: Mode = %v, username = %v, Password = %v", ApiMode, username, passwordTest)
}

func GetHostIp() (string) {

    netInterfaceAddresses, err := net.InterfaceAddrs()

    if err != nil { return "" }

    for _, netInterfaceAddress := range netInterfaceAddresses {

        networkIp, ok := netInterfaceAddress.(*net.IPNet)

        if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {

            ip := networkIp.IP.String()

            //fmt.Println("Resolved Host IP: " + ip)

            return ip
        }
    }
    return ""
}

func GenUUID() string {
  id := uuid.NewV4()
  return id.String()
}

func setup(){
  client = redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
  })
}

func shutdown() {

}

func TestMain(m *testing.M) {
  setup()
  code := m.Run()
  shutdown()
  os.Exit(code)
}