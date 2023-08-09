package edns

import (
	"fmt"
	"net"
	_ "net"
	"net/dns"
)

func SetupEDNS() {
	// Create a DNS server.
	server := &dns.Server{Addr: ":53", Net: "udp"}

	// Create a map of records.
	records := map[string]string{
		"www.example.com":  "192.168.0.1",
		"mail.example.com": "192.168.0.2",
	}

	// Register the records with the DNS server.
	for name, ip := range records {
		dns.HandleFunc(name, func(w dns.ResponseWriter, r *dns.Msg) {
			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", name, ip))
			if err != nil {
				fmt.Println(err)
				return
			}

			w.WriteMsg(dns.Msg{Answer: []dns.RR{rr}})
		})
	}

	// Listen for DNS requests.
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func serveDNS(u *net.UDPConn, clientAddr net.Addr, request *layers.DNS) {
	replyMess := request
	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	var ip string
	var err error
	var ok bool
	ip, ok = records[string(request.Questions[0].Name)]
	if !ok {
		//Todo: Log no data present for the IP and handle:todo
	}
	a, _, _ := net.ParseCIDR(ip + "/24")
	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.IP = a
	dnsAnswer.Name = []byte(request.Questions[0].Name)
	fmt.Println(request.Questions[0].Name)
	dnsAnswer.Class = layers.DNSClassIN
	replyMess.QR = true
	replyMess.ANCount = 1
	replyMess.OpCode = layers.DNSOpCodeNotify
	replyMess.AA = true
	replyMess.Answers = append(replyMess.Answers, dnsAnswer)
	replyMess.ResponseCode = layers.DNSResponseCodeNoErr
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
	err = replyMess.SerializeTo(buf, opts)
	if err != nil {
		panic(err)
	}
	u.WriteTo(buf.Bytes(), clientAddr)
}
