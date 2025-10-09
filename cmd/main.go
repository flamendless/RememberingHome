package main

import (
	"errors"
	"os"
	"remembering-home/src/assets"
	"remembering-home/src/conf"
	"remembering-home/src/errs"
	"remembering-home/src/game"
	"remembering-home/src/logger"
	"remembering-home/src/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	conf.Log()

	logger.Log().Info("Setting up resources loader...")
	loader := assets.NewAssetsLoader()

	logger.Log().Info("Launching game...")
	sceneManager := scenes.NewSceneManager()
	inputSystem := game.NewInputSystem()
	gameState := game.NewGame(loader, sceneManager, inputSystem)

	// sceneManager.GoTo(scenes.NewDummyScene(gameState.Context, sceneManager))
	// sceneManager.GoTo(scenes.NewSplashScene(gameState.Context, sceneManager))
	sceneManager.GoTo(scenes.NewMainMenuScene(gameState.Context, sceneManager))

	if err := ebiten.RunGame(gameState); err != nil {
		if errors.Is(err, errs.ErrQuit) {
			logger.Log().Info("Successfully exited the game")
			os.Exit(0)
			return
		}
		logger.Log().Fatal(err.Error())
	}
}
