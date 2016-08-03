#!/bin/bash
set -e
source /bd_build/buildconfig
set -x

cd /tmp
curl -L https://git.aydo.com/binary/agentcontroller2/repository/archive.tar.gz -o agentcontroller2.tar.gz
tar xzvf agentcontroller2.tar.gz
cd agentcontroller2.git
mv agentcontroller2 /usr/bin/
rm /usr/bin/agentcontroller2/agentcontroller2.toml
cd ..
rm -rf agentcontroller2.git
rm -f agentcontroller2.tar.gz
cp /bd_build/services/agentcontroller/agentcontroller.toml /usr/bin/agentcontroller2/agentcontroller.toml

mkdir /etc/service/agentcontroller
cp /bd_build/services/agentcontroller/agentcontroller.runit /etc/service/agentcontroller/run
