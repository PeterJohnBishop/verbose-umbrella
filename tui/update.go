package tui

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case responseMsg:
		if msg.err != nil {
			errJSON := fmt.Sprintf(`{"error": "%s"}`, msg.err.Error())
			m.req.Response = []byte(errJSON)
			m.viewport.SetContent(formatJSONWithLineNumbers(errJSON))
		} else {
			m.req.Response = msg.body
			m.viewport.SetContent(formatJSONWithLineNumbers(string(msg.body)))
		}
		return m, nil

	case tea.WindowSizeMsg:
		formHeight := 55
		if !m.ready {
			m.viewport = viewport.New(
				viewport.WithWidth(msg.Width),
				viewport.WithHeight(msg.Height-formHeight),
			)
			m.viewport.SetContent(formatJSONWithLineNumbers(""))
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height - formHeight)
		}
	case tea.KeyMsg:
		key := msg.String()

		// GLOBAL NAVIGATION (Overrides everything)
		switch key {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+b":
			m.bodyType = (m.bodyType + 1) % 3
			m.moveToBodyFocus()
			return m, m.updateFocus()

		case ">":
			m.handleForwardNav()
			return m, m.updateFocus()

		case "<":
			m.handleBackwardNav()
			return m, m.updateFocus()
		}

		// COMPONENT LOGIC (Only runs if not navigating)
		switch m.focus {
		case focusName:
			if key == "enter" {
				m.focus = focusMethod
				return m, m.updateFocus()
			}

		case focusMethod:
			switch key {
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
			case "enter":
				m.focus = focusEndpoint
				return m, m.updateFocus()
			}

		case focusEndpoint:
			if key == "enter" {
				m.focus = focusHeaderKey
				return m, m.updateFocus()
			}

		case focusHeaderKey:
			if key == "enter" {
				m.focus = focusHeaderValue
				return m, m.updateFocus()
			}

		case focusHeaderValue:
			if key == "enter" {
				m.focus = focusHeaderSubmit
				return m, m.updateFocus()
			}

		case focusHeaderSubmit:
			if key == "enter" {
				k := m.inputs[inputHeadersKeyIdx].Value()
				v := m.inputs[inputHeadersValueIdx].Value()
				if k != "" && v != "" {
					m.req.Headers.Add(k, v)
				}
				m.inputs[inputHeadersKeyIdx].SetValue("")
				m.inputs[inputHeadersValueIdx].SetValue("")
				m.focus = focusHeaderKey
				return m, m.updateFocus()
			}

		case focusHeaderList:
			m.handleHeaderListKeys(key)

		case focusParamKey:
			if key == "enter" {
				m.focus = focusParamValue
				return m, m.updateFocus()
			}

		case focusParamValue:
			if key == "enter" {
				m.focus = focusParamSubit
				return m, m.updateFocus()
			}

		case focusParamSubit:
			if key == "enter" {
				k := m.inputs[inputParamsKeyIdx].Value()
				v := m.inputs[inputParamsValueIdx].Value()
				if k != "" && v != "" {
					m.req.Params.Add(k, v)
					m.syncEndpoint()
				}
				m.inputs[inputParamsKeyIdx].SetValue("")
				m.inputs[inputParamsValueIdx].SetValue("")
				m.focus = focusParamKey
				return m, m.updateFocus()
			}

		case focusParamList:
			m.handleParamListKeys(key)

		case focusBodyJSON:
			if key == "tab" {
				m.bodyTextArea.InsertString("    ")
				return m, nil
			}
			// Manual update to ensure component gets the message
			m.bodyTextArea, cmd = m.bodyTextArea.Update(msg)
			return m, cmd

		case focusBodyKey:
			if key == "enter" {
				m.focus = focusBodyValue
				return m, m.updateFocus()
			}

		case focusBodyValue:
			if key == "enter" {
				m.focus = focusBodySubmit
				return m, m.updateFocus()
			}

		case focusBodySubmit:
			if key == "enter" {
				k := m.inputs[inputFormBodyKeyIdx].Value()
				v := m.inputs[inputFormBodyValueIdx].Value()
				if k != "" && v != "" {
					m.req.FormData.Add(k, v)
				}
				m.inputs[inputFormBodyKeyIdx].SetValue("")
				m.inputs[inputFormBodyValueIdx].SetValue("")
				m.focus = focusBodyKey
				return m, m.updateFocus()
			}

		case focusBodyList:
			m.handleBodyListKeys(key)

		case focusBodyFile:
			if key == "enter" {
				m.focus = focusSendReq
				return m, m.updateFocus()
			}
		case focusSendReq:
			if msg.String() == "enter" {
				endpoint := m.inputs[inputEndpointIdx].Value()
				if endpoint != "" && !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
					endpoint = "http://" + endpoint
				}
				m.req.Endpoint = endpoint
				if m.bodyType == BodyRaw {
					m.req.Body = m.bodyTextArea.Value()
				} else if m.bodyType == BodyFile {
					if !m.fileExists {
						return m, nil
					}
					m.req.Body = m.inputs[inputFileBodyIdx].Value()
				}

				return m, sendRequestCmd(&m.req)
			}
		}
	}

	var tiCmd tea.Cmd
	for i := range m.inputs {
		m.inputs[i], tiCmd = m.inputs[i].Update(msg)
		cmds = append(cmds, tiCmd)

		if i == inputFileBodyIdx && m.focus == focusBodyFile {
			m.fileExists, m.fileError = checkFileExists(m.inputs[i].Value())
		}
	}

	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}

