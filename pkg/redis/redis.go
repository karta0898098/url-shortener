package redis

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Config redis config
type Config struct {
	Address      string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	MaxRetries   int           `mapstructure:"max_retries"`
	PoolSize     int           `mapstructure:"pool_size"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// ClusterConfig redis cluster config
type ClusterConfig struct {
	Address      []string      `mapstructure:"addrs"`
	Password     string        `mapstructure:"password"`
	MaxRetries   int           `mapstructure:"max_retries"`
	PoolSize     int           `mapstructure:"pool_size"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// NewRedis new redis client
func NewRedis(config Config) (*redis.Client, error) {
	return setupRedis(config)
}

// NewRedisCluster new redis cluster
func NewRedisCluster(config ClusterConfig) (*redis.ClusterClient, error) {
	return setupRedisCluster(config)
}

func setupRedis(config Config) (*redis.Client, error) {
	var (
		ctx context.Context
		rdb *redis.Client
	)
	ctx = context.Background()

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	err := backoff.Retry(func() error {
		rdb = redis.NewClient(&redis.Options{
			Addr:         config.Address,
			Password:     config.Password,
			MaxRetries:   config.MaxRetries,
			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			PoolSize:     config.PoolSize,
		})
		err := rdb.Ping(ctx).Err()
		if err != nil {
			return err
		}
		return nil
	}, bo)

	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Ping to redis %s success", config.Address)

	return rdb, nil
}

func setupRedisCluster(config ClusterConfig) (*redis.ClusterClient, error) {
	var (
		ctx context.Context
		rdb *redis.ClusterClient
	)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	ctx = context.Background()

	err := backoff.Retry(func() error {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			NewClient: func(opt *redis.Options) *redis.Client {
				return redis.NewClient(opt)
			},
			Addrs:        config.Address,
			MaxRedirects: config.MaxRetries,
			Password:     config.Password,
			MaxRetries:   config.MaxRetries,
			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			PoolSize:     config.PoolSize,
		})

		err := rdb.Ping(ctx).Err()
		if err != nil {
			return err
		}

		return nil
	}, bo)

	if err != nil {
		return nil, err
	}

	log.Info().
		Strs("addrs", config.Address).
		Msg("Ping to redis cluster success")

	return rdb, nil
}
