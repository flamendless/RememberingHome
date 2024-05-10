package overlays

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	mask     *ebiten.Image
	vertices []ebiten.Vertex
	curCount int

	FadeAlphaMaxCount int
	FadeAlpha         float32
)

func init() {
	img := ebiten.NewImage(3, 3)
	img.Fill(color.White)
	mask = img.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	vertices = make([]ebiten.Vertex, 4)

	curCount = 0 //0 = totally black
	FadeAlphaMaxCount = 100
}

func IsFadeInFinished() bool {
	return curCount >= FadeAlphaMaxCount
}

func IsFadeOutFinished() bool {
	return curCount <= 0
}

func UpdateFade(fadeDir int) {
	curCount += fadeDir

	if fadeDir == 1 && curCount > FadeAlphaMaxCount {
		fadeDir = -1
	} else if fadeDir == -1 && curCount < 0 {
		fadeDir = 1
	}
}

func DrawFade(screen *ebiten.Image) {
	FadeAlpha = 1 - float32(curCount)/float32(FadeAlphaMaxCount)

	for i := range vertices {
		vertices[i].SrcX = 1.0
		vertices[i].SrcY = 1.0
		vertices[i].ColorA = float32(FadeAlpha)
	}
	bounds := screen.Bounds()
	vertices[1].DstX = float32(bounds.Dx())
	vertices[2].DstY = float32(bounds.Dy())
	vertices[3].DstX = vertices[1].DstX
	vertices[3].DstY = vertices[2].DstY

	screen.DrawTriangles(vertices, []uint16{0, 1, 2, 1, 2, 3}, mask, nil)
}
