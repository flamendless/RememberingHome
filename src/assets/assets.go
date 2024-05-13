package assets

import (
	"embed"
	"io"
	"nowhere-home/src/logger"

	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
)

//go:embed data
var gameAssets embed.FS

type FrameData struct {
	W, H, MaxCols int
}

const (
	ImageNone resource.ImageID = iota
	ImageWindowIcon
	ImageFlamLogo
	ImageSheetWits
	ImageSheetDesk
	ImageBGDoor
	ImageBGHallway
)

var (
	SheetWitsFrameData FrameData
	SheetDeskFrameData FrameData
)

func init() {
	SheetWitsFrameData = FrameData{W: 256, H: 128, MaxCols: 3}
	SheetDeskFrameData = FrameData{W: 256, H: 64, MaxCols: 3}
}

func NewAssetsLoader() *resource.Loader {
	audioCtx := audio.NewContext(44100)
	loader := resource.NewLoader(audioCtx)

	loader.OpenAssetFunc = func(path string) io.ReadCloser {
		f, err := gameAssets.Open("data/" + path)
		if err != nil {
			logger.Log().Error(err.Error())
		}
		return f
	}

	SetImageResources(loader)

	return loader
}

func SetImageResources(loader *resource.Loader) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageWindowIcon: {Path: "icon.png"},

		ImageFlamLogo:  {Path: "splash/logo_flamendless.png"},
		ImageSheetWits: {Path: "splash/sheet_wits.png"},

		ImageSheetDesk: {Path: "main_menu/sheet_desk.png"},
		ImageBGDoor:    {Path: "main_menu/bg_door.png"},
		ImageBGHallway: {Path: "main_menu/bg_hallway.png"},
	}
	for id, res := range imageResources {
		loader.ImageRegistry.Set(id, res)
	}
}
