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

cat <<EOF >/etc/nginx/nginx.conf
# For more information on configuration, see:
#   * Official English Documentation: http://nginx.org/en/docs/
#   * Official Russian Documentation: http://nginx.org/ru/docs/

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

# Settings for a TLS enabled server.
#
#    server {
#        listen       443 ssl http2;
#        listen       [::]:443 ssl http2;
#        server_name  _;
#        root         /usr/share/nginx/html;
#
#        ssl_certificate "/etc/pki/nginx/server.crt";
#        ssl_certificate_key "/etc/pki/nginx/private/server.key";
#        ssl_session_cache shared:SSL:1m;
#        ssl_session_timeout  10m;
#        ssl_ciphers PROFILE=SYSTEM;
#        ssl_prefer_server_ciphers on;
#
#        # Load configuration files for the default server block.
#        include /etc/nginx/default.d/*.conf;
#
#        error_page 404 /404.html;
#            location = /40x.html {
#        }
#
#        error_page 500 502 503 504 /50x.html;
#            location = /50x.html {
#        }
#    }

}
EOF

cd /tmp/voip/
git config --global user.email "hq@schmied.us"
git config --global user.name "eper.io installer"

git add .
git commit -m local

touch /etc/containers/nodocker
