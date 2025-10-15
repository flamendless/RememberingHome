package ui

import (
	"fmt"
	"remembering-home/src/assets/shaders"
	"remembering-home/src/conf"
	"remembering-home/src/context"
	"remembering-home/src/graphics"
	"remembering-home/src/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	input "github.com/quasilyte/ebitengine-input"
	"golang.org/x/image/font"
)

const HOLD_THRESHOLD = 9

type MenuLevel int

const (
	MenuLevelMain MenuLevel = iota
	MenuLevelSubmenu
)

const (
	SettingsIdxGraphics = iota
	SettingsIdxAudio
	SettingsIdxBack
)

const (
	GraphicsIdxQuality = iota
	GraphicsIdxWindow
	GraphicsIdxBack
)

const (
	AudioIdxVolume = iota
	AudioIdxMusic
	AudioIdxBack
)

type SettingsMenu struct {
	Context             *context.GameContext
	MainList            *TextList
	GraphicsSubmenu     *TextList
	AudioSubmenu        *TextList
	CurrentLevel        MenuLevel
	ActiveSubmenu       int
	Enabled             bool
	Visible             bool
	backRequested       bool
	submenuBaseX        float64
	dimAlpha            float32
	mainListDimmedAlpha float32
	primaryFace         *font.Face
	secondaryFace       *font.Face
	leftHoldFrames      int
	rightHoldFrames     int
	holdThreshold       int
}

type SettingsMenuConfig struct {
	Context        *context.GameContext
	BaseX          float64
	BaseY          float64
	SubmenuOffsetX float64
	Gap            float64
	BannerConfig   BannerConfig
	PrimaryFace    *font.Face
	SecondaryFace  *font.Face
}

func NewSettingsMenu(config SettingsMenuConfig) *SettingsMenu {
	sm := &SettingsMenu{
		Context:             config.Context,
		CurrentLevel:        MenuLevelMain,
		ActiveSubmenu:       -1,
		Enabled:             true,
		Visible:             true,
		submenuBaseX:        config.BaseX + config.SubmenuOffsetX,
		dimAlpha:            0.5,
		mainListDimmedAlpha: 0.6,
		primaryFace:         config.PrimaryFace,
		secondaryFace:       config.SecondaryFace,
		leftHoldFrames:      0,
		rightHoldFrames:     0,
		holdThreshold:       15,
	}

	sm.MainList = NewTextList(TextListConfig{
		Layout:       LayoutVertical,
		BaseX:        config.BaseX,
		BaseY:        config.BaseY,
		Gap:          config.Gap,
		Alignment:    text.AlignStart,
		SecAlignment: text.AlignStart,
		BannerConfig: config.BannerConfig,
	})

	for _, txt := range []string{"Graphics", "Audio", "Back"} {
		txtMenu := graphics.NewText(config.PrimaryFace, txt, true)
		utils.SetColor(txtMenu.DO, 1, 1, 1, 1)
		sm.MainList.AddText(txtMenu)
	}

	sm.MainList.SetBannerWidth(sm.MainList.GetBannerWidth() * 1.5)

	sm.GraphicsSubmenu = NewTextList(TextListConfig{
		Layout:       LayoutVertical,
		BaseX:        sm.submenuBaseX,
		BaseY:        config.BaseY,
		Gap:          config.Gap,
		Alignment:    text.AlignStart,
		SecAlignment: text.AlignStart,
		BannerConfig: config.BannerConfig,
	})

	sm.GraphicsSubmenu.SetVisible(true)
	sm.GraphicsSubmenu.SetEnabled(false)
	sm.addGraphicsSubmenuTexts()

	sm.AudioSubmenu = NewTextList(TextListConfig{
		Layout:       LayoutVertical,
		BaseX:        sm.submenuBaseX,
		BaseY:        config.BaseY,
		Gap:          config.Gap,
		Alignment:    text.AlignStart,
		SecAlignment: text.AlignStart,
		BannerConfig: config.BannerConfig,
	})

	sm.AudioSubmenu.SetVisible(true)
	sm.AudioSubmenu.SetEnabled(false)
	sm.addAudioSubmenuTexts()

	return sm
}

