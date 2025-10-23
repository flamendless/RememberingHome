package game

import (
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
		context.ActionBack:      {input.KeyGamepadBack, input.KeyEscape, input.KeyBackspace},
		context.ActionZoomIn:    {input.KeyEqual},
		context.ActionZoomOut:   {input.KeyMinus},
		context.ActionZoomReset: {input.Key0},
	}
	handler := system.NewHandler(0, keymap)
	return handler
}
