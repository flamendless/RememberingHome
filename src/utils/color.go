package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func DIOReplaceAlpha(dio *ebiten.DrawImageOptions, a float32) {
	dio.ColorScale = ebiten.ColorScale{}
	dio.ColorScale.ScaleAlpha(a)
}

func DOReplaceAlpha(do *text.DrawOptions, a float32) {
	do.ColorScale = ebiten.ColorScale{}
	do.ColorScale.ScaleAlpha(a)
}

func SetColor(do *text.DrawOptions, r, g, b, a float32) {
	do.ColorScale.SetR(r)
	do.ColorScale.SetG(g)
	do.ColorScale.SetB(b)
	do.ColorScale.SetA(a)
}

func DRSOSetColor(drso *ebiten.DrawRectShaderOptions, r, g, b, a float32) {
	drso.ColorScale.SetR(r)
	drso.ColorScale.SetG(g)
	drso.ColorScale.SetB(b)
	drso.ColorScale.SetA(a)
}
