package game

import (
	"nowhere-home/src/overlays"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	GetName() string
	Update() error
	Draw(screen *ebiten.Image)
}

const transitionMaxCount = 20

type Scene_Manager struct {
	current         Scene
	next            Scene
	transitionCount int
}

func (s *Scene_Manager) Update() error {
	overlays.UpdateFade()

	if s.transitionCount == 0 {
		return s.current.Update()
	}

	s.transitionCount--
	if s.transitionCount > 0 {
		return nil
	}

	s.current = s.next
	s.next = nil
	return nil
}

func (s *Scene_Manager) Draw(screen *ebiten.Image) {
	if s.transitionCount == 0 {
		s.current.Draw(screen)
		return
	}

	overlays.DrawFade(screen)
}

func (s *Scene_Manager) GoTo(scene Scene) {
	if s.current == nil {
		s.current = scene
	} else {
		s.next = scene
		s.transitionCount = transitionMaxCount
	}
}
