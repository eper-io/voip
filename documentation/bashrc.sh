#!/bin/bash

# To the extent possible under law, the author(s) have dedicated all copyright
# neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.
# You should have received a copy of the CC0 Public Domain Dedication along wi
# If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

# Periodic service script to update the code for staging
# Setup
# Add this to 'crontab -e'
# @reboot sudo /tmp/voip/documentation/bashrc.sh

systemctl disable nginx.service
systemctl stop nginx.service

# This can do an entire cleanup on boot
# rm /tmp/voipautoscale; cat /var/lib/voip_gc_* | sed --expression='s/sleep/echo/g' | bash
# rm -f /var/lib/voip_*

nohup bash -c 'while true; do timeout 600 /tmp/voip/documentation/cicd.sh; done' &

# Production workloads will want to snapshot the staging once it is proven and replicate it.