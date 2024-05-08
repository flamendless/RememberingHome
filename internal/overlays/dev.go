package overlays

import (
	"fmt"
	"nowhere-home/internal/conf"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	strGameVer string
	strFPS     string
	strTPS     string
)

func init() {
	strGameVer = fmt.Sprintf("VERSION: %s\n", conf.GAME_VERSION)
}

func UpdateDebug() {
	strFPS = fmt.Sprintf("FPS: %f\n", ebiten.ActualFPS())
	strTPS = fmt.Sprintf("TPS: %f\n", ebiten.ActualTPS())
}

func DrawDebug(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, strGameVer)
	ebitenutil.DebugPrintAt(screen, strFPS, 0, 16)
	ebitenutil.DebugPrintAt(screen, strTPS, 0, 32)
}
