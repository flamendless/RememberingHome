//credit: https://github.com/setanarut/bulimia/blob/main/engine/animplayer.go

package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	Frames []*ebiten.Image
	FPS    float64
	Name   string
}

type AnimationPlayer struct {
	SpriteSheet       *ebiten.Image
	CurrentFrame      *ebiten.Image
	DIO               *ebiten.DrawImageOptions
	Animations        map[string]*Animation
	Paused            bool
	CurrentFrameIndex int
	CurrentState      string
	Tick              float64
}

func NewAnimationPlayer(spriteSheet *ebiten.Image) *AnimationPlayer {
	return &AnimationPlayer{
		SpriteSheet:       spriteSheet,
		Paused:            false,
		Animations:        make(map[string]*Animation),
		CurrentFrameIndex: 0,
		DIO:               &ebiten.DrawImageOptions{},
	}

}

func (ap *AnimationPlayer) AddStateAnimation(
	stateName string,
	x, y, w, h, frameCount int,
	pingpong bool,
) *Animation {
	subImages := []*ebiten.Image{}
	frameRect := image.Rect(x, y, x+w, y+h)
	for i := 0; i < frameCount; i++ {
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
	ap.CurrentFrame = ap.Animations[ap.CurrentState].Frames[ap.CurrentFrameIndex]

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
	return ap.Animations[ap.State()].FPS
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
	if frameIndex < len(ap.Animations[ap.State()].Frames) && frameIndex >= 0 {
		ap.Paused = true
		ap.CurrentFrameIndex = frameIndex
	}
}

func (ap *AnimationPlayer) GetLastFrameCount() int {
	return len(ap.Animations[ap.CurrentState].Frames) - 1
}

func (ap *AnimationPlayer) IsInLastFrame() bool {
	return ap.CurrentFrameIndex == ap.GetLastFrameCount()
}

func (ap *AnimationPlayer) Update() {
	if !ap.Paused {
		ap.Tick += ap.Animations[ap.CurrentState].FPS / 60.0
		ap.CurrentFrameIndex = int(math.Floor(ap.Tick))
		if ap.CurrentFrameIndex >= len(ap.Animations[ap.CurrentState].Frames) {
			ap.Tick = 0
			ap.CurrentFrameIndex = 0
		}
	}

	ap.CurrentFrame = ap.Animations[ap.CurrentState].Frames[ap.CurrentFrameIndex]
}
