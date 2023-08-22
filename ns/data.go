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

// EntryPoint is the ip address of the main entry point.
// We need to be more distributed in the future.
// This allows lower latency on the other hand.
var EntryPoint = [4]byte{0, 0, 0, 0}

// HostNames are sample easy to pronounce host names. Cerbot can handle ~70
var HostNames = `apple
bicycle
sunshine
whisper
ocean
lantern
rainbow
adventure
galaxy
chocolate
dragon
elephant
moonlight
jazz
telescope
enigma
firefly
harmony
carousel
blossom
wanderlust
serenade
mountain
symphony
velvet
bubble
stardust
flamingo
aurora
puzzle
horizon
quasar
twilight
radiant
meadow
tornado
velvet
whisper
avalanche
carnival
rhapsody
mirage
monsoon
universe
jubilee
sapphire
zephyr
`
