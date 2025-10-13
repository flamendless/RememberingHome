package ui

import (
	"remembering-home/src/assets/shaders"
	"remembering-home/src/graphics"
	"remembering-home/src/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type BannerText struct {
	Text         *graphics.Text
	BannerConfig BannerConfig
	Visible      bool
	bannerWidth  float64
}

func NewBannerText(text *graphics.Text, config BannerConfig) *BannerText {
	bt := &BannerText{
		Text:         text,
		BannerConfig: config,
		Visible:      true,
	}

	if config.AutoWidth {
		bt.calculateBannerWidth()
	} else {
		bt.bannerWidth = config.FixedWidth
	}

	return bt
}

func (bt *BannerText) calculateBannerWidth() {
	width := float64(len(bt.Text.Txt)) * 15.0
	bt.bannerWidth = width + bt.BannerConfig.Padding*3
}

func (bt *BannerText) SetVisible(visible bool) {
	bt.Visible = visible
}

func (bt *BannerText) IsVisible() bool {
	return bt.Visible
}

func (bt *BannerText) UpdateBannerPosition(uniforms *shaders.SilentHillRedShaderUniforms) {
	if !bt.Visible {
		return
	}

	bt.UpdateBannerPositionForce(uniforms)
}

func (bt *BannerText) UpdateBannerPositionForce(uniforms *shaders.SilentHillRedShaderUniforms) {
	textH := bt.Text.DO.LineSpacing
	bannerHeight := float64(textH) * bt.BannerConfig.HeightMultiplier
	posX, posY, sizeX, sizeY := graphics.CalculateBannerPosition(
		bt.Text,
		bt.bannerWidth,
		bannerHeight,
	)
	uniforms.BannerPos[0] = posX
	uniforms.BannerPos[1] = posY
	uniforms.BannerSize[0] = sizeX
	uniforms.BannerSize[1] = sizeY
}

func (bt *BannerText) ApplyShaderAlpha(uniforms *shaders.SilentHillRedShaderUniforms) {
	if !bt.Visible {
		utils.DOReplaceAlpha(bt.Text.DO, 0)
		return
	}

	textAlpha := uniforms.GetTextAlpha()
	utils.DOReplaceAlpha(bt.Text.DO, textAlpha)
}

func (bt *BannerText) Draw(canvas *ebiten.Image) {
	if !bt.Visible {
		return
	}

	bt.Text.Draw(canvas)
}

func (bt *BannerText) GetBannerWidth() float64 {
	return bt.bannerWidth
}

func (bt *BannerText) SetBannerWidth(width float64) {
	bt.bannerWidth = width
}

func (bt *BannerText) UpdateText(newText string) {
	bt.Text.Txt = newText
	if bt.BannerConfig.AutoWidth {
		bt.calculateBannerWidth()
	}
}
