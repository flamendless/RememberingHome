package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Intro_Scene struct {
	GameState        *Game_State
	CurrentStateName string
}

func (scene *Intro_Scene) GetName() string {
	return "Intro"
}

func (scene *Intro_Scene) GetStateName() string {
	return scene.CurrentStateName
}

func NewIntroScene(gameState *Game_State) *Intro_Scene {
	scene := Intro_Scene{
		GameState: gameState,
	}
	return &scene
}

func (scene *Intro_Scene) Update() error {
	return nil
}

func (scene *Intro_Scene) Draw(screen *ebiten.Image) {
}

var _ Scene = (*Intro_Scene)(nil)
