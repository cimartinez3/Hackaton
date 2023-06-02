package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"kushki/hackaton/gateway"
	"kushki/hackaton/schemas"
)

type TokenService struct {
	Mongo        gateway.IMongo
	Decoder      IRSA
	DecoderVault IRSA
}

func NewTokenService() *TokenService {
	return &TokenService{
		Mongo:        gateway.NewMongoService(),
		Decoder:      NewRSAService("./keys/id_rsa_ksk", "./keys/id_rsa_ksk.pub"),
		DecoderVault: NewRSAService("./keys/id_rsa_vault_ksk", "./keys/id_rsa_vault_ksk.pub"),
	}
}

func (t *TokenService) ProcessToken(req *schemas.TokenRequest) (*schemas.TokenResponse, error) {
	tokenRequest := &schemas.TokenRequestSource{}
	/*_ = t.Decoder.DecodeRSA(req.Token, tokenRequest)*/
	tokenRequest.Card = schemas.Card{
		Name:        "Anthony Torres",
		Number:      "453473482462342",
		ExpiryMonth: "02",
		ExpiryYear:  "03",
		Cvv:         "123",
	}
	tokenRequest.Client.Email = "anthonybastidas48@gmail.com"
	tokenRequest.Source = "001"
	switch tokenRequest.Source {
	case "001":
		res := t.generateVaultToken(*tokenRequest)
		vault := fmt.Sprintf("%s%s%s%s%s%s%s%s%s", res[0], "ksk|", res[1], "ksk|", res[2], "ksk|", res[3], "ksk|", res[4])
		uEnc := base64.URLEncoding.EncodeToString([]byte(vault))

		save := map[string]string{
			"vault":  uEnc,
			"number": res[0],
			"cvv":    res[1],
			"month":  res[2],
			"year":   res[3],
			"email":  res[4],
			"name":   res[5],
			"id":     uuid.Must(uuid.NewV4(), nil).String(),
		}
		t.Mongo.PutDocument(save)
		response := &schemas.TokenResponse{
			Token: uEnc,
			Type:  "001",
		}
		return response, nil
	case "002":
		cvv, _ := t.DecoderVault.EncodeRSA([]byte(tokenRequest.Card.Cvv))
		t.Mongo.UpdateDocument(tokenRequest.CardId, string(cvv))
	default:
		return nil, errors.New("invalid source")

	}
	return nil, nil
}

func (t *TokenService) generateVaultToken(req schemas.TokenRequestSource) []string {
	numberEncode, _ := t.DecoderVault.EncodeRSA([]byte(req.Card.Number))
	cvvEncode, _ := t.DecoderVault.EncodeRSA([]byte(req.Card.Cvv))
	monthEncode, _ := t.DecoderVault.EncodeRSA([]byte(req.Card.ExpiryMonth))
	yearEncode, _ := t.DecoderVault.EncodeRSA([]byte(req.Card.ExpiryYear))
	return []string{string(numberEncode), string(cvvEncode), string(monthEncode), string(yearEncode), req.Client.Email, req.Card.Name}
}
