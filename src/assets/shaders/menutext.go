package shaders

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type MenuTextShaderUniforms struct {
	Time              float64
	Pos               [2]float32
	Size              [2]float32
	StartingAmplitude float32 // 0.0 - 0.5
	StartingFreq      float32
	Shift             float32 //-1.0 - 1.0
	WhiteCutoff       float32 // 0.0 - 1.0
	Velocity          [2]float32
	Color             [4]float32
}

func (mtsu *MenuTextShaderUniforms) ToShaders(dtso *ebiten.DrawTrianglesShaderOptions) {
	dtso.Uniforms = map[string]any{
		"Time":              mtsu.Time,
		"Pos":               mtsu.Pos,
		"Size":              mtsu.Size,
		"StartingAmplitude": mtsu.StartingAmplitude,
		"StartingFreq":      mtsu.StartingFreq,
		"Shift":             mtsu.Shift,
		"WhiteCutoff":       mtsu.WhiteCutoff,
		"Velocity":          mtsu.Velocity,
		"Color":             mtsu.Color,
	}
}

var _ ShaderUniforms = (*MenuTextShaderUniforms)(nil)
