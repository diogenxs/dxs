package utils

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"
)

func MyHTTP(url string, method string, verbose bool) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	var start, connect, dns, tlsHandshake time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			if verbose {
				fmt.Printf("DNS Done (%s): %v\n", url, time.Since(dns))
			}
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			if verbose {
				fmt.Printf("TLS Handshake (%s): %v\n", url, time.Since(tlsHandshake))
			}
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			if verbose {
				fmt.Printf("Connect time (%s): %v\n", url, time.Since(connect))
			}
		},

		GotFirstResponseByte: func() {
			if verbose {
				fmt.Printf("Time from start to first byte (%s): %v\n", url, time.Since(start))
			}
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		return nil, err
	}

	if verbose {
		fmt.Printf("Total time (%s): %v\n", url, time.Since(start))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
