package tui

import (
	"net/url"
	"sort"

	tea "charm.land/bubbletea/v2"
)

// the Update method returns an updated model state
// and optionally sends a command 'Cmd'.

// Commands perform some I/O and return a message 'Msg'.
// The message defines the update to be made by the Update method.

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "tab":
			switch m.focus {
			case focusHeaderValue:
				key := m.inputs[inputHeadersKeyIdx].Value()
				val := m.inputs[inputHeadersValueIdx].Value()

				// If EITHER key or value is empty, skip the submit button
				if key == "" || val == "" {
					if len(m.req.Headers) == 0 {
						m.focus = focusParamKey
					} else {
						m.focus = focusHeaderList
					}
				} else {
					// Both have text, advance normally to focusHeaderSubmit
					m.focus++
				}

			case focusParamValue:
				key := m.inputs[inputParamsKeyIdx].Value()
				val := m.inputs[inputParamsValueIdx].Value()

				// If EITHER key or value is empty, skip the submit button
				if key == "" || val == "" {
					if len(m.req.Params) == 0 {
						m.focus = focusBody
					} else {
						m.focus = focusParamList
					}
				} else {
					// Both have text, advance normally to focusParamSubit
					m.focus++
				}

			default:
				m.focus++
			}

			if m.focus > maxFocus {
				m.focus = 0
			}
			cmd = m.updateFocus()
			cmds = append(cmds, cmd)
		case "shift+tab":
			m.focus--
			if m.focus < 0 {
				m.focus = maxFocus
			}
			cmd = m.updateFocus()
			cmds = append(cmds, cmd)
		}

		switch m.focus {

		case focusName:
			if msg.String() == "esc" || msg.String() == "enter" {
				m.focus = focusMethod
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}

		case focusMethod:
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
				if msg.String() == "esc" || msg.String() == "enter" {
					m.focus = focusHeaderKey
					cmd = m.updateFocus()
					cmds = append(cmds, cmd)
				}
			}

		case focusEndpoint:
			if msg.String() == "esc" {
				m.focus = focusMethod
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
			if msg.String() == "enter" {
				m.focus = focusHeaderKey
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}

		case focusHeaderKey:
			if msg.String() == "esc" {
				m.focus = focusMethod
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
			if msg.String() == "enter" {
				m.focus = focusHeaderValue
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
		case focusHeaderValue:
			if msg.String() == "esc" {
				m.focus = focusHeaderKey
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
			if msg.String() == "enter" {
				m.focus = focusHeaderSubmit
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
		case focusHeaderSubmit:
			if msg.String() == "esc" {
				m.focus = focusHeaderValue
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
			if msg.String() == "enter" {

				key := m.inputs[inputHeadersKeyIdx].Value()
				val := m.inputs[inputHeadersValueIdx].Value()

				if key != "" && val != "" {
					m.req.Headers.Add(key, val)
				}
				m.inputs[inputHeadersKeyIdx].SetValue("")
				m.inputs[inputHeadersValueIdx].SetValue("")

				m.focus = focusHeaderKey
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
		case focusHeaderList:
			maxIndex := len(m.req.Headers) - 1

			switch msg.String() {
			case "up", "k":
				if m.headerCursor > 0 {
					m.headerCursor--
				}
			case "down", "j":
				if m.headerCursor < maxIndex {
					m.headerCursor++
				}
			case "delete", "backspace":
				if len(m.req.Headers) > 0 {
					var keys []string
					for k := range m.req.Headers {
						keys = append(keys, k)
					}
					sort.Strings(keys)

					keyToDelete := keys[m.headerCursor]
					m.req.Headers.Del(keyToDelete)

					if m.headerCursor > len(m.req.Headers)-1 && m.headerCursor > 0 {
						m.headerCursor--
					}
				}
			}
		case focusParamKey:
			if msg.String() == "esc" {
				m.focus = focusHeaderList
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
			if msg.String() == "enter" {
				m.focus = focusParamValue
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
		case focusParamValue:
			if msg.String() == "esc" {
				m.focus = focusParamKey
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
			if msg.String() == "enter" {
				m.focus = focusParamSubit
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
		case focusParamSubit:
			if msg.String() == "esc" {
				m.focus = focusParamValue
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
			if msg.String() == "enter" {

				key := m.inputs[inputParamsKeyIdx].Value()
				val := m.inputs[inputParamsValueIdx].Value()

				if key != "" && val != "" {
					m.req.Params.Add(key, val)

					endpointStr := m.inputs[inputEndpointIdx].Value()
					u, err := url.Parse(endpointStr)
					if err == nil {
						q := u.Query()
						q.Add(key, val)
						u.RawQuery = q.Encode()
						m.inputs[inputEndpointIdx].SetValue(u.String())
					}
				}

				m.inputs[inputParamsKeyIdx].SetValue("")
				m.inputs[inputParamsValueIdx].SetValue("")

				m.focus = focusParamKey
				cmd = m.updateFocus()
				cmds = append(cmds, cmd)
			}
		case focusParamList:
			maxIndex := len(m.req.Params) - 1

			switch msg.String() {
			case "up", "k":
				if m.paramCursor > 0 {
					m.paramCursor--
				}
			case "down", "j":
				if m.paramCursor < maxIndex {
					m.paramCursor++
				}
			case "delete", "backspace":
				if len(m.req.Params) > 0 {
					var keys []string
					for k := range m.req.Params {
						keys = append(keys, k)
					}
					sort.Strings(keys)

					keyToDelete := keys[m.paramCursor]
					m.req.Params.Del(keyToDelete)

					endpointStr := m.inputs[inputEndpointIdx].Value()
					u, err := url.Parse(endpointStr)
					if err == nil {
						q := u.Query()
						q.Del(keyToDelete)
						u.RawQuery = q.Encode()

						m.inputs[inputEndpointIdx].SetValue(u.String())
					}

					if m.paramCursor > len(m.req.Params)-1 && m.paramCursor > 0 {
						m.paramCursor--
					}
				}
			}
		}
	}

	var tiCmd tea.Cmd
	for i := range m.inputs {
		m.inputs[i], tiCmd = m.inputs[i].Update(msg)
		cmds = append(cmds, tiCmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *model) updateFocus() tea.Cmd {
	var cmds []tea.Cmd

	var activeIndex int
	switch m.focus {
	case focusName:
		activeIndex = inputNameIdx
	case focusEndpoint:
		activeIndex = inputEndpointIdx
	case focusHeaderKey:
		activeIndex = inputHeadersKeyIdx
	case focusHeaderValue:
		activeIndex = inputHeadersValueIdx
	case focusParamKey:
		activeIndex = inputParamsKeyIdx
	case focusParamValue:
		activeIndex = inputParamsValueIdx
	case focusBody:
		activeIndex = inputBodyIdx
	case focusMethod:
		activeIndex = -1
	case focuseResponse:
		activeIndex = -1
	}

	for i := range m.inputs {
		if i == activeIndex {
			cmds = append(cmds, m.inputs[i].Focus())
		} else {
			m.inputs[i].Blur()
		}
	}

	return tea.Batch(cmds...)
}
