#!/bin/bash
set -e
source /bd_build/buildconfig
set -x


## Install mongo.
[ "$DISABLE_SYSLOG" -eq 0 ] && /bd_build/services/mongo/mongo.sh || true
