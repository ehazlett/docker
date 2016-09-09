package daemon

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/cloudflare/cfssl/log"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/secrets"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/volume"
	"github.com/docker/engine-api/types/secret"
)

const (
	secretsVolumeName  = "secrets"
	secretsVolumeLabel = "com.docker.secrets.container"
)

var (
	ErrSecretStoreNotInitialized = errors.New("secret store is not initialized")
)

func secretsMountPath(c *container.Container) string {
	return filepath.Join(c.Root, "secrets")
}

func (daemon *Daemon) SetSecretStore(s secrets.SecretStore) {
	daemon.secretStore = s
}

func (daemon *Daemon) getSecrets(container *container.Container) (io.Reader, error) {
	logrus.Debugf("mounting secrets for container %s", container.ID)
	logrus.Debugf("container %s requested secrets %v", container.ID, container.Config.Secrets)

	logrus.Debugf("secret store: %+v", daemon.secretStore)

	// generate and return tar of secrets
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	for _, s := range container.Config.Secrets {
		logrus.Debugf("requesting secret %q for container %s", s.Name, container.ID)
		secret, err := daemon.InspectSecret(s.Name)
		if err != nil {
			logrus.Warnf("secret: unable to find secret %s in backend: %s", s.Name, err)
			continue
		}

		logrus.Debugf("received secret %s for container %s", secret.Name, container.ID)

		h := &tar.Header{
			Name: s.Mountpoint,
			Mode: 0600,
			Size: int64(len(secret.Data)),
		}
		if err := tw.WriteHeader(h); err != nil {
			return nil, err

		}
		if _, err := tw.Write(secret.Data); err != nil {
			return nil, err
		}
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func (daemon *Daemon) setupSecrets(c *container.Container) error {
	rdr, err := daemon.getSecrets(c)
	if err != nil {
		return err
	}

	if rdr == nil {
		return nil
	}

	secretsMountPath := secrets.SecretsContainerMountpath()
	m, ok := c.MountPoints[secretsMountPath]
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
			Destination: secretsMountPath,
			RW:          false,
		}
		c.MountPoints[secretsMountPath] = m

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

func (daemon *Daemon) InspectSecret(id string) (*secret.Secret, error) {
	if daemon.secretStore == nil {
		return nil, fmt.Errorf("secret store is not initialized")
	}

	return daemon.secretStore.InspectSecret(id)
}
