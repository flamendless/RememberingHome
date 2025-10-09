package scenes

import (
	"fmt"
	"math"
	"remembering-home/src/assets"
	"remembering-home/src/assets/shaders"
	"remembering-home/src/conf"
	"remembering-home/src/context"
	"remembering-home/src/debug"
	"remembering-home/src/effects"
	"remembering-home/src/errs"
	"remembering-home/src/graphics"
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
)

type Main_Menu_Scene struct {
	Context       *context.GameContext
	SceneManager  SceneManager
	TimerSys      *ebitick.TimerSystem
	Routine       *routine.Routine
	TextSubtitle  *graphics.Text
	TextQuit      *graphics.Text
	TextVersion   *graphics.Text
	FaderSubtitle *effects.Fader
	AnimDesk      *graphics.AnimationPlayer
	AnimHallway   *graphics.AnimationPlayer
	LayerColorize *graphics.Layer
	LayerText     *graphics.Layer
	// AnimTitle            *graphics.AnimationPlayer
	CurrentStateName         string
	TextsMenu                []*graphics.Text
	TextsQuit                []*graphics.Text
	CurrentIdx               int
	CurrentQuitIdx           int
	SelectedIdx              int
	ShowMenuTexts            bool
	ShowQuitSubMenuTexts     bool
	CanInteract              bool
	FixedBannerWidth         float64
	FixedQuitBannerWidth     float64
	FixedSubtitleBannerWidth float64
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
		TextsMenu:    make([]*graphics.Text, 0, 3),
		TextsQuit:    make([]*graphics.Text, 0, 2),
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

	// titleFrameW, titleFrameH := assets.SheetTitleFrameData.W, assets.SheetTitleFrameData.H
	// scaleTitle := float64(min(conf.GAME_W/titleFrameW, conf.GAME_H/titleFrameH)) * 0.5
	// scene.AnimTitle = graphics.NewAnimationPlayer(ctx.Loader.LoadImage(assets.ImageSheetTitle).Data)
	// scene.AnimTitle.AddStateAnimation("row1", 0, 0, titleFrameW, titleFrameH, assets.SheetTitleFrameData.MaxCols, false)
	// scene.AnimTitle.SetFPS(0)
	// utils.DIOReplaceAlpha(scene.AnimTitle.DIO, 0)
	// scene.AnimTitle.Update()
	// scene.AnimTitle.DIO.GeoM.Scale(scaleTitle, scaleTitle)
	// scene.AnimTitle.DIO.GeoM.Translate(
	// 	BASE_X,
	// 	conf.GAME_H/2-float64(titleFrameH)*scaleTitle/2,
	// )

	txtSubtitle := graphics.NewText(&resFontJamboree18.Face, fmt.Sprintf("press <%s> to continue", keys[0]), true)
	txtSubtitle.SetPos(conf.GAME_W/2, conf.GAME_H-txtSubtitle.Face.Metrics().HAscent*2)
	txtSubtitle.SetAlign(text.AlignCenter, text.AlignCenter)
	utils.SetColor(txtSubtitle.DO, 1, 1, 1, 1)
	scene.TextSubtitle = txtSubtitle
	scene.FaderSubtitle = effects.NewFader(0, 1, 1)
	scene.FaderSubtitle.Stopped = true

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

	for _, txt := range []string{"Start", "Settings", "Quit"} {
		txtMenu := graphics.NewText(&resFontJamboree26.Face, txt, true)
		txtMenu.SetPos(BASE_X, baseY)
		txtMenu.SetAlign(text.AlignStart, text.AlignStart)
		utils.SetColor(txtMenu.DO, 1, 1, 1, 1)
		scene.TextsMenu = append(scene.TextsMenu, txtMenu)
		baseY += txtMenu.DO.LineSpacing
	}

	txtQuit := graphics.NewText(&resFontJamboree18.Face, "Are you sure you want to quit?", true)
	txtQuit.SetPos(conf.GAME_W/2, conf.GAME_H-txtQuit.Face.Metrics().HAscent*4)
	txtQuit.SetAlign(text.AlignCenter, text.AlignCenter)
	utils.SetColor(txtQuit.DO, 1, 1, 1, 1)
	scene.TextQuit = txtQuit

	{
		txtNo := graphics.NewText(&resFontJamboree26.Face, "No", true)
		txtNo.SetPos(
			float64(conf.GAME_W/2)-GAP,
			conf.GAME_H-txtNo.Face.Metrics().HAscent,
		)
		txtNo.SetAlign(text.AlignEnd, text.AlignCenter)
		utils.SetColor(txtNo.DO, 1, 1, 1, 1)
		scene.TextsQuit = append(scene.TextsQuit, txtNo)
	}

	{
		txtYes := graphics.NewText(&resFontJamboree26.Face, "Yes", true)
		txtYes.SetPos(
			float64(conf.GAME_W/2)+GAP,
			conf.GAME_H-txtYes.Face.Metrics().HAscent,
		)
		txtYes.SetAlign(text.AlignStart, text.AlignCenter)
		utils.SetColor(txtYes.DO, 1, 1, 1, 1)
		scene.TextsQuit = append(scene.TextsQuit, txtYes)
	}

	scene.calculateFixedBannerWidth()

	txt0 := scene.TextsMenu[0]
	textH := scene.TextsMenu[0].DO.LineSpacing
	scene.LayerText.DTSO.Images[1] = ctx.Loader.LoadImage(assets.TexturePaper).Data
	scene.LayerText.Disabled = true
	scene.LayerText.Uniforms = &shaders.SilentHillRedShaderUniforms{
		Time:           0,
		BannerPos:      [2]float64{txt0.X - BANNER_PADDING, txt0.Y},
		BannerSize:     [2]float64{scene.FixedBannerWidth, textH * 1.2},
		BaseRedColor:   [4]float64{1.0, 0.2, 0.2, 1.0},
		GlowIntensity:  1.0,
		MetallicShine:  0.4,
		EdgeDarkness:   0.3,
		TextGlowRadius: 2.5,
		NoiseScale:     0.5,
		NoiseIntensity: 0.5,
	}

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

	scene.Routine = routine.New()
	scene.Routine.Define(
		"main menu scene",
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "title text animation"
			if scene.SceneManager.IsFadeInFinished() {
				scene.FaderSubtitle.Stopped = false
				return routine.FlowNext
			}
			return routine.FlowIdle
		}),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "waiting input"
			inputHandler := scene.Context.InputHandler
			if inputHandler.ActionIsJustPressed(context.ActionEnter) {
				scene.CurrentStateName = "showing menu..."
				scene.FaderSubtitle.Alpha = 0
				scene.FaderSubtitle.Stopped = true

				scene.AnimDesk.SetStateReset("row1")
				scene.AnimDesk.Update()

				utils.DOReplaceAlpha(scene.TextSubtitle.DO, 1)

				// scene.AnimTitle.SetStateReset("row1")
				// utils.DIOReplaceAlpha(scene.AnimTitle.DIO, 1)
				// scene.AnimTitle.Update()
				scene.LayerText.Disabled = false

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
				// scene.RandomTitleFrame()
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

func (scene *Main_Menu_Scene) calculateFixedBannerWidth() {
	menuMaxWidth := float64(0)
	for _, txt := range scene.TextsMenu {
		width := float64(len(txt.Txt)) * 20.0
		if width > menuMaxWidth {
			menuMaxWidth = width
		}
	}
	scene.FixedBannerWidth = menuMaxWidth + BANNER_PADDING*3

	quitMaxWidth := float64(0)
	for _, txt := range scene.TextsQuit {
		width := float64(len(txt.Txt)) * 20.0
		if width > quitMaxWidth {
			quitMaxWidth = width
		}
	}
	scene.FixedQuitBannerWidth = quitMaxWidth + BANNER_PADDING*2

	subtitleWidth := float64(len(scene.TextSubtitle.Txt)) * 15.0
	scene.FixedSubtitleBannerWidth = subtitleWidth + BANNER_PADDING*3
}

func (scene *Main_Menu_Scene) RandomFlicker() {
	waitFor := utils.IntRandRange(3, 6)
	scene.TimerSys.After(time.Second*time.Duration(waitFor), func() {
		scene.AnimDesk.Paused = false
		scene.AnimDesk.SetStateReset("row1")
		scene.AnimDesk.SetFPS(float64(utils.IntRandRange(4, 8)))
	})
}

// func (scene *Main_Menu_Scene) RandomTitleFrame() {
// 	waitFor := utils.IntRandRange(3, 6)
// 	scene.TimerSys.After(time.Second*time.Duration(waitFor), func() {
// 		scene.AnimTitle.Paused = false
// 		scene.AnimTitle.SetStateReset("row1")
// 		scene.AnimTitle.SetFPS(float64(utils.IntRandRange(1, 2)))
// 	})
// }

func (scene *Main_Menu_Scene) Update() error {
	scene.TimerSys.Update()
	scene.Routine.Update()

	scene.FaderSubtitle.Update()
	cs := scene.FaderSubtitle.GetCS()
	scene.TextSubtitle.DO.ColorScale = *cs

	if scene.SceneManager.IsFading() {
		return nil
	}

	// if scene.AnimTitle.DIO.ColorScale.A() < 1.0 {
	// 	scene.AnimTitle.DIO.ColorScale = *cs
	// }

	scene.AnimDesk.Update()
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

	// scene.AnimTitle.Update()
	// if scene.AnimTitle.IsInLastFrame() {
	// 	scene.AnimTitle.SetStateReset("row1")
	// 	scene.AnimTitle.PauseAtFrame(0)
	// 	scene.AnimTitle.SetFPS(0)
	// 	scene.RandomTitleFrame()
	// }

	uniform, ok := scene.LayerText.Uniforms.(*shaders.SilentHillRedShaderUniforms)
	if !ok {
		panic("incorrect casting")
	}
	uniform.Time += 0.01
	v := (math.Sin(uniform.Time) + 1) / 2
	v = utils.ClampFloat64(v, 0.4, 0.6)
	uniform.GlowIntensity = float64(v)

	if !scene.ShowMenuTexts && scene.TextSubtitle.GetAlpha() >= 0.9 {
		scene.LayerText.Disabled = false
		th := scene.TextSubtitle.Face.Metrics().HAscent
		uniform.BannerPos[0] = float64(scene.TextSubtitle.X) - scene.FixedSubtitleBannerWidth/2
		uniform.BannerPos[1] = float64(scene.TextSubtitle.Y - th/2 - 4)
		uniform.BannerSize[0] = scene.FixedSubtitleBannerWidth
		uniform.BannerSize[1] = float64(th) * 1.2
	}

	if scene.ShowMenuTexts {
		if scene.CanInteract {
			inputHandler := scene.Context.InputHandler
			switch {
			case inputHandler.ActionIsJustReleased(context.ActionMoveUp):
				scene.CurrentIdx--
			case inputHandler.ActionIsJustReleased(context.ActionMoveDown):
				scene.CurrentIdx++
			case inputHandler.ActionIsJustReleased(context.ActionEnter):
				switch scene.CurrentIdx {
				case MENU_START: //TODO: (Brandon) - go to game
					scene.SceneManager.GoTo(NewIntroScene(scene.Context, scene.SceneManager))
					return nil
				case MENU_SETTINGS:
					scene.SelectedIdx = MENU_SETTINGS
					return errs.ErrNotYetImpl
				case MENU_QUIT:
					scene.ShowMenuTexts = false
					scene.CurrentQuitIdx = MENU_QUIT_CANCEL
					scene.SelectedIdx = MENU_QUIT
					return nil
				default:
					panic(scene.CurrentIdx)
				}
			}
		}

		scene.CurrentIdx = utils.ClampInt(scene.CurrentIdx, 0, len(scene.TextsMenu)-1)
		curText := scene.TextsMenu[scene.CurrentIdx]
		th := curText.Face.Metrics().HAscent
		uniform.BannerPos[0] = float64(curText.X) - BANNER_PADDING*2
		uniform.BannerPos[1] = float64(curText.Y)
		uniform.BannerSize[0] = scene.FixedBannerWidth
		uniform.BannerSize[1] = float64(th) * 1.2
	}

	if !scene.ShowMenuTexts && scene.CurrentIdx == MENU_QUIT {
		if scene.CanInteract {
			inputHandler := scene.Context.InputHandler
			switch {
			case inputHandler.ActionIsJustReleased(context.ActionMoveLeft):
				scene.CurrentQuitIdx = MENU_QUIT_CANCEL
			case inputHandler.ActionIsJustReleased(context.ActionMoveRight):
				scene.CurrentQuitIdx = MENU_QUIT_CONFIRM
			case inputHandler.ActionIsJustReleased(context.ActionEnter):
				switch scene.CurrentQuitIdx {
				case MENU_QUIT_CANCEL:
					scene.ShowMenuTexts = true
					scene.CurrentQuitIdx = MENU_QUIT_CANCEL
					scene.SelectedIdx = 0
					scene.CurrentIdx = MENU_QUIT
					return nil
				case MENU_QUIT_CONFIRM:
					return errs.ErrQuit
				default:
					panic(scene.CurrentQuitIdx)
				}
			}
		}

		curText := scene.TextsQuit[scene.CurrentQuitIdx]
		th := curText.Face.Metrics().HAscent
		uniform.BannerPos[0] = float64(curText.X) - scene.FixedQuitBannerWidth/2
		uniform.BannerPos[1] = float64(curText.Y) - curText.Face.Metrics().HAscent/1.5
		uniform.BannerSize[0] = scene.FixedQuitBannerWidth
		uniform.BannerSize[1] = float64(th) * 1.2
	}

	scene.LayerText.ApplyTransformation()

	return nil
}

func (scene *Main_Menu_Scene) Draw(screen *ebiten.Image) {
	canvas := scene.LayerColorize.Canvas
	canvas.Clear()

	canvas2 := scene.LayerText.Canvas
	canvas2.Clear()

	uniform, ok := scene.LayerText.Uniforms.(*shaders.SilentHillRedShaderUniforms)
	if !ok {
		panic("incorrect casting")
	}

	uniform.ToShaders(scene.LayerText.DTSO)

	switch scene.SelectedIdx {
	case MENU_SETTINGS:
		canvas.DrawImage(scene.AnimHallway.CurrentFrame, scene.AnimHallway.DIO)
	case MENU_QUIT:
		canvas.DrawImage(scene.AnimDesk.CurrentFrame, scene.AnimDesk.DIO)
		// canvas.DrawImage(scene.AnimTitle.CurrentFrame, scene.AnimTitle.DIO)
		for _, txt := range scene.TextsQuit {
			utils.SetColor(txt.DO, 1, 1, 1, 1)
			txt.Draw(canvas2)
		}
		scene.TextQuit.Draw(canvas2)
	default:
		canvas.DrawImage(scene.AnimDesk.CurrentFrame, scene.AnimDesk.DIO)
		// canvas.DrawImage(scene.AnimTitle.CurrentFrame, scene.AnimTitle.DIO)
		scene.TextSubtitle.Draw(canvas2)
		if scene.ShowMenuTexts {
			scene.TextVersion.Draw(canvas)
			for _, txt := range scene.TextsMenu {
				utils.SetColor(txt.DO, 1, 1, 1, 1)
				txt.Draw(canvas2)
			}
		}
	}

	scene.LayerText.RenderWithShader(canvas)
	scene.LayerColorize.RenderWithShader(screen)
}

var _ Scene = (*Main_Menu_Scene)(nil)
