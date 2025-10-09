package shaders

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type WaterShaderUniforms struct {
	ScreenSize   [2]float32
	WaveShift    float32
	WaveOffset   [2]float32
	WaveEmersion float32
	Scale        float32
	FastPeriod   float64
	SlowPeriod   float64
}

func (wsu *WaterShaderUniforms) ToShaders(dtso *ebiten.DrawTrianglesShaderOptions) {}

var _ ShaderUniforms = (*WaterShaderUniforms)(nil)
