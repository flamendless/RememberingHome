package scenes

import (
	"remembering-home/src/context"
	"remembering-home/src/graphics"

	"github.com/hajimehoshi/ebiten/v2"
)

type Intro_Scene struct {
	Context          *context.GameContext
	SceneManager     SceneManager
	CurrentStateName string
}

func (scene *Intro_Scene) GetName() string {
	return "Intro"
}

func (scene *Intro_Scene) GetStateName() string {
	return scene.CurrentStateName
}

func NewIntroScene(ctx *context.GameContext, sceneManager SceneManager) *Intro_Scene {
	scene := Intro_Scene{
		Context:      ctx,
		SceneManager: sceneManager,
	}
	return &scene
}

func (scene *Intro_Scene) Update() error {
	return nil
}

func (scene *Intro_Scene) Draw(screen *ebiten.Image) {
}

func (scene *Intro_Scene) GetItemRenderer() *graphics.ItemRenderer {
	return nil
}

var _ Scene = (*Intro_Scene)(nil)
