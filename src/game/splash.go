package game

import (
	"fmt"
	"nowhere-home/src/assets"
	"nowhere-home/src/conf"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

type Splash_Scene struct {
	GameState   *Game_State
	Routine     *routine.Routine
}

func (scene Splash_Scene) GetName() string {
	return "Splash"
}

func NewSplashScene(gameState *Game_State) *Splash_Scene {
	testRoutine := routine.New()

	testRoutine.Define(
		"test",
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			fmt.Println(11111)
			return routine.FlowNext
		}),
		actions.NewWait(time.Second * 2),
		actions.NewFunction(func(block *routine.Block) routine.Flow {
			fmt.Println(22222)
			return routine.FlowNext
		}),
	)
	testRoutine.Run()

	splashScene := Splash_Scene{
		GameState: gameState,
		Routine: testRoutine,
	}
	return &splashScene
}

func (scene *Splash_Scene) Update() error {
	if scene.Routine.Running() {
		scene.Routine.Update()
	}
	return nil
}

func (scene *Splash_Scene) Draw(screen *ebiten.Image) {
	logo := scene.GameState.Loader.LoadImage(assets.ImageFlamendlessLogo)
	op := &ebiten.DrawImageOptions{}
	size := logo.Data.Bounds().Size()
	op.GeoM.Translate(
		float64(conf.GAME_W/2-size.X/2),
		float64(conf.GAME_H/2-size.Y/2),
	)
	screen.DrawImage(logo.Data, op)
}
