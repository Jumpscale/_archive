set -e
source /bd_build/buildconfig
set -x


curl -L https://git.aydo.com/binary/skydns/raw/master/skydns -o /usr/bin/skydns
chmod a+x /usr/bin/skydns


mkdir /etc/service/skydns
cp /bd_build/services/skydns/skydns.runit /etc/service/skydns/run
