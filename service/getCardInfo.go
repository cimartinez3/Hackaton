package service

import (
	"fmt"
	"kushki/hackaton/gateway"
	"kushki/hackaton/schemas"
	"log"
)

func GetCardsInfo(email string) []schemas.CardInfoResponse {
	log.Println("Get Card Info Service")

	response := make([]schemas.CardInfoResponse, 0)

	decoder := NewRSAService("./keys/id_rsa_vault_ksk", "./keys/id_rsa_vault_ksk.pub")

	mg := gateway.NewMongoService()

	items, err := mg.GetDocument(email)

	if err != nil {
		log.Println(err)
		return nil
	}

	for _, item := range items {
		card := schemas.CardInfoResponse{
			ID:   item.ID,
			Name: item.Name,
		}
		cardDecoded, err := decoder.DecodeRSA(item.Number)
		if err != nil {
			log.Println("ERROR DECODING", err)
			continue
		}

		if cardDecoded != "" && len(cardDecoded) > 8 {
			card.Number = cardDecoded[:4] + "XXXX" + cardDecoded[8:]
		}

		response = append(response, card)
	}

	fmt.Println(response)

	return response
}
