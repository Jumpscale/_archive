#!/bin/bash
set -e
source /bd_build/buildconfig
set -x

groupadd -r influxdb && useradd -r -g influxdb influxdb

wget https://s3.amazonaws.com/influxdb/influxdb_0.9.4.2_amd64.deb
dpkg -i influxdb_0.9.4.2_amd64.deb

mkdir -p /data/influxdb && chown -R influxdb:influxdb /data/influxdb


mkdir /etc/service/influxdb
cp /bd_build/services/influxdb/influxdb.runit /etc/service/influxdb/run

mkdir /etc/influxdb
cp /bd_build/services/influxdb/influxdb.conf /etc/influxdb/
