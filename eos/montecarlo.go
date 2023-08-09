package eos

import (
	"gitlab.com/eper.io/engine/metadata"
	"math/rand"
	"os"
	"time"
)

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

func redactPublicKey(uq string) string {
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

func generateUniqueKey() string {
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
