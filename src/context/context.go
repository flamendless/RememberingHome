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
)

const (
	DevToggleTexts input.Action = iota
	DevToggleLines
	DevGoToDummy
	DevGoToSplash
	DevGoToMainMenu
)

type GameContext struct {
	Loader          *resource.Loader
	InputSystem     *input.System
	InputHandler    *input.Handler
	InputHandlerDev *input.Handler
	Settings        *conf.Settings
}

func NewGameContext(
	loader *resource.Loader,
	inputSystem *input.System,
	inputHandler *input.Handler,
	inputHandlerDev *input.Handler,
	settings *conf.Settings,
) *GameContext {
	return &GameContext{
		Loader:          loader,
		InputSystem:     inputSystem,
		InputHandler:    inputHandler,
		InputHandlerDev: inputHandlerDev,
		Settings:        settings,
	}
}
