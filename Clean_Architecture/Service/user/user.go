package user

import (
	"context"
	"errors"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Repository"
	"strconv"
)

func NewAuthService(repository repository.Repository) *UserServiceImpl {
	return &UserServiceImpl{repo: repository}
}

func (u *UserServiceImpl) GetUser(ctx context.Context, id string) (mod.NewUserResponse, error) {
	var userResponse mod.NewUserResponse
	userID, _ := strconv.Atoi(id)
	user, err := u.repo.GetUser(ctx, userID)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("пользователь с таким id не найден")
	}
	userToken := ctx.Value("jwt_token").(string)
	userResponse.Email = user.Email
	userResponse.Role = user.Role
	userResponse.Token.Token = userToken

	return userResponse, nil
}
