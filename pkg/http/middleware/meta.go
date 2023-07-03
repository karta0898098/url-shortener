package middleware

import "context"

const (
	ServerInfo = "ServerInfo"
)

func DomainFromContext(ctx context.Context) string {
	return "http://localhost:8080"
}
