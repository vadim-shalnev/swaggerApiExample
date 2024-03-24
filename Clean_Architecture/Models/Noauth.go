package Models

type TokenString struct {
	Token string `json:"auth"`
}
type NewUserResponse struct {
	Email string      `json:"email"`
	Role  string      `json:"role"`
	Token TokenString `json:"token"`
}
type NewUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
