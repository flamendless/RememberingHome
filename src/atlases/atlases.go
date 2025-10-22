package atlases

import (
	"remembering-home/src/common"
	"remembering-home/src/enums"
)

type FrameData struct {
	ID   enums.Item
	Pos  common.Vec2
	Size common.Vec2
}

type Metadata struct {
	Padding   int
	Extrude   int
	Size      common.Vec2
	QuadCount int
}

type AtlasData struct {
	Frames   []FrameData
	Metadata Metadata
}
