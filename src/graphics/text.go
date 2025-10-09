package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
)

type Text struct {
	Face   *text.GoXFace
	DO     *text.DrawOptions
	Txt    string
	X, Y   float64
	Static bool
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

func CalculateBannerPosition(textObj *Text, bannerWidth, bannerHeight float64) (posX, posY, sizeX, sizeY float64) {
	textH := textObj.DO.LineSpacing

	var bannerY float64
	if textObj.DO.SecondaryAlign == text.AlignCenter {
		bannerY = float64(textObj.Y) - bannerHeight/2
	} else {
		bannerY = float64(textObj.Y) - bannerHeight/2 + float64(textH)/2
	}

	var bannerX float64
	switch textObj.DO.PrimaryAlign {
	case text.AlignCenter:
		bannerX = float64(textObj.X) - bannerWidth/2
	case text.AlignEnd:
		textWidth := float64(len(textObj.Txt)) * 20.0
		textCenterX := float64(textObj.X) - textWidth/2
		bannerX = textCenterX - bannerWidth/2
	default:
		textWidth := float64(len(textObj.Txt)) * 20.0
		textCenterX := float64(textObj.X) + textWidth/2
		bannerX = textCenterX - bannerWidth/2
	}

	return bannerX, bannerY, bannerWidth, bannerHeight
}
