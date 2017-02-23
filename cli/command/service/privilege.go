package service

import "github.com/docker/docker/api/types/swarm"

// parsePrivilegeOptions parses privilege options passed in the form of
// key=val.  for example, cap=CAP_NET_ADMIN or credentialspec=/path/to/foo
func parsePrivilegeOptions(o map[string]string) (*swarm.PrivilegeProfiles, error) {
	p := &swarm.PrivilegeProfiles{
		Windows: &swarm.WindowsRawPrivilegeProfile{},
		Linux:   &swarm.LinuxRawPrivilegeProfile{},
	}

	for k, v := range o {
		switch k {
		case "cap":
			p.Linux.Capabilities = append(p.Linux.Capabilities, v)
		}
	}

	return p, nil
}
