package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func handleGRPCResponse(resp *http.Response) (*http.Response, error) {

	code := resp.Header.Get(headerGRPCStatusCode)
	if code != "0" && code != "" {
		buff := bytes.NewBuffer(nil)
		grpcMessage := resp.Header.Get(headerGRPCMessage)
		j, _ := json.Marshal(grpcMessage)
		buff.WriteString(`{"error":` + string(j) + ` ,"code":` + code + `}`)

		resp.Body = ioutil.NopCloser(buff)
		resp.StatusCode = 500

		return resp, nil
	}

	prefix := make([]byte, 5)
	_, _ = resp.Body.Read(prefix)

	resp.Header.Del(headerContentLength)

	return resp, nil

}
