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
func (s *SwarmSecretStore) CreateSecret(secret secret.Secret) (*secret.Secret, error) {
	//if secret.ID == "" {
	//	secret.ID = secret.Name
	//}
	//s.secrets[secret.ID] = secret
	//return &secret, nil
	return nil, nil
}

func (s *SwarmSecretStore) ListSecrets() ([]secret.Secret, error) {
	all := []secret.Secret{}
	//for _, v := range s.secrets {
	//	all = append(all, v)
	//}
	return all, nil
}
func (s *SwarmSecretStore) InspectSecret(id string) (*secret.Secret, error) {
	sec, err := s.c.GetSecret(id)
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

func (s *SwarmSecretStore) UpdateSecret(id string, secret *secret.Secret) error {
	return nil
}

func (s *SwarmSecretStore) RemoveSecret(id string) error {
	//if _, ok := s.secrets[id]; ok {
	//	delete(s.secrets, id)
	//}
	return nil
}
