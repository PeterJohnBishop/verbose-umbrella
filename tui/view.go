package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

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
	str := lipgloss.JoinVertical(
		lipgloss.Left,
		m.headerView(),
		m.requestNameInput(),
		m.methodSelector(),
		m.requestEndpointInput(),
		m.headersInputView(),
		m.headersListView(),
		m.paramsInputView(),
		m.paramsListView(),
		m.bodyView(),
		m.sendRequestView(),
		m.footerView(),   // Keybindings
		m.responseView(), // Response Data below keybindings
	)
	return tea.NewView(str)
}

func (m model) headerView() string {
	return "API Testing\n"
}

func (m model) requestNameInput() string {
	var labelStyle lipgloss.Style
	if m.focus == focusName {
		labelStyle = labelStyleSelected
	} else {
		labelStyle = labelStyleUnselected
	}
	return labelStyle.Render("Name") + "\n" + m.inputs[inputNameIdx].View() + "\n"
}

func (m model) requestEndpointInput() string {
	var labelStyle lipgloss.Style
	if m.focus == focusEndpoint {
		labelStyle = labelStyleSelected
	} else {
		labelStyle = labelStyleUnselected
	}
	return labelStyle.Render("Endpoint") + "\n" + m.inputs[inputEndpointIdx].View() + "\n"
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

	title := keyLabelStyle.Render("Headers") + "\n"

	keyView := m.inputs[inputHeadersKeyIdx].View()
	valueView := m.inputs[inputHeadersValueIdx].View()

	keyContainer := lipgloss.NewStyle().Width(30).Render(keyView)

	key := m.inputs[inputHeadersKeyIdx].Value()
	val := m.inputs[inputHeadersValueIdx].Value()

	inputRow := lipgloss.JoinHorizontal(lipgloss.Top, keyContainer, valueView)

	if key == "" || val == "" {
		return title + inputRow + "\n"
	}

	return title + lipgloss.JoinHorizontal(lipgloss.Top, inputRow, lipgloss.NewStyle().MarginLeft(2).Render(buttonView)) + "\n"
}

func (m model) headersListView() string {
	if len(m.req.Headers) == 0 {
		return ""
	}
	var s strings.Builder
	var keys []string
	for k := range m.req.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		val := m.req.Headers.Get(k)

		row := fmt.Sprintf("%-30s %s", k, val)

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

	title := keyLabelStyle.Render("Params") + "\n"

	keyView := m.inputs[inputParamsKeyIdx].View()
	valueView := m.inputs[inputParamsValueIdx].View()

	keyContainer := lipgloss.NewStyle().Width(30).Render(keyView)

	key := m.inputs[inputParamsKeyIdx].Value()
	val := m.inputs[inputParamsValueIdx].Value()

	inputRow := lipgloss.JoinHorizontal(lipgloss.Top, keyContainer, valueView)

	if key == "" || val == "" {
		return title + inputRow + "\n"
	}
	return title + lipgloss.JoinHorizontal(lipgloss.Top, inputRow, lipgloss.NewStyle().MarginLeft(2).Render(buttonView)) + "\n"
}

func (m model) paramsListView() string {
	if len(m.req.Params) == 0 {
		return ""
	}
	var s strings.Builder
	var keys []string
	for k := range m.req.Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		val := m.req.Params.Get(k)

		row := fmt.Sprintf("%-30s %s", k, val)

		if m.focus == focusParamList && m.paramCursor == i {
			s.WriteString(labelStyleSelected.Render("> "+row) + "\n")
		} else {
			s.WriteString(labelStyleUnselected.Render("  "+row) + "\n")
		}
	}
	return s.String()
}

