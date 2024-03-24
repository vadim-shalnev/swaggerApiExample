package auth

import (
	"context"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Repository"
	"log"
	"time"
)

func NewAuthService(repository repository.Repository) *AuthServiceImpl {
	return &AuthServiceImpl{repo: repository}
}

func (a *AuthServiceImpl) Register(ctx context.Context, regData mod.NewUserRequest) (mod.NewUserResponse, error) {
	tokenAuth, err := a.TokenGenerate(ctx, regData.Email, regData.Password)
	if err != nil {
		return mod.NewUserResponse{}, err
	}
	var userResponse mod.NewUserResponse
	userResponse.Email = regData.Email
	userResponse.Token.Token = tokenAuth
	userResponse.Role = regData.Role

	// Хэшируем пароль и добавляем его в запрос к БД
	hashPassword, err := a.crypto.HashPassword(regData.Password)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to hash password")
	}
	regData.Password = hashPassword

	err = a.repo.CreateUser(ctx, regData)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to add new user to the database")
	}

	return userResponse, nil
}

func (a *AuthServiceImpl) Login(ctx context.Context, loginData mod.NewUserRequest) (mod.NewUserResponse, error) {
	var userResponse mod.NewUserResponse
	userResponse.Email = loginData.Email
	userResponse.Role = loginData.Role
	userToken := ctx.Value("jwt_token").(string)
	emailValid, passwordValid, tokenValid := a.UserInfoChecker(ctx, loginData.Email, loginData.Password, userToken)
	if !emailValid {
		return mod.NewUserResponse{}, errors.New("invalid email")
	}
	if !passwordValid {
		return mod.NewUserResponse{}, errors.New("invalid password")
	}
	if !tokenValid {
		freshToken := a.RefreshToken(ctx, loginData.Email, loginData.Password)
		userResponse.Token.Token = freshToken
		return mod.NewUserResponse{}, errors.New("you have successfully logged out of the service")
	}
	return userResponse, nil
}

func (a *AuthServiceImpl) UserInfoChecker(ctx context.Context, email, password, token string) (bool, bool, bool) {
	emailValid := a.repo.CheckEmail(ctx, email)
	passwordValid := a.repo.CheckPassword(ctx, password)
	_, tokenValid := a.VerifyToken(token, "exp")
	return emailValid, passwordValid, tokenValid
}

func (a *AuthServiceImpl) TokenGenerate(ctx context.Context, email, password string) (string, error) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"Username": email,
		"Password": password,
		"Exp":      time.Now().Add(time.Second * 60).Unix(),
	})
	if err != nil {
		return "", errors.New("token generation error")
	}
	return tokenString, nil
}

func (a *AuthServiceImpl) VerifyToken(tokenString, searchIntoken string) (string, bool) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Возвращаем секретный ключ для проверки подписи
		return []byte("secret"), nil
	})
	if err != nil {
		return "", false
	}
	if !token.Valid {
		return "", false
	}

	var search string
	if searchIntoken == "exp" {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := int64(claims[searchIntoken].(float64))
			if time.Now().Unix() > exp {
				return "", false
			}
		}
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			search = claims[searchIntoken].(string)
		}
	}

	return search, true
}

func (a *AuthServiceImpl) RefreshToken(ctx context.Context, email, password string) string {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"Username": email,
		"Password": password,
		"Exp":      time.Now().Add(time.Second * 60).Unix(),
	})
	if err != nil {
		log.Println(err)
	}

	return tokenString
}
