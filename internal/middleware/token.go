package middleware

import (
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type Tokenmanager struct {
	SecretKey string
}

type TokenManager interface {
	TokenGenerate(email, password string) (string, error)
	VerifyToken(tokenString string) bool
	RefreshToken(email, password string) string
	AuthMiddleware(next http.Handler) http.Handler
}

func NewTokenManager(secretKey string) *Tokenmanager {
	return &Tokenmanager{SecretKey: secretKey}
}

func (s *Tokenmanager) TokenGenerate(email, password string) (string, error) {
	tokenAuth := jwtauth.New("HS256", []byte(s.SecretKey), nil)
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

func (s *Tokenmanager) VerifyToken(tokenString string) bool {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Возвращаем секретный ключ для проверки подписи
		return []byte(s.SecretKey), nil
	})
	if err != nil {
		return false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем срок действия токена
		exp := int64(claims["Exp"].(float64))
		if time.Now().Unix() > exp {
			return false
		}
		return true
	}
	return false
}

func (s *Tokenmanager) RefreshToken(email, password string) string {
	tokenAuth := jwtauth.New("HS256", []byte(s.SecretKey), nil)
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

func (s *Tokenmanager) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		token := s.VerifyToken(Usertoken)
		log.Println("midleware")
		if !token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
