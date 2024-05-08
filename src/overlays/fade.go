package overlays

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	mask     *ebiten.Image
	vertices []ebiten.Vertex

	curCount int
	maxCount int
	fadeDir  int
)

func init() {
	img := ebiten.NewImage(3, 3)
	img.Fill(color.White)
	mask = img.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	vertices = make([]ebiten.Vertex, 4)

	curCount = 0
	maxCount = 100
	fadeDir = 1
}

func UpdateFade() {
	curCount += fadeDir

	if fadeDir == 1 && curCount > maxCount {
		fadeDir = -1
	} else if fadeDir == -1 && curCount < 0 {
		fadeDir = 1
	}
}

func DrawFade(screen *ebiten.Image) {
	alpha := 1 - float32(curCount)/float32(maxCount)

	for i := range vertices {
		vertices[i].SrcX = 1.0
		vertices[i].SrcY = 1.0
		vertices[i].ColorA = float32(alpha)
	}
	bounds := screen.Bounds()
	vertices[1].DstX = float32(bounds.Dx())
	vertices[2].DstY = float32(bounds.Dy())
	vertices[3].DstX = vertices[1].DstX
	vertices[3].DstY = vertices[2].DstY

	screen.DrawTriangles(vertices, []uint16{0, 1, 2, 1, 2, 3}, mask, nil)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", alpha), 0, 96)
}
