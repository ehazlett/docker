package convert

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	container "github.com/docker/docker/api/types/container"
	mounttypes "github.com/docker/docker/api/types/mount"
	types "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
)

func containerSpecFromGRPC(c *swarmapi.ContainerSpec) types.ContainerSpec {
	containerSpec := types.ContainerSpec{
		Image:     c.Image,
		Labels:    c.Labels,
		Command:   c.Command,
		Args:      c.Args,
		Hostname:  c.Hostname,
		Env:       c.Env,
		Dir:       c.Dir,
		User:      c.User,
		Groups:    c.Groups,
		TTY:       c.TTY,
		OpenStdin: c.OpenStdin,
		ReadOnly:  c.ReadOnly,
		Hosts:     c.Hosts,
		Secrets:   secretReferencesFromGRPC(c.Secrets),
	}

	if c.DNSConfig != nil {
		containerSpec.DNSConfig = &types.DNSConfig{
			Nameservers: c.DNSConfig.Nameservers,
			Search:      c.DNSConfig.Search,
			Options:     c.DNSConfig.Options,
		}
	}

	// Mounts
	for _, m := range c.Mounts {
		mount := mounttypes.Mount{
			Target:   m.Target,
			Source:   m.Source,
			Type:     mounttypes.Type(strings.ToLower(swarmapi.Mount_MountType_name[int32(m.Type)])),
			ReadOnly: m.ReadOnly,
		}

		if m.BindOptions != nil {
			mount.BindOptions = &mounttypes.BindOptions{
				Propagation: mounttypes.Propagation(strings.ToLower(swarmapi.Mount_BindOptions_MountPropagation_name[int32(m.BindOptions.Propagation)])),
			}
		}

		if m.VolumeOptions != nil {
			mount.VolumeOptions = &mounttypes.VolumeOptions{
				NoCopy: m.VolumeOptions.NoCopy,
				Labels: m.VolumeOptions.Labels,
			}
			if m.VolumeOptions.DriverConfig != nil {
				mount.VolumeOptions.DriverConfig = &mounttypes.Driver{
					Name:    m.VolumeOptions.DriverConfig.Name,
					Options: m.VolumeOptions.DriverConfig.Options,
				}
			}
		}

		if m.TmpfsOptions != nil {
			mount.TmpfsOptions = &mounttypes.TmpfsOptions{
				SizeBytes: m.TmpfsOptions.SizeBytes,
				Mode:      m.TmpfsOptions.Mode,
			}
		}
		containerSpec.Mounts = append(containerSpec.Mounts, mount)
	}

	if c.StopGracePeriod != nil {
		grace, _ := gogotypes.DurationFromProto(c.StopGracePeriod)
		containerSpec.StopGracePeriod = &grace
	}

	if c.Healthcheck != nil {
		containerSpec.Healthcheck = healthConfigFromGRPC(c.Healthcheck)
	}

	profiles := &types.PrivilegeProfiles{
		Windows: &types.WindowsRawPrivilegeProfile{},
		Linux:   &types.LinuxRawPrivilegeProfile{},
	}

	// linux profile
	if c.SecurityProfiles != nil && c.SecurityProfiles.Linux != nil {
		p := c.SecurityProfiles.Linux
		var ver types.PrivilegeProfileVersion
		switch p.Version {
		case swarmapi.LinuxRawPrivilegeProfile_V0:
			ver = types.PrivilegeProfileVersionLinuxV0

		}
		lp := &types.LinuxRawPrivilegeProfile{
			Version:    ver,
			AllDevices: p.AllDevices,
		}

		devices := []*types.Device{}
		for _, d := range p.Devices {
			devices = append(devices, &types.Device{
				Path: d.Path,
				Rwm:  d.Rwm,
			})
		}
		lp.Devices = devices
		pCaps := []string{}
		for _, c := range p.Capabilities {
			pCaps = append(pCaps, c.String())
		}
		lp.Capabilities = pCaps

		profiles.Linux = lp
		logrus.Debugf("LINUX PROFILE: %+v", lp)
	}

	// TODO (ehazlett): windows privilege profile

	containerSpec.SecurityProfiles = profiles

	return containerSpec
}

func secretReferencesToGRPC(sr []*types.SecretReference) []*swarmapi.SecretReference {
	refs := make([]*swarmapi.SecretReference, 0, len(sr))
	for _, s := range sr {
		ref := &swarmapi.SecretReference{
			SecretID:   s.SecretID,
			SecretName: s.SecretName,
		}
		if s.File != nil {
			ref.Target = &swarmapi.SecretReference_File{
				File: &swarmapi.SecretReference_FileTarget{
					Name: s.File.Name,
					UID:  s.File.UID,
					GID:  s.File.GID,
					Mode: s.File.Mode,
				},
			}
		}

		refs = append(refs, ref)
	}

	return refs
}
func secretReferencesFromGRPC(sr []*swarmapi.SecretReference) []*types.SecretReference {
	refs := make([]*types.SecretReference, 0, len(sr))
	for _, s := range sr {
		target := s.GetFile()
		if target == nil {
			// not a file target
			logrus.Warnf("secret target not a file: secret=%s", s.SecretID)
			continue
		}
		refs = append(refs, &types.SecretReference{
			File: &types.SecretReferenceFileTarget{
				Name: target.Name,
				UID:  target.UID,
				GID:  target.GID,
				Mode: target.Mode,
			},
			SecretID:   s.SecretID,
			SecretName: s.SecretName,
		})
	}

	return refs
}

