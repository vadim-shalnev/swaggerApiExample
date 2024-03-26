package userService

import (
	"context"
	"errors"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/internal/Repository"
	"strconv"
)

func NewAuthService(repository repository.Repository) *Userservice {
	return &Userservice{repo: repository}
}

func (u *Userservice) GetUser(ctx context.Context, id string) (mod.NewUserResponse, error) {
	var userResponse mod.NewUserResponse
	userID, _ := strconv.Atoi(id)
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("пользователь с таким id не найден")
	}
	userToken := ctx.Value("jwt_token").(string)
	userResponse.Email = user.Email
	userResponse.Role = user.Role
	userResponse.Token.Token = userToken

	return userResponse, nil
}

func (u *Userservice) DelUser(ctx context.Context, id string) error {
	userID, _ := strconv.Atoi(id)
	err := u.repo.Delete(ctx, userID)
	if err != nil {
		return errors.New("пользователь с таким id не найден")
	}
	return nil
}

func (u *Userservice) ListUsers(ctx context.Context) ([]mod.User, error) {
	// добавить логику для роли админа других ролей
	users, err := u.repo.List(ctx)
	if err != nil {
		return nil, errors.New("пользователи не найдены")
	}
	return users, nil
}
