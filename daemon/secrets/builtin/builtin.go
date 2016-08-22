package builtin

import (
	"errors"
	"fmt"
	"time"

	"github.com/docker/engine-api/types"
)

var ErrSecretNotFound = errors.New("secret not found")

// BuiltinSecretStore is currently just an in memory store for debug
type BuiltinSecretStore struct {
	sharedKey string
	secrets   map[string]types.Secret
}

func NewSecretStore(sharedKey string) BuiltinSecretStore {
	return BuiltinSecretStore{
		sharedKey: sharedKey,
		secrets:   map[string]types.Secret{},
	}
}

func (s BuiltinSecretStore) CreateSecret(secret types.Secret) (*types.Secret, error) {
	if secret.ID == "" {
		secret.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	s.secrets[secret.ID] = secret
	return &secret, nil
}

func (s BuiltinSecretStore) ListSecrets() ([]types.Secret, error) {
	all := []types.Secret{}
	for _, v := range s.secrets {
		all = append(all, v)
	}
	return all, nil
}
func (s BuiltinSecretStore) InspectSecret(id string) (*types.Secret, error) {
	v, ok := s.secrets[id]
	if !ok {
		return nil, ErrSecretNotFound
	}
	return &v, nil
}
func (s BuiltinSecretStore) UpdateSecret(id string, secret *types.Secret) error {
	s.secrets[id] = *secret
	return nil
}

func (s BuiltinSecretStore) RemoveSecret(id string) error {
	if _, ok := s.secrets[id]; ok {
		delete(s.secrets, id)
	}
	return nil
}
