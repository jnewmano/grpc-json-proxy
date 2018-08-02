package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	globalClient    *http.Client
	globalTLSClient *http.Client

	shutdownTimeout    = time.Second * 3
	defaultIdleTimeout = time.Second * 60
)

func main() {
	addr := flag.String("http-addr", ":7001", "listen addr for HTTP server")

	flag.Parse()

	ctx := context.Background()

	proxyListenAddr := *addr

	p := NewProxy()

	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx, s, err := StartProxy(ctx, proxyListenAddr, p)
	if err != nil {
		log.Println("unable to start proxy", err)
		os.Exit(1)
	}

	wait(ctx, s)
}

// StartProxy starts a httputil.ReverseProxy listening on addr
func StartProxy(ctx context.Context, addr string, p *httputil.ReverseProxy) (context.Context, stopFunction, error) {

	idleTimeout := defaultIdleTimeout

	s := &http.Server{
		Addr:        addr,
		Handler:     p,
		IdleTimeout: idleTimeout,
	}

	go func() {
		var done func()
		ctx, done = context.WithCancel(ctx)
		log.Println("starting HTTP server", addr)
		err := s.ListenAndServe()
		if err != nil {
			log.Println("listen and serve ended with error", err)
		}
		done()
	}()

	sf := func(ctx context.Context) error {
		log.Println("shutting down server")
		return s.Close()
	}

	return ctx, sf, nil
}

type stopFunction func(ctx context.Context) error

func wait(ctx context.Context, stoppers ...stopFunction) {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-c:
	}

	// call stop functions in reverse order
	for i := len(stoppers) - 1; i >= 0; i-- {
		s := stoppers[i]
		if s == nil {
			continue
		}

		ctx, done := context.WithTimeout(ctx, shutdownTimeout)

		err := s(ctx)
		if err != nil {
			log.Println("did not exit", err)
		}

		done()
	}
}
