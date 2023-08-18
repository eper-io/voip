#!/bin/sh

# use oci setup config
# tenant ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq

oci compute instance launch --user-data-file /tmp/voip/oraclecloud/metadata/user-data --compartment-id ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq --availability-domain lynu:US-SANJOSE-1-AD-1 --shape VM.Standard.A1.Flex --image-id ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq --subnet-id ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q --shape-config '{"ocpus":"4"}'

