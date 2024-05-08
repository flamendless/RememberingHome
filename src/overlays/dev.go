package overlays

import (
	"fmt"
	"nowhere-home/src/conf"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	strGameVer     string
	strFPS         string
	strTPS         string
	strScene       string
	strFadeOverlay string
)

func init() {
	strGameVer = fmt.Sprintf("VERSION: %s\n", conf.GAME_VERSION)
}

func UpdateDebug(sceneName string) {
	strFPS = fmt.Sprintf("FPS: %.2f\n", ebiten.ActualFPS())
	strTPS = fmt.Sprintf("TPS: %.2f\n", ebiten.ActualTPS())
	strScene = fmt.Sprintf("Scene: %s\n", sceneName)
	strFadeOverlay = fmt.Sprintf("Overlay: %.2f\n", FadeAlpha)
}

func DrawDebug(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, strGameVer)
	ebitenutil.DebugPrintAt(screen, strFPS, 0, 16)
	ebitenutil.DebugPrintAt(screen, strTPS, 0, 32)
	ebitenutil.DebugPrintAt(screen, strScene, 0, 48)
	ebitenutil.DebugPrintAt(screen, strFadeOverlay, 0, 64)
}
