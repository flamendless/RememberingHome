package graphics

import (
	"image"
	"math"
	"remembering-home/src/errs"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	Name   string
	Frames []*ebiten.Image
	FPS    float64
}

type AnimationPlayer struct {
	SpriteSheet       *ebiten.Image
	CurrentFrame      *ebiten.Image
	DIO               *ebiten.DrawImageOptions
	DRSO              *ebiten.DrawRectShaderOptions
	Animations        map[string]*Animation
	CurrentState      string
	CurrentFrameIndex int
	Tick              float64
	Paused            bool
}

func NewAnimationPlayer(spriteSheet *ebiten.Image) *AnimationPlayer {
	return &AnimationPlayer{
		SpriteSheet:       spriteSheet,
		Paused:            false,
		Animations:        make(map[string]*Animation),
		CurrentFrameIndex: 0,
		DIO: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterNearest,
		},
	}
}

func (ap *AnimationPlayer) AddStateAnimation(
	stateName string,
	x, y, w, h, frameCount int,
	pingpong bool,
) *Animation {
	subImages := []*ebiten.Image{}
	frameRect := image.Rect(x, y, x+w, y+h)
	for range frameCount {
		subImages = append(subImages, ap.SpriteSheet.SubImage(frameRect).(*ebiten.Image))
		frameRect.Min.X += w
		frameRect.Max.X += w
	}

	if pingpong {
		for i := frameCount - 2; i > 1; i-- {
			subImages = append(subImages, subImages[i])
		}
	}

	anim := &Animation{
		FPS:    15.0,
		Frames: subImages,
		Name:   stateName,
	}

	ap.CurrentState = stateName
	ap.Animations[stateName] = anim

	dataAnim, ok := ap.Animations[ap.CurrentState]
	if !ok {
		panic(errs.ErrNilElem)
	}
	ap.CurrentFrame = dataAnim.Frames[ap.CurrentFrameIndex]

	return anim
}

func (ap *AnimationPlayer) SetFPS(fps float64) {
	for _, anim := range ap.Animations {
		anim.FPS = fps
	}
}

func (ap *AnimationPlayer) AddAnimation(a *Animation) {
	ap.Animations[a.Name] = a
}

func (ap *AnimationPlayer) State() string {
	return ap.CurrentState
}

func (ap *AnimationPlayer) CurrentStateFPS() float64 {
	dataAnim, ok := ap.Animations[ap.State()]
	if !ok {
		panic(errs.ErrNilElem)
	}
	return dataAnim.FPS
}

func (ap *AnimationPlayer) SetStateReset(state string) {
	if ap.CurrentState != state {
		ap.CurrentState = state
		ap.Tick = 0
		ap.CurrentFrameIndex = 0
	}
}

func (ap *AnimationPlayer) SetState(state string) {
	if ap.CurrentState != state {
		ap.CurrentState = state
	}
}

func (ap *AnimationPlayer) PauseAtFrame(frameIndex int) {
	dataAnim, ok := ap.Animations[ap.State()]
	if !ok {
		panic(errs.ErrNilElem)
	}

	if frameIndex < len(dataAnim.Frames) && frameIndex >= 0 {
		ap.Paused = true
		ap.CurrentFrameIndex = frameIndex
	}
}

func (ap *AnimationPlayer) GetLastFrameCount() int {
	dataAnim, ok := ap.Animations[ap.CurrentState]
	if !ok {
		panic(errs.ErrNilElem)
	}
	return len(dataAnim.Frames) - 1
}

func (ap *AnimationPlayer) IsInLastFrame() bool {
	return ap.CurrentFrameIndex == ap.GetLastFrameCount()
}

func (ap *AnimationPlayer) Update() {
	dataAnim, ok := ap.Animations[ap.CurrentState]
	if !ok {
		panic(errs.ErrNilElem)
	}

	if !ap.Paused {
		ap.Tick += dataAnim.FPS / 60.0
		ap.CurrentFrameIndex = int(math.Floor(ap.Tick))
		if ap.CurrentFrameIndex >= len(dataAnim.Frames) {
			ap.Tick = 0
			ap.CurrentFrameIndex = 0
		}
	}

	ap.CurrentFrame = dataAnim.Frames[ap.CurrentFrameIndex]
}
