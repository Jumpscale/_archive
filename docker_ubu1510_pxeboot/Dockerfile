FROM ubuntu:15.10

MAINTAINER JumpScale

# EXPOSE 9000 8080 3389

ADD scripts /bd_build/

RUN /bd_build/prepare.sh && \
	/bd_build/cleanup.sh 

COPY files/dnsmasq.conf /etc





ENTRYPOINT ["/usr/sbin/dnsmasq"]
