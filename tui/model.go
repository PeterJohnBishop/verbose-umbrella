package tui

import (
	"net/http"
	"net/url"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	focus           FocusState
	inputs          []textinput.Model
	req             Request
	reqType         []string
	reqCusror       int
	reqTypeSelected int
	headerCursor    int
	paramCursor     int
	bodyCursor      int
	bodyType        BodyType
	bodyTextArea    textarea.Model
	fp              filepicker.Model
	selectedFile    string
}

func InitialModel() model {
	inputs := make([]textinput.Model, 8)

	inputs[inputNameIdx] = textinput.New()
	inputs[inputNameIdx].Placeholder = "Name/Description"
	inputs[inputNameIdx].Focus()
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
	inputs[inputHeadersValueIdx].SetWidth(120)

	inputs[inputParamsKeyIdx] = textinput.New()
	inputs[inputParamsKeyIdx].Placeholder = "key"
	inputs[inputParamsKeyIdx].SetWidth(120)

	inputs[inputParamsValueIdx] = textinput.New()
	inputs[inputParamsValueIdx].Placeholder = "value"
	inputs[inputParamsValueIdx].SetWidth(120)

	inputs[inputFormBodyKeyIdx] = textinput.New()
	inputs[inputFormBodyKeyIdx].Placeholder = "key"
	inputs[inputFormBodyKeyIdx].SetWidth(120)

	inputs[inputFormBodyValueIdx] = textinput.New()
	inputs[inputFormBodyValueIdx].Placeholder = "value"
	inputs[inputFormBodyValueIdx].SetWidth(120)

	bodyTextArea := textarea.New()
	bodyTextArea.Placeholder = "{ \"key\": \"value\" }"
	bodyTextArea.SetWidth(60)
	bodyTextArea.SetHeight(10)

	fp := filepicker.New()
	fp.AllowedTypes = []string{".json", ".txt", ".csv", ".bin"}
	fp.CurrentDirectory = "."

	return model{
		focus:  focusName,
		inputs: inputs,
		req: Request{
			Headers:  make(http.Header),
			Params:   url.Values{},
			FormData: url.Values{},
		},
		reqType: []string{
			"GET",
			"PUT",
			"PATCH",
			"POST",
			"DELETE",
		},
		reqTypeSelected: 0,
		bodyType:        BodyRaw,
		bodyTextArea:    bodyTextArea,
		fp:              fp,
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, textinput.Blink)
	cmds = append(cmds, m.updateFocus())
	return tea.Batch(cmds...)
}

type BodyType int

const (
	BodyRaw BodyType = iota
	BodyForm
	BodyFile
)

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
	focusBodyKey
	focusBodyValue
	focusBodySubmit
	focusBodyList
	focusBodyJSON
	focusBodyFile
	focusSendReq
	focuseResponse
)

const (
	inputNameIdx          = 0
	inputEndpointIdx      = 1
	inputHeadersKeyIdx    = 2
	inputHeadersValueIdx  = 3
	inputParamsKeyIdx     = 4
	inputParamsValueIdx   = 5
	inputFormBodyKeyIdx   = 6
	inputFormBodyValueIdx = 7
)

const maxFocus = focuseResponse

type Request struct {
	Name     string
	Method   string
	Endpoint string
	Params   url.Values
	Headers  http.Header
	FormData url.Values
	BodyType string
	Body     any
	Response []byte
}