func (sm *SettingsMenu) addGraphicsSubmenuTexts() {
	qualityText := sm.formatCyclingOption("Quality", sm.Context.Settings.Quality.String())
	txtQuality := graphics.NewText(sm.secondaryFace, qualityText, true)
	utils.SetColor(txtQuality.DO, 1, 1, 1, 1)
	sm.GraphicsSubmenu.AddText(txtQuality)

	windowText := sm.formatCyclingOption("Window", sm.Context.Settings.Window.String())
	txtWindow := graphics.NewText(sm.secondaryFace, windowText, true)
	utils.SetColor(txtWindow.DO, 1, 1, 1, 1)
	sm.GraphicsSubmenu.AddText(txtWindow)

	txtBack := graphics.NewText(sm.secondaryFace, "Back", true)
	utils.SetColor(txtBack.DO, 1, 1, 1, 1)
	sm.GraphicsSubmenu.AddText(txtBack)

	sm.GraphicsSubmenu.SetBannerWidth(sm.GraphicsSubmenu.GetBannerWidth() * 1.5)
}

func (sm *SettingsMenu) addAudioSubmenuTexts() {
	volumeText := sm.formatCyclingOption("Volume", fmt.Sprintf("%d", sm.Context.Settings.Volume))
	txtVolume := graphics.NewText(sm.secondaryFace, volumeText, true)
	utils.SetColor(txtVolume.DO, 1, 1, 1, 1)
	sm.AudioSubmenu.AddText(txtVolume)

	musicText := sm.formatCyclingOption("Music", fmt.Sprintf("%d", sm.Context.Settings.Music))
	txtMusic := graphics.NewText(sm.secondaryFace, musicText, true)
	utils.SetColor(txtMusic.DO, 1, 1, 1, 1)
	sm.AudioSubmenu.AddText(txtMusic)

	txtBack := graphics.NewText(sm.secondaryFace, "Back", true)
	utils.SetColor(txtBack.DO, 1, 1, 1, 1)
	sm.AudioSubmenu.AddText(txtBack)

	sm.AudioSubmenu.SetBannerWidth(sm.AudioSubmenu.GetBannerWidth() * 1.5)
}

func (sm *SettingsMenu) formatCyclingOption(label, value string) string {
	return fmt.Sprintf("%s: < %s >", label, value)
}

func (sm *SettingsMenu) updateGraphicsSubmenuTexts() {
	if len(sm.GraphicsSubmenu.Texts) < 2 {
		return
	}

	sm.GraphicsSubmenu.Texts[GraphicsIdxQuality].Txt = sm.formatCyclingOption("Quality", sm.Context.Settings.Quality.String())
	sm.GraphicsSubmenu.Texts[GraphicsIdxWindow].Txt = sm.formatCyclingOption("Window", sm.Context.Settings.Window.String())
}

func (sm *SettingsMenu) cycleQuality(forward bool) {
	if forward {
		sm.Context.Settings.Quality = (sm.Context.Settings.Quality + 1) % 3
	} else {
		sm.Context.Settings.Quality = (sm.Context.Settings.Quality + 2) % 3
	}
	sm.updateGraphicsSubmenuTexts()
}

func (sm *SettingsMenu) cycleWindowMode(forward bool) {
	if forward {
		sm.Context.Settings.Window = (sm.Context.Settings.Window + 1) % 2
	} else {
		sm.Context.Settings.Window = (sm.Context.Settings.Window + 1) % 2
	}
	sm.updateGraphicsSubmenuTexts()
	sm.applyWindowMode()
}

func (sm *SettingsMenu) applyWindowMode() {
	ebiten.SetFullscreen(sm.Context.Settings.Window == conf.WindowModeFullscreen)
}

func (sm *SettingsMenu) adjustAudio(setting *int, amount int) {
	*setting = utils.ClampInt(*setting+amount, 0, 100)
	sm.AudioSubmenu.Texts[AudioIdxVolume].Txt = sm.formatCyclingOption("Volume", fmt.Sprintf("%d", sm.Context.Settings.Volume))
	sm.AudioSubmenu.Texts[AudioIdxMusic].Txt = sm.formatCyclingOption("Music", fmt.Sprintf("%d", sm.Context.Settings.Music))
}

func (sm *SettingsMenu) Update(inputHandler *input.Handler) {
	sm.backRequested = false

	if !sm.Enabled || !sm.Visible {
		return
	}

	switch sm.CurrentLevel {
	case MenuLevelMain:
		sm.updateMainMenu(inputHandler)
	case MenuLevelSubmenu:
		sm.updateSubmenu(inputHandler)
	}
}

func (sm *SettingsMenu) updateMainMenu(inputHandler *input.Handler) {
	sm.MainList.Update(inputHandler)

	if inputHandler.ActionIsJustPressed(context.ActionBack) {
		sm.backRequested = true
		return
	}

	if sm.MainList.IsConfirmed() {
		switch sm.MainList.GetSelectedIndex() {
		case SettingsIdxGraphics:
			sm.enterSubmenu(SettingsIdxGraphics)
		case SettingsIdxAudio:
			sm.enterSubmenu(SettingsIdxAudio)
		case SettingsIdxBack:
			sm.backRequested = true
		}
	}
}

