package game

import (
	"fmt"
	"image"
	"remembering-home/src/assets"
	"remembering-home/src/conf"
	"remembering-home/src/context"
	"remembering-home/src/debug"
	"remembering-home/src/scenes"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Game_State struct {
	Context      *context.GameContext
	SceneManager *scenes.Scene_Manager
}

func NewGame(
	loader *resource.Loader,
	sceneManager *scenes.Scene_Manager,
	inputSystem *input.System,
	settings *conf.Settings,
) *Game_State {
	iconImg := loader.LoadImage(assets.ImageWindowIcon)
	icons := []image.Image{iconImg.Data}

	title := fmt.Sprintf("%s v%s", conf.GAME_TITLE, conf.GAME_VERSION)
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowIcon(icons)
	ebiten.SetWindowSize(conf.WINDOW_W, conf.WINDOW_H)
	ebiten.SetFullscreen(settings.Window == conf.WindowModeFullscreen)

	inputHandler := NewInputHandler(inputSystem)
	var inputHandlerDev *input.Handler
	if conf.DEV {
		inputHandlerDev = NewInputHandlerDev(inputSystem)
	}

	ctx := context.NewGameContext(loader, inputSystem, inputHandler, inputHandlerDev, settings)

	gs := &Game_State{
		Context:      ctx,
		SceneManager: sceneManager,
	}

	return gs
}

func (g *Game_State) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return conf.GAME_W, conf.GAME_H
}

func (g *Game_State) Update() error {
	g.Context.InputSystem.Update()

	if conf.DEV {
		debug.FixWSLWindow()

		var sceneName, sceneState string
		if currentScene := g.SceneManager.GetCurrentScene(); currentScene != nil {
			sceneName = currentScene.GetName()
			sceneState = currentScene.GetStateName()
		} else {
			sceneName = "None"
			sceneState = "None"
		}

		if err := debug.UpdateDebugUI(g.Context, sceneName, sceneState); err != nil {
			return err
		}
	}

	if !ebiten.IsFocused() {
		return nil
	}

	return g.SceneManager.Update()
}

func (g *Game_State) Draw(screen *ebiten.Image) {
	if !ebiten.IsFocused() && runtime.GOARCH != "wasm" {
		return
	}

	g.SceneManager.Draw(screen)

	if conf.DEV {
		debug.DrawDebugOverlay(screen)
	}
}
