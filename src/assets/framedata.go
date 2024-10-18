package assets

type FrameData struct {
	W, H, MaxCols int
}

var (
	SheetWitsFrameData = FrameData{W: 256, H: 128, MaxCols: 3}
	SheetDeskFrameData = FrameData{W: 256, H: 64, MaxCols: 3}
	BGHallwayFrameData = FrameData{W: 256, H: 64, MaxCols: 1}
)