// NAVIGATION HELPERS

func (m *model) handleForwardNav() {
	switch m.focus {
	case focusHeaderValue:
		if m.inputs[inputHeadersKeyIdx].Value() == "" || m.inputs[inputHeadersValueIdx].Value() == "" {
			if len(m.req.Headers) == 0 {
				m.focus = focusParamKey
			} else {
				m.focus = focusHeaderList
			}
		} else {
			m.focus = focusHeaderSubmit
		}
	case focusHeaderList:
		m.focus = focusParamKey
	case focusParamValue:
		if m.inputs[inputParamsKeyIdx].Value() == "" || m.inputs[inputParamsValueIdx].Value() == "" {
			if len(m.req.Params) == 0 {
				m.moveToBodyFocus()
			} else {
				m.focus = focusParamList
			}
		} else {
			m.focus = focusParamSubit
		}
	case focusParamList:
		m.moveToBodyFocus()
	case focusBodyValue:
		if m.inputs[inputFormBodyKeyIdx].Value() == "" || m.inputs[inputFormBodyValueIdx].Value() == "" {
			if len(m.req.FormData) == 0 {
				m.focus = focusSendReq
			} else {
				m.focus = focusBodyList
			}
		} else {
			m.focus = focusBodySubmit
		}
	case focusBodyJSON, focusBodyFile, focusBodyList, focusBodySubmit:
		m.focus = focusSendReq
	case focusSendReq:
		m.focus = focusName
	default:
		m.focus++
	}
}

func (m *model) handleBackwardNav() {
	switch m.focus {
	case focusSendReq:
		switch m.bodyType {
		case BodyRaw:
			m.focus = focusBodyJSON
		case BodyFile:
			m.focus = focusBodyFile
		case BodyForm:
			if len(m.req.FormData) > 0 {
				m.focus = focusBodyList
			} else {
				m.focus = focusBodyValue
			}
		}

	case focusBodyList, focusBodySubmit:
		m.focus = focusBodyValue
	case focusBodyValue:
		m.focus = focusBodyKey

	case focusBodyJSON, focusBodyKey, focusBodyFile:
		if len(m.req.Params) > 0 {
			m.focus = focusParamList
		} else {
			m.focus = focusParamValue
		}

	case focusParamList, focusParamSubit:
		m.focus = focusParamValue
	case focusParamValue:
		m.focus = focusParamKey

	case focusParamKey:
		if len(m.req.Headers) > 0 {
			m.focus = focusHeaderList
		} else {
			m.focus = focusHeaderValue
		}

	// Headers reverse navigation
	case focusHeaderList, focusHeaderSubmit:
		m.focus = focusHeaderValue
	case focusHeaderValue:
		m.focus = focusHeaderKey

	// Top section reverse navigation
	case focusHeaderKey:
		m.focus = focusEndpoint
	case focusEndpoint:
		m.focus = focusMethod
	case focusMethod:
		m.focus = focusName
	case focusName:
		m.focus = focusSendReq

	default:
		m.focus--
	}
}

