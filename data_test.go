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

func TestResetCurrentState(t *testing.T) {
	mockRedis()

	userId := 111
	data := botstate.Data{
		"product_name":     "Product Test",
		"product_price":    "10",
		"product_quantity": "20",
	}

	botData := &botstate.BotData{}

	botData.User(userId)
	botData.SetCurrentState("test_state")
	botData.SetStateWithCallback("test_state")
	botData.SetData(data)

	err := botData.ResetCurrentState()
	assert.Nil(t, err)
	assert.Empty(t, botData.Current["current_state"])
	assert.Empty(t, botData.Current["state_with_callback"])

	for key := range data {
		assert.NotEmpty(t, botData.Current[key])
	}
}
