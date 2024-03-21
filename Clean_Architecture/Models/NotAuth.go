package Models

type TokenString struct {
	Token string `json:"auth"`
}
type NewUserResponse struct {
	Email       string      `json:"email"`
	Password    string      `json:"password"`
	TokenString TokenString `json:"token"`
}
