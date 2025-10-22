package enums

type ItemClass string

const (
	ItemClassUndefined ItemClass = ""
	ItemClassDoor      ItemClass = "door"
	ItemClassStair     ItemClass = "stair"
)

func (ic ItemClass) Constant() string {
	switch ic {
	case ItemClassDoor:
		return "ItemClassDoor"
	case ItemClassStair:
		return "ItemClassStair"
	default:
		return "ItemClassUndefined"
	}
}
