package swarm

import (
	"errors"
	"fmt"

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
	sec, err := s.c.GetSecretFromNode(id)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("secrets: found secret %s in backend", sec.Name)

	rev, ok := sec.SecretData[sec.LatestVersion]
	if !ok {
		return nil, fmt.Errorf("unable to find latest revision of secret")
	}

	return &secret.Secret{
		Name: sec.Name,
		Data: rev.Spec.Data,
	}, nil
}
