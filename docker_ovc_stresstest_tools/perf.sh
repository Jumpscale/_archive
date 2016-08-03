#!/bin/bash

set -e
source /bd_build/buildconfig
set -x

echo 'deb http://archive.ubuntu.com/ubuntu wily multiverse' > /etc/apt/sources.list.d/multiverse.list
echo 'deb http://archive.ubuntu.com/ubuntu wily-updates multiverse' >> /etc/apt/sources.list.d/multiverse.list
apt-get update

$minimal_apt_get_install iozone3 btrfs-tools xfsprogs

pip install statsd
