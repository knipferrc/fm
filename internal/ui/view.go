package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of the entire application UI.
func (m Model) View() string {
	// If the viewport on the panes is not ready display the spinner.
	if !m.ready {
		return fmt.Sprintf("%s%s", m.loader.View(), "loading...")
	}

	horizontalView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.fileTree.View(),
		m.renderer.View(),
	)

	if m.appConfig.Settings.SimpleMode {
		horizontalView = lipgloss.JoinHorizontal(lipgloss.Top, m.fileTree.View())
	}

	// Return the UI with the two panes side by side and
	// the status bar at the bottom of the screen.
	return lipgloss.JoinVertical(
		lipgloss.Top,
		horizontalView,
		m.statusBar.View(),
	)
}
