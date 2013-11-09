# DRAFT

# Docker Official Installers

This project is to create the "official" installers that will be used to create the various Docker releases and make it easier to run Docker on multiple platforms.  These include EC2, Qemu (Mac/Windows), virtual disk for Virtualbox/VMware and eventually other providers such as Digital Ocean, Google Compute, etc.

# Preface
Initial requirements:

* Use Docker to build
* Pre-bake images (no launch instance & build)
* Images must be small
* Installs ideally will be upgradeable

# Phase 1
For the initial release the Docker images are built using Ubuntu 12.04.  There
is a base Docker image (currently `ehazlett/docker-image-base`) that was created
as follows:

* `vmbuilder kvm ubuntu --suite precise --arch amd64 --flavour virtual --add ca-certificates --chroot-dir /tmp/chroot --only-chroot`
* `tar -C /tmp/chroot -c . | docker import - docker-base-image`

You can change `docker-base-image` to a repository and push to the index as well.

This base image is then used in `hack/make/disk_image` as the core.  A 
`docker export` exports that root filesystem and The .deb files created from 
the `binary` make target are then installed into the chroot and packaged as
disk images.

# Phase 2 (planning)
The base image will be based on Tiny Core Linux a very minimal distro.  This 
will drastically reduce the image size and make for easier distribution and 
faster boot times.
