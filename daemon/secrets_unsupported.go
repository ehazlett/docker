// +build !linux

package daemon

const (
	secretsContainerMountPath = ""
)

func secretsSupported() bool {
	return false
}
