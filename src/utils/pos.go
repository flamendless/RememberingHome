package utils

import "github.com/hajimehoshi/ebiten/v2"

func DIOSetPos(dio *ebiten.DrawImageOptions, x float64, y float64) {
	dio.GeoM.Reset()
	dio.GeoM.Translate(x, y)
}

func DRSOSetPos(drso *ebiten.DrawRectShaderOptions, x float64, y float64) {
	drso.GeoM.Reset()
	drso.GeoM.Translate(x, y)
}

func DRSOSetPosX(drso *ebiten.DrawRectShaderOptions, x float64) {
	drso.GeoM.SetElement(0, 2, x)
}

func DRSOSetPosY(drso *ebiten.DrawRectShaderOptions, y float64) {
	drso.GeoM.SetElement(1, 2, y)
}
