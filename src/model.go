package main

import (
	"io/fs"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

type model struct {
	Files         []fs.FileInfo
	Cursor        int
	Viewport      viewport.Model
	TextInput     textinput.Model
	ViewportReady bool
	ScreenWidth   int
	Move          bool
	Rename        bool
}

func createInitialModel() model {
	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250
	input.Width = 50

	return model{
		make([]fs.FileInfo, 0),
		0,
		viewport.Model{},
		input,
		false,
		0,
		false,
		false,
	}
}
