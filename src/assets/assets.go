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
	W, H, Count int
}

const (
	ImageNone resource.ImageID = iota
	ImageWindowIcon
	ImageFlamLogo
	ImageSheetWits
)

var (
	SheetWitsFrameData FrameData
)

func init() {
	SheetWitsFrameData = FrameData{W: 256, H: 128, Count: 3}
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
		ImageFlamLogo:   {Path: "logo_flamendless.png"},
		ImageSheetWits:  {Path: "sheet_wits.png"},
	}
	for id, res := range imageResources {
		loader.ImageRegistry.Set(id, res)
	}
}
