package geocodeRepository

import (
	"context"
	"database/sql"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Cache"
)

type Geocoderepository struct {
	DB    *sql.DB
	Cache Cache.Cache
}

type GeocodeRepository interface {
	CacheChecker(ctx context.Context, query Models.RequestQuery) (Models.RequestAddress, error)
	Insert(ctx context.Context, email string, query Models.RequestQuery, requestQuery Models.RequestAddress) error
	GetByEmail(ctx context.Context, email string) (Models.User, error)
}
