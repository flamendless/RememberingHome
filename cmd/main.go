package main

import (
	"nowhere-home/src/assets"
	"nowhere-home/src/conf"
	"nowhere-home/src/game"
	"nowhere-home/src/logger"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	logger.InitLog()
	conf.LogConf()

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
		logger.Log().Fatal(err.Error())
	}
}
