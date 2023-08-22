#!/bin/bash

# Licensed under Creative Commons CC0.
#
# To the extent possible under law, the author(s) have dedicated all copyright
# neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.
# You should have received a copy of the CC0 Public Domain Dedication along wi
# If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

certbot --standalone -m hq@schmied.us --cert-name oracle.eper.io -d oracle.eper.io,sunshine.oracle.eper.io,whisper.oracle.eper.io,ocean.oracle.eper.io,jazz.oracle.eper.io,dragon.oracle.eper.io,symphony.oracle.eper.io,telescope.oracle.eper.io,lantern.oracle.eper.io,harmony.oracle.eper.io,blossom.oracle.eper.io,firefly.oracle.eper.io,apple.oracle.eper.io,galaxy.oracle.eper.io,bicycle.oracle.eper.io,mountain.oracle.eper.io,rainbow.oracle.eper.io,elephant.oracle.eper.io,serenade.oracle.eper.io,enigma.oracle.eper.io,adventure.oracle.eper.io,moonlight.oracle.eper.io,wanderlust.oracle.eper.io,chocolate.oracle.eper.io,carousel.oracle.eper.io --https-port 4443 certonly
go run main.go dns
shutdown -r now
