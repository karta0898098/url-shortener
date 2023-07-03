package bloom

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Namespace       string
	ErrorRate       float64
	DefaultCapacity int64
}

type RedisImpl struct {
	rds       *redis.Client
	namespace string
}

func NewRedisFilter(namespace string, rds *redis.Client) Filter {
	cfg := &Config{
		Namespace:       namespace,
		ErrorRate:       0.001,
		DefaultCapacity: 100000,
	}

	ctx := context.Background()

	result, err := rds.Do(ctx, "BF.INFO", cfg.Namespace).Result()
	if err != nil {

		if err.Error() == "ERR not found" {
			log.Info().
				Msgf("create new redis bloom filter")
			err := rds.Do(ctx, "BF.RESERVE", cfg.Namespace, cfg.ErrorRate, cfg.DefaultCapacity).Err()
			if err != nil {
				log.Panic().
					Err(err).
					Msgf("failed init redis bloom filter")
			}
		} else {
			log.Panic().
				Err(err).
				Msgf("failed init redis bloom filter")
		}
	} else {
		log.Info().
			Interface("info", result).
			Msgf("check redis bloom info")
	}

	return &RedisImpl{
		rds: rds,
	}
}

func (r *RedisImpl) Add(ctx context.Context, item interface{}) {
	result, err := r.rds.Do(ctx, "BF.ADD", r.namespace, item).Int()
	if err != nil || result != 1 {
		log.Error().
			Err(err).
			Msgf("failed to add %v into %v bloom filter", item, r.namespace)
	}
}

func (r *RedisImpl) Exist(ctx context.Context, item interface{}) bool {
	result, err := r.rds.Do(ctx, "BF.EXISTS", r.namespace, item).Int()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to check %v in %v bloom filter", item, r.namespace)

		return false
	}

	return result == 1
}

func (r *RedisImpl) GetFilterNamespace() string {
	return r.namespace
}
