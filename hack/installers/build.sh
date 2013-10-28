#!/bin/bash
#
#   Initial concept for building release installers

RELEASE="$1"
DOCKER_ROOT=`pwd`
DOCKER_VERSION=`cat ./VERSION`
DOCKER_BUNDLE_PATH=$DOCKER_ROOT/bundles/$DOCKER_VERSION/binary/docker-$DOCKER_VERSION
RELEASE_DIR=$DOCKER_ROOT/releases
RELEASE_VERSION="release-$DOCKER_VERSION"
RELEASE_FILE_BASE=$RELEASE_DIR/$RELEASE_VERSION
BUILD_DIR=/tmp/$RELEASE_VERSION
BUILD_MIRROR=${BUILD_MIRROR:-http://archive.ubuntu.com/ubuntu}
mkdir -p $RELEASE_DIR

function create_docker_release {
    echo "Creating build for $DOCKER_VERSION..."
    docker run -lxc-conf=lxc.aa_profile=unconfined -privileged -v $DOCKER_ROOT:/go/src/github.com/dotcloud/docker docker hack/make.sh binary
}
function build_tarball {
    create_docker_release
    echo "Creating tarball for $DOCKER_VERSION..."
    #CORE_ID=$(docker run -lxc-conf=lxc.aa_profile=unconfined -privileged -v $DOCKER_ROOT:/go/src/github.com/dotcloud/docker docker hack/installers/build_core.sh)
    CORE_ID=$(docker run -i -t -d -v $DOCKER_ROOT:/docker ubuntu:12.04 /bin/bash -c /docker/hack/installers/build_core.sh)
    if [ "`docker wait $CORE_ID`" != "0" ] ; then
        echo "Error building.  Check logs for $CORE_ID"
        exit 1
    fi
    docker export $CORE_ID > $RELEASE_FILE_BASE.tar
}

function build_image {
    RELEASE_FILE=$RELEASE_FILE_BASE.img
    if [ ! -e $RELEASE_FILE_BASE.tar ] ; then
        build_tarball
    else
        echo "Using existing `basename $RELEASE_FILE_BASE.tar`"
    fi
    mkdir $BUILD_DIR
    pushd $BUILD_DIR > /dev/null  && tar xf $RELEASE_FILE_BASE.tar
    popd > /dev/null
    # get usage of base system for disk image
    SIZE=`du -ls $BUILD_DIR | awk '{ print $1; }'`
    # convert to MB and add 200M for spare
    NEW_SIZE=$(( ( ${SIZE#0} / 1024 ) + 200 ))
    # create and mount base image
    qemu-img create -f raw $RELEASE_FILE ${NEW_SIZE}M
    # partition
    LOOP_DEV=/dev/loop7
    losetup $LOOP_DEV $RELEASE_FILE
    parted --script $LOOP_DEV mklabel msdos
    parted --script $LOOP_DEV mkpart primary ext3 0 $NEW_SIZE
    parted --script $LOOP_DEV set 1 boot on
    # format
    mkfs.ext3 -L / -m 0 ${LOOP_DEV}p1
    # mount new partition and sync
    TMP_MNT='/tmp/build_root'
    mkdir -p $TMP_MNT
    mount ${LOOP_DEV}p1 $TMP_MNT
    rsync -a $BUILD_DIR/ $TMP_MNT/
    # bind mounts
    mount -o bind /dev $TMP_MNT/dev
    mount -o bind /proc $TMP_MNT/proc
    echo "root:docker" | chroot $TMP_MNT chpasswd
    # create /etc/fstab
    cat << EOF > $TMP_MNT/etc/fstab
LABEL=/ / ext3 defaults 0 0
EOF
    # create /etc/network/interfaces
    cat << EOF > $TMP_MNT/etc/network/interfaces
auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp
EOF
    # create /etc/resolv.conf
    cat << EOF > $TMP_MNT/etc/resolv.conf
nameserver 8.8.8.8
EOF
    # create /etc/hostname
    cat << EOF > $TMP_MNT/etc/hostname
docker
EOF
    KERNEL=`ls $TMP_MNT/boot | grep vmlinuz*`
    RAMDISK=`ls $TMP_MNT/boot | grep initrd.*`
    # create /boot/grub/menu.lst
    cat << EOF > $TMP_MNT/boot/grub/menu.lst
default 0
timeout 1
title=Docker
root (hd0,0)
kernel /boot/$KERNEL root=LABEL=.
initrd /boot/$RAMDISK
EOF
    chroot $TMP_MNT grub-install --recheck /dev/sda
    sleep 3 # wait for procs to exit
    fuser -k $TMP_MNT
    umount $TMP_MNT/dev
    umount $TMP_MNT/proc
    umount -f $TMP_MNT
    losetup -d $LOOP_DEV
    rm -rf $BUILD_DIR
    rm -rf $TMP_MNT
}

case "$RELEASE" in
	*tarball)
            build_tarball
            ;;
        *image)
            build_image
            ;;
	*ec2)
	    ;;
	*)
    	    echo >&2 "Usage: $0 <tarball|ec2>"
	    exit 1
	    ;;
esac
