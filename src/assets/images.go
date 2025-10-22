package assets

import (
	"remembering-home/src/logger"

	resource "github.com/quasilyte/ebitengine-resource"
	"go.uber.org/zap"
)

const (
	ImageNone resource.ImageID = iota
	ImageDummy
	ImageWindowIcon
	ImageFlamLogo
	ImageSheetWits
	ImageSheetDesk
	ImageSheetTitle
	ImageBGDoor
	ImageBGHallway
	ImageBGStorageRoom
	ImageBGUtilityRoom
	TextureNoise
	TextureFog
	TexturePaper
)

func setImageResources(loader *resource.Loader) {
	logger.Log().Info("Setting image resources...")
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageDummy:      {Path: "dummy.png"},
		ImageWindowIcon: {Path: "icon.png"},

		ImageFlamLogo:  {Path: "splash/logo_flamendless.png"},
		ImageSheetWits: {Path: "splash/sheet_wits.png"},

		ImageSheetDesk:  {Path: "main_menu/sheet_desk_colored.png"},
		ImageSheetTitle: {Path: "main_menu/sheet_title.png"},
		ImageBGDoor:     {Path: "main_menu/bg_door.png"},
		ImageBGHallway:  {Path: "main_menu/bg_hallway.png"},

		ImageBGStorageRoom: {Path: "atlases/storage_room/bg.png"},
		ImageBGUtilityRoom: {Path: "atlases/utility_room/bg.png"},

		TextureNoise: {Path: "textures/noise.png"},
		TextureFog:   {Path: "textures/fog.png"},
		TexturePaper: {Path: "textures/paper.png"},
	}

	_ = map[resource.ImageID]bool{
		ImageNone:          false,
		ImageDummy:         true,
		ImageWindowIcon:    true,
		ImageFlamLogo:      true,
		ImageSheetWits:     true,
		ImageSheetDesk:     true,
		ImageSheetTitle:    true,
		ImageBGDoor:        true,
		ImageBGHallway:     true,
		ImageBGStorageRoom: true,
		ImageBGUtilityRoom: true,
		TextureNoise:       true,
		TextureFog:         true,
		TexturePaper:       true,
	}

	for id, res := range imageResources {
		logger.Log().Info("Loading image", zap.String("path", res.Path))
		loader.ImageRegistry.Set(id, res)
		loader.LoadImage(id)
	}
}
