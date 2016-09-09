package secret

type SecretMode string

const (
	SecretModeFile SecretMode = "file"
	SecretModeEnv  SecretMode = "env"
)

type Secret struct {
	ID         string `json:"Id,omitempty"`
	Name       string
	Data       []byte     `json:",omitempty"`
	Mountpoint string     `json:",omitempty"`
	Mode       SecretMode `json:",omitempty"`
}
