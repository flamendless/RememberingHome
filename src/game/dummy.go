package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Dummy_Scene struct {
	GameState *Game_State
}

func (scene Dummy_Scene) GetName() string {
	return "Dummy"
}

func (scene Dummy_Scene) GetStateName() string {
	return ""
}

func NewDummyScene(gameState *Game_State) *Dummy_Scene {
	scene := Dummy_Scene{
		GameState: gameState,
	}
	return &scene
}

func (scene *Dummy_Scene) Update() error {
	return nil
}

func (scene *Dummy_Scene) Draw(screen *ebiten.Image) {
}
