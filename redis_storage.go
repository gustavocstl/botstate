package botstate

import (
	"github.com/go-redis/redis/v7"
)

//RedisStorage is a default storage of botstate to manipulate user data.
type RedisStorage struct{}

//RedisClient receives an instance of Redis
var RedisClient *redis.Client

//SetValues uses hset to save multiple values
func (r RedisStorage) SetValues(key string, values ...interface{}) error {
	return RedisClient.HSet(key, values...).Err()
}

//GetValue return specific value
func (r RedisStorage) GetValue(key string, valueName string) (string, error) {
	return RedisClient.HGet(key, valueName).Result()
}

//GetAllValues return all values from user
func (r RedisStorage) GetAllValues(key string) (Data, error) {
	return RedisClient.HGetAll(key).Result()
}

//ResetCurrentState clear current state and callback state from user data
func (r RedisStorage) ResetCurrentState(key string) error {
	var err error

	err = r.SetValues(key, "current_state", "")
	err = r.SetValues(key, "state_with_callback", "")

	return err
}
