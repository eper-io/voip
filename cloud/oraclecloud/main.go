package main

import (
	"fmt"
	"github.com/eper.io/cloud/oraclecloud/cloud"
	"github.com/eper.io/cloud/oraclecloud/ns"
	"time"
)

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

func main() {
	ns.SetupNameServer()
	cloud.SetupOracleComputeCluster()
	time.Sleep(10 * time.Second)
	id, host, ip := cloud.LaunchInstance()
	fmt.Println(id, host, ip)
	time.Sleep(3 * time.Minute)
	fmt.Println("Terminating.")
	cloud.TerminateInstance(id, host)
	time.Sleep(1 * time.Minute)
}
