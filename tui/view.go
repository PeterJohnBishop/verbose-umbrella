package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// the View method displays the current state of the model

func (m model) View() tea.View {
	var s string

	s += "\n	API TESTING TUI\n\n"

	for i, choice := range m.reqType {

		cursor := " "
		if m.reqCusror == i {
			cursor = ">"
		}

		checked := " "
		if m.reqTypeSelected == i {
			checked = "x"
		}

		s += fmt.Sprintf(" %s [%s] %s ", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return tea.NewView(s)
}
