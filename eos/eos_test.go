package eos

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// Some may say to run the launcher in a non-privileged container and have a separate container that handles docker.
// Our approach is to make the :443 listener get bundled with privileged docker runners.
// Running the docker stuff separately would also require opening a socket or pipe that may be tampered with by other containers.
// However, we must ensure that it remains simple, and it is easy to security review.
// This is what this test does.

func TestProductManager(t *testing.T) {
	var lines int
	list, _ := os.ReadDir(".")
	for _, f := range list {
		if strings.HasSuffix(f.Name(), ".go") {
			content, _ := os.ReadFile(f.Name())
			scanner := bufio.NewScanner(bytes.NewBuffer(content))
			for scanner.Scan() {
				if !strings.HasPrefix(strings.TrimSpace(scanner.Text()), "//") {
					lines++
				}
			}
		}
	}
	if lines > 300 {
		t.Error("Important security consideration. This is running podman or docker with admin rights. We must make sure it stays simple so that security reviews catch issues before deployment.")
	}
}
