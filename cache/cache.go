package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Init()
	Get(key string) ([]byte, error)
	Set(key string, bytes []byte, expiration time.Duration) error
	Expire(key string) error
}

type RedisCache struct {
	RedisConnectionUrl string
	redisClient        *redis.Client
}

func (rc *RedisCache) Init() {
	if rc.redisClient != nil {
		log.Println("redis client has already been initialized")
		return
	}

	opt, err := redis.ParseURL(rc.RedisConnectionUrl)
	if err != nil {
		panic(err)
	}

	rc.redisClient = redis.NewClient(opt)
}

func (rc *RedisCache) Get(key string) ([]byte, error) {
	ctx := context.Background()
	return rc.redisClient.Get(ctx, key).Bytes()
}

func (rc *RedisCache) Set(key string, bytes []byte, expiration time.Duration) error {
	ctx := context.Background()
	return rc.redisClient.Set(ctx, key, bytes, expiration).Err()
}

func (rc *RedisCache) Expire(key string) error {
	ctx := context.Background()
	return rc.redisClient.Expire(ctx, key, 0).Err()
}

func Store[V interface{}](cache Cache, key string, value V, expiration time.Duration) error {
	bytes, err := json.Marshal(&value)

	if err != nil {
		log.Println("unable to marshall a value: [", value, "] to bytes array for a key[", key, "]")
		return errors.New("unable to marshall a value")
	}

	err = cache.Set(key, bytes, expiration)

	if err != nil {
		log.Println("unable to store in Redis a value for a key[", key, "]")
		return errors.New("unable to store in Redis a value")
	}

	return err
}

func Get[V any](cache Cache, key string) (*V, error) {
	bytes, err := cache.Get(key)

	if err != nil {
		log.Println("value for a key [", key, "] is not found: ", err)
		return nil, errors.New(fmt.Sprintf("value is not found for a key[%s]", key))
	}

	value := new(V)
	err = json.Unmarshal(bytes, &value)

	if err != nil {
		log.Println("error during unmarshalling value for a key [", key, "]", err)
		return nil, errors.New(fmt.Sprintf("value is not found for a key[%s]", key))
	}

	return value, nil
}

func GetAndStoreIfMissed[V any](cache Cache, key string, supplier func() (*V, time.Duration, error)) (*V, error) {
	value, err := Get[V](cache, key)

	if err != nil {
		log.Println("value for a key [", key, "] is not found, using a supplier: ", err)
		supplierValue, expiration, err := supplier()

		if err != nil {
			log.Println("supplier returned an error for a key [", key, "]: ", err)
			return nil, errors.New(err.Error())
		}

		err = Store[V](cache, key, *supplierValue, expiration)

		if err != nil {
			log.Println("unable to store a value for a key [", key, "]: ", err)
			return nil, errors.New(err.Error())
		}

		value = supplierValue
	}

	return value, nil
}
