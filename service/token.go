package service

import (
	"errors"
	"kushki/hackaton/gateway"
	"kushki/hackaton/schemas"
)

type TokenService struct {
	Mongo   gateway.IMongo
	Decoder IRSA
}

func NewTokenService() *TokenService {
	return &TokenService{
		Mongo:   gateway.NewMongoService(),
		Decoder: NewRSAService(),
	}
}

func (t *TokenService) ProcessToken(req schemas.TokenRequest) (bool, error) {
	tokenRequest := &schemas.TokenRequestSource{}
	_ = t.Decoder.DecodeRSA([]byte(req.Token), tokenRequest)
	switch tokenRequest.Source {
	case "001":
	case "002":
	default:
		return false, errors.New("invalid source")

	}
	return true, nil
}

func generateVaultken() {

}
