package game

import (
	"fmt"
	"math"
	"remembering-home/src/assets"
	"remembering-home/src/common"
	"remembering-home/src/conf"
	"remembering-home/src/effects"
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
	BASE_X = 32.0
	SIZE_X = float32(128.0)
)

type Main_Menu_Scene struct {
	GameState            *Game_State
	TimerSys             *ebitick.TimerSystem
	CurrentStateName     string
	Routine              *routine.Routine
	TextSubtitle         *Text
	TextQuit             *Text
	FaderSubtitle        *effects.Fader
	ShowMenuTexts        bool
	ShowQuitSubMenuTexts bool
	CanInteract          bool
	AnimDesk             *AnimationPlayer
	AnimHallway          *AnimationPlayer
	AnimTitle            *AnimationPlayer
	Texts                []*Text
	QuitTexts            []*Text
	CurrentIdx           int
	CurrentQuitIdx       int
	SelectedIdx          int
	LayerColorize        *Layer
	LayerText            *Layer
}

func (scene *Main_Menu_Scene) GetName() string {
	return "Main Menu"
}

func (scene *Main_Menu_Scene) GetStateName() string {
	return scene.CurrentStateName
}

func NewMainMenuScene(gameState *Game_State) *Main_Menu_Scene {
	scene := Main_Menu_Scene{
		GameState: gameState,
		TimerSys:  ebitick.NewTimerSystem(),
		Texts:     make([]*Text, 0, 3),
		QuitTexts: make([]*Text, 0, 2),
		LayerColorize: NewLayerWithShader(
			"colorize layer",
			conf.GAME_W,
			conf.GAME_H,
			gameState.Loader.LoadShader(assets.ShaderColorize).Data,
		),
		LayerText: NewLayerWithTriangleShader(
			"menu texts layer",
			conf.GAME_W,
			conf.GAME_H,
			gameState.Loader.LoadShader(assets.ShaderMenuText).Data,
		),
	}

	scene.LayerColorize.DRSO.Uniforms = map[string]any{"Color": [4]float32{1, 1, 1, 1}}
	scene.LayerColorize.Disabled = true

	resFontJamboree18 := gameState.Loader.LoadFont(assets.FontJamboree18)
	keys := scene.GameState.InputHandler.ActionKeyNames(ActionEnter, input.KeyboardDevice)
	if len(keys) == 0 {
		panic(fmt.Sprintf("No valid '%d' in action key names", ActionEnter))
	}

	titleFrameW, titleFrameH := assets.SheetTitleFrameData.W, assets.SheetTitleFrameData.H
	scaleTitle := float64(min(conf.GAME_W/titleFrameW, conf.GAME_H/titleFrameH)) * 0.5
	scene.AnimTitle = NewAnimationPlayer(gameState.Loader.LoadImage(assets.ImageSheetTitle).Data)
	scene.AnimTitle.AddStateAnimation("row1", 0, 0, titleFrameW, titleFrameH, assets.SheetTitleFrameData.MaxCols, false)
	scene.AnimTitle.SetFPS(0)
	utils.DIOReplaceAlpha(scene.AnimTitle.DIO, 0)
	scene.AnimTitle.Update()
	scene.AnimTitle.DIO.GeoM.Scale(scaleTitle, scaleTitle)
	scene.AnimTitle.DIO.GeoM.Translate(
		BASE_X,
		conf.GAME_H/2-float64(titleFrameH)*scaleTitle/2,
	)

	subtitleTxt := NewText(&resFontJamboree18.Face, fmt.Sprintf("press <%s> to continue", keys[0]), true)
	subtitleTxt.SetPos(conf.GAME_W/2, conf.GAME_H-subtitleTxt.Face.Metrics().HAscent*2)
	subtitleTxt.SetAlign(text.AlignCenter, text.AlignCenter)
	utils.SetColor(subtitleTxt.DO, 1, 1, 1, 1)
	scene.TextSubtitle = subtitleTxt
	scene.FaderSubtitle = effects.NewFader(0, 1, 1)
	scene.FaderSubtitle.Stopped = true

	deskFrameW, deskFrameH := assets.SheetDeskFrameData.W, assets.SheetDeskFrameData.H
	scaleDesk := float64(min(conf.GAME_W/deskFrameW, conf.GAME_H/deskFrameH))
	scene.AnimDesk = NewAnimationPlayer(gameState.Loader.LoadImage(assets.ImageSheetDesk).Data)
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
	scene.AnimHallway = NewAnimationPlayer(gameState.Loader.LoadImage(assets.ImageBGHallway).Data)
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
	resFontJamboree26 := gameState.Loader.LoadFont(assets.FontJamboree26)

	for _, txt := range []string{"Start", "Settings", "Quit"} {
		newTxt := NewText(&resFontJamboree26.Face, txt, true)
		newTxt.SetPos(BASE_X, baseY)
		newTxt.SetAlign(text.AlignStart, text.AlignStart)
		utils.SetColor(newTxt.DO, 1, 1, 1, 1)
		scene.Texts = append(scene.Texts, newTxt)
		baseY += newTxt.DO.LineSpacing
	}

	quitTxt := NewText(&resFontJamboree18.Face, "Are you sure you want to quit?", true)
	quitTxt.SetPos(conf.GAME_W/2, conf.GAME_H-quitTxt.Face.Metrics().HAscent*4)
	quitTxt.SetAlign(text.AlignCenter, text.AlignCenter)
	utils.SetColor(quitTxt.DO, 1, 1, 1, 1)
	scene.TextQuit = quitTxt

	const GAP = 64.0
	{
		newTxt := NewText(&resFontJamboree26.Face, "No", true)
		newTxt.SetPos(
			float64(conf.GAME_W/2)-GAP,
			conf.GAME_H-newTxt.Face.Metrics().HAscent,
		)
		newTxt.SetAlign(text.AlignEnd, text.AlignCenter)
		utils.SetColor(newTxt.DO, 1, 1, 1, 1)
		scene.QuitTexts = append(scene.QuitTexts, newTxt)
	}

	{
		newTxt := NewText(&resFontJamboree26.Face, "Yes", true)
		newTxt.SetPos(
			float64(conf.GAME_W/2)+GAP,
			conf.GAME_H-newTxt.Face.Metrics().HAscent,
		)
		newTxt.SetAlign(text.AlignStart, text.AlignCenter)
		utils.SetColor(newTxt.DO, 1, 1, 1, 1)
		scene.QuitTexts = append(scene.QuitTexts, newTxt)
	}

	txt0 := scene.Texts[0]
	textH := scene.Texts[0].DO.LineSpacing
	scene.LayerText.DTSO.Images[1] = gameState.Loader.LoadImage(assets.TexturePaper).Data
	scene.LayerText.Disabled = true
	scene.LayerText.Uniforms = &assets.MenuTextShaderUniforms{
		Time:              0,
		Pos:               [2]float32{float32(txt0.X), float32(txt0.Y)},
		Size:              [2]float32{conf.GAME_W*0.1 + float32(len(txt0.Txt))*8, float32(textH)},
		StartingAmplitude: 0.5,
		StartingFreq:      1.0,
		Shift:             0.0,
		WhiteCutoff:       0.999,
		Velocity:          [2]float32{-8.0, -32.0},
		Color:             [4]float32{0.0, 1.0, 1.0, 1.0},
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
			if scene.GameState.SceneManager.IsFadeInFinished() {
				scene.FaderSubtitle.Stopped = false
				return routine.FlowNext
			}
			return routine.FlowIdle
		}),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "waiting input"
			inputHandler := scene.GameState.InputHandler
			if inputHandler.ActionIsJustPressed(ActionEnter) {
				scene.CurrentStateName = "showing menu..."
				scene.FaderSubtitle.Alpha = 0
				scene.FaderSubtitle.Stopped = true

				scene.AnimDesk.SetStateReset("row1")
				scene.AnimDesk.Update()

				utils.DOReplaceAlpha(scene.TextSubtitle.DO, 1)

				scene.AnimTitle.SetStateReset("row1")
				utils.DIOReplaceAlpha(scene.AnimTitle.DIO, 1)
				scene.AnimTitle.Update()
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
				scene.RandomTitleFrame()
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
	})
}

