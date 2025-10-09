package scenes

import "github.com/hajimehoshi/ebiten/v2"

type SceneManager interface {
	Update() error
	Draw(screen *ebiten.Image)
	GoTo(scene Scene)
	IsFadeInFinished() bool
	IsFadeOutFinished() bool
	IsFading() bool
}
