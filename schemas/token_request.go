package schemas

type TokenRequest struct {
	Token string `json:"token"`
}

type TokenRequestSource struct {
	Card        Card    `json:"card"`
	TotalAmount float64 `json:"totalAmount"`
	Source      string  `json:"source"`
	Client      Client  `json:"clientInfo"`
	CardId      string  `json:"cardId"`
}

type Card struct {
	Name        string `json:"name"`
	Number      string `json:"number"`
	ExpiryMonth string `json:"expiryMonth"`
	ExpiryYear  string `json:"expiryYear"`
	Cvv         string `json:"totalAmount"`
}

type Client struct {
	Email string `json:"email"`
}
