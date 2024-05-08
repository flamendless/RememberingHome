package scenes

import (
	// "nowhere-home/src/assets"
	// "nowhere-home/src/conf"

	// "fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SplashScene struct {}

func (scene *SplashScene) Update(state *GameState) error {
	// fmt.Println(state)
	return nil
}

func (scene *SplashScene) Draw(screen *ebiten.Image) {
	// logo := scene.game.Loader.LoadImage(assets.ImageFlamendlessLogo)
	// op := &ebiten.DrawImageOptions{}
	// size := logo.Data.Bounds().Size()
	// op.GeoM.Translate(
	// 	float64(conf.GAME_W/2-size.X/2),
	// 	float64(conf.GAME_H/2-size.Y/2),
	// )
	// screen.DrawImage(logo.Data, op)
	ebitenutil.DebugPrintAt(screen, "SCENE", 160, 160)
}
