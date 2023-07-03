package service

import (
	"context"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"

	"url-shortener/pkg/app/urlshortener/entity"
	"url-shortener/pkg/app/urlshortener/repository"
	"url-shortener/pkg/bloom"
	"url-shortener/pkg/errors"
	"url-shortener/pkg/utils"
)

var _ ShortenedURLService = &shortenedURLServiceImpl{}

// ShortURLOption ShortURL option
type ShortURLOption struct {
	// ExpiredAt is optional expire time at
	ExpiredAt *time.Time
}

// ShortenedURLService define shortened URL service
type ShortenedURLService interface {
	// ShortURL create short key and store to datastore
	ShortURL(ctx context.Context, url string, opts *ShortURLOption) (shortenedURL *entity.ShortenedURL, err error)

	// RetrieveShortenedURL retrieve short url
	RetrieveShortenedURL(ctx context.Context, short string) (shortenedURL *entity.ShortenedURL, err error)
}

// shortenedURLServiceImpl implement ShortenedURLService
type shortenedURLServiceImpl struct {
	repo        repository.Repository
	bloomFilter bloom.Filter
}

func New(repo repository.Repository, bloomFilter bloom.Filter) ShortenedURLService {
	return &shortenedURLServiceImpl{
		repo:        repo,
		bloomFilter: bloomFilter,
	}
}

// ShortURL is implement ShortenedURLService method
func (srv *shortenedURLServiceImpl) ShortURL(ctx context.Context, originalURL string, opts *ShortURLOption) (shortenedURL *entity.ShortenedURL, err error) {
	var (
		expiredAt time.Time

		short string
	)

	logger := log.Ctx(ctx)

	if originalURL == "" {
		return nil, errors.Wrap(errors.ErrInvalidInput, "input originalURL is empty")
	}

	_, err = url.Parse(originalURL)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrInvalidInput, "input originalURL is %v err = %v", originalURL, err)
	}

	now := time.Now().UTC()
	if opts.ExpiredAt == nil {
		expiredAt = now.Add(entity.DefaultShortenedURLExpireDur)
	} else {
		expiredAt = *opts.ExpiredAt
	}

	// using bloom filter prevent direct to hit database
	// accept missing rate then reduce direct access datastore
	for true {
		short = srv.NewShortID(ctx)
		if !srv.IsInBloomFilter(ctx, short) {
			break
		}

		logger.Trace().Msgf("short id =%v is in bloom, need create new one", short)
	}

	shortenedURL = entity.NewShortenedURL(
		short,
		originalURL,
		expiredAt,
	)

	if err := srv.repo.StoreShortenedURL(ctx, shortenedURL); err != nil {
		return nil, err
	}

	// Add to bloom filter
	srv.AddToBloomFilter(ctx, short)

	return
}

// RetrieveShortenedURL is implement ShortenedURLService method
func (srv *shortenedURLServiceImpl) RetrieveShortenedURL(ctx context.Context, short string) (shortenedURL *entity.ShortenedURL, err error) {
	now := time.Now().UTC()

	if short == "" {
		return nil, errors.Wrap(errors.ErrInvalidInput, "input shot url is empty")
	}

	if !srv.IsInBloomFilter(ctx, short) {
		return nil, errors.Wrap(errors.ErrPageNotFound, "shortened url not in bloom filter")
	}

	shortenedURL, err = srv.repo.FindShortenedURL(ctx, short)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrPageNotFound, "shortened =%v url not found ", short)
	}

	if shortenedURL.ExpiredAt.UnixMilli() <= now.UnixMilli() {
		return nil, errors.Wrapf(
			errors.ErrShortenedURLExpire,
			"the short = %v is expire, now=%v expireAt=%v",
			short,
			now.Format(time.RFC3339),
			shortenedURL.ExpiredAt.UTC().Format(time.RFC3339),
		)
	}

	return
}

// AddToBloomFilter wrap bloomFilter add method
func (srv *shortenedURLServiceImpl) AddToBloomFilter(ctx context.Context, item string) {
	srv.bloomFilter.Add(ctx, item)
}

// IsInBloomFilter wrap bloomFilter exist method
func (srv *shortenedURLServiceImpl) IsInBloomFilter(ctx context.Context, item string) bool {
	return srv.bloomFilter.Exist(ctx, item)
}

// NewShortID new short id method

func (srv *shortenedURLServiceImpl) NewShortID(ctx context.Context) string {
	length := 8
	// TODO: using key generator service or snowflake
	// In high concurrency , the token will create too many token in same time
	// consider using key generator service prepare some tokens, then can quick response token.
	return utils.RandStringBytesMaskImprSrc(length)
}
