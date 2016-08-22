package secrets

import "github.com/docker/engine-api/types"

type SecretStore interface {
	CreateSecret(types.Secret) (*types.Secret, error)
	ListSecrets() ([]types.Secret, error)
	InspectSecret(name string) (*types.Secret, error)
	UpdateSecret(name string, s *types.Secret) error
	RemoveSecret(name string) error
}
