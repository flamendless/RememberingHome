package conf

import (
	"flag"
	"remembering-home/src/logger"

	"go.uber.org/zap"
)

var (
	DEV = true
)

const (
	GAME_TITLE   = "Remembering Home"
	GAME_VERSION = "v0.0.1"
	WINDOW_W     = 1280
	WINDOW_H     = 640
	GAME_W       = 1280
	GAME_H       = 640
)

type QualityLevel int

const (
	QualityLow QualityLevel = iota
	QualityMedium
	QualityHigh
)

func (q QualityLevel) String() string {
	switch q {
	case QualityLow:
		return "Low"
	case QualityMedium:
		return "Medium"
	case QualityHigh:
		return "High"
	default:
		return "Unknown"
	}
}

func (q QualityLevel) ToShaderValue() float64 {
	switch q {
	case QualityLow:
		return 0.0
	case QualityMedium:
		return 1.0
	case QualityHigh:
		return 2.0
	default:
		return 2.0
	}
}

type WindowMode int

const (
	WindowModeFullscreen WindowMode = iota
	WindowModeWindowed
)

func (w WindowMode) String() string {
	switch w {
	case WindowModeFullscreen:
		return "Fullscreen"
	case WindowModeWindowed:
		return "Windowed"
	default:
		return "Unknown"
	}
}

type Settings struct {
	Quality QualityLevel
	Window  WindowMode
	Volume  int
	Music   int
}

func NewSettings() *Settings {
	window := WindowModeFullscreen
	if DEV {
		window = WindowModeWindowed
	}

	return &Settings{
		Quality: QualityHigh,
		Window:  window,
		Volume:  100,
		Music:   100,
	}
}

func init() {
	flag.BoolVar(&DEV, "dev", false, "developer mode")
	flag.Parse()
}

func Log(settings *Settings) {
	logger.Log().Info(
		"Game Config",
		zap.Bool("dev", DEV),
		zap.String("title", GAME_TITLE),
		zap.String("version", GAME_VERSION),
		zap.Int("window width", WINDOW_W),
		zap.Int("window height", WINDOW_H),
		zap.Int("game width", GAME_W),
		zap.Int("game height", GAME_H),
		zap.String("window mode", settings.Window.String()),
	)
}
