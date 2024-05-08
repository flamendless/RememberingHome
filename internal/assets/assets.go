package assets

import (
	"embed"
	"io"
	"nowhere-home/internal/logger"

	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
)

//go:embed data
var gameAssets embed.FS

const (
	ImageNone resource.ImageID = iota
	ImageWindowIcon
	ImageFlamendlessLogo
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

	return loader
}

func SetImageResources(loader *resource.Loader) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageWindowIcon:      {Path: "icon.png"},
		ImageFlamendlessLogo: {Path: "logo_flamendless.png"},
	}
	for id, res := range imageResources {
		loader.ImageRegistry.Set(id, res)
	}
}
