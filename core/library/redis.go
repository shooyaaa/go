package library

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	rdb redis.UniversalClient
}

func (r *Redis) Init(params map[string]interface{}) {
	options := redis.UniversalOptions{}
	options.Addrs = params["address"].([]string)
	master, ok := params["master"].(string)
	if ok {
		options.MasterName = master
	}
	r.rdb = redis.NewUniversalClient(&options)
}

func (r Redis) GetString(key string) string {
	return r.Get(key).String()
}

func (r Redis) GetInt(key string) (int, error) {
	return r.Get(key).Int()
}

func (r Redis) GetInt64(key string) (int64, error) {
	return r.Get(key).Int64()
}

func (r Redis) Get(key string) *redis.StringCmd {
	ctx := context.Background()
	return r.rdb.Get(ctx, key)
}

func (r Redis) SetWithTimeout(key string, value interface{}, duration time.Duration) (string, error) {
	ctx := context.Background()
	return r.rdb.Set(ctx, key, value, duration).Result()
}

func (r Redis) Set(key string, value interface{}) (string, error) {
	return r.SetWithTimeout(key, value, 0)
}

func (r Redis) SetInt(key string, value int) error {
	_, err := r.SetWithTimeout(key, value, 0)
	return err
}

func (r Redis) SetInt64(key string, value int64) error {
	_, err := r.SetWithTimeout(key, value, 0)
	return err
}

func (r Redis) SetString(key string, value string) error {
	_, err := r.SetWithTimeout(key, value, 0)
	return err
}
func (r Redis) Delete(key string) error {
	return r.Delete(key)
}
