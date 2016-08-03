FROM jumpscale/ubuntu1510

MAINTAINER JumpScale

# etcd ports
EXPOSE 2379 2380

# skydns ports
EXPOSE 53

# Agentcontroller ports
EXPOSE 8966

ADD . /bd_build

RUN /bd_build/prepare.sh && \
	/bd_build/system_services.sh && \
	/bd_build/utilities.sh && \
	/bd_build/cleanup.sh

ENV HOME /root/


VOLUME /data/db


CMD ["/sbin/my_init"]
