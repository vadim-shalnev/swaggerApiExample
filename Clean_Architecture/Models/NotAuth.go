package Models

type TokenString struct {
	Token string `json:"auth"`
}
type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type NewUserResponse struct {
	Email       string      `json:"email"`
	Password    string      `json:"password"`
	TokenString TokenString `json:"token"`
}
