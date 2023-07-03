package bloom

import (
	"context"
)

//go:generate mockery --name Filter --with-expecter
type Filter interface {
	GetFilterNamespace() string
	Add(ctx context.Context, item interface{})
	Exist(ctx context.Context, item interface{}) bool
}
