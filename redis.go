package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"time"
)

type RedisWrapper struct {
	Client         *redis.Client
	ExpirationTime time.Duration
}

type RedisConfig struct {
	Addr string
	DB   int
	Pass string
}

func NewRedisWrapper(config *RedisConfig) (*RedisWrapper, error) {
	opt := &redis.Options{
		Addr:     config.Addr,
		DB:       config.DB,
		Password: config.Pass,
	}
	client := redis.NewClient(opt)

	return &RedisWrapper{
		Client: client,
	}, nil
}

func (wrapper *RedisWrapper) Get(
	ctx context.Context,
	key string,
	value interface{},
) (bool, error) {
	result, err := wrapper.Client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if len(result) == 0 {
		return false, nil
	}

	err = msgpack.Unmarshal(result, value)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (wrapper *RedisWrapper) Set(
	ctx context.Context,
	key string,
	value interface{},
	expTime *time.Duration,
) error {
	if expTime == nil {
		expTime = &wrapper.ExpirationTime
	}
	bdata, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}

	_, err = wrapper.Client.Set(ctx, key, bdata, *expTime).Result()
	if err != nil {
		return err
	}

	return nil
}

func (wrapper *RedisWrapper) MGet(ctx context.Context, cKeys []string) (map[string]string, error) {
	client := wrapper.Client
	results := map[string]string{}
	result := client.MGet(ctx, cKeys...)
	if result.Err() != nil {
		return nil, nil
	}
	values := result.Val()
	for i := 0; i < len(cKeys); i++ {
		ckey := cKeys[i]
		if values[i] == nil {
			results[ckey] = ""
		} else {
			results[ckey] = fmt.Sprintf("%s", values[i])
		}
	}
	return results, nil
}

//AcquireLock return value of expire time if can acquire lock
func (wrapper *RedisWrapper) AcquireLock(
	ctx context.Context,
	lockKey string,
	lockTime time.Duration,
) string {
	client := wrapper.Client

	cKey := "rlock:" + lockKey
	nowTs := time.Now().Add(lockTime)

	lockValue := ParseToISOTime(nowTs)
	success := client.SetNX(ctx, cKey, lockValue, lockTime).Val()
	if success {
		val := client.Get(ctx, cKey).Val()
		if ParseToISOTime(time.Now()) > val {
			return ""
		}

		return lockValue
	}
	return ""
}

func (wrapper *RedisWrapper) ReleaseLock(ctx context.Context, lockKey string) {
	client := wrapper.Client

	cKey := "rlock:" + lockKey
	client.Del(ctx, cKey)
}

func (wrapper *RedisWrapper) applyTimeout(ctx context.Context, cKey string, timeout time.Duration) {
	_, err := wrapper.Client.Expire(ctx, cKey, timeout).Result()
	if err != nil {
		log.Printf("error when want timeout %v", err)
	}
}
