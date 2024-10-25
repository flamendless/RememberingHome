package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func DIOReplaceAlpha(DIO *ebiten.DrawImageOptions, a float32) {
	DIO.ColorScale = ebiten.ColorScale{}
	DIO.ColorScale.ScaleAlpha(a)
}

func DOReplaceAlpha(DO *text.DrawOptions, a float32) {
	DO.ColorScale = ebiten.ColorScale{}
	DO.ColorScale.ScaleAlpha(a)
}

func SetColor(DO *text.DrawOptions, r, g, b, a float32) {
	DO.ColorScale.SetR(r)
	DO.ColorScale.SetG(g)
	DO.ColorScale.SetB(b)
	DO.ColorScale.SetA(a)
}

func DRSOSetColor(DRSO *ebiten.DrawRectShaderOptions, r, g, b, a float32) {
	DRSO.ColorScale.SetR(r)
	DRSO.ColorScale.SetG(g)
	DRSO.ColorScale.SetB(b)
	DRSO.ColorScale.SetA(a)
}
