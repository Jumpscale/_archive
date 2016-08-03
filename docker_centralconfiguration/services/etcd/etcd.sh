#!/bin/bash
set -e
source /bd_build/buildconfig
set -x

cd /tmp
curl -L  https://github.com/coreos/etcd/releases/download/v2.2.2/etcd-v2.2.2-linux-amd64.tar.gz -o etcd-v2.2.2-linux-amd64.tar.gz
tar xzvf etcd-v2.2.2-linux-amd64.tar.gz
cd etcd-v2.2.2-linux-amd64
mv etcd /usr/bin/
mv etcdctl /usr/bin/
cd ..
rm -rf etcd-v2.2.2-linux-amd64
rm -f etcd-v2.2.2-linux-amd64.tar.gz

mkdir /etc/service/etcd
cp /bd_build/services/etcd/etcd.runit /etc/service/etcd/run
