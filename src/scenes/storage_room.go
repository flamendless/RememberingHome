package scenes

import (
	"remembering-home/src/assets"
	"remembering-home/src/atlases"
	"remembering-home/src/common"
	"remembering-home/src/conf"
	"remembering-home/src/context"
	"remembering-home/src/debug"
	"remembering-home/src/graphics"
	"remembering-home/src/items"

	"github.com/hajimehoshi/ebiten/v2"
)

type StorageRoom_Scene struct {
	Context          *context.GameContext
	SceneManager     SceneManager
	CurrentStateName string
	BackgroundImage  *ebiten.Image
	BackgroundDIO    *ebiten.DrawImageOptions
	ItemRenderer     *graphics.ItemRenderer
}

func (scene *StorageRoom_Scene) GetName() string {
	return "Storage Room"
}

func (scene *StorageRoom_Scene) GetStateName() string {
	return scene.CurrentStateName
}

func (scene *StorageRoom_Scene) GetItemRenderer() *graphics.ItemRenderer {
	return scene.ItemRenderer
}

func NewStorageRoomScene(ctx *context.GameContext, sceneManager SceneManager) *StorageRoom_Scene {
	scene := StorageRoom_Scene{
		Context:      ctx,
		SceneManager: sceneManager,
	}

	resBG := ctx.Loader.LoadImage(assets.ImageBGStorageRoom)
	scene.BackgroundImage = resBG.Data
	scene.BackgroundDIO = &ebiten.DrawImageOptions{}

	bounds := resBG.Data.Bounds()
	imageWidth := bounds.Dx()
	imageHeight := bounds.Dy()

	scaleX := float64(conf.GAME_W) / float64(imageWidth)
	scaleY := float64(conf.GAME_H) / float64(imageHeight)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	scene.BackgroundDIO.GeoM.Scale(scale, scale)
	scene.BackgroundDIO.GeoM.Translate(
		conf.GAME_W/2-float64(imageWidth)*scale/2,
		conf.GAME_H/2-float64(imageHeight)*scale/2,
	)

	scene.ItemRenderer = graphics.NewItemRenderer(ctx, assets.ImageAtlasStorageRoom, atlases.AtlasStorageRoom, items.ItemStorageRoom)

	return &scene
}

func (scene *StorageRoom_Scene) Update() error {
	return nil
}

func (scene *StorageRoom_Scene) Draw(screen *ebiten.Image) {
	screen.DrawImage(scene.BackgroundImage, scene.BackgroundDIO)

	bounds := scene.BackgroundImage.Bounds()
	imageWidth := bounds.Dx()
	imageHeight := bounds.Dy()

	scaleX := float64(conf.GAME_W) / float64(imageWidth)
	scaleY := float64(conf.GAME_H) / float64(imageHeight)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	baseOffset := common.Vec2{
		X: conf.GAME_W/2 - float64(imageWidth)*scale/2,
		Y: conf.GAME_H/2 - float64(imageHeight)*scale/2,
	}

	scene.ItemRenderer.DrawAllItems(screen, scale, baseOffset)

	if debug.ShowItemSelection {
		zoomLevel := debug.GetZoomLevel()
		scene.ItemRenderer.DrawItemSelection(screen, scale, baseOffset, zoomLevel)
	}
}

var _ Scene = (*StorageRoom_Scene)(nil)
