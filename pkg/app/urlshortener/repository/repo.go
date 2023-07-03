package repository

import (
	"context"
	"encoding/json"

	"github.com/coocood/freecache"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"url-shortener/pkg/app/urlshortener/entity"
	"url-shortener/pkg/db"
	"url-shortener/pkg/errors"
)

// Repository define url shortener repository layer
type Repository interface {
	// StoreShortenedURL store shortened URL
	StoreShortenedURL(
		ctx context.Context,
		shortenedURL *entity.ShortenedURL,
	) (err error)

	// FindShortenedURL find shortened URL by short id
	FindShortenedURL(
		ctx context.Context,
		short string,
	) (shortenedURL *entity.ShortenedURL, err error)
}

// RepoImpl is implementation for Repository
type RepoImpl struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
	cache   *freecache.Cache
}

// New Repository constructor
func New(conn db.Connection) Repository {
	// In bytes, where 1024 * 1024 represents a single Megabyte, and 100 * 1024*1024 represents 100 Megabytes.
	cacheSize := 100 * 1024 * 1024

	return &RepoImpl{
		readDB:  conn.ReadDB(),
		writeDB: conn.WriteDB(),
		cache:   freecache.NewCache(cacheSize),
	}
}

// StoreShortenedURL method is implementation for Repository
func (repo *RepoImpl) StoreShortenedURL(ctx context.Context, shortenedURL *entity.ShortenedURL) (err error) {
	const (
		sql = `INSERT INTO "shortened_urls" ("short","original_url","created_at","expired_at") VALUES (?,?,?,?)`
	)

	err = repo.writeDB.WithContext(ctx).Exec(
		sql,
		shortenedURL.Short,
		shortenedURL.OriginalURL,
		shortenedURL.CreatedAt.UnixMilli(),
		shortenedURL.ExpiredAt.UnixMilli(),
	).Error
	if err != nil {
		return errors.Wrapf(
			errors.ErrInternal,
			"failed to store shortenedURL %v", shortenedURL,
		)
	}

	return nil
}

// FindShortenedURL method is implementation for Repository
func (repo *RepoImpl) FindShortenedURL(ctx context.Context, short string) (shortenedURL *entity.ShortenedURL, err error) {
	keyPrefix := "ShortenedURL:"
	shortenedURL = &entity.ShortenedURL{}
	expireSeconds := 600

	logger := log.Ctx(ctx)

	// try to get data from local cache
	data, err := repo.cache.Get([]byte(keyPrefix + short))
	if err != nil {
		// not found in cache
		// load data from datastore
		if errors.Is(err, freecache.ErrNotFound) {
			const (
				sql = `SELECT 
    				short,
    				original_url,
    				to_timestamp(created_at/1000) as created_at,
    				to_timestamp(expired_at/1000) as expired_at
       		   FROM shortened_urls 
       		   WHERE short = ? LIMIT 1`
			)

			rtn := make([]*entity.ShortenedURL, 0)
			err = repo.readDB.
				WithContext(ctx).
				Raw(sql, short).
				Scan(&rtn).
				Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.Wrapf(errors.ErrResourceNotFound, "short = %v not found", short)
				}

				return nil, errors.Wrapf(errors.ErrInternal, "short = %v not found, err = %v", short, err)
			}

			if len(rtn) == 0 {
				return nil, errors.Wrapf(errors.ErrResourceNotFound, "short = %v not found", short)
			}

			// write entity into cache
			data, _ = json.Marshal(rtn[0])
			err = repo.cache.Set([]byte(keyPrefix+short), data, expireSeconds)
			if err != nil {
				logger.Warn().Err(err).Msgf("fail to write short=%v into cache", short)
			}

			return rtn[0], nil
		}

		return nil, errors.Wrapf(errors.ErrInternal, "failed to get data from cache %v", err)
	}

	_ = json.Unmarshal(data, shortenedURL)

	return shortenedURL, nil
}
