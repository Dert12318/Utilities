package redis_universal

import (
	"context"
	"encoding"
	"fmt"
	"sync"
	"time"

	"github.com/Dert12318/Utilities/cache"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type (
	Option struct {
		Address      []string
		Password     string
		DB           int
		PoolSize     int
		MinIdleConns int
		ReadOnly     bool
		DialTimeout  time.Duration
		PoolTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		MaxConnAge   time.Duration
	}

	redisUniversalClient struct {
		r        redis.UniversalClient
		mu       sync.Mutex
		channels map[string]cache.PubSub
	}
)

func New(option *Option) (cache.Cache, error) {
	var client redis.UniversalClient

	client = redis.NewUniversalClient(&redis.UniversalOptions{
		DB:           option.DB,
		Addrs:        option.Address,
		Password:     option.Password,
		PoolSize:     option.PoolSize,
		PoolTimeout:  option.PoolTimeout,
		ReadTimeout:  option.ReadTimeout,
		WriteTimeout: option.WriteTimeout,
		DialTimeout:  option.DialTimeout,
		MinIdleConns: option.MinIdleConns,
		MaxConnAge:   option.MaxConnAge,
		ReadOnly:     option.ReadOnly,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "Failed to connect to redis!")
	}

	return &redisUniversalClient{r: client}, nil
}

func (c *redisUniversalClient) Ping(ctx *context.Context) error {
	if _, err := c.r.Ping().Result(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *redisUniversalClient) SetWithExpiration(ctx *context.Context, key string, value interface{}, duration time.Duration) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.Set(key, value, duration).Result(); err != nil {
		return errors.Wrapf(err, "failed to set cache with key %s!", key)
	}
	return nil
}

func (c *redisUniversalClient) Set(ctx *context.Context, key string, value interface{}) error {
	return c.SetWithExpiration(ctx, key, value, 0)
}

func (c *redisUniversalClient) Get(ctx *context.Context, key string, data interface{}) error {
	if _, ok := data.(encoding.BinaryUnmarshaler); !ok {
		return errors.New(fmt.Sprintf("failed to get cache with key %s!: redis: can't unmarshal (implement encoding.BinaryUnmarshaler)", key))
	}

	if err := check(c); err != nil {
		return err
	}

	val, err := c.r.Get(key).Result()

	if err == redis.Nil {
		return errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to get key %s!", key)
	}

	if err := data.(encoding.BinaryUnmarshaler).UnmarshalBinary([]byte(val)); err != nil {
		return err
	}

	return nil
}

func (c *redisUniversalClient) Keys(ctx *context.Context, pattern string) ([]string, error) {
	if err := check(c); err != nil {
		return []string{}, err
	}

	return c.r.Keys(pattern).Result()
}

func (c *redisUniversalClient) HMSetWithExpiration(ctx *context.Context, key string, value map[string]interface{}, ttl time.Duration) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.HMSet(key, value).Result(); err != nil {
		return errors.Wrapf(err, "failed to HMSet cache with key %s!", key)
	}

	if _, err := c.r.Expire(key, ttl).Result(); err != nil {
		c.r.Del(key)
		return errors.Wrapf(err, "failed to HMSet cache with key %s!", key)
	}
	return nil
}

func (c *redisUniversalClient) HMSet(ctx *context.Context, key string, value map[string]interface{}) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.HMSet(key, value).Result(); err != nil {
		return errors.Wrapf(err, "failed to HMSet cache with key %s!", key)
	}
	return nil
}

func (c *redisUniversalClient) HSetWithExpiration(ctx *context.Context, key, field string, value interface{}, ttl time.Duration) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.HSet(key, field, value).Result(); err != nil {
		return errors.Wrapf(err, "failed to HSet cache with key %s!", key)
	}
	if _, err := c.r.Expire(key, ttl).Result(); err != nil {
		c.r.Del(key)
		return errors.Wrapf(err, "failed to HMSet cache with key %s!", key)
	}
	return nil
}

func (c *redisUniversalClient) HSet(ctx *context.Context, key, field string, value interface{}) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.HSet(key, field, value).Result(); err != nil {
		return errors.Wrapf(err, "failed to HSet cache with key %s!", key)
	}
	return nil
}

func (c *redisUniversalClient) HMGet(ctx *context.Context, key string, fields ...string) ([]interface{}, error) {
	if err := check(c); err != nil {
		return nil, err
	}

	val, err := c.r.HMGet(key, fields...).Result()
	if err == redis.Nil {
		return nil, errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get key %s!", key)
	}

	return val, nil
}

