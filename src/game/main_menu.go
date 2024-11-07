package game

import (
	"fmt"
	"math"
	"math/rand/v2"
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

type Main_Menu_Scene struct {
	GameState        *Game_State
	TimerSys         *ebitick.TimerSystem
	CurrentStateName string
	Routine          *routine.Routine
	TextTitle        *Text
	TextSubtitle     *Text
	FaderSubtitle    *effects.Fader
	ShowMenuTexts    bool
	CanInteract      bool
	AnimDesk         *AnimationPlayer
	AnimHallway      *AnimationPlayer
	Flickering       bool
	Texts            []*Text
	CurrentIdx       int
	SelectedIdx      int
	LayerColorize    *Layer
	LayerText        *Layer
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
		Texts:     make([]*Text, 0, 16),
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

	scene.LayerColorize.DRSO.Uniforms = map[string]any{
		"Color": [4]float32{1, 1, 1, 1},
	}
	scene.LayerColorize.Disabled = true

	resFontJamboree46 := gameState.Loader.LoadFont(assets.FontJamboree46)
	titleTxt := NewText(&resFontJamboree46.Face, "Remembering Home", true)
	titleTxt.SetPos(conf.GAME_W/2, conf.GAME_H/2)
	titleTxt.SetAlign(text.AlignCenter, text.AlignCenter)
	utils.SetColor(titleTxt.DO, 1, 1, 1, 1)
	scene.TextTitle = titleTxt

	resFontJamboree18 := gameState.Loader.LoadFont(assets.FontJamboree18)
	keys := scene.GameState.InputHandler.ActionKeyNames(ActionEnter, input.KeyboardDevice)
	if len(keys) == 0 {
		panic(fmt.Sprintf("No valid '%d' in action key names", ActionEnter))
	}

	subtitleTxt := NewText(&resFontJamboree18.Face, fmt.Sprintf("press <%s> to continue", keys[0]), true)
	subtitleTxt.SetPos(conf.GAME_W/2, conf.GAME_H/2+titleTxt.DO.LineSpacing*2)
	subtitleTxt.SetAlign(text.AlignCenter, text.AlignCenter)
	subtitleTxt.SpaceY = titleTxt.DO.LineSpacing
	utils.SetColor(subtitleTxt.DO, 1, 1, 1, 1)
	scene.TextSubtitle = subtitleTxt
	scene.FaderSubtitle = effects.NewFader(0, 1, 1)
	scene.FaderSubtitle.Stopped = true

	deskFrameW, deskFrameH := assets.SheetDeskFrameData.W, assets.SheetDeskFrameData.H
	scaleDesk := float64(min(conf.GAME_W/deskFrameW, conf.GAME_H/deskFrameH))
	resDesk := gameState.Loader.LoadImage(assets.ImageSheetDesk)
	scene.AnimDesk = NewAnimationPlayer(resDesk.Data)
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
	resHallway := gameState.Loader.LoadImage(assets.ImageBGHallway)
	scene.AnimHallway = NewAnimationPlayer(resHallway.Data)
	scene.AnimHallway.AddStateAnimation("row1", 0, 0, hallwayFrameW, hallwayFrameH, assets.BGHallwayFrameData.MaxCols, false)
	scene.AnimHallway.SetFPS(0)
	utils.DIOReplaceAlpha(scene.AnimHallway.DIO, 1)
	scene.AnimHallway.Update()
	scene.AnimHallway.DIO.GeoM.Scale(scaleHallway, scaleHallway)
	scene.AnimHallway.DIO.GeoM.Translate(
		conf.GAME_W/2-float64(hallwayFrameW)*scaleHallway/2,
		conf.GAME_H/2-float64(hallwayFrameH)*scaleHallway/2,
	)

	spaceY := float64(8)
	baseX := float64(conf.GAME_W / 2)
	baseY := conf.GAME_H/2 + float64(deskFrameH)*scaleDesk/2 + spaceY
	resFontJamboree26 := gameState.Loader.LoadFont(assets.FontJamboree26)

	texts := []string{"Start", "Settings", "Quit"}
	for _, txt := range texts {
		newTxt := NewText(&resFontJamboree26.Face, txt, true)
		newTxt.SetPos(baseX, baseY)
		newTxt.SetAlign(text.AlignCenter, text.AlignStart)
		newTxt.SpaceY = spaceY
		utils.SetColor(newTxt.DO, 1, 1, 1, 1)
		scene.Texts = append(scene.Texts, newTxt)
		baseY += newTxt.DO.LineSpacing
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

	sceneRoutine := routine.New()
	sceneRoutine.Define(
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

				utils.DOReplaceAlpha(scene.TextTitle.DO, 0)
				scene.AnimDesk.SetStateReset("row1")
				scene.AnimDesk.Update()

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

			return routine.FlowIdle
		}),
	)
	sceneRoutine.Run()
	scene.Routine = sceneRoutine

	return &scene
}

