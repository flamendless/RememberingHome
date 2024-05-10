package game

import (
	"nowhere-home/src/assets"
	"nowhere-home/src/conf"
	"nowhere-home/src/overlays"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

type Splash_Scene struct {
	GameState    *Game_State
	Routine      *routine.Routine
	FlamLogoAnim *AnimationPlayer
	WitsAnim     *AnimationPlayer
	ShowWits     bool
}

func (scene Splash_Scene) GetName() string {
	return "Splash"
}

func NewSplashScene(gameState *Game_State) *Splash_Scene {
	splashScene := Splash_Scene{
		GameState: gameState,
	}
	resFlamLogo := gameState.Loader.LoadImage(assets.ImageFlamLogo)
	splashScene.FlamLogoAnim = NewAnimationPlayer(resFlamLogo.Data)
	splashScene.FlamLogoAnim.AddStateAnimation("static", 0, 0, resFlamLogo.Data.Bounds().Max.X, resFlamLogo.Data.Bounds().Max.Y, 1, false)

	resWits := gameState.Loader.LoadImage(assets.ImageSheetWits)
	witsFrameX, witsFrameY := assets.SheetWitsFrameData.W, assets.SheetWitsFrameData.H
	splashScene.WitsAnim = NewAnimationPlayer(resWits.Data)
	splashScene.WitsAnim.AddStateAnimation("row1", 0, 0, witsFrameX, witsFrameY, assets.SheetWitsFrameData.Count, false)
	splashScene.WitsAnim.AddStateAnimation("row2", 0, 128, witsFrameX, witsFrameY, assets.SheetWitsFrameData.Count, false)
	splashScene.WitsAnim.AddStateAnimation("row3", 0, 256, witsFrameX, witsFrameY, assets.SheetWitsFrameData.Count, false)
	splashScene.WitsAnim.AddStateAnimation("row4", 0, 384, witsFrameX, witsFrameY, assets.SheetWitsFrameData.Count, false)
	splashScene.WitsAnim.SetFPS(7)
	splashScene.WitsAnim.SetStateReset("row1")

	sceneRoutine := routine.New()
	sceneRoutine.Define(
		"scene",
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			if overlays.IsFadeInFinished() {
				return routine.FlowNext
			}
			return routine.FlowIdle
		}),

		actions.NewWait(time.Second*2),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			splashScene.ShowWits = true
			return routine.FlowIdle
		}),
	)
	sceneRoutine.Run()
	splashScene.Routine = sceneRoutine

	return &splashScene
}

func (scene *Splash_Scene) Update() error {
	if scene.Routine.Running() {
		scene.Routine.Update()
	}

	if scene.ShowWits {
		scene.WitsAnim.Update()
		if scene.WitsAnim.CurrentFrameIndex == 2 {
			switch scene.WitsAnim.State() {
			case "row1":
				scene.WitsAnim.SetStateReset("row2")
			case "row2":
				scene.WitsAnim.SetStateReset("row3")
			case "row3":
				scene.WitsAnim.SetStateReset("row4")
			case "row4":
				scene.WitsAnim.PauseAtFrame(2)
			}
		}
	}
	return nil
}

func (scene *Splash_Scene) Draw(screen *ebiten.Image) {
	resFlam := scene.GameState.Loader.LoadImage(assets.ImageFlamLogo)
	sizeFlam := resFlam.Data.Bounds().Size()
	opFlam := &ebiten.DrawImageOptions{}
	opFlam.GeoM.Translate(
		float64(conf.GAME_W/2-sizeFlam.X/2),
		float64(conf.GAME_H/2-sizeFlam.Y/2),
	)
	screen.DrawImage(scene.FlamLogoAnim.CurrentFrame, opFlam)

	if scene.ShowWits {
		witsFrameX, witsFrameY := assets.SheetWitsFrameData.W, assets.SheetWitsFrameData.H
		scale := min(float64(conf.GAME_W/witsFrameX), float64(conf.GAME_H/witsFrameY))
		opWits := &ebiten.DrawImageOptions{}
		opWits.GeoM.Scale(scale, scale)
		opWits.GeoM.Translate(
			float64(conf.GAME_W/2-float64(witsFrameX)*scale/2),
			float64(conf.GAME_H/2-float64(witsFrameY)*scale/2),
		)

		screen.DrawImage(scene.WitsAnim.CurrentFrame, opWits)
	}
}
