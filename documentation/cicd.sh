#!/bin/bash

# To the extent possible under law, the author(s) have dedicated all copyright
# neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.
# You should have received a copy of the CC0 Public Domain Dedication along wi
# If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

# Usage
# See almalinux.sh for Almalinux

# Example: Find if a changelist is included:
# https://l.eper.io/englang#:~:text=security
# 26881f8f71cf2e7e1ec37c1552cfeb9142be86fa security hardening

# Short glitch in the service only when updated. Acceptable.
cd /tmp/voip
git pull -r > /var/log/voip

# Save some logs
TZ='America/Los_Angeles' date >> /var/log/voip
echo Next update check is in thirty seconds >> /var/log/voip
git status >> /var/log/voip
git log --format=oneline >> /var/log/voip

# Build voip broker
(cat /var/log/voip | grep 'up to date') || docker build -t line.eper.io/line /tmp/voip
echo "build result $?" >> /var/log/voip

# We may want to run a privileged container in the future. It is difficult to mix podman and golang
cd /tmp/voip
(cat /var/log/voip | grep 'up to date') || go build -o /opt/voipbroker ./eos/main/main.go >> /var/log/voip || true
cat /var/log/voip

# All the worker containers keep running for their respective customers until they shut down themselves. (~1-2 hours)

# Run test container
# docker stop -t 2 voiptest
# sleep 6
# docker run -d --rm --name voiptest -e SITEURL=https://l.eper.io -e PORT=7777 -p 7777:443 -v /etc/letsencrypt/live/l.eper.io/fullchain.pem:/tmp/fullchain.pem:ro -v /etc/letsencrypt/live/l.eper.io/privkey.pem:/tmp/privkey.pem:ro line.eper.io/line

# Running as root on the network is dangerous so we need to be very lean with the codebase in /eos

# Run the broker if needed. It launches containers that do the call lines
#pgrep voipbroker || (DOCKERIMAGE=line.eper.io/line SITEURL=https://l.eper.io APIKEY=JVPSVWUIUTSXGPTWOVEWMHBUFJMVIALPQDMXQZROKZLYPYQGMBRQZMRWSQZIACQDKIFVWYQBWGGHQLGALYBQTAQNLHDR nohup /opt/voipbroker >>/var/log/voipbroker &)
(cat /var/log/voip | grep 'up to date') || kill -9 `pgrep voipbroker`
sleep 2
pgrep voipbroker || (nohup /opt/voipbroker no-proxy >>/var/log/voipbroker &) || true

sleep 30;
