package daemon

import (
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/volume"
)

const (
	secretsVolumeName  = "secrets"
	secretsVolumeLabel = "com.docker.secrets.container"
)

func (daemon *Daemon) setupSecrets(c *container.Container) error {
	rdr, err := daemon.getSecrets(c)
	if err != nil {
		return err
	}

	m, ok := c.MountPoints[secretsContainerMountPath]
	if !ok {
		opts := map[string]string{
			"type": "ramfs",
		}
		labels := map[string]string{
			secretsVolumeLabel: c.ID,
		}
		vol, err := daemon.VolumeCreate("", "", opts, labels)
		if err != nil {
			return err
		}
		m = &volume.MountPoint{
			Name:        vol.Name,
			Source:      vol.Mountpoint,
			Destination: secretsContainerMountPath,
			RW:          false,
		}
		if c.Config.Labels == nil {
			c.Config.Labels = map[string]string{}
		}
		c.Config.Labels[secretsVolumeLabel] = vol.Name
	}
	c.MountPoints[secretsVolumeName] = m
	if c.Config.Volumes == nil {
		c.Config.Volumes = map[string]struct{}{}
	}
	c.Config.Volumes[m.Name] = struct{}{}

	// TODO: populate volume
	if err := archive.Untar(rdr, m.Source, nil); err != nil {
		return err
	}
	return nil
}
