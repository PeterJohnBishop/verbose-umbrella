package tui

import tea "charm.land/bubbletea/v2"

// the Update method returns an updated model state
// and optionally sends a command 'Cmd'.

// Commands perform some I/O and return a message 'Msg'.
// The message defines the update to be made by the Update method.

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.focus++
			if m.focus > maxFocus {
				m.focus = 0
			}
			return m, m.updateFocus()

		case "shift+tab":
			m.focus--
			if m.focus < 0 {
				m.focus = maxFocus
			}
			return m, m.updateFocus()

		case "ctrl+1":
			m.focus = focusNameInput
			return m, m.updateFocus()
		case "ctrl+2":
			m.focus = focusMethodSelector
			return m, m.updateFocus()
		case "ctrl+3":
			m.focus = focusHeaders
			return m, m.updateFocus()
		}

		switch m.focus {

		case focusNameInput:
			// Pressing Enter or Esc shifts focus away from the text box
			if msg.String() == "esc" || msg.String() == "enter" {
				m.focus = focusMethodSelector
				cmd = m.updateFocus()
			}

		case focusMethodSelector:
			switch msg.String() {
			case "left", "h":
				if m.reqTypeSelected > 0 {
					m.reqTypeSelected--
					m.req.Method = m.reqType[m.reqTypeSelected]
				}
			case "right", "l":
				if m.reqTypeSelected < len(m.reqType)-1 {
					m.reqTypeSelected++
					m.req.Method = m.reqType[m.reqTypeSelected]
				}
			}

		case focusHeaders:
			// Placeholder for later
		}
	}

	var tiCmd tea.Cmd
	m.reqNameInput, tiCmd = m.reqNameInput.Update(msg)

	return m, tea.Batch(cmd, tiCmd)
}

func (m *model) updateFocus() tea.Cmd {
	if m.focus == focusNameInput {
		return m.reqNameInput.Focus()
	}

	m.reqNameInput.Blur()
	return nil
}
