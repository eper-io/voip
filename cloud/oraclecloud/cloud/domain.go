package cloud

import (
	"fmt"
	"github.com/eper.io/cloud/oraclecloud/metadata"
	"github.com/eper.io/cloud/oraclecloud/ns"
	"math/rand"
	"net"
	"strings"
	"time"
)

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

func SetupOracleComputeCluster() {
	ips, err := net.LookupHost(metadata.DomainNS)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	ip := [4]byte{127, 0, 0, 1}
	if len(ips) > 0 {
		ipv4d := [4]byte{127, 0, 0, 1}
		n, _ := fmt.Sscanf(ips[0], "%d.%d.%d.%d", &ipv4d[0], &ipv4d[1], &ipv4d[2], &ipv4d[3])
		if n == 4 {
			ip = ipv4d
		}
	}
	EntryPoint = ip
	fmt.Println("Host", metadata.DomainNS, ip)

	split := strings.Split(metadata.HostNames, "\n")
	command := fmt.Sprintf("certbot --nginx -m hq@schmied.us --cert-name %s -d %s,", metadata.Domain, metadata.Domain)
	list := command
	shuffled := split
	rand.Seed(time.Now().UnixNano())
	length := len(shuffled)
	rand.Shuffle(length, func(i, j int) {
		t := shuffled[i%length]
		shuffled[i%length] = shuffled[j%length]
		shuffled[j%length] = t
	})

	ns.Nodes[metadata.Domain] = EntryPoint
	for _, v := range shuffled {
		host := strings.TrimSpace(v) + "." + metadata.Domain
		ns.Nodes[host] = EntryPoint
		ns.Candidates = append(ns.Candidates, host)
		if list != command {
			list = list + ","
		}
		list = list + host
	}

	// Just print the cert command. See documentation/almalinux.sh
	list = list + fmt.Sprintf(" --https-port 4443 certonly")
	fmt.Println(list)
}
