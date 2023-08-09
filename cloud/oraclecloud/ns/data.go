package ns

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// Nodes are the main DNS entries. ipv4 only for now.
var Nodes = map[string][4]byte{}

// Candidates are the list of available host names.
var Candidates = make([]string, 0)
