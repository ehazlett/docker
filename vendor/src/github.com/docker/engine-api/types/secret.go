package types

type Secret struct {
	ID         string `json:"Id,omitempty"`
	Name       string
	Data       string
	Mountpoint string
}

// SecretsListResponse contains the response for the remote API:
// GET "/secrets"
type SecretsListResponse []*Secret

// SecretCreateRequest contains the response for the remote API:
// POST "/secrets/create"
type SecretCreateRequest struct {
	Name       string // Name is the requested name of the secret
	Data       string // Data is the data of the secret
	Mountpoint string // Mountpoint is the dest filename in the ramfs
}
