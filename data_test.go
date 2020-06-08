package botstate_test

import (
	"strconv"
	"testing"

	"github.com/gucastiliao/botstate"
	"github.com/stretchr/testify/assert"
)

func TestInitializeData(t *testing.T) {
	mockRedis()

	userId := 111
	botData := &botstate.BotData{}

	err := botData.User(userId)

	assert.Nil(t, err)
	assert.Equal(t, strconv.Itoa(userId), botData.Current["user_id"])
}

func TestCurrentSate(t *testing.T) {
	mockRedis()

	userId := 111
	botData := &botstate.BotData{}

	botData.User(userId)
	current, err := botData.SetCurrentState("state_test")

	assert.Nil(t, err)
	assert.Equal(t, "state_test", current)

	state, err := botData.GetCurrentState()

	assert.Nil(t, err)
	assert.Equal(t, "state_test", state)
}

func TestStateWithCallback(t *testing.T) {
	mockRedis()

	userId := 111
	botData := &botstate.BotData{}

	botData.User(userId)
	err := botData.SetStateWithCallback("state_test")

	assert.Nil(t, err)

	state, err := botData.GetStateWithCallback()

	assert.Nil(t, err)
	assert.Equal(t, "state_test", state)
}

func TestSetData(t *testing.T) {
	mockRedis()

	userId := 111
	botData := &botstate.BotData{}

	botData.User(userId)
	err := botData.SetData(botstate.Data{
		"test_data": "test_value",
	})

	assert.Nil(t, err)
	assert.Equal(t, "test_value", botData.Current["test_data"])
}

func TestResetAllData(t *testing.T) {
	mockRedis()

	userId := 111
	data := botstate.Data{
		"product_name":     "Product Test",
		"product_price":    "10",
		"product_quantity": "20",
	}

	botData := &botstate.BotData{}

	botData.User(userId)

	err := botData.SetData(data)
	assert.Nil(t, err)

	for key, value := range data {
		assert.Equal(t, value, botData.Current[key])
	}

	err = botData.ResetAll()
	assert.Nil(t, err)

	for key := range data {
		assert.Empty(t, botData.Current[key])
	}
}
