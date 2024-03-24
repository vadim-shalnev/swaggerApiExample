package geocoder

import (
	"context"
	"github.com/ekomobile/dadata/v2/api/model"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Repository"
	"io"
)

type DadataWorkerImpl struct {
	repo repository.Repository
}

type DadataWorker interface {
	Register(ctx context.Context, userBody io.ReadCloser) (mod.NewUserResponse, error)
	Login(ctx context.Context, userBody io.ReadCloser) (mod.NewUserResponse, error)
	Search(ctx context.Context, userRequest io.ReadCloser) (interface{}, error)
	Address(ctx context.Context, userRequest io.ReadCloser) (interface{}, error)
	UserInfoChecker(ctx context.Context, email, password, token string) (bool, bool, bool)
	HandleWorker(ctx context.Context, query mod.RequestQuery) (mod.RequestAddress, error)
	RefreshToken(ctx context.Context, email, password string) string
	Geocode(query mod.RequestQuery) ([]*model.Address, error)
	GetUser(ctx context.Context, id string) (mod.NewUserResponse, error)
}
