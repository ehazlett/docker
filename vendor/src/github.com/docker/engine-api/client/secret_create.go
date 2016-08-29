package client

import (
	"encoding/json"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// SecretCreate creates a secret in the docker host.
func (cli *Client) SecretCreate(ctx context.Context, options types.SecretCreateRequest) (types.Secret, error) {
	var secret types.Secret
	resp, err := cli.post(ctx, "/secrets/create", nil, options, nil)
	if err != nil {
		return secret, err
	}
	err = json.NewDecoder(resp.body).Decode(&secret)
	ensureReaderClosed(resp)
	return secret, err
}
