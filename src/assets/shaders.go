package assets

import (
	"nowhere-home/src/logger"

	resource "github.com/quasilyte/ebitengine-resource"
	"go.uber.org/zap"
)

const (
	ShaderNone resource.ShaderID = iota
	ShaderTest
	ShaderColorize
	ShaderWater
	ShaderFog
)

func SetShaderResources(loader *resource.Loader) {
	logger.Log().Info("Setting shader resources...")
	shaderResources := map[resource.ShaderID]resource.ShaderInfo{
		// ShaderTest:     {Path: "shaders/test.kage"},
		ShaderColorize: {Path: "shaders/colorize.kage"},
		ShaderWater:    {Path: "shaders/water.kage"},
		ShaderFog:      {Path: "shaders/fog.kage"},
	}
	for id, res := range shaderResources {
		logger.Log().Info("Loading shader", zap.String("path", res.Path))
		loader.ShaderRegistry.Set(id, res)
		loader.LoadShader(id)
	}
}

type ShaderUniforms interface{}

type WaterShader struct {
	ScreenSize   [2]float32
	WaveShift    float32
	WaveOffset   [2]float32
	WaveEmersion float32
	Scale        float32
	FastPeriod   float64
	SlowPeriod   float64
}

type FogShader struct {
	Time              float64
	StartingAmplitude float32
	StartingFreq      float32
	Shift             float32
	WhiteCutoff       float32
	Velocity          [2]float32
	FogColor          [4]float32
}
