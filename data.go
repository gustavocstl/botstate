package botstate

import (
	"strconv"

	"github.com/go-redis/redis/v7"
)

//Data are used to store all current user data
type Data map[string]string

//Storager is the interface that wraps the methods to manipulate current user data
type Storager interface {
	SetValues(key string, values ...interface{}) error
	GetValue(key string, valueName string) (string, error)
	GetAllValues(key string) (Data, error)
	ResetCurrentState(key string) error
}

//StorageClient is a global variable that receives current Storager
//All functions in BotData struct use StorageClient
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

//BotData are used to control and manipulate current user data
//Passed to Bot struct
type BotData struct {
	UserID string

	//Current has default keys with values in map[string]:
	//user_id is used to save current user id
	//current_state is used to save user's current state
	//state_with_callback is used to save name of state with callback to be executed in next call
	//messages is used to save all messages during the execution flow
	Current Data
}

//User initialize user's data if not exists and set the user_id in data with the argument userID.
//Returns error if exists
func (bd *BotData) User(userID int) error {
	var err error

	bd.UserID = strconv.Itoa(userID)

	err = bd.SetData(Data{
		"user_id": bd.UserID,
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
	return StorageClient.GetValue(bd.UserID, "current_state")
}

//SetStateWithCallback define new state to state_with_callback in user's data.
//This state is used to define callback in state execution.
func (bd *BotData) SetStateWithCallback(name string) error {
	err := bd.SetData(Data{
		"state_with_callback": name,
	})

	return err
}

//GetStateWithCallback get the state_with_callback value from the user's data.
//This state is used to define callback in state execution.
func (bd *BotData) GetStateWithCallback() (string, error) {
	return StorageClient.GetValue(bd.UserID, "state_with_callback")
}

//SetData accepts type Data struct to define multiple values in user's data.
func (bd *BotData) SetData(values Data) error {
	var err error

	for key, value := range values {
		err = StorageClient.SetValues(bd.UserID, key, value)

		if err != nil {
			return err
		}

		bd.Current, err = bd.GetData()
	}

	return err
}

//GetData get all user's data.
func (bd *BotData) GetData() (Data, error) {
	return StorageClient.GetAllValues(bd.UserID)
}

//ResetCurrentState clear current state and callback state from user data
func (bd *BotData) ResetCurrentState() error {
	var err error

	err = StorageClient.ResetCurrentState(bd.UserID)

	if err != nil {
		return err
	}

	bd.Current, err = bd.GetData()

	return err
}
