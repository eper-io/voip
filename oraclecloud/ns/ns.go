package ns

import (
	"bufio"
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"net"
	"net/http"
	"os"
	"strings"
)

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// SetupNameServer starts a traditional name server.
// This might require running the tool as root.
// Ports lower than 1024 are usually protected.
func SetupNameServer() {
	// TODO Debugging
	Nodes["www.example.com"] = [4]byte{1, 2, 3, 4}

	http.HandleFunc("/dns", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPut {
			if strings.HasPrefix(request.RemoteAddr, "127.0.0.1") {
				// Internal use
				host := request.URL.Query().Get("a")
				ipv4 := request.URL.Query().Get("ipv4")
				ipv4d := [4]byte{0, 0, 0, 0}
				n, _ := fmt.Sscanf(ipv4, "%d.%d.%d.%d", &ipv4d[0], &ipv4d[1], &ipv4d[2], &ipv4d[3])
				if host != "" && n == 4 {
					// Only once
					_, ok := Nodes[host]
					if !ok {
						Nodes[host] = ipv4d
					}
					return
				}
			}
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if request.Method == http.MethodGet {
			host := request.URL.Query().Get("a")
			if host == "" {
				buffered := bufio.NewWriter(writer)
				for k, v := range Nodes {
					_, _ = buffered.WriteString(fmt.Sprintf("%s %s\n", k, v))
				}
				_ = buffered.Flush()
				return
			}
			ip, ok := Nodes[host]
			if ok {
				_, _ = writer.Write([]byte(fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])))
				return
			}

		}
	})

	go func() {
		fmt.Println("Try\ncurl -X PUT 'http://127.0.0.1:7777/dns?a=www.hello.com&ipv4=127.127.127.127' && curl -X GET 'http://127.0.0.1:7777/dns?a=www.hello.com' && dig @127.0.0.1 www.hello.com")
		_ = http.ListenAndServe(":7777", nil)
	}()

	go func() {
		listener, err := net.ListenPacket("udp", ":53")
		if err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("%v\n", err.Error()))
			return
		}
		for {
			buf := make([]byte, 1024)
			n, addr, err := listener.ReadFrom(buf)
			if err != nil {
				continue
			}
			var m dnsmessage.Message
			err = m.Unpack(buf[0:n])
			if err != nil {
				continue
			}
			//fmt.Printf("%v %v r %v a %v\n", addr, m, m.RecursionDesired, m.Authoritative)
			for _, q := range m.Questions {
				{
					var a = dnsmessage.Resource{}
					a.Header.Type = q.Type
					a.Header.Class = q.Class
					a.Header.Name = q.Name
					a.Header.TTL = 300
					m.RecursionAvailable = true
					m.Authoritative = true
					m.Response = true
					m.Truncated = false
					for k, v := range Nodes {
						if q.Type == dnsmessage.TypeA && strings.ToLower(q.Name.String()) == k+"." {
							a.Body = &dnsmessage.AResource{A: v}
							m.Answers = append(m.Answers, a)
						}
					}
				}

				// TODO
				//{
				// var a = dnsmessage.Resource{}
				// a.Header.Type = dnsmessage.TypeCNAME
				// a.Header.Class = dnsmessage.ClassANY
				// a.Header.Name = q.Name
				// name, _ := dnsmessage.NewName("example.com")
				// a.Body = &dnsmessage.CNAMEResource{name}
				// m.Answers = append(m.Answers, a)
				//}
			}
			ret, err := m.Pack()
			if err != nil {
				continue
			}
			_, err = listener.WriteTo(ret, addr)
			if err != nil {
				continue
			}
		}
	}()

}
