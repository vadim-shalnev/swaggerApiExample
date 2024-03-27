package authService

import (
	"context"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Cryptografi"
	repository "github.com/vadim-shalnev/swaggerApiExample/internal/Repository"
	"log"
	"strings"
	"time"
)

func NewAuthService(repository repository.Repository) *Authservice {
	return &Authservice{repo: repository}
}

func (a *Authservice) Register(ctx context.Context, regData mod.NewUserRequest) (mod.NewUserResponse, error) {
	tokenAuth, err := a.TokenGenerate(ctx, regData.Email, regData.Password)
	if err != nil {
		return mod.NewUserResponse{}, err
	}
	var userResponse mod.NewUserResponse
	userResponse.Email = regData.Email
	userResponse.Token.Token = tokenAuth
	userResponse.Role = regData.Role

	// Хэшируем пароль и добавляем его в запрос к БД
	hashPass, err := Cryptografi.HashPassword(regData)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to hash password")
	}
	err = a.repo.CreateUser(ctx, hashPass)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to add new userController to the database")
	}

	return userResponse, nil
}

func (a *Authservice) Login(ctx context.Context, loginData mod.NewUserRequest) (mod.NewUserResponse, error) {
	var userResponse mod.NewUserResponse
	userResponse.Email = loginData.Email
	userResponse.Role = loginData.Role
	userToken := ctx.Value("jwt_token").(string)
	log.Println("ligintoken", userToken)
	emailValid, passwordValid, tokenValid := a.UserInfoChecker(ctx, loginData.Email, loginData.Password, userToken)
	if !emailValid {
		return mod.NewUserResponse{}, errors.New("invalid email")
	}
	if !passwordValid {
		return mod.NewUserResponse{}, errors.New("invalid password")
	}
	if !tokenValid {
		freshToken := a.RefreshToken(ctx, loginData.Email, loginData.Password)
		log.Println("freshtoken", freshToken)
		userResponse.Token.Token = freshToken
		return mod.NewUserResponse{}, errors.New("you have successfully logged out of the service")
	}
	return userResponse, nil
}

func (a *Authservice) UserInfoChecker(ctx context.Context, email, password, token string) (bool, bool, bool) {
	email, _, tokenValid := a.VerifyToken(token)
	log.Println("UserInfoChecker email", email)
	user, _ := a.repo.GetByEmail(ctx, email)
	if !tokenValid {
		return false, false, false
	}
	if strings.TrimSpace(user.Email) != strings.TrimSpace(email) {
		log.Println("!"+user.Email+"!", "!"+email+"!", "false")
		return false, false, false
	} else {
		log.Println("!"+user.Email+"!", "!"+email+"!", "true")
	}

	if err := Cryptografi.CheckPassword(user.Password, password); err != nil {
		log.Println("check pass is false")
		return false, false, false
	} else {
		log.Println("check is ok", password, user.Password)
	}

	return true, true, tokenValid
}

func (a *Authservice) TokenGenerate(ctx context.Context, email, password string) (string, error) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"Email":    email,
		"Password": password,
		"Exp":      time.Now().Add(time.Second * 60).Unix(),
	})
	if err != nil {
		return "", errors.New("token generation error")
	}
	return tokenString, nil
}

func (a *Authservice) VerifyToken(tokenString string) (string, string, bool) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Возвращаем секретный ключ для проверки подписи
		return []byte("secret"), nil
	})
	if err != nil {
		return "", "", false
	}
	if !token.Valid {
		return "", "", false
	}

	var username string
	var password string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем срок действия токена
		exp := int64(claims["Exp"].(float64))
		if time.Now().Unix() > exp {
			return "", "", false
		}

		username = claims["Email"].(string)
		password = claims["Password"].(string)

		return username, password, true
	}
	return "", "", false
}

func (a *Authservice) RefreshToken(ctx context.Context, email, password string) string {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"Email":    email,
		"Password": password,
		"Exp":      time.Now().Add(time.Second * 60).Unix(),
	})
	if err != nil {
		log.Println(err)
	}

	return tokenString
}
