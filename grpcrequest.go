package main

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"net/http"
)

const (
	// header to detect if it is a grpc+json request
	contentTypeGRPCJSON = "application/grpc+json"

	grpcNoCompression byte = 0x00
)

func modifyRequestToJSONgRPC(r *http.Request) *http.Request {
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md

	var body []byte
	// read body so we can add the grpc prefix
	if r.Body != nil {
		body, _ = ioutil.ReadAll(r.Body)
	}

	b := make([]byte, 0, len(body)+5)
	buff := bytes.NewBuffer(b)

	// grpc prefix is
	// 1 byte: compression indicator
	// 4 bytes: content length (excluding prefix)
	_ = buff.WriteByte(grpcNoCompression) // 0 or 1, indicates compressed payload

	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(body)))

	_, _ = buff.Write(lenBytes)
	_, _ = buff.Write(body)

	// create new request
	req, _ := http.NewRequest(r.Method, r.URL.String(), buff)
	req.Header = r.Header

	// remove content length header
	req.Header.Del(headerContentLength)

	return req

}

func isJSONGRPC(r *http.Request) bool {

	h := r.Header.Get("Content-Type")

	if h == contentTypeGRPCJSON {
		return true
	}

	return false
}
