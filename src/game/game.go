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

var WSLTricked bool

type Game_State struct {
	Loader       *resource.Loader
	SceneManager *Scene_Manager
}

func NewGame(
	loader *resource.Loader,
	sceneManager *Scene_Manager,
) *Game_State {
	iconImg := loader.LoadImage(assets.ImageWindowIcon)
	icons := []image.Image{iconImg.Data}

	ebiten.SetWindowIcon(icons)
	ebiten.SetWindowSize(conf.WINDOW_W, conf.WINDOW_H)

	title := fmt.Sprintf("%s v%s", conf.GAME_TITLE, conf.GAME_VERSION)
	ebiten.SetWindowTitle(title)

	return &Game_State{
		Loader:       loader,
		SceneManager: sceneManager,
	}
}

func (g *Game_State) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return conf.GAME_W, conf.GAME_H
}

func (g *Game_State) Update() error {
	g.SceneManager.Update()

	if conf.DEV {
		//trick for WSL2
		if !WSLTricked && !ebiten.IsFocused() {
			ebiten.MinimizeWindow()
			ebiten.MaximizeWindow()
			ebiten.RestoreWindow()
			WSLTricked = true
		}

		overlays.UpdateDebug(g.SceneManager.current.GetName())
	}

	return nil
}

func (g *Game_State) Draw(screen *ebiten.Image) {
	g.SceneManager.Draw(screen)

	if conf.DEV {
		overlays.DrawDebug(screen)
	}
}
