package shaders

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ShaderUniforms interface {
	ToShaders(DTSO *ebiten.DrawTrianglesShaderOptions)
}
