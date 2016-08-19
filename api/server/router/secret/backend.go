package secret

import (
	enginetypes "github.com/docker/engine-api/types"
)

// Backend for Secret
type Backend interface {
	CreateSecret(enginetypes.Secret) error
	ListSecrets() ([]enginetypes.Secret, error)
	InspectSecret(id string) (*enginetypes.Secret, error)
	UpdateSecret(id string, s *enginetypes.Secret) error
	RemoveSecret(id string) error
}
