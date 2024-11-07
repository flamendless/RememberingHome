package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
)

type Text struct {
	Face   *text.GoXFace
	DO     *text.DrawOptions
	Txt    string
	Static bool
	X, Y   float64
	Show   bool
}

func NewText(face *font.Face, str string, static bool) *Text {
	xface := text.NewGoXFace(*face)

	op := &text.DrawOptions{}
	op.LineSpacing = xface.Metrics().HLineGap + xface.Metrics().HAscent + xface.Metrics().HDescent

	return &Text{
		Face: xface,
		DO:   op,
		Txt:  str,
		Show: true,
	}
}

func (txt *Text) GetAlpha() float32 {
	return txt.DO.ColorScale.A()
}

func (txt *Text) SetShow(show bool) {
	txt.Show = show
}

func (txt *Text) SetPos(x, y float64) {
	txt.X = x
	txt.Y = y
	txt.DO.GeoM.Translate(x, y)
}

func (txt *Text) SetAlign(pa, sa text.Align) {
	txt.DO.PrimaryAlign = pa
	txt.DO.SecondaryAlign = sa
}

func (txt *Text) Draw(screen *ebiten.Image) {
	if !txt.Show {
		return
	}

	if !txt.Static {
		txt.DO.GeoM.Reset()
		txt.SetPos(txt.X, txt.Y)
	}

	text.Draw(screen, txt.Txt, txt.Face, txt.DO)
}
