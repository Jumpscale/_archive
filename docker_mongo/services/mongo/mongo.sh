#!/bin/bash
set -e
source /bd_build/buildconfig
set -x

groupadd -r mongodb && useradd -r -g mongodb mongodb

apt-key adv --keyserver ha.pool.sks-keyservers.net --recv-keys "DFFA3DCF326E302C4787673A01C4E7FAAAB2461C"
apt-key adv --keyserver ha.pool.sks-keyservers.net --recv-keys "42F3E95A2C4F08279C4960ADD68FA50FEA312927"


MONGO_MAJOR=3.2

echo "deb http://repo.mongodb.org/apt/ubuntu trusty/mongodb-org/$MONGO_MAJOR multiverse" | tee /etc/apt/sources.list.d/mongodb-org.list
set -x \
	&& apt-get update \
	&& $minimal_apt_get_install \
		mongodb-org \
		mongodb-org-server \
		mongodb-org-shell \
		mongodb-org-mongos \
		mongodb-org-tools \
	&& rm -rf /var/lib/apt/lists/* \
	&& rm -rf /var/lib/mongodb \
	&& mv /etc/mongod.conf /etc/mongod.conf.orig

mkdir -p /data/db && chown -R mongodb:mongodb /data/db


mkdir /etc/service/mongo
cp /bd_build/services/mongo/mongo.runit /etc/service/mongo/run
