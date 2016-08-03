#!/bin/bash
set -e
source /bd_build/buildconfig
set -x



## Install influxdb.
[ "$DISABLE_SYSLOG" -eq 0 ] && /bd_build/services/influxdb/influxdb.sh || true

## Install grafana.
[ "$DISABLE_SYSLOG" -eq 0 ] && /bd_build/services/grafana/grafana.sh || true
