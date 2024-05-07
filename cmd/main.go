package main

import (
	"somewhere-home/internal/assets"
	"somewhere-home/internal/conf"
	"somewhere-home/internal/game"
	"somewhere-home/internal/logger"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	logger.InitLog()
	conf.LogConf()

	logger.Log().Info("Setting up resources loader...")
	loader := assets.NewAssetsLoader()

	logger.Log().Info("Launching game...")
	game := game.NewGame(loader)

	if err := ebiten.RunGame(game); err != nil {
		logger.Log().Fatal(err.Error())
	}
}
