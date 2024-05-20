package assets

import (
	"nowhere-home/src/logger"

	resource "github.com/quasilyte/ebitengine-resource"
	"go.uber.org/zap"
)

const (
	ShaderNone resource.ShaderID = iota
	ShaderColorize
)

func SetShaderResources(loader *resource.Loader) {
	logger.Log().Info("Setting shader resources...")
	shaderResources := map[resource.ShaderID]resource.ShaderInfo{
		ShaderColorize: {Path: "shaders/colorize.kage"},
	}
	for id, res := range shaderResources {
		logger.Log().Info("Loading shader", zap.String("path", res.Path))
		loader.ShaderRegistry.Set(id, res)
		loader.LoadShader(id)
	}
}
