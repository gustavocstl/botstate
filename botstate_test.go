package botstate_test

import (
	"strconv"
	"testing"

	"github.com/gucastiliao/botstate"
	"github.com/stretchr/testify/assert"
)

func TestStateExecution(t *testing.T) {
	mockRedis()

	userID := 111
	defaultFunc := func(bot *botstate.Bot) bool {
		u, _ := strconv.Atoi(bot.Data.UserID)
		return u == userID
	}

	states := []botstate.State{
		{
			Name:     "start",
			Executes: defaultFunc,
		},
		{
			Name:     "end",
			Executes: defaultFunc,
		},
	}

	bot := botstate.New(states)
	bot.Data.User(userID)

	for _, state := range states {
		execute, err := bot.ExecuteState(state.Name)

		assert.Nil(t, err)
		assert.True(t, execute)
	}
}

func TestStateWithCallbacksExecution(t *testing.T) {
	mockRedis()

	userID := 111
	timesToSimulateFailCallback := 2
	failCallbackTimes := 0
	productName := "Product Teste"

	addProduct := func(bot *botstate.Bot) bool {
		return true
	}
	getProductName := func(bot *botstate.Bot) bool {
		if failCallbackTimes < timesToSimulateFailCallback {
			failCallbackTimes++
			return false
		}

		bot.Data.SetData(botstate.Data{
			"product_name": productName,
		})

		return true
	}
	confirmation := func(bot *botstate.Bot) bool {
		return true
	}

	states := []botstate.State{
		{
			Name:     "add_product",
			Executes: addProduct,
			Callback: getProductName,
			Next:     "confirmation",
		},
		{
			Name:     "confirmation",
			Executes: confirmation,
			Next:     "end",
		},
	}

	bot := botstate.New(states)

	bot.Data.User(userID)

	//First execution
	//Only to set state_with_callback to execute callback getProductName in next ExecuteState call
	execute, err := bot.ExecuteState("add_product")

	assert.Nil(t, err)
	assert.True(t, execute)
	assert.Equal(t, "add_product", bot.Data.Current["state_with_callback"])
	assert.Equal(t, "confirmation", bot.Data.Current["current_state"])

	//Loop to simulate fail state in getProductName
	//When the callback fails, ExecuteState must execute the state again
	//Until the callback is valid
	for failCallbackTimes < timesToSimulateFailCallback {
		execute, err = bot.ExecuteState(bot.Data.Current["current_state"])

		assert.Nil(t, err)
		assert.False(t, execute)
		assert.Equal(t, "confirmation", bot.Data.Current["current_state"])
	}

	//When state execution returns true, the value of product_name must be valid
	execute, err = bot.ExecuteState(bot.Data.Current["current_state"])

	assert.Nil(t, err)
	assert.True(t, execute)
	assert.Equal(t, productName, bot.Data.Current["product_name"])
	assert.Equal(t, "end", bot.Data.Current["current_state"])
}

func TestForceChangeStateExecution(t *testing.T) {
	mockRedis()

	userID := 111

	start := func(bot *botstate.Bot) bool {
		return true
	}
	end := func(bot *botstate.Bot) bool {
		bot.Data.ResetCurrentState()
		bot.ExecuteState("lost_state")

		return false
	}
	lost := func(bot *botstate.Bot) bool {
		bot.Data.SetData(botstate.Data{
			"lost_value": "ok",
		})

		return true
	}

	states := []botstate.State{
		{
			Name:     "start",
			Executes: start,
			Next:     "end",
		},
		{
			Name:     "end",
			Executes: end,
			Next:     "start",
		},
		{
			Name:     "lost_state",
			Executes: lost,
		},
	}

	bot := botstate.New(states)

	bot.Data.User(userID)

	_, err := bot.ExecuteState("start")

	assert.Nil(t, err)

	_, err = bot.ExecuteState(bot.Data.Current["current_state"])

	assert.Nil(t, err)
	assert.Equal(t, strconv.Itoa(userID), bot.Data.Current["user_id"])
	assert.Equal(t, "lost_state", bot.Data.Current["current_state"])
	assert.Equal(t, "ok", bot.Data.Current["lost_value"])
}

func TestAddMessageInStateExecution(t *testing.T) {
	mockRedis()

	userID := 111
	messages := []string{
		"Hello!",
		"See, I can send multiple messages.",
		"That's nice!",
	}

	bot := botstate.New([]botstate.State{
		{
			Name: "start",
			Executes: func(bot *botstate.Bot) bool {
				for _, m := range messages {
					//Loop to test multiple call
					bot.AddMessage([]string{m})
				}

				return true
			},
		},
	})

	bot.Data.User(userID)

	execute, err := bot.ExecuteState("start")

	assert.Nil(t, err)
	assert.True(t, execute)

	botMessages := bot.GetMessages()
	assert.Equal(t, 3, len(botMessages))

	for k, m := range messages {
		assert.Equal(t, m, botMessages[k])
	}
}

func TestAddMessageFail(t *testing.T) {
	mockRedis()

	userID := 111

	bot := botstate.New([]botstate.State{
		{
			Name: "start",
			Executes: func(bot *botstate.Bot) bool {
				err := bot.AddMessage([]string{})

				assert.Equal(t, "undefined messages", err.Error())

				return true
			},
		},
	})

	bot.Data.User(userID)

	execute, err := bot.ExecuteState("start")

	assert.Nil(t, err)
	assert.True(t, execute)

	botMessages := bot.GetMessages()
	assert.Equal(t, 0, len(botMessages))
}

func TestStateExecutionFail(t *testing.T) {
	mockRedis()

	userID := 111

	states := []botstate.State{}

	bot := botstate.New(states)
	err := bot.Data.User(userID)

	assert.Nil(t, err)

	execute, err := bot.ExecuteState("add_product")

	assert.False(t, execute)
	assert.Equal(t, "No state to execute with name add_product.", err.Error())
}

func TestStateExecutionWithoutUser(t *testing.T) {
	states := []botstate.State{
		{Name: "add_product"},
	}

	bot := botstate.New(states)

	execute, err := bot.ExecuteState("add_product")

	assert.False(t, execute)
	assert.Equal(t, "Undefined user to execute state add_product.", err.Error())
}

func TestStateExecutionWithEmptyMethod(t *testing.T) {
	states := []botstate.State{
		{Name: "add_product"},
	}

	userID := 111
	bot := botstate.New(states)
	err := bot.Data.User(userID)

	assert.Nil(t, err)

	execute, err := bot.ExecuteState("add_product")

	assert.False(t, execute)
	assert.Equal(t, "Method to execute in the add_product state is nil.", err.Error())
}
