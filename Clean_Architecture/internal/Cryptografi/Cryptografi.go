package Cryptografi

import (
	"github.com/agnivade/levenshtein"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Метод для шифрования пароля
func HashPassword(user Models.NewUserRequest) (Models.NewUserRequest, error) {
	password := user.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Models.NewUserRequest{}, err
	}
	user.Password = string(hashedPassword)
	return user, nil
}
func CheckPassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("Неверный пароль", hashedPassword, password)
		return err
	}
	return nil
}

// Метод для сравнения запроса с кэшом
func Levanshtain(searchHistory []Models.SearchHistory, qwery string) (Models.RequestAddress, bool) {
	var result Models.RequestAddress
	var found bool
	threshold := 0.7

	for _, v := range searchHistory {
		distance := levenshtein.ComputeDistance(qwery, v.Address)
		maxLen := maxStringLength(qwery, v.Address)
		similarity := 1 - float64(distance)/float64(maxLen)
		if similarity >= threshold {
			result.Addres = v.Address
			result.RequestSearch.Lat = v.Lat
			result.RequestSearch.Lng = v.Lng
			found = true
			break
		}
	}
	return result, found
}

func maxStringLength(s1, s2 string) int {
	return max(len([]rune(s1)), len([]rune(s2)))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
