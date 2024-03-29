package geocodeRepository

import (
	"context"
	"database/sql"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
)

type Geocoderepository struct {
	DB *sql.DB
}

type GeocodeRepository interface {
	CacheChecker(ctx context.Context, email string, historyCount int) ([]Models.SearchHistory, error)
	Insert(ctx context.Context, email string, query Models.RequestQuery, requestQuery Models.RequestAddress) error
	GetByEmail(ctx context.Context, email string) (Models.User, error)
}
