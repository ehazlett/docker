// +build linux

package daemon

const (
	secretsContainerMountPath = "/run/secrets"
)

func secretsSupported() bool {
	return true
}
