package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"url-shortener/pkg/app/urlshortener/entity"
	"url-shortener/pkg/app/urlshortener/mocks"
	"url-shortener/pkg/app/urlshortener/repository"
	. "url-shortener/pkg/app/urlshortener/service"
	"url-shortener/pkg/bloom"
	bm "url-shortener/pkg/bloom/mocks"
	"url-shortener/pkg/errors"
	"url-shortener/pkg/logging"
)

func Test_shortenedURLServiceImpl_ShortURL(t *testing.T) {
	type args struct {
		url  string
		opts *ShortURLOption
	}
	tests := []struct {
		name        string
		repo        repository.Repository
		bloomFilter bloom.Filter
		args        args
		expected    *entity.ShortenedURL
		err         error
	}{
		{
			name: "Success",
			repo: func() repository.Repository {
				repo := mocks.NewRepository(t)
				repo.EXPECT().StoreShortenedURL(mock.Anything, mock.Anything).Return(nil)
				return repo
			}(),
			bloomFilter: func() bloom.Filter {
				bf := bm.NewFilter(t)
				bf.EXPECT().Exist(mock.Anything, mock.Anything).Return(false)
				bf.EXPECT().Add(mock.Anything, mock.Anything)
				return bf
			}(),
			args: args{
				url: "https://www.dcard.tw/f",
				opts: &ShortURLOption{
					ExpiredAt: func() *time.Time {
						t := time.Date(2023, time.June, 30, 11, 00, 00, 000, time.UTC)
						return &t
					}(),
				},
			},
			expected: &entity.ShortenedURL{
				OriginalURL: "https://www.dcard.tw/f",
				ExpiredAt:   time.Date(2023, time.June, 30, 11, 00, 00, 000, time.UTC),
			},
			err: nil,
		},
		{
			name: "SuccessWithoutExpireTime",
			repo: func() repository.Repository {
				repo := mocks.NewRepository(t)
				repo.EXPECT().StoreShortenedURL(mock.Anything, mock.Anything).Return(nil)
				return repo
			}(),
			bloomFilter: func() bloom.Filter {
				bf := bm.NewFilter(t)
				bf.EXPECT().Exist(mock.Anything, mock.Anything).Return(false)
				bf.EXPECT().Add(mock.Anything, mock.Anything)
				return bf
			}(),
			args: args{
				url: "https://www.dcard.tw/f",
				opts: &ShortURLOption{
					ExpiredAt: nil,
				},
			},
			expected: &entity.ShortenedURL{
				OriginalURL: "https://www.dcard.tw/f",
			},
			err: nil,
		},
		{
			name: "SuccessWithShortIDInBloom",
			repo: func() repository.Repository {
				repo := mocks.NewRepository(t)
				repo.EXPECT().StoreShortenedURL(mock.Anything, mock.Anything).Return(nil)
				return repo
			}(),
			bloomFilter: func() bloom.Filter {
				bf := bm.NewFilter(t)
				bf.EXPECT().Exist(mock.Anything, mock.Anything).Return(true).Twice()
				bf.EXPECT().Exist(mock.Anything, mock.Anything).Return(false)
				bf.EXPECT().Add(mock.Anything, mock.Anything)
				return bf
			}(),
			args: args{
				url: "https://www.dcard.tw/f",
				opts: &ShortURLOption{
					ExpiredAt: func() *time.Time {
						t := time.Date(2023, time.June, 30, 11, 00, 00, 000, time.UTC)
						return &t
					}(),
				},
			},
			expected: &entity.ShortenedURL{
				OriginalURL: "https://www.dcard.tw/f",
				ExpiredAt:   time.Date(2023, time.June, 30, 11, 00, 00, 000, time.UTC),
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.SetupWithOption(
				logging.WithDebug(true),
				logging.WithLevel(logging.TraceLevel),
			)
			ctx := context.Background()
			ctx = logger.WithContext(ctx)

			srv := New(
				tt.repo,
				tt.bloomFilter,
			)

			actual, err := srv.ShortURL(ctx, tt.args.url, tt.args.opts)
			if (err != nil) || tt.err != nil {
				assert.Truef(
					t,
					errors.Is(err, tt.err),
					"ShortURL() error = %v, expected error %v",
					err,
					tt.err,
				)
				return
			}

			t.Logf("actual.Short = %v", actual.Short)
			assert.NotEmpty(t, actual.Short)
			assert.Equalf(
				t,
				tt.expected.OriginalURL,
				actual.OriginalURL,
				"ShortURL() actual OriginalURL = %v, expected %v",
				actual.OriginalURL,
				tt.expected.OriginalURL,
			)
			if tt.args.opts.ExpiredAt != nil {
				assert.Equalf(
					t,
					tt.expected.ExpiredAt,
					actual.ExpiredAt,
					"ShortURL() actual ExpiredAt = %v, expected %v",
					actual.ExpiredAt,
					tt.expected.ExpiredAt,
				)
			}
		})
	}
}

