package scenes

import (
	"remembering-home/src/assets"
	"remembering-home/src/conf"
	"remembering-home/src/context"

	"github.com/hajimehoshi/ebiten/v2"
)

type StorageRoom_Scene struct {
	Context          *context.GameContext
	SceneManager     SceneManager
	CurrentStateName string
	BackgroundImage  *ebiten.Image
	BackgroundDIO    *ebiten.DrawImageOptions
}

func (scene *StorageRoom_Scene) GetName() string {
	return "Storage Room"
}

func (scene *StorageRoom_Scene) GetStateName() string {
	return scene.CurrentStateName
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

	return &scene
}

func (scene *StorageRoom_Scene) Update() error {
	return nil
}

func (scene *StorageRoom_Scene) Draw(screen *ebiten.Image) {
	screen.DrawImage(scene.BackgroundImage, scene.BackgroundDIO)
}

var _ Scene = (*StorageRoom_Scene)(nil)
