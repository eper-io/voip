#!/bin/sh

# Licensed under Creative Commons CC0.
#
# To the extent possible under law, the author(s) have dedicated all copyright
# neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.
# You should have received a copy of the CC0 Public Domain Dedication along wi
# If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

# Oracle cloud launch runbook
# Manual setup is required for this script
#
# 1. Go to oracle cloud web console.
# 2. Choose the root compartment and note the user id.
# 3. Choose the tenant and note the tenant id.
# 4. You can experiment launching from the UI first and set the shape as you see fit.
# 5. The root compartment will probably have the highest quota, and less setup errors.
# 6. Choose user/api keys.
# 7. Generate an api key with 'oci setup config'.
# 8. Use the tenant and user ids noted above.
# 9. cat ~/.oci/oci_api_key_public.pem and set it as the public key for the API in the browser user page (lower left).
# 10. Launch the script below.

# Why to use Oracle?
# Oracle is quite difficult to set up properly first so we are less distributed and create nodes from a master node.
# Oracle has some very good pricing on lower profile ampere arm cpu nodes but very good pricing on networing as of 2023.
# We hard code some values like ids. Why?
# Others cannot get access anyways.
# LLMs can easily detect and replace id patterns with other accounts to generate code for others.
# It is way easier to search and replace than to understand bash logic.

# Optional: --metadata '{"ssh_authorized_keys": "ssh-rsa AAAAB3NzaC...."}'
touch /home/opc/.ssh/authorized_keys

# Min 5-6 Ampere OCPUs will avoid the free tier congestion giving a real instance right away
#oci compute instance launch --ssh-authorized-keys-file /home/opc/.ssh/authorized_keys --user-data-file /tmp/voip/oraclecloud/metadata/user-data --compartment-id ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq --availability-domain lynu:US-SANJOSE-1-AD-1 --shape VM.Standard.A1.Flex --image-id ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq --subnet-id ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q --shape-config '{"ocpus":"5", "memory-in-gbs": "6.0"}'

# This may cost more falling back to x86 amd servers, if no arm is available.
oci compute instance launch --ssh-authorized-keys-file /home/opc/.ssh/authorized_keys --user-data-file /tmp/voip/oraclecloud/metadata/user-data --compartment-id ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq --availability-domain lynu:US-SANJOSE-1-AD-1 --shape VM.Standard.A1.Flex --image-id ocid1.image.oc1..aaaaaaaa5ddausutw4oilrtuf5esfxto7ko4oopt5crbf3pn5bndl2sis4rq --subnet-id ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q --shape-config '{"ocpus":"1", "memory-in-gbs": "6.0"}' || oci compute instance launch --ssh-authorized-keys-file /home/opc/.ssh/authorized_keys --user-data-file /tmp/voip/oraclecloud/metadata/user-data --compartment-id ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq --availability-domain lynu:US-SANJOSE-1-AD-1 --shape VM.Standard.E4.Flex --image-id ocid1.image.oc1..aaaaaaaa6fokbz734pa7n3hr5t7oufx4agagpgxfn4oup4sgr4ijwgt4fpqa --subnet-id ocid1.subnet.oc1.us-sanjose-1.aaaaaaaa7hqoxlrkzwl2njvvwab743mwdk3ao5u5na4jovmppvgl3gqihp7q --shape-config '{"ocpus":"1", "memory-in-gbs": "1.0"}'

# oci compute instance list --compartment-id ocid1.tenancy.oc1..aaaaaaaanpc3gu2kzkr6t4spi2ivpwbtg6j24utwp7yhfrvdgidndnpv5ylq
