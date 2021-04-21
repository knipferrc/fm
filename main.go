package main

import (
	"log"
	"os"

	"github.com/knipferrc/fm/app"
	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/filesystem"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knipferrc/fm/components"
)

func main() {
	config.SetDefaults()
	config.LoadConfig()

	cfg := config.GetConfig()
	m := app.CreateModel()

	if cfg.Settings.StartDir == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		m.Files = filesystem.GetDirectoryListing(home)
	} else {
		m.Files = filesystem.GetDirectoryListing(cfg.Settings.StartDir)
	}

	m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
	m.SecondaryViewport.SetContent(components.Help())
	p := tea.NewProgram(m)

	p.EnableMouseCellMotion()
	defer p.DisableMouseCellMotion()

	if err := p.Start(); err != nil {
		log.Fatal("Failed to start fm", err)
		os.Exit(1)
	}
}