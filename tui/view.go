package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// the View method displays the current state of the model
var (
	labelStyleSelected    = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true) // pink
	labelStyleUnselected  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))            // gray
	buttonStyleSelected   = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true) // Pink
	buttonStyleUnselected = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	unselectedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	getStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)  // green
	putStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Bold(true) // orange
	patchStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)  // yellow
	postStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)  // blue
	deleteStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)   // red
)

func (m model) View() tea.View {

	str := lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.requestNameInput(), m.methodSelector(), m.requestEndpointInput(), m.headersInputView(), m.headersListView(), m.paramsInputView(), m.paramsListView(), m.footerView())
	return tea.NewView(str)
}

func (m model) headerView() string {
	return "\nAPI Testing\n"
}

func (m model) requestNameInput() string {
	var labelStyle lipgloss.Style

	if m.focus == focusName {
		labelStyle = labelStyleSelected
	} else {
		labelStyle = labelStyleUnselected
	}
	return labelStyle.Render("Name") + "\n" + m.inputs[0].View() + "\n"
}

func (m model) requestEndpointInput() string {
	var labelStyle lipgloss.Style

	if m.focus == focusEndpoint {
		labelStyle = labelStyleSelected
	} else {
		labelStyle = labelStyleUnselected
	}
	return labelStyle.Render("Endpoint") + "\n" + m.inputs[1].View() + "\n"
}

func (m model) methodSelector() string {
	var labelStyle lipgloss.Style
	if m.focus == focusMethod {
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
	return labelStyle.Render("Method") + "\n" + s.String() + "\n"
}

func (m model) headersInputView() string {
	var keyLabelStyle lipgloss.Style

	if m.focus == focusHeaderKey || m.focus == focusHeaderValue {
		keyLabelStyle = labelStyleSelected
	} else {
		keyLabelStyle = labelStyleUnselected
	}

	var buttonView string
	if m.focus == focusHeaderSubmit {
		buttonView = buttonStyleSelected.Render("[ add ]")
	} else {
		buttonView = buttonStyleUnselected.Render("[ add ]")
	}

	keyView := keyLabelStyle.Render("Headers") + "\n" + m.inputs[inputHeadersKeyIdx].View()

	valueView := m.inputs[inputHeadersValueIdx].View()

	key := m.inputs[inputHeadersKeyIdx].Value()
	val := m.inputs[inputHeadersValueIdx].Value()

	if key == "" || val == "" {
		return lipgloss.JoinVertical(lipgloss.Left, keyView, valueView) + "\n"

	} else {
		return lipgloss.JoinVertical(lipgloss.Left, keyView, valueView, buttonView) + "\n"

	}

}

func (m model) headersListView() string {
	// If there are no headers, don't render anything
	if len(m.req.Headers) == 0 {
		return ""
	}

	var s strings.Builder
	s.WriteString("\ndelete || backspace|remove selected\n")

	var keys []string
	for k := range m.req.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		val := m.req.Headers.Get(k)
		row := fmt.Sprintf("%s: %s", k, val)

		if m.focus == focusHeaderList && m.headerCursor == i {
			s.WriteString(labelStyleSelected.Render("> "+row) + "\n")
		} else {
			s.WriteString(labelStyleUnselected.Render("  "+row) + "\n")
		}
	}

	return s.String()
}

func (m model) paramsInputView() string {
	var keyLabelStyle lipgloss.Style

	if m.focus == focusParamKey || m.focus == focusParamValue {
		keyLabelStyle = labelStyleSelected
	} else {
		keyLabelStyle = labelStyleUnselected
	}

	var buttonView string
	if m.focus == focusParamSubit {
		buttonView = buttonStyleSelected.Render("[ add ]")
	} else {
		buttonView = buttonStyleUnselected.Render("[ add ]")
	}

	keyView := keyLabelStyle.Render("Params") + "\n" + m.inputs[inputParamsKeyIdx].View()

	valueView := m.inputs[inputParamsValueIdx].View()

	key := m.inputs[inputParamsKeyIdx].Value()
	val := m.inputs[inputParamsValueIdx].Value()

	if key == "" || val == "" {
		return lipgloss.JoinVertical(lipgloss.Left, keyView, valueView) + "\n"

	} else {
		return lipgloss.JoinVertical(lipgloss.Left, keyView, valueView, buttonView) + "\n"

	}

}

func (m model) paramsListView() string {
	if len(m.req.Params) == 0 {
		return ""
	}

	var s strings.Builder
	s.WriteString("\ndelete || backspace|remove selected\n")

	var keys []string
	for k := range m.req.Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		val := m.req.Params.Get(k)
		row := fmt.Sprintf("%s: %s", k, val)

		if m.focus == focusParamList && m.paramCursor == i {
			s.WriteString(labelStyleSelected.Render("> "+row) + "\n")
		} else {
			s.WriteString(labelStyleUnselected.Render("  "+row) + "\n")
		}
	}

	return s.String()
}

func (m model) footerView() string {
	return "\n\ntab|next shift+tab|previous [ ctrl+c|quit ]\n"
}
