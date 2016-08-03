#!/bin/bash
set -e
source /bd_build/buildconfig
set -x

$minimal_apt_get_install libfontconfig wget adduser openssl ca-certificates

wget https://grafanarel.s3.amazonaws.com/builds/grafana_2.5.0_amd64.deb -O /tmp/grafana.deb && \
    dpkg -i /tmp/grafana.deb && \
    rm /tmp/grafana.deb

mkdir -p /data/grafana && chown -R grafana:grafana /data/grafana


mkdir /etc/service/grafana
cp /bd_build/services/grafana/grafana.runit /etc/service/grafana/run
