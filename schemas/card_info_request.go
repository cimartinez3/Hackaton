package schemas

type CardInfoRequest struct {
	Email string `json:"email"`
}

type CardInfoResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Number string `json:"number"`
}
