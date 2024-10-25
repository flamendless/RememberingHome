package game

import (
	"fmt"
	// "math"
	"math/rand/v2"
	"nowhere-home/src/assets"
	"nowhere-home/src/common"
	"nowhere-home/src/conf"
	"nowhere-home/src/effects"
	"nowhere-home/src/utils"
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
	texFogFD := assets.TexFogFrameData
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
		LayerText: NewLayerWithShader(
			"menu texts layer",
			texFogFD.W,
			texFogFD.H,
			gameState.Loader.LoadShader(assets.ShaderFog).Data,
		),
	}

	scene.LayerColorize.DRSO.Uniforms = map[string]any{
		"Color": [4]float32{1, 1, 1, 1},
	}

	resFontJamboree46 := gameState.Loader.LoadFont(assets.FontJamboree46)
	titleTxt := NewText(&resFontJamboree46.Face, "Nowhere Home", true)
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

	utils.DRSOSetColor(scene.LayerText.DRSO, 1, 0, 0, 1)
	scene.LayerText.DRSO.Images[1] = gameState.Loader.LoadImage(assets.TextureFog).Data
	scene.LayerText.Uniforms = &assets.FogShader{
		Time:              0,
		StartingAmplitude: 0.5, //0.0 - 0.5
		StartingFreq:      1.0,
		Shift:             0.0,   //-1.0 - 1.0
		WhiteCutoff:       0.999, //0.0 - 1.0
		Velocity:          [2]float32{0.0, 32.0},
		// FogColor:          [4]float32{0.0, 1.0, 1.0, 1.0},
		FogColor:          [4]float32{0.0, 0.0, 0.0, 1.0},
	}
	textH := scene.Texts[0].DO.LineSpacing
	scaleTextBG := float64(textH / float64(texFogFD.H))
	scene.LayerText.ScaleX = 4
	scene.LayerText.ScaleY = scaleTextBG
	scene.LayerText.ScaleY = 4
	scene.LayerText.X = conf.GAME_W/2 - float64(texFogFD.W*int(scene.LayerText.ScaleX))/2
	scene.LayerText.Y = scene.Texts[0].Y - spaceY/2 - float64(texFogFD.H*int(scene.LayerText.ScaleY))/2

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

	uniform, ok := scene.LayerText.Uniforms.(*assets.FogShader)
	if !ok {
		panic("incorrect casting")
	}

	uniform.Time += 0.01
	// v := (math.Sin(uniform.Time) + 1) / 2
	// uniform.StartingAmplitude = float32(v)

	if scene.ShowMenuTexts && scene.CanInteract {
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
		scene.CurrentIdx = utils.ClampInt(scene.CurrentIdx, 0, len(scene.Texts)-1)

		texNoiseFD := assets.TexFogFrameData
		curText := scene.Texts[scene.CurrentIdx]
		scene.LayerText.Y = curText.Y - (curText.SpaceY / 2) - float64((texNoiseFD.H*int(scene.LayerText.ScaleY))/2)
	}

	scene.LayerText.ApplyTransformation()

	return nil
}

func (scene *Main_Menu_Scene) Draw(screen *ebiten.Image) {
	canvas := scene.LayerColorize.Canvas
	canvas.Clear()

	if scene.SelectedIdx == MENU_SETTINGS {
		canvas.DrawImage(scene.AnimHallway.CurrentFrame, scene.AnimHallway.DIO)
	} else {
		canvas.DrawImage(scene.AnimDesk.CurrentFrame, scene.AnimDesk.DIO)
		scene.TextTitle.Draw(canvas)
		scene.TextSubtitle.Draw(canvas)
	}

	if scene.ShowMenuTexts {
		canvas2 := scene.LayerText.Canvas
		canvas2.Clear()
		uniform, ok := scene.LayerText.Uniforms.(*assets.FogShader)
		if !ok {
			panic("incorrect casting")
		}
		scene.LayerText.DRSO.Uniforms = map[string]any{
			"Time":              uniform.Time,
			"StartingAmplitude": uniform.StartingAmplitude,
			"StartingFreq":      uniform.StartingFreq,
			"Shift":             uniform.Shift,
			"WhiteCutoff":       uniform.WhiteCutoff,
			"Velocity":          uniform.Velocity,
			"FogColor":          uniform.FogColor,
		}
		scene.LayerText.ApplyShader(canvas)

		for i, txt := range scene.Texts {
			if i == scene.CurrentIdx {
				utils.SetColor(txt.DO, 1, 0, 0, 1)
			} else {
				utils.SetColor(txt.DO, 1, 1, 1, 1)
			}
			txt.Draw(canvas)
		}
	}

	canvas2 := scene.LayerText.Canvas
	canvas2.Clear()
	uniform, ok := scene.LayerText.Uniforms.(*assets.FogShader)
	if !ok {
		panic("incorrect casting")
	}
	scene.LayerText.DRSO.Uniforms = map[string]any{
		"Time":              uniform.Time,
		"StartingAmplitude": uniform.StartingAmplitude,
		"StartingFreq":      uniform.StartingFreq,
		"Shift":             uniform.Shift,
		"WhiteCutoff":       uniform.WhiteCutoff,
		"Velocity":          uniform.Velocity,
		"FogColor":          uniform.FogColor,
	}
	scene.LayerText.ApplyShader(canvas)

	scene.LayerColorize.ApplyShader(screen)
}

var _ Scene = (*Main_Menu_Scene)(nil)
