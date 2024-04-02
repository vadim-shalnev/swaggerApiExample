package geocodService

import (
	"context"
	"github.com/ekomobile/dadata/v2/api/model"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/config"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geo/geocodeRepository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/middleware"
)

type Geocodeservice struct {
	repo         geocodeRepository.GeocodeRepository
	Tokenmanager middleware.TokenManager
	Config       config.AppConf
}

type GeocodeService interface {
	Search(ctx context.Context, userRequest mod.RequestQuery) (mod.RequestQuery, error)
	Address(ctx context.Context, userRequest mod.RequestQuery) (mod.RequestQuery, error)
	HandleWorker(ctx context.Context, query mod.RequestQuery) (mod.RequestAddress, error)
	Geocode(query mod.RequestQuery) ([]*model.Address, error)
}
