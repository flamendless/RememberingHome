package game

import (
	"fmt"
	"image"
	"remembering-home/src/assets"
	"remembering-home/src/assets/shaders"
	"remembering-home/src/conf"
	"remembering-home/src/context"
	"remembering-home/src/debug"
	"remembering-home/src/enums"
	"remembering-home/src/graphics"
	"remembering-home/src/scenes"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Game_State struct {
	Context          *context.GameContext
	SceneManager     *scenes.Scene_Manager
	finalShaderLayer *graphics.Layer
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
	ctx := context.NewGameContext(loader, inputSystem, inputHandler, settings)

	gs := &Game_State{
		Context:      ctx,
		SceneManager: sceneManager,
	}

	if conf.DEV {
		debug.SetSceneNavigator(gs)
	}

	gs.finalShaderLayer = graphics.NewLayerWithShader(
		"quality",
		conf.GAME_W,
		conf.GAME_H,
		loader.LoadShader(assets.ShaderGraphicsQuality).Data,
	)
	gs.finalShaderLayer.Disabled = false
	gs.finalShaderLayer.Uniforms = shaders.NewGraphicsQualityUniform(settings)

	if conf.DEV {
		debug.AddDebugShader(gs.finalShaderLayer.Uniforms)
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
		var currentScene scenes.Scene
		var itemRenderer *graphics.ItemRenderer
		if currentScene = g.SceneManager.GetCurrentScene(); currentScene != nil {
			sceneName = currentScene.GetName()
			sceneState = currentScene.GetStateName()
			itemRenderer = currentScene.GetItemRenderer()
		} else {
			sceneName = "None"
			sceneState = "None"
			itemRenderer = nil
		}

		if g.Context.InputHandler.ActionIsJustPressed(context.ActionZoomIn) {
			currentZoom := debug.GetZoomLevel()
			newZoom := currentZoom + 0.1
			if newZoom <= 3.0 {
				debug.SetZoomLevel(newZoom)
			}
		} else if g.Context.InputHandler.ActionIsJustPressed(context.ActionZoomOut) {
			currentZoom := debug.GetZoomLevel()
			newZoom := currentZoom - 0.1
			if newZoom >= 0.1 {
				debug.SetZoomLevel(newZoom)
			}
		} else if g.Context.InputHandler.ActionIsJustPressed(context.ActionZoomReset) {
			debug.SetZoomLevel(1.0)
		}

		if err := debug.UpdateDebugUI(g.Context, sceneName, sceneState, itemRenderer); err != nil {
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

	// Draw scene to the layer canvas first
	g.finalShaderLayer.Canvas.Clear()
	g.SceneManager.Draw(g.finalShaderLayer.Canvas)

	if conf.DEV {
		zoomLevel := debug.GetZoomLevel()
		g.finalShaderLayer.ScaleX = zoomLevel
		g.finalShaderLayer.ScaleY = zoomLevel

		gameW, gameH := conf.GAME_W, conf.GAME_H
		g.finalShaderLayer.X = (float64(gameW) - float64(gameW)*zoomLevel) / 2
		g.finalShaderLayer.Y = (float64(gameH) - float64(gameH)*zoomLevel) / 2
		g.finalShaderLayer.ApplyTransformation()
	}

	// Update and apply graphics quality shader
	uniform := graphics.MustCastUniform[*shaders.GraphicsQualityUniforms](g.finalShaderLayer)
	uniform.Update()
	uniform.ToShadersDRSO(g.finalShaderLayer.DRSO)
	g.finalShaderLayer.RenderWithShader(screen)

	if conf.DEV {
		debug.DrawDebugOverlay(screen)
	}
}

func (g *Game_State) NavigateTo(scene enums.Scene) {
	switch scene {
	case enums.SceneDummy:
		g.SceneManager.GoTo(scenes.NewDummyScene(g.Context, g.SceneManager))
	case enums.SceneSplash:
		g.SceneManager.GoTo(scenes.NewSplashScene(g.Context, g.SceneManager))
	case enums.SceneMainMenu:
		g.SceneManager.GoTo(scenes.NewMainMenuScene(g.Context, g.SceneManager))
	case enums.SceneStorageRoom:
		g.SceneManager.GoTo(scenes.NewStorageRoomScene(g.Context, g.SceneManager))
	default:
		panic(fmt.Sprintf("Unknown scene: %s", scene))
	}
}
