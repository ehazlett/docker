// +build linux

package secrets

const (
	secretsContainerMountPath = "/run/secrets"
)

func SecretsSupported() bool {
	return true
}
