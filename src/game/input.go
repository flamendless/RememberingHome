package game

import (
	"remembering-home/src/conf"
	"remembering-home/src/context"

	input "github.com/quasilyte/ebitengine-input"
)

func NewInputSystem() *input.System {
	system := input.System{}
	system.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	return &system
}

func NewInputHandler(system *input.System) *input.Handler {
	keymap := input.Keymap{
		context.ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
		context.ActionMoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
		context.ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		context.ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		context.ActionEnter:     {input.KeyGamepadStart, input.KeyEnter, input.KeySpace, input.KeyE},
	}
	handler := system.NewHandler(0, keymap)
	return handler
}

func NewInputHandlerDev(system *input.System) *input.Handler {
	if !conf.DEV {
		panic("DEV mode is not activated. Can't use this feature")
	}
	keymap := input.Keymap{
		context.DevToggleTexts:  {input.KeyD},
		context.DevToggleLines:  {input.KeyL},
		context.DevGoToDummy:    {input.Key1},
		context.DevGoToSplash:   {input.Key2},
		context.DevGoToMainMenu: {input.Key3},
	}
	handler := system.NewHandler(1, keymap)
	return handler
}
