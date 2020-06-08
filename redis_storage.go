package botstate

import (
	"github.com/go-redis/redis/v7"
)

type RedisStorage struct{}

var RedisClient *redis.Client

func (r RedisStorage) SetValues(key string, values ...interface{}) error {
	return RedisClient.HSet(key, values...).Err()
}

func (r RedisStorage) GetValue(key string, valueName string) (string, error) {
	return RedisClient.HGet(key, valueName).Result()
}

func (r RedisStorage) GetAllValues(key string) (Data, error) {
	return RedisClient.HGetAll(key).Result()
}

func (r RedisStorage) ResetAll(key string) error {
	return RedisClient.FlushAll().Err()
}
