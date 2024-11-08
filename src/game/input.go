package game

import (
	"remembering-home/src/conf"

	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionMoveUp input.Action = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionEnter
)

const (
	DevToggleTexts input.Action = iota
	DevToggleLines
	DevGoToDummy
	DevGoToSplash
	DevGoToMainMenu
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
		ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
		ActionMoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		ActionEnter:     {input.KeyGamepadStart, input.KeyEnter, input.KeySpace, input.KeyE},
	}
	handler := system.NewHandler(0, keymap)
	return handler
}

func NewInputHandlerDev(system *input.System) *input.Handler {
	if !conf.DEV {
		panic("DEV mode is not activated. Can't use this feature")
	}
	keymap := input.Keymap{
		DevToggleTexts:  {input.KeyD},
		DevToggleLines:  {input.KeyL},
		DevGoToDummy:    {input.Key1},
		DevGoToSplash:   {input.Key2},
		DevGoToMainMenu: {input.Key3},
	}
	handler := system.NewHandler(1, keymap)
	return handler
}
