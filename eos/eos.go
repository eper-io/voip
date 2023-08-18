package eos

import (
	"fmt"
	"gitlab.com/eper.io/engine/metadata"
	"io"
	"math/rand"
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

	//go Mitosis()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(writer, request, "./eos/res/launch.html")
	})
	http.HandleFunc("/englang", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(writer, request, "/var/log/voip")
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
		lock.Lock()
		time.Sleep(15 * time.Millisecond)
		if metadata.ActivationKey == "" || apiKey != metadata.ActivationKey {
			//writer.WriteHeader(http.StatusPaymentRequired)
		}
		lock.Unlock()
		lock.Lock()
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
		lock.Unlock()

		key := generateUniqueKey()
		startCommand := exec.Command("podman", "run", "--timeout", fmt.Sprintf("%d", int(MaxContainerTime.Seconds())), "-d", "--rm", "--name", redactPublicKey(key), "-e", fmt.Sprintf("PORT=%d", port), "-e", "APIKEY="+key, "-p", fmt.Sprintf("%d:443", port), "-v", metadata.Certificate+":/tmp/fullchain.pem:ro", "-v", metadata.PrivateKey+":/tmp/privkey.pem:ro", metadata.ContainerRuntime)
		fmt.Println(startCommand.String())
		y, _ := startCommand.CombinedOutput()
		fmt.Println(string(y))

		go func() {
			lock.Lock()
			x := make([]string, 0)
			for k := range launches {
				x = append(x, k)
			}
			if len(x) > 0 {
				pick := x[rand.Intn(len(x))]
				launches[pick]++
			}
			lock.Unlock()
		}()

		//TODO proxy

		time.Sleep(DockerDelay)
		mobile := request.URL.Query().Get("mobile")
		if mobile != "" {
			mobile = "&mobile=1"
		}
		redirect := request.URL.Query().Get("redirect")
		newLine := fmt.Sprintf("%s:%d/line.html?apikey=%s%s#generate_leaf", metadata.SiteUrl, port, key, mobile)
		if redirect != "1" {
			_, _ = io.WriteString(writer, newLine)
		} else {
			writer.Header().Set("Location", newLine)
			writer.WriteHeader(http.StatusTemporaryRedirect)
		}
	})
}
