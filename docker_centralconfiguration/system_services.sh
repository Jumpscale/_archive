#!/bin/bash
set -e
source /bd_build/buildconfig
set -x


## Install etcd.
[ "$DISABLE_SYSLOG" -eq 0 ] && /bd_build/services/etcd/etcd.sh || true

## Install skydns.
[ "$DISABLE_SYSLOG" -eq 0 ] && /bd_build/services/skydns/skydns.sh || true

## Install redis (needed for the agentcontroller)
[ "$DISABLE_SYSLOG" -eq 0 ] && /bd_build/services/redis/redis.sh || true

## Install agentcontroller
[ "$DISABLE_SYSLOG" -eq 0 ] && /bd_build/services/agentcontroller/agentcontroller.sh || true
