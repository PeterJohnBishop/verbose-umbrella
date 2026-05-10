package tui

import (
	"net/http"
	"net/url"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// the model defines the state of the app
// which is updated by the Update method
// and displayed by the View method

type model struct {
	focus           FocusState
	req             Request
	reqNameInput    textinput.Model
	reqType         []string
	reqCusror       int
	reqTypeSelected int
}

func InitialModel() model {

	rnTi := textinput.New()
	rnTi.Placeholder = "Name/Description"
	rnTi.SetVirtualCursor(false)
	rnTi.Focus()
	rnTi.CharLimit = 156
	rnTi.SetWidth(30)

	return model{
		focus:        focusNameInput,
		reqNameInput: rnTi,
		reqType: []string{
			"GET",
			"PUT",
			"PATCH",
			"POST",
			"DELETE",
		},
		reqTypeSelected: 0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

type FocusState int

const (
	focusNameInput FocusState = iota
	focusMethodSelector
	focusHeaders
	// Add focusBody, focusParams, etc.
)

const maxFocus = focusMethodSelector // or whichever is last

type Request struct {
	Name     string
	Method   string
	Endpoint string
	Params   url.Values
	Headers  http.Header
	Body     any
	Response []byte
}
