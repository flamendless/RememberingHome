package shaders

import (
	"remembering-home/src/conf"

	"github.com/hajimehoshi/ebiten/v2"
)

type GraphicsQualityUniforms struct {
	initial    *GraphicsQualityUniforms
	Settings   *conf.Settings
	Quality    float64
	Resolution [2]float64
}

func NewGraphicsQualityUniform(settings *conf.Settings) *GraphicsQualityUniforms {
	uniforms := &GraphicsQualityUniforms{
		Settings:   settings,
		Quality:    settings.Quality.ToShaderValue(),
		Resolution: [2]float64{conf.GAME_W, conf.GAME_H},
	}
	initialCopy := *uniforms
	uniforms.initial = &initialCopy
	return uniforms
}

func (gqu *GraphicsQualityUniforms) ToShaders(dtso *ebiten.DrawTrianglesShaderOptions) {
	dtso.Uniforms = map[string]any{
		"Quality":    float32(gqu.Quality),
		"Resolution": [2]float32{float32(gqu.Resolution[0]), float32(gqu.Resolution[1])},
	}
}

func (gqu *GraphicsQualityUniforms) ToShadersDRSO(drso *ebiten.DrawRectShaderOptions) {
	drso.Uniforms = map[string]any{
		"Quality":    float32(gqu.Quality),
		"Resolution": [2]float32{float32(gqu.Resolution[0]), float32(gqu.Resolution[1])},
	}
}

func (gqu *GraphicsQualityUniforms) Update() {
	gqu.Quality = gqu.Settings.Quality.ToShaderValue()
}

func (gqu *GraphicsQualityUniforms) ResetToInitial() {
	*gqu = *gqu.initial
}

var _ ShaderUniforms = (*GraphicsQualityUniforms)(nil)
