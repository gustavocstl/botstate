package botstate

import (
	"strconv"

	"github.com/go-redis/redis/v7"
)

type Data map[string]string

type Storager interface {
	SetValues(key string, values ...interface{}) error
	GetValue(key string, valueName string) (string, error)
	GetAllValues(key string) (Data, error)
	ResetAll(key string) error
}

var StorageClient Storager

//SetStorageClient accept as argument any storage client that meets the needs of the Storager interface.
func SetStorageClient(client Storager) {
	StorageClient = client
}

//DefaultStorage get Redis Storage to manage bot data.
//Returns Storager
func DefaultStorage(redisClient *redis.Client) Storager {
	RedisClient = redisClient
	_, err := RedisClient.Ping().Result()

	if err != nil {
		panic(err)
	}

	return RedisStorage{}
}

type BotData struct {
	UserId  string
	Current Data
}

//User initialize user's data if not exists and set the user_id in data with the argument userId.
//Returns error if exists
func (bd *BotData) User(userId int) error {
	var err error

	bd.UserId = strconv.Itoa(userId)

	err = bd.SetData(Data{
		"user_id": bd.UserId,
	})

	if err != nil {
		return err
	}

	bd.Current, err = bd.GetData()

	if err != nil {
		return err
	}

	return nil
}

//SetCurrentState define new state to current_state in user's data.
//This state is used to control where the user is in the bot flow.
func (bd *BotData) SetCurrentState(name string) (string, error) {
	err := bd.SetData(Data{
		"current_state": name,
	})

	return name, err
}

//GetCurrentState get the current_state value from the user's data.
//This state is used to control where the user is in the bot flow.
func (bd *BotData) GetCurrentState() (string, error) {
	return StorageClient.GetValue(bd.UserId, "current_state")
}

//SetStateWithCallback define new state to state_with_callback in user's data.
//This state is used to define callback in state execution.
func (bd *BotData) SetStateWithCallback(name string) error {
	err := bd.SetData(Data{
		"state_with_callback": name,
	})

	return err
}

//GetCurrentState get the state_with_callback value from the user's data.
//This state is used to define callback in state execution.
func (bd *BotData) GetStateWithCallback() (string, error) {
	return StorageClient.GetValue(bd.UserId, "state_with_callback")
}

//SetData accepts type Data struct to define multiple values in user's data.
func (bd *BotData) SetData(values Data) error {
	var err error

	for key, value := range values {
		err = StorageClient.SetValues(bd.UserId, key, value)

		if err != nil {
			return err
		}

		bd.Current, err = bd.GetData()
	}

	return err
}

//GetData get all user's data.
func (bd *BotData) GetData() (Data, error) {
	return StorageClient.GetAllValues(bd.UserId)
}

//ResetAll reset all current user's data
func (bd *BotData) ResetAll() error {
	var err error

	err = StorageClient.ResetAll(bd.UserId)

	if err != nil {
		return err
	}

	bd.Current, err = bd.GetData()

	return err
}
