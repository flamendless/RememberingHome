package game

import (
	"math/rand/v2"
	"nowhere-home/src/assets"
	"nowhere-home/src/conf"
	"nowhere-home/src/utils"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ebitick"
)

type Main_Menu_Scene struct {
	GameState        *Game_State
	TimerSys         *ebitick.TimerSystem
	CurrentStateName string
	DeskAnim         *AnimationPlayer
	Flickering       bool
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
	}

	resDesk := gameState.Loader.LoadImage(assets.ImageSheetDesk)

	deskFrameW, deskFrameH := assets.SheetDeskFrameData.W, assets.SheetDeskFrameData.H
	scaleDesk := float64(min(conf.GAME_W/deskFrameW, conf.GAME_H/deskFrameH))

	scene.DeskAnim = NewAnimationPlayer(resDesk.Data)
	scene.DeskAnim.AddStateAnimation("row1", 0, 0, deskFrameW, deskFrameH, assets.SheetDeskFrameData.MaxCols, false)
	scene.DeskAnim.AddStateAnimation("row2", 0, 64, deskFrameW, deskFrameH, assets.SheetDeskFrameData.MaxCols, false)
	scene.DeskAnim.AddStateAnimation("row3", 0, 128, deskFrameW, deskFrameH, 1, false)
	scene.DeskAnim.SetFPS(7)
	scene.DeskAnim.SetStateReset("row1")

	scene.DeskAnim.DIO.GeoM.Scale(scaleDesk, scaleDesk)
	scene.DeskAnim.DIO.GeoM.Translate(
		float64(conf.GAME_W/2-float64(deskFrameW)*scaleDesk/2),
		float64(conf.GAME_H/2-float64(deskFrameH)*scaleDesk/2),
	)

	waitFor := utils.IntRandRange(1, 3)
	scene.CurrentStateName = "waiting flicker... " + strconv.Itoa(waitFor)
	scene.TimerSys.After(time.Second*time.Duration(waitFor), func() {
		scene.RandomFlicker()
	})

	return &scene
}

func (scene *Main_Menu_Scene) RandomFlicker() {
	scene.Flickering = rand.IntN(100) >= 60
	if scene.Flickering {
		scene.CurrentStateName = "flickering"
		scene.DeskAnim.SetFPS(float64(utils.IntRandRange(4, 8)))
	} else {
		waitFor := utils.IntRandRange(1, 3)
		scene.CurrentStateName = "waiting flicker... " + strconv.Itoa(waitFor)
		scene.TimerSys.After(time.Second*time.Duration(waitFor), func() {
			scene.RandomFlicker()
		})
		scene.DeskAnim.SetStateReset("row1")
	}
}

func (scene *Main_Menu_Scene) Update() error {
	scene.TimerSys.Update()
	if scene.Flickering {
		scene.DeskAnim.Update()
		if scene.DeskAnim.CurrentFrameIndex == scene.DeskAnim.GetLastFrameCount() {
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

	return nil
}

func (scene *Main_Menu_Scene) Draw(screen *ebiten.Image) {
	screen.DrawImage(scene.DeskAnim.CurrentFrame, scene.DeskAnim.DIO)
}
