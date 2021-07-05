package statusbar

import (
	"github.com/charmbracelet/lipgloss"
)

type Color struct {
	Background string
	Foreground string
}

type Model struct {
	Width               int
	FirstColumnContent  string
	SecondColumnContent string
	ThirdColumnContent  string
	FourthColumnContent string
	FirstColumnColors   Color
	SecondColumnColors  Color
	ThirdColumnColors   Color
	FourthColumnColors  Color
}

// Create a new instance of a status bar
func NewModel(firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors Color) Model {
	return Model{
		FirstColumnColors:  firstColumnColors,
		SecondColumnColors: secondColumnColors,
		ThirdColumnColors:  thirdColumnColors,
		FourthColumnColors: fourthColumnColors,
	}
}

// Set the content of the 4 colums of the status bar
func (m *Model) SetContent(firstColumnContent, secondColumnContent, thirdColumnContent, fourthColumnContent string) {
	m.FirstColumnContent = firstColumnContent
	m.SecondColumnContent = secondColumnContent
	m.ThirdColumnContent = thirdColumnContent
	m.FourthColumnContent = fourthColumnContent
}

// Set the size of the status bar, useful for when screen size changes
func (m *Model) SetSize(width int) {
	m.Width = width
}

// Return the statusbar and all its content
func (m Model) View() string {
	width := lipgloss.Width

	// First column of the status bar displayed on the left with configurable
	// foreground and background colors and some padding
	firstColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.FirstColumnColors.Foreground)).
		Background(lipgloss.Color(m.FirstColumnColors.Background)).
		Padding(0, 1).
		Render(m.FirstColumnContent)

	// Third column of the status bar displayed on the left with configurable
	// foreground and background colors and some padding
	thirdColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.ThirdColumnColors.Foreground)).
		Background(lipgloss.Color(m.ThirdColumnColors.Background)).
		Align(lipgloss.Right).
		Padding(0, 1).
		Render(m.ThirdColumnContent)

	// Fourth column of the status bar displayed on the left with configurable
	// foreground and background colors and some padding
	fourthColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.FourthColumnColors.Foreground)).
		Background(lipgloss.Color(m.FourthColumnColors.Background)).
		Padding(0, 1).
		Render(m.FourthColumnContent)

	// Second column of the status bar displayed on the left with configurable
	// foreground and background colors and some padding. Also calculate the
	// width of the other three columns so that this one can take up the rest of the space
	// in the center of the bar
	secondColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.SecondColumnColors.Foreground)).
		Background(lipgloss.Color(m.SecondColumnColors.Background)).
		Padding(0, 1).
		Width(m.Width - width(firstColumn) - width(thirdColumn) - width(fourthColumn)).
		Render(m.SecondColumnContent)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		firstColumn,
		secondColumn,
		thirdColumn,
		fourthColumn,
	)
}