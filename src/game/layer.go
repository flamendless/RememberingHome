package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Layer struct {
	ID     string
	Canvas *ebiten.Image
	DIO    *ebiten.DrawImageOptions
	Shader *ebiten.Shader
	DRSO   *ebiten.DrawRectShaderOptions
}

func NewLayer(id string, width, height int) *Layer {
	return &Layer{
		ID:     id,
		Canvas: ebiten.NewImage(width, height),
		DIO:    &ebiten.DrawImageOptions{},
	}
}

func NewLayerWithShader(id string, width, height int, shader *ebiten.Shader) *Layer {
	return &Layer{
		ID:     id,
		Canvas: ebiten.NewImage(width, height),
		DRSO:   &ebiten.DrawRectShaderOptions{},
		Shader: shader,
	}
}

func (layer *Layer) Apply(screen *ebiten.Image) {
	screen.DrawImage(layer.Canvas, layer.DIO)
}

func (layer *Layer) ApplyShader(screen *ebiten.Image) {
	if layer.Shader == nil {
		panic(fmt.Sprintf("No shader was given to layer %s", layer.ID))
	}
	layer.DRSO.Images[0] = layer.Canvas
	w, h := layer.Canvas.Bounds().Dx(), layer.Canvas.Bounds().Dy()
	screen.DrawRectShader(w, h, layer.Shader, layer.DRSO)
}
