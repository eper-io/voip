package eos

import (
	"fmt"
	"gitlab.com/eper.io/engine/oraclecloud"
	"math/rand"
	"time"
)

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// The industry-leading Mitosis Algorithm of eper.io

// Problem. We need a distributed scaling logic that is
// fail proof for quantum grade security.
// i.e. no single point or hot spot of failure in the age of quantum computers
// Solution.
// Each machine has an expected runtime and an expected number of sessions,
// after which it terminates normally.

// Do a mitosis check at 25%, 50%, and 75% of expected sessions.
// If the runtime that elapsed is less than 25%, 50% and 75% respectively, we create a new node (mitosis).
// We terminate the node at 100%+session runtime (25%) not starting new sessions after 100%.
// We do not terminate on the system that is marked as root, but recreate expired nodes.

// Discussion.
// An automated termination logic helps to resolve sporadic issues and data leaking.
// Mitosis checks can help to address quick demand increases.
// An automated termination also handles automated scaling down.
// The logic is nice as it needs only knowledge of local demand (sessions) to scale the cluster up and down.

// We need a random load balancer to help with this of the session launches.

// Problem. Kubernetes already does this.
// Solution. Yes, but it uses millions of lines of code.
// This is easier to learn as a starter before moving on.
// It is easier to security review, if it is shorter.
// Trick: it saves time and money for security reviewers who work for a fixed fee vs. feature count based fee.

// Where does the name come from?
// Bacteria do mitosis and grow in numbers if they find enough food, e.g. sugar.
// This algorithm is similar relying only on local information that is easy to measure.

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
// This will avoid a sudden load a single node.
// This will make it hard to predict the next node for malware.
// This will be debugged easily assuming any node can be targeted next.
// This helps architects to assume any node can be loaded next.
// This reduces costs keeping the algorithm proven and low cost to maintain.
// Example
// oci compute instance launch --compartment-id 'ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq' --availability-domain 'lynu:US-SANJOSE-1-AD-1' --shape 'VM.Standard.A1.Flex' --image-id 'ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq' --subnet-id 'ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q' --shape-config '{"ocpus":"4"}'
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
//
// It is similar to how amoebas and bacteria are splitting, if they are getting enough nutrients.
// Setup. Here is the basic calculation to set the two parameters.
// A single 1 cpu node is capable of handling 1000 MBps.
// Assume we handle sessions that require 10 MBps dedicated lines.
// We also assume that sessions last max 2 hours on the line,
// and we have to reserve for the entire session.
// We limit launches for at most eight hours to avoid any garbage piling up.
// We assume a safety utility ratio of 50%, that is we actually keep the containers half unused.
// This is especially good for very low latency communication like Audio or User Interface.
// Such a container can handle a load of
// (8 hrs / 2 hrs * 1000 MBps / 10 MBps) * 50% = 200 sessions of the entire 8+2 hrs lifetime.
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
//
// One benefit is a logic simple to understand, teach and debug.
// Another benefit is a quick scale up, if it is needed i.e. we trigger a check on the session count not the time.
// If we get a thousand sudden requests we scale up quickly.
// TODO Ideally the nodes should terminate themselves and remove from the cloud in the future to be super reliable.
func Mitosis() {
	id, host, ip := oraclecloud.LaunchInstance(maxRuntime)
	if id == "" {
		fmt.Println("Failed Launch")
		return
	}
	fmt.Println("Launched", id, host, ip)

	launches[id] = 0
	fqdn[id] = host
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

// Terminate terminates the instance.
func Terminate(id string, host string) {
	singleton := len(launches) <= 1
	if singleton {
		go Mitosis()
	}
	oraclecloud.TerminateInstance(id, host)
	delete(launches, id)
}
