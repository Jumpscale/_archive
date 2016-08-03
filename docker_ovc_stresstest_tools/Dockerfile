FROM jumpscale/ubuntu1510_python

MAINTAINER JumpScale

ADD . /bd_build

ENV HOME /root/

RUN /bd_build/perf.sh
 
CMD ["/sbin/my_init"]
