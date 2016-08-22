package builtin

import (
	"fmt"
	"time"

	enginetypes "github.com/docker/engine-api/types"
)

// BuiltinSecretStore is currently just an in memory store for debug
type BuiltinSecretStore struct {
	sharedKey string
	secrets   map[string]enginetypes.Secret
}

func NewSecretStore(sharedKey string) BuiltinSecretStore {
	return BuiltinSecretStore{
		sharedKey: sharedKey,
		secrets:   map[string]enginetypes.Secret{},
	}
}

func (s BuiltinSecretStore) CreateSecret(secret enginetypes.Secret) error {
	if secret.ID == "" {
		secret.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	s.secrets[secret.Name] = secret
	return nil
}

func (s BuiltinSecretStore) ListSecrets() ([]enginetypes.Secret, error) {
	var all []enginetypes.Secret
	for _, v := range s.secrets {
		all = append(all, v)
	}
	return all, nil
}
func (s BuiltinSecretStore) InspectSecret(id string) (*enginetypes.Secret, error) {
	if v, ok := s.secrets[id]; ok {
		return &v, nil
	}
	return nil, nil
}
func (s BuiltinSecretStore) UpdateSecret(id string, secret *enginetypes.Secret) error {
	s.secrets[id] = *secret
	return nil
}

func (s BuiltinSecretStore) RemoveSecret(id string) error {
	if _, ok := s.secrets[id]; ok {
		delete(s.secrets, id)
	}
	return nil
}
