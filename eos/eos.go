package eos

import (
	"bytes"
	"fmt"
	"gitlab.com/eper.io/engine/metadata"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

func Setup() {
	key := os.Getenv("APIKEY")
	if key != "" {
		metadata.ActivationKey = key
	}
	siteUrl := os.Getenv("SITEURL")
	if siteUrl != "" {
		metadata.SiteUrl = siteUrl
	}
	runt := os.Getenv("DOCKERIMAGE")
	if runt != "" {
		metadata.ContainerRuntime = runt
	}
	if metadata.ContainerRuntime == "" {
		fmt.Println("docker image not specified")
	}

	url1, _ := url.Parse(metadata.SiteUrl)
	metadata.Certificate = fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", url1.Hostname())
	metadata.PrivateKey = fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", url1.Hostname())

	fmt.Printf("Launch %s as %s?apikey=%s\n", metadata.ContainerRuntime, metadata.SiteUrl, metadata.ActivationKey)
	fmt.Printf("")

	SetupMitosis()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(writer, request, "./eos/res/launch.html")
	})
	http.HandleFunc("/englang", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		// http.ServeFile(writer, request, "/var/log/voip")

		ret := bytes.Buffer{}
		x := launches
		for k, v := range x {
			hostname := fqdn[k]
			ret.WriteString(fmt.Sprintf("%s had %d launches\n", hostname, v))
		}
		_, _ = io.Copy(writer, &ret)
		startupInfo, _ := os.ReadFile("/var/log/voip")
		_, _ = io.Copy(writer, bytes.NewBuffer(startupInfo))
	})
	http.HandleFunc("/main.css", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		writer.Header().Set("Content-Type", "text/css")
		http.ServeFile(writer, request, "./eos/res/main.css")
	})
	http.HandleFunc("/moose.png", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(writer, request, "./eos/res/moose.png")
	})
	http.HandleFunc("/documentation.html", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(writer, request, "./line/res/documentation.html")
	})
	http.HandleFunc("/terms.html", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(writer, request, "./metadata/terms.html")
	})
	http.HandleFunc("/contact.html", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(writer, request, "./metadata/contact.html")
	})

	http.HandleFunc("/idle", func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.URL.Query().Get("apikey")
		mobile := request.URL.Query().Get("mobile")
		redirect := request.URL.Query().Get("redirect")

		proxy := false
		for _, x := range os.Args {
			if x == "proxy" {
				proxy = true
			}
		}
		if proxy {
			pick := LaunchSite()
			if pick != "" {
				hostFqdn := fqdn[pick]
				if hostFqdn != "" {
					mobile1 := ""
					if mobile == "1" {
						// This one is forced to low bandwidth, but it is not studio quality
						mobile1 = "&mobile=1"
					}

					pickedUrl := fmt.Sprintf("https://%s%s?apikey=%s&redirect=0%s", hostFqdn, request.URL.Path, apiKey, mobile1)
					retx, _ := http.Get(pickedUrl)
					container, _ := io.ReadAll(retx.Body)
					newLine := string(container)
					if newLine != "" {
						fmt.Println("proxy to", hostFqdn, pickedUrl)
						launches[pick]++
						// Proxy
						if redirect != "1" {
							_, _ = io.WriteString(writer, newLine)
						} else {
							writer.Header().Set("Location", newLine)
							writer.WriteHeader(http.StatusTemporaryRedirect)
						}
						return
					} else {
						fmt.Println("cannot proxy connect to", hostFqdn, pickedUrl)
					}
				} else {
					fmt.Println("cannot proxy address to", hostFqdn)
				}
			} else {
				fmt.Println("proxy list error. local fallback")
			}
		} else {
			fmt.Println("cannot proxy. local fallback")
		}
		// Local fallback

		fmt.Println("line creation")
		time.Sleep(15 * time.Millisecond)
		if metadata.ActivationKey == "" || apiKey != metadata.ActivationKey {
			// TODO
			//writer.WriteHeader(http.StatusPaymentRequired)
		}

		port := BasePort
		for ; lastContainer < LastPort; lastContainer++ {
			x, err := net.Listen("tcp", fmt.Sprintf(":%d", lastContainer))
			if err == nil {
				_ = x.Close()
				port = lastContainer
				break
			}
			if x != nil {
				_ = x.Close()
			}
		}

		fmt.Println(port)
		key := generateUniqueKey()
		startCommand := exec.Command("podman", "run", "--timeout", fmt.Sprintf("%d", int(MaxContainerTime.Seconds())), "-d", "--rm", "--name", redactPublicKey(key), "-e", fmt.Sprintf("PORT=%d", port), "-e", "APIKEY="+key, "-p", fmt.Sprintf("%d:443", port), "-v", metadata.Certificate+":"+metadata.Certificate+":ro", "-v", metadata.PrivateKey+":"+metadata.PrivateKey+":ro", metadata.ContainerRuntime)
		fmt.Println(startCommand.String())
		returned, _ := startCommand.CombinedOutput()
		fmt.Println(string(returned))

		time.Sleep(DockerDelay)
		if mobile == "1" {
			// This one is forced to low bandwidth, but it is not studio quality
			mobile = "&mobile=1"
		}
		newLine := fmt.Sprintf("%s:%d/line.html?apikey=%s%s#generate_leaf", metadata.SiteUrl, port, key, mobile)
		if redirect != "1" {
			_, _ = io.WriteString(writer, newLine)
		} else {
			writer.Header().Set("Location", newLine)
			writer.WriteHeader(http.StatusTemporaryRedirect)
		}
	})
}
