package daemon

import (
	"github.com/docker/docker/container"
	"github.com/docker/libnetwork"
)

// LogContainerEvent generates an event related to a container.
func (daemon *Daemon) LogContainerEvent(container *container.Container, action string) {
	daemon.EventsService.Log(
		action,
		container.ID,
		container.Config.Image,
	)
}

// LogNetworkEvent generates an event related to a network.
func (daemon *Daemon) LogNetworkEvent(network libnetwork.Network, action string) {
	daemon.EventsService.Log(
		action,
		network.ID(),
		network.Type(),
	)
}
