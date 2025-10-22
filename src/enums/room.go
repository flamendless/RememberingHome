package enums

type Room string

const (
	RoomUndefined   Room = ""
	RoomStorageRoom Room = "storage_room"
)

func (r Room) Constant() string {
	switch r {
	case RoomStorageRoom:
		return "RoomStorageRoom"
	default:
		return "RoomUndefined"
	}
}
