// +build !windows

package swarm

type WindowsRawPrivilegeProfile struct{}

type Device struct {
	Path string
	Rwm  string
}

type LinuxRawPrivilegeProfile struct {
	Version      PrivilegeProfileVersion
	Capabilities []string
	Devices      []*Device
	AllDevices   bool
}
