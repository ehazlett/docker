package daemon

import (
	"fmt"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/cloudflare/cfssl/log"
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/volume"
)

const (
	secretsVolumeName  = "secrets"
	secretsVolumeLabel = "com.docker.secrets.container"
)

func secretsMountPath(c *container.Container) string {
	return filepath.Join(c.Root, "secrets")
}

func (daemon *Daemon) setupSecrets(c *container.Container) error {
	rdr, err := daemon.getSecrets(c)
	if err != nil {
		return err
	}

	m, ok := c.MountPoints[secretsContainerMountPath]
	if !ok {
		opts := map[string]string{
			"type":   "tmpfs",
			"device": "tmpfs",
		}
		labels := map[string]string{
			secretsVolumeLabel: c.ID,
		}
		vol, err := daemon.VolumeCreate("", "local", opts, labels)
		if err != nil {
			return err
		}

		logrus.Debugf("secrets: setting up secrets at %s", vol.Mountpoint)
		if err := mount.Mount("tmpfs", vol.Mountpoint, "tmpfs", "nodev"); err != nil {
			return fmt.Errorf("secrets: unable to setup mount: %s", err)
		}

		m = &volume.MountPoint{
			Name:        vol.Name,
			Source:      vol.Mountpoint,
			Destination: secretsContainerMountPath,
			RW:          false,
		}
		c.MountPoints[secretsContainerMountPath] = m

		if c.Config.Labels == nil {
			c.Config.Labels = map[string]string{}
		}
		c.Config.Labels[secretsVolumeLabel] = vol.Name
	}

	log.Debugf("secrets: mountpoint %+v", m)
	if c.Config.Volumes == nil {
		c.Config.Volumes = map[string]struct{}{}
	}
	c.Config.Volumes[m.Name] = struct{}{}

	// populate volume
	logrus.Debugf("secrets: populating volume data for %s -> %s", c.ID, m.Name)
	if err := archive.Untar(rdr, m.Source, nil); err != nil {
		return err
	}
	return nil
}
