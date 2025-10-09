package debug

import (
	"image/color"
	"remembering-home/src/assets/shaders"
	"remembering-home/src/conf"
	"runtime"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	WSLTricked         bool
	ShowTexts          bool
	ShowLines          bool
	DebugUI            debugui.DebugUI
	CurrentDebugShader shaders.ShaderUniforms
)

const (
	VERSION int = iota
	FPS
	TPS
	OVERLAY_FADE_ALPHA
	SCENE_NAME
	SCENE_STATE
)

func init() {
	ShowTexts = conf.DEV
}

func FixWSLWindow() {
	if !WSLTricked && !ebiten.IsFocused() && runtime.GOARCH != "wasm" {
		ebiten.MinimizeWindow()
		ebiten.MaximizeWindow()
		ebiten.RestoreWindow()
		WSLTricked = true
	}
}

// SetDebugShader sets the current shader to debug
func SetDebugShader(uniforms shaders.ShaderUniforms) {
	CurrentDebugShader = uniforms
}

// ClearDebugShader clears the current debug shader
func ClearDebugShader() {
	CurrentDebugShader = nil
}

func DrawDebugOverlay(screen *ebiten.Image) {
	if ShowTexts {
		DebugUI.Draw(screen)
	}

	if ShowLines {
		clr := color.White
		vector.StrokeLine(screen, conf.GAME_W/2, 0, conf.GAME_W/2, conf.GAME_H, 1, clr, false)
		vector.StrokeLine(screen, 0, conf.GAME_H/2, conf.GAME_W, conf.GAME_H/2, 1, clr, false)
	}
}
