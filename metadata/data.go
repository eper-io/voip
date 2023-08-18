package metadata

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// data.go - Relay websocket configuration

// You can set this as a default or load from os.Getenv("APIKEY")
var ActivationKey = ""
var SiteUrl = "http://127.0.0.1:7777"
var Certificate = ""
var PrivateKey = ""
var ContainerRuntime = ""

// This requires a Certificate and PrivateKey in place
// Prof I.: it is non-compliant and extremely dangerous to put tls keys to temp directory for bare metal servers.
// Prof II.: temp directory gets cleaned up fast, so it is less likely to get scraped after container exit.
// Prof III.: Do one thing, do it well. Make sure it does not leak, open and safely run it as root.
// 2:1

//var SiteUrl = "https://127.0.0.1:443"
//var Certificate = "/tmp/fullchain.pem"
//var PrivateKey = "/tmp/privkey.pem"

// Use this for local unit testing
var Info = ""
var Bandwidth = ""
var RandomSalt = "ZXNTYCBQINHOIRVUIHNTLADWNZTZVEJDUYAUVKKPUDTYWBONSSRFAOKYFNVMQCZGAEQQLBKGKQHOIIJVOKYAXKONYIBR"
