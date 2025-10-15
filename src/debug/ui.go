package debug

import (
	"fmt"
	"image"
	"remembering-home/src/assets/shaders"
	"remembering-home/src/conf"
	"remembering-home/src/context"
	"runtime"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type SceneNavigator interface {
	NavigateToDummy()
	NavigateToSplash()
	NavigateToMainMenu()
}

var sceneNavigator SceneNavigator

func SetSceneNavigator(navigator SceneNavigator) {
	sceneNavigator = navigator
}

func UpdateDebugUI(context *context.GameContext, sceneName, sceneState string) error {
	if !ShowTexts {
		return nil
	}

	_, err := DebugUI.Update(func(ctx *debugui.Context) error {
		ctx.Window("Debug Overlay", image.Rect(10, 10, 300, 300), func(layout debugui.ContainerLayout) {
			ctx.Text(fmt.Sprintf("OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
			ctx.Text(fmt.Sprintf("Version/Dev: %s/%v", conf.GAME_VERSION, conf.DEV))
			ctx.Text(fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()))
			ctx.Text(fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()))
			ctx.Text(fmt.Sprintf("Mode/Quality: %s, %s", context.Settings.Window.String(), context.Settings.Quality.String()))
			ctx.Text(fmt.Sprintf("Volume/Music: %d, %d", context.Settings.Volume, context.Settings.Music))

			// Scene information
			ctx.Text(fmt.Sprintf("Scene: %s", sceneName))
			ctx.Text(fmt.Sprintf("State: %s", sceneState))

			// Debug controls
			ctx.Header("Debug Controls", false, func() {
				ctx.Checkbox(&ShowLines, "Show Center Lines")
				if sceneNavigator != nil {
					ctx.Button("Go to Dummy").On(func() {
						sceneNavigator.NavigateToDummy()
					})
					ctx.Button("Go to Splash").On(func() {
						sceneNavigator.NavigateToSplash()
					})
					ctx.Button("Go to Main Menu").On(func() {
						sceneNavigator.NavigateToMainMenu()
					})
				}
			})

			AddShaderDebugControls(ctx)
		})
		return nil
	})

	return err
}

func AddShaderDebugControls(ctx *debugui.Context) {
	ctx.Header("Shader Debug", false, func() {
		if CurrentDebugShader == nil {
			ctx.Text("No shader set for debugging")
			return
		}

		if uniforms, ok := CurrentDebugShader.(*shaders.SilentHillRedShaderUniforms); ok {
			ctx.Header("Silent Hill Red Shader", false, func() {
				ctx.Header("Base Color", false, func() {
					ctx.Text("Red:")
					ctx.SliderF((*float64)(&uniforms.BaseRedColor[0]), 0.0, 1.0, 0.01, 2)
					ctx.Text("Green:")
					ctx.SliderF((*float64)(&uniforms.BaseRedColor[1]), 0.0, 1.0, 0.01, 2)
					ctx.Text("Blue:")
					ctx.SliderF((*float64)(&uniforms.BaseRedColor[2]), 0.0, 1.0, 0.01, 2)
					ctx.Text("Alpha:")
					ctx.SliderF((*float64)(&uniforms.BaseRedColor[3]), 0.0, 1.0, 0.01, 2)
				})

				ctx.Header("Effect Parameters", false, func() {
					ctx.Text("Glow Intensity:")
					ctx.SliderF((*float64)(&uniforms.GlowIntensity), 0.0, 2.0, 0.01, 2)
					ctx.Text("Metallic Shine:")
					ctx.SliderF((*float64)(&uniforms.MetallicShine), 0.0, 1.0, 0.01, 2)
					ctx.Text("Edge Darkness:")
					ctx.SliderF((*float64)(&uniforms.EdgeDarkness), 0.0, 1.0, 0.01, 2)
					ctx.Text("Text Glow Radius:")
					ctx.SliderF((*float64)(&uniforms.TextGlowRadius), 0.0, 5.0, 0.01, 2)
				})

				ctx.Header("Noise Parameters", false, func() {
					ctx.Text("Noise Scale:")
					ctx.SliderF((*float64)(&uniforms.NoiseScale), 0.0, 1.0, 0.001, 3)
					ctx.Text("Noise Intensity:")
					ctx.SliderF((*float64)(&uniforms.NoiseIntensity), 0.0, 1.0, 0.01, 2)
				})

				ctx.Header("Cell Movement", false, func() {
					ctx.Text("Movement Speed 1:")
					ctx.SliderF((*float64)(&uniforms.MovementSpeed1), 0.0, 3.0, 0.01, 2)
					ctx.Text("Movement Speed 2:")
					ctx.SliderF((*float64)(&uniforms.MovementSpeed2), 0.0, 3.0, 0.01, 2)
					ctx.Text("Movement Speed 3:")
					ctx.SliderF((*float64)(&uniforms.MovementSpeed3), 0.0, 3.0, 0.01, 2)
					ctx.Text("Movement Range 1:")
					ctx.SliderF((*float64)(&uniforms.MovementRange1), 0.0, 50.0, 0.1, 1)
					ctx.Text("Movement Range 2:")
					ctx.SliderF((*float64)(&uniforms.MovementRange2), 0.0, 50.0, 0.1, 1)
					ctx.Text("Movement Range 3:")
					ctx.SliderF((*float64)(&uniforms.MovementRange3), 0.0, 50.0, 0.1, 1)
				})

				ctx.Header("Spot Generation", false, func() {
					ctx.Text("Large Spot Scale:")
					ctx.SliderF((*float64)(&uniforms.LargeSpotScale), 0.01, 0.5, 0.001, 3)
					ctx.Text("Medium Spot Scale:")
					ctx.SliderF((*float64)(&uniforms.MediumSpotScale), 0.01, 0.5, 0.001, 3)
					ctx.Text("Small Spot Scale:")
					ctx.SliderF((*float64)(&uniforms.SmallSpotScale), 0.01, 0.5, 0.001, 3)
					ctx.Text("Large Spot Threshold:")
					ctx.SliderF((*float64)(&uniforms.LargeSpotThreshold), 0.0, 1.0, 0.01, 2)
					ctx.Text("Small Spot Threshold:")
					ctx.SliderF((*float64)(&uniforms.SmallSpotThreshold), 0.0, 1.0, 0.01, 2)
					ctx.Text("Pulse Speed:")
					ctx.SliderF((*float64)(&uniforms.PulseSpeed), 0.0, 5.0, 0.01, 2)
					ctx.Text("Pulse Intensity:")
					ctx.SliderF((*float64)(&uniforms.PulseIntensity), 0.0, 1.0, 0.01, 2)

					currentState := uniforms.GetCurrentFadeState()
					var stateText string
					switch currentState {
					case shaders.FadeStateVisible:
						stateText = "Visible (can fade out)"
					case shaders.FadeStateHidden:
						stateText = "Hidden (can fade in)"
					default:
						stateText = "Unknown"
					}
					ctx.Text(fmt.Sprintf("Current State: %s", stateText))

					ctx.Button("Trigger Fade In").On(func() {
						uniforms.TriggerFadeIn(4.0)
					})

					ctx.Button("Trigger Fade Out").On(func() {
						uniforms.TriggerFadeOut(4.0)
					})
				})

				ctx.Header("Banner Properties", false, func() {
					ctx.Text("Banner X:")
					ctx.SliderF((*float64)(&uniforms.BannerPos[0]), 0.0, 800.0, 1.0, 0)
					ctx.Text("Banner Y:")
					ctx.SliderF((*float64)(&uniforms.BannerPos[1]), 0.0, 600.0, 1.0, 0)
					ctx.Text("Banner Width:")
					ctx.SliderF((*float64)(&uniforms.BannerSize[0]), 10.0, 300.0, 1.0, 0)
					ctx.Text("Banner Height:")
					ctx.SliderF((*float64)(&uniforms.BannerSize[1]), 5.0, 100.0, 1.0, 0)
				})

				ctx.Button("Reset to Defaults").On(func() {
					uniforms.BaseRedColor = [4]float64{1.0, 0.2, 0.2, 1.0}
					uniforms.GlowIntensity = 1.0
					uniforms.MetallicShine = 0.4
					uniforms.EdgeDarkness = 0.3
					uniforms.TextGlowRadius = 2.5
					uniforms.NoiseScale = 0.2
					uniforms.NoiseIntensity = 0.1
					// Cell movement defaults
					uniforms.MovementSpeed1 = 0.7
					uniforms.MovementSpeed2 = 0.5
					uniforms.MovementSpeed3 = 0.3
					uniforms.MovementRange1 = 15.0
					uniforms.MovementRange2 = 12.0
					uniforms.MovementRange3 = 18.0
					// Spot generation defaults
					uniforms.LargeSpotScale = 0.06
					uniforms.MediumSpotScale = 0.1
					uniforms.SmallSpotScale = 0.15
					uniforms.LargeSpotThreshold = 0.85
					uniforms.SmallSpotThreshold = 0.95
					uniforms.PulseSpeed = 2.0
					uniforms.PulseIntensity = 0.3
				})
			})
		} else {
			ctx.Text("Unknown shader type")
		}
	})
}
