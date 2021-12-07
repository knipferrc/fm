package ui

import (
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/filetree"
	"github.com/knipferrc/fm/internal/renderer"
	"github.com/knipferrc/fm/internal/statusbar"
	"github.com/knipferrc/fm/internal/theme"
	"github.com/knipferrc/fm/internal/treepreview"
)

// Model represents the state of the UI.
type Model struct {
	fileTree    filetree.Model
	treePreview treepreview.Model
	statusBar   statusbar.Model
	renderer    renderer.Model
	appConfig   config.Config
	theme       theme.Theme
	showPreview bool
}

// NewModel create an instance of the entire application model.
func NewModel() Model {
	cfg := config.GetConfig()
	theme := theme.GetTheme(cfg.Settings.Theme)

	fileTree := filetree.NewModel(
		!cfg.Settings.SimpleMode && cfg.Settings.ShowIcons,
		cfg.Settings.Borderless,
		true,
		true,
		theme.SelectedTreeItemColor,
		theme.UnselectedTreeItemColor,
		theme.ActivePaneBorderColor,
		theme.InactivePaneBorderColor,
		cfg,
	)

	treePreview := treepreview.NewModel(
		!cfg.Settings.SimpleMode && cfg.Settings.ShowIcons,
		cfg.Settings.Borderless,
		true,
		true,
		theme.SelectedTreeItemColor,
		theme.UnselectedTreeItemColor,
		theme.ActivePaneBorderColor,
		theme.InactivePaneBorderColor,
		cfg,
	)

	renderer := renderer.NewModel(
		cfg.Settings.Borderless,
		false,
		theme.ActivePaneBorderColor,
		theme.InactivePaneBorderColor,
	)

	statusBar := statusbar.NewModel(
		statusbar.Color{
			Background: theme.StatusBarSelectedFileBackgroundColor,
			Foreground: theme.StatusBarSelectedFileForegroundColor,
		},
		statusbar.Color{
			Background: theme.StatusBarBarBackgroundColor,
			Foreground: theme.StatusBarBarForegroundColor,
		},
		statusbar.Color{
			Background: theme.StatusBarTotalFilesBackgroundColor,
			Foreground: theme.StatusBarTotalFilesForegroundColor,
		},
		statusbar.Color{
			Background: theme.StatusBarLogoBackgroundColor,
			Foreground: theme.StatusBarLogoForegroundColor,
		},
		!cfg.Settings.SimpleMode && cfg.Settings.ShowIcons,
		cfg.Settings.SimpleMode,
	)

	return Model{
		fileTree:    fileTree,
		treePreview: treePreview,
		statusBar:   statusBar,
		renderer:    renderer,
		appConfig:   cfg,
		theme:       theme,
		showPreview: false,
	}
}
