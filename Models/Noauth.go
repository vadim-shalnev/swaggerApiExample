package Models

import "time"

type TokenString struct {
	Token string `json:"authController"`
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
type User struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
