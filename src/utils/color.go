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
