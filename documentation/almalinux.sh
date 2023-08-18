#!/bin/bash

# To the extent possible under law, the author(s) have dedicated all copyright
# neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.
# You should have received a copy of the CC0 Public Domain Dedication along wi
# If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

# Installation:
# curl https://gitlab.com/eper.io/voip/-/raw/main/documentation/almalinux.sh?ref_type=heads >/tmp/setup.sh
# Review code
# bash /tmp/setup.sh

# Once set up a CI/CD from git
# Add this to 'crontab -e'
# @reboot sudo /tmp/voip/documentation/bashrc.sh

# Install moose-assist on AlmaLinux on Oracle Cloud
touch /root/moose-assist || (echo *** run this as root ***; exit 1;)

yum install -y docker
yum install -y nginx
yum install -y git
yum install -y golang
systemctl disable nginx.service
systemctl stop nginx.service
systemctl disable firewalld
systemctl stop firewalld

dnf install -y epel-release
dnf -y upgrade
dnf install -y snapd
systemctl enable --now snapd.socket
ln -s /var/lib/snapd/snap /snap
snap install core; sudo snap refresh core

snap install --classic certbot
ln -s /snap/bin/certbot /usr/bin/certbot

if ! [ -d /etc/letsencrypt/live/moose-assist.eper.io/ ]; then
  echo *** SETUP NEEDED *** Replace all occurrences of moose-assist.eper.io in this file to your domain.
  certbot --nginx -m hq@schmied.us --cert-name l.eper.io -d l.eper.io --https-port 4443 certonly
  cd /tmp
  git clone https://gitlab.com/eper.io/voip.git
fi

cat <<EOF >/tmp/voip/metadata/data.go
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
var ActivationKey = "QTOPZNNEYGPBKUQEKJYLHBIJVHAJEOOXWMFMXWJDDWNOIJVHXFXRLFJXAAOGEBRBUMQJEYYNDHBTWJUYVNEKZJMJTHHR"
var SiteUrl = "https://l.eper.io"
var Certificate = "/etc/letsencrypt/live/l.eper.io/fullchain.pem"
var PrivateKey = "/etc/letsencrypt/live/l.eper.io/privkey.pem"
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
var RandomSalt = "XBGXTNTKIAVWBNHGODJGSSNUFBDISPRYVKCFLYBFHPEWBRHQHYUWQLHHOPZLDZREJIAVPGEQMHOJFICSXNWADFHIHFRR"

EOF

cat <<EOF >/tmp/voip/Dockerfile
FROM golang:1.19.3
ADD . /go/src
WORKDIR /go/src
RUN apt update; apt install -y docker-compose;
# This will listen to tcp port metadata.Http11Port externally.
CMD go run main.go
EOF

cd /tmp/voip/
git config --global user.email "hq@schmied.us"
git config --global user.name "eper.io installer"

git add .
git commit -m local

touch /etc/containers/nodocker
