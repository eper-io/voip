package cloud

import (
	"encoding/json"
	"fmt"
	"gitlab.com/eper.io/engine/line"
	"gitlab.com/eper.io/engine/oraclecloud/metadata"
	"gitlab.com/eper.io/engine/oraclecloud/ns"
	"os"
	"os/exec"
	"path/filepath"
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

// Commands to launch and terminate instances on Oracle cloud.
// We assume oci api key is installed.
// oci compute instance launch --compartment-id 'ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq' --availability-domain 'lynu:US-SANJOSE-1-AD-1' --shape 'VM.Standard.A1.Flex' --image-id 'ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq' --subnet-id 'ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q' --shape-config '{"ocpus":"4"}'

func ParseInstanceId(jsonData string) string {
	type Data struct {
		ID string `json:"id"`
	}

	type JSONData struct {
		Data Data `json:"data"`
	}

	var parsedData JSONData
	err := json.Unmarshal([]byte(jsonData), &parsedData)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	return parsedData.Data.ID
}

type InstanceVNIC struct {
	Data []struct {
		PublicIP string `json:"public-ip"`
	} `json:"data"`
}

func ParsePublicIP(jsonData string) string {
	var instance InstanceVNIC
	err := json.Unmarshal([]byte(jsonData), &instance)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	if len(instance.Data) > 0 {
		publicIP := instance.Data[0].PublicIP
		return publicIP
	} else {
		fmt.Println("No data found")
	}
	return ""
}

func OciCommand(name string, arg []string) *exec.Cmd {
	cmd := &exec.Cmd{
		Path: name,
		Args: append([]string{name}, arg...),
	}
	if filepath.Base(name) == name {
		lp, err := exec.LookPath(name)
		if err == nil && lp != "" {
			cmd.Path = lp
		}
		if err != nil {
			cmd.Err = err
		}
	}
	return cmd
}

func LaunchInstance() (instanceId string, host string, ip string) {
	cmdx := metadata.OracleLaunchCommand
	c := OciCommand("sh", strings.Split(cmdx, " ")[1:])
	ret, _ := c.Output()
	if len(ret) == 0 {
		fmt.Println(strings.Join(c.Args, " "))
	}
	fmt.Println(string(ret))
	fmt.Println(string(ret))
	id := ParseInstanceId(string(ret))
	if id != "" {
		start := time.Now()
		for time.Now().Sub(start).Seconds() < 5*time.Minute.Seconds() {
			time.Sleep(10 * time.Second)
			ipv4d := [4]byte{0, 0, 0, 0}
			ipv4 := GetInstancePublicIp(id)
			n, _ := fmt.Sscanf(ipv4, "%d.%d.%d.%d", &ipv4d[0], &ipv4d[1], &ipv4d[2], &ipv4d[3])
			x := ns.Candidates
			for _, host := range x {
				value, ok := ns.Nodes[host]
				if !ok || value == EntryPoint {
					if host != "" && n == 4 {
						ns.Nodes[host] = ipv4d
						return id, host, ipv4
					}
				}
			}
		}
	}
	return id, "", ""
}

func GetInstancePublicIp(instance string) string {
	cmdx := fmt.Sprintf("oci compute instance list-vnics --instance-id %s", instance)
	c := OciCommand("oci", strings.Split(cmdx, " ")[1:])
	ret, _ := c.Output()
	if len(ret) == 0 {
		fmt.Println(strings.Join(c.Args, " "))
	}
	fmt.Println("ip", string(ret))
	return ParsePublicIP(string(ret))
}

func TerminateInstance(id string, host string) {
	cmdx := fmt.Sprintf("oci compute instance terminate --force --instance-id %s", id)
	c := OciCommand("oci", strings.Split(cmdx, " ")[1:])
	_, _ = c.Output()
	delete(ns.Nodes, host)
}

func CleanupInstance(id string, host string, duration time.Duration) {
	cmd := fmt.Sprintf("nohup sleep %d && oci compute instance terminate --force --instance-id %s &\n", int64(duration.Seconds()), id)

	name := "/tmp/cleanup_" + line.GenerateUniqueKey()
	_ = os.WriteFile(name, []byte(cmd), 0700)

	go func() {
		fmt.Println("cleanup", name, host, cmd)
		_ = exec.Command("bash", name).Start()
	}()
}
