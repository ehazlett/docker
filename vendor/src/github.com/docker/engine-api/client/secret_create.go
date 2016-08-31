package client

import (
	"encoding/json"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/secret"
	"golang.org/x/net/context"
)

// SecretCreate creates a secret in the docker host.
func (cli *Client) SecretCreate(ctx context.Context, options types.SecretCreateRequest) (secret.Secret, error) {
	var secret secret.Secret
	resp, err := cli.post(ctx, "/secrets/create", nil, options, nil)
	if err != nil {
		return secret, err
	}
	err = json.NewDecoder(resp.body).Decode(&secret)
	ensureReaderClosed(resp)
	return secret, err
}
