package main

import (
	"errors"
	"os"
	"remembering-home/src/assets"
	"remembering-home/src/conf"
	"remembering-home/src/errs"
	"remembering-home/src/game"
	"remembering-home/src/logger"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	conf.Log()

	logger.Log().Info("Setting up resources loader...")
	loader := assets.NewAssetsLoader()

	logger.Log().Info("Launching game...")
	sceneManager := game.NewSceneManager()
	inputSystem := game.NewInputSystem()
	gameState := game.NewGame(loader, sceneManager, inputSystem)

	// sceneManager.GoTo(game.NewDummyScene(gameState))
	// sceneManager.GoTo(game.NewSplashScene(gameState))
	sceneManager.GoTo(game.NewMainMenuScene(gameState))

	if err := ebiten.RunGame(gameState); err != nil {
		if errors.Is(err, errs.ERR_QUIT) {
			logger.Log().Info("Successfully exited the game")
			os.Exit(0)
			return
		}
		logger.Log().Fatal(err.Error())
	}
}
