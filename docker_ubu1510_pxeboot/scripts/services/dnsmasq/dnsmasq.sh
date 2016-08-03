#!/bin/bash
set -e
source /bd_build/buildconfig
set -x

BUILD_PATH=/bd_build/services/dnsmasq

## Install the SSH server.
$minimal_apt_get_install dnsmasq
mkdir /var/run/dnsmasq
mkdir /etc/service/dnsmasq
touch /etc/service/dnsmasq/down

cp $BUILD_PATH/dnsmasq.runit /etc/service/dnsmasq/run

