package context

import (
	"remembering-home/src/conf"

	input "github.com/quasilyte/ebitengine-input"
	resource "github.com/quasilyte/ebitengine-resource"
)

const (
	ActionMoveUp input.Action = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionEnter
	ActionBack
	ActionZoomIn
	ActionZoomOut
	ActionZoomReset
)

type GameContext struct {
	Loader       *resource.Loader
	InputSystem  *input.System
	InputHandler *input.Handler
	Settings     *conf.Settings
}

func NewGameContext(
	loader *resource.Loader,
	inputSystem *input.System,
	inputHandler *input.Handler,
	settings *conf.Settings,
) *GameContext {
	return &GameContext{
		Loader:       loader,
		InputSystem:  inputSystem,
		InputHandler: inputHandler,
		Settings:     settings,
	}
}
