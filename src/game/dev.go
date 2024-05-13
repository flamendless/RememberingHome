package game

import (
	"fmt"
	"nowhere-home/src/conf"
	"nowhere-home/src/overlays"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	WSLTricked bool
	Show       bool
	debugTexts [6]string
)

const (
	VERSION int = iota
	FPS
	TPS
	OVERLAY_FADE_ALPHA
	SCENE_NAME
	SCENE_STATE
)

func init() {
	Show = conf.DEV
	if conf.DEV {
		debugTexts[VERSION] = fmt.Sprintf("VERSION: %s\n", conf.GAME_VERSION)
	}
}

func FixWSLWindow() {
	if !WSLTricked && !ebiten.IsFocused() {
		ebiten.MinimizeWindow()
		ebiten.MaximizeWindow()
		ebiten.RestoreWindow()
		WSLTricked = true
	}
}

func UpdateDebugInput(g *Game_State) {
	if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		Show = !Show
	} else if inpututil.IsKeyJustReleased(ebiten.Key1) {
		g.SceneManager.GoTo(NewDummyScene(g))
	} else if inpututil.IsKeyJustReleased(ebiten.Key2) {
		g.SceneManager.GoTo(NewSplashScene(g))
	} else if inpututil.IsKeyJustReleased(ebiten.Key3) {
		g.SceneManager.GoTo(NewMainMenuScene(g))
	}
}

func UpdateDebugOverlay(g *Game_State) {
	if !Show {
		return
	}
	debugTexts[FPS] = fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS())
	debugTexts[TPS] = fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS())
	debugTexts[OVERLAY_FADE_ALPHA] = fmt.Sprintf("Overlay: %.2f", overlays.FadeAlpha)

	sceneName := g.SceneManager.current.GetName()
	if g.SceneManager.next != nil {
		nextSceneName := g.SceneManager.next.GetName()
		debugTexts[SCENE_NAME] = fmt.Sprintf("Scene: %s, Next: %s", sceneName, nextSceneName)
	} else {
		debugTexts[SCENE_NAME] = fmt.Sprintf("Scene: %s", sceneName)
	}

	debugTexts[SCENE_STATE] = fmt.Sprintf(
		"Scene State: %s",
		g.SceneManager.current.GetStateName(),
	)
}

func DrawDebugOverlay(screen *ebiten.Image) {
	if !Show {
		return
	}
	for i := 0; i < len(debugTexts); i++ {
		ebitenutil.DebugPrintAt(screen, debugTexts[i], 0, 16*i)
	}
}
