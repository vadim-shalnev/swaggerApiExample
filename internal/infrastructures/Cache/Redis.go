package Cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
	"go.uber.org/zap"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

type Cache interface {
	Get(ctx context.Context, key string, query *Models.RequestAddress) error
	Set(ctx context.Context, key string, value Models.RequestAddress, expire time.Duration)
}

func NewRedisCache(address, password string, logger *zap.Logger) *RedisCache {
	// В goboilerplate структура Options немного другая, хотя версия тут и там одинаковые
	// добавили от себя в redis?
	//
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		// добавляем порт
		DB: 0, // use default DB
	})
	ctx := context.Background()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Error("Error connection to redis", zap.Error(err))
	}
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string, query *Models.RequestAddress) error {
	b, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return err
		}
		// loger err
	}
	buffer := make([]string, 3)
	json.Unmarshal(b, &buffer)
	query.Addres = buffer[0]
	query.RequestSearch.Lat = buffer[1]
	query.RequestSearch.Lng = buffer[2]
	return nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value Models.RequestAddress, expire time.Duration) {
	var buffer []string
	buffer = append(buffer, value.Addres, value.RequestSearch.Lat, value.RequestSearch.Lng)
	b, _ := json.Marshal(buffer)
	r.client.Set(ctx, key, b, expire)

}
