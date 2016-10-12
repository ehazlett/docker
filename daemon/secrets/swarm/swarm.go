package swarm

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/secret"
	"github.com/docker/docker/daemon/cluster"
)

const (
	storeName = "swarm"
)

var ErrSecretNotFound = errors.New("unable to find secret")

// SwarmSecretStore is currently just an in memory store for debug
type SwarmSecretStore struct {
	c *cluster.Cluster
}

func NewSecretStore(c *cluster.Cluster) *SwarmSecretStore {
	logrus.Debugf("secrets: initializing swarm secret store %+v", c)
	return &SwarmSecretStore{
		c: c,
	}
}

func (s *SwarmSecretStore) Name() string {
	return storeName
}
func (s *SwarmSecretStore) InspectSecret(id string) (*secret.Secret, error) {
	sec, err := s.c.GetSecret(id)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("secrets: found secret %s in backend", sec.ID)

	return &secret.Secret{
		Name: sec.Spec.Annotations.Name,
		Data: sec.Spec.Data,
	}, nil
}