func Test_shortenedURLServiceImpl_RetrieveShortenedURL(t *testing.T) {
	type args struct {
		short string
	}
	tests := []struct {
		name        string
		repo        repository.Repository
		bloomFilter bloom.Filter
		args        args
		expected    *entity.ShortenedURL
		err         error
	}{
		{
			name: "Success",
			repo: func() repository.Repository {
				repo := mocks.NewRepository(t)
				repo.EXPECT().
					FindShortenedURL(mock.Anything, mock.Anything).
					Return(&entity.ShortenedURL{
						Short:       "6Xme5Xwp",
						OriginalURL: "https://www.dcard.tw/f",
						CreatedAt:   time.Now(),
						ExpiredAt:   time.Now().Add(entity.DefaultShortenedURLExpireDur),
					}, nil)
				return repo
			}(),
			bloomFilter: func() bloom.Filter {
				bf := bm.NewFilter(t)
				bf.EXPECT().Exist(mock.Anything, mock.Anything).Return(true)
				return bf
			}(),
			args: args{
				short: "6Xme5Xwp",
			},
			expected: &entity.ShortenedURL{
				Short:       "6Xme5Xwp",
				OriginalURL: "https://www.dcard.tw/f",
			},
			err: nil,
		},
		{
			name: "NotInBloom",
			repo: func() repository.Repository {
				repo := mocks.NewRepository(t)
				return repo
			}(),
			bloomFilter: func() bloom.Filter {
				bf := bm.NewFilter(t)
				bf.EXPECT().Exist(mock.Anything, mock.Anything).Return(false)
				return bf
			}(),
			args: args{
				short: "6Xme5Xwp",
			},
			expected: nil,
			err:      errors.ErrPageNotFound,
		},
		{
			name: "BloomHasError",
			repo: func() repository.Repository {
				repo := mocks.NewRepository(t)
				repo.EXPECT().
					FindShortenedURL(mock.Anything, mock.Anything).
					Return(nil, errors.ErrResourceNotFound)
				return repo
			}(),
			bloomFilter: func() bloom.Filter {
				bf := bm.NewFilter(t)
				bf.EXPECT().Exist(mock.Anything, mock.Anything).Return(true)
				return bf
			}(),
			args: args{
				short: "6Xme5Xwp",
			},
			expected: nil,
			err:      errors.ErrPageNotFound,
		},
		{
			name: "Expire",
			repo: func() repository.Repository {
				repo := mocks.NewRepository(t)
				repo.EXPECT().
					FindShortenedURL(mock.Anything, mock.Anything).
					Return(&entity.ShortenedURL{
						Short:       "6Xme5Xwp",
						OriginalURL: "https://www.dcard.tw/f",
						CreatedAt:   time.Now(),
						ExpiredAt:   time.Now().Add(-entity.DefaultShortenedURLExpireDur),
					}, nil)
				return repo
			}(),
			bloomFilter: func() bloom.Filter {
				bf := bm.NewFilter(t)
				bf.EXPECT().Exist(mock.Anything, mock.Anything).Return(true)
				return bf
			}(),
			args: args{
				short: "6Xme5Xwp",
			},
			expected: nil,
			err:      errors.ErrShortenedURLExpire,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.SetupWithOption(
				logging.WithDebug(true),
				logging.WithLevel(logging.TraceLevel),
			)
			ctx := context.Background()
			ctx = logger.WithContext(ctx)

			srv := New(
				tt.repo,
				tt.bloomFilter,
			)

			actual, err := srv.RetrieveShortenedURL(ctx, tt.args.short)
			if (err != nil) || tt.err != nil {
				assert.Truef(
					t,
					errors.Is(err, tt.err),
					"RetrieveShortenedURL(ctx, %v), error = %v, expected error %v",
					tt.args.short,
					err,
					tt.err,
				)
				return
			}

			assert.Equalf(
				t,
				tt.expected.OriginalURL,
				actual.OriginalURL,
				"RetrieveShortenedURL(%v, %v)",
				ctx,
				tt.args.short,
			)
		})
	}
}
