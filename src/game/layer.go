package game

import (
	"fmt"
	"nowhere-home/src/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Layer struct {
	ID       string
	Canvas   *ebiten.Image
	X        float64
	Y        float64
	ScaleX   float64
	ScaleY   float64
	DIO      *ebiten.DrawImageOptions
	Shader   *ebiten.Shader
	DRSO     *ebiten.DrawRectShaderOptions
	Uniforms assets.ShaderUniforms
}

func NewLayer(id string, width, height int) *Layer {
	return &Layer{
		ID:     id,
		Canvas: ebiten.NewImage(width, height),
		DIO: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterNearest,
		},
		ScaleX: 1,
		ScaleY: 1,
	}
}

func NewLayerWithShader(id string, width, height int, shader *ebiten.Shader) *Layer {
	return &Layer{
		ID:     id,
		Canvas: ebiten.NewImage(width, height),
		DRSO:   &ebiten.DrawRectShaderOptions{},
		Shader: shader,
		ScaleX: 1,
		ScaleY: 1,
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

func (layer *Layer) ApplyTransformation() {
	if layer.Shader == nil {
		layer.DIO.GeoM.Reset()
		layer.DIO.GeoM.Scale(layer.ScaleX, layer.ScaleY)
		layer.DIO.GeoM.Translate(layer.X, layer.Y)
	} else {
		layer.DRSO.GeoM.Reset()
		layer.DRSO.GeoM.Scale(layer.ScaleX, layer.ScaleY)
		layer.DRSO.GeoM.Translate(layer.X, layer.Y)
	}
}
