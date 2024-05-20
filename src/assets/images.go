package assets

import (
	"nowhere-home/src/logger"

	resource "github.com/quasilyte/ebitengine-resource"
	"go.uber.org/zap"
)

const (
	ImageNone resource.ImageID = iota
	ImageWindowIcon
	ImageFlamLogo
	ImageSheetWits
	ImageSheetDesk
	ImageBGDoor
	ImageBGHallway
)

func SetImageResources(loader *resource.Loader) {
	logger.Log().Info("Setting image resources...")
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageWindowIcon: {Path: "icon.png"},

		ImageFlamLogo:  {Path: "splash/logo_flamendless.png"},
		ImageSheetWits: {Path: "splash/sheet_wits.png"},

		ImageSheetDesk: {Path: "main_menu/sheet_desk_colored.png"},
		ImageBGDoor:    {Path: "main_menu/bg_door.png"},
		ImageBGHallway: {Path: "main_menu/bg_hallway.png"},
	}
	for id, res := range imageResources {
		logger.Log().Info("Loading image", zap.String("path", res.Path))
		loader.ImageRegistry.Set(id, res)
		loader.LoadImage(id)
	}
}
