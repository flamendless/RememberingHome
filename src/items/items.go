package items

import (
	"remembering-home/src/common"
	"remembering-home/src/dialogues"
	"remembering-home/src/enums"
)

type ItemData struct {
	Item      enums.Item
	Name      string
	Pos       common.Vec2
	Z         int
	Scale     float64
	NoCol     bool
	ReqColDir enums.Direction
	Tags      []enums.ItemClass
	Dialogue  dialogues.DialogueKeys
	Grouped   bool
	SubItems  []ItemData
}
