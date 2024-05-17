package game

import (
	"image"
	"image/color"
	"nowhere-home/src/effects"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	GetName() string
	GetStateName() string
	Update() error
	Draw(screen *ebiten.Image)
}

type Scene_Manager struct {
	GameState *Game_State
	current   Scene
	next      Scene
	fader     *effects.Fader
	mask      *ebiten.Image
	vertices  []ebiten.Vertex
}

func NewSceneManager() *Scene_Manager {
	img := ebiten.NewImage(3, 3)
	img.Fill(color.White)
	mask := img.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	vertices := make([]ebiten.Vertex, 4)

	return &Scene_Manager{
		mask:     mask,
		vertices: vertices,
		fader:    effects.NewFader(100, 1, 0),
	}
}

func (s *Scene_Manager) Update() error {
	if s.fader.Dir != 0 {
		s.fader.Update()
	}

	if s.fader.Dir == -1 && s.IsFadeInFinished() {
		if s.next == nil {
			s.fader.Dir = 0
		} else {
			s.fader.Dir = 1
			s.current = s.next
			s.next = nil
		}
	} else if s.fader.Dir == 1 && s.IsFadeOutFinished() {
		if s.next == nil {
			s.fader.Dir = 0
		} else {
			s.fader.Dir = -1
			s.current = s.next
			s.next = nil
		}
	}

	return s.current.Update()
}

func (s *Scene_Manager) Draw(screen *ebiten.Image) {
	s.current.Draw(screen)

	for i := range s.vertices {
		s.vertices[i].SrcX = 1.0
		s.vertices[i].SrcY = 1.0
		s.vertices[i].ColorA = float32(s.fader.Alpha / 100)
	}
	bounds := screen.Bounds()
	s.vertices[1].DstX = float32(bounds.Dx())
	s.vertices[2].DstY = float32(bounds.Dy())
	s.vertices[3].DstX = s.vertices[1].DstX
	s.vertices[3].DstY = s.vertices[2].DstY
	screen.DrawTriangles(s.vertices, []uint16{0, 1, 2, 1, 2, 3}, s.mask, nil)
}

func (s *Scene_Manager) GoTo(scene Scene) {
	if s.current == nil {
		s.fader.Dir = -1
		s.current = scene
	} else {
		s.fader.Dir = 1
		s.next = scene
	}
}

func (s *Scene_Manager) IsFadeInFinished() bool {
	return s.fader.Alpha <= 0
}

func (s *Scene_Manager) IsFadeOutFinished() bool {
	return s.fader.Alpha >= s.fader.MaxAlpha
}
