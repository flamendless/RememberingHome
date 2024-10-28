package game

import (
	"fmt"
	"nowhere-home/src/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Layer struct {
	ID       string
	Canvas   *ebiten.Image
	Disabled bool

	X        float64
	Y        float64
	ScaleX   float64
	ScaleY   float64
	DIO      *ebiten.DrawImageOptions
	Shader   *ebiten.Shader
	DRSO     *ebiten.DrawRectShaderOptions
	Uniforms assets.ShaderUniforms

	DTSO     *ebiten.DrawTrianglesShaderOptions
	Vertices []ebiten.Vertex
	Indices  []uint16
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

func NewLayerWithTriangleShader(id string, width, height int, shader *ebiten.Shader) *Layer {
	return &Layer{
		ID:     id,
		Canvas: ebiten.NewImage(width, height),
		DTSO:   &ebiten.DrawTrianglesShaderOptions{},
		Shader: shader,
		ScaleX: 1,
		ScaleY: 1,
	}
}

func (layer *Layer) Render(screen *ebiten.Image) {
	screen.DrawImage(layer.Canvas, layer.DIO)
}

func (layer *Layer) RenderWithShader(screen *ebiten.Image) {
	if layer.Shader == nil {
		panic(fmt.Sprintf("No shader was given to layer %s", layer.ID))
	}

	if layer.Disabled {
		screen.DrawImage(layer.Canvas, nil)
		return
	}

	if len(layer.Vertices) == 0 {
		layer.DRSO.Images[0] = layer.Canvas
		w, h := layer.Canvas.Bounds().Dx(), layer.Canvas.Bounds().Dy()
		screen.DrawRectShader(w, h, layer.Shader, layer.DRSO)
		return
	}

	layer.DTSO.Images[0] = layer.Canvas
	screen.DrawTrianglesShader(layer.Vertices, layer.Indices, layer.Shader, layer.DTSO)
}

func (layer *Layer) ApplyTransformation() {
	if layer.Shader == nil {
		layer.DIO.GeoM.Reset()
		layer.DIO.GeoM.Scale(layer.ScaleX, layer.ScaleY)
		layer.DIO.GeoM.Translate(layer.X, layer.Y)
	} else if len(layer.Vertices) == 0 {
		layer.DRSO.GeoM.Reset()
		layer.DRSO.GeoM.Scale(layer.ScaleX, layer.ScaleY)
		layer.DRSO.GeoM.Translate(layer.X, layer.Y)
	}
}