func (sm *SettingsMenu) updateSubmenu(inputHandler *input.Handler) {
	var currentSubmenu *TextList
	switch sm.ActiveSubmenu {
	case SettingsIdxGraphics:
		currentSubmenu = sm.GraphicsSubmenu
	case SettingsIdxAudio:
		currentSubmenu = sm.AudioSubmenu
	default:
		return
	}

	if inputHandler.ActionIsJustPressed(context.ActionBack) {
		sm.exitSubmenu()
		return
	}

	currentSubmenu.Update(inputHandler)

	if sm.ActiveSubmenu == SettingsIdxGraphics {
		selectedIdx := currentSubmenu.GetSelectedIndex()

		if inputHandler.ActionIsJustReleased(context.ActionMoveLeft) {
			switch selectedIdx {
			case GraphicsIdxQuality:
				sm.cycleQuality(false)
			case GraphicsIdxWindow:
				sm.cycleWindowMode(false)
			}
		} else if inputHandler.ActionIsJustReleased(context.ActionMoveRight) {
			switch selectedIdx {
			case GraphicsIdxQuality:
				sm.cycleQuality(true)
			case GraphicsIdxWindow:
				sm.cycleWindowMode(true)
			}
		}
	} else if sm.ActiveSubmenu == SettingsIdxAudio {
		selectedIdx := currentSubmenu.GetSelectedIndex()
		leftPressed := inputHandler.ActionIsPressed(context.ActionMoveLeft)
		rightPressed := inputHandler.ActionIsPressed(context.ActionMoveRight)

		if leftPressed && rightPressed {
			sm.resetHoldFrames()
		} else if leftPressed {
			sm.leftHoldFrames++

			if sm.leftHoldFrames == 1 {
				switch selectedIdx {
				case AudioIdxVolume:
					sm.adjustAudio(&sm.Context.Settings.Volume, -1)
				case AudioIdxMusic:
					sm.adjustAudio(&sm.Context.Settings.Music, -1)
				}
			} else if sm.leftHoldFrames > sm.holdThreshold && sm.leftHoldFrames%HOLD_THRESHOLD == 0 {
				switch selectedIdx {
				case AudioIdxVolume:
					sm.adjustAudio(&sm.Context.Settings.Volume, -10)
				case AudioIdxMusic:
					sm.adjustAudio(&sm.Context.Settings.Music, -10)
				}
			}
		} else if rightPressed {
			sm.rightHoldFrames++

			if sm.rightHoldFrames == 1 {
				switch selectedIdx {
				case AudioIdxVolume:
					sm.adjustAudio(&sm.Context.Settings.Volume, 1)
				case AudioIdxMusic:
					sm.adjustAudio(&sm.Context.Settings.Music, 1)
				}
			} else if sm.rightHoldFrames > sm.holdThreshold && sm.rightHoldFrames%HOLD_THRESHOLD == 0 {
				switch selectedIdx {
				case AudioIdxVolume:
					sm.adjustAudio(&sm.Context.Settings.Volume, 10)
				case AudioIdxMusic:
					sm.adjustAudio(&sm.Context.Settings.Music, 10)
				}
			}
		} else {
			sm.resetHoldFrames()
		}
	}

	if currentSubmenu.IsConfirmed() {
		selectedIdx := currentSubmenu.GetSelectedIndex()

		if selectedIdx == len(currentSubmenu.Texts)-1 {
			sm.exitSubmenu()
		}
	}
}

func (sm *SettingsMenu) enterSubmenu(submenuIdx int) {
	sm.CurrentLevel = MenuLevelSubmenu
	sm.ActiveSubmenu = submenuIdx

	switch submenuIdx {
	case SettingsIdxGraphics:
		sm.GraphicsSubmenu.SetEnabled(true)
		sm.GraphicsSubmenu.SetCurrentIndex(0)
	case SettingsIdxAudio:
		sm.AudioSubmenu.SetEnabled(true)
		sm.AudioSubmenu.SetCurrentIndex(0)
	}
}

func (sm *SettingsMenu) exitSubmenu() {
	sm.CurrentLevel = MenuLevelMain

	switch sm.ActiveSubmenu {
	case SettingsIdxGraphics:
		sm.GraphicsSubmenu.SetEnabled(false)
	case SettingsIdxAudio:
		sm.AudioSubmenu.SetEnabled(false)
	}

	sm.ActiveSubmenu = -1
	sm.resetHoldFrames()
}

