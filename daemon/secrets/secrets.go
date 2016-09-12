package secrets

import (
	"path/filepath"

	"github.com/docker/docker/api/types/secret"
)

type SecretStore interface {
	Name() string
	InspectSecret(name string) (*secret.Secret, error)
}

func SecretsContainerMountpath() string {
	return secretsContainerMountPath
}

func GetContainerMountpoint(target string) string {
	return filepath.Join(SecretsContainerMountpath(), target)
}
