package main

import (
	"github.com/knipferrc/fm/src/components"
	"github.com/knipferrc/fm/src/filesystem"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) fixViewport(moveCursor bool) {
	top := m.Viewport.YOffset
	bottom := m.Viewport.Height + m.Viewport.YOffset - 1

	if moveCursor {
		if m.Cursor < top {
			m.Cursor = top
		} else if m.Cursor > bottom {
			m.Cursor = bottom
		}
		return
	}

	if m.Cursor < top {
		m.Viewport.LineUp(1)
	} else if m.Cursor > bottom {
		m.Viewport.LineDown(1)
	}
}

func (m *model) fixCursor() {
	if m.Cursor > len(m.Files)-1 {
		m.Cursor = 0
	} else if m.Cursor < 0 {
		m.Cursor = len(m.Files) - 1
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		textInputCmd tea.Cmd
		cmds         []tea.Cmd
	)

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()

		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case fileStatus:
		m.Files = msg
		return m, nil
	case tea.WindowSizeMsg:
		m.ScreenWidth = msg.Width
		m.ScreenHeight = msg.Height
		m.Viewport = viewport.Model{
			Width:  msg.Width,
			Height: msg.Height - 1,
		}
		m.Viewport.YPosition = 0

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if !m.TextInput.Focused() {
				m.Cursor--
				m.fixCursor()
				m.fixViewport(false)
			}
		case "down", "j":
			if !m.TextInput.Focused() {
				m.Cursor++
				m.fixCursor()
				m.fixViewport(false)
			}
		case "enter", " ":
			if m.Files[m.Cursor].IsDir() && !m.TextInput.Focused() {
				m.Files = filesystem.GetDirectoryListing(m.Files[m.Cursor].Name())
				m.Cursor = 0
			} else if m.Rename {
				filesystem.RenameDirOrFile(m.Files[m.Cursor].Name(), m.TextInput.Value())
				m.Files = filesystem.GetDirectoryListing("./")
				m.TextInput.Blur()
				m.Rename = false
			} else if m.Move {
				if m.Files[m.Cursor].IsDir() {
					filesystem.MoveDir(m.Files[m.Cursor].Name(), m.TextInput.Value())
					m.Files = filesystem.GetDirectoryListing("./")
					m.TextInput.Blur()
					m.Move = false
				} else {
					filesystem.CopyFile(m.Files[m.Cursor].Name(), m.TextInput.Value(), true)
					m.Files = filesystem.GetDirectoryListing("./")
					m.TextInput.Blur()
					m.Move = false
				}
			} else if m.Delete {
				if m.Files[m.Cursor].IsDir() {
					if m.TextInput.Value() == "y" {
						filesystem.DeleteDirectory(m.Files[m.Cursor].Name())
						m.Files = filesystem.GetDirectoryListing("./")
						m.TextInput.Blur()
						m.Delete = false
					} else {
						m.Files = filesystem.GetDirectoryListing("./")
						m.TextInput.Blur()
						m.Delete = false
					}
				} else {
					if m.TextInput.Value() == "y" {
						filesystem.DeleteFile(m.Files[m.Cursor].Name())
						m.Files = filesystem.GetDirectoryListing("./")
						m.TextInput.Blur()
						m.Delete = false
					} else {
						m.Files = filesystem.GetDirectoryListing("./")
						m.TextInput.Blur()
						m.Delete = false
					}
				}
			} else {
				return m, nil
			}
		case "h", "backspace":
			if !m.TextInput.Focused() {
				m.Cursor = 0
				m.Files = filesystem.GetDirectoryListing("..")
			}
		case "m":
			if !m.TextInput.Focused() {
				m.Move = true
				m.TextInput.Placeholder = "/usr/share/"
				m.TextInput.Focus()
			}
		case "r":
			if !m.TextInput.Focused() {
				m.Rename = true
				m.TextInput.Placeholder = "newfilename.ex"
				m.TextInput.Focus()
			}
		case "d":
			if !m.TextInput.Focused() {
				m.Delete = true
				m.TextInput.Placeholder = "[y/n]"
				m.TextInput.Focus()
			}
		case "i":
			m.Viewport.SetContent(components.Help())
			m.ShowHelp = true
		case "esc":
			m.Move = false
			m.Rename = false
			m.Delete = false
			m.ShowHelp = false
			m.TextInput.Blur()
		}
	}

	m.TextInput, textInputCmd = m.TextInput.Update(msg)
	cmds = append(cmds, textInputCmd)

	return m, tea.Batch(cmds...)
}
