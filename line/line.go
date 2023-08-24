package line

import (
	"bytes"
	"fmt"
	"gitlab.com/eper.io/engine/metadata"
	"gitlab.com/eper.io/engine/websocket"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// main.go - Relay websocket connections for streaming applications. See README.md for usage

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
	port := os.Getenv("PORT")
	if port != "" && !strings.HasSuffix(metadata.SiteUrl, "ondigitalocean.app") {
		metadata.SiteUrl = metadata.SiteUrl + ":" + port
		fmt.Printf("Try wired on DigitalOcean App %s\n", metadata.SiteUrl+"/line.html?apikey="+metadata.ActivationKey+"#generate_leaf")
		fmt.Printf("Try mobile on DigitalOcean App %s\n", metadata.SiteUrl+"/line.html?apikey="+metadata.ActivationKey+"&mobile=1#generate_leaf")
	} else {
		if strings.Contains(metadata.SiteUrl, "example.com") {
			fmt.Printf("Set SITEURL In DigitalOcean App Spec like")
			fmt.Printf("- key: SITEURL")
			fmt.Printf("  scope: RUN_AND_BUILD_TIME")
			fmt.Printf("  value: https://voip-2-g9p9u.ondigitalocean.app")
		}

		fmt.Printf("Try wired %s\n", metadata.SiteUrl+"/line.html?apikey="+metadata.ActivationKey+"#generate_leaf")
		fmt.Printf("Try mobile %s\n", metadata.SiteUrl+"/line.html?apikey="+metadata.ActivationKey+"&mobile=1#generate_leaf")
	}

	http.Handle("/ws", websocket.Handler(relayLine))

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		http.Redirect(writer, request, "/line.html", http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/line.html", func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(15 * time.Millisecond)
		var mobile string
		if request.URL.Query().Get("mobile") != "" {
			mobile = "&mobile=1"
		}
		if metadata.ActivationKey == "" {
			// First participant
			metadata.ActivationKey = GenerateUniqueKey()
			http.Redirect(writer, request, fmt.Sprintf("/line.html?apikey=%s%s#generate_leaf", metadata.ActivationKey, mobile), http.StatusTemporaryRedirect)
			return
		}
		if "" == request.URL.Query().Get("apikey") {
			if strings.Contains(metadata.SiteUrl, "example.com") {
				// Just a test example
				http.Redirect(writer, request, fmt.Sprintf("/line.html?apikey="+metadata.ActivationKey), http.StatusTemporaryRedirect)
			} else {
				http.Redirect(writer, request, fmt.Sprintf("/documentation.html"), http.StatusTemporaryRedirect)
			}
			return
		}
		if mobile != "" {
			buffer, _ := os.ReadFile("./line/res/line.html")
			_, _ = io.WriteString(writer, strings.Replace(string(buffer), "200000", "6000", 1))
			return
		}
		dynamicContentPart(writer, request, path.Join("./line/res/line.html"))
	})

	http.HandleFunc("/news.html", func(writer http.ResponseWriter, request *http.Request) {
		//dynamicContentPart(writer, request, path.Join("./line/res/news.html"))
		content, _ := os.ReadFile(path.Join("./line/res/news.html"))
		var host, participant string

		var i int
		lock.Lock()
		i = len(peer)
		lock.Unlock()
		if i > 0 {
			host = "🎧"
		}
		if i > 1 {
			participant = "🎧"
		}

		metadata.Info = fmt.Sprintf("%s %s %s", host, metadata.Bandwidth, participant)

		text := strings.Replace(string(content), "Info.", metadata.Info, 1)
		_, _ = io.WriteString(writer, text)
	})

	http.HandleFunc("/info", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "PUT" {
			apiKey := request.URL.Query().Get("apikey")
			if apiKey == metadata.ActivationKey {
				share, _ := io.ReadAll(request.Body)
				_ = os.WriteFile(path.Join("/tmp", metadata.ActivationKey), share, 0700)
			}
		}
		if request.Method == "GET" {
			share, _ := os.ReadFile(path.Join("/tmp", metadata.ActivationKey))
			if share != nil && len(share) > 0 {
				_, _ = io.Copy(writer, bytes.NewBuffer(share))
			}
		}
	})

	http.HandleFunc("/moose.png", func(writer http.ResponseWriter, request *http.Request) {
		dynamicContentPart(writer, request, path.Join("./line/res/moose.png"))
	})

	http.HandleFunc("/main.css", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/css")
		dynamicContentPart(writer, request, path.Join("./line/res/main.css"))
	})

	http.HandleFunc("/e2e.js", func(writer http.ResponseWriter, request *http.Request) {
		dynamicContentPart(writer, request, path.Join("./line/res/e2e.js"))
	})

	http.HandleFunc("/pcm.js", func(writer http.ResponseWriter, request *http.Request) {
		dynamicContentPart(writer, request, path.Join("./line/res/pcm.js"))
	})

	http.HandleFunc("/documentation.html", func(writer http.ResponseWriter, request *http.Request) {
		dynamicContentPart(writer, request, path.Join("./line/res/documentation.html"))
	})

	http.HandleFunc("/terms.html", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(writer, request, "./metadata/terms.html")
	})

	http.HandleFunc("/contact.html", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(writer, request, "./metadata/contact.html")
	})
}

