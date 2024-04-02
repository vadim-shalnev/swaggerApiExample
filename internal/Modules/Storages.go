package Modules

import (
	"database/sql"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authRepository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geo/geocodeRepository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userRepository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Cache"
)

type Storages struct {
	Auth authRepository.AuthRepository
	User userRepository.UserRepository
	Geo  geocodeRepository.GeocodeRepository
}

func NewStorages(db *sql.DB, cache Cache.Cache) *Storages {
	return &Storages{
		Auth: authRepository.NewAuthrepository(db),
		User: userRepository.NewUserRepository(db),
		Geo:  geocodeRepository.NewGeocodeRepository(db, cache),
	}
}
