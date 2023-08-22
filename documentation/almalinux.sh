#!/bin/bash

# To the extent possible under law, the author(s) have dedicated all copyright
# neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.
# You should have received a copy of the CC0 Public Domain Dedication along wi
# If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

# Installation:
# 1. Run install script
# curl https://gitlab.com/eper.io/voip/-/raw/main/documentation/almalinux.sh?ref_type=heads >/var/lib/voip.sh
# Review code
# bash /var/lib/voip.sh

# 2. Enable the service
# One time setup a CI/CD pipeline from git
# Add this to 'crontab -e'
# @reboot sudo /tmp/voip/documentation/bashrc.sh

# 3. Fine tune
# Replace all occurrences of l.eper.io in this file to your domain.

# 4. Reboot
# shutdown -r now

# Make sure we are running as root
touch /root/voip || (echo *** run this as root ***; exit 1;)

mkdir /var/lib/voip

yum install -y docker
yum install -y nginx
yum install -y git
yum install -y golang

# Professor I. Only open the ports that are needed
# Professor II. You must make sure that we do not have backdoors even if all the ports are open
# Draw 1:1
systemctl disable firewalld
systemctl stop firewalld

# // TODO
# Certbot requires this
dnf install -y epel-release
dnf -y upgrade
dnf install -y snapd
systemctl enable --now snapd.socket
ln -s /var/lib/snapd/snap /snap
snap install core; sudo snap refresh core

# Certbot is good enough to make the browsers happy
# Enterprise projects must use corporate signed certificates
snap install --classic certbot
ln -s /snap/bin/certbot /usr/bin/certbot

if ! [ -d /etc/letsencrypt/live/l.eper.io/ ]; then
  echo *** SETUP NEEDED *** Replace all occurrences of l.eper.io in this file to your domain.
  certbot --nginx -m hq@schmied.us --cert-name l.eper.io -d l.eper.io --https-port 4443 --http-01-port 4444 certonly
  cd /tmp
  git clone https://gitlab.com/eper.io/voip.git
fi

cat <<EOF >/tmp/voip/metadata/data.go
package metadata

var ActivationKey = "QTOPZNNEYGPBKUQEKJYLHBIJVHAJEOOXWMFMXWJDDWNOIJVHXFXRLFJXAAOGEBRBUMQJEYYNDHBTWJUYVNEKZJMJTHHR"
var SiteUrl = "https://l.eper.io"
var Certificate = "/etc/letsencrypt/live/l.eper.io/fullchain.pem"
var PrivateKey = "/etc/letsencrypt/live/l.eper.io/privkey.pem"
var ContainerRuntime = "line.eper.io/line"
var Info = ""
var Bandwidth = ""
var RandomSalt = "XBGXTNTKIAVWBNHGODJGSSNUFBDISPRYVKCFLYBFHPEWBRHQHYUWQLHHOPZLDZREJIAVPGEQMHOJFICSXNWADFHIHFRR"

EOF

cat <<EOF >/tmp/voip/Dockerfile
FROM golang:1.19.3

ADD . /go/src

WORKDIR /go/src

# This will listen to tcp port metadata.Http11Port externally.
CMD go run main.go
EOF

# We just use a different port than 80 rather than disabling
#systemctl disable nginx.service
#systemctl stop nginx.service

cp -f /etc/nginx/nginx.conf /var/lib/voip/nginx.conf.$RANDOM
cat <<EOF >/etc/nginx/nginx.conf
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;

# Load dynamic modules. See /usr/share/doc/nginx/README.dynamic.
include /usr/share/nginx/modules/*.conf;

events {
    worker_connections 1024;
}

http {
    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile            on;
    tcp_nopush          on;
    tcp_nodelay         on;
    keepalive_timeout   65;
    types_hash_max_size 4096;

    include             /etc/nginx/mime.types;
    default_type        application/octet-stream;

    # Load modular configuration files from the /etc/nginx/conf.d directory.
    # See http://nginx.org/en/docs/ngx_core_module.html#include
    # for more information.
    include /etc/nginx/conf.d/*.conf;

    server {
        listen       4442;
        listen       [::]:4442;
        server_name  _;
        root         /usr/share/nginx/html;

        # Load configuration files for the default server block.
        include /etc/nginx/default.d/*.conf;

        error_page 404 /404.html;
        location = /404.html {
        }

        error_page 500 502 503 504 /50x.html;
        location = /50x.html {
        }
    }
}
EOF

cd /tmp/voip/
git config --global user.email "hq@example.com"
git config --global user.name "voip installer"

git add .
git commit -m local

touch /etc/containers/nodocker
