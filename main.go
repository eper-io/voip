package main

import (
	"crypto/tls"
	"fmt"
	"gitlab.com/eper.io/engine/line"
	"gitlab.com/eper.io/engine/metadata"
	"net/http"
	"net/url"
)

// main.go - Relay websocket connections for streaming applications. See README.md for usage

// Licensed under Creative Commons CC0.
// To the extent possible under law, the author(s) have dedicated all copyright and related and
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along with this software.
// If not, see <https://creativecommons.org/publicdomain/zero/1.0/legalcode>.

// This example demonstrates a trivial echo server.
func main() {
	line.Setup()
	if metadata.PrivateKey != "" {
		err := listenAndServeTLS(":443", metadata.Certificate, metadata.PrivateKey, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
		err = http.ListenAndServe(":80", nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func listenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	url1, _ := url.Parse(metadata.SiteUrl)
	fmt.Println(url1.Hostname())
	config := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
	server := &http.Server{Addr: addr, Handler: handler, TLSConfig: config}
	return server.ListenAndServeTLS(certFile, keyFile)
}
