package cloud

import (
	"testing"
)

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

func TestOracle(t *testing.T) {
	jsonData := `
	{
		"data": {
			"id": "ocid1.instance.oc1.us-sanjose-1.anzwuljra36or2qc3tr5zn5t627exr26nschfblbbbtpjd7daopuxeid5kqa"
		}
	}
	`

	id := ParseInstanceId(jsonData)

	if id != "ocid1.instance.oc1.us-sanjose-1.anzwuljra36or2qc3tr5zn5t627exr26nschfblbbbtpjd7daopuxeid5kqa" {
		t.Error(id)
	}

	jsonData = `
	{
	  "data": [
		{
		  "availability-domain": "lynu:US-SANJOSE-1-AD-1",
		  "compartment-id": "ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq",
		  "defined-tags": {
			"Oracle-Tags": {
			  "CreatedBy": "oracleidentitycloudservice/hq@schmied.us",
			  "CreatedOn": "2023-06-13T17:40:12.748Z"
			}
		  },
		  "display-name": "instance-20230613-1037",
		  "freeform-tags": {},
		  "hostname-label": "instance-20230613-1037",
		  "id": "ocid1.vnic.oc1.us-sanjose-1.abzwuljr2qnpwmsdwb3yqzhfu3ygkrdmy4zquzw3kx2hae5c52m4iciwxwgq",
		  "is-primary": true,
		  "lifecycle-state": "AVAILABLE",
		  "mac-address": "02:00:17:01:29:B7",
		  "nsg-ids": [],
		  "private-ip": "10.0.0.141",
		  "public-ip": "152.67.231.197",
		  "skip-source-dest-check": false,
		  "subnet-id": "ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q",
		  "time-created": "2023-06-13T17:40:21.971000+00:00",
		  "vlan-id": null
		}
	  ]
	}
	`

	ip := ParsePublicIP(jsonData)

	if ip != "152.67.231.197" {
		t.Error(ip)
	}

}
