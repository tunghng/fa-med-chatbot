package rediscmd

import (
	"encoding/json"
	"github.com/spf13/viper"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
)

// RedisCmd is wrapper of "github.com/go-redis/redis/v7".
type RedisCmd interface {
	Set(key string, value interface{}, expireTime int64) error
	Get(key string, value interface{}) (interface{}, error)
	Del(key string) (int64, error)
	Expire(key string, expire int64) (bool, error)
	Ping() error
	Exist(key string) (bool, error)
	Incr(key string) (int64, error)
	TTL(key string) (int64, error)
	DeletePattern(pattern string) error
}

type Config struct {
	Address  string
	Password string
	DB       int
}

type redisCmd struct {
	client *redis.Client
}

// redisCmd ...
func NewRedisCmd(cfgReader *viper.Viper) (RedisCmd, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfgReader.GetString(cfg.KeyDBRedisAddr),
		Password: cfgReader.GetString(cfg.KeyDBRedisPassword),
		DB:       cfgReader.GetInt(cfg.KeyDBRedisDB),
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &redisCmd{
		client: client,
	}, nil
}

func (_this *redisCmd) Ping() error {
	if err := _this.client.Ping().Err(); err != nil {
		return errors.Wrap(err, "checking redis client")
	}

	return nil
}

func (_this *redisCmd) Set(key string, value interface{}, expireTime int64) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = _this.client.Set(key, data, time.Duration(expireTime)*time.Second).Result()
	return err
}

func (_this *redisCmd) DeletePattern(pattern string) error {
	keys, err := _this.client.Keys(pattern).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = _this.client.Del(key).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (_this *redisCmd) Get(key string, value interface{}) (interface{}, error) {
	data, err := _this.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), value)
	return value, err
}

func (_this *redisCmd) Del(key string) (int64, error) {
	return _this.client.Del(key).Result()
}

func (_this *redisCmd) Expire(key string, expire int64) (bool, error) {
	value, err := _this.client.Expire(key, time.Duration(expire)*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return value, nil
}

func (_this *redisCmd) Exist(key string) (bool, error) {
	value, err := _this.client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return value == 1, nil
}

func (_this *redisCmd) Incr(key string) (int64, error) {
	result, err := _this.client.Incr(key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (_this *redisCmd) TTL(key string) (int64, error) {
	result, err := _this.client.TTL(key).Result()
	if err != nil {
		return 0, err
	}
	return int64(result.Seconds()), nil
}
