package Cryptografi

import "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"

// Метод для шифрования пароля
func (c *UserServiceImpl) HashPassword(user Models.NewUserRequest) (Models.NewUserRequest, error) {
	return func(user Models.NewUserRequest) (Models.NewUserRequest, error) {
		user.Password = c.hashPassword(user.Password)
		return user, nil
	}
}

//Метод для сравнения запроса с имеющимся в базе
