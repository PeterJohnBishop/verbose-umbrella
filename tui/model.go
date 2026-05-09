package tui

import tea "charm.land/bubbletea/v2"

// the model defines the state of the app
// which is updated by the Update method
// and displayed by the View method

type model struct {
	reqType         []string
	reqCusror       int
	reqTypeSelected int
}

func InitialModel() model {
	return model{
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
	return nil
}