func (scene *Main_Menu_Scene) RandomFlicker() {
	scene.Flickering = rand.IntN(100) >= 60
	if scene.Flickering {
		scene.CurrentStateName = "flickering"
		scene.AnimDesk.SetFPS(float64(utils.IntRandRange(4, 8)))
	} else {
		waitFor := utils.IntRandRange(2, 4)
		scene.CurrentStateName = "waiting flicker... " + strconv.Itoa(waitFor)
		scene.TimerSys.After(time.Second*time.Duration(waitFor), func() {
			scene.RandomFlicker()
		})
		scene.AnimDesk.SetStateReset("row1")
	}
}

func (scene *Main_Menu_Scene) Update() error {
	scene.TimerSys.Update()
	scene.Routine.Update()

	scene.FaderSubtitle.Update()
	scene.TextSubtitle.DO.ColorScale = *scene.FaderSubtitle.GetCS()

	if scene.Flickering {
		scene.AnimDesk.Update()
		if scene.AnimDesk.IsInLastFrame() {
			switch scene.AnimDesk.State() {
			case "row1":
				scene.AnimDesk.SetStateReset("row2")
			case "row2":
				scene.AnimDesk.SetStateReset("row3")
			case "row3":
				scene.RandomFlicker()
			}
		}
	}

	uniform, ok := scene.LayerText.Uniforms.(*assets.MenuTextShaderUniforms)
	if !ok {
		panic("incorrect casting")
	}
	uniform.Time += 0.01
	v := (math.Sin(uniform.Time) + 1) / 2
	v = utils.ClampFloat64(v, 0.4, 0.6)
	uniform.StartingAmplitude = float32(v)
	const speed = 4
	uniform.Velocity[0] = float32(math.Sin(uniform.Time) * speed)
	uniform.Velocity[1] = float32(math.Cos(uniform.Time) * speed)

	if !scene.ShowMenuTexts && scene.TextSubtitle.GetAlpha() >= 0.9 {
		scene.LayerText.Disabled = false
		uniform.Pos[0] = float32(scene.TextSubtitle.X)
		uniform.Pos[1] = float32(scene.TextSubtitle.Y - scene.TextSubtitle.SpaceY/3.5)
		uniform.Size[0] = conf.GAME_W*0.1 + float32(len(scene.TextSubtitle.Txt))*1.5
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
				case MENU_SETTINGS:
					scene.SelectedIdx = MENU_SETTINGS
				case MENU_QUIT:
					return common.ERR_QUIT
				default:
					panic(scene.CurrentIdx)
				}
			}
		}
		scene.CurrentIdx = utils.ClampInt(scene.CurrentIdx, 0, len(scene.Texts)-1)
		curText := scene.Texts[scene.CurrentIdx]
		uniform.Pos[0] = float32(curText.X)
		uniform.Pos[1] = float32(curText.Y)
		uniform.Size[0] = conf.GAME_W*0.1 + float32(len(curText.Txt))*8
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

	if scene.SelectedIdx == MENU_SETTINGS {
		canvas.DrawImage(scene.AnimHallway.CurrentFrame, scene.AnimHallway.DIO)
	} else {
		canvas.DrawImage(scene.AnimDesk.CurrentFrame, scene.AnimDesk.DIO)
		scene.TextTitle.Draw(canvas)
		scene.TextSubtitle.Draw(canvas2)
	}

	if scene.ShowMenuTexts {
		for _, txt := range scene.Texts {
			utils.SetColor(txt.DO, 1, 1, 1, 1)
			txt.Draw(canvas2)
		}
	}

	scene.LayerText.RenderWithShader(canvas)

	scene.LayerColorize.RenderWithShader(screen)
}

var _ Scene = (*Main_Menu_Scene)(nil)
