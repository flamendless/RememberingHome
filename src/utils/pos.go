package utils

import "github.com/hajimehoshi/ebiten/v2"

func DIOSetPos(DIO *ebiten.DrawImageOptions, x float64, y float64) {
	DIO.GeoM.Reset()
	DIO.GeoM.Translate(x, y)
}

func DRSOSetPos(DRSO *ebiten.DrawRectShaderOptions, x float64, y float64) {
	DRSO.GeoM.Reset()
	DRSO.GeoM.Translate(x, y)
}

func DRSOSetPosX(DRSO *ebiten.DrawRectShaderOptions, x float64) {
	DRSO.GeoM.SetElement(0, 2, x)
}

func DRSOSetPosY(DRSO *ebiten.DrawRectShaderOptions, y float64) {
	DRSO.GeoM.SetElement(1, 2, y)
}