func (c *redisUniversalClient) HGetAll(ctx *context.Context, key string) (map[string]string, error) {
	if err := check(c); err != nil {
		return nil, err
	}

	val, err := c.r.HGetAll(key).Result()
	if err == redis.Nil {
		return nil, errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get key %s!", key)
	}

	return val, nil
}

func (c *redisUniversalClient) HGet(ctx *context.Context, key, field string, response interface{}) error {
	if _, ok := response.(encoding.BinaryUnmarshaler); !ok {
		return errors.New(fmt.Sprintf("failed to get cache with key %s!: redis: can't unmarshal (implement encoding.BinaryUnmarshaler)", key))
	}

	if err := check(c); err != nil {
		return err
	}

	val, err := c.r.HGet(key, field).Result()
	if err == redis.Nil {
		return errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to get key %s!", key)
	}

	if err := response.(encoding.BinaryUnmarshaler).UnmarshalBinary([]byte(val)); err != nil {
		return err
	}

	return nil
}

func (c *redisUniversalClient) MGet(ctx *context.Context, key []string) ([]interface{}, error) {
	if err := check(c); err != nil {
		return nil, err
	}

	val, err := c.r.MGet(key...).Result()
	if err == redis.Nil {
		return nil, errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get key %s!", key)
	}

	return val, nil
}

func (c *redisUniversalClient) SAdd(ctx context.Context, key string, values ...interface{}) error {
	if err := check(c); err != nil {
		return err
	}

	if err := c.r.SAdd(key, values...).Err(); err != nil {
		return errors.Wrapf(err, "failed to set cache with key %s!", key)
	}
	return nil
}

func (c *redisUniversalClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	if err := check(c); err != nil {
		return false, err
	}

	val, err := c.r.SIsMember(key, member).Result()
	if err == redis.Nil {
		return false, errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return false, errors.Wrapf(err, "failed to get key %s!", key)
	}

	return val, nil
}

func (c *redisUniversalClient) SMembers(ctx context.Context, key string) ([]string, error) {
	if err := check(c); err != nil {
		return nil, err
	}

	val, err := c.r.SMembers(key).Result()
	if err == redis.Nil {
		return nil, errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get key %s!", key)
	}
	return val, nil
}

func (c *redisUniversalClient) Remove(ctx *context.Context, key string) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.Del(key).Result(); err != nil {
		return errors.Wrapf(err, "failed to remove key %s!", key)
	}

	return nil
}

func (c *redisUniversalClient) RemoveByPattern(ctx *context.Context, pattern string, countPerLoop int64) error {
	if err := check(c); err != nil {
		return err
	}

	iteration := 1
	for {
		keys, _, err := c.r.Scan(0, pattern, countPerLoop).Result()
		if err != nil {
			return errors.Wrapf(err, "failed to scan redis pattern %s!", pattern)
		}

		if len(keys) == 0 {
			break
		}

		if _, err := c.r.Del(keys...).Result(); err != nil {
			return errors.Wrapf(err, "failed iteration-%d to remove key with pattern %s", iteration, pattern)
		}

		iteration++
	}

	return nil
}

func (c *redisUniversalClient) FlushDatabase(ctx *context.Context) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.FlushDB().Result(); err != nil {
		return errors.Wrap(err, "failed to flush db!")
	}

	return nil
}

func (c *redisUniversalClient) FlushAll(ctx *context.Context) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.FlushAll().Result(); err != nil {
		return errors.Wrap(err, "failed to flush db!")
	}

	return nil
}

func (c *redisUniversalClient) Close() error {
	if err := c.r.Close(); err != nil {
		return errors.Wrap(err, "failed to close redis client")
	}

	return nil
}

func check(c *redisUniversalClient) error {
	if c.r == nil {
		return errors.New("redis client is not connected")
	}

	return nil
}

func (c *redisUniversalClient) Client() cache.Cache {
	return c
}

func (c *redisUniversalClient) Pipeline() cache.Pipe {
	return &pipe{instance: c.r.Pipeline()}
}

func (c *redisUniversalClient) Subscribe(channel string) (cache.PubSub, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for c, p := range c.channels {
		if c == channel {
			return p, nil
		}
	}

	p := c.r.Subscribe(channel)
	c.channels[channel] = &pubsub{r: c.r, p: p, cn: channel}
	return c.channels[channel], nil
}
