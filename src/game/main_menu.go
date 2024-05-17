package game

import (
	"fmt"
	"math/rand/v2"
	"nowhere-home/src/assets"
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

type Main_Menu_Scene struct {
	GameState        *Game_State
	TimerSys         *ebitick.TimerSystem
	CurrentStateName string
	Routine          *routine.Routine
	TitleText        *Text
	SubtitleText     *Text
	SubtitleFader    *effects.Fader
	ShowMenuTexts    bool
	DeskAnim         *AnimationPlayer
	Flickering       bool
	Texts            []*Text
	CurrentIdx       int
}

func (scene Main_Menu_Scene) GetName() string {
	return "Main Menu"
}

func (scene Main_Menu_Scene) GetStateName() string {
	return scene.CurrentStateName
}

func NewMainMenuScene(gameState *Game_State) *Main_Menu_Scene {
	scene := Main_Menu_Scene{
		GameState: gameState,
		TimerSys:  ebitick.NewTimerSystem(),
		Texts:     make([]*Text, 0, 16),
	}

	resFontJamboree46 := gameState.Loader.LoadFont(assets.FontJamboree46)
	titleTxt := NewText(&resFontJamboree46.Face, "Nowhere Home", true)
	titleTxt.SetPos(conf.GAME_W/2, conf.GAME_H/2)
	titleTxt.SetAlign(text.AlignCenter, text.AlignCenter)
	titleTxt.SetColor(1, 1, 1, 1)
	scene.TitleText = titleTxt

	resFontJamboree18 := gameState.Loader.LoadFont(assets.FontJamboree18)
	enterKey := scene.GameState.InputHandler.ActionKeyNames(ActionEnter, input.KeyboardDevice)[0]
	subtitleTxt := NewText(&resFontJamboree18.Face, fmt.Sprintf("press <%s> to continue", enterKey), true)
	subtitleTxt.SetPos(conf.GAME_W/2, conf.GAME_H/2+titleTxt.DO.LineSpacing*2)
	subtitleTxt.SetAlign(text.AlignCenter, text.AlignCenter)
	subtitleTxt.SetColor(1, 1, 1, 1)
	scene.SubtitleText = subtitleTxt
	scene.SubtitleFader = effects.NewFader(0, 1, 1)
	scene.SubtitleFader.Stopped = true

	resDesk := gameState.Loader.LoadImage(assets.ImageSheetDesk)

	deskFrameW, deskFrameH := assets.SheetDeskFrameData.W, assets.SheetDeskFrameData.H
	scaleDesk := float64(min(conf.GAME_W/deskFrameW, conf.GAME_H/deskFrameH))

	scene.DeskAnim = NewAnimationPlayer(resDesk.Data)
	scene.DeskAnim.AddStateAnimation("static", deskFrameW, 128, deskFrameW, deskFrameH, 1, false)
	scene.DeskAnim.AddStateAnimation("row1", 0, 0, deskFrameW, deskFrameH, assets.SheetDeskFrameData.MaxCols, false)
	scene.DeskAnim.AddStateAnimation("row2", 0, 64, deskFrameW, deskFrameH, assets.SheetDeskFrameData.MaxCols, false)
	scene.DeskAnim.AddStateAnimation("row3", 0, 128, deskFrameW, deskFrameH, 1, false)
	scene.DeskAnim.SetFPS(7)
	scene.DeskAnim.SetStateReset("static")
	scene.DeskAnim.DIO.ColorScale.ScaleAlpha(0)

	scene.DeskAnim.DIO.GeoM.Scale(scaleDesk, scaleDesk)
	scene.DeskAnim.DIO.GeoM.Translate(
		conf.GAME_W/2-float64(deskFrameW)*scaleDesk/2,
		conf.GAME_H/2-float64(deskFrameH)*scaleDesk/2,
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
		newTxt.SetColor(1, 1, 1, 1)
		scene.Texts = append(scene.Texts, newTxt)
		baseY += newTxt.DO.LineSpacing
	}

	sceneRoutine := routine.New()
	sceneRoutine.Define(
		"main menu scene",
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "title text animation"
			if scene.GameState.SceneManager.IsFadeInFinished() {
				scene.SubtitleFader.Stopped = false

				scene.DeskAnim.SetStateReset("static")
				utils.DIOReplaceAlpha(scene.DeskAnim.DIO, 1)
				scene.DeskAnim.Update()

				return routine.FlowNext
			}
			return routine.FlowIdle
		}),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "waiting input"
			if scene.GameState.InputHandler.ActionIsJustPressed(ActionEnter) {
				scene.CurrentStateName = "showing menu..."
				scene.SubtitleFader.Alpha = 0
				scene.SubtitleFader.Stopped = true

				utils.DOReplaceAlpha(scene.TitleText.DO, 0)
				scene.DeskAnim.SetStateReset("row1")
				scene.DeskAnim.Update()

				return routine.FlowNext
			}
			return routine.FlowIdle
		}),
		// actions.NewFunction(func(block *routine.Block) routine.Flow {
		// 	scene.CurrentStateName = "moving title text"
		// 	baseY := conf.GAME_H/2-float64(deskFrameH)*scaleDesk/2
		// 	const speed = 8
		// 	if scene.TitleText.Y > baseY - scene.TitleText.DO.LineSpacing/2 {
		// 		scene.TitleText.SetPos(scene.TitleText.X, scene.TitleText.Y - speed)
		// 	} else {
		// 		return routine.FlowNext
		// 	}
		// 	return routine.FlowIdle
		// }),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "finished"
			scene.ShowMenuTexts = true

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
		scene.DeskAnim.SetFPS(float64(utils.IntRandRange(4, 8)))
	} else {
		waitFor := utils.IntRandRange(2, 4)
		scene.CurrentStateName = "waiting flicker... " + strconv.Itoa(waitFor)
		scene.TimerSys.After(time.Second*time.Duration(waitFor), func() {
			scene.RandomFlicker()
		})
		scene.DeskAnim.SetStateReset("row1")
	}
}

func (scene *Main_Menu_Scene) Update() error {
	scene.TimerSys.Update()
	scene.Routine.Update()

	scene.SubtitleFader.Update()
	scene.SubtitleText.DO.ColorScale = *scene.SubtitleFader.GetCS()

	if scene.Flickering {
		scene.DeskAnim.Update()
		if scene.DeskAnim.IsInLastFrame() {
			switch scene.DeskAnim.State() {
			case "row1":
				scene.DeskAnim.SetStateReset("row2")
			case "row2":
				scene.DeskAnim.SetStateReset("row3")
			case "row3":
				scene.RandomFlicker()
			}
		}
	}

	if scene.ShowMenuTexts {
		if scene.GameState.InputHandler.ActionIsJustReleased(ActionMoveUp) {
			scene.CurrentIdx--
		} else if scene.GameState.InputHandler.ActionIsJustReleased(ActionMoveDown) {
			scene.CurrentIdx++
		}
		scene.CurrentIdx = utils.ClampInt(scene.CurrentIdx, 0, len(scene.Texts)-1)
	}

	return nil
}

func (scene *Main_Menu_Scene) Draw(screen *ebiten.Image) {
	screen.DrawImage(scene.DeskAnim.CurrentFrame, scene.DeskAnim.DIO)
	scene.TitleText.Draw(screen)
	scene.SubtitleText.Draw(screen)

	if scene.ShowMenuTexts {
		for i, txt := range scene.Texts {
			if i == scene.CurrentIdx {
				txt.SetColor(1, 0, 0, 1)
			} else {
				txt.SetColor(1, 1, 1, 1)
			}
			txt.Draw(screen)
		}
	}
}
