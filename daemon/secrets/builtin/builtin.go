package builtin

import (
	"errors"

	"github.com/docker/engine-api/types/secret"
)

var ErrSecretNotFound = errors.New("unable to find secret")

// BuiltinSecretStore is currently just an in memory store for debug
type BuiltinSecretStore struct {
	sharedKey string
	secrets   map[string]secret.Secret
}

func NewSecretStore(sharedKey string) BuiltinSecretStore {
	return BuiltinSecretStore{
		sharedKey: sharedKey,
		secrets:   map[string]secret.Secret{},
	}
}

func (s BuiltinSecretStore) CreateSecret(secret secret.Secret) (*secret.Secret, error) {
	if secret.ID == "" {
		secret.ID = secret.Name
	}
	s.secrets[secret.Name] = secret
	return nil
}

func (s BuiltinSecretStore) ListSecrets() ([]secret.Secret, error) {
	all := []secret.Secret{}
	for _, v := range s.secrets {
		all = append(all, v)
	}
	return all, nil
}
func (s BuiltinSecretStore) InspectSecret(id string) (*secret.Secret, error) {
	v, ok := s.secrets[id]
	if !ok {
		return nil, ErrSecretNotFound
	}
	return nil, nil
}
func (s BuiltinSecretStore) UpdateSecret(id string, secret *secret.Secret) error {
	s.secrets[id] = *secret
	return nil
}

func (s BuiltinSecretStore) RemoveSecret(name string) error {
	if _, ok := s.secrets[name]; ok {
		delete(s.secrets, name)
	}
	return nil
}
