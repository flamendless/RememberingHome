package ui

import (
	"remembering-home/src/assets/shaders"
	"remembering-home/src/context"
	"remembering-home/src/graphics"
	"remembering-home/src/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	input "github.com/quasilyte/ebitengine-input"
)

type LayoutType int

const (
	LayoutVertical LayoutType = iota
	LayoutHorizontal
)

type BannerConfig struct {
	Padding          float64
	HeightMultiplier float64
	FixedWidth       float64
	AutoWidth        bool
}

type TextListConfig struct {
	Layout       LayoutType
	BaseX        float64
	BaseY        float64
	Gap          float64
	Alignment    text.Align
	SecAlignment text.Align
	BannerConfig BannerConfig
}

type TextList struct {
	Config       TextListConfig
	Texts        []*graphics.Text
	CurrentIdx   int
	Enabled      bool
	Visible      bool
	confirmed    bool
	cancelled    bool
	indexChanged bool
	bannerWidth  float64
}

func NewTextList(config TextListConfig) *TextList {
	return &TextList{
		Config:      config,
		Texts:       make([]*graphics.Text, 0),
		CurrentIdx:  0,
		Enabled:     true,
		Visible:     true,
		confirmed:   false,
		cancelled:   false,
		bannerWidth: config.BannerConfig.FixedWidth,
	}
}

func (tl *TextList) AddText(txt *graphics.Text) {
	var x, y float64
	if tl.Config.Layout == LayoutVertical {
		x = tl.Config.BaseX
		y = tl.Config.BaseY
		if len(tl.Texts) > 0 {
			lastText := tl.Texts[len(tl.Texts)-1]
			y = lastText.Y + lastText.DO.LineSpacing
		}
	} else {
		x = tl.Config.BaseX
		y = tl.Config.BaseY
		if len(tl.Texts) > 0 {
			x += tl.Config.Gap * float64(len(tl.Texts))
		}
	}

	txt.SetPos(x, y)
	txt.SetAlign(tl.Config.Alignment, tl.Config.SecAlignment)
	tl.Texts = append(tl.Texts, txt)

	if tl.Config.BannerConfig.AutoWidth {
		tl.calculateBannerWidth()
	}
}

func (tl *TextList) calculateBannerWidth() {
	maxWidth := float64(0)
	for _, txt := range tl.Texts {
		width := float64(len(txt.Txt)) * 20.0
		if width > maxWidth {
			maxWidth = width
		}
	}
	tl.bannerWidth = maxWidth + tl.Config.BannerConfig.Padding*3
}

func (tl *TextList) Update(inputHandler *input.Handler) {
	tl.confirmed = false
	tl.cancelled = false
	tl.indexChanged = false

	if !tl.Enabled || !tl.Visible {
		return
	}

	if len(tl.Texts) == 0 {
		return
	}

	prevIdx := tl.CurrentIdx

	if tl.Config.Layout == LayoutVertical {
		if inputHandler.ActionIsJustReleased(context.ActionMoveUp) {
			tl.CurrentIdx--
		} else if inputHandler.ActionIsJustReleased(context.ActionMoveDown) {
			tl.CurrentIdx++
		}
	} else {
		if inputHandler.ActionIsJustReleased(context.ActionMoveLeft) {
			tl.CurrentIdx--
		} else if inputHandler.ActionIsJustReleased(context.ActionMoveRight) {
			tl.CurrentIdx++
		}
	}

	tl.CurrentIdx = utils.ClampInt(tl.CurrentIdx, 0, len(tl.Texts)-1)

	if prevIdx != tl.CurrentIdx {
		tl.indexChanged = true
	}

	if inputHandler.ActionIsJustReleased(context.ActionEnter) {
		tl.confirmed = true
	}
}

func (tl *TextList) IsConfirmed() bool {
	return tl.confirmed
}

func (tl *TextList) IsCancelled() bool {
	return tl.cancelled
}

func (tl *TextList) SetCancelled(cancelled bool) {
	tl.cancelled = cancelled
}

func (tl *TextList) GetSelectedIndex() int {
	return tl.CurrentIdx
}

func (tl *TextList) GetSelectedText() *graphics.Text {
	if len(tl.Texts) == 0 {
		return nil
	}
	return tl.Texts[tl.CurrentIdx]
}

func (tl *TextList) SetEnabled(enabled bool) {
	tl.Enabled = enabled
}

func (tl *TextList) SetVisible(visible bool) {
	tl.Visible = visible
}

func (tl *TextList) IsVisible() bool {
	return tl.Visible
}

func (tl *TextList) SetCurrentIndex(idx int) {
	tl.CurrentIdx = utils.ClampInt(idx, 0, len(tl.Texts)-1)
}

func (tl *TextList) IndexChanged() bool {
	return tl.indexChanged
}

func (tl *TextList) UpdateBannerPosition(uniforms *shaders.SilentHillRedShaderUniforms) {
	if len(tl.Texts) == 0 {
		return
	}

	curText := tl.Texts[tl.CurrentIdx]
	textH := curText.DO.LineSpacing
	bannerHeight := float64(textH) * tl.Config.BannerConfig.HeightMultiplier

	posX, posY, sizeX, sizeY := graphics.CalculateBannerPosition(
		curText,
		tl.bannerWidth,
		bannerHeight,
	)

	uniforms.BannerPos[0] = posX
	uniforms.BannerPos[1] = posY
	uniforms.BannerSize[0] = sizeX
	uniforms.BannerSize[1] = sizeY
}

func (tl *TextList) ApplyShaderAlpha(uniforms *shaders.SilentHillRedShaderUniforms) {
	if !tl.Visible {
		for _, txt := range tl.Texts {
			utils.DOReplaceAlpha(txt.DO, 0)
		}
		return
	}

	textAlpha := uniforms.GetTextAlpha()
	for _, txt := range tl.Texts {
		utils.DOReplaceAlpha(txt.DO, textAlpha)
	}
}

func (tl *TextList) Draw(canvas *ebiten.Image) {
	if !tl.Visible {
		return
	}

	for _, txt := range tl.Texts {
		txt.Draw(canvas)
	}
}

func (tl *TextList) GetBannerWidth() float64 {
	return tl.bannerWidth
}

func (tl *TextList) SetBannerWidth(width float64) {
	tl.bannerWidth = width
}
