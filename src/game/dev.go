package game

import (
	"fmt"
	"image/color"
	"nowhere-home/src/conf"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	WSLTricked bool
	ShowTexts       bool
	ShowLines  bool
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
	ShowTexts = conf.DEV
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
	if g.InputHandlerDev.ActionIsJustReleased(DevToggleTexts) {
		ShowTexts = !ShowTexts
	} else if g.InputHandlerDev.ActionIsJustReleased(DevToggleLines) {
		ShowLines = !ShowLines
	} else if g.InputHandlerDev.ActionIsJustReleased(DevGoToDummy) {
		g.SceneManager.GoTo(NewDummyScene(g))
	} else if g.InputHandlerDev.ActionIsJustReleased(DevGoToSplash) {
		g.SceneManager.GoTo(NewSplashScene(g))
	} else if g.InputHandlerDev.ActionIsJustReleased(DevGoToMainMenu) {
		g.SceneManager.GoTo(NewMainMenuScene(g))
	}
}

func UpdateDebugOverlay(g *Game_State) {
	if !ShowTexts {
		return
	}
	debugTexts[FPS] = fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS())
	debugTexts[TPS] = fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS())
	debugTexts[OVERLAY_FADE_ALPHA] = fmt.Sprintf("Fade: Alpha = %.2f, Dir = %d", g.SceneManager.fader.Alpha, g.SceneManager.fader.Dir)

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
	if ShowTexts {
		for i := 0; i < len(debugTexts); i++ {
			ebitenutil.DebugPrintAt(screen, debugTexts[i], 0, 16*i)
		}
	}

	if ShowLines {
		clr := color.White
		vector.StrokeLine(screen, conf.GAME_W/2, 0, conf.GAME_W/2, conf.GAME_H, 1, clr, false)
		vector.StrokeLine(screen, 0, conf.GAME_H/2, conf.GAME_W, conf.GAME_H/2, 1, clr, false)
	}
}
