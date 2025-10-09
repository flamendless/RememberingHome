package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	GetName() string
	GetStateName() string
	Update() error
	Draw(screen *ebiten.Image)
}
