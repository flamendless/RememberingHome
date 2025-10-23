package scenes

import (
	"remembering-home/src/context"
	"remembering-home/src/graphics"

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

func (scene *Dummy_Scene) GetItemRenderer() *graphics.ItemRenderer {
	return nil
}

var _ Scene = (*Dummy_Scene)(nil)
