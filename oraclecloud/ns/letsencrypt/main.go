package main

import (
	"gitlab.com/eper.io/engine/oraclecloud/ns"
)

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// This is useful to run a single instance track its lifetime and shut it down.
func main() {
	ns.SetupNameServer()
	ns.SetupComputeCluster()
}
