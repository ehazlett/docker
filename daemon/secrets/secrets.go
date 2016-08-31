package secrets

import "github.com/docker/engine-api/types/secret"

type SecretStore interface {
	CreateSecret(secret.Secret) (*secret.Secret, error)
	ListSecrets() ([]secret.Secret, error)
	InspectSecret(name string) (*secret.Secret, error)
	UpdateSecret(name string, s *secret.Secret) error
	RemoveSecret(name string) error
}
