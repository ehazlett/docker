package swarm

// Secret represents a secret.
type Secret struct {
	ID string
	Meta
	Spec       *SecretSpec `json:",omitempty"`
	Digest     string      `json:",omitempty"`
	SecretSize int32       `json:",omitempty"`
}

type SecretSpec struct {
	Annotations
	Data []byte `json",omitempty"`
}

type SecretReferenceMode int

const (
	SecretReferenceFile SecretReferenceMode = 0
	SecretReferenceEnv  SecretReferenceMode = 1
)

type SecretReference struct {
	SecretID   string              `json:",omitempty"`
	Mode       SecretReferenceMode `json:",omitempty"`
	Name       string              `json:",omitempty"`
	Target     string              `json:",omitempty"`
	SecretName string              `json:",omitempty"`
}
