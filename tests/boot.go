package gopayments

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-redis/redis"
	"github.com/mobilemindtec/go-payments/api"
	"github.com/mobilemindtec/go-payments/payzen/v4"
	"github.com/mobilemindtec/go-utils/support"
	"github.com/satori/go.uuid"
)

var (
	CacheClient *redis.Client

	//payzen v4
	Authentication *v4.Authentication
	ApiMode        = v4.Test

	// pagarme
	ApiKey    = ""
	CryptoKey = ""
	PublicKey = ""
	SecretKey = ""

	// payzen soap
	Mode   = "TEST"
	ShopId = ""
	Cert   = ""

	// picpay
	Token       = ""
	SallerToken = ""

	// asaas
	AsaasAccessToken = ""
	AsaasApiMode     = api.AsaasModeTest
)

func init() {

	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	logs.Debug("apppath = ", apppath, "file", file)

	certsFile := fmt.Sprintf("%v/certs.json", apppath)
	data, err := os.ReadFile(certsFile)
	if err != nil {
		panic(fmt.Sprintf("error on open file %v: %v\n", certsFile, err))
		return
	}

	jsonData := make(map[string]interface{})

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		panic(fmt.Sprintf("error load configs json: %v\n", err))
		return
	}

	jsonParser := new(support.JsonParser)

	// payzen v4
	mobilemindObj := jsonParser.GetJsonObject(jsonData, "mobilemind")
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
	pagarmeV1 := jsonParser.GetJsonObject(pagarmeObj, "v1")
	pagarmeV5 := jsonParser.GetJsonObject(pagarmeObj, "v5")
	ApiKey = jsonParser.GetJsonString(pagarmeV1, "api_key")
	CryptoKey = jsonParser.GetJsonString(pagarmeV1, "crypto_key")
	PublicKey = jsonParser.GetJsonString(pagarmeV5, "public_key")
	SecretKey = jsonParser.GetJsonString(pagarmeV5, "secret_key")

	// picpay
	picpayObj := jsonParser.GetJsonObject(mobilemindObj, "picpay")
	Token = jsonParser.GetJsonString(picpayObj, "token")
	SallerToken = jsonParser.GetJsonString(picpayObj, "sallerToken")

	// asaas

	asaasObj := jsonParser.GetJsonObject(mobilemindObj, "asaas")
	AsaasAccessToken = jsonParser.GetJsonString(asaasObj, "api_key")

	fmt.Printf("init picpay token = %v, sallerToken = %v\n", Token, SallerToken)
	fmt.Printf("init payzen data: Mode = %v, ShopId = %v, Cert = %v\n", Mode, ShopId, Cert)
	fmt.Printf("init pagarme v1 toApiKey = %v, CryptoKey = %v\n", Token, CryptoKey)
	fmt.Printf("init pagarme v5 public key = %v, secret key = %v\n", PublicKey, SecretKey)
	fmt.Printf("init payzen v4 data: Mode = %v, username = %v, Password = %v\n", ApiMode, username, passwordTest)
	fmt.Printf("init asaas: Token = %v\n", AsaasAccessToken)
}

func GetHostIp() string {

	netInterfaceAddresses, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

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

func setup() {
	CacheClient = redis.NewClient(&redis.Options{
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
