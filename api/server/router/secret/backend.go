package secret

import (
	"github.com/docker/engine-api/types/secret"
)

// Backend for Secret
type Backend interface {
	CreateSecret(secret.Secret) (*secret.Secret, error)
	ListSecrets() ([]secret.Secret, error)
	InspectSecret(id string) (*secret.Secret, error)
	UpdateSecret(id string, s *secret.Secret) error
	RemoveSecret(id string) error
}
