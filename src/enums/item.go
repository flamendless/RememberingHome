package enums

type Item string

const (
	ItemUndefined     Item = ""
	ItemBulb          Item = "bulb"
	ItemLeftDoor      Item = "left_door"
	ItemRightDoor     Item = "right_door"
	ItemLightSwitch   Item = "light_switch"
	ItemLadder        Item = "ladder"
	ItemShelf         Item = "shelf"
	ItemShelfSide     Item = "shelf_side"
	ItemTable         Item = "table"
	ItemTires         Item = "tires"
	ItemFilingCabinet Item = "filing_cabinet"
	ItemBarrell       Item = "barrell"
)

func (i Item) Constant() string {
	switch i {
	case ItemBulb:
		return "ItemBulb"
	case ItemLeftDoor:
		return "ItemLeftDoor"
	case ItemRightDoor:
		return "ItemRightDoor"
	case ItemLightSwitch:
		return "ItemLightSwitch"
	case ItemLadder:
		return "ItemLadder"
	case ItemShelf:
		return "ItemShelf"
	case ItemShelfSide:
		return "ItemShelfSide"
	case ItemTable:
		return "ItemTable"
	case ItemTires:
		return "ItemTires"
	case ItemFilingCabinet:
		return "ItemFilingCabinet"
	case ItemBarrell:
		return "ItemBarrell"
	default:
		return "ItemUndefined"
	}
}