func (m model) bodyView() string {
	switch m.bodyType {
	case BodyRaw:
		var labelStyle lipgloss.Style
		if m.focus == focusBodyJSON {
			labelStyle = labelStyleSelected
		} else {
			labelStyle = labelStyleUnselected
		}
		return labelStyle.Render("Body (JSON)\n") + "\n" + m.bodyTextArea.View() + "\n"

	case BodyForm:
		return m.bodyFormInputView() + m.bodyFormListView()

	case BodyFile:
		var labelStyle lipgloss.Style
		if m.focus == focusBodyFile {
			labelStyle = labelStyleSelected
		} else {
			labelStyle = labelStyleUnselected
		}

		statusText := ""
		if m.inputs[inputFileBodyIdx].Value() != "" {
			if m.fileExists {
				statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("✓ " + m.fileError)
			} else {
				statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("✗ " + m.fileError)
			}
		}

		return labelStyle.Render("Body (File Path)") + "\n" +
			m.inputs[inputFileBodyIdx].View() + "\n" +
			statusText + "\n"
	}
	return ""
}

func (m model) bodyFormInputView() string {
	var keyLabelStyle lipgloss.Style
	if m.focus == focusBodyKey || m.focus == focusBodyValue {
		keyLabelStyle = labelStyleSelected
	} else {
		keyLabelStyle = labelStyleUnselected
	}

	var buttonView string
	if m.focus == focusBodySubmit {
		buttonView = buttonStyleSelected.Render("[ add ]")
	} else {
		buttonView = buttonStyleUnselected.Render("[ add ]")
	}

	title := keyLabelStyle.Render("Body (Form-Data)") + "\n"

	keyView := m.inputs[inputFormBodyKeyIdx].View()
	valueView := m.inputs[inputFormBodyValueIdx].View()

	keyContainer := lipgloss.NewStyle().Width(30).Render(keyView)

	key := m.inputs[inputFormBodyKeyIdx].Value()
	val := m.inputs[inputFormBodyValueIdx].Value()

	inputRow := lipgloss.JoinHorizontal(lipgloss.Top, keyContainer, valueView)

	if key == "" || val == "" {
		return title + inputRow + "\n"
	}
	return title + lipgloss.JoinHorizontal(lipgloss.Top, inputRow, lipgloss.NewStyle().MarginLeft(2).Render(buttonView)) + "\n"
}

func (m model) bodyFormListView() string {
	if len(m.req.FormData) == 0 {
		return ""
	}
	var s strings.Builder
	var keys []string
	for k := range m.req.FormData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		val := m.req.FormData.Get(k)

		row := fmt.Sprintf("%-30s %s", k, val)

		if m.focus == focusBodyList && m.bodyCursor == i {
			s.WriteString(labelStyleSelected.Render("> "+row) + "\n")
		} else {
			s.WriteString(labelStyleUnselected.Render("  "+row) + "\n")
		}
	}
	return s.String()
}

func (m model) sendRequestView() string {
	if m.focus == focusSendReq {
		return "\n" + buttonStyleSelected.Render("[ SEND REQUEST ]") + "\n"
	}
	return "\n" + buttonStyleUnselected.Render("[ SEND REQUEST ]") + "\n"
}

func (m model) responseView() string {
	var labelStyle lipgloss.Style

	// Using your spelling of 'focuseResponse' from the previous iota
	if m.focus == focuseResponse {
		labelStyle = labelStyleSelected
	} else {
		labelStyle = labelStyleUnselected
	}

	header := labelStyle.Render("Response Data") + "\n\n"

	return header + m.viewport.View()
}

func (m model) footerView() string {
	return "\n> [next] < [previous] ctrl+b [toggle: JSON, FormData, File] ctrl+c [quit]\n"
}

func formatJSONWithLineNumbers(rawJSON string) string {
	if rawJSON == "" {
		return unselectedStyle.Render("Awaiting response...")
	}

	var pretty bytes.Buffer
	err := json.Indent(&pretty, []byte(rawJSON), "", "  ")

	textToFormat := pretty.String()
	if err != nil {
		textToFormat = rawJSON
	}

	lines := strings.Split(textToFormat, "\n")
	var b strings.Builder

	numWidth := len(fmt.Sprintf("%d", len(lines)))
	lineNumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginRight(1)

	for i, line := range lines {
		lineNum := fmt.Sprintf("%*d", numWidth, i+1)
		b.WriteString(lineNumStyle.Render(lineNum) + "│ " + line + "\n")
	}

	return b.String()
}
