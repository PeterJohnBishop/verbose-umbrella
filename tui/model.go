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
	inputs          []textinput.Model
	req             Request
	reqType         []string
	reqCusror       int
	reqTypeSelected int
	headerCursor    int
	paramCursor     int
}

func InitialModel() model {

	inputs := make([]textinput.Model, 7)

	inputs[inputNameIdx] = textinput.New()
	inputs[inputNameIdx].Placeholder = "Name/Description"
	inputs[inputNameIdx].Focus() // Default focus on startup
	inputs[inputNameIdx].CharLimit = 156
	inputs[inputNameIdx].SetWidth(120)

	inputs[inputEndpointIdx] = textinput.New()
	inputs[inputEndpointIdx].Placeholder = "http://"
	inputs[inputEndpointIdx].SetWidth(120)

	inputs[inputHeadersKeyIdx] = textinput.New()
	inputs[inputHeadersKeyIdx].Placeholder = "key"
	inputs[inputHeadersKeyIdx].SetWidth(120)

	inputs[inputHeadersValueIdx] = textinput.New()
	inputs[inputHeadersValueIdx].Placeholder = "value"
	inputs[inputHeadersValueIdx].SetWidth(3120)

	inputs[inputParamsKeyIdx] = textinput.New()
	inputs[inputParamsKeyIdx].Placeholder = "key"
	inputs[inputParamsKeyIdx].SetWidth(120)

	inputs[inputParamsValueIdx] = textinput.New()
	inputs[inputParamsValueIdx].Placeholder = "value"
	inputs[inputParamsValueIdx].SetWidth(120)

	inputs[inputBodyIdx] = textinput.New()
	inputs[inputBodyIdx].Placeholder = "{ \"key\": \"value\" }"
	inputs[inputBodyIdx].SetWidth(120)

	return model{
		focus:  focusName,
		inputs: inputs,
		req: Request{
			Headers: make(http.Header),
			Params:  url.Values{},
		},
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
	var cmds []tea.Cmd
	cmds = append(cmds, textinput.Blink)
	cmds = append(cmds, m.updateFocus())
	return tea.Batch(cmds...)
}

type FocusState int

const (
	focusName FocusState = iota
	focusMethod
	focusEndpoint
	focusHeaderKey
	focusHeaderValue
	focusHeaderSubmit
	focusHeaderList
	focusParamKey
	focusParamValue
	focusParamSubit
	focusParamList
	focusBody
	focuseResponse
)

const (
	inputNameIdx         = 0
	inputEndpointIdx     = 1
	inputHeadersKeyIdx   = 2
	inputHeadersValueIdx = 3
	inputParamsKeyIdx    = 4
	inputParamsValueIdx  = 5
	inputBodyIdx         = 6
)

const maxFocus = focuseResponse // or whichever is last

type Request struct {
	Name     string
	Method   string
	Endpoint string
	Params   url.Values
	Headers  http.Header
	Body     any
	Response []byte
}
