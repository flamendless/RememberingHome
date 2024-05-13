package game

import (
	"fmt"
	"image"
	"nowhere-home/src/assets"
	"nowhere-home/src/conf"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
)

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

	title := fmt.Sprintf("%s v%s", conf.GAME_TITLE, conf.GAME_VERSION)
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowIcon(icons)
	ebiten.SetWindowSize(conf.WINDOW_W, conf.WINDOW_H)
	ebiten.SetFullscreen(conf.FULLSCREEN)

	return &Game_State{
		Loader:       loader,
		SceneManager: sceneManager,
	}
}

func (g *Game_State) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return conf.GAME_W, conf.GAME_H
}

func (g *Game_State) Update() error {
	if conf.DEV {
		FixWSLWindow()
		UpdateDebugInput(g)
		UpdateDebugOverlay(g)
	}

	if !ebiten.IsFocused() {
		return nil
	}

	g.SceneManager.Update()

	return nil
}

func (g *Game_State) Draw(screen *ebiten.Image) {
	if !ebiten.IsFocused() {
		return
	}

	g.SceneManager.Draw(screen)

	if conf.DEV {
		DrawDebugOverlay(screen)
	}
}
