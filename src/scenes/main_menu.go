package scenes

import (
	"fmt"
	"math"
	"remembering-home/src/assets"
	"remembering-home/src/assets/shaders"
	"remembering-home/src/conf"
	"remembering-home/src/context"
	"remembering-home/src/debug"
	"remembering-home/src/errs"
	"remembering-home/src/graphics"
	"remembering-home/src/ui"
	"remembering-home/src/utils"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/ebitick"
	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

const (
	MENU_START = iota
	MENU_SETTINGS
	MENU_QUIT
)

const (
	MENU_QUIT_CANCEL = iota
	MENU_QUIT_CONFIRM
)

const (
	BASE_X         = 32.0
	SIZE_X         = float32(128.0)
	GAP            = 64.0
	BANNER_PADDING = 20.0
	BANNER_HEIGHT  = 1.4
)

type Main_Menu_Scene struct {
	Context              *context.GameContext
	SceneManager         SceneManager
	TimerSys             *ebitick.TimerSystem
	Routine              *routine.Routine
	BannerSubtitle       *ui.BannerText
	TextQuit             *graphics.Text
	TextVersion          *graphics.Text
	AnimDesk             *graphics.AnimationPlayer
	AnimHallway          *graphics.AnimationPlayer
	LayerColorize        *graphics.Layer
	LayerText            *graphics.Layer
	AnimTitle            *graphics.AnimationPlayer
	CurrentStateName     string
	MenuMain             *ui.TextList
	MenuQuit             *ui.TextList
	MenuSettings         *ui.SettingsMenu
	SelectedIdx          int
	ShowMenuTexts        bool
	ShowQuitSubMenuTexts bool
	CanInteract          bool
	TitleGlitchCooldown  bool
}

func (scene *Main_Menu_Scene) GetName() string {
	return "Main Menu"
}

func (scene *Main_Menu_Scene) GetStateName() string {
	return scene.CurrentStateName
}

func NewMainMenuScene(ctx *context.GameContext, sceneManager SceneManager) *Main_Menu_Scene {
	scene := Main_Menu_Scene{
		Context:      ctx,
		SceneManager: sceneManager,
		TimerSys:     ebitick.NewTimerSystem(),
		LayerColorize: graphics.NewLayerWithShader(
			"colorize layer",
			conf.GAME_W,
			conf.GAME_H,
			ctx.Loader.LoadShader(assets.ShaderColorize).Data,
		),
		LayerText: graphics.NewLayerWithTriangleShader(
			"texts layer",
			conf.GAME_W,
			conf.GAME_H,
			ctx.Loader.LoadShader(assets.ShaderSilentHillRed).Data,
		),
	}

	scene.LayerColorize.DRSO.Uniforms = map[string]any{"Color": [4]float32{1, 1, 1, 1}}
	scene.LayerColorize.Disabled = true

	resFontJamboree18 := ctx.Loader.LoadFont(assets.FontJamboree18)
	keys := ctx.InputHandler.ActionKeyNames(context.ActionEnter, input.KeyboardDevice)
	if len(keys) == 0 {
		panic(fmt.Errorf("no valid '%d' in action key names", context.ActionEnter))
	}

	titleFrameW, titleFrameH := assets.SheetTitleFrameData.W, assets.SheetTitleFrameData.H
	scaleTitle := float64(min(conf.GAME_W/titleFrameW, conf.GAME_H/titleFrameH)) * 0.5
	scene.AnimTitle = graphics.NewAnimationPlayer(ctx.Loader.LoadImage(assets.ImageSheetTitle).Data)
	scene.AnimTitle.AddStateAnimation("row1", 0, 0, titleFrameW, titleFrameH, assets.SheetTitleFrameData.MaxCols, false)
	scene.AnimTitle.SetFPS(0)
	utils.DIOReplaceAlpha(scene.AnimTitle.DIO, 0)
	scene.AnimTitle.Update()
	scene.AnimTitle.DIO.GeoM.Scale(scaleTitle, scaleTitle)
	scene.AnimTitle.DIO.GeoM.Translate(
		BASE_X,
		conf.GAME_H/2-float64(titleFrameH)*scaleTitle/2,
	)

	txtSubtitle := graphics.NewText(&resFontJamboree18.Face, fmt.Sprintf("press <%s> to continue", keys[0]), true)
	txtSubtitle.SetPos(conf.GAME_W/2, conf.GAME_H-txtSubtitle.Face.Metrics().HAscent*2)
	txtSubtitle.SetAlign(text.AlignCenter, text.AlignCenter)
	utils.SetColor(txtSubtitle.DO, 1, 1, 1, 0)
	scene.BannerSubtitle = ui.NewBannerText(txtSubtitle, ui.BannerConfig{
		Padding:          BANNER_PADDING,
		HeightMultiplier: BANNER_HEIGHT,
		AutoWidth:        true,
	})
	scene.BannerSubtitle.SetVisible(false)

	versionText := graphics.NewText(&resFontJamboree18.Face, "version: "+conf.GAME_VERSION, true)
	versionText.SetPos(conf.GAME_W-BASE_X*0.3, conf.GAME_H-versionText.Face.Metrics().HAscent)
	versionText.SetAlign(text.AlignEnd, text.AlignCenter)
	utils.SetColor(versionText.DO, 1, 1, 1, 1)
	scene.TextVersion = versionText

	deskFrameW, deskFrameH := assets.SheetDeskFrameData.W, assets.SheetDeskFrameData.H
	scaleDesk := float64(min(conf.GAME_W/deskFrameW, conf.GAME_H/deskFrameH))
	scene.AnimDesk = graphics.NewAnimationPlayer(ctx.Loader.LoadImage(assets.ImageSheetDesk).Data)
	scene.AnimDesk.AddStateAnimation("row1", 0, 0, deskFrameW, deskFrameH, assets.SheetDeskFrameData.MaxCols, false)
	scene.AnimDesk.AddStateAnimation("row2", 0, 64, deskFrameW, deskFrameH, assets.SheetDeskFrameData.MaxCols, false)
	scene.AnimDesk.AddStateAnimation("row3", 0, 128, deskFrameW, deskFrameH, 1, false)
	scene.AnimDesk.AddStateAnimation("static", deskFrameW, 128, deskFrameW, deskFrameH, 1, false)
	scene.AnimDesk.SetFPS(7)
	utils.DIOReplaceAlpha(scene.AnimDesk.DIO, 1)
	scene.AnimDesk.Update()
	scene.AnimDesk.DIO.GeoM.Scale(scaleDesk, scaleDesk)
	scene.AnimDesk.DIO.GeoM.Translate(
		conf.GAME_W/2-float64(deskFrameW)*scaleDesk/2,
		conf.GAME_H/2-float64(deskFrameH)*scaleDesk/2,
	)

	hallwayFrameW, hallwayFrameH := assets.BGHallwayFrameData.W, assets.BGHallwayFrameData.H
	scaleHallway := float64(min(conf.GAME_W/hallwayFrameW, conf.GAME_H/hallwayFrameH))
	scene.AnimHallway = graphics.NewAnimationPlayer(ctx.Loader.LoadImage(assets.ImageBGHallway).Data)
	scene.AnimHallway.AddStateAnimation("row1", 0, 0, hallwayFrameW, hallwayFrameH, assets.BGHallwayFrameData.MaxCols, false)
	scene.AnimHallway.SetFPS(0)
	utils.DIOReplaceAlpha(scene.AnimHallway.DIO, 1)
	scene.AnimHallway.Update()
	scene.AnimHallway.DIO.GeoM.Scale(scaleHallway, scaleHallway)
	scene.AnimHallway.DIO.GeoM.Translate(
		conf.GAME_W/2-float64(hallwayFrameW)*scaleHallway/2,
		conf.GAME_H/2-float64(hallwayFrameH)*scaleHallway/2,
	)

	baseY := conf.GAME_H/2 + float64(deskFrameH)*scaleDesk/2 + 8.0
	resFontJamboree26 := ctx.Loader.LoadFont(assets.FontJamboree26)

	scene.MenuMain = ui.NewTextList(ui.TextListConfig{
		Layout:       ui.LayoutVertical,
		BaseX:        BASE_X,
		BaseY:        baseY,
		Gap:          0,
		Alignment:    text.AlignStart,
		SecAlignment: text.AlignStart,
		BannerConfig: ui.BannerConfig{
			Padding:          BANNER_PADDING,
			HeightMultiplier: BANNER_HEIGHT * 0.7,
			AutoWidth:        true,
		},
	})

	for _, txt := range []string{"Start", "Settings", "Quit"} {
		txtMenu := graphics.NewText(&resFontJamboree26.Face, txt, true)
		utils.SetColor(txtMenu.DO, 1, 1, 1, 1)
		scene.MenuMain.AddText(txtMenu)
	}

	scene.MenuMain.SetBannerWidth(scene.MenuMain.GetBannerWidth() * 1.5)

	txtQuit := graphics.NewText(&resFontJamboree18.Face, "Are you sure you want to quit?", true)
	txtQuit.SetPos(conf.GAME_W/2, conf.GAME_H-txtQuit.Face.Metrics().HAscent*4)
	txtQuit.SetAlign(text.AlignCenter, text.AlignCenter)
	utils.SetColor(txtQuit.DO, 1, 1, 1, 1)
	scene.TextQuit = txtQuit

	firstMenuText := scene.MenuMain.Texts[0]
	quitBaseY := conf.GAME_H - firstMenuText.Face.Metrics().HAscent
	scene.MenuQuit = ui.NewTextList(ui.TextListConfig{
		Layout:       ui.LayoutHorizontal,
		BaseX:        float64(conf.GAME_W/2) - GAP,
		BaseY:        quitBaseY,
		Gap:          GAP * 2,
		Alignment:    text.AlignEnd,
		SecAlignment: text.AlignCenter,
		BannerConfig: ui.BannerConfig{
			Padding:          BANNER_PADDING,
			HeightMultiplier: BANNER_HEIGHT,
			AutoWidth:        true,
		},
	})

	txtNo := graphics.NewText(&resFontJamboree26.Face, "No", true)
	utils.SetColor(txtNo.DO, 1, 1, 1, 1)
	scene.MenuQuit.AddText(txtNo)

	txtYes := graphics.NewText(&resFontJamboree26.Face, "Yes", true)
	txtYes.SetAlign(text.AlignStart, text.AlignCenter)
	utils.SetColor(txtYes.DO, 1, 1, 1, 1)
	scene.MenuQuit.AddText(txtYes)

	scene.MenuQuit.SetVisible(false)
	scene.MenuQuit.SetEnabled(false)

	scene.MenuSettings = ui.NewSettingsMenu(ui.SettingsMenuConfig{
		Context:        ctx,
		BaseX:          BASE_X,
		BaseY:          baseY,
		SubmenuOffsetX: float64(SIZE_X + GAP*3),
		Gap:            0,
		BannerConfig: ui.BannerConfig{
			Padding:          BANNER_PADDING,
			HeightMultiplier: BANNER_HEIGHT * 0.7,
			AutoWidth:        true,
		},
		PrimaryFace:   &resFontJamboree26.Face,
		SecondaryFace: &resFontJamboree18.Face,
	})
	scene.MenuSettings.SetVisible(false)
	scene.MenuSettings.SetEnabled(false)

	scene.LayerText.DTSO.Images[1] = ctx.Loader.LoadImage(assets.TexturePaper).Data
	scene.LayerText.Disabled = false
	scene.LayerText.Uniforms = shaders.NewSilentHillRedShaderUniforms(shaders.FadeStateHidden)

	uniforms := graphics.MustCastUniform[*shaders.SilentHillRedShaderUniforms](scene.LayerText)
	scene.BannerSubtitle.UpdateBannerPositionForce(uniforms)

	if conf.DEV {
		debug.SetDebugShader(scene.LayerText.Uniforms)
	}

	vx, ix := utils.AppendRectVerticesIndices(
		nil,
		nil,
		0,
		&utils.RectOpts{
			DstX:      0,
			DstY:      0,
			SrcX:      0,
			SrcY:      0,
			DstWidth:  conf.GAME_W,
			DstHeight: conf.GAME_H,
			SrcWidth:  conf.GAME_W,
			SrcHeight: conf.GAME_H,
			R:         1,
			G:         1,
			B:         1,
			A:         1,
		},
	)
	scene.LayerText.Vertices = vx
	scene.LayerText.Indices = ix

	uniform := graphics.MustCastUniform[*shaders.SilentHillRedShaderUniforms](scene.LayerText)

	scene.Routine = routine.New()
	scene.Routine.Define(
		"main menu scene",
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "title text animation"
			uniform.TriggerFadeInWithCallback(2, func() {
				scene.BannerSubtitle.SetVisible(true)
			})
			if scene.SceneManager.IsFadeInFinished() {
				return routine.FlowNext
			}
			return routine.FlowIdle
		}),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "waiting input"
			inputHandler := scene.Context.InputHandler

			if inputHandler.ActionIsJustPressed(context.ActionEnter) && !uniform.IsAnimating() {
				scene.CurrentStateName = "showing menu..."
				scene.AnimDesk.SetStateReset("row1")
				scene.AnimDesk.Update()
				scene.BannerSubtitle.SetVisible(false)
				scene.AnimTitle.SetStateReset("row1")
				utils.DIOReplaceAlpha(scene.AnimTitle.DIO, 1)
				scene.AnimTitle.Update()
				return routine.FlowNext
			}
			return routine.FlowIdle
		}),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "finished"
			scene.ShowMenuTexts = true

			scene.TimerSys.After(time.Second/2, func() {
				scene.CanInteract = true
			})

			waitFor := utils.IntRandRange(1, 3)
			scene.CurrentStateName = "waiting flicker... " + strconv.Itoa(waitFor)
			scene.TimerSys.After(time.Second*time.Duration(waitFor), func() {
				scene.RandomFlicker()
			})
			return routine.FlowNext
		}),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "in main menu"
			return routine.FlowIdle
		}),
	)
	scene.Routine.Run()
	return &scene
}

