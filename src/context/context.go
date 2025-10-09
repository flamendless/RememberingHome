package context

import (
	input "github.com/quasilyte/ebitengine-input"
	resource "github.com/quasilyte/ebitengine-resource"
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

type GameContext struct {
	Loader          *resource.Loader
	InputSystem     *input.System
	InputHandler    *input.Handler
	InputHandlerDev *input.Handler
}

func NewGameContext(
	loader *resource.Loader,
	inputSystem *input.System,
	inputHandler *input.Handler,
	inputHandlerDev *input.Handler,
) *GameContext {
	return &GameContext{
		Loader:          loader,
		InputSystem:     inputSystem,
		InputHandler:    inputHandler,
		InputHandlerDev: inputHandlerDev,
	}
}
