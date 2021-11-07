package statusbar

import (
	"fmt"
	"os"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/icons"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Color is a struct that contains the foreground and background colors of the statusbar.
type Color struct {
	Background lipgloss.AdaptiveColor
	Foreground lipgloss.AdaptiveColor
}

// Model is a struct that contains all the properties of the statusbar.
type Model struct {
	Width              int
	Height             int
	TotalFiles         int
	Cursor             int
	ShowIcons          bool
	ShowCommandBar     bool
	InMoveMode         bool
	ItemSize           string
	FilePaths          []string
	SelectedFile       os.FileInfo
	ItemToMove         os.FileInfo
	FirstColumnColors  Color
	SecondColumnColors Color
	ThirdColumnColors  Color
	FourthColumnColors Color
	Textinput          textinput.Model
	Spinner            spinner.Model
}

// NewModel creates an instance of a statusbar.
func NewModel(
	firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors Color, showIcons bool,
) Model {
	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250
	input.Placeholder = "enter a name"
	input.PlaceholderStyle.Background(secondColumnColors.Background)

	s := spinner.NewModel()
	s.Spinner = spinner.Dot

	return Model{
		Height:             1,
		TotalFiles:         0,
		Cursor:             0,
		ShowIcons:          showIcons,
		ShowCommandBar:     false,
		InMoveMode:         false,
		ItemSize:           "",
		SelectedFile:       nil,
		ItemToMove:         nil,
		FirstColumnColors:  firstColumnColors,
		SecondColumnColors: secondColumnColors,
		ThirdColumnColors:  thirdColumnColors,
		FourthColumnColors: fourthColumnColors,
		Textinput:          input,
		Spinner:            s,
	}
}

// GetHeight returns the height of the statusbar.
func (m Model) GetHeight() int {
	return m.Height
}

// BlurCommandBar blurs the textinput used for the command bar.
func (m *Model) BlurCommandBar() {
	m.Textinput.Blur()
}

// ResetCommandBar resets the textinput used for the command bar.
func (m *Model) ResetCommandBar() {
	m.Textinput.Reset()
}

// CommandBarValue returns the value of the command bar.
func (m Model) CommandBarValue() string {
	return m.Textinput.Value()
}

// CommandBarFocused returns true if the command bar is focused.
func (m Model) CommandBarFocused() bool {
	return m.Textinput.Focused()
}

// FocusCommandBar focuses the command bar.
func (m *Model) FocusCommandBar() {
	m.Textinput.Focus()
}

// SetContent sets the content of the statusbar.
func (m *Model) SetContent(
	totalFiles, cursor int,
	showCommandBar, inMoveMode bool,
	selectedFile, itemToMove os.FileInfo, filePaths []string,
) {
	m.TotalFiles = totalFiles
	m.Cursor = cursor
	m.ShowCommandBar = showCommandBar
	m.InMoveMode = inMoveMode
	m.SelectedFile = selectedFile
	m.ItemToMove = itemToMove
	m.FilePaths = filePaths
}

// SetItemSize sets the size of the currently selected
// directory item as a formatted size string.
func (m *Model) SetItemSize(itemSize string) {
	m.ItemSize = itemSize
}

// SetSize sets the size of the statusbar, useful when the terminal is resized.
func (m *Model) SetSize(width int) {
	m.Width = width
}

// Update updates the statusbar.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width)
	}

	m.Textinput, cmd = m.Textinput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the statusbar.
func (m Model) View() string {
	var logo string
	var status string

	width := lipgloss.Width
	selectedFile := "N/A"
	fileCount := "0/0"
	fileSize := m.Spinner.View()

	if m.TotalFiles > 0 && m.SelectedFile != nil {
		selectedFile = m.SelectedFile.Name()
		fileCount = fmt.Sprintf("%d/%d", m.Cursor+1, m.TotalFiles)

		currentPath, err := dirfs.GetWorkingDirectory()
		if err != nil {
			currentPath = dirfs.CurrentDirectory
		}

		if m.ItemSize != "" {
			fileSize = m.ItemSize
		}

		if len(m.FilePaths) > 0 {
			currentPath = m.FilePaths[m.Cursor]
		}

		// Display some information about the currently seleted file including
		// its size, the mode and the current path.
		status = fmt.Sprintf("%s %s %s",
			fileSize,
			m.SelectedFile.Mode().String(),
			currentPath,
		)
	}

	if m.ShowCommandBar {
		status = m.Textinput.View()
	}

	if m.InMoveMode {
		status = fmt.Sprintf("%s %s", "Currently moving:", m.ItemToMove.Name())
	}

	if m.ShowIcons {
		logo = fmt.Sprintf("%s %s", icons.IconDef["dir"].GetGlyph(), "FM")
	} else {
		logo = "FM"
	}

	selectedFileColumn := lipgloss.NewStyle().
		Foreground(m.FirstColumnColors.Foreground).
		Background(m.FirstColumnColors.Background).
		Padding(0, 1).
		Height(m.Height).
		Render(truncate.StringWithTail(selectedFile, 30, "..."))

	fileCountColumn := lipgloss.NewStyle().
		Foreground(m.ThirdColumnColors.Foreground).
		Background(m.ThirdColumnColors.Background).
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(m.Height).
		Render(fileCount)

	logoColumn := lipgloss.NewStyle().
		Foreground(m.FourthColumnColors.Foreground).
		Background(m.FourthColumnColors.Background).
		Padding(0, 1).
		Height(m.Height).
		Render(logo)

	statusColumn := lipgloss.NewStyle().
		Foreground(m.SecondColumnColors.Foreground).
		Background(m.SecondColumnColors.Background).
		Padding(0, 1).
		Height(m.Height).
		Width(m.Width - width(selectedFileColumn) - width(fileCountColumn) - width(logoColumn)).
		Render(truncate.StringWithTail(
			status,
			uint(m.Width-width(selectedFileColumn)-width(fileCountColumn)-width(logoColumn)-3),
			"..."),
		)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFileColumn,
		statusColumn,
		fileCountColumn,
		logoColumn,
	)
}
