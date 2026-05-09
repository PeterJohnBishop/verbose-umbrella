package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// the View method displays the current state of the model
var (
	unselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")) // Gray
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")). // Pink
			Bold(true)
)

func (m model) View() tea.View {
	var s strings.Builder

	s.WriteString("\n	API TESTING TUI\n\n")

	for i, choice := range m.reqType {

		if m.reqTypeSelected == i {
			s.WriteString(selectedStyle.Render(choice) + "  ")
		} else {
			s.WriteString(unselectedStyle.Render(choice) + "  ")
		}
	}

	s.WriteString("\nPress q to quit.\n")

	return tea.NewView(s.String())
}
