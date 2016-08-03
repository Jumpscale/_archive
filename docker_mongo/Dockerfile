FROM jumpscale/ubuntu1510

MAINTAINER JumpScale

EXPOSE 27017

ADD . /bd_build

RUN /bd_build/prepare.sh && \
	/bd_build/system_services.sh && \
	/bd_build/utilities.sh && \
	/bd_build/cleanup.sh

ENV HOME /root/


VOLUME /data/db


CMD ["/sbin/my_init"]
