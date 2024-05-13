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
	GameState        *Game_State
	Routine          *routine.Routine
	FlamLogoAnim     *AnimationPlayer
	WitsAnim         *AnimationPlayer
	ShowWits         bool
	FinishedWits     bool
	CurrentStateName string
}

func (scene Splash_Scene) GetName() string {
	return "Splash"
}

func (scene Splash_Scene) GetStateName() string {
	return scene.CurrentStateName
}

func NewSplashScene(gameState *Game_State) *Splash_Scene {
	splashScene := Splash_Scene{
		GameState: gameState,
	}

	resFlamLogo := gameState.Loader.LoadImage(assets.ImageFlamLogo)
	resWits := gameState.Loader.LoadImage(assets.ImageSheetWits)

	sizeFlamLogo := resFlamLogo.Data.Bounds()
	scaleFlamLogo := float64(min(conf.GAME_W*0.7/sizeFlamLogo.Max.X, conf.GAME_H*0.7/sizeFlamLogo.Max.Y))

	witsFrameX, witsFrameY := assets.SheetWitsFrameData.W, assets.SheetWitsFrameData.H
	scaleWitsAnim := float64(min(conf.GAME_W/witsFrameX, conf.GAME_H/witsFrameY))

	splashScene.FlamLogoAnim = NewAnimationPlayer(resFlamLogo.Data)
	splashScene.FlamLogoAnim.AddStateAnimation("static", 0, 0, sizeFlamLogo.Max.X, sizeFlamLogo.Max.Y, 1, false)

	splashScene.FlamLogoAnim.DIO.GeoM.Scale(scaleFlamLogo, scaleFlamLogo)
	splashScene.FlamLogoAnim.DIO.GeoM.Translate(
		float64(conf.GAME_W/2-sizeFlamLogo.Max.X*int(scaleFlamLogo)/2),
		float64(conf.GAME_H/2-sizeFlamLogo.Max.Y*int(scaleFlamLogo)/2),
	)

	splashScene.WitsAnim = NewAnimationPlayer(resWits.Data)
	splashScene.WitsAnim.AddStateAnimation("row1", 0, 0, witsFrameX, witsFrameY, assets.SheetWitsFrameData.Count, false)
	splashScene.WitsAnim.AddStateAnimation("row2", 0, 128, witsFrameX, witsFrameY, assets.SheetWitsFrameData.Count, false)
	splashScene.WitsAnim.AddStateAnimation("row3", 0, 256, witsFrameX, witsFrameY, assets.SheetWitsFrameData.Count, false)
	splashScene.WitsAnim.AddStateAnimation("row4", 0, 384, witsFrameX, witsFrameY, assets.SheetWitsFrameData.Count, false)
	splashScene.WitsAnim.SetFPS(7)
	splashScene.WitsAnim.SetStateReset("row1")

	splashScene.WitsAnim.DIO.GeoM.Scale(scaleWitsAnim, scaleWitsAnim)
	splashScene.WitsAnim.DIO.GeoM.Translate(
		float64(conf.GAME_W/2-float64(witsFrameX)*scaleWitsAnim/2),
		float64(conf.GAME_H/2-float64(witsFrameY)*scaleWitsAnim/2),
	)

	sceneRoutine := routine.New()
	sceneRoutine.Define(
		"splash scene",
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			splashScene.CurrentStateName = "flamendless logo fading in"
			if overlays.IsFadeInFinished() {
				return routine.FlowNext
			}
			return routine.FlowIdle
		}),
		actions.NewWait(time.Second*2),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			splashScene.CurrentStateName = "wits animation showing"
			splashScene.ShowWits = true
			return routine.FlowNext
		}),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			if splashScene.FinishedWits {
				splashScene.CurrentStateName = "wits waiting"
				return routine.FlowNext
			}
			splashScene.CurrentStateName = "wits animating"
			return routine.FlowIdle
		}),
		actions.NewWait(time.Second/2),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			splashScene.CurrentStateName = "wits fading out"
			splashScene.GameState.SceneManager.GoTo(&Dummy_Scene{GameState: splashScene.GameState})
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
				scene.FinishedWits = true
			}
		}
	}
	return nil
}

func (scene *Splash_Scene) Draw(screen *ebiten.Image) {
	screen.DrawImage(scene.FlamLogoAnim.CurrentFrame, scene.FlamLogoAnim.DIO)
	if scene.ShowWits {
		screen.DrawImage(scene.WitsAnim.CurrentFrame, scene.WitsAnim.DIO)
	}
}
