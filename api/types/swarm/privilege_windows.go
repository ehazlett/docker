// +build windows

package swarm

type WindowsRawPrivilegeProfile struct {
	Version        PrivilegeProfileVersion
	CredentialSpec string
}
type LinuxRawPrivilegeProfile struct{}
