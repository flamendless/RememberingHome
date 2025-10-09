package shaders

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type SilentHillRedShaderUniforms struct {
	ID             string
	Time           float64
	BannerPos      [2]float64
	BannerSize     [2]float64
	BaseRedColor   [4]float64
	GlowIntensity  float64
	MetallicShine  float64
	EdgeDarkness   float64
	TextGlowRadius float64
	NoiseScale     float64
	NoiseIntensity float64
}

func (shrsu *SilentHillRedShaderUniforms) ToShaders(dtso *ebiten.DrawTrianglesShaderOptions) {
	dtso.Uniforms = map[string]any{
		"ID":             "SilentHillRed",
		"Time":           shrsu.Time,
		"BannerPos":      shrsu.BannerPos,
		"BannerSize":     shrsu.BannerSize,
		"BaseRedColor":   shrsu.BaseRedColor,
		"GlowIntensity":  shrsu.GlowIntensity,
		"MetallicShine":  shrsu.MetallicShine,
		"EdgeDarkness":   shrsu.EdgeDarkness,
		"TextGlowRadius": shrsu.TextGlowRadius,
		"NoiseScale":     shrsu.NoiseScale,
		"NoiseIntensity": shrsu.NoiseIntensity,
	}
}

var _ ShaderUniforms = (*SilentHillRedShaderUniforms)(nil)
