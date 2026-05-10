package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// the View method displays the current state of the model
var (
	labelStyleSelected   = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true) // pink
	labelStyleUnselected = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	unselectedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))            // Gray
	getStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)  // green
	putStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Bold(true) // orange
	patchStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)  // yellow
	postStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)  // blue
	deleteStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)   // red
)

func (m model) View() tea.View {

	str := lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.requestNameInput(), m.methodSelector(), m.footerView())
	return tea.NewView(str)
}

func (m model) headerView() string {
	return "\nAPI Testing\n"
}

func (m model) requestNameInput() string {
	var labelStyle lipgloss.Style

	if m.focus == focusNameInput {
		labelStyle = labelStyleSelected
	} else {
		labelStyle = labelStyleUnselected
	}
	return labelStyle.Render("Request Name") + "\n" + m.reqNameInput.View() + "\n"
}

func (m model) methodSelector() string {
	var labelStyle lipgloss.Style
	if m.focus == focusMethodSelector {
		labelStyle = labelStyleSelected
	} else {
		labelStyle = labelStyleUnselected
	}

	var s strings.Builder

	for i, choice := range m.reqType {
		if m.reqTypeSelected == i {
			switch i {
			case 0:
				s.WriteString(getStyle.Render(choice) + "  ")
			case 1:
				s.WriteString(putStyle.Render(choice) + "  ")
			case 2:
				s.WriteString(patchStyle.Render(choice) + "  ")
			case 3:
				s.WriteString(postStyle.Render(choice) + "  ")
			case 4:
				s.WriteString(deleteStyle.Render(choice) + "  ")
			}
		} else {
			s.WriteString(unselectedStyle.Render(choice) + "  ")
		}
	}

	// 3. Render the Title, add a newline, and then render the horizontal methods
	return labelStyle.Render("Method") + "\n" + s.String()
}

func (m model) footerView() string {
	return "\n\ntab|next shift+tab|previous ctrl+1|name ctrl+2|method ctrl+3|headers ctrl+4|params ctrl+5|body ctrl+6|response [ ctrl+c|quit ]\n"
}
