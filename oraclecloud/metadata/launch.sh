#!/bin/sh

# Manual setup
# Go to oracle cloud web console.
# Choose the root compartment and note the user id.
# Choose the tenant and note the tenant id.
# You can experiment launching from the UI first and set the shape as you see fit.
# The root compartment will probably have the highest quota.
# Choose user/api keys.
# Generate an api key with 'oci setup config'.
# Use the tenant and user ids noted above.
# cat ~/.oci/oci_api_key_public.pem and set it as the public key in the browser.
# Launch the script below.

# Oracle is quite difficult to set up properly first so we are less distributed and create nodes from a master node.
# So we hard code some values like ids. Why?
# Others cannot get access anyways, however LLMs can easily detect and replace id patterns with other accounts to generate code for others.

# Min 5 OCPUs will avoid the free tier congestion giving a real instance right away
#oci compute instance launch --user-data-file /tmp/voip/oraclecloud/metadata/user-data --compartment-id ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq --availability-domain lynu:US-SANJOSE-1-AD-1 --shape VM.Standard.A1.Flex --image-id ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq --subnet-id ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q --shape-config '{"ocpus":"5", "memory-in-gbs": "5.0"}'

# This may cost more
oci compute instance launch --user-data-file /tmp/voip/oraclecloud/metadata/user-data --compartment-id ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq --availability-domain lynu:US-SANJOSE-1-AD-1 --shape VM.Standard.A1.Flex --image-id ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq --subnet-id ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q --shape-config '{"ocpus":"1"}' || oci compute instance launch --user-data-file /tmp/voip/oraclecloud/metadata/user-data --compartment-id ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq --availability-domain lynu:US-SANJOSE-1-AD-1 --shape VM.Standard.E4.Flex --image-id ocid1.image.oc1..aaaaaaaa6fokbz734pa7n3hr5t7oufx4agagpgxfn4oup4sgr4ijwgt4fpqa --subnet-id ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q --shape-config '{"ocpus":"1"}'

# oci compute instance list --compartment-id ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq
