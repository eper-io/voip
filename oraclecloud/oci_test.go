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
        "CreatedOn": "2023-08-18T17:18:09.115Z"
      }
    },
    "display-name": "instance20230818171810",
    "extended-metadata": {},
    "fault-domain": "FAULT-DOMAIN-1",
    "freeform-tags": {},
    "id": "ocid1.instance.oc1.us-sanjose-1.anzwuljra36or2qc7kijp2tnx42qaj6kgvyjyjvwp2wgdwjddf3u456t5qoq",
    "image-id": "ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq",
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
      "user_data": "IyEvYmluL3NoCgplY2hvICJIZWxsbyBXb3JsZC4gIFRoZSB0aW1lIGlzIG5vdyAkKGRhdGUgLVIpISIgfCB0ZWUgL3Jvb3Qvb3V0cHV0LnR4dAo="
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
      "max-vnic-attachments": 4,
      "memory-in-gbs": 24.0,
      "networking-bandwidth-in-gbps": 4.0,
      "ocpus": 4.0,
      "processor-description": "3.0 GHz Ampere\u00ae Altra\u2122"
    },
    "source-details": {
      "boot-volume-size-in-gbs": null,
      "boot-volume-vpus-per-gb": null,
      "image-id": "ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq",
      "kms-key-id": null,
      "source-type": "image"
    },
    "system-tags": {},
    "time-created": "2023-08-18T17:18:10.101000+00:00",
    "time-maintenance-reboot-due": null
  },
  "etag": "deafd0c7107399ab7e855d9c9347b9897411edc5a86c27d803cd5413e6444dff",
  "opc-work-request-id": "ocid1.coreservicesworkrequest.oc1.us-sanjose-1.abzwuljromhqgbuz2vw4rrx6gogayvpcdudm3lsjo4bxx3t54n2xnmomyrba"
}

	`

	id := ParseInstanceId(jsonData)

	if id != "ocid1.instance.oc1.us-sanjose-1.anzwuljra36or2qc7kijp2tnx42qaj6kgvyjyjvwp2wgdwjddf3u456t5qoq" {
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
          "CreatedBy": "hq@schmied.us",
          "CreatedOn": "2023-08-20T02:20:34.369Z"
        }
      },
      "display-name": "instance20230820022035",
      "freeform-tags": {},
      "hostname-label": "instance20230820022035",
      "id": "ocid1.vnic.oc1.us-sanjose-1.abzwuljr5jg7eoyy2lmm5vikv5edjft5wmsz7dvi4ughfcyazvkjho2nbmqa",
      "is-primary": true,
      "lifecycle-state": "AVAILABLE",
      "mac-address": "02:00:17:00:D8:AD",
      "nsg-ids": [],
      "private-ip": "10.0.0.161",
      "public-ip": "138.2.237.100",
      "skip-source-dest-check": false,
      "subnet-id": "ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q",
      "time-created": "2023-08-20T02:20:41.078000+00:00",
      "vlan-id": null
    }
  ]
}

	`

	ip := ParsePublicIP(jsonData)

	if ip != "138.2.237.100" {
		t.Error(ip)
	}

}
