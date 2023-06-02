package schemas

type TokenResponse struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}
