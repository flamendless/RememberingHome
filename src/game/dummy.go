package game

import (
	"nowhere-home/src/assets"
	"nowhere-home/src/conf"

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

func (scene *Dummy_Scene) Update() error {
	return nil
}

func (scene *Dummy_Scene) Draw(screen *ebiten.Image) {
	logo := scene.GameState.Loader.LoadImage(assets.ImageFlamLogo)
	op := &ebiten.DrawImageOptions{}
	size := logo.Data.Bounds().Size()
	op.GeoM.Rotate(0.5)
	op.GeoM.Translate(
		float64(conf.GAME_W/2-size.X/2),
		float64(conf.GAME_H/2-size.Y/2),
	)
	screen.DrawImage(logo.Data, op)
}
