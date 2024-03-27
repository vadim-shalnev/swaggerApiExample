package geocodService

import (
	"context"
	"github.com/ekomobile/dadata/v2/api/model"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authService"
	repository "github.com/vadim-shalnev/swaggerApiExample/internal/Repository"
)

type Geocodeworker struct {
	repo repository.Repository
	auth authService.AuthService
}

type GeocodeWorker interface {
	Search(ctx context.Context, userRequest mod.RequestQuery) (mod.RequestQuery, error)
	Address(ctx context.Context, userRequest mod.RequestQuery) (mod.RequestQuery, error)
	HandleWorker(ctx context.Context, query mod.RequestQuery) (mod.RequestAddress, error)
	CacheChecker(ctx context.Context, query mod.RequestQuery, ttl int) (bool, mod.RequestAddress, string, error)
	Geocode(query mod.RequestQuery) ([]*model.Address, error)
}
