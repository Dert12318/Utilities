package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type (
	Pipe interface {
		Set(key string, value interface{}) error
		SetWithExpiration(key string, value interface{}, expired time.Duration) error
		Get(key string, object interface{}) error
		Exec() error
	}

	PubSub interface {
		Receive() error
		Publish(message string) error
		Channel() <-chan *redis.Message
		Close() error
	}

	Cache interface {
		Ping(ctx *context.Context) error

		// SetWithExpiration value must implement encoding.BinaryMarshaler
		SetWithExpiration(ctx *context.Context, key string, value interface{}, duration time.Duration) error
		// Set value must implement encoding.BinaryMarshaler
		Set(ctx *context.Context, key string, value interface{}) error
		// Set data must implement encoding.BinaryUnmarshaler
		Get(ctx *context.Context, key string, data interface{}) error

		HMSetWithExpiration(ctx *context.Context, key string, value map[string]interface{}, ttl time.Duration) error
		HMSet(ctx *context.Context, key string, value map[string]interface{}) error
		HMGet(ctx *context.Context, key string, fields ...string) ([]interface{}, error)

		HGetAll(ctx *context.Context, key string) (map[string]string, error)

		HSetWithExpiration(ctx *context.Context, key string, field string, value interface{}, ttl time.Duration) error
		HSet(ctx *context.Context, key string, field string, value interface{}) error
		HGet(ctx *context.Context, key string, field string, response interface{}) error

		SAdd(ctx context.Context, key string, values ...interface{}) error
		SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
		SMembers(ctx context.Context, key string) ([]string, error)

		MGet(ctx *context.Context, key []string) ([]interface{}, error)

		Keys(ctx *context.Context, pattern string) ([]string, error)

		Remove(ctx *context.Context, key string) error
		RemoveByPattern(ctx *context.Context, pattern string, countPerLoop int64) error
		FlushDatabase(ctx *context.Context) error
		FlushAll(ctx *context.Context) error
		Close() error

		Pipeline() Pipe
		Client() Cache

		// Implemented by Redis to redis pubsub
		Subscribe(channel string) (PubSub, error)
	}

	PoolCallback func(client Cache)

	Pool interface {
		Use(callback PoolCallback)
		Client() Cache
		Close() error
	}
)
