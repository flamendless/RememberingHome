package game

import (
	"fmt"
	"image"
	"nowhere-home/src/assets"
	"nowhere-home/src/conf"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Game_State struct {
	Loader          *resource.Loader
	SceneManager    *Scene_Manager
	InputSystem     *input.System
	InputHandler    *input.Handler
	InputHandlerDev *input.Handler
}

func NewGame(
	loader *resource.Loader,
	sceneManager *Scene_Manager,
	inputSystem *input.System,
) *Game_State {
	iconImg := loader.LoadImage(assets.ImageWindowIcon)
	icons := []image.Image{iconImg.Data}

	title := fmt.Sprintf("%s v%s", conf.GAME_TITLE, conf.GAME_VERSION)
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowIcon(icons)
	ebiten.SetWindowSize(conf.WINDOW_W, conf.WINDOW_H)
	ebiten.SetFullscreen(conf.FULLSCREEN)

	gs := &Game_State{
		Loader:       loader,
		SceneManager: sceneManager,
		InputSystem:  inputSystem,
		InputHandler: NewInputHandler(inputSystem),
	}

	if conf.DEV {
		gs.InputHandlerDev = NewInputHandlerDev(inputSystem)
	}

	return gs
}

func (g *Game_State) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return conf.GAME_W, conf.GAME_H
}

func (g *Game_State) Update() error {
	g.InputSystem.Update()

	if conf.DEV {
		FixWSLWindow()
		UpdateDebugInput(g)
		UpdateDebugOverlay(g)
	}

	if !ebiten.IsFocused() {
		return nil
	}

	return g.SceneManager.Update()
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
