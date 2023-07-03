-- +goose Up
CREATE TABLE IF NOT EXISTS shortened_urls
(
    short        varchar(8)   NOT NULL UNIQUE PRIMARY KEY,
    original_url varchar(500) NOT NULL,
    created_at   BIGINT       NOT NULL,
    expired_at   BIGINT       NOT NULL
);

COMMENT ON COLUMN shortened_urls.short IS 'Short is shortened URL and Primary key';
COMMENT ON COLUMN shortened_urls.original_url IS 'OriginalURL original URL which input by the user';
COMMENT ON COLUMN shortened_urls.created_at IS 'CreatedAt the url created at';
COMMENT ON COLUMN shortened_urls.expired_at IS 'ExpiredAt the url expired at';

