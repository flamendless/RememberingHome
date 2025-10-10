package shaders

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type FadeState int

const (
	FadeStateVisible FadeState = iota // Fully visible (faded in) - can only trigger fade out
	FadeStateHidden                   // Fully hidden (faded out) - can only trigger fade in
)

type SilentHillRedShaderUniforms struct {
	ID             string
	Time           float64
	BannerPos      [2]float64
	BannerSize     [2]float64
	BaseRedColor   [4]float64
	GlowIntensity  float64
	MetallicShine  float64
	EdgeDarkness   float64
	TextGlowRadius float64
	NoiseScale     float64
	NoiseIntensity float64
	// Cell movement parameters
	MovementSpeed1 float64
	MovementSpeed2 float64
	MovementSpeed3 float64
	MovementRange1 float64
	MovementRange2 float64
	MovementRange3 float64
	// Spot generation parameters
	LargeSpotScale     float64
	MediumSpotScale    float64
	SmallSpotScale     float64
	LargeSpotThreshold float64
	SmallSpotThreshold float64
	PulseSpeed         float64
	PulseIntensity     float64
	FadeProgress       float64
	FadeState          int
	// Internal state to prevent spamming and track current state
	currentFadeState    FadeState
	initialFadeState    FadeState
	lastFadeTriggerTime float64
	fadeCooldown        float64
	// Fade animation state
	isAnimating   bool
	fadeStartTime float64
	fadeDuration  float64
	fadeDirection bool // true = fade in, false = fade out
}

func (shrsu *SilentHillRedShaderUniforms) ToShaders(dtso *ebiten.DrawTrianglesShaderOptions) {
	dtso.Uniforms = map[string]any{
		"ID":             "SilentHillRed",
		"Time":           shrsu.Time,
		"BannerPos":      shrsu.BannerPos,
		"BannerSize":     shrsu.BannerSize,
		"BaseRedColor":   shrsu.BaseRedColor,
		"GlowIntensity":  shrsu.GlowIntensity,
		"MetallicShine":  shrsu.MetallicShine,
		"EdgeDarkness":   shrsu.EdgeDarkness,
		"TextGlowRadius": shrsu.TextGlowRadius,
		"NoiseScale":     shrsu.NoiseScale,
		"NoiseIntensity": shrsu.NoiseIntensity,
		// Cell movement parameters
		"MovementSpeed1": shrsu.MovementSpeed1,
		"MovementSpeed2": shrsu.MovementSpeed2,
		"MovementSpeed3": shrsu.MovementSpeed3,
		"MovementRange1": shrsu.MovementRange1,
		"MovementRange2": shrsu.MovementRange2,
		"MovementRange3": shrsu.MovementRange3,
		// Spot generation parameters
		"LargeSpotScale":     shrsu.LargeSpotScale,
		"MediumSpotScale":    shrsu.MediumSpotScale,
		"SmallSpotScale":     shrsu.SmallSpotScale,
		"LargeSpotThreshold": shrsu.LargeSpotThreshold,
		"SmallSpotThreshold": shrsu.SmallSpotThreshold,
		"PulseSpeed":         shrsu.PulseSpeed,
		"PulseIntensity":     shrsu.PulseIntensity,
		"FadeProgress":       shrsu.FadeProgress,
		"FadeDirection":      shrsu.fadeDirectionToFloat(),
		"InitialFadeState":   float64(shrsu.initialFadeState),
	}
}

func NewSilentHillRedShaderUniforms(initialFadeState FadeState) *SilentHillRedShaderUniforms {
	uniforms := &SilentHillRedShaderUniforms{
		ID:             "SilentHillRed",
		Time:           0,
		BannerPos:      [2]float64{0, 0},
		BannerSize:     [2]float64{100, 50},
		BaseRedColor:   [4]float64{1.0, 0.2, 0.2, 1.0},
		GlowIntensity:  1.0,
		MetallicShine:  0.4,
		EdgeDarkness:   0.3,
		TextGlowRadius: 2.5,
		NoiseScale:     0.2,
		NoiseIntensity: 0.1,
		MovementSpeed1: 0.7,
		MovementSpeed2: 0.5,
		MovementSpeed3: 0.3,
		MovementRange1: 15.0,
		MovementRange2: 12.0,
		MovementRange3: 18.0,
		LargeSpotScale:     0.06,
		MediumSpotScale:    0.1,
		SmallSpotScale:     0.15,
		LargeSpotThreshold: 0.85,
		SmallSpotThreshold: 0.95,
		PulseSpeed:         2.0,
		PulseIntensity:     0.3,
		FadeProgress:       0.0,
		fadeCooldown:        0.1,
		lastFadeTriggerTime: -1.0,
		isAnimating:   false,
		fadeStartTime: 0.0,
		fadeDuration:  0.0,
		fadeDirection: false,
	}

	uniforms.currentFadeState = initialFadeState
	uniforms.initialFadeState = initialFadeState

	return uniforms
}

func (shrsu *SilentHillRedShaderUniforms) TriggerFadeOut(duration float64) {
	if shrsu.currentFadeState != FadeStateVisible {
		return
	}
	if shrsu.Time-shrsu.lastFadeTriggerTime < shrsu.fadeCooldown {
		return
	}

	shrsu.FadeProgress = 0.0
	shrsu.isAnimating = true
	shrsu.fadeStartTime = shrsu.Time
	shrsu.fadeDuration = duration
	shrsu.fadeDirection = false
	shrsu.lastFadeTriggerTime = shrsu.Time
	shrsu.currentFadeState = FadeStateHidden
}

func (shrsu *SilentHillRedShaderUniforms) TriggerFadeIn(duration float64) {
	if shrsu.currentFadeState != FadeStateHidden {
		return
	}
	if shrsu.Time-shrsu.lastFadeTriggerTime < shrsu.fadeCooldown {
		return
	}

	shrsu.FadeProgress = 0.0
	shrsu.isAnimating = true
	shrsu.fadeStartTime = shrsu.Time
	shrsu.fadeDuration = duration
	shrsu.fadeDirection = true
	shrsu.lastFadeTriggerTime = shrsu.Time
	shrsu.currentFadeState = FadeStateVisible
}

func (shrsu *SilentHillRedShaderUniforms) GetCurrentFadeState() FadeState {
	return shrsu.currentFadeState
}

func (shrsu *SilentHillRedShaderUniforms) fadeDirectionToFloat() float64 {
	if shrsu.fadeDirection {
		return 1.0
	}
	return 0.0
}

// TODO: (Brandon) add Update to interface
func (shrsu *SilentHillRedShaderUniforms) Update() {
	if !shrsu.isAnimating {
		return
	}

	elapsed := shrsu.Time - shrsu.fadeStartTime
	if elapsed >= shrsu.fadeDuration {
		shrsu.FadeProgress = 1.0
		shrsu.isAnimating = false
	} else {
		progress := elapsed / shrsu.fadeDuration
		shrsu.FadeProgress = progress
	}
}

var _ ShaderUniforms = (*SilentHillRedShaderUniforms)(nil)
