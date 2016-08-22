package secret

import (
	"github.com/docker/engine-api/types"
)

// Backend for Secret
type Backend interface {
	CreateSecret(types.Secret) (*types.Secret, error)
	ListSecrets() ([]types.Secret, error)
	InspectSecret(id string) (*types.Secret, error)
	UpdateSecret(id string, s *types.Secret) error
	RemoveSecret(id string) error
}