func (sm *SettingsMenu) resetHoldFrames() {
	sm.leftHoldFrames = 0
	sm.rightHoldFrames = 0
}

func (sm *SettingsMenu) IsBackRequested() bool {
	return sm.backRequested
}

func (sm *SettingsMenu) SetEnabled(enabled bool) {
	sm.Enabled = enabled
}

func (sm *SettingsMenu) SetVisible(visible bool) {
	sm.Visible = visible
}

func (sm *SettingsMenu) Reset() {
	sm.CurrentLevel = MenuLevelMain
	sm.MainList.SetCurrentIndex(0)
	sm.exitSubmenu()
}

func (sm *SettingsMenu) ApplyShaderAlpha(uniforms *shaders.SilentHillRedShaderUniforms) {
	if !sm.Visible {
		sm.MainList.ApplyShaderAlpha(uniforms)
		sm.GraphicsSubmenu.ApplyShaderAlpha(uniforms)
		sm.AudioSubmenu.ApplyShaderAlpha(uniforms)
		return
	}

	textAlpha := uniforms.GetTextAlpha()

	for i, txt := range sm.MainList.Texts {
		if sm.CurrentLevel == MenuLevelSubmenu {
			if i != sm.MainList.GetSelectedIndex() {
				utils.SetColor(txt.DO, 1, 1, 1, textAlpha*sm.dimAlpha)
			} else {
				utils.SetColor(txt.DO, 1, 0.2, 0.2, textAlpha)
			}
		} else {
			utils.SetColor(txt.DO, 1, 1, 1, textAlpha)
		}
	}

	selectedIdx := sm.MainList.GetSelectedIndex()
	var previewSubmenu *TextList
	var isPreview bool

	if sm.CurrentLevel == MenuLevelMain {
		isPreview = true
		switch selectedIdx {
		case SettingsIdxGraphics:
			previewSubmenu = sm.GraphicsSubmenu
		case SettingsIdxAudio:
			previewSubmenu = sm.AudioSubmenu
		}
	} else {
		isPreview = false
		switch sm.ActiveSubmenu {
		case SettingsIdxGraphics:
			previewSubmenu = sm.GraphicsSubmenu
		case SettingsIdxAudio:
			previewSubmenu = sm.AudioSubmenu
		}
	}

	if previewSubmenu != nil {
		for _, txt := range previewSubmenu.Texts {
			if isPreview {
				utils.DOReplaceAlpha(txt.DO, textAlpha*sm.dimAlpha)
			} else {
				utils.DOReplaceAlpha(txt.DO, textAlpha)
			}
		}
	}
}

func (sm *SettingsMenu) UpdateBannerPosition(uniforms *shaders.SilentHillRedShaderUniforms) {
}

func (sm *SettingsMenu) Draw(canvas *ebiten.Image) {
	if !sm.Visible {
		return
	}

	sm.MainList.Draw(canvas)

	if sm.CurrentLevel == MenuLevelMain {
		selectedIdx := sm.MainList.GetSelectedIndex()
		switch selectedIdx {
		case SettingsIdxGraphics:
			sm.GraphicsSubmenu.Draw(canvas)
		case SettingsIdxAudio:
			sm.AudioSubmenu.Draw(canvas)
		}
	} else {
		switch sm.ActiveSubmenu {
		case SettingsIdxGraphics:
			sm.GraphicsSubmenu.Draw(canvas)
		case SettingsIdxAudio:
			sm.AudioSubmenu.Draw(canvas)
		}
	}
}

func (sm *SettingsMenu) GetMainList() *TextList {
	return sm.MainList
}

func (sm *SettingsMenu) GetActiveSubmenu() *TextList {
	if sm.CurrentLevel == MenuLevelMain {
		selectedIdx := sm.MainList.GetSelectedIndex()
		switch selectedIdx {
		case SettingsIdxGraphics:
			return sm.GraphicsSubmenu
		case SettingsIdxAudio:
			return sm.AudioSubmenu
		default:
			return nil
		}
	}

	switch sm.ActiveSubmenu {
	case SettingsIdxGraphics:
		return sm.GraphicsSubmenu
	case SettingsIdxAudio:
		return sm.AudioSubmenu
	default:
		return nil
	}
}

func (sm *SettingsMenu) GetCurrentLevel() MenuLevel {
	return sm.CurrentLevel
}
