package daemon

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/volume"
)

func secretsMountPath(c *container.Container) string {
	return filepath.Join(c.Root, "secrets")
}

func (daemon *Daemon) setupSecrets(c *container.Container) error {
	secretsPath := secretsMountPath(c)
	// TODO: check for existing mount and unmount
	if err := os.MkdirAll(secretsPath, 0700); err != nil {
		return fmt.Errorf("secrets: unable to create ramfs mountpoint: %s", err)
	}
	logrus.Debugf("secrets: setting up secret ramfs at %s", secretsPath)
	if err := mount.Mount("ramfs", secretsPath, "ramfs", "nodev"); err != nil {
		return fmt.Errorf("secrets: unable to setup ramfs mount: %s", err)
	}

	// TODO: inject data
	rdr, err := daemon.getSecrets(c)
	if err != nil {
		return err
	}

	if err := archive.Untar(rdr, secretsPath, nil); err != nil {
		return err
	}

	// add mountpoint
	c.MountPoints["secrets"] = &volume.MountPoint{
		Source:      secretsPath,
		Destination: secretsContainerMountPath,
		RW:          false,
	}

	return nil
}
