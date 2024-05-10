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

type Scene_Manager struct {
	GameState *Game_State
	current   Scene
	next      Scene
	fadeDir   int //1 = fading in, -1 fading out
}

func (s *Scene_Manager) Update() error {
	overlays.UpdateFade(s.fadeDir)

	if s.fadeDir == 1 && overlays.IsFadeInFinished() {
		if s.next == nil {
			s.fadeDir = 0
		} else {
			s.fadeDir = -1
			s.current = s.next
			s.next = nil
		}
	} else if s.fadeDir == -1 && overlays.IsFadeOutFinished() {
		if s.next == nil {
			s.fadeDir = 0
		} else {
			s.fadeDir = 1
			s.current = s.next
			s.next = nil
		}
	}

	return s.current.Update()
}

func (s *Scene_Manager) Draw(screen *ebiten.Image) {
	s.current.Draw(screen)
	overlays.DrawFade(screen)
}

func (s *Scene_Manager) GoTo(scene Scene) {
	if s.current == nil {
		s.fadeDir = 1
		s.current = scene
	} else {
		s.fadeDir = -1
		s.next = scene
	}
}
