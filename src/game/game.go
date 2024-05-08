package game

import (
	"fmt"
	"image"
	"nowhere-home/src/assets"
	"nowhere-home/src/conf"
	"nowhere-home/src/overlays"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Game struct {
	Loader       *resource.Loader
}

func NewGame(loader *resource.Loader) *Game {
	iconImg := loader.LoadImage(assets.ImageWindowIcon)
	icons := []image.Image{iconImg.Data}

	ebiten.SetWindowIcon(icons)
	ebiten.SetWindowSize(conf.WINDOW_W, conf.WINDOW_H)

	title := fmt.Sprintf("%s v%s", conf.GAME_TITLE, conf.GAME_VERSION)
	ebiten.SetWindowTitle(title)
	return &Game{
		Loader: loader,
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return conf.GAME_W, conf.GAME_H
}

func (g *Game) Update() error {
	overlays.UpdateFade()

	if conf.DEV {
		overlays.UpdateDebug()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	logo := g.Loader.LoadImage(assets.ImageFlamendlessLogo)
	op := &ebiten.DrawImageOptions{}
	size := logo.Data.Bounds().Size()
	op.GeoM.Translate(
		float64(conf.GAME_W/2-size.X/2),
		float64(conf.GAME_H/2-size.Y/2),
	)
	screen.DrawImage(logo.Data, op)
	overlays.DrawFade(screen)
	if conf.DEV {
		overlays.DrawDebug(screen)
	}
}
