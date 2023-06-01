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
		DecoderVault: NewRSAService("./keys/id_rsa_vault", "./keys/id_rsa_vault.pub"),
	}
}

func (t *TokenService) ProcessToken(req schemas.TokenRequest) (string, error) {
	tokenRequest := &schemas.TokenRequestSource{}
	_ = t.Decoder.DecodeRSA([]byte(req.Token), tokenRequest)
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
			"id":     uuid.Must(uuid.NewV4(), nil).String(),
		}
		t.Mongo.PutDocument(save)
		return uEnc, nil
	case "002":
	default:
		return "false", errors.New("invalid source")

	}
	return "", nil
}

func (t *TokenService) generateVaultToken(req schemas.TokenRequestSource) []string {
	numberEncode, _ := t.DecoderVault.EncodeRSA([]byte(req.Card.Number))
	cvvEncode, _ := t.DecoderVault.EncodeRSA([]byte(req.Card.Cvv))
	monthEncode, _ := t.DecoderVault.EncodeRSA([]byte(req.Card.ExpiryMonth))
	yearEncode, _ := t.DecoderVault.EncodeRSA([]byte(req.Card.ExpiryYear))
	return []string{string(numberEncode), string(cvvEncode), string(monthEncode), string(yearEncode), req.Client.Email}
}
