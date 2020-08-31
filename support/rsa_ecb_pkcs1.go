package support

import (
	_ "bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "encoding/base64"
	"encoding/pem"
	"errors"
	_ "fmt"
	"io"
	_ "io/ioutil"
	_ "os"
)

/*
func main() {
	var (
		plainText = "hello world"
	)

	fmt.Println("=== start generating RSA key pair")
	pubKey := bytes.NewBuffer([]byte{})
	priKey := bytes.NewBuffer([]byte{})
	xrsa.CreateKeys(pubKey, priKey, 2048)
	fmt.Println(pubKey.String())
	fmt.Println(priKey.String())
	// saveToFile("filename", []byte(priKey.String()))

	fmt.Println("=== start encrypting ===")
	encryptedData, err := RsaEncrypt([]byte(plainText), []byte(pubKey.String()))
	checkError(err)
	encryptedText := base64.StdEncoding.EncodeToString(encryptedData)
	fmt.Println(encryptedText)

	fmt.Println("=== start decrypting ===")
	encryptedData, err = base64.StdEncoding.DecodeString(encryptedText)
	checkError(err)
	decryptedData, err := RsaDecrypt([]byte(encryptedData), []byte(priKey.String()))
	fmt.Println(string(decryptedData))
}*/

func CreateKeys(publicKeyWriter, privateKeyWriter io.Writer, keyLength int) error {
	// generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(privateKeyWriter, block)
	if err != nil {
		return err
	}

	// generate public key
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	err = pem.Encode(publicKeyWriter, block)
	if err != nil {
		return err
	}

	return nil
}

func RsaDecrypt(ciphertext []byte, privKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func RsaEncrypt(ciphertext []byte, pubKey []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("public key error!")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), ciphertext)
}
