package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (r *Request) BuildHTTPRequest() (*http.Request, error) {
	var reqBody io.Reader

	// Content-Type header determines the body type
	contentType := r.Headers.Get("Content-Type")

	if r.Body != nil {
		switch {
		// body is JSON
		case strings.Contains(contentType, "application/json"):
			jsonBytes, err := json.Marshal(r.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal json body: %w", err)
			}
			reqBody = bytes.NewReader(jsonBytes)

			// body is Form Data
		case strings.Contains(contentType, "application/x-www-form-urlencoded"):
			if formData, ok := r.Body.(url.Values); ok {
				reqBody = strings.NewReader(formData.Encode())
			} else {
				return nil, fmt.Errorf("form body must be of type url.Values")
			}

		default:
			// Raw String/Bytes fallback
			if strBody, ok := r.Body.(string); ok {
				reqBody = strings.NewReader(strBody)
			} else if byteBody, ok := r.Body.([]byte); ok {
				reqBody = bytes.NewReader(byteBody)
			}
		}
	}

	req, err := http.NewRequest(r.Method, r.Endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header = r.Headers

	if len(r.Params) > 0 {
		req.URL.RawQuery = r.Params.Encode()
	}

	return req, nil
}
