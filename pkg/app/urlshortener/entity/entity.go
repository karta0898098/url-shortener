package entity

import (
	"time"
)

const (
	DefaultShortenedURLExpireDur = time.Hour * 24 * 3 // DefaultShortenedURLExpireDur
)

// ShortenedURL define shortened URL entity
type ShortenedURL struct {
	// Short is shortened URL and Primary key
	Short string `gorm:"column:short"`

	// OriginalURL original URL which input by the user
	OriginalURL string `gorm:"column:original_url"`

	// CreatedAt the url created at
	CreatedAt time.Time `gorm:"column:created_at"`

	// ExpiredAt the url expired at
	ExpiredAt time.Time `gorm:"column:expired_at"`
}

func NewShortenedURL(short string, originalURL string, expiredAt time.Time) *ShortenedURL {
	t := time.Now()

	return &ShortenedURL{
		Short:       short,
		OriginalURL: originalURL,
		CreatedAt:   t,
		ExpiredAt:   expiredAt,
	}
}
