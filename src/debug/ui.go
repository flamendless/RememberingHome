package debug

import (
	"fmt"
	"image"
	"remembering-home/src/assets/shaders"
	"remembering-home/src/conf"
	"remembering-home/src/context"
	"remembering-home/src/enums"
	"remembering-home/src/graphics"
	"runtime"
	"sort"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type SceneNavigator interface {
	NavigateTo(scene enums.Scene)
}

var sceneNavigator SceneNavigator
var zoomLevel float64 = 1.0

func SetSceneNavigator(navigator SceneNavigator) {
	sceneNavigator = navigator
}

func GetZoomLevel() float64 {
	return zoomLevel
}

func SetZoomLevel(level float64) {
	zoomLevel = level
}

func UpdateDebugUI(context *context.GameContext, sceneName, sceneState string, itemRenderer *graphics.ItemRenderer) error {
	if !ShowTexts {
		return nil
	}

	_, err := DebugUI.Update(func(ctx *debugui.Context) error {
		ctx.Window("Debug Overlay", image.Rect(10, 10, 300, 500), func(layout debugui.ContainerLayout) {
			ctx.Text(fmt.Sprintf("OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
			ctx.Text(fmt.Sprintf("Version/Dev: %s/%v", conf.GAME_VERSION, conf.DEV))
			ctx.Text(fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()))
			ctx.Text(fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()))
			ctx.Text(fmt.Sprintf("Mode/Quality: %s, %s", context.Settings.Window.String(), context.Settings.Quality.String()))
			ctx.Text(fmt.Sprintf("Volume/Music: %d, %d", context.Settings.Volume, context.Settings.Music))
			ctx.Text(fmt.Sprintf("Zoom Level: %.2fx", zoomLevel))

			// Mouse coordinates
			cursorX, cursorY := ebiten.CursorPosition()
			ctx.Text(fmt.Sprintf("Mouse: (%d, %d)", cursorX, cursorY))

			// Scene information
			ctx.Text(fmt.Sprintf("Scene: %s", sceneName))
			ctx.Text(fmt.Sprintf("State: %s", sceneState))

			// ItemRenderer information
			if itemRenderer != nil {
				ctx.Header("ItemRenderer Info", false, func() {
					ctx.Text(fmt.Sprintf("Item Count: %d", itemRenderer.GetItemCount()))

					itemIDs := make([]enums.Item, 0, len(itemRenderer.FrameDataMap))
					for itemID := range itemRenderer.FrameDataMap {
						itemIDs = append(itemIDs, itemID)
					}

					sort.Slice(itemIDs, func(i, j int) bool {
						return string(itemIDs[i]) < string(itemIDs[j])
					})

					for _, itemID := range itemIDs {
						ctx.Header(string(itemID), false, func() {
							if item, exists := itemRenderer.ItemMap[itemID]; exists {
								if item.Name != "" {
									ctx.Text(fmt.Sprintf("Name: %s", item.Name))
								}
								ctx.Text(fmt.Sprintf("X: %.0f", item.Pos.X))
								ctx.Text(fmt.Sprintf("Y: %.0f", item.Pos.Y))
								ctx.Text(fmt.Sprintf("Scale: %.2f", item.Scale))
							}
						})
					}
				})
			} else {
				ctx.Text("No ItemRenderer in current scene")
			}

			// Debug controls
			ctx.Header("Debug Controls", false, func() {
				ctx.Checkbox(&ShowLines, "Center Lines")
				ctx.Checkbox(&ShowItemSelection, "Item Selection")

				// Zoom controls
				ctx.Header("Zoom Controls", false, func() {
					ctx.Button("Zoom In (+0.1)").On(func() {
						newZoom := zoomLevel + 0.1
						if newZoom <= 3.0 {
							zoomLevel = newZoom
						}
					})
					ctx.Button("Zoom Out (-0.1)").On(func() {
						newZoom := zoomLevel - 0.1
						if newZoom >= 0.1 {
							zoomLevel = newZoom
						}
					})
					ctx.Button("Reset Zoom (1.0)").On(func() {
						zoomLevel = 1.0
					})
					ctx.Text("Zoom Range: 0.1x - 3.0x")
					ctx.Text("Keyboard: =/- to zoom, 0 to reset")
				})

				if sceneNavigator != nil {
					ctx.Header("Go to", false, func() {
						ctx.Button("Dummy").On(func() { sceneNavigator.NavigateTo(enums.SceneDummy) })
						ctx.Button("Splash").On(func() { sceneNavigator.NavigateTo(enums.SceneSplash) })
						ctx.Button("Main Menu").On(func() { sceneNavigator.NavigateTo(enums.SceneMainMenu) })
						ctx.Button("Storage Room").On(func() { sceneNavigator.NavigateTo(enums.SceneStorageRoom) })
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
		if len(CurrentDebugShaders) == 0 {
			ctx.Text("No shaders set for debugging")
			return
		}

		for _, shader := range CurrentDebugShaders {
			if uniforms, ok := shader.(*shaders.GraphicsQualityUniforms); ok {
				debugGraphicsQualityShader(ctx, uniforms)
			}

			if uniforms, ok := shader.(*shaders.SilentHillRedShaderUniforms); ok {
				debugSilentHillShader(ctx, uniforms)
			}
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

func debugGraphicsQualityShader(ctx *debugui.Context, uniforms *shaders.GraphicsQualityUniforms) {
	ctx.Header("Graphics Quality Shader", false, func() {
		qualityOptions := []string{"Low", "Medium", "High"}
		currentIndex := int(uniforms.Settings.Quality)

		ctx.Text("Graphics Quality:")
		ctx.Dropdown(&currentIndex, qualityOptions).On(func() {
			uniforms.Settings.Quality = conf.QualityLevel(currentIndex)
		})

		ctx.Text(fmt.Sprintf("Shader Value: %.1f", uniforms.Quality))

		ctx.Button("Reset to Defaults").On(func() {
			uniforms.ResetToInitial()
		})
	})
}
