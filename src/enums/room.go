package enums

type Room string

const (
	RoomUndefined   Room = ""
	RoomStorageRoom Room = "storage_room"
	RoomUtilityRoom Room = "utility_room"
)

func (r Room) Constant() string {
	switch r {
	case RoomStorageRoom:
		return "RoomStorageRoom"
	case RoomUtilityRoom:
		return "RoomUtilityRoom"
	default:
		return "RoomUndefined"
	}
}
