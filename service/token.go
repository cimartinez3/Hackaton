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

func (t *TokenService) ProcessToken(req *schemas.TokenRequestSource) (*schemas.TokenResponse, error) {

	req.Client.Email = "anthonybastidas48@gmail.com"
	req.Source = "001"
	switch req.Source {
	case "001":
		res := t.generateVaultToken(*req)
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
		res, _ := t.Mongo.GetItem(req.CardId)
		cvv, _ := t.DecoderVault.EncodeRSA([]byte(req.Card.Cvv))
		t.Mongo.UpdateDocument(req.CardId, string(cvv))
		response := &schemas.TokenResponse{
			Token: res.Vault,
			Type:  "001",
		}
		return response, nil
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
