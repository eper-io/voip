package eos

import (
	"time"
)

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

var launches = map[string]int64{}

// maxSessions is the maximum number of sessions during the lifetime of the node.
// See Mitosis() for explanation.
const maxSessions = int64(200)

// maxRuntime is the maximum time a node is available for launch requests
const maxRuntime = 1 * time.Hour //TODO 8
const InitialNodes = 2

const BasePort = 49999
const LastPort = 60000

var lastContainer = BasePort
var MaxContainerTime = 2 * time.Hour
var DockerDelay = 4 * time.Second
