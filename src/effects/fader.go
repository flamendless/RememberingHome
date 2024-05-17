package effects

import "github.com/hajimehoshi/ebiten/v2"

type Fader struct {
	Alpha    float32
	Amount   float32
	Dir      int
	MaxAlpha float32
	Stopped  bool
}

func NewFader(alpha, amount float32, dir int) *Fader {
	if amount < 0 {
		panic("amount must positive")
	}

	return &Fader{
		Alpha:    alpha,
		Amount:   amount,
		Dir:      dir,
		MaxAlpha: 100,
	}
}

func (f *Fader) Update() {
	if f.Stopped {
		return
	}

	if f.Alpha >= f.MaxAlpha {
		f.Dir = -1
	} else if f.Alpha <= 0 {
		f.Dir = 1
	}

	f.Alpha += float32(f.Dir) * f.Amount
}

func (f *Fader) GetCS() *ebiten.ColorScale {
	cs := ebiten.ColorScale{}
	cs.ScaleAlpha(f.Alpha / f.MaxAlpha)
	return &cs
}
