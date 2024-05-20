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
	SetShaderResources(loader)

	return loader
}
