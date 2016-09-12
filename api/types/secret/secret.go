package secret

type SecretMode int

const (
	ModeEnv  SecretMode = 0
	ModeFile SecretMode = 1
)

type Secret struct {
	ID         string `json:"Id,omitempty"`
	Name       string
	Data       []byte     `json:",omitempty"`
	Mode       SecretMode `json:",omitempty"`
	Mountpoint string     `json:",omitempty"`
}
