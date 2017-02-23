package swarm

type PrivilegeProfileVersion int

const (
	PrivilegeProfileVersionLinuxV0   PrivilegeProfileVersion = 0
	PrivilegeProfileVersionWindowsV0 PrivilegeProfileVersion = 0
)

type PrivilegeProfiles struct {
	Windows *WindowsRawPrivilegeProfile `json:",omitempty"`
	Linux   *LinuxRawPrivilegeProfile   `json:",omitempty"`
}
