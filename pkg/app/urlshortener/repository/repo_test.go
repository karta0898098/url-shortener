package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"url-shortener/pkg/app/urlshortener/entity"
	"url-shortener/pkg/app/urlshortener/repository"
	"url-shortener/pkg/db"
	"url-shortener/pkg/db/conn"
	"url-shortener/pkg/errors"
	"url-shortener/pkg/logging"
)

func TestRepoImpl_StoreShortenedURL(t *testing.T) {
	connection, err := db.NewConnection(db.Config{
		Conn: conn.Database{
			Debug:    true,
			Name:     "url_shortener",
			Host:     "localhost",
			User:     "url_shortener_service",
			Port:     5432,
			Password: "pqV7EJ8bYJpFDXXJtw66s6JKG4xpZb4v",
			Type:     conn.Postgres,
		},
	})
	assert.NoError(t, err)

	type args struct {
		shortenedURL *entity.ShortenedURL
	}
	tests := []struct {
		name     string
		args     args
		setup    func()
		teardown func()
		err      error
	}{
		{
			name: "Success",
			args: args{
				shortenedURL: &entity.ShortenedURL{
					Short:       "LsE2ypFI",
					OriginalURL: "https://www.dcard.tw/",
					CreatedAt:   time.Now(),
					ExpiredAt:   time.Now().Add(time.Hour),
				},
			},
			setup: func() {},
			teardown: func() {
				sql := `DELETE FROM shortened_urls WHERE short = 'LsE2ypFI'`
				err = connection.WriteDB().Exec(sql).Error
				assert.NoError(t, err)
			},
			err: nil,
		},
		{
			name: "Duplicate key",
			args: args{
				shortenedURL: &entity.ShortenedURL{
					Short:       "LsE2ypFI",
					OriginalURL: "https://www.dcard.tw/",
					CreatedAt:   time.Now(),
					ExpiredAt:   time.Now().Add(time.Hour),
				},
			},
			setup: func() {
				sql := `INSERT INTO "shortened_urls" ("short","original_url","created_at","expired_at") VALUES ('LsE2ypFI','https://www.dcard.tw/',1688221548150,1688225148150)`
				err = connection.WriteDB().Exec(sql).Error
				assert.NoError(t, err)
			},
			teardown: func() {
				sql := `DELETE FROM shortened_urls WHERE short = 'LsE2ypFI'`
				err = connection.WriteDB().Exec(sql).Error
				assert.NoError(t, err)
			},
			err: errors.ErrInternal,
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
			repo := repository.New(connection)

			tt.setup()
			defer tt.teardown()

			err := repo.StoreShortenedURL(ctx, tt.args.shortenedURL)
			if err != nil && tt.err != nil {
				assert.Truef(
					t,
					errors.Is(err, tt.err),
					"StoreShortenedURL() error = %v, expected error %v",
					err,
					tt.err,
				)
			}
		})
	}
}

func TestRepoImpl_FindShortenedURL(t *testing.T) {
	connection, err := db.NewConnection(db.Config{
		Conn: conn.Database{
			Debug:    true,
			Name:     "url_shortener",
			Host:     "localhost",
			User:     "url_shortener_service",
			Port:     5432,
			Password: "pqV7EJ8bYJpFDXXJtw66s6JKG4xpZb4v",
			Type:     conn.Postgres,
		},
	})
	assert.NoError(t, err)

	type args struct {
		short string
	}
	tests := []struct {
		name     string
		args     args
		setup    func()
		teardown func()
		expected *entity.ShortenedURL
		err      error
	}{
		{
			name: "Success",
			args: args{
				short: "LsE2ypFI",
			},
			setup: func() {
				now := time.Now()
				sql := `INSERT INTO "shortened_urls" ("short","original_url","created_at","expired_at") VALUES ('LsE2ypFI','https://www.dcard.tw/',?,?)`
				err = connection.WriteDB().Exec(sql, now.UnixMilli(), now.Add(entity.DefaultShortenedURLExpireDur).UnixMilli()).Error
				assert.NoError(t, err)
			},
			teardown: func() {
				sql := `DELETE FROM shortened_urls WHERE short = 'LsE2ypFI'`
				err = connection.WriteDB().Exec(sql).Error
				assert.NoError(t, err)
			},
			expected: &entity.ShortenedURL{
				Short:       "LsE2ypFI",
				OriginalURL: "https://www.dcard.tw/",
			},
			err: nil,
		},
		{
			name: "ResourceNotFound",
			args: args{
				short: "LsE2ypFI",
			},
			setup:    func() {},
			teardown: func() {},
			expected: nil,
			err:      errors.ErrResourceNotFound,
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

			repo := repository.New(connection)

			tt.setup()
			defer tt.teardown()

			actual, err := repo.FindShortenedURL(ctx, tt.args.short)
			if err != nil && tt.err != nil {
				assert.Truef(
					t,
					errors.Is(err, tt.err),
					"FindShortenedURL(ctx,%v) error = %v, expected error %v",
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
				"FindShortenedURL(%v, %v)",
				ctx,
				tt.args.short,
			)
		})
	}
}
