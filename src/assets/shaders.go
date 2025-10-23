package assets

import (
	"remembering-home/src/logger"

	resource "github.com/quasilyte/ebitengine-resource"
	"go.uber.org/zap"
)

const (
	ShaderNone resource.ShaderID = iota
	ShaderTest
	ShaderColorize
	ShaderWater
	ShaderTextRedBG
	ShaderSilentHillRed
	ShaderGraphicsQuality
)

func setShaderResources(loader *resource.Loader) {
	logger.Log().Info("Setting shader resources...")
	shaderResources := map[resource.ShaderID]resource.ShaderInfo{
		// ShaderTest:     {Path: "shaders/test.kage"},
		ShaderColorize:        {Path: "shaders/colorize.kage"},
		ShaderWater:           {Path: "shaders/water.kage"},
		ShaderTextRedBG:       {Path: "shaders/text_red_bg.kage"},
		ShaderSilentHillRed:   {Path: "shaders/silent_hill_red.kage"},
		ShaderGraphicsQuality: {Path: "shaders/graphics_quality.kage"},
	}

	for id, res := range shaderResources {
		logger.Log().Info("Loading shader", zap.String("path", res.Path))
		loader.ShaderRegistry.Set(id, res)
		loader.LoadShader(id)
	}
}
