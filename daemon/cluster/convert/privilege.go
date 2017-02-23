package convert

import (
	"fmt"

	swarmapi "github.com/docker/swarmkit/api"
)

func getSwarmCapability(c string) (*swarmapi.LinuxRawPrivilegeProfile_Capabilities, error) {
	v, ok := swarmapi.LinuxRawPrivilegeProfile_Capabilities_value[c]
	if !ok {
		return nil, fmt.Errorf("capability %s is not supported", c)
	}

	var cap swarmapi.LinuxRawPrivilegeProfile_Capabilities
	n := swarmapi.LinuxRawPrivilegeProfile_Capabilities_name[v]

	// TODO (ehazlett): i hate this. there has to be a better way....
	// can we add a lookup func in swarmkit.  pretty please.
	switch n {
	case "CAP_CHOWN":
		//return swarmapi.LinuxRawPrivilegeProfile_CAP_CHOWN, nil
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_CHOWN
	case "CAP_DAC_OVERRIDE":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_DAC_OVERRIDE
	case "CAP_DAC_READ_SEARCH":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_DAC_READ_SEARCH
	case "CAP_FOWNER":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_FOWNER
	case "CAP_FSETID":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_FSETID
	case "CAP_KILL":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_KILL
	case "CAP_SETGID":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SETGID
	case "CAP_SETUID":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SETUID
	case "CAP_SETPCAP":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SETPCAP
	case "CAP_LINUX_IMMUTABLE":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_LINUX_IMMUTABLE
	case "CAP_NET_BIND_SERVICE":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_NET_BIND_SERVICE
	case "CAP_NET_BROADCAST":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_NET_BROADCAST
	case "CAP_NET_ADMIN":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_NET_ADMIN
	case "CAP_NET_RAW":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_NET_RAW
	case "CAP_IPC_LOCK":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_IPC_LOCK
	case "CAP_IPC_OWNER":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_IPC_OWNER
	case "CAP_SYS_MODULE":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_MODULE
	case "CAP_SYS_RAWIO":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_RAWIO
	case "CAP_SYS_CHROOT":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_CHROOT
	case "CAP_SYS_PTRACE":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_PTRACE
	case "CAP_SYS_PACCT":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_PACCT
	case "CAP_SYS_ADMIN":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_ADMIN
	case "CAP_SYS_BOOT":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_BOOT
	case "CAP_SYS_NICE":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_NICE
	case "CAP_SYS_RESOURCE":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_RESOURCE
	case "CAP_SYS_TIME":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_TIME
	case "CAP_SYS_TTY_CONFIG":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYS_TTY_CONFIG
	case "CAP_MKNOD":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_MKNOD
	case "CAP_LEASE":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_LEASE
	case "CAP_AUDIT_WRITE":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_AUDIT_WRITE
	case "CAP_AUDIT_CONTROL":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_AUDIT_CONTROL
	case "CAP_SETFCAP":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SETFCAP
	case "CAP_MAC_OVERRIDE":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_MAC_OVERRIDE
	case "CAP_MAC_ADMIN":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_MAC_ADMIN
	case "CAP_SYSLOG":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_SYSLOG
	case "CAP_WAKE_ALARM":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_WAKE_ALARM
	case "CAP_BLOCK_SUSPEND":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_BLOCK_SUSPEND
	case "CAP_AUDIT_READ":
		cap = swarmapi.LinuxRawPrivilegeProfile_CAP_AUDIT_READ
	default:
		return nil, fmt.Errorf("capability %s is not supported", n)
	}

	return &cap, nil
}

func getSwarmCapabilities(c []string) ([]swarmapi.LinuxRawPrivilegeProfile_Capabilities, error) {
	caps := []swarmapi.LinuxRawPrivilegeProfile_Capabilities{}
	for _, x := range c {
		v, err := getSwarmCapability(x)
		if err != nil {
			return nil, err
		}
		if v != nil {
			caps = append(caps, *v)
		}
	}

	return caps, nil
}
