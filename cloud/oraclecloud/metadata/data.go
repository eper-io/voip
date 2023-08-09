package metadata

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// Domain of the cluster.
var Domain = "oracle.eper.io"

// DomainNS is the main entrypoint and name service
var DomainNS = "l.eper.io"

var OracleLaunchCommand = "oci compute instance launch --user-data-file ./metadata/user-data --compartment-id ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq --availability-domain lynu:US-SANJOSE-1-AD-1 --shape VM.Standard.A1.Flex --image-id ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq --subnet-id ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q --shape-config {\"ocpus\":\"4\"}"

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
	symphony`

//
//velvet
//bubble
//stardust
//flamingo
//aurora
//puzzle
//horizon
//quasar
//twilight
//radiant
//meadow
//tornado
//velvet
//whisper
//avalanche
//carnival
//rhapsody
//mirage
//monsoon
//universe
//jubilee
//sapphire
//zephyr
//orchid
//infinity
//delight
//enchanted
//mirage
//firework
//blossom
//cascade
//celestial
//enigma
//serendipity
//lullaby
//harmony
//galaxy
//whirlwind
//enchanted
//cascade
//rainbow
//nebula
//odyssey
//velvet
//butterfly
//whisper
//reflection
//symphony
//moonbeam
//kaleidoscope
//serenade
//radiance
//secret
//elixir
//aurora
//wanderlust
//twilight
//crystal
//horizon
//starlight
//echo
//zephyr
//jubilee
//breeze
//luminous
//solitude
//ethereal
//galaxy
//purity
//blossom
//whimsical
//euphoria
//ocean
//serenity
//illuminate
//tranquil`
