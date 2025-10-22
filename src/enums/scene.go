package enums

type Scene int

const (
	SceneDummy Scene = iota
	SceneSplash
	SceneMainMenu
	SceneStorageRoom
)

func (s Scene) String() string {
	switch s {
	case SceneDummy:
		return "dummy"
	case SceneSplash:
		return "splash"
	case SceneMainMenu:
		return "mainmenu"
	case SceneStorageRoom:
		return "storageroom"
	default:
		return "unknown"
	}
}
