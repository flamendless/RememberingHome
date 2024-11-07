package assets

import (
	"remembering-home/src/logger"

	resource "github.com/quasilyte/ebitengine-resource"
	"go.uber.org/zap"
)

const (
	FontNone resource.FontID = iota
	FontJamboree18
	FontJamboree26
	FontJamboree46
)

func SetFontResources(loader *resource.Loader) {
	logger.Log().Info("Setting font resources...")
	fontResources := map[resource.FontID]resource.FontInfo{
		FontJamboree18: {Path: "fonts/Jamboree.ttf", Size: 18},
		FontJamboree26: {Path: "fonts/Jamboree.ttf", Size: 26},
		FontJamboree46: {Path: "fonts/Jamboree.ttf", Size: 46},
	}
	for id, res := range fontResources {
		logger.Log().Info("Loading font", zap.String("path", res.Path))
		loader.FontRegistry.Set(id, res)
		loader.LoadFont(id)
	}
}