func (scene *Main_Menu_Scene) RandomFlicker() {
	waitFor := utils.IntRandRange(3, 6)
	scene.TimerSys.After(time.Second*time.Duration(waitFor), func() {
		scene.AnimDesk.Paused = false
		scene.AnimDesk.SetStateReset("row1")
		scene.AnimDesk.SetFPS(float64(utils.IntRandRange(4, 8)))
		scene.TitleGlitchCooldown = false
	})
}

func (scene *Main_Menu_Scene) TriggerTitleGlitch() {
	if scene.TitleGlitchCooldown {
		return
	}

	scene.TitleGlitchCooldown = true
	scene.AnimTitle.Paused = false
	scene.AnimTitle.SetStateReset("row1")
	scene.AnimTitle.SetFPS(float64(utils.IntRandRange(20, 35)))

	glitchDuration := time.Millisecond * time.Duration(utils.IntRandRange(100, 300))
	scene.TimerSys.After(glitchDuration, func() {
		scene.AnimTitle.PauseAtFrame(0)
		scene.AnimTitle.SetFPS(0)
	})
}

func (scene *Main_Menu_Scene) Update() error {
	scene.TimerSys.Update()
	scene.Routine.Update()

	if scene.SceneManager.IsFading() {
		return nil
	}

	scene.AnimDesk.Update()

	if scene.AnimDesk.State() == "row2" && scene.AnimDesk.IsInLastFrame() {
		scene.TriggerTitleGlitch()
	}

	if scene.AnimDesk.IsInLastFrame() {
		switch scene.AnimDesk.State() {
		case "row1":
			scene.AnimDesk.SetStateReset("row2")
		case "row2":
			scene.AnimDesk.SetStateReset("row3")
		case "row3":
			scene.AnimDesk.SetStateReset("row1")
			scene.AnimDesk.PauseAtFrame(0)
			scene.AnimDesk.SetFPS(0)
			scene.RandomFlicker()
		}
	}

	scene.AnimTitle.Update()
	if scene.AnimTitle.IsInLastFrame() && !scene.AnimTitle.Paused {
		scene.AnimTitle.SetStateReset("row1")
	}

	uniform := graphics.MustCastUniform[*shaders.SilentHillRedShaderUniforms](scene.LayerText)
	uniform.Time += 0.01
	v := (math.Sin(uniform.Time) + 1) / 2
	v = v*0.5 + 1.0
	uniform.GlowIntensity = float64(v)
	uniform.Update()

	textAlpha := uniform.GetTextAlpha()

	scene.BannerSubtitle.ApplyShaderAlpha(uniform)
	scene.MenuMain.ApplyShaderAlpha(uniform)
	scene.MenuQuit.ApplyShaderAlpha(uniform)
	scene.MenuSettings.ApplyShaderAlpha(uniform)
	utils.DOReplaceAlpha(scene.TextQuit.DO, textAlpha)
	utils.DOReplaceAlpha(scene.TextVersion.DO, textAlpha)

	if !scene.ShowMenuTexts {
		scene.BannerSubtitle.UpdateBannerPosition(uniform)
	}

	if scene.ShowMenuTexts {
		if scene.CanInteract {
			scene.MenuMain.Update(scene.Context.InputHandler)

			if scene.MenuMain.IsConfirmed() {
				switch scene.MenuMain.GetSelectedIndex() {
				case MENU_START: //TODO: (Brandon) - go to game
					uniform.TriggerFadeOutWithCallback(2, func() {
						scene.SceneManager.GoTo(NewIntroScene(scene.Context, scene.SceneManager))
					})
					return nil
				case MENU_SETTINGS:
					scene.SelectedIdx = MENU_SETTINGS
					scene.ShowMenuTexts = false
					scene.MenuSettings.SetVisible(true)
					scene.MenuSettings.SetEnabled(true)
					scene.MenuSettings.Reset()
					return nil
				case MENU_QUIT:
					scene.ShowMenuTexts = false
					scene.MenuQuit.SetCurrentIndex(MENU_QUIT_CANCEL)
					scene.MenuQuit.SetVisible(true)
					scene.MenuQuit.SetEnabled(true)
					scene.SelectedIdx = MENU_QUIT
					return nil
				default:
					panic(scene.MenuMain.GetSelectedIndex())
				}
			}
		}

		curText := scene.MenuMain.GetSelectedText()
		textH := curText.DO.LineSpacing
		bannerHeight := float64(textH) * BANNER_HEIGHT * 0.7

		_, posY, _, sizeY := graphics.CalculateBannerPosition(curText, scene.MenuMain.GetBannerWidth(), bannerHeight)
		fixedPosX := BASE_X - BANNER_PADDING*3

		uniform.SetBannerBounds(fixedPosX, posY, scene.MenuMain.GetBannerWidth(), sizeY)
	}

	if !scene.ShowMenuTexts && scene.SelectedIdx == MENU_QUIT {
		if scene.CanInteract {
			scene.MenuQuit.Update(scene.Context.InputHandler)

			if scene.MenuQuit.IsConfirmed() {
				switch scene.MenuQuit.GetSelectedIndex() {
				case MENU_QUIT_CANCEL:
					scene.ShowMenuTexts = true
					scene.MenuQuit.SetVisible(false)
					scene.MenuQuit.SetEnabled(false)
					scene.SelectedIdx = 0
					scene.MenuMain.SetCurrentIndex(MENU_QUIT)
					return nil
				case MENU_QUIT_CONFIRM:
					return errs.ErrQuit
				default:
					panic(scene.MenuQuit.GetSelectedIndex())
				}
			}
		}

		scene.MenuQuit.UpdateBannerPosition(uniform)
	}

	if !scene.ShowMenuTexts && scene.SelectedIdx == MENU_SETTINGS {
		if scene.CanInteract {
			scene.MenuSettings.Update(scene.Context.InputHandler)

			if scene.MenuSettings.IsBackRequested() {
				scene.ShowMenuTexts = true
				scene.MenuSettings.SetVisible(false)
				scene.MenuSettings.SetEnabled(false)
				scene.SelectedIdx = 0
				scene.MenuMain.SetCurrentIndex(MENU_SETTINGS)
				return nil
			}
		}

		if scene.MenuSettings.GetCurrentLevel() == ui.MenuLevelMain {
			mainList := scene.MenuSettings.GetMainList()
			curText := mainList.GetSelectedText()
			textH := curText.DO.LineSpacing
			bannerHeight := float64(textH) * BANNER_HEIGHT * 0.7

			_, posY, _, sizeY := graphics.CalculateBannerPosition(curText, mainList.GetBannerWidth(), bannerHeight)
			fixedPosX := BASE_X - BANNER_PADDING*3

			uniform.SetBannerBounds(fixedPosX, posY, mainList.GetBannerWidth(), sizeY)
		} else {
			submenu := scene.MenuSettings.GetActiveSubmenu()
			if submenu != nil {
				curText := submenu.GetSelectedText()
				textH := curText.DO.LineSpacing
				bannerHeight := float64(textH) * BANNER_HEIGHT * 0.7

				_, posY, _, sizeY := graphics.CalculateBannerPosition(curText, submenu.GetBannerWidth(), bannerHeight)
				fixedPosX := scene.MenuSettings.GetMainList().Texts[0].X + float64(SIZE_X+GAP*2) - BANNER_PADDING*3

				uniform.SetBannerBounds(fixedPosX, posY, submenu.GetBannerWidth(), sizeY)
			}
		}
	}

	scene.LayerText.ApplyTransformation()

	return nil
}