// LIST HANDLERS

func (m *model) handleHeaderListKeys(key string) {
	max := len(m.req.Headers) - 1
	if key == "up" || key == "k" {
		if m.headerCursor > 0 {
			m.headerCursor--
		}
	}
	if key == "down" || key == "j" {
		if m.headerCursor < max {
			m.headerCursor++
		}
	}
	if key == "delete" || key == "backspace" {
		if len(m.req.Headers) > 0 {
			keys := getSortedKeys(m.req.Headers)
			m.req.Headers.Del(keys[m.headerCursor])
			if m.headerCursor > 0 && m.headerCursor >= len(m.req.Headers) {
				m.headerCursor--
			}
		}
	}
}

func (m *model) handleParamListKeys(key string) {
	max := len(m.req.Params) - 1
	if key == "up" || key == "k" {
		if m.paramCursor > 0 {
			m.paramCursor--
		}
	}
	if key == "down" || key == "j" {
		if m.paramCursor < max {
			m.paramCursor++
		}
	}
	if key == "delete" || key == "backspace" {
		if len(m.req.Params) > 0 {
			keys := getSortedKeysParams(m.req.Params)
			m.req.Params.Del(keys[m.paramCursor])
			m.syncEndpoint()
			if m.paramCursor > 0 && m.paramCursor >= len(m.req.Params) {
				m.paramCursor--
			}
		}
	}
}

func (m *model) handleBodyListKeys(key string) {
	max := len(m.req.FormData) - 1
	if key == "up" || key == "k" {
		if m.bodyCursor > 0 {
			m.bodyCursor--
		}
	}
	if key == "down" || key == "j" {
		if m.bodyCursor < max {
			m.bodyCursor++
		}
	}
	if key == "delete" || key == "backspace" {
		if len(m.req.FormData) > 0 {
			keys := getSortedKeysParams(m.req.FormData)
			m.req.FormData.Del(keys[m.bodyCursor])
			if m.bodyCursor > 0 && m.bodyCursor >= len(m.req.FormData) {
				m.bodyCursor--
			}
		}
	}
}

func (m *model) syncEndpoint() {
	u, err := url.Parse(m.inputs[inputEndpointIdx].Value())
	if err == nil {
		u.RawQuery = m.req.Params.Encode()
		m.inputs[inputEndpointIdx].SetValue(u.String())
	}
}

func getSortedKeys(h http.Header) []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getSortedKeysParams(v url.Values) []string {
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (m *model) moveToBodyFocus() {
	switch m.bodyType {
	case BodyRaw:
		m.focus = focusBodyJSON
	case BodyForm:
		m.focus = focusBodyKey
	case BodyFile:
		m.focus = focusBodyFile
	}
}

func (m *model) updateFocus() tea.Cmd {
	var cmds []tea.Cmd
	activeIndex := -1

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
	case focusBodyKey:
		activeIndex = inputFormBodyKeyIdx
	case focusBodyValue:
		activeIndex = inputFormBodyValueIdx
	case focusBodyFile:
		activeIndex = inputFileBodyIdx
	}

	for i := range m.inputs {
		if i == activeIndex {
			cmds = append(cmds, m.inputs[i].Focus())
		} else {
			m.inputs[i].Blur()
		}
	}

	if m.focus == focusBodyJSON {
		cmds = append(cmds, m.bodyTextArea.Focus())
	} else {
		m.bodyTextArea.Blur()
	}

	return tea.Batch(cmds...)
}

func checkFileExists(path string) (bool, string) {
	if path == "" {
		return false, ""
	}

	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, "File does not exist"
	}
	if err != nil {
		return false, "Cannot access: " + err.Error()
	}
	if info.IsDir() {
		return false, "Path is a directory, not a file"
	}

	sizeBytes := float64(info.Size())

	// 1 MB = 1024 * 1024 bytes
	if sizeBytes >= 1048576 {
		sizeMB := sizeBytes / 1048576
		return true, fmt.Sprintf("File found (%.2f MB)", sizeMB)
	}

	sizeKB := sizeBytes / 1024.0
	return true, fmt.Sprintf("File found (%.2f KB)", sizeKB)
}
