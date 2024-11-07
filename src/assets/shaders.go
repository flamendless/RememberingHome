package assets

import (
	"remembering-home/src/logger"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"go.uber.org/zap"
)

const (
	ShaderNone resource.ShaderID = iota
	ShaderTest
	ShaderColorize
	ShaderWater
	ShaderMenuText
)

func SetShaderResources(loader *resource.Loader) {
	logger.Log().Info("Setting shader resources...")
	shaderResources := map[resource.ShaderID]resource.ShaderInfo{
		// ShaderTest:     {Path: "shaders/test.kage"},
		ShaderColorize: {Path: "shaders/colorize.kage"},
		ShaderWater:    {Path: "shaders/water.kage"},
		ShaderMenuText: {Path: "shaders/menutext.kage"},
	}
	for id, res := range shaderResources {
		logger.Log().Info("Loading shader", zap.String("path", res.Path))
		loader.ShaderRegistry.Set(id, res)
		loader.LoadShader(id)
	}
}

type ShaderUniforms interface{
	ToShaders(DTSO *ebiten.DrawTrianglesShaderOptions)
}

type WaterShaderUniforms struct {
	ScreenSize   [2]float32
	WaveShift    float32
	WaveOffset   [2]float32
	WaveEmersion float32
	Scale        float32
	FastPeriod   float64
	SlowPeriod   float64
}

func (wsu *WaterShaderUniforms) ToShaders(DTSO *ebiten.DrawTrianglesShaderOptions) {}

type MenuTextShaderUniforms struct {
	Time              float64
	Pos               [2]float32
	Size              [2]float32
	StartingAmplitude float32 //0.0 - 0.5
	StartingFreq      float32
	Shift             float32 //-1.0 - 1.0
	WhiteCutoff       float32 //0.0 - 1.0
	Velocity          [2]float32
	Color             [4]float32
}

func (mtsu *MenuTextShaderUniforms) ToShaders(DTSO *ebiten.DrawTrianglesShaderOptions) {
	DTSO.Uniforms = map[string]any{
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

var _ ShaderUniforms = (*WaterShaderUniforms)(nil)
var _ ShaderUniforms = (*MenuTextShaderUniforms)(nil)
