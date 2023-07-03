package bloom

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestNewRedisFilter(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
	}{
		{
			name:      "Success",
			namespace: "bf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rds := redis.NewClient(&redis.Options{})
			NewRedisFilter(tt.namespace, rds)
			rds.Del(ctx, tt.namespace)

		})
	}
}
