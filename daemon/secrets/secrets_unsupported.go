// +build !linux

package secrets

const (
	secretsContainerMountPath = ""
)

func SecretsSupported() bool {
	return false
}
