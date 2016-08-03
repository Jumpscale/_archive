FROM jumpscale/ubuntu1510

MAINTAINER JumpScale

EXPOSE 3000 8083 8086 8088

ADD . /bd_build

RUN /bd_build/prepare.sh && \
	/bd_build/system_services.sh && \
	/bd_build/utilities.sh && \
	/bd_build/cleanup.sh

ENV HOME /root/

VOLUME ["/data/influxdb", "/data/grafana", "/var/log/grafana", "/etc/grafana"]

CMD ["/sbin/my_init"]
