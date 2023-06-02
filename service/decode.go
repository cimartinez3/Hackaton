package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
)

type IRSA interface {
	DecodeRSA(data string) (string, error)
	EncodeRSA(data []byte) ([]byte, error)
}

type RSA struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewRSAService(private, public string) IRSA {
	priBytes, _ := ioutil.ReadFile(private)
	pubBytes, _ := ioutil.ReadFile(public)
	privateKey := BytesToPrivateKey(priBytes)
	publicKey := BytesToPublicKey(pubBytes)
	return &RSA{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

func (r *RSA) DecodeRSA(data string) (string, error) {
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, r.PrivateKey, []byte(data), nil)
	if err != nil {
		fmt.Print("Error in decode RSA: ", err)
		return "", err
	}

	return string(plaintext), nil
}
func (r *RSA) EncodeRSA(data []byte) ([]byte, error) {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, r.PublicKey, data, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ciphertext, nil
}

func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			fmt.Println("Error in DecryptPENBloc")
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		fmt.Println("Error in ParsePKCS1 Private key")
	}
	return key
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			fmt.Println("Error in decryptPEM BLOCK")
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		fmt.Println("Error in ParsePKIXPublicKey")
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		fmt.Println("IFC")
	}
	return key
}
