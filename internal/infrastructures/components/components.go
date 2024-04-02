package components

import (
	"github.com/vadim-shalnev/swaggerApiExample/config"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Cache"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Cryptografi"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Responder"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/middleware"
	"go.uber.org/zap"
)

type Components struct {
	Conf         config.AppConf
	TokenManager middleware.TokenManager
	Hash         Cryptografi.Hasher
	Cache        Cache.Cache
	Logger       *zap.Logger
	Responder    Responder.Responder
}

func NewComponents(conf config.AppConf, tokenManager middleware.TokenManager, hash Cryptografi.Hasher, cache Cache.Cache, responder Responder.Responder) *Components {
	return &Components{
		Conf:         conf,
		TokenManager: tokenManager,
		Hash:         hash,
		Cache:        cache,
		Responder:    responder,
	}
}
