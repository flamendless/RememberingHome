package atlases

import (
	"remembering-home/src/common"
	"remembering-home/src/dialogues"
	"remembering-home/src/enums"
)

type Data struct {
	Item      enums.Item
	Name      string
	Pos       common.Vec2
	NoCol     bool
	ReqColDir enums.Direction
	Tags      []enums.ItemClass
	Dialogue  dialogues.DialogueKeys
}
