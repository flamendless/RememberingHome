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

const (
	ImageNone resource.ImageID = iota
	ImageWindowIcon
	ImageFlamLogo
	ImageSheetWits
	ImageSheetDesk
	ImageBGDoor
	ImageBGHallway
)

const (
	FontNone resource.FontID = iota
	FontJamboree18
	FontJamboree26
	FontJamboree46
)

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
	SetFontResources(loader)

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
		loader.LoadImage(id)
	}
}

func SetFontResources(loader *resource.Loader) {
	fontResources := map[resource.FontID]resource.FontInfo{
		FontJamboree18: {Path: "fonts/Jamboree.ttf", Size: 18},
		FontJamboree26: {Path: "fonts/Jamboree.ttf", Size: 26},
		FontJamboree46: {Path: "fonts/Jamboree.ttf", Size: 46},
	}
	for id, res := range fontResources {
		loader.FontRegistry.Set(id, res)
		loader.LoadFont(id)
	}
}
