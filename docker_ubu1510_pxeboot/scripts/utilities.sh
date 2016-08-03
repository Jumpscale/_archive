#!/bin/bash
set -e
source /bd_build/buildconfig
set -x


$minimal_apt_get_install curl less mc python3.5 iproute2 iputils-arping inetutils-telnet inetutils-ftp rsync inetutils-traceroute iputils-ping iputils-tracepath iputils-clockdiff

$minimal_apt_get_install net-tools sudo

$minimal_apt_get_install mc git wget tmux

rm -rf /usr/bin/python
ln /usr/bin/python3.5 /usr/bin/python

## This tool runs a command as another user and sets $HOME.
cp /bd_build/bin/setuser /sbin/setuser
