package debug

import (
	"fmt"
	"image"
	"remembering-home/src/assets/shaders"
	"remembering-home/src/conf"
	"runtime"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

func UpdateDebugUI(sceneName, sceneState string) error {
	if !ShowTexts {
		return nil
	}

	_, err := DebugUI.Update(func(ctx *debugui.Context) error {
		ctx.Window("Debug Overlay", image.Rect(10, 10, 300, 300), func(layout debugui.ContainerLayout) {
			// Performance metrics
			ctx.Text(fmt.Sprintf("OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
			ctx.Text(fmt.Sprintf("Version: %s", conf.GAME_VERSION))
			ctx.Text(fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()))
			ctx.Text(fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()))

			// Scene information
			ctx.Text(fmt.Sprintf("Scene: %s", sceneName))
			ctx.Text(fmt.Sprintf("State: %s", sceneState))

			// Debug controls
			ctx.Header("Debug Controls", false, func() {
				ctx.Button("Toggle Lines").On(func() {
					ShowLines = !ShowLines
				})
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
					uniforms.NoiseScale = 0.1
					uniforms.NoiseIntensity = 0.5
				})
			})
		} else {
			ctx.Text("Unknown shader type")
		}
	})
}
