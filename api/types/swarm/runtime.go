package swarm

type RuntimeType string

const (
	RuntimeContainer RuntimeType = "com.docker.runtime.container"
	RuntimePlugin    RuntimeType = "com.docker.runtime.plugin"
)
