package eos

import (
	"fmt"
	"gitlab.com/eper.io/engine/oraclecloud/cloud"
	"math/rand"
	"time"
)

// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// Problem. We need a distributed scaling logic that is
// fail proof for quantum grade security.
// i.e. no single point or hot spot of failure in the age of quantum computers
// Solution.
// Each machine has an expected runtime, after which it terminates normally.
// Do a mitosis check at 25%, 50%, and 75% of expected runtime.
// If load is at least 25%, 50% and 75% at these points, we create a new node.
// We terminate the node at 100% not issuing new sessions after 75%.
// We do not terminate on the system that is marked as root, as a long haul test.

// Discussion.
// An automated termination logic helps to resolve sporadic issues and data leaking.
// Mitosis checks can help to address quick demand increases.
// An automated termination also handles automated scaling down.
// The logic is nice as it needs only local demand knowledge to scale the cluster up and down.
// We need a random load balancer to achieve this for the session starts.
// Problem. Kubernetes already does this.
// Solution. Yes, but it uses millions of lines of code. This is easier to learn as a starter before moving on.
// It is easier to security review, if it is shorter.
// Trick: it saves time and money for security reviewers who work for a fixed fee vs. feature count based fee.

//TODO
//oci compute instance launch --compartment-id 'ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq' --availability-domain 'lynu:US-SANJOSE-1-AD-1' --shape 'VM.Standard.A1.Flex' --image-id 'ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq' --subnet-id 'ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q' --shape-config '{"ocpus":"4"}'

// SetupMitosis is an algorithm that sets up the Mitosis algorithm.
// It checks if the load is big enough to trigger the scale out logic.
// It also scales down, once the machine runtime has expired.
func SetupMitosis() {
	for i := 0; i < InitialNodes; i++ {
		go Mitosis()
	}
}

// LaunchSite gets a site for the current request.
// We generally choose a random method to discourage malware and regular hotspots.
// This will avoid sudden loads on single nodes.
// This will make it hard to predict the next node for malware.
// This will be debugged easily assuming any node can be targeted next.
// This helps architects to assume any node can be directed next.
// This reduces costs keeping the algorithm proven and low cost to maintain.
func LaunchSite() {
	x := make([]string, 0)
	for k := range launches {
		x = append(x, k)
	}
	if len(x) > 0 {
		pick := x[rand.Intn(len(x))]
		launches[pick]++
	}
}

// Mitosis checks and executes a scaling logic inspired by millions of years of evolution.
// It is similar to how amoebas and bacteria are splitting getting enough nutrients.
// Setup. Here is the basic calculation to set the two parameters.
// A single 1 VCPU node is capable of handling 1000 MBps.
// Assume we handle sessions that require 10 MBps dedicated lines.
// We also assume that sessions last max 2 hours on the line,
// and we have to reserve for the entire session.
// We limit launches for at most eight hours to avoid any garbage piling up.
// We assume a utility ratio of 50%.
// Such a container can handle a load of
// 8 hrs / 2 hrs * 1000 MBps / 10 MBps * 50% = 200 sessions of the entire 8+2 hrs lifetime.
// That allows 1000 MBps / 10 MBps * 50% = 50 concurrent sessions.
//
// The final logic then is the following.
// If we have not passed 2 hrs, when 50 sessions have been launched, we carry out mitosis.
// If we have not passed 4 hrs, when 100 sessions have been launched, we carry out mitosis.
// If we have not passed 6 hrs, when 150 sessions have been launched, we carry out mitosis.
// If we have not passed 8 hrs, when 200 sessions have been launched, we carry out mitosis.
// We carry out an outstanding mitosis at 8 hours to a followup node.
// We forward all local launches to the followup node after 8 hours.
// We terminate the current instance at 10 hours, since all session reservations expired.
// We terminate the current instance at max sessions, since the state is statistically dirty/exhausted.

func Mitosis() {
	id, host, ip := cloud.LaunchInstance()
	if id == "" {
		fmt.Println("Failed Launch")
		return
	}
	fmt.Println("Launched", id, host, ip)
	cloud.CleanupInstance(id, host, maxRuntime)
	if id == "" {
		fmt.Println("failed mitosis")
		return
	}

	launches[id] = 0
	go func() {
		start := time.Now()
		last := launches[id]
		for {
			time.Sleep(10 * time.Minute)
			current := launches[id]
			fmt.Println("Mitosis from", last, "to", current, "at", int64(time.Now().Sub(start).Seconds()), "seconds")
			for i := int64(1); i <= 4; i++ {
				if current > maxSessions || time.Now().Sub(start) > maxRuntime {
					Terminate(id, host)
					singleton := len(launches) == 1
					if singleton {
						go Mitosis()
					}
					return
				}
				if last < maxSessions*i/4 && current >= maxSessions*i/4 {
					if int64(time.Now().Sub(start).Seconds()) < int64(maxRuntime.Seconds())*i/4 {
						// We used the quota too fast: need more
						fmt.Println("Mitosis needed from", last, "to", current, "at", int64(time.Now().Sub(start).Seconds())*i/4, "seconds", i, "/4")
						go Mitosis()
						break
					}
				}
			}

			last = current
		}
	}()
}

func Terminate(id string, host string) {
	cloud.TerminateInstance(id, host)
	delete(launches, id)
}
