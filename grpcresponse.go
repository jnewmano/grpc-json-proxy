package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func handleGRPCResponse(resp *http.Response) (*http.Response, error) {

	code := metadata(resp, headerGRPCStatusCode)

	if code != "0" && code != "" {
		r := struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		}{
			Error: metadata(resp, headerGRPCMessage),
			Code:  code,
		}

		buff := bytes.NewBuffer(nil)
		_ = json.NewEncoder(buff).Encode(r)

		resp.StatusCode = 500
		resp.Body = io.NopCloser(buff)

		return resp, nil
	}

	prefix := make([]byte, 5)
	_, _ = resp.Body.Read(prefix)

	resp.Header.Del(headerContentLength)

	return resp, nil

}

func metadata(resp *http.Response, field string) string {
	v := resp.Header.Get(field)
	if v != "" {
		return v
	}
	return resp.Trailer.Get(field)
}
