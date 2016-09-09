package swarm

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types/secret"
	swarmapi "github.com/docker/swarmkit/api"
)

const (
	storeName = "swarm"
)

var ErrSecretNotFound = errors.New("unable to find secret")

// SwarmSecretStore is currently just an in memory store for debug
type SwarmSecretStore struct {
	client swarmapi.SecretsClient
}

func NewSecretStore(c swarmapi.SecretsClient) *SwarmSecretStore {
	logrus.Debugf("secrets: initializing swarm secret store %+v", c)
	return &SwarmSecretStore{
		client: c,
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
	logrus.Debugf("secret store: looking up secret %s in swarm", id)
	r, err := s.client.GetSecret(context.Background(), &swarmapi.GetSecretRequest{Name: id})
	if err != nil {
		return nil, err
	}

	logrus.Debugf("secrets: found secret in backend: %+v", r)

	sec := r.Secret
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
