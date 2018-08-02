package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"golang.org/x/net/http2"
)

const (
	headerContentLength  = "Content-Length"
	headerGRPCMessage    = "Grpc-Message"
	headerGRPCStatusCode = "Grpc-Status"
	headerUseInsecure    = "Grpc-Insecure"

	defaultClientTimeout = time.Second * 60
)

// Transport struct for intercepting grpc+json requests
type Transport struct {
	HTTPClient    *http.Client
	H2Client      *http.Client
	H2NoTLSClient *http.Client
}

/*
	NewProxy returns a configured reverse proxy
	to handle grpc+json requests
*/
func NewProxy() *httputil.ReverseProxy {

	h2NoTLSClient := &http.Client{
		// Skip TLS dial
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
		},
		Timeout: defaultClientTimeout,
	}

	h2Client := &http.Client{
		Transport: &http2.Transport{},
		Timeout:   defaultClientTimeout,
	}

	client := &http.Client{
		Timeout: defaultClientTimeout,
	}

	t := &Transport{
		HTTPClient:    client,
		H2Client:      h2Client,
		H2NoTLSClient: h2NoTLSClient,
	}

	u := url.URL{}
	p := httputil.NewSingleHostReverseProxy(&u)
	p.Director = t.director
	p.Transport = t

	return p
}

func (t Transport) director(r *http.Request) {
}

/*
  RoundTrip handles processing the incoming request
  and outgoing response for grpc+json detection
*/
func (t Transport) RoundTrip(r *http.Request) (*http.Response, error) {

	isGRPC := false
	if isJSONGRPC(r) {
		isGRPC = true
		r = modifyRequestToJSONgRPC(r)
	}

	client := t.HTTPClient
	if isGRPC {
		if r.Header.Get(headerUseInsecure) != "" {
			client = t.H2NoTLSClient
		} else {
			client = t.H2Client
		}
	}

	// clear requestURI, set in call to director
	r.RequestURI = ""

	log.Printf("proxying request url=[%s] isJSONGRPC=[%t]\n", r.URL.String(), isGRPC)

	resp, err := client.Do(r)
	if err != nil {
		log.Printf("unable to do request err=[%s]", err)

		buff := bytes.NewBuffer(nil)
		buff.WriteString(err.Error())
		resp = &http.Response{
			StatusCode: 502,
			Body:       ioutil.NopCloser(buff),
		}

		return resp, nil
	}

	if isGRPC {
		return handleGRPCResponse(resp)
	}

	return resp, err
}
