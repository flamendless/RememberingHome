package main

import (
	"nowhere-home/internal/assets"
	"nowhere-home/internal/conf"
	"nowhere-home/internal/game"
	"nowhere-home/internal/logger"

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
