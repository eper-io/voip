#!/bin/bash

# To the extent possible under law, the author(s) have dedicated all copyright
# neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.
# You should have received a copy of the CC0 Public Domain Dedication along wi
# If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

# Usage
# Check out Moose Assist
# git clone https://gitlab.com/eper.io/voip.git /tmp/voip
# Run this command from /tmp/voip:
cd /tmp/voip

date

git status
# Short glitch in the service. Acceptable.
(git pull -r | grep 'up to date') || kill -9 `pgrep voipbroker`

docker build -t line.eper.io/line /tmp/voip

# We may want to run a privileged container in the future. It is difficult to mix podman and golang
cd /tmp/voip
go build -o /opt/voipbroker ./eos/main/main.go

# All the worker containers keep running for their respective customers until they shut down themselves. (~2 hours)

# Run test container

docker stop -t 2 voiptest

sleep 6

# Running docker as root on the network is dangerous so we need to be very lean and careful with the codebase in /eos

docker run -d --rm --name voiptest -e SITEURL=https://l.eper.io -e PORT=7777 -p 7777:443 -v /etc/letsencrypt/live/moose-assist.eper.io/fullchain.pem:/tmp/fullchain.pem:ro -v /etc/letsencrypt/live/moose-assist.eper.io/privkey.pem:/tmp/privkey.pem:ro line.eper.io/line

# Run the broker if needed.

pgrep voipbroker || (DOCKERIMAGE=line.eper.io/line SITEURL=https://l.eper.io APIKEY=JVPSVWUIUTSXGPTWOVEWMHBUFJMVIALPQDMXQZROKZLYPYQGMBRQZMRWSQZIACQDKIFVWYQBWGGHQLGALYBQTAQNLHDR /opt/voipbroker >>/var/log/voipbroker)

# WARNING blocking call. Do not extend here.