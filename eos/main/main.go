package main

import (
	"gitlab.com/eper.io/engine/eos"
	"gitlab.com/eper.io/engine/metadata"
	"gitlab.com/eper.io/engine/ns"
	"net/http"
	"strings"
)

// TODO rename to autoscale

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

func main() {
	ns.SetupNameServer()
	ns.SetupComputeCluster()
	eos.Setup()
	if strings.HasPrefix(metadata.SiteUrl, "https://") {
		go func() {
			err := http.ListenAndServe(":80", http.RedirectHandler(metadata.SiteUrl, http.StatusTemporaryRedirect))
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		}()
		err := http.ListenAndServeTLS(":443", metadata.Certificate, metadata.PrivateKey, nil)
		if err != nil {
			panic("ListenAndServeTLS: " + err.Error())
		}
	}
}