func (scene *Main_Menu_Scene) Draw(screen *ebiten.Image) {
	canvas := scene.LayerColorize.Canvas
	canvas.Clear()

	canvas2 := scene.LayerText.Canvas
	canvas2.Clear()

	uniform := graphics.MustCastUniform[*shaders.SilentHillRedShaderUniforms](scene.LayerText)

	uniform.ToShaders(scene.LayerText.DTSO)

	switch scene.SelectedIdx {
	case MENU_SETTINGS:
		canvas.DrawImage(scene.AnimHallway.CurrentFrame, scene.AnimHallway.DIO)
		scene.MenuSettings.Draw(canvas2)
	case MENU_QUIT:
		canvas.DrawImage(scene.AnimDesk.CurrentFrame, scene.AnimDesk.DIO)
		canvas.DrawImage(scene.AnimTitle.CurrentFrame, scene.AnimTitle.DIO)
		scene.MenuQuit.Draw(canvas2)
		scene.TextQuit.Draw(canvas2)
	default:
		canvas.DrawImage(scene.AnimDesk.CurrentFrame, scene.AnimDesk.DIO)
		canvas.DrawImage(scene.AnimTitle.CurrentFrame, scene.AnimTitle.DIO)
		scene.BannerSubtitle.Draw(canvas2)
		if scene.ShowMenuTexts {
			scene.TextVersion.Draw(canvas)
			scene.MenuMain.Draw(canvas2)
		}
	}

	scene.LayerText.RenderWithShader(canvas)
	scene.LayerColorize.RenderWithShader(screen)
}

var _ Scene = (*Main_Menu_Scene)(nil)
