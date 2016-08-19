package secrets

import enginetypes "github.com/docker/engine-api/types"

type SecretStore interface {
	CreateSecret(enginetypes.Secret) error
	ListSecrets() ([]enginetypes.Secret, error)
	InspectSecret(name string) (*enginetypes.Secret, error)
	UpdateSecret(name string, s *enginetypes.Secret) error
	RemoveSecret(name string) error
}
