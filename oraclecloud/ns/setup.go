package ns

import (
	"fmt"
	"gitlab.com/eper.io/engine/oraclecloud/metadata"
	"math/rand"
	"net"
	"strings"
	"time"
)

func SetupComputeCluster() {
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

	split := strings.Split(HostNames, "\n")
	command := fmt.Sprintf("certbot --standalone -m hq@schmied.us --cert-name %s -d %s,", metadata.Domain, metadata.Domain)
	list := command
	shuffled := split
	rand.Seed(time.Now().UnixNano())
	length := len(shuffled)
	rand.Shuffle(length, func(i, j int) {
		t := shuffled[i%length]
		shuffled[i%length] = shuffled[j%length]
		shuffled[j%length] = t
	})

	Nodes[metadata.Domain] = ip
	Nodes["example."+metadata.Domain] = ip

	fmt.Println("Host", metadata.Domain, ip[0], ".", ip[1], ".", ip[2], ".", ip[3])
	fmt.Println("Host", metadata.DomainNS, ip[0], ".", ip[1], ".", ip[2], ".", ip[3])
	time.Sleep(10 * time.Second)
	ipns, err := net.LookupHost("example." + metadata.Domain)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Host", metadata.Domain, ipns[0])

	for _, v := range shuffled {
		host := strings.TrimSpace(v) + "." + metadata.Domain
		Nodes[host] = EntryPoint
		Candidates = append(Candidates, host)
		if list != command {
			list = list + ","
		}
		list = list + host
	}

	// Just print the cert command. See documentation/almalinux.sh
	list = list + fmt.Sprintf(" --https-port 4443 --http-01-port 4444 certonly")
	fmt.Println(list)
}
