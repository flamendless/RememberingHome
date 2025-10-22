package enums

type Item string

const (
	ItemUndefined      Item = ""
	ItemBulb           Item = "bulb"
	ItemLeftDoor       Item = "left_door"
	ItemRightDoor      Item = "right_door"
	ItemLightSwitch    Item = "light_switch"
	ItemLadder         Item = "ladder"
	ItemShelf          Item = "shelf"
	ItemShelfSide      Item = "shelf_side"
	ItemTable          Item = "table"
	ItemTires          Item = "tires"
	ItemFilingCabinet  Item = "filing_cabinet"
	ItemBarrell        Item = "barrell"
	ItemIroningBoard   Item = "ironing_board"
	ItemWood           Item = "wood"
	ItemWood2          Item = "wood2"
	ItemWood3          Item = "wood3"
	ItemBasket         Item = "basket"
	ItemWashingMachine Item = "washing_machine"
	ItemBroom          Item = "broom"
	ItemElectricalBox  Item = "electrical_box"
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
	case ItemIroningBoard:
		return "ItemIroningBoard"
	case ItemWood:
		return "ItemWood"
	case ItemWood2:
		return "ItemWood2"
	case ItemWood3:
		return "ItemWood3"
	case ItemBasket:
		return "ItemBasket"
	case ItemWashingMachine:
		return "ItemWashingMachine"
	case ItemBroom:
		return "ItemBroom"
	case ItemElectricalBox:
		return "ItemElectricalBox"
	default:
		return "ItemUndefined"
	}
}
