package geocodService

import (
	"context"
	"github.com/ekomobile/dadata/v2/api/model"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Repository"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Service/authService"
)

type GeocodeWorkerImpl struct {
	repo repository.Repository
	auth authService.AuthService
}

type GeocodeWorker interface {
	Search(ctx context.Context, userRequest mod.RequestQuery) (interface{}, error)
	Address(ctx context.Context, userRequest mod.RequestQuery) (interface{}, error)
	HandleWorker(ctx context.Context, query mod.RequestQuery) (mod.RequestAddress, error)
	CacheChecker(ctx context.Context, query mod.RequestQuery, ttl int) (bool, mod.RequestAddress, string, error)
	Geocode(query mod.RequestQuery) ([]*model.Address, error)
}
