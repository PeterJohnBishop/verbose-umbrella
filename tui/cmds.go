package tui

import (
	"fmt"
	"io"
	"net/http"
	"time"

	tea "charm.land/bubbletea/v2"
)

func sendRequestCmd(r *Request) tea.Cmd {
	return func() tea.Msg {
		req, err := r.BuildHTTPRequest()
		if err != nil {
			return responseMsg{err: fmt.Errorf("failed to build request: %w", err)}
		}

		client := &http.Client{
			Timeout: 60 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			return responseMsg{err: fmt.Errorf("request failed: %w", err)}
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return responseMsg{err: fmt.Errorf("failed to read response body: %w", err)}
		}

		return responseMsg{body: bodyBytes}
	}
}
