package graphics

import (
	"fmt"
	"image"
	"image/color"
	"remembering-home/src/assets"
	"remembering-home/src/atlases"
	"remembering-home/src/common"
	"remembering-home/src/context"
	"remembering-home/src/enums"
	"remembering-home/src/items"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	resource "github.com/quasilyte/ebitengine-resource"
)

type ItemRenderer struct {
	AtlasTexture *ebiten.Image
	AtlasData    atlases.AtlasData
	FrameMap     map[enums.Item]*ebiten.Image
	FrameDataMap map[enums.Item]atlases.FrameData
	Items        []items.ItemData
	ItemMap      map[enums.Item]items.ItemData
	DIO          *ebiten.DrawImageOptions
	debugFace    *text.GoXFace
	debugText    *text.DrawOptions
}

func NewItemRenderer(ctx *context.GameContext, atlasImageID resource.ImageID, atlasData atlases.AtlasData, itemsData []items.ItemData) *ItemRenderer {
	atlasTexture := ctx.Loader.LoadImage(atlasImageID).Data

	frameMap := make(map[enums.Item]*ebiten.Image)
	frameDataMap := make(map[enums.Item]atlases.FrameData)
	itemMap := make(map[enums.Item]items.ItemData)

	for _, frame := range atlasData.Frames {
		frameRect := image.Rect(
			int(frame.Pos.X), int(frame.Pos.Y),
			int(frame.Pos.X+frame.Size.X), int(frame.Pos.Y+frame.Size.Y),
		)
		subTexture := atlasTexture.SubImage(frameRect).(*ebiten.Image)
		frameMap[frame.ID] = subTexture
		frameDataMap[frame.ID] = frame
	}

	for _, item := range itemsData {
		itemMap[item.Item] = item
	}

	debugFace := text.NewGoXFace(ctx.Loader.LoadFont(assets.FontJamboree18).Face)
	debugText := &text.DrawOptions{DrawImageOptions: ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}}
	debugText.ColorScale.SetR(1.0)
	debugText.ColorScale.SetG(1.0)
	debugText.ColorScale.SetB(1.0)
	debugText.ColorScale.SetA(1.0)
	debugText.PrimaryAlign = text.AlignCenter
	debugText.SecondaryAlign = text.AlignEnd

	return &ItemRenderer{
		AtlasTexture: atlasTexture,
		AtlasData:    atlasData,
		FrameMap:     frameMap,
		FrameDataMap: frameDataMap,
		Items:        itemsData,
		ItemMap:      itemMap,
		DIO: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterNearest,
		},
		debugFace: debugFace,
		debugText: debugText,
	}
}

func (ir *ItemRenderer) DrawAllItems(screen *ebiten.Image, gameScale float64, gameOffset common.Vec2) {
	for _, itemData := range ir.Items {
		ir.DrawItem(screen, itemData, gameScale, gameOffset)
	}
}

func (ir *ItemRenderer) DrawItem(screen *ebiten.Image, itemData items.ItemData, gameScale float64, gameOffset common.Vec2) {
	subTexture, exists := ir.FrameMap[itemData.Item]
	if !exists {
		return
	}

	finalPos := itemData.Pos.Add(gameOffset).Mult(gameScale)
	finalScale := gameScale * itemData.Scale

	ir.DIO.GeoM.Reset()
	ir.DIO.GeoM.Scale(finalScale, finalScale)
	ir.DIO.GeoM.Translate(finalPos.X, finalPos.Y)
	screen.DrawImage(subTexture, ir.DIO)
}

func (ir *ItemRenderer) GetItemAtPosition(screenX, screenY float64, gameScale float64, gameOffset common.Vec2) *items.ItemData {
	for i := len(ir.Items) - 1; i >= 0; i-- {
		itemData := ir.Items[i]
		frameData, exists := ir.FrameDataMap[itemData.Item]
		if !exists {
			continue
		}

		finalPos := itemData.Pos.Add(gameOffset).Mult(gameScale)
		finalScale := gameScale * itemData.Scale
		itemWidth := float64(frameData.Size.X) * finalScale
		itemHeight := float64(frameData.Size.Y) * finalScale

		if screenX >= finalPos.X && screenX <= finalPos.X+itemWidth &&
			screenY >= finalPos.Y && screenY <= finalPos.Y+itemHeight {
			return &itemData
		}
	}
	return nil
}

func (ir *ItemRenderer) DrawDebug(screen *ebiten.Image, itemData items.ItemData, gameScale float64, gameOffset common.Vec2) {
	frameData, exists := ir.FrameDataMap[itemData.Item]
	if !exists {
		return
	}

	finalPos := itemData.Pos.Add(gameOffset).Mult(gameScale)
	finalScale := gameScale * itemData.Scale
	itemWidth := float64(frameData.Size.X) * finalScale
	itemHeight := float64(frameData.Size.Y) * finalScale

	redColor := color.RGBA{255, 0, 0, 255}
	vector.StrokeRect(screen, float32(finalPos.X), float32(finalPos.Y), float32(itemWidth), float32(itemHeight), 2, redColor, false)

	var displayText string
	if itemData.Name != "" {
		displayText = itemData.Name
	} else {
		displayText = string(itemData.Item)
	}

	textX := finalPos.X + itemWidth/2
	textY := finalPos.Y - 5

	ir.debugText.GeoM.Reset()
	ir.debugText.GeoM.Translate(textX, textY)
	text.Draw(screen, displayText, ir.debugFace, ir.debugText)

	ir.debugText.GeoM.Reset()
	ir.debugText.GeoM.Translate(textX, textY+16)
	text.Draw(screen, fmt.Sprintf("(X: %.0f, Y: %.0f)", finalPos.X, finalPos.Y), ir.debugFace, ir.debugText)
}

func (ir *ItemRenderer) DrawItemSelection(screen *ebiten.Image, gameScale float64, gameOffset common.Vec2, zoomLevel float64) {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorXFloat := float64(cursorX)
	cursorYFloat := float64(cursorY)

	if zoomLevel != 1.0 {
		gameW, gameH := 1920, 1080 // Use constants instead of conf package to avoid import cycle
		offsetX := (float64(gameW) - float64(gameW)*zoomLevel) / 2
		offsetY := (float64(gameH) - float64(gameH)*zoomLevel) / 2
		cursorXFloat = (cursorXFloat - offsetX) / zoomLevel
		cursorYFloat = (cursorYFloat - offsetY) / zoomLevel
	}

	hoveredItem := ir.GetItemAtPosition(cursorXFloat, cursorYFloat, gameScale, gameOffset)
	if hoveredItem != nil {
		ir.DrawDebug(screen, *hoveredItem, gameScale, gameOffset)
	}
}

func (ir *ItemRenderer) GetItemCount() int {
	return len(ir.FrameMap)
}

func (ir *ItemRenderer) GetItems() []items.ItemData {
	return ir.Items
}
