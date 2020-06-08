package botstate

import (
	"errors"
)

type State struct {
	Name     string
	Executes func(bot *Bot) bool
	Callback func(bot *Bot) bool
	Next     string
}

type Bot struct {
	States []State
	Data   *BotData
}

//New returns new Bot struct with BotData
func New(states []State) *Bot {
	return &Bot{
		States: states,
		Data:   &BotData{},
	}
}

//ExecuteState define the current_state using argument name in the user's data.
//
//If exists callback in user's data (state_with_callback), execute it first with the executeCallback method.
//Terminates execution if callback returns false.
//
//If there is no callback, the flow continues and if the current state has a method in the item State.Callback, this value will be defined in the user's current data (state_with_callback) to be executed later.
//
//After all checks, the method in State.Executes is executed.
//The current state is defined using the value of State.Next if execution return true.
//
//Return execution boolean and error if exists.
func (b *Bot) ExecuteState(name string) (bool, error) {
	for _, state := range b.States {
		if state.Name == name {
			if b.Data.UserId == "" {
				return false, errors.New("Undefined user to execute state " + state.Name + ".")
			}

			if state.Executes == nil {
				return false, errors.New("Method to execute in the " + state.Name + " state is nil.")
			}

			b.Data.SetCurrentState(state.Name)

			callbackResp, err := b.executeCallback()

			if err != nil {
				return false, err
			}

			if callbackResp == false {
				return false, nil
			}

			if state.Callback != nil {
				err := b.Data.SetStateWithCallback(state.Name)

				if err != nil {
					return false, err
				}
			}

			execute := state.Executes(b)

			if execute == true && state.Next != "" {
				b.Data.SetCurrentState(state.Next)
			}

			return execute, nil
		}
	}

	return false, errors.New("No state to execute with name " + name + ".")
}

//executeCallback will get state_with_callback from user's data.
//And execute the executeCallbackFromState method passing state name as argument.
//Return callback execution boolean and error if exists.
func (b *Bot) executeCallback() (bool, error) {
	stateWithCallback, _ := b.Data.GetStateWithCallback()

	if stateWithCallback != "" {
		return b.executeCallbackFromState(stateWithCallback), nil
	}

	return true, nil
}

//executeCallbackFromState will loop through all states to find the state with argument name.
//Checks if the state has method in State.Callback.
//Execute method in State.Callback.
//Return callback boolean response.
func (b *Bot) executeCallbackFromState(name string) bool {
	for _, state := range b.States {
		if state.Name == name {
			if state.Callback != nil {
				return state.Callback(b)
			}
		}
	}

	return true
}
