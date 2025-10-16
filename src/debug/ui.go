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
			debugSilentHillShader(ctx, uniforms)
		}
	})
}

func debugSilentHillShader(ctx *debugui.Context, uniforms *shaders.SilentHillRedShaderUniforms) {
	ctx.Header("Silent Hill Red Shader", false, func() {
		ctx.Header("Base Color", false, func() {
			ctx.Text("Red:")
			ctx.SliderF(&uniforms.BaseRedColor[0], 0.0, 1.0, 0.01, 2)
			ctx.Text("Green:")
			ctx.SliderF(&uniforms.BaseRedColor[1], 0.0, 1.0, 0.01, 2)
			ctx.Text("Blue:")
			ctx.SliderF(&uniforms.BaseRedColor[2], 0.0, 1.0, 0.01, 2)
			ctx.Text("Alpha:")
			ctx.SliderF(&uniforms.BaseRedColor[3], 0.0, 1.0, 0.01, 2)
		})

		ctx.Header("Effect Parameters", false, func() {
			ctx.Text("Glow Intensity:")
			ctx.SliderF(&uniforms.GlowIntensity, 0.0, 2.0, 0.01, 2)
			ctx.Text("Metallic Shine:")
			ctx.SliderF(&uniforms.MetallicShine, 0.0, 1.0, 0.01, 2)
			ctx.Text("Edge Darkness:")
			ctx.SliderF(&uniforms.EdgeDarkness, 0.0, 1.0, 0.01, 2)
			ctx.Text("Text Glow Radius:")
			ctx.SliderF(&uniforms.TextGlowRadius, 0.0, 5.0, 0.01, 2)
		})

		ctx.Header("Noise Parameters", false, func() {
			ctx.Text("Noise Scale:")
			ctx.SliderF(&uniforms.NoiseScale, 0.0, 1.0, 0.001, 3)
			ctx.Text("Noise Intensity:")
			ctx.SliderF(&uniforms.NoiseIntensity, 0.0, 1.0, 0.01, 2)
		})

		ctx.Header("Cell Movement", false, func() {
			ctx.Text("Movement Speed 1:")
			ctx.SliderF(&uniforms.MovementSpeed1, 0.0, 3.0, 0.01, 2)
			ctx.Text("Movement Speed 2:")
			ctx.SliderF(&uniforms.MovementSpeed2, 0.0, 3.0, 0.01, 2)
			ctx.Text("Movement Speed 3:")
			ctx.SliderF(&uniforms.MovementSpeed3, 0.0, 3.0, 0.01, 2)
			ctx.Text("Movement Range 1:")
			ctx.SliderF(&uniforms.MovementRange1, 0.0, 50.0, 0.1, 1)
			ctx.Text("Movement Range 2:")
			ctx.SliderF(&uniforms.MovementRange2, 0.0, 50.0, 0.1, 1)
			ctx.Text("Movement Range 3:")
			ctx.SliderF(&uniforms.MovementRange3, 0.0, 50.0, 0.1, 1)
		})

		ctx.Header("Spot Generation", false, func() {
			ctx.Text("Large Spot Scale:")
			ctx.SliderF(&uniforms.LargeSpotScale, 0.01, 0.5, 0.001, 3)
			ctx.Text("Medium Spot Scale:")
			ctx.SliderF(&uniforms.MediumSpotScale, 0.01, 0.5, 0.001, 3)
			ctx.Text("Small Spot Scale:")
			ctx.SliderF(&uniforms.SmallSpotScale, 0.01, 0.5, 0.001, 3)
			ctx.Text("Large Spot Threshold:")
			ctx.SliderF(&uniforms.LargeSpotThreshold, 0.0, 1.0, 0.01, 2)
			ctx.Text("Small Spot Threshold:")
			ctx.SliderF(&uniforms.SmallSpotThreshold, 0.0, 1.0, 0.01, 2)
			ctx.Text("Pulse Speed:")
			ctx.SliderF(&uniforms.PulseSpeed, 0.0, 5.0, 0.01, 2)
			ctx.Text("Pulse Intensity:")
			ctx.SliderF(&uniforms.PulseIntensity, 0.0, 1.0, 0.01, 2)

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
			ctx.SliderF(&uniforms.BannerPos[0], 0.0, 800.0, 1.0, 0)
			ctx.Text("Banner Y:")
			ctx.SliderF(&uniforms.BannerPos[1], 0.0, 600.0, 1.0, 0)
			ctx.Text("Banner Width:")
			ctx.SliderF(&uniforms.BannerSize[0], 10.0, 300.0, 1.0, 0)
			ctx.Text("Banner Height:")
			ctx.SliderF(&uniforms.BannerSize[1], 5.0, 100.0, 1.0, 0)
		})

		ctx.Button("Reset to Defaults").On(func() {
			uniforms.ResetToInitial()
		})
	})
}
