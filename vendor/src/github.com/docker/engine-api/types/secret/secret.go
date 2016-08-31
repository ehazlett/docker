package secret

type Secret struct {
	ID         string `json:"Id,omitempty"`
	Name       string
	Data       string
	Mountpoint string
	Required   bool
}
