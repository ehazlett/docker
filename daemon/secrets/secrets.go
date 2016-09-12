package secrets

import (
	"path/filepath"

	"github.com/docker/docker/api/types/secret"
)

type SecretStore interface {
	Name() string
	CreateSecret(secret.Secret) (*secret.Secret, error)
	ListSecrets() ([]secret.Secret, error)
	InspectSecret(name string) (*secret.Secret, error)
	UpdateSecret(name string, s *secret.Secret) error
	RemoveSecret(name string) error
}

func SecretsContainerMountpath() string {
	return secretsContainerMountPath
}

func GetContainerMountpoint(target string) string {
	return filepath.Join(SecretsContainerMountpath(), target)
}