func containerToGRPC(c types.ContainerSpec) (*swarmapi.ContainerSpec, error) {
	containerSpec := &swarmapi.ContainerSpec{
		Image:     c.Image,
		Labels:    c.Labels,
		Command:   c.Command,
		Args:      c.Args,
		Hostname:  c.Hostname,
		Env:       c.Env,
		Dir:       c.Dir,
		User:      c.User,
		Groups:    c.Groups,
		TTY:       c.TTY,
		OpenStdin: c.OpenStdin,
		ReadOnly:  c.ReadOnly,
		Hosts:     c.Hosts,
		Secrets:   secretReferencesToGRPC(c.Secrets),
	}

	if c.DNSConfig != nil {
		containerSpec.DNSConfig = &swarmapi.ContainerSpec_DNSConfig{
			Nameservers: c.DNSConfig.Nameservers,
			Search:      c.DNSConfig.Search,
			Options:     c.DNSConfig.Options,
		}
	}

	if c.StopGracePeriod != nil {
		containerSpec.StopGracePeriod = gogotypes.DurationProto(*c.StopGracePeriod)
	}

	// Mounts
	for _, m := range c.Mounts {
		mount := swarmapi.Mount{
			Target:   m.Target,
			Source:   m.Source,
			ReadOnly: m.ReadOnly,
		}

		if mountType, ok := swarmapi.Mount_MountType_value[strings.ToUpper(string(m.Type))]; ok {
			mount.Type = swarmapi.Mount_MountType(mountType)
		} else if string(m.Type) != "" {
			return nil, fmt.Errorf("invalid MountType: %q", m.Type)
		}

		if m.BindOptions != nil {
			if mountPropagation, ok := swarmapi.Mount_BindOptions_MountPropagation_value[strings.ToUpper(string(m.BindOptions.Propagation))]; ok {
				mount.BindOptions = &swarmapi.Mount_BindOptions{Propagation: swarmapi.Mount_BindOptions_MountPropagation(mountPropagation)}
			} else if string(m.BindOptions.Propagation) != "" {
				return nil, fmt.Errorf("invalid MountPropagation: %q", m.BindOptions.Propagation)
			}
		}

		if m.VolumeOptions != nil {
			mount.VolumeOptions = &swarmapi.Mount_VolumeOptions{
				NoCopy: m.VolumeOptions.NoCopy,
				Labels: m.VolumeOptions.Labels,
			}
			if m.VolumeOptions.DriverConfig != nil {
				mount.VolumeOptions.DriverConfig = &swarmapi.Driver{
					Name:    m.VolumeOptions.DriverConfig.Name,
					Options: m.VolumeOptions.DriverConfig.Options,
				}
			}
		}

		if m.TmpfsOptions != nil {
			mount.TmpfsOptions = &swarmapi.Mount_TmpfsOptions{
				SizeBytes: m.TmpfsOptions.SizeBytes,
				Mode:      m.TmpfsOptions.Mode,
			}
		}

		containerSpec.Mounts = append(containerSpec.Mounts, mount)
	}

	if c.Healthcheck != nil {
		containerSpec.Healthcheck = healthConfigToGRPC(c.Healthcheck)
	}

	profiles := &swarmapi.PrivilegeProfiles{
		Windows: &swarmapi.WindowsRawPrivilegeProfile{},
		Linux:   &swarmapi.LinuxRawPrivilegeProfile{},
	}
	if p := c.SecurityProfiles.Linux; p != nil {
		lp := &swarmapi.LinuxRawPrivilegeProfile{
			AllDevices: p.AllDevices,
		}
		switch p.Version {
		case types.PrivilegeProfileVersionLinuxV0:
			lp.Version = swarmapi.LinuxRawPrivilegeProfile_V0
		}

		devices := []*swarmapi.LinuxRawPrivilegeProfile_Device{}
		for _, d := range p.Devices {
			devices = append(devices, &swarmapi.LinuxRawPrivilegeProfile_Device{
				Path: d.Path,
				Rwm:  d.Rwm,
			})
		}
		pCaps, err := getSwarmCapabilities(p.Capabilities)
		if err != nil {
			return nil, err
		}
		lp.Capabilities = pCaps

		profiles.Linux = lp
	}
	// TODO (ehazlett): windows privilege profile

	containerSpec.SecurityProfiles = profiles

	return containerSpec, nil
}

func healthConfigFromGRPC(h *swarmapi.HealthConfig) *container.HealthConfig {
	interval, _ := gogotypes.DurationFromProto(h.Interval)
	timeout, _ := gogotypes.DurationFromProto(h.Timeout)
	return &container.HealthConfig{
		Test:     h.Test,
		Interval: interval,
		Timeout:  timeout,
		Retries:  int(h.Retries),
	}
}

func healthConfigToGRPC(h *container.HealthConfig) *swarmapi.HealthConfig {
	return &swarmapi.HealthConfig{
		Test:     h.Test,
		Interval: gogotypes.DurationProto(h.Interval),
		Timeout:  gogotypes.DurationProto(h.Timeout),
		Retries:  int32(h.Retries),
	}
}
