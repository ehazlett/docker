package swarm

// Secret represents a secret.
type Secret struct {
	ID string
	Meta
	Name          string                 `json:",omitempty"`
	LatestVersion string                 `json:",omitempty"`
	SecretData    map[string]*SecretData `json:",omitempty"`
}

type SecretData struct {
	ID string
	Meta
	Spec       SecretSpec `json:",omitempty"`
	Digest     string     `json:",omitempty"`
	SecretSize int32      `json:",omitempty"`
}

type SecretType string

const (
	ContainerSecret SecretType = "ContainerSecret"
	NodeSecret      SecretType = "NodeSecret"
)

type SecretSpec struct {
	Annotations
	Type SecretType `json:",omitempty"`
	Data []byte     `json",omitempty"`
}
