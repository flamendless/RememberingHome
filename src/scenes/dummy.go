package scenes

import (
	"remembering-home/src/context"

	"github.com/hajimehoshi/ebiten/v2"
)

type Dummy_Scene struct {
	Context      *context.GameContext
	SceneManager SceneManager
}

func (scene *Dummy_Scene) GetName() string {
	return "Dummy"
}

func (scene *Dummy_Scene) GetStateName() string {
	return ""
}

func NewDummyScene(ctx *context.GameContext, sceneManager SceneManager) *Dummy_Scene {
	scene := Dummy_Scene{
		Context:      ctx,
		SceneManager: sceneManager,
	}
	return &scene
}

func (scene *Dummy_Scene) Update() error {
	return nil
}

func (scene *Dummy_Scene) Draw(screen *ebiten.Image) {
}

var _ Scene = (*Dummy_Scene)(nil)
