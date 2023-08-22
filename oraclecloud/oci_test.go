package oraclecloud

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

func TestNameService(t *testing.T) {
	jsonData := `
{
  "data": {
    "agent-config": {
      "are-all-plugins-disabled": false,
      "is-management-disabled": false,
      "is-monitoring-disabled": false,
      "plugins-config": null
    },
    "availability-config": {
      "is-live-migration-preferred": null,
      "recovery-action": "RESTORE_INSTANCE"
    },
    "availability-domain": "lynu:US-SANJOSE-1-AD-1",
    "capacity-reservation-id": null,
    "compartment-id": "ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq",
    "dedicated-vm-host-id": null,
    "defined-tags": {
      "Oracle-Tags": {
        "CreatedBy": "hq@schmied.us",
        "CreatedOn": "2023-08-22T02:54:36.054Z"
      }
    },
    "display-name": "instance20230822025436",
    "extended-metadata": {},
    "fault-domain": "FAULT-DOMAIN-2",
    "freeform-tags": {},
    "id": "ocid1.instance.oc1.us-sanjose-1.anzwuljra36or2qcxhja7dkkuqfg5lfdzl45mb4bath4ac7yzdj4vcl5e2gq",
    "image-id": "ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq",
    "instance-configuration-id": null,
    "instance-options": {
      "are-legacy-imds-endpoints-disabled": false
    },
    "ipxe-script": null,
    "is-cross-numa-node": false,
    "launch-mode": "PARAVIRTUALIZED",
    "launch-options": {
      "boot-volume-type": "PARAVIRTUALIZED",
      "firmware": "UEFI_64",
      "is-consistent-volume-naming-enabled": true,
      "is-pv-encryption-in-transit-enabled": false,
      "network-type": "PARAVIRTUALIZED",
      "remote-data-volume-type": "PARAVIRTUALIZED"
    },
    "lifecycle-state": "PROVISIONING",
    "metadata": {
      "user_data": "IyEvYmluL3NoCgplY2hvICJIZWxsbyBXb3JsZC4gIFRoZSB0aW1lIGlzIG5vdyAkKGRhdGUgLVIpISIgfCB0ZWUgL3Jvb3Qvb3V0cHV0LnR4dAoKeXVtIGluc3RhbGwgLXkgZ29sYW5nCgpjYXQgPDxFT0YgPi90bXAvbWFpbi5nbwpwYWNrYWdlIG1haW4KCmltcG9ydCAoCiAgICAiZm10IgogICAgIm5ldC9odHRwIgopCgpmdW5jIGhhbmRsZXIodyBodHRwLlJlc3BvbnNlV3JpdGVyLCByICpodHRwLlJlcXVlc3QpIHsKICAgIGZtdC5GcHJpbnRmKHcsICJIZWxsbywgV29ybGQhIikKfQoKZnVuYyBtYWluKCkgewogICAgaHR0cC5IYW5kbGVGdW5jKCIvIiwgaGFuZGxlcikKICAgIHBvcnQgOj0gIjo3Nzc3IgoKICAgIGVyciA6PSBodHRwLkxpc3RlbkFuZFNlcnZlKHBvcnQsIG5pbCkKICAgIGlmIGVyciAhPSBuaWwgewogICAgICAgIGZtdC5QcmludGxuKGVycikKICAgIH0KfQpFT0YKCmdvIHJ1biAvdG1wL21haW4uZ28K"
    },
    "platform-config": null,
    "preemptible-instance-config": null,
    "region": "us-sanjose-1",
    "shape": "VM.Standard.A1.Flex",
    "shape-config": {
      "baseline-ocpu-utilization": null,
      "gpu-description": null,
      "gpus": 0,
      "local-disk-description": null,
      "local-disks": 0,
      "local-disks-total-size-in-gbs": null,
      "max-vnic-attachments": 2,
      "memory-in-gbs": 6.0,
      "networking-bandwidth-in-gbps": 1.0,
      "ocpus": 1.0,
      "processor-description": "3.0 GHz Ampere\u00ae Altra\u2122",
      "vcpus": 1
    },
    "source-details": {
      "boot-volume-size-in-gbs": null,
      "boot-volume-vpus-per-gb": null,
      "image-id": "ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq",
      "instance-source-image-filter-details": null,
      "kms-key-id": null,
      "source-type": "image"
    },
    "system-tags": {},
    "time-created": "2023-08-22T02:54:36.813000+00:00",
    "time-maintenance-reboot-due": null
  },
  "etag": "679b0d39c271ef84698624d9c9a8e3489d4dd085317ec08bca3b2d8b77ac3ee4",
  "opc-work-request-id": "ocid1.coreservicesworkrequest.oc1.us-sanjose-1.abzwuljrfmf5mroubbpdtpd53mrx6dxedypc2fhs7zbik5zz37csqrhnvjta"
}
	`

	id := ParseInstanceId(jsonData)

	if id != "ocid1.instance.oc1.us-sanjose-1.anzwuljra36or2qcxhja7dkkuqfg5lfdzl45mb4bath4ac7yzdj4vcl5e2gq" {
		t.Error(id)
	}

	jsonData = `ip {
  "data": [
    {
      "availability-domain": "lynu:US-SANJOSE-1-AD-1",
      "compartment-id": "ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq",
      "defined-tags": {
        "Oracle-Tags": {
          "CreatedBy": "hq@schmied.us",
          "CreatedOn": "2023-08-22T03:19:38.242Z"
        }
      },
      "display-name": "instance20230822031939",
      "freeform-tags": {},
      "hostname-label": "instance20230822031939-670401",
      "id": "ocid1.vnic.oc1.us-sanjose-1.abzwuljruukaxjm5nqw75ogg4ve4rmgn5cahue5n5no2ihgoyslkd2s35k6q",
      "is-primary": true,
      "lifecycle-state": "AVAILABLE",
      "mac-address": "02:00:17:00:83:2F",
      "nsg-ids": [],
      "private-ip": "10.0.0.236",
      "public-ip": "146.235.211.85",
      "skip-source-dest-check": false,
      "subnet-id": "ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q",
      "time-created": "2023-08-22T03:19:42.809000+00:00",
      "vlan-id": null
    }
  ]
}`

	ip := ParsePublicIP(jsonData)

	if ip != "146.235.211.85" {
		t.Error(ip)
	}

}
