package scenes

import (
	"remembering-home/src/assets"
	"remembering-home/src/conf"
	"remembering-home/src/context"
	"remembering-home/src/graphics"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

type Splash_Scene struct {
	Context          *context.GameContext
	SceneManager     SceneManager
	Routine          *routine.Routine
	FlamLogoAnim     *graphics.AnimationPlayer
	WitsAnim         *graphics.AnimationPlayer
	CurrentStateName string
	ShowWits         bool
	FinishedWits     bool
}

func (scene *Splash_Scene) GetName() string {
	return "Splash"
}

func (scene *Splash_Scene) GetStateName() string {
	return scene.CurrentStateName
}

func NewSplashScene(ctx *context.GameContext, sceneManager SceneManager) *Splash_Scene {
	scene := Splash_Scene{
		Context:      ctx,
		SceneManager: sceneManager,
	}

	resFlamLogo := ctx.Loader.LoadImage(assets.ImageFlamLogo)
	sizeFlamLogo := resFlamLogo.Data.Bounds()
	scene.FlamLogoAnim = graphics.NewAnimationPlayer(resFlamLogo.Data)
	scene.FlamLogoAnim.AddStateAnimation("static", 0, 0, sizeFlamLogo.Max.X, sizeFlamLogo.Max.Y, 1, false)

	scaleFlamLogo := float64(min(conf.GAME_W*0.7/sizeFlamLogo.Max.X, conf.GAME_H*0.7/sizeFlamLogo.Max.Y))
	scene.FlamLogoAnim.DIO.GeoM.Scale(scaleFlamLogo, scaleFlamLogo)
	scene.FlamLogoAnim.DIO.GeoM.Translate(
		conf.GAME_W/2-float64(sizeFlamLogo.Max.X)*scaleFlamLogo/2,
		conf.GAME_H/2-float64(sizeFlamLogo.Max.Y)*scaleFlamLogo/2,
	)

	resWits := ctx.Loader.LoadImage(assets.ImageSheetWits)
	witsFrameW, witsFrameH := assets.SheetWitsFrameData.W, assets.SheetWitsFrameData.H
	scene.WitsAnim = graphics.NewAnimationPlayer(resWits.Data)
	scene.WitsAnim.AddStateAnimation("row1", 0, 0, witsFrameW, witsFrameH, assets.SheetWitsFrameData.MaxCols, false)
	scene.WitsAnim.AddStateAnimation("row2", 0, 128, witsFrameW, witsFrameH, assets.SheetWitsFrameData.MaxCols, false)
	scene.WitsAnim.AddStateAnimation("row3", 0, 256, witsFrameW, witsFrameH, assets.SheetWitsFrameData.MaxCols, false)
	scene.WitsAnim.AddStateAnimation("row4", 0, 384, witsFrameW, witsFrameH, assets.SheetWitsFrameData.MaxCols, false)
	scene.WitsAnim.SetFPS(7)
	scene.WitsAnim.SetStateReset("row1")

	scaleWitsAnim := float64(min(conf.GAME_W/witsFrameW, conf.GAME_H/witsFrameH))
	scene.WitsAnim.DIO.GeoM.Scale(scaleWitsAnim, scaleWitsAnim)
	scene.WitsAnim.DIO.GeoM.Translate(
		float64(conf.GAME_W/2-float64(witsFrameW)*scaleWitsAnim/2),
		float64(conf.GAME_H/2-float64(witsFrameH)*scaleWitsAnim/2),
	)

	scene.Routine = routine.New()
	scene.Routine.Define(
		"splash scene",
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "flamendless logo fading in"
			if scene.SceneManager.IsFadeInFinished() {
				return routine.FlowNext
			}
			return routine.FlowIdle
		}),
		actions.NewWait(time.Second*2),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "wits animation showing"
			scene.ShowWits = true
			return routine.FlowNext
		}),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			if scene.FinishedWits {
				scene.CurrentStateName = "wits waiting"
				return routine.FlowNext
			}
			scene.CurrentStateName = "wits animating"
			return routine.FlowIdle
		}),
		actions.NewWait(time.Second/2),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			scene.CurrentStateName = "wits fading out"
			scene.SceneManager.GoTo(NewMainMenuScene(scene.Context, scene.SceneManager))
			return routine.FlowIdle
		}),
	)
	scene.Routine.Run()
	return &scene
}

func (scene *Splash_Scene) Update() error {
	if scene.Routine.Running() {
		scene.Routine.Update()
	}

	if scene.ShowWits {
		scene.WitsAnim.Update()
		if scene.WitsAnim.IsInLastFrame() {
			switch scene.WitsAnim.State() {
			case "row1":
				scene.WitsAnim.SetStateReset("row2")
			case "row2":
				scene.WitsAnim.SetStateReset("row3")
			case "row3":
				scene.WitsAnim.SetStateReset("row4")
			case "row4":
				scene.WitsAnim.PauseAtFrame(assets.SheetWitsFrameData.MaxCols - 1)
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

var _ Scene = (*Splash_Scene)(nil)
