#!/bin/bash
#
#   Initial concept for building release installers

RELEASE="$1"
#DOCKER_ROOT=/go/src/github.com/dotcloud/docker
DOCKER_ROOT=/docker
DOCKER_VERSION=`cat $DOCKER_ROOT/VERSION`
DOCKER_BUNDLE_PATH=$DOCKER_ROOT/bundles/$DOCKER_VERSION/binary/docker-$DOCKER_VERSION
RELEASE_DIR=$DOCKER_ROOT/releases
RELEASE_VERSION="release-$DOCKER_VERSION"
# create release dir
mkdir -p $RELEASE_DIR
DISTRO=`cat /etc/lsb-release | grep DISTRIB_CODENAME | awk -F "=" '{ print $2; }'`
# install packages
cat << EOF > /etc/apt/sources.list
deb http://us.archive.ubuntu.com/ubuntu $DISTRO main universe multiverse
deb-src http://us.archive.ubuntu.com/ubuntu $DISTRO main universe multiverse
deb http://us.archive.ubuntu.com/ubuntu $DISTRO-updates main universe multiverse
deb-src http://us.archive.ubuntu.com/ubuntu $DISTRO-updates main universe multiverse
deb http://us.archive.ubuntu.com/ubuntu $DISTRO-backports main universe multiverse
deb-src http://us.archive.ubuntu.com/ubuntu $DISTRO-backports main universe multiverse
EOF
apt-get update
RUNLEVEL=1 DEBIAN_FRONTEND=noninteractive apt-get install -y linux-image-generic-lts-raring linux-headers-generic-lts-raring aufs-tools dhcpcd lxc

# add docker binary
cp $DOCKER_BUNDLE_PATH /usr/local/bin/docker

echo "$DOCKER_VERSION" >> /etc/docker-version

