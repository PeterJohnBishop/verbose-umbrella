package tui

import tea "charm.land/bubbletea/v2"

// the Update method returns an updated model state
// and optionally sends a command 'Cmd'.

// Commands perform some I/O and return a message 'Msg'.
// The message defines the update to be made by the Update method.

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "h":
			if m.reqCusror > 0 {
				m.reqCusror--
			}
		case "right", "l":
			if m.reqCusror < len(m.reqType)-1 {
				m.reqCusror++
			}
		case "enter", "space":
			m.reqTypeSelected = m.reqCusror
		}
	}

	return m, nil
}