func dynamicContentPart(writer http.ResponseWriter, request *http.Request, file string) {
	writer.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(writer, request, file)
}

func relayLine(ws *websocket.Conn) {
	//TODO TCP is not exposed here. set maxbuffer, timeout
	//if tcpConn, ok := ws.Request().TLS.(*net.TCPConn); ok {
	//	err := tcpConn.SetNoDelay(true)
	//}
	start := time.Now()
	var lastMillis float64
	var amount float64
	ws.PayloadType = websocket.BinaryFrame
	defer func() { _ = ws.Close() }()
	apiKey := ws.Request().URL.Query().Get("apikey")
	if metadata.ActivationKey != apiKey {
		time.Sleep(15 * time.Millisecond)
		fmt.Println("apikey mismatch error")
		return
	}
	lock.Lock()
	peers := len(peer)
	lock.Unlock()
	if peers > 1 {
		fmt.Println("connections cannot exceed 2")
		return
	}
	lock.Lock()
	peer[ws] = 1
	lock.Unlock()
	fmt.Println("connection received")

	finished := false
	for !finished {
		buf := make([]byte, math.MaxUint16)
		n, _ := ws.Read(buf)
		if n == 0 {
			buf = nil
			fmt.Println("connection lost")
			finished = true
			break
		}
		amount = amount + float64(n)
		lock.Lock()
		idleSince = time.Now()
		for cr := range peer {
			if cr != ws || len(peer) == 1 {
				m, _ := cr.Write(buf[0:n])
				currentMillis := float64(time.Now().Sub(start).Milliseconds())
				if currentMillis > lastMillis+10000.0 {
					mbps := 8 * amount / 1000 / (currentMillis - lastMillis)
					metadata.Bandwidth = "❚️"
					for i := 0; i < 5; i++ {
						if mbps > float64(i)*0.25 {
							metadata.Bandwidth = metadata.Bandwidth + "❚️"
						}
					}
					amount = 0
					lastMillis = currentMillis
				}
				if m == 0 {
					finished = true
					fmt.Println("connection closed")
					break
				}
			}
		}
		lock.Unlock()
	}

	lock.Lock()
	delete(peer, ws)
	peers = len(peer)
	lock.Unlock()

	if peers == 0 {
		idleSince = time.Now()
		fmt.Println("lost all connections")
		metadata.Info = ""
		go func() {
			const timeout = 3 * time.Minute
			time.Sleep(timeout + 1*time.Minute)
			if time.Now().Sub(idleSince) > timeout {
				fmt.Printf("Idle shutdown after %d minutes. Restarting or recycling based on docker rules.", int(timeout.Minutes()))
				os.Exit(0)
			}
		}()
	}
}

func RedactPublicKey(uq string) string {
	if uq == "" {
		return ""
	}
	return uq[0:6]
}

func Random() uint32 {
	buf := make([]byte, 4)
	n, err := rand.Read(buf)
	if err != nil || n != 4 {
		return 0
	}
	return uint32(buf[0])<<24 | uint32(buf[0])<<16 | uint32(buf[0])<<8 | uint32(buf[0])<<0
}

func GenerateUniqueKey() string {
	// So we do not add much of a header suggesting it is the best solution.
	// Adding a header would increase the chance of randomly testing the
	// private key with sites to verify it works, practically leaking it.
	// Your internal context should tell where an api key is valid.

	// TODO Need to get a better seed from the internet
	x, _ := os.Stat(os.Args[0])
	seed := time.Now().UnixNano() ^ x.ModTime().UnixNano()
	rand.Seed(seed)

	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	key := make([]rune, 92)
	salt := metadata.RandomSalt
	for i := 0; i < 1000; i++ {
		for i := 0; i < 92; i++ {
			key[i] = letters[(((Random() ^ rand.Uint32()) + uint32(salt[i])) % uint32(len(letters)))]
			time.Sleep(550 * time.Nanosecond)
		}
		if key[91] == 'A' {
			break
		}
	}
	key[91] = 'R'
	return string(key)
}