func (scene *Main_Menu_Scene) RandomTitleFrame() {
	waitFor := utils.IntRandRange(3, 6)
	scene.TimerSys.After(time.Second*time.Duration(waitFor), func() {
		scene.AnimTitle.Paused = false
		scene.AnimTitle.SetStateReset("row1")
		scene.AnimTitle.SetFPS(float64(utils.IntRandRange(1, 2)))
	})
}

func (scene *Main_Menu_Scene) Update() error {
	scene.TimerSys.Update()
	scene.Routine.Update()

	scene.FaderSubtitle.Update()
	cs := scene.FaderSubtitle.GetCS()
	scene.TextSubtitle.DO.ColorScale = *cs

	if scene.GameState.SceneManager.IsFading() {
		return nil
	}

	if scene.AnimTitle.DIO.ColorScale.A() < 1.0 {
		scene.AnimTitle.DIO.ColorScale = *cs
	}

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

	scene.AnimTitle.Update()
	if scene.AnimTitle.IsInLastFrame() {
		scene.AnimTitle.SetStateReset("row1")
		scene.AnimTitle.PauseAtFrame(0)
		scene.AnimTitle.SetFPS(0)
		scene.RandomTitleFrame()
	}

	uniform, ok := scene.LayerText.Uniforms.(*assets.MenuTextShaderUniforms)
	if !ok {
		panic("incorrect casting")
	}
	uniform.Time += 0.01
	v := (math.Sin(uniform.Time) + 1) / 2
	v = utils.ClampFloat64(v, 0.4, 0.6)
	uniform.StartingAmplitude = float32(v)
	const SPEED = 4
	uniform.Velocity[0] = float32(math.Sin(uniform.Time) * SPEED)
	uniform.Velocity[1] = float32(math.Cos(uniform.Time) * SPEED)

	if !scene.ShowMenuTexts && scene.TextSubtitle.GetAlpha() >= 0.9 {
		scene.LayerText.Disabled = false
		th := scene.TextSubtitle.Face.Metrics().HAscent
		uniform.Pos[0] = float32(scene.TextSubtitle.X)
		uniform.Pos[1] = float32(scene.TextSubtitle.Y - th/2 - 4)
		uniform.Size[0] = conf.GAME_W*0.1 + float32(len(scene.TextSubtitle.Txt))*1.5
		uniform.Size[1] = float32(th) * 1.2
	}

	if scene.ShowMenuTexts {
		if scene.CanInteract {
			inputHandler := scene.GameState.InputHandler
			if inputHandler.ActionIsJustReleased(ActionMoveUp) {
				scene.CurrentIdx--
			} else if inputHandler.ActionIsJustReleased(ActionMoveDown) {
				scene.CurrentIdx++
			} else if inputHandler.ActionIsJustReleased(ActionEnter) {
				switch scene.CurrentIdx {
				case MENU_START: //TODO: (Brandon) - go to game
					scene.GameState.SceneManager.GoTo(NewSplashScene(scene.GameState))
					return nil
				case MENU_SETTINGS:
					scene.SelectedIdx = MENU_SETTINGS
					return common.ERR_NOT_YET_IMPL
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

		scene.CurrentIdx = utils.ClampInt(scene.CurrentIdx, 0, len(scene.Texts)-1)
		curText := scene.Texts[scene.CurrentIdx]
		th := curText.Face.Metrics().HAscent
		uniform.Pos[0] = BASE_X*0.3 + SIZE_X
		uniform.Pos[1] = float32(curText.Y)
		uniform.Size[0] = SIZE_X
		uniform.Size[1] = float32(th) * 1.2
	}

	if !scene.ShowMenuTexts && scene.CurrentIdx == MENU_QUIT {
		if scene.CanInteract {
			inputHandler := scene.GameState.InputHandler
			if inputHandler.ActionIsJustReleased(ActionMoveLeft) {
				scene.CurrentQuitIdx = MENU_QUIT_CANCEL
			} else if inputHandler.ActionIsJustReleased(ActionMoveRight) {
				scene.CurrentQuitIdx = MENU_QUIT_CONFIRM
			} else if inputHandler.ActionIsJustReleased(ActionEnter) {
				switch scene.CurrentQuitIdx {
				case MENU_QUIT_CANCEL:
					scene.ShowMenuTexts = true
					scene.CurrentQuitIdx = MENU_QUIT_CANCEL
					scene.SelectedIdx = 0
					scene.CurrentIdx = MENU_QUIT
					return nil
				case MENU_QUIT_CONFIRM:
					return common.ERR_QUIT
				default:
					panic(scene.CurrentQuitIdx)
				}
			}
		}

		sizex := conf.GAME_W * 0.05
		curText := scene.QuitTexts[scene.CurrentQuitIdx]
		th := curText.Face.Metrics().HAscent
		if curText.DO.PrimaryAlign == text.AlignEnd {
			uniform.Pos[0] = float32(curText.X) - float32(curText.Face.Metrics().HAscent)/2
		} else if curText.DO.PrimaryAlign == text.AlignStart {
			uniform.Pos[0] = float32(curText.X) + float32(curText.Face.Metrics().HAscent)/2
		}
		uniform.Pos[1] = float32(curText.Y) - float32(curText.Face.Metrics().HAscent)/1.5
		uniform.Size[0] = float32(sizex)
		uniform.Size[1] = float32(th) * 1.2
	}

	scene.LayerText.ApplyTransformation()

	return nil
}

func (scene *Main_Menu_Scene) Draw(screen *ebiten.Image) {
	canvas := scene.LayerColorize.Canvas
	canvas.Clear()

	canvas2 := scene.LayerText.Canvas
	canvas2.Clear()

	uniform, ok := scene.LayerText.Uniforms.(*assets.MenuTextShaderUniforms)
	if !ok {
		panic("incorrect casting")
	}

	uniform.ToShaders(scene.LayerText.DTSO)

	switch scene.SelectedIdx {
	case MENU_SETTINGS:
		canvas.DrawImage(scene.AnimHallway.CurrentFrame, scene.AnimHallway.DIO)
	case MENU_QUIT:
		canvas.DrawImage(scene.AnimDesk.CurrentFrame, scene.AnimDesk.DIO)
		canvas.DrawImage(scene.AnimTitle.CurrentFrame, scene.AnimTitle.DIO)
		for _, txt := range scene.QuitTexts {
			utils.SetColor(txt.DO, 1, 1, 1, 1)
			txt.Draw(canvas2)
		}
		scene.TextQuit.Draw(canvas2)
	default:
		canvas.DrawImage(scene.AnimDesk.CurrentFrame, scene.AnimDesk.DIO)
		canvas.DrawImage(scene.AnimTitle.CurrentFrame, scene.AnimTitle.DIO)
		scene.TextSubtitle.Draw(canvas2)
		if scene.ShowMenuTexts {
			for _, txt := range scene.Texts {
				utils.SetColor(txt.DO, 1, 1, 1, 1)
				txt.Draw(canvas2)
			}
		}
	}

	scene.LayerText.RenderWithShader(canvas)
	scene.LayerColorize.RenderWithShader(screen)
}

var _ Scene = (*Main_Menu_Scene)(nil)
